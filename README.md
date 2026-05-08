# Cryplio

> **Trade Crypto. Trust the Process.**

P2P cryptocurrency exchange platform — built with Clean Architecture and Domain-Driven Design in Go.

---

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Environment Variables](#environment-variables)
- [Running the App](#running-the-app)
- [Database Migrations](#database-migrations)
- [API Reference](#api-reference)
- [WebSocket](#websocket)
- [Domain Guide](#domain-guide)
- [Smart Contracts](#smart-contracts)
- [Contributing](#contributing)

---

## Overview

Cryplio is a global P2P crypto exchange where buyers and sellers trade directly with each other using a blockchain escrow system. The platform holds no custody of funds — every trade is secured by a smart contract escrow deployed on EVM-compatible chains.

**Key capabilities**

- Escrow-backed P2P trades (USDT, USDC, ETH, BTC)
- Multi-payment method support — bKash, Nagad, Bank Transfer, Wise, PayPal, UPI
- Real-time trade chat and notifications via WebSocket
- Google OAuth 2.0 login alongside email/password
- TOTP two-factor authentication
- Admin dispute resolution with evidence upload
- Merchant portal for high-volume traders
- Referral programme with commission tracking
- Rate-limited, session-based API with HttpOnly cookie + Bearer token support

---

## Architecture

The codebase follows **Clean Architecture** with **Domain-Driven Design**. Dependencies flow strictly inward — outer layers may depend on inner layers, never the reverse.

```
┌────────────────────────────────────────────────────────┐
│  interfaces/        HTTP handlers · WebSocket          │
│  ┌──────────────────────────────────────────────────┐  │
│  │  application/   Use cases · orchestration        │  │
│  │  ┌────────────────────────────────────────────┐  │  │
│  │  │  domain/    Entities · rules · interfaces  │  │  │
│  │  └────────────────────────────────────────────┘  │  │
│  └──────────────────────────────────────────────────┘  │
│  infrastructure/    Postgres · Redis · blockchain      │
└────────────────────────────────────────────────────────┘
```

| Layer | Package | Responsibility |
|---|---|---|
| **Domain** | `internal/domain/` | Pure business logic. No external imports. Defines entities, domain services, and repository interfaces. |
| **Application** | `internal/application/` | Use-case orchestration. Delegates to domain services. |
| **Interfaces** | `internal/interfaces/` | HTTP handlers grouped by domain, DTOs, middleware, WebSocket hub. |
| **Infrastructure** | `internal/infrastructure/` | Postgres repositories, Redis, blockchain clients, email, object storage. |
| **Shared** | `pkg/` | JWT, crypto helpers, structured logger, config, pagination. No domain imports. |

---

## Tech Stack

| Concern | Technology |
|---|---|
| Language | Go 1.25 |
| HTTP Framework | Gin |
| Database | PostgreSQL 17 |
| Cache / Rate Limiting | Redis 7 (go-redis/v9) |
| Job Queue | Asynq |
| Migrations | golang-migrate |
| Auth | JWT (golang-jwt/v5) + TOTP (pquerna/otp) |
| OAuth | Google OAuth 2.0 |
| WebSocket | gorilla/websocket |
| Object Storage | MinIO (S3-compatible) |
| Blockchain | go-ethereum |
| Smart Contracts | Solidity + Foundry |
| Email | SMTP (configurable provider) |
| Logging | Zerolog |
| Hot Reload | Air |
| Containerisation | Docker Compose (local dev) |

---

## Project Structure

```
cryplio/
├── cmd/
│   ├── api/            # HTTP server entrypoint
│   ├── migrate/        # Migration CLI (apply / rollback / status)
│   └── seed/           # Database seeder
│
├── internal/
│   │
│   ├── domain/                     # Layer 1 — pure business logic
│   │   ├── identity/               # User, Auth, Sessions, OAuth, 2FA, Payment methods
│   │   │   ├── service.go          # AuthService interface + constructor
│   │   │   ├── auth.go             # Register, Login, Logout, Refresh, Complete2FA
│   │   │   ├── oauth.go            # Google OAuth flow
│   │   │   ├── profile.go          # Profile CRUD + stats
│   │   │   ├── email.go            # Email verification + password reset
│   │   │   ├── twofactor.go        # TOTP setup / verify / disable
│   │   │   ├── session.go          # Session CRUD
│   │   │   ├── payment.go          # User payment method profiles
│   │   │   ├── admin.go            # Admin user management
│   │   │   ├── user.go             # User entity + value objects
│   │   │   ├── repository.go       # Segregated repository interfaces
│   │   │   └── helpers.go          # hashToken, validatePasswordComplexity
│   │   │
│   │   ├── trading/                # Trade advertisements + Trade lifecycle
│   │   │   ├── service.go          # TradeService interface + constructor
│   │   │   ├── ad.go               # Ad management (create, update, list, toggle)
│   │   │   ├── trade.go            # Trade lifecycle (initiate → pay → release → dispute)
│   │   │   ├── chat.go             # Trade chat messages + feedback
│   │   │   ├── models.go           # TradeAd, Trade, TradeMessage entities
│   │   │   ├── feedback.go         # TradeFeedback entity
│   │   │   ├── blockchain.go       # EscrowContractClient interface
│   │   │   └── repository.go       # TradeRepository interface
│   │   │
│   │   ├── wallet/                 # Wallets + transactions
│   │   ├── dispute/                # Dispute lifecycle + evidence
│   │   ├── notification/           # In-app notifications + preferences
│   │   ├── platform/               # CryptoAsset, FiatCurrency, PaymentMethod catalogue
│   │   ├── market/                 # Exchange rate feed
│   │   ├── referral/               # Referral tracking + payouts
│   │   └── events/                 # Domain event types + dispatcher interface
│   │
│   ├── application/                # Layer 2 — use cases
│   │   ├── bootstrap.go            # App wiring (DI root)
│   │   ├── user/                   # login, register, session use cases
│   │   ├── trade/                  # initiate_trade, create_ad, mark_paid, release_escrow, cancel_trade
│   │   ├── dispute/                # raise, assign, resolve
│   │   ├── wallet/                 # deposit, withdraw, balance
│   │   ├── merchant/               # apply, verify (scaffolded)
│   │   ├── referral/               # track, payout
│   │   └── identity/               # (reserved)
│   │
│   ├── interfaces/                 # Layer 3 — delivery
│   │   ├── http/
│   │   │   ├── handler/            # Handlers grouped by domain
│   │   │   │   ├── auth/           # AuthHandler (login, register, OAuth, 2FA, sessions, payment methods)
│   │   │   │   │                   # AdminHandler (user management, dashboard stats)
│   │   │   │   ├── trade/          # TradeHandler (ads, lifecycle, chat, feedback)
│   │   │   │   ├── wallet/         # WalletHandler
│   │   │   │   ├── dispute/        # DisputeHandler
│   │   │   │   ├── notification/   # NotificationHandler
│   │   │   │   ├── platform/       # PlatformHandler (crypto assets, fiat currencies, payment methods)
│   │   │   │   ├── market/         # MarketHandler
│   │   │   │   ├── helper.go       # Shared: handleError, getUserIDFromContext
│   │   │   │   └── health.go       # Health + readiness probes
│   │   │   ├── middleware/         # AuthMiddleware, CORSMiddleware, RateLimitMiddleware
│   │   │   ├── dto/                # Request / response structs
│   │   │   ├── validator/          # Input normalization
│   │   │   └── router.go           # Route registration
│   │   └── websocket/              # WebSocket hub (trade chat, notifications)
│   │
│   └── infrastructure/             # Layer 4 — external adapters
│       ├── blockchain/             # EVM escrow client, USDT wallet client, mock escrow
│       ├── email/                  # Template-based email queue (DB-backed)
│       ├── market/                 # Mock market data provider (rate feed placeholder)
│       ├── notification/           # SMTP email client, WebSocket notifier
│       ├── payment/                # bKash, Nagad, Wise, PayPal stubs
│       ├── persistence/
│       │   ├── postgres/           # Repository implementations per domain
│       │   │   ├── identity/       # user, oauth, email_verification, password_reset,
│       │   │   │                   # session, twofactor, payment_method
│       │   │   ├── trading/        # ad, trade, message, feedback
│       │   │   ├── wallet/         # wallet, transaction
│       │   │   ├── platform/       # crypto_asset, fiat_currency, payment_method
│       │   │   ├── dispute/        # dispute
│       │   │   └── notification/   # notification
│       │   └── redis/              # Redis client, rate limiter, session store, cache
│       ├── storage/                # MinIO / S3 object storage
│       └── worker/                 # Asynq worker + scheduler (trade reconciliation)
│
├── pkg/                            # Shared utilities — no domain imports
│   ├── apperrors/                  # Typed error codes + HTTP status mapping
│   ├── config/                     # Environment-based config loader
│   ├── crypto/                     # bcrypt password hashing
│   ├── database/                   # sql.DB factory + migrator
│   ├── jwt/                        # JWT issue + parse helpers
│   ├── logger/                     # Zerolog wrapper
│   └── pagination/                 # Generic pagination types
│
├── migrations/                     # 25 numbered SQL migrations (000–025)
├── seeder/                         # Development data seeder
├── contracts/                      # Solidity escrow contracts (Foundry)
├── apps/
│   ├── web/                        # Next.js frontend
│   └── mobile/                     # Mobile app
├── tests/
│   ├── unit/
│   ├── integration/
│   └── e2e/
├── .air.toml                       # Air hot-reload config
├── docker-compose.yml              # Local dev: Postgres + Redis + MinIO
├── Makefile
└── go.mod
```

---

## Getting Started

### Prerequisites

| Tool | Version | Purpose |
|---|---|---|
| Go | 1.25+ | Backend |
| Docker & Compose | Latest | Local infrastructure |
| Make | Any | Task runner |
| Air | Latest | Hot reload (`go install github.com/air-verse/air@latest`) |
| Foundry | Latest | Smart contract dev |

### 1 — Clone and install dependencies

```bash
git clone https://github.com/parvejhossain55/cryplio.git
cd cryplio
go mod download
```

### 2 — Start infrastructure services

```bash
make env-up
# Starts: PostgreSQL 17, Redis 7, MinIO
```

### 3 — Configure environment

```bash
cp .env.example .env
# Edit .env — see Environment Variables section below
```

### 4 — Run database migrations

```bash
make migrate-up
```

### 5 — Seed development data (optional)

```bash
make seed
# Creates: 1 admin, 2 merchants, 5 traders, wallets, ads, trades, disputes
```

### 6 — Start the API server

```bash
# Development (with hot reload)
make dev-backend

# Or plain Go run
make run
```

API available at: `http://localhost:8080`  
MinIO console: `http://localhost:9001` (cryplio / cryplio123)

---

## Environment Variables

Copy `.env.example` to `.env` and fill in the values. Required fields are marked *.

```env
# ── App ──────────────────────────────────────────────────
APP_ENV=development          # development | production
SERVER_PORT=8080
FRONTEND_URL=http://localhost:3000

# ── JWT / Auth * ─────────────────────────────────────────
JWT_SECRET=your-32-char-secret       # * Required; must be set in production
JWT_EXPIRY=24h
REFRESH_TOKEN_EXPIRY=168h            # 7 days

# Cookie settings
COOKIE_NAME=auth_token
COOKIE_SECURE=false                  # Set true in production (HTTPS)
COOKIE_SAME_SITE=strict
ISSUER_NAME=Cryplio                  # Shown in TOTP authenticator apps

# ── Database * ───────────────────────────────────────────
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=parvej
DB_NAME=cryplio_db

# ── Redis ────────────────────────────────────────────────
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# ── Object Storage (MinIO / S3) * ────────────────────────
S3_ENDPOINT=localhost:9000
S3_ACCESS_KEY_ID=cryplio
S3_SECRET_ACCESS_KEY=cryplio123
S3_USE_SSL=false
S3_BUCKET_NAME=cryplio-storage
S3_PUBLIC_BASE_URL=

# ── Email (SMTP) ─────────────────────────────────────────
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=
SMTP_PASSWORD=
EMAIL_FROM=noreply@cryplio.io

# ── Google OAuth (optional) ──────────────────────────────
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
OAUTH_REDIRECT_URL=http://localhost:8080/api/v1/auth/oauth/google/callback

# ── Blockchain (optional) ────────────────────────────────
ETH_RPC_URL=http://localhost:8545
ETH_PRIVATE_KEY=                     # Platform wallet private key
ESCROW_CONTRACT_ADDRESS=             # Deployed escrow contract address

# ── Rate Limiting ────────────────────────────────────────
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m
# When REDIS_ADDR is reachable, a distributed Redis sliding-window limiter
# is used automatically. Falls back to in-process limiter otherwise.

# ── CORS ─────────────────────────────────────────────────
CORS_ALLOWED_ORIGINS=http://localhost:3000
```

---

## Running the App

```bash
make run            # Start API server (go run)
make dev            # Start backend (Air) + frontend (npm) in parallel
make dev-backend    # Backend only with hot reload
make dev-frontend   # Frontend only (Next.js)
make build          # Compile binary → ./bin/api
make prod           # Run with APP_ENV=production
make fmt            # Format all Go code
make lint           # Run golangci-lint
make test           # Run all tests
make seed           # Seed development data
make env-up         # docker compose up -d
make env-down       # docker compose down
```

---

## Database Migrations

Migrations live in `migrations/` and are applied with `golang-migrate`. There are currently **25 migration pairs** (000–025).

```bash
# Apply all pending migrations
make migrate-up

# Roll back the most recent migration
make migrate-down

# Show current migration status
make migrate-status

# Create a new migration pair
make migrate-create name=add_merchant_tier_column
# → migrations/026_add_merchant_tier_column.up.sql
# → migrations/026_add_merchant_tier_column.down.sql
```

### Schema overview

| Migration | Tables created |
|---|---|
| 000 | PostgreSQL enum types |
| 001 | `users`, `user_stats`, `user_sessions`, token tables |
| 002 | `crypto_assets`, `fiat_currencies`, `payment_methods`, `fee_tiers` |
| 003 | `trade_status_log`, `email_templates`, `email_queue` |
| 004 | `trade_ads` |
| 005 | `trades`, `trade_messages`, `trade_attachments` |
| 006 | `trade_feedback` |
| 007 | `disputes`, `dispute_messages` |
| 008 | `wallets`, `wallet_transactions` |
| 009 | `notifications`, `notification_preferences` |
| 010 | `referrals` |
| 011 | `merchant_applications`, `merchant_analytics` |
| 012 | `audit_logs`, `admin_actions`, `platform_config`, `announcements` |
| 013 | `rate_limit_counts`, `login_attempts`, `api_request_logs` |
| 014 | Database views (`active_trade_ads`, `completed_trades`, etc.) |
| 015 | Auto-update triggers for `updated_at` columns |
| 016 | Seed data (crypto assets, fiat currencies, payment methods, config) |
| 017–025 | OAuth, 2FA, user payment methods, withdrawal approvals, notification type fix |

---

## API Reference

**Base URL:** `http://localhost:8080/api/v1`

Authenticated endpoints accept the JWT either as:
- HttpOnly cookie `auth_token`
- `Authorization: Bearer <token>` header

### Authentication

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/auth/register` | Public | Register and auto-login |
| `POST` | `/auth/login` | Public | Login (returns tokens; 2FA challenge if enabled) |
| `POST` | `/auth/logout` | Public | Invalidate session |
| `POST` | `/auth/refresh` | Cookie | Rotate refresh token |
| `GET` | `/auth/oauth/google` | Public | Redirect to Google OAuth |
| `GET` | `/auth/oauth/google/callback` | Public | OAuth callback |
| `POST` | `/auth/email/request` | Public | Send email verification link |
| `POST` | `/auth/email/verify` | Public | Verify email with token |
| `POST` | `/auth/password/reset-request` | Public | Send password reset link |
| `POST` | `/auth/password/reset` | Public | Reset password with token |
| `POST` | `/auth/2fa/complete-login` | Public | Complete 2FA challenge |

### User Profile

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/users/me` | ✅ | Get own profile + trade stats |
| `PUT` | `/users/me` | ✅ | Update username / bio |
| `POST` | `/users/me/avatar` | ✅ | Upload avatar (JPEG/PNG ≤ 2 MB) |
| `GET` | `/users/username/:username` | Public | Public user profile |
| `POST` | `/auth/2fa/setup` | ✅ | Generate TOTP secret + QR URI |
| `POST` | `/auth/2fa/verify` | ✅ | Confirm TOTP setup |
| `POST` | `/auth/2fa/disable` | ✅ | Disable 2FA (password required) |
| `GET` | `/sessions` | ✅ | List active sessions |
| `DELETE` | `/sessions/:tokenId` | ✅ | Revoke a session (force sign-out) |
| `GET` | `/users/me/payment-methods` | ✅ | List saved payment methods |
| `POST` | `/users/me/payment-methods` | ✅ | Add payment method profile |
| `PUT` | `/users/me/payment-methods/:id` | ✅ | Update payment method |
| `DELETE` | `/users/me/payment-methods/:id` | ✅ | Remove payment method |
| `PATCH` | `/users/me/payment-methods/:id/default` | ✅ | Set default payment method |

### Marketplace — Trade Advertisements

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/marketplace/ads` | Public | List active ads (filter by crypto, fiat, payment, type) |
| `POST` | `/marketplace/ads` | ✅ | Create ad |
| `GET` | `/marketplace/my-ads` | ✅ | List own ads |
| `PUT` | `/marketplace/ads/:id` | ✅ | Update own ad |
| `DELETE` | `/marketplace/ads/:id` | ✅ | Soft-delete own ad |
| `PATCH` | `/marketplace/ads/:id/status` | ✅ | Pause / resume ad |

### Marketplace — Trades

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/marketplace/ads/:id/trades` | ✅ | Initiate trade from ad |
| `GET` | `/marketplace/trades` | ✅ | List own trades |
| `GET` | `/marketplace/trades/:id` | ✅ | Trade detail |
| `PATCH` | `/marketplace/trades/:id/status` | ✅ | `pay` / `release` / `cancel` |
| `POST` | `/marketplace/trades/:id/dispute` | ✅ | Raise dispute |
| `POST` | `/marketplace/trades/:id/feedback` | ✅ | Leave post-trade feedback |
| `GET` | `/marketplace/trades/:id/messages` | ✅ | Fetch chat history |
| `POST` | `/marketplace/trades/:id/messages` | ✅ | Send text or file message |

### Wallet

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/wallet` | ✅ | Create wallet for a crypto asset |
| `GET` | `/wallet/balance` | ✅ | All wallet balances |
| `GET` | `/wallet/deposit/:crypto` | ✅ | Get deposit address |
| `POST` | `/wallet/withdraw` | ✅ (2FA required) | Request withdrawal |
| `GET` | `/wallet/transactions` | ✅ | Transaction history |

### Notifications

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/notifications` | ✅ | List notifications |
| `PATCH` | `/notifications/:id/read` | ✅ | Mark as read |
| `GET` | `/notifications/preferences` | ✅ | Get notification preferences |
| `POST` | `/notifications/preferences` | ✅ | Save notification preferences |

### Market Rates

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/market/rates` | Public | All rates (all crypto × all fiat) |
| `GET` | `/market/rates/:crypto` | Public | All fiat rates for one crypto |
| `GET` | `/market/rates/:crypto/:fiat` | Public | Single crypto/fiat rate |

### Dispute Evidence

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/disputes/:id/evidence` | ✅ | Upload evidence file (≤ 10 MB) |

### Admin (role: `admin` required)

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/admin/dashboard/stats` | Aggregated platform stats |
| `GET` | `/admin/users` | Paginated user list |
| `POST` | `/admin/users/:id/suspend` | Suspend user |
| `POST` | `/admin/users/:id/unsuspend` | Lift suspension |
| `POST` | `/admin/users/:id/ban` | Permanent ban |
| `POST` | `/admin/users/:id/unban` | Lift ban |
| `GET` | `/admin/trades` | All trades (admin view) |
| `POST/GET/PUT/DELETE` | `/admin/crypto-assets[/:id]` | Crypto asset catalogue |
| `POST/GET/PUT/DELETE` | `/admin/fiat-currencies[/:id]` | Fiat currency catalogue |
| `POST/GET/PUT/DELETE` | `/admin/payment-methods[/:id]` | Platform payment method catalogue |
| `GET` | `/admin/withdrawals/pending` | Pending withdrawal approvals |
| `POST` | `/admin/withdrawals/:id/approve` | Approve withdrawal |
| `POST` | `/admin/withdrawals/:id/reject` | Reject withdrawal |
| `GET` | `/admin/disputes` | All disputes |
| `GET` | `/admin/disputes/:id` | Dispute detail |
| `POST` | `/admin/disputes/:id/assign` | Assign dispute to admin |
| `POST` | `/admin/disputes/:id/resolve` | Resolve dispute |

### Health

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/health` | Liveness check |
| `GET` | `/live` | Liveness probe |
| `GET` | `/ready` | Readiness probe (pings database) |

---

## WebSocket

Connect to `/ws` for real-time trade chat and push notifications.

```
ws://localhost:8080/ws
```

**Message types** pushed from server:

| Type | Payload | When |
|---|---|---|
| `chat_message` | `ChatMessage` | A new trade message is sent |
| `trade_update` | `TradeUpdate` | Trade status changes |
| `notification` | `NotificationEvent` | Any in-app notification |
| `market_update` | `MarketUpdate` | Rate feed update |

Authentication: pass the JWT as query param `?token=<jwt>` or as a cookie on the upgrade request.

---

## Domain Guide

### Trade lifecycle

```
AD_ACTIVE
  └─▶ TRADE_PENDING    (escrow locked on-chain)
        └─▶ TRADE_ACTIVE
                ├─▶ PAID             (buyer marks paid)
                │     ├─▶ COMPLETED  (seller releases escrow → crypto sent to buyer)
                │     └─▶ DISPUTED   (buyer or seller raises dispute)
                │               └─▶ RESOLVED  (admin resolves)
                ├─▶ CANCELLED        (either party, before payment)
                └─▶ EXPIRED          (payment window elapsed, auto-cancelled)
```

### Escrow flow

1. Buyer initiates trade → `escrowClient.Lock(trade)` — seller's crypto is locked in the smart contract
2. Buyer sends fiat via the agreed payment method and marks the trade as **paid**
3. Seller verifies receipt and calls **release** → `escrowClient.Release(trade)` — crypto is sent to buyer's wallet
4. If the seller does not release within the payment window, an auto-dispute is flagged by the background worker
5. Admin reviews evidence and resolves: **release to buyer** or **return to seller**

### Payment window auto-cancellation

A background Asynq worker (`worker/`) runs every 5 minutes to:
- Expire `pending`/`active` trades whose payment window has elapsed → `EXPIRED`
- Flag `paid` trades past the grace period for auto-dispute

---

## Smart Contracts

Contracts are written in Solidity and developed with [Foundry](https://book.getfoundry.sh). The compiled ABI is consumed by the Go blockchain client in `internal/infrastructure/blockchain/`.

```bash
cd contracts

# Build
forge build

# Test (with gas report)
forge test --gas-report

# Deploy to Sepolia testnet
forge script script/DeployEscrow.s.sol \
  --rpc-url $SEPOLIA_RPC_URL \
  --private-key $DEPLOYER_PRIVATE_KEY \
  --broadcast --verify

# Regenerate Go bindings after ABI changes
abigen --abi contracts/abi/Escrow.json \
       --pkg blockchain \
       --out internal/infrastructure/blockchain/escrow_binding.go
```

When `ESCROW_CONTRACT_ADDRESS` or `ETH_RPC_URL` are not set, the application falls back to a **mock escrow client** that simulates blockchain calls in memory — safe for local development without a running node.

---

## Contributing

1. Branch from `main`: `git checkout -b feat/your-feature`
2. Follow the layer rules — domain code must never import infrastructure packages
3. Add tests for any new use case or service method
4. Run `make fmt && make lint && make test` before opening a PR
5. PRs require at least one reviewer approval

### Commit convention

```
feat(trade):    add 30-minute payment timer auto-cancel
fix(wallet):    handle insufficient balance on withdrawal
refactor(auth): split authService into focused domain files
chore(deps):    upgrade go-ethereum to v1.17
test(identity): add unit tests for password complexity validation
```

---

*Built by ❤️ Parvej Hossain*
