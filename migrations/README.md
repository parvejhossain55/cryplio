# Cryplio Database Migrations

This directory contains the database migration files for the Cryplio P2P Exchange platform.

## Migration Structure

Migrations are intentionally **small and focused** — each file represents a single logical change. This makes code review, rollback, and debugging easier.

### Migration Files

| # | Name | Purpose |
|---|------|---------|
| `000_types.sql` | ENUM Type Definitions | All custom PostgreSQL ENUM types used across tables |
| `001_core_tables.sql` | Core User & Auth Tables | `users`, `user_stats`, `user_sessions`, `password_reset_tokens`, `email_verification_tokens`, `user_blocks` |
| `002_lookup_tables.sql` | Lookup Tables | `crypto_assets`, `fiat_currencies`, `payment_methods`, `dispute_reasons`, `fee_tiers` + seed data |
| `004_trade_ads.sql` | Trade Ads | `trade_ads` (P2P advertisements) |
| `005_trades.sql` | Trades | `trades`, `trade_messages`, `trade_attachments` |
| `006_feedback.sql` | Feedback | `trade_feedback` (post-trade ratings) |
| `007_disputes.sql` | Disputes | `disputes`, `dispute_messages` |
| `008_wallets.sql` | Wallets & Transactions | `wallets`, `wallet_transactions`, `withdrawals` |
| `009_notifications.sql` | Notifications | `notifications`, `notification_preferences` |
| `010_referrals.sql` | Referrals | `referrals` |
| `011_merchant.sql` | Merchant System | `merchant_applications`, `merchant_analytics` |
| `012_audit_admin.sql` | Audit & Admin | `audit_logs`, `admin_actions`, `platform_config`, `announcements` |
| `013_monitoring.sql` | Monitoring | `rate_limit_counts`, `login_attempts`, `api_request_logs` |
| `014_views.sql` | Views | `active_trade_ads`, `completed_trades`, `public_user_profiles`, `open_disputes` |
| `015_triggers.sql` | Triggers | Timestamp auto-update triggers for all tables |
| `016_seed_data.sql` | Seed Data | Initial data for lookup tables (idempotent) |

**Total:** 17 migration files (000–016), 23 tables, 4 views, 60+ indexes.

## Running Migrations

### Option 1: Using psql (recommended for production)

```bash
# Run all migrations in order
cd /home/parvej/Project/Cryplio/migrations
psql -U <user> -d <database> -f 000_types.sql
psql -U <user> -d <database> -f 001_core_tables.sql
# ... continue in order
```

Or run all at once:

```bash
for f in $(ls -v *.sql); do
  echo "Applying $f..."
  psql -U <user> -d <database> -f "$f"
done
```

### Option 2: Using migrate CLI (recommended for CI/CD)

Install: `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`

```bash
# Create migration source URL
export DATABASE_URL="postgres://user:pass@localhost:5432/cryplio_db?sslmode=disable"

# Apply all migrations
migrate -path ./migrations -database "$DATABASE_URL" up

# Check status
migrate -path ./migrations -database "$DATABASE_URL" version

# Rollback last batch
migrate -path ./migrations -database "$DATABASE_URL" down 1
```

### Option 3: Using Golang program (embedded in application)

The application will embed migrations using `github.com/golang-migrate/migrate/v4` and apply them on startup in development mode. Production should use explicit migration CLI.

## Schema Design Principles

### 1. UUID Primary Keys
- All primary keys are UUIDs (GenRandomUUID) for distributed system safety
- Never expose sequential IDs publicly
- Trade IDs use UUID for obfuscation

### 2. Soft Deletes
- Tables: `users`, `trade_ads`, `trades`, `messages` have `deleted_at`
- GDPR compliance: user data retained 5 years, then purged
- Never physically delete financial records

### 3. Denormalized Statistics
- `user_stats` table aggregates trade counts, ratings, volume
- Updated by triggers or application logic
- Improves profile page performance

### 4. Index Strategy
- Foreign keys indexed automatically
- Composite indexes for common query patterns
- Partial indexes for filtered queries (e.g., only active records)
- GIN indexes for array columns (`payment_methods`)

### 5. JSONB for Flexible Data
- `evidence_links` in disputes
- `metadata` in transactions
- `data` in notifications
- Indexed when needed with GIN

