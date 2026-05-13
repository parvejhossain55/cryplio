# CRYPLIO
## Software Requirements Specification
### P2P Cryptocurrency Exchange Platform — MVP Edition
*Trade Crypto. Trust the Process.*

| Field | Detail |
|-------|--------|
| Document Version | 1.0 — MVP |
| Prepared By | Parvej Hossain |
| Date | Dec 2025 |
| Scope | MVP only — Phase 2+ features excluded |

---

> **MVP Scope Statement:** This SRS covers only the features required to launch Cryplio MVP. KYC, multi-chain support, mobile app, SMS notifications, and premium membership are **excluded from MVP**.

---

# 1. Introduction

## 1.1 Purpose
This SRS defines the complete functional and non-functional requirements for Cryplio MVP — a P2P USDT trading platform with escrow protection, trade chat, and dispute resolution. Intended for the development team, QA, and stakeholders.

## 1.2 MVP Scope
The MVP includes:
- Web application only (no mobile app)
- Email and Social registration
- USDT ERC20 trading on Ethereum only
- Payment methods: Bkash, Nagad, Bank Transfer
- Escrow-protected P2P trades
- Trade chat with file upload
- Basic dispute resolution (admin-mediated)
- Email + in-app notifications
- Admin panel

**Explicitly out of MVP scope:**
- KYC / identity verification
- TRON, BNB Chain, BTC, ETH, USDC support
- Mobile app (iOS/Android)
- SMS notifications
- Premium membership
- Native token
- API access for third parties
- Multi-language support

## 1.3 Definitions & Acronyms

| Term | Definition |
|------|-----------|
| P2P | Peer-to-Peer — direct trade between two users |
| Escrow | System that holds USDT until trade conditions are met |
| Maker | User who creates a trade advertisement |
| Taker | User who responds to and initiates a trade from an ad |
| USDT | Tether (ERC20) — the only trading asset in MVP |
| Fiat | Local currency (BDT, PKR, USD, etc.) |
| 2FA | TOTP-based Two-Factor Authentication (Google Authenticator) |
| JWT | JSON Web Token — used for session authentication |
| WAF | Web Application Firewall |
| MVP | Minimum Viable Product |
| ERC20 | Token standard on Ethereum blockchain |

## 1.4 Priority Conventions

| Priority | Meaning |
|----------|---------|
| **Critical** | Must ship in MVP — platform cannot launch without it |
| **High** | Ships in MVP — strongly needed for good UX |
| **Post-MVP** | Excluded from MVP — Phase 2 or later |

## 1.5 References
- Cryplio Business Plan — MVP Edition (2025)
- OWASP Top 10 Security Standards
- Ethereum ERC20 Token Standard

---

# 2. Overall System Description

## 2.1 Product Perspective
Cryplio MVP is a web-based P2P crypto trading platform. It interfaces with:
- **Ethereum blockchain** — USDT ERC20 escrow smart contract
- **SMTP** — transactional email notifications
- **AWS/DigitalOcean** — cloud hosting
- **Cloudflare** — CDN and DDoS protection
- **MinIO** — self-hosted file storage for chat uploads

## 2.2 MVP Module Overview

| Module | Core Function |
|--------|--------------|
| User Auth | Registration, login, 2FA, password reset |
| User Profile | Public profile, stats, trade history |
| Trade Ads | Create, browse, filter, manage buy/sell ads |
| Trade Engine | Initiate, escrow, chat, pay, release, complete |
| Wallet | USDT ERC20 deposit, withdrawal, balance |
| Dispute System | Raise, evidence upload, admin resolution |
| Notifications | Email + in-app for all trade events |
| Admin Panel | User mgmt, trade monitor, dispute management |

## 2.3 User Types

| User Type | Description | Access |
|-----------|-------------|--------|
| Guest | Unregistered — can browse ads only | Read-only |
| Registered User | Email-verified — can trade, deposit, withdraw | Full trading |
| Admin | Platform staff — manages disputes, users | Admin panel |
| Super Admin | System configuration access | All + settings |

## 2.4 Operating Environment
- Web: Chrome, Firefox, Safari, Edge (latest 2 versions)
- Mobile Web: iOS Safari 15+, Android Chrome 90+ (responsive)
- Backend: Golang 1.21+ on Ubuntu 22.04
- Database: PostgreSQL 14+ + Redis 7+
- Blockchain: Ethereum mainnet (USDT ERC20)

## 2.5 Assumptions
- Users have a compatible browser and internet connection
- Ethereum network remains operational
- Bkash/Nagad payment confirmation is manual (user uploads receipt)
- Platform operates under Dubai offshore entity

---

# 3. Functional Requirements

## 3.1 User Authentication

### 3.1.1 Registration

| ID | Requirement | Priority |
|----|------------|---------|
| FR-101 | User registers with email + password, Social Login (Google) | Critical |
| FR-102 | System sends email verification link; user cannot trade until verified | Critical |
| FR-105 | User chooses a unique public username at registration | Critical |
| FR-106 | Password: min 8 chars, must include uppercase + number + special char | Critical |
| FR-107 | System rejects duplicate email or username | Critical |

**Excluded from MVP:** Phone registration (FR-103), Social login (Facebook, Whatsapp) (FR-104)

### 3.1.2 Login & Session

| ID | Requirement | Priority |
|----|------------|---------|
| FR-111 | User logs in with email + password | Critical |
| FR-112 | TOTP 2FA via Google Authenticator — mandatory for withdrawals | Critical |
| FR-114 | JWT sessions expire after 24h inactivity | Critical |
| FR-116 | Account locked after 5 failed login attempts; unlock via email | Critical |
| FR-117 | Password reset via time-limited email link (expires in 15 min) | Critical |

**Excluded from MVP:** SMS 2FA (FR-113), Device management (FR-115), Remember device (FR-118)

### 3.1.3 User Profile

| ID | Requirement | Priority |
|----|------------|---------|
| FR-131 | Public profile: username, join date, trade count, completion rate, rating | Critical |
| FR-132 | Add or Update user profile photo | Critical |
| FR-133 | User can set their bio | Critical |
| FR-134 | Trade stats shown publicly: total trades, completion %, positive feedback % | High |
| FR-136 | User can block another user — blocked user cannot initiate trades with them | High |

**Excluded from MVP:**  Online status (FR-135)

---

## 3.2 Trade Advertisement System

### 3.2.1 Create Trade Ad

| ID | Requirement | Priority |
|----|------------|---------|
| FR-201 | Maker creates buy or sell ad: crypto (USDT), fiat, amount range, price, payment method | Critical |
| FR-202 | Price type: fixed (manual) OR floating (% above/below live market rate) | Critical |
| FR-203 | Maker sets min and max fiat amount per trade | Critical |
| FR-204 | Maker sets payment time limit: 15, 30, 45, or 60 minutes | Critical |
| FR-205 | Maker selects accepted payment methods: Bkash / Nagad / Bank Transfer | Critical |
| FR-206 | Maker can add optional trade instructions (max 500 characters) | High |
| FR-207 | Sell ad is only active when maker has sufficient USDT in wallet | Critical |
| FR-208 | Maker can pause, edit, or delete own ads at any time | Critical |

**Excluded from MVP:** KYC requirement on taker (FR-209)

### 3.2.2 Browse & Filter Ads

| ID | Requirement | Priority |
|----|------------|---------|
| FR-211 | All active ads shown in a paginated list with price, limits, payment badges | Critical |
| FR-212 | Filter by fiat currency (BDT, PKR, USD, etc.) | Critical |
| FR-214 | Filter by payment method (Bkash / Nagad / Bank) | Critical |
| FR-215 | Sort by: best price, completion rate, newest | High |

**Excluded from MVP:** Filter by crypto (only USDT in MVP), Search by username (FR-216), Trusted trader filter (FR-217)