### 6. Timestamp Tracking
- `created_at` — row creation
- `updated_at` — auto-updated by trigger
- Specific timestamps: `published_at`, `expires_at`, `completed_at`, `resolved_at`, etc.

### 7. Compliance & Audit
- `audit_logs` — generic change history for all tables
- `admin_actions` — specific admin operations with IP/UA
- `platform_config` — auditable configuration changes
- All foreign keys have appropriate ON DELETE actions (RESTRICT for financial, CASCADE for dependent)

## Relationship Map

```
users
  ├─ user_stats (1:1)
  ├─ user_sessions (1:N)
  ├─ trade_ads (1:N)
  ├─ wallets (1:N per crypto)
  │   └─ wallet_transactions (1:N)
  ├─ withdrawals (1:N)
  ├─ referrals (1:N as referrer)
  ├─ referrals (1:N as referee)
  ├─ user_blocks (1:N)
  ├─ notifications (1:N)
  ├─ merchant_applications (1:1)
  └─ audit_logs (N:1)
         └─ admin_actions (N:1)

trade_ads
  └─ trades (1:N)
         ├─ trade_messages (1:N)
         ├─ trade_attachments (N:1 via message)
         ├─ trade_feedback (1:1)
         └─ disputes (1:1)
                └─ dispute_messages (1:N)
```

## Data Retention Policy (SRS 7.2)

| Data | Retention | Action |
|------|-----------|---------|
| User Account Data | 5 years post-closure | Soft delete → purge job |
| Trade Records | 7 years | Never delete (archive to cold storage after 2 years) |
| Chat Messages | 2 years | Auto-purge job |
| Audit Logs | Permanent | Append-only |
| Session Tokens | 24h inactive | Auto-expire |
| Rate Limit Counts | 2 minutes | TTL index |

## Migration Best Practices

1. **Each migration is idempotent** — can be re-run safely (uses `IF NOT EXISTS`, `ON CONFLICT DO NOTHING`)
2. **Never remove columns** — only add new ones; deprecate old ones by renaming
3. **Always add indexes in same transaction** as table creation to avoid lock contention
4. **Small batches** — max 10-15 tables per migration file
5. **Document trade-offs** in comments (e.g., denormalization choices)

## Adding a New Migration

```bash
# 1. Create file with next sequence number
nano 017_add_new_table.sql

# 2. Include proper BEGIN/COMMIT, indexes, and seed data if needed
# 3. Document the change in file header
# 4. Run locally to verify
psql -U postgres -d cryplio -f 017_add_new_table.sql

# 5. Commit with message
git add 017_add_new_table.sql
git commit -m "Add new_table for X feature (SRS FR-XXX)"
```

## Rollback Procedure

To revert a migration:

```bash
# Find the migration you want to rollback
migrate -path ./migrations -database "$DATABASE_URL" version
# => 16 (current)

# Roll back 1 migration
migrate -path ./migrations -database "$DATABASE_URL" down 1

# Or manually drop objects (dangerous, ensure no data loss)
psql -U postgres -d cryplio_db -c "DROP VIEW IF EXISTS problematic_view;"
```

**Warning:** Never rollback a migration that has already been deployed to production without thorough review. Instead, create a new migration that undoes the change.

## Testing Migrations

```bash
# Spin up a fresh Postgres instance
docker run -d --name cryplio-db -e POSTGRES_PASSWORD=secret -p 5432:5432 postgres:15-alpine

# Wait for DB to be ready
sleep 5

# Apply all migrations
export DATABASE_URL="postgres://postgres:secret@localhost:5432/postgres?sslmode=disable"
migrate -path ./migrations -database "$DATABASE_URL" up

# Verify
psql -U postgres -d postgres -c "\dt"  # Should list all tables
```

## Production Checklist

- [ ] Review migration SQL for security (SQL injection via dynamic SQL is impossible here)
- [ ] Ensure all indexes are present (trade-off: space vs query speed)
- [ ] Verify foreign key ON DELETE actions (CASCADE vs RESTRICT)
- [ ] Check column nullability — `NOT NULL` with defaults preferred
- [ ] Confirm ENUM values match application constants
- [ ] Run `ANALYZE` after migration on large tables
- [ ] Monitor query plans after deployment (may need additional indexes)

---

**Total Tables:** 23  
**Total Indexes:** 60+  
**Total Views:** 4