### 3.2.3 Execute Trade

| ID | Requirement | Priority |
|----|------------|---------|
| FR-221 | Taker enters fiat amount → system calculates USDT amount at current ad rate | Critical |
| FR-222 | On trade start, USDT is locked in Ethereum escrow smart contract | Critical |
| FR-223 | Private real-time WebSocket chat opens between buyer and seller | Critical |
| FR-224 | Both users can upload files in chat (max 5MB, jpg/png/pdf) — for payment proof | Critical |
| FR-225 | Countdown timer shown; trade auto-cancels and escrow releases if buyer doesn't pay in time | Critical |
| FR-226 | Buyer clicks "Mark as Paid" after sending fiat | Critical |
| FR-227 | Seller clicks "Release" after confirming fiat received — USDT sent to buyer | Critical |
| FR-228 | Trade cannot be cancelled by users once USDT is locked in escrow; it must either be released by seller or refunded by admin after expiry | Critical |
| FR-229 | Every trade gets a unique alphanumeric Trade ID | Critical |
| FR-230 | Email + in-app notification sent on trade completion | Critical |

### 3.2.4 Trade History & Feedback

| ID | Requirement | Priority |
|----|------------|---------|
| FR-241 | User views full list of own trades: completed, cancelled, disputed | Critical |
| FR-242 | Each trade record shows: ID, date, USDT amount, fiat amount, counterparty, status | Critical |
| FR-244 | After trade, both parties can leave rating (positive / neutral / negative) + optional comment | Critical |

**Excluded from MVP:** CSV export (FR-243)

---

## 3.3 Wallet — USDT ERC20 Only

| ID | Requirement | Priority |
|----|------------|---------|
| FR-301 | Each user has a system-generated Ethereum wallet for USDT ERC20 | Critical |
| FR-302 | Unique deposit address shown with QR code | Critical |
| FR-303 | System detects on-chain USDT deposits and credits wallet after 12 confirmations | Critical |
| FR-304 | User can withdraw USDT to external address; requires 2FA + email confirmation | Critical |
| FR-305 | Daily withdrawal limit: $500 USD equivalent (no KYC in MVP) | Critical |
| FR-306 | Ethereum gas fee + platform withdrawal margin shown to user before confirming | Critical |
| FR-307 | Wallet shows: Available USDT, In-Escrow USDT, Pending deposits | Critical |
| FR-308 | Full transaction history: deposits and withdrawals with txn hash + Etherscan link | Critical |

**Excluded from MVP:** BTC, ETH, USDC, BNB wallets; address whitelisting (FR-309); auto-convert (FR-310)

---

## 3.4 Dispute Resolution

| ID | Requirement | Priority |
|----|------------|---------|
| FR-401 | After buyer marks paid, either party can raise a dispute before seller releases | Critical |
| FR-402 | Disputing party selects reason: payment not received / wrong amount / fraud / other | Critical |
| FR-403 | Both parties can upload up to 10 evidence files (jpg/png/pdf) | Critical |
| FR-404 | Dispute auto-assigned to an admin moderator within 1 business hour | Critical |
| FR-405 | Admin must issue ruling within 48 hours of dispute creation | Critical |
| FR-406 | 3-way dispute chat: admin + buyer + seller | High |
| FR-407 | Admin can: release escrow to buyer (Force Release) / refund to seller (Standard or Flexible) | Critical |
| FR-408 | Flexible Admin Refund: Authorized admins can refund seller at any time to resolve backend disputes | High |
| FR-410 | If seller does not release within 1h after buyer marks paid with no action → auto-dispute triggered | High |

**Excluded from MVP:** Appeal process (FR-408), dispute stats on public profile (FR-409)

---

## 3.5 Notifications

| ID | Requirement | Priority |
|----|------------|---------|
| FR-501 | Transactional emails for: registration verify, trade started, trade complete, dispute raised, withdrawal confirmed | Critical |
| FR-502 | In-app bell icon with real-time notifications for all trade events | Critical |
| FR-506 | Real-time notification when counterparty sends a trade chat message | Critical |

**Excluded from MVP:** SMS notifications (FR-503), Push notifications (FR-504), Notification preferences UI (FR-505)

---

## 3.6 Admin Panel

| ID | Requirement | Priority |
|----|------------|---------|
| FR-801 | Admin can view, search, filter, suspend, or ban any user | Critical |
| FR-802 | Admin can view all active trades | Critical |
| FR-803 | Admin can view open disputes, read evidence, and issue resolution | Critical |
| FR-804 | Admin can configure platform fee % | Critical |
| FR-809 | All admin actions are logged with timestamp and admin ID | Critical |
| FR-807 | Basic dashboard: daily volume, active users, trade count, open disputes | High |
| FR-808 | Withdrawals above $1,000 require admin approval | High |

**Excluded from MVP:** IP ban / geo-block (FR-810), Announcement system (FR-806)

---

## Excluded MVP Modules (Phase 2)

| Module | Reason |
|--------|--------|
| KYC / AML verification | Adds friction, not needed for MVP trust model |
| Premium membership | Revenue model expansion — Phase 2 |
| Native token | Phase 3 — needs platform maturity |
| Mobile app | Web-first; responsive design covers mobile MVP |
| Multi-chain (TRON, BNB) | Complexity — Ethereum USDT ERC20 is sufficient for MVP |

---

# 4. Non-Functional Requirements

## 4.1 Performance

| Metric | Target |
|--------|--------|
| Page load | < 2s on 4G |
| API response | 95% under 300ms |
| Escrow lock | < 5s after trade initiation |
| Chat message latency | < 500ms |
| Concurrent users (MVP launch) | 1,000 — scalable to 10,000 |
| Uptime | 99.9% |
| DB query | < 100ms under normal load |
| Email delivery | < 60 seconds |

## 4.2 Security

| Requirement | Detail |
|------------|--------|
| TLS 1.3 | All traffic encrypted in transit |
| AES-256 | Data encrypted at rest |
| bcrypt | Passwords hashed — cost factor 12 |
| Parameterized queries | No raw SQL — prevents injection |
| CSP headers | XSS prevention |
| CSRF tokens | All state-changing requests |
| Rate limiting | Login: 5/15min; API: 100 req/min/IP |
| 2FA mandatory | Required for all withdrawals |
| Smart contract audit | Ethereum escrow audited before mainnet |
| Cloudflare WAF | DDoS and bot protection |
| Private key storage | HashiCorp Vault / AWS/DigitalOcean Secrets |
| Admin access | Separate subdomain + IP whitelist + hardware 2FA |
| Audit logs | Immutable logs for all financial and admin actions |

## 4.3 Scalability (MVP)
- AWS/DigitalOcean managed Kubernetes for horizontal scaling
- Redis for session and rate limit caching
- Asynq (Golang) for async job queue
- Cloudflare CDN for static assets
- PostgreSQL read replica for reporting

## 4.4 Usability
- Fully responsive — 360px mobile to 1440px desktop
- Trade flow: max 5 steps for a first-time user
- English only at MVP launch
- Human-readable error messages with corrective action hints
- Onboarding tooltip on first trade

## 4.5 Reliability
- DB backups every 6 hours, 30-day retention
- Zero-downtime deployments (rolling or blue-green)
- Health monitoring + auto-alerts (PagerDuty or OpsGenie)
- Graceful error if Ethereum node is unavailable

## 4.6 Compliance
- GDPR: user data deletion within 30 days on request
- Privacy Policy and Terms of Service before launch
- No OFAC-sanctioned country users
- All logs retained per data retention policy

---

# 5. System Architecture

## 5.1 Tech Stack

| Layer | Technology |
|-------|-----------|
| Frontend | React.js / Next.js |
| Backend | Golang 1.21+ + Gin |
| Real-Time | WebSocket (Gorilla) |
| Database | PostgreSQL 14+ |
| Cache | Redis 7+ |
| Queue | Asynq (Golang) |
| Blockchain | Web3.go + Solidity (Ethereum) |
| File Storage | MinIO (self-hosted) |
| Auth | JWT + TOTP (2FA) |
| Email | (SMTP) |
| Hosting | AWS/DigitalOcean (managed K8s + managed PG) |
| CDN + Security | Cloudflare + WAF |
| Monitoring | Grafana + Prometheus |

## 5.2 MVP Blockchain

| Network | Asset | Standard | Notes |
|---------|-------|---------|-------|
| Ethereum | USDT | ERC20 | Only supported asset in MVP |

> TRON (TRC20), BNB Chain, BTC, ETH, USDC — all Phase 2.

## 5.3 MVP Payment Methods

| Method | Region | Confirmation |
|--------|--------|-------------|
| Bkash | Bangladesh | Manual (user uploads receipt) |
| Nagad | Bangladesh | Manual (user uploads receipt) |
| Bank Transfer | Global | Manual (user uploads receipt) |

> Wise, PayPal, SEPA, UPI, M-Pesa — all Phase 2.

---

# 6. Core Data Entities

| Entity | Key Fields |
|--------|-----------|
| User | user_id, email, username, password_hash, email_verified, status, 2fa_secret, created_at |
| Trade Ad | ad_id, user_id, type, fiat_currency, price_type, price, min_amount, max_amount, payment_methods, payment_window, instructions, status |
| Trade | trade_id, ad_id, buyer_id, seller_id, usdt_amount, fiat_amount, payment_method, status, timer_expires_at, created_at, completed_at |
| Wallet | wallet_id, user_id, usdt_balance, locked_balance, eth_address |
| Transaction | txn_id, wallet_id, type, amount, txn_hash, status, created_at |
| Dispute | dispute_id, trade_id, raised_by, reason, status, assigned_admin_id, resolution, created_at |
| Message | msg_id, context_type (trade/dispute), context_id, sender_id, content, file_url, created_at |
| Feedback | feedback_id, trade_id, from_user_id, to_user_id, rating, comment, created_at |
| Notification | notif_id, user_id, type, content, read, created_at |
| Admin Log | log_id, admin_id, action, target_type, target_id, created_at |

## Data Retention

| Data | Retention |
|------|----------|
| User accounts | 5 years post-closure |
| Trade records | 7 years |
| Chat messages | 2 years |
| Audit logs | Permanent |
| Deleted users | 30-day grace then purge |
| JWT sessions | 24h inactivity |

---

# 7. MVP API Endpoints

All APIs: RESTful, JSON, versioned at `/api/v1/`, JWT auth, standard error format `{ code, message, details }`.

| Method | Endpoint | Description | Auth |
|--------|---------|------------|------|
| POST | /api/v1/auth/register | Register new user | Public |
| POST | /api/v1/auth/verify-email | Verify email with token | Public |
| POST | /api/v1/auth/login | Login → receive JWT | Public |
| POST | /api/v1/auth/forgot-password | Send reset email | Public |
| POST | /api/v1/auth/reset-password | Reset password with token | Public |
| GET | /api/v1/users/me | Get own profile | Required |
| GET | /api/v1/users/:username | Get public profile | Public |
| GET | /api/v1/ads | List/filter trade ads | Public |
| POST | /api/v1/ads | Create trade ad | Required |
| PUT | /api/v1/ads/:id | Edit own ad | Required |
| DELETE | /api/v1/ads/:id | Delete own ad | Required |
| POST | /api/v1/trades | Initiate trade from ad | Required |
| GET | /api/v1/trades/:id | Get trade details | Required |
| POST | /api/v1/trades/:id/paid | Mark trade as paid | Required |
| POST | /api/v1/trades/:id/release | Release escrow | Required |
| POST | /api/v1/trades/:id/cancel | Cancel trade | Required |
| POST | /api/v1/trades/:id/dispute | Raise dispute | Required |
| POST | /api/v1/trades/:id/feedback | Leave feedback | Required |
| GET | /api/v1/trades | Get own trade history | Required |
| GET | /api/v1/wallet | Get USDT wallet info + balance | Required |
| POST | /api/v1/wallet/withdraw | Request USDT withdrawal | Required + 2FA |
| GET | /api/v1/wallet/transactions | Get deposit/withdrawal history | Required |
| GET | /api/v1/messages/:tradeId | Get trade chat messages | Required |
| POST | /api/v1/messages/:tradeId | Send trade chat message | Required |
| GET | /api/v1/notifications | Get notifications | Required |
| POST | /api/v1/notifications/read | Mark notifications read | Required |
| GET | /api/v1/market/rates | Get live USDT/fiat rates | Public |

---

# 8. Key Use Cases

## UC-01: User Registers and Starts Trading

| | |
|-|-|
| **Actor** | New User |
| **Flow** | 1. Opens cryplio.io → clicks Register 2. Enters email, username, password 3. Clicks verify link in email 4. Deposits USDT ERC20 to wallet address 5. Creates sell ad or finds buy ad 6. Initiates trade |
| **Postcondition** | User is trading within 5 minutes of registration |

## UC-02: Complete a Sell Trade (Bkash)

| | |
|-|-|
| **Actor** | Seller + Buyer |
| **Precondition** | Seller has USDT in wallet |
| **Flow** | 1. Seller posts sell ad: 100 USDT, BDT, Bkash, 30-min timer 2. USDT locked in Ethereum escrow 3. Buyer finds ad, enters BDT amount, starts trade 4. Buyer sends Bkash payment, uploads screenshot in chat 5. Buyer clicks "Mark as Paid" 6. Seller verifies screenshot and Bkash SMS, clicks "Release" 7. USDT released to buyer's wallet 8. Both rate each other |
| **Postcondition** | Trade complete; fee collected; ratings updated |

## UC-03: Dispute — Seller Not Releasing

| | |
|-|-|
| **Actor** | Buyer, Seller, Admin |
| **Flow** | 1. Buyer marks paid, seller goes quiet 2. After 1h — auto-dispute triggered 3. Buyer uploads Bkash transaction as evidence 4. Admin assigned within 1 hour 5. Admin reviews evidence in 3-way chat 6. Admin confirms payment is valid → releases escrow to buyer |
| **Postcondition** | Buyer receives USDT; seller completion rate drops |

---

# 9. Testing

| Type | Scope | Target |
|------|-------|--------|
| Unit | Golang handlers and business logic | 80% coverage |
| Integration | API + DB end-to-end | All critical paths |
| Smart Contract | Escrow on Ethereum Goerli testnet | 100% functions |
| Security | OWASP Top 10 scan | Pre-launch mandatory |
| Load | k6 — 1,000 concurrent users | Before launch |
| UAT | 50-user closed beta | Month 5 |
| Regression | Automated on every deploy | All core flows |

---

# 10. Constraints

## Technical
- Web only — no mobile app in MVP
- USDT ERC20 only — no other assets
- Manual payment confirmation (no bank API)
- MinIO requires self-hosted maintenance

## Business
- No KYC — $500/day withdrawal limit per user
- No fiat custody — all payments are user-to-user
- MVP budget: ~$30,000–$50,000 USD
- Launch in 6 months

## Regulatory
- GDPR compliant for EU users
- Offshore entity required (Dubai)
- Terms of Service and Privacy Policy before launch

---

# 11. Appendix

## Revision History

| Version | Date | Author | Notes |
|---------|------|--------|-------|
| 1.0 | Dec 2025 | Parvej Hossain | Full SRS |

## Open Issues

---

*cryplio.io | Trade Crypto. Trust the Process |*
*SRS v1.0 — MVP Edition | Dec 2025*