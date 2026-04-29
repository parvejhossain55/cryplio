  
**CRYPLIO**

**Software Requirements Specification**

P2P Cryptocurrency Exchange Platform

*Trade Crypto. Trust the Process.*

| Document Version | 1.0 |
| :---- | :---- |
| **Status** | Draft |
| **Prepared By** | Cryplio Product Team |
| **Date** | 2025 |
| **Classification** | Confidential |

# **1\. Introduction**

## **1.1 Purpose**

This Software Requirements Specification (SRS) document describes the complete functional and non-functional requirements for Cryplio — a global Peer-to-Peer (P2P) cryptocurrency exchange platform. This document is intended for the development team, QA engineers, stakeholders, and investors to clearly understand what the system must do and how it must behave.

## **1.2 Scope**

Cryplio is a web and mobile-based P2P crypto trading platform that enables users to buy and sell cryptocurrencies directly with each other using a secure escrow mechanism. The platform serves global users with support for local payment methods including Bkash, Nagad, bank transfers, and international payment rails.

The system includes:

* A web application (desktop and mobile browser)

* Native mobile apps (iOS & Android — Phase 2\)

* An admin panel for platform management

* A merchant portal for high-volume traders

* RESTful APIs for third-party integrations

## **1.3 Definitions & Acronyms**

| Term | Definition |
| :---- | :---- |
| P2P | Peer-to-Peer — direct trade between two users without a central intermediary |
| Escrow | A smart contract or system that holds crypto until trade conditions are met |
| KYC | Know Your Customer — identity verification process |
| AML | Anti-Money Laundering — compliance checks against illegal financial activity |
| USDT | Tether — a USD-pegged stablecoin widely used in P2P trades |
| Maker | The user who creates a trade advertisement on the platform |
| Taker | The user who responds and initiates a trade from an existing ad |
| Merchant | A verified high-volume trader with a special badge and features |
| Fiat | Government-issued currency (BDT, USD, EUR, etc.) |
| 2FA | Two-Factor Authentication for additional account security |
| WAF | Web Application Firewall — protection against web attacks |
| API | Application Programming Interface |
| MVP | Minimum Viable Product — the first launch-ready version |
| TRC20 | Token standard on the TRON blockchain (used for USDT transfers) |

## **1.4 Document Conventions**

| Priority | Meaning |
| :---- | :---- |
| MUST / Critical | Mandatory — system cannot launch without this |
| SHOULD / High | Important — strongly recommended for launch |
| MAY / Medium | Nice to have — can be added post-launch |
| FUTURE / Low | Planned for later phases |

## **1.5 References**

* Cryplio Business Plan v1.0 (2025)

* IEEE 830-1998 SRS Standard

* FATF Virtual Asset Guidelines (2021)

* Binance P2P & Paxful platform analysis

* OWASP Top 10 Security Standards

# **2\. Overall System Description**

## **2.1 Product Perspective**

Cryplio is a standalone SaaS platform accessible via web browsers and mobile applications. It interfaces with:

* Blockchain networks (Ethereum) for escrow contracts

* Third-party KYC providers (Sumsub)

* Local payment gateways (Bkash, Nagad) via API

* Email and SMS notification services (SendGrid, Twilio)

* Cloud infrastructure (DigitalOcean)

## **2.2 Product Functions — High Level**

| Module | Core Function |
| :---- | :---- |
| User Management | Registration, login, KYC, profile, 2FA |
| Trade Engine | Create ads, match trades, escrow management |
| Wallet System | Crypto deposit, withdrawal, balance tracking |
| Payment Gateway | Fiat payment confirmation & multi-method support |
| Dispute System | Raise, manage, and resolve trade disputes |
| Merchant System | Merchant verification, dashboard, analytics |
| Notification System | Email, SMS, in-app push notifications |
| Admin Panel | Platform management, user control, analytics |
| Referral System | Invite tracking and commission payout |
| API Layer | RESTful APIs for mobile app and third-party access |

## **2.3 User Classes & Characteristics**

| User Type | Description | Access Level |
| :---- | :---- | :---- |
| Guest | Unregistered visitor — can browse trade ads only | Read-only |
| Basic User | Registered & email-verified user | Trade, wallet, chat |
| KYC Verified User | Identity-verified — higher limits & features | Full trading |
| Merchant | Verified high-volume trader with subscription | Advanced dashboard |
| Admin | Platform staff managing users, disputes, settings | Full system access |
| Super Admin | Top-level access for system configuration | All \+ configuration |

## **2.4 Operating Environment**

* Web: Chrome, Firefox, Safari, Edge (latest 2 versions)

* Mobile Web: iOS Safari 15+, Android Chrome 90+

* Backend: Golang

* Database: PostgreSQL 14+

* Blockchain: EVM-compatible chains (ETH)

## **2.5 Assumptions & Dependencies**

* Users have access to a compatible smartphone or computer

* Third-party KYC API (Sumsub) remains available and compliant

* Blockchain networks (Ethereum) remain operational

* Local payment gateways (Bkash/Nagad) provide merchant API access

* Platform will operate under an offshore legal entity (Dubai/Estonia)

# **3\. Functional Requirements**

## **3.1 User Registration & Authentication**

### **3.1.1 Registration**

| Req. ID | Requirement | Priority | Description |
| :---- | :---- | :---- | :---- |
| FR-101 | Email Registration | Critical | User must register with a valid email address and password |
| FR-102 | Email Verification | Critical | System must send a verification link; user cannot trade until email is verified |
| FR-103 | Phone Registration | High | User can optionally register/link a phone number for SMS 2FA |
| FR-104 | Social Login | Medium | Optional sign-up via Google OAuth 2.0 |
| FR-105 | Username Selection | Critical | User selects a unique public username shown on trade ads |
| FR-106 | Password Policy | Critical | Minimum 8 characters, must include uppercase, number, and special character |
| FR-107 | Duplicate Check | Critical | System must reject duplicate email or username at registration |

### **3.1.2 Login & Session**

| Req. ID | Requirement | Priority | Description |
| :---- | :---- | :---- | :---- |
| FR-111 | Email/Password Login | Critical | User authenticates with email and password |
| FR-112 | Two-Factor Auth (2FA) | Critical | TOTP-based 2FA via Google Authenticator; mandatory for withdrawals |
| FR-113 | SMS 2FA | High | Optional SMS OTP as secondary 2FA method |
| FR-114 | Session Management | Critical | JWT tokens; sessions expire after 24h of inactivity |
| FR-115 | Device Management | High | User can view and revoke active sessions/devices |
| FR-116 | Login Attempt Limit | Critical | Account locked after 5 failed attempts; unlock via email |
| FR-117 | Password Reset | Critical | Secure password reset via email with time-limited token (15 min) |
| FR-118 | Remember Device | Medium | Optional trusted device flag to skip 2FA for 30 days |

### **3.1.3 KYC Verification**

| Req. ID | Requirement | Priority | Description |
| :---- | :---- | :---- | :---- |
| FR-121 | KYC Level 0 | Critical | No verification — limited to browsing only; no trading |
| FR-122 | KYC Level 1 | Critical | Email verified — daily trade limit $500 USD equivalent |
| FR-123 | KYC Level 2 | Critical | Government ID \+ selfie — daily limit $10,000 USD equivalent |
| FR-124 | KYC Level 3 | High | Address proof \+ enhanced due diligence — unlimited trading |
| FR-125 | KYC Provider | Critical | Integration with Sumsub or Veriff API for document verification |
| FR-126 | KYC Status Display | High | User profile shows current KYC level and pending status |
| FR-127 | KYC Rejection Handling | High | User notified with reason and allowed to resubmit once within 7 days |
| FR-128 | AML Screening | Critical | Automated AML check against OFAC and UN sanctions lists on KYC approval |

### **3.1.4 User Profile**

| Req. ID | Requirement | Priority | Description |
| :---- | :---- | :---- | :---- |
| FR-131 | Profile Page | Critical | Public profile with username, join date, trade count, and rating |
| FR-132 | Profile Photo | Medium | User uploads avatar (max 2MB, jpg/png) |
| FR-133 | Bio / Description | Low | Optional 200-character bio on public profile |
| FR-134 | Trade Statistics | High | Total trades, completion rate %, positive feedback % shown publicly |
| FR-135 | Online Status | Medium | Last seen indicator visible to other users during active trade |
| FR-136 | Block User | High | User can block another user — blocked user cannot initiate trades |

## **3.2 Trade Advertisement System**

### **3.2.1 Create Trade Ad**

| Req. ID | Requirement | Priority | Description |
| :---- | :---- | :---- | :---- |
| FR-201 | Buy/Sell Ad | Critical | Maker creates buy or sell ad specifying: crypto, fiat currency, amount, price, payment methods |
| FR-202 | Price Type | Critical | Fixed price OR floating (% above/below live market rate) |
| FR-203 | Trade Limits | Critical | Maker sets min/max fiat amount per trade |
| FR-204 | Payment Window | Critical | Maker sets payment time limit (15–90 minutes) |
| FR-205 | Payment Methods | Critical | Maker selects accepted payment methods (Bkash, Nagad, Bank, etc.) |
| FR-206 | Trade Terms | High | Maker can add trade terms/instructions (500 char max) |
| FR-207 | Ad Visibility | Critical | Ad is active only when maker has sufficient crypto in escrow wallet (sell) or balance (buy) |
| FR-208 | Ad Management | Critical | Maker can pause, edit, or delete own ads at any time |
| FR-209 | KYC Requirement | High | Maker can require taker to be KYC Level 2 before initiating trade |

### **3.2.2 Browse & Filter Ads**

| Req. ID | Requirement | Priority | Description |
| :---- | :---- | :---- | :---- |
| FR-211 | Ad Listing | Critical | All active ads displayed in filterable list with price, limits, payment methods |
| FR-212 | Filter by Crypto | Critical | Filter ads by cryptocurrency (BTC, ETH, USDT, BNB, etc.) |
| FR-213 | Filter by Currency | Critical | Filter by fiat currency (BDT, USD, EUR, INR, NGN, etc.) |
| FR-214 | Filter by Payment | Critical | Filter by accepted payment method |
| FR-215 | Sort Options | High | Sort by: best price, completion rate, trade volume, newest |
| FR-216 | Search by Username | Medium | Search for specific trader's ads by username |
| FR-217 | Trusted Trader Filter | Medium | Filter to show only verified merchants or KYC-verified traders |

### **3.2.3 Execute Trade**

| Req. ID | Requirement | Priority | Description |
| :---- | :---- | :---- | :---- |
| FR-221 | Initiate Trade | Critical | Taker enters fiat amount and clicks trade — system calculates crypto amount at current rate |
| FR-222 | Escrow Lock | Critical | On trade initiation, crypto is locked in smart contract escrow from seller's wallet |
| FR-223 | Trade Chat | Critical | Private real-time chat opened between buyer and seller for the duration of the trade |
| FR-224 | File Upload in Chat | High | Users can upload payment receipts/screenshots (max 5MB, jpg/png/pdf) |
| FR-225 | Payment Timer | Critical | Countdown timer shown — if buyer doesn't mark paid before expiry, trade auto-cancels and escrow released |
| FR-226 | Mark as Paid | Critical | Buyer clicks 'Paid' button after sending fiat payment |
| FR-227 | Release Escrow | Critical | Seller clicks 'Release' after confirming payment received — crypto sent to buyer |
| FR-228 | Cancel Trade | Critical | Either party can cancel before buyer marks paid; after marking paid only seller can cancel or dispute must be raised |
| FR-229 | Trade ID | Critical | Every trade gets a unique alphanumeric trade ID for reference |
| FR-230 | Trade Confirmation | Critical | Email \+ in-app notification on trade completion with summary |

### **3.2.4 Trade History**

| Req. ID | Requirement | Priority | Description |
| :---- | :---- | :---- | :---- |
| FR-241 | Trade History | Critical | User views full list of completed, cancelled, and disputed trades |
| FR-242 | Trade Details | Critical | Each trade record shows: ID, date, crypto, amount, fiat, counterparty, status |
| FR-243 | Export History | Medium | User can export trade history as CSV for tax/accounting purposes |
| FR-244 | Feedback System | Critical | After trade, both parties can leave rating (positive/neutral/negative) and optional comment |

## **3.3 Wallet & Asset Management**

| Req. ID | Requirement | Priority | Description |
| :---- | :---- | :---- | :---- |
| FR-301 | Multi-Asset Wallet | Critical | Each user has a system-generated wallet for BTC, ETH, USDT (TRC20 & ERC20), BNB |
| FR-302 | Deposit Address | Critical | Unique deposit address per user per asset; QR code displayed |
| FR-303 | Deposit Detection | Critical | System detects on-chain deposits and credits user wallet after required confirmations |
| FR-304 | Withdrawal | Critical | User withdraws to external address; requires 2FA and email confirmation |
| FR-305 | Withdrawal Limits | Critical | Daily withdrawal limits per KYC level; Level 1: $500, Level 2: $10,000, Level 3: unlimited |
| FR-306 | Withdrawal Fee | Critical | Network fee displayed to user before confirmation; platform can add margin |
| FR-307 | Balance Display | Critical | Available balance, in-escrow balance, and pending deposits shown separately |
| FR-308 | Transaction History | Critical | Full on-chain deposit and withdrawal history with txn hash and block explorer link |
| FR-309 | Address Whitelisting | Medium | User can whitelist trusted withdrawal addresses; new addresses require 24h delay |
| FR-310 | Auto-Convert (future) | Future | Auto-convert incoming crypto to USDT stablecoin option |

## **3.4 Dispute Resolution System**

| Req. ID | Requirement | Priority | Description |
| :---- | :---- | :---- | :---- |
| FR-401 | Raise Dispute | Critical | After buyer marks paid, either party can raise a dispute before seller releases |
| FR-402 | Dispute Reason | Critical | Disputing party selects reason: payment not received / payment wrong amount / fraud / other |
| FR-403 | Evidence Upload | Critical | Both parties can upload up to 10 files (images/PDFs) as evidence |
| FR-404 | Admin Assignment | Critical | Dispute auto-assigned to available admin moderator within 1 business hour |
| FR-405 | 48-Hour Rule | Critical | Admin must issue ruling within 48 hours of dispute creation |
| FR-406 | Dispute Chat | High | Separate 3-way chat: admin \+ buyer \+ seller for the dispute |
| FR-407 | Admin Resolution | Critical | Admin can: release escrow to buyer / return escrow to seller / split escrow (partial) |
| FR-408 | Appeal Process | Medium | Losing party has 24h to appeal to senior moderator; one appeal per trade |
| FR-409 | Dispute Statistics | High | User dispute history visible on public profile (% resolved in their favor not shown) |
| FR-410 | Auto-Dispute Trigger | High | If seller does not release within 1h after buyer marks paid with no action, auto-dispute is created |

## **3.5 Notification System**

| Req. ID | Requirement | Priority | Description |
| :---- | :---- | :---- | :---- |
| FR-501 | Email Notifications | Critical | Transactional emails for: registration, trade start, trade complete, dispute, withdrawal |
| FR-502 | In-App Notifications | Critical | Real-time bell icon notifications for all trade events |
| FR-503 | SMS Notifications | High | Optional SMS for trade start and dispute raised (user opt-in) |
| FR-504 | Push Notifications | Medium | Mobile push notifications (Phase 2 mobile app) |
| FR-505 | Notification Preferences | High | User controls which notification types are enabled per channel |
| FR-506 | Trade Chat Notifications | Critical | Real-time notification when counterparty sends a message in trade chat |

## **3.6 Merchant System**

| Req. ID | Requirement | Priority | Description |
| :---- | :---- | :---- | :---- |
| FR-601 | Merchant Application | High | Verified (KYC L2+) users with 50+ completed trades can apply for merchant status |
| FR-602 | Merchant Badge | High | Verified merchant badge displayed on profile and all trade ads |
| FR-603 | Merchant Subscription | High | Monthly fee ($30–$100) for merchant benefits; auto-renewal |
| FR-604 | Merchant Dashboard | High | Advanced analytics: daily volume, completion rate, revenue, top currencies |
| FR-605 | Priority Support | Medium | Merchant disputes escalated to senior moderators with 6-hour SLA |
| FR-606 | Bulk Ad Management | Medium | Merchants can create/edit multiple ads simultaneously |
| FR-607 | Merchant API Access | Medium | API keys for merchants to integrate Cryplio into their own tools |

## **3.7 Referral & Rewards System**

| Req. ID | Requirement | Priority | Description |
| :---- | :---- | :---- | :---- |
| FR-701 | Referral Link | High | Every user has a unique referral link to share |
| FR-702 | Referral Tracking | High | System tracks successful referrals (referee must complete first trade) |
| FR-703 | Commission Payout | High | Referrer earns 20% of platform fee on referee's trades for first 90 days |
| FR-704 | Referral Dashboard | High | User views: total referrals, active referrals, total earned |
| FR-705 | Commission Withdrawal | High | Referral commission credited to wallet; can be withdrawn or used for trading |

## **3.8 Admin Panel**

| Req. ID | Requirement | Priority | Description |
| :---- | :---- | :---- | :---- |
| FR-801 | User Management | Critical | Admin views, searches, filters, suspends, or bans any user account |
| FR-802 | Trade Monitoring | Critical | Admin views all active trades and intervenes if needed |
| FR-803 | Dispute Management | Critical | Admin processes and resolves open disputes with evidence viewer |
| FR-804 | Fee Configuration | Critical | Admin configures platform fee percentages per crypto and trade type |
| FR-805 | KYC Review | Critical | Admin manually reviews flagged KYC submissions |
| FR-806 | Announcement System | High | Admin posts platform-wide announcements (banner, email blast) |
| FR-807 | Analytics Dashboard | High | Real-time charts: daily volume, active users, trade count, revenue, disputes |
| FR-808 | Withdrawal Approval | High | Large withdrawals above threshold require admin approval |
| FR-809 | Audit Logs | Critical | All admin actions logged with timestamp, admin ID, and action details |
| FR-810 | IP Ban / Geo-Block | Medium | Admin can block IPs or restrict access from specific countries |
| FR-811 | Merchant Approval | High | Admin reviews and approves merchant applications |

# **4\. Non-Functional Requirements**

## **4.1 Performance Requirements**

| Metric | Requirement |
| :---- | :---- |
| Page Load Time | Initial page load \< 2 seconds on 4G connection |
| API Response Time | 95% of API requests respond within 300ms |
| Trade Initiation | Escrow lock confirmation within 5 seconds of trade initiation |
| Chat Latency | Trade chat messages delivered within 500ms |
| Concurrent Users | System supports 10,000 concurrent users at launch; scalable to 100,000 |
| Uptime SLA | 99.9% uptime — maximum 8.7 hours downtime per year |
| Database Query | Core queries execute within 100ms under normal load |
| Notification Delivery | Email delivered within 60 seconds; SMS within 30 seconds |

## **4.2 Security Requirements**

| Requirement | Description |
| :---- | :---- |
| Data Encryption | All data encrypted in transit (TLS 1.3) and at rest (AES-256) |
| Password Hashing | Passwords hashed with bcrypt (cost factor 12\) |
| SQL Injection | All database queries use parameterized statements / ORM |
| XSS Prevention | All user inputs sanitized; Content Security Policy headers enforced |
| CSRF Protection | CSRF tokens required on all state-changing requests |
| Rate Limiting | Login: 5 attempts/15min; API: 100 req/min per IP; Trade: 30 per hour per user |
| 2FA Enforcement | 2FA mandatory for all withdrawals and KYC changes |
| Smart Contract Audit | Escrow smart contracts must pass third-party audit before deployment |
| DDoS Protection | Cloudflare WAF and DDoS mitigation at DNS level |
| Private Key Security | Platform private keys stored in HSM (Hardware Security Module) or AWS KMS |
| Penetration Testing | Annual third-party pen test; critical vulnerabilities patched within 24h |
| Admin Access | Admin panel on separate subdomain with IP whitelist and hardware 2FA |
| Audit Trail | Immutable audit logs for all financial transactions and admin actions |

## **4.3 Scalability Requirements**

* Horizontal scaling: application tier must scale via load balancer to additional nodes

* Database: read replicas for analytics and reporting queries

* Caching: Redis cache layer for market data, user sessions, and rate limits

* CDN: static assets (JS, CSS, images) served via CloudFront/Cloudflare CDN

* Queue system: RabbitMQ or AWS SQS for async tasks (notifications, blockchain monitoring)

* Microservices-ready: modular architecture allows individual services to be extracted

## **4.4 Usability Requirements**

* Platform must be fully responsive — usable on screens from 360px to 4K

* Trade flow must be completable in 5 steps or fewer for a first-time user

* Error messages must be clear, human-readable, and suggest corrective action

* Multi-language support: English at launch; Bengali, Arabic, French in Phase 2

* WCAG 2.1 Level AA accessibility compliance

* Onboarding tooltip/guide for first-time traders

## **4.5 Reliability & Availability**

* Automated database backups every 6 hours with 30-day retention

* Disaster recovery plan with RTO \< 4 hours and RPO \< 1 hour

* Health monitoring with automated alerts (PagerDuty / OpsGenie)

* Graceful degradation: if blockchain node is down, system shows maintenance notice rather than crashing

* Zero-downtime deployments using blue-green or rolling deployment strategy

## **4.6 Compliance Requirements**

* GDPR: user data deletion on request within 30 days; privacy policy required

* FATF guidelines: transaction monitoring and suspicious activity reporting

* KYC/AML: records kept for minimum 5 years as per financial regulations

* Cookie consent: GDPR-compliant cookie banner for EU users

* Terms of Service and Privacy Policy legally reviewed before launch

# **5\. System Architecture Overview**

## **5.1 Architecture Layers**

| Layer | Technology | Responsibility |
| :---- | :---- | :---- |
| Presentation | React.js / Next.js | Web UI — SSR for SEO, CSR for dynamic trading views |
| Mobile (Phase 2\) | React Native | Cross-platform iOS & Android app |
| API Gateway | Nginx \+ Golang | Request routing, rate limiting, auth middleware |
| Business Logic | Golang \+ Gin | Trade engine, escrow management, user logic |
| Real-Time | ebSocket | Live trade chat, price feeds, notifications |
| Blockchain Layer | Web3.js \+ Solidity | Smart contract interaction for escrow |
| Database | PostgreSQL | Primary relational data store |
| Cache | Redis | Sessions, rate limits, market data caching |
| Queue | Asyncq | Async jobs: notifications, blockchain monitoring |
| Object Storage | MinIO | KYC documents, chat file uploads |
| CDN | Cloudflare | Static assets, DDoS, WAF |
| Monitoring | Grafana, Prometheus | System metrics, error tracking, alerts |

## **5.2 Supported Cryptocurrencies (Launch)**

| Cryptocurrency | Network | Standard | Priority |
| :---- | :---- | :---- | :---- |
| USDT (Tether) | Ethereum | ERC20 | Critical — primary chain |
| USDC | Ethereum | ERC20  | Critical |

## **5.3 Supported Payment Methods (Launch)**

| Payment Method | Region | Type |
| :---- | :---- | :---- |
| Bkash | Bangladesh | Mobile Money |
| Nagad | Bangladesh | Mobile Money |
| Bank Transfer | Global | Fiat Banking |
| Wise (TransferWise) | Global | International Transfer |
| PayPal | Global | Online Wallet |
| SEPA Transfer | Europe | Bank Transfer |
| UPI | India | Mobile Payment |
| GCash | Philippines | Mobile Money |
| M-Pesa | Kenya / Africa | Mobile Money |

# **6\. Key Use Cases**

## **UC-01: User Registers & Completes KYC**

| Field | Detail |
| :---- | :---- |
| Actor | New User |
| Precondition | User has a valid email address |
| Main Flow | 1\. User visits cryplio.io and clicks Register2. Enters email, username, password3. Receives verification email and clicks link4. Navigates to KYC section5. Uploads government ID and selfie6. System sends to Sumsub for verification 7\. Verification approved — KYC Level 2 granted |
| Alternative Flow | KYC rejected: User notified with reason; can resubmit within 7 days |
| Postcondition | User can now trade up to $10,000/day |

## **UC-02: Seller Creates & Completes a Sell Trade**

| Field | Detail |
| :---- | :---- |
| Actor | KYC-verified Seller (Maker) \+ KYC-verified Buyer (Taker) |
| Precondition | Seller has USDT in Cryplio wallet; Buyer has Bkash account |
| Main Flow | 1\. Seller creates sell ad: 100 USDT at market \+1%, min 500 BDT, max 10,000 BDT, Bkash only2. System locks 100 USDT in escrow when ad goes live 3\. Buyer finds ad, enters 5,000 BDT, clicks Trade 4\. Trade created — 30-minute payment timer starts 5\. Buyer sends 5,000 BDT to seller's Bkash, uploads receipt 6\. Buyer clicks Mark as Paid 7\. Seller confirms receipt, clicks Release 8\. Smart contract releases USDT to buyer's wallet 9\. Both parties leave feedback |
| Alternative Flow | Seller does not release within 1h after Mark Paid: auto-dispute triggered |
| Postcondition | Buyer receives USDT; seller receives BDT; platform fee collected |

## **UC-03: Dispute Resolution**

| Field | Detail |
| :---- | :---- |
| Actor | Buyer, Seller, Admin Moderator |
| Precondition | Active trade where Buyer has marked as Paid |
| Main Flow | 1\. Buyer raises dispute: 'Seller not releasing despite payment' 2\. Buyer uploads Bkash transaction screenshot as evidence 3\. Admin assigned within 1 hour 4\. Admin reviews evidence in dispute chat with both parties 5\. Admin confirms payment is valid 6\. Admin releases escrow to buyer 7\. Dispute resolved — seller's completion rate updated |
| Alternative Flow | Seller provides evidence of non-payment: escrow returned to seller |
| Postcondition | Trade resolved; appropriate party receives funds; dispute logged |

# **7\. Data Requirements**

## **7.1 Core Data Entities**

| Entity | Key Attributes |
| :---- | :---- |
| User | user\_id, email, username, password\_hash, kyc\_level, status, created\_at, 2fa\_secret |
| Trade Ad | ad\_id, user\_id, type (buy/sell), crypto, fiat, price\_type, price, min, max, payment\_methods, status |
| Trade | trade\_id, ad\_id, buyer\_id, seller\_id, crypto\_amount, fiat\_amount, payment\_method, status, created\_at, completed\_at |
| Wallet | wallet\_id, user\_id, crypto, balance, locked\_balance, deposit\_address |
| Transaction | txn\_id, wallet\_id, type, amount, txn\_hash, status, created\_at |
| Dispute | dispute\_id, trade\_id, raised\_by, reason, status, assigned\_admin, resolution, created\_at |
| Message | msg\_id, trade\_id (or dispute\_id), sender\_id, content, file\_url, created\_at |
| KYC Record | kyc\_id, user\_id, level, document\_type, provider\_ref, status, created\_at |
| Feedback | feedback\_id, trade\_id, from\_user\_id, to\_user\_id, rating, comment, created\_at |
| Referral | referral\_id, referrer\_id, referee\_id, status, commission\_earned, created\_at |

## **7.2 Data Retention Policy**

| Data Type | Retention Period | Reason |
| :---- | :---- | :---- |
| User Account Data | 5 years after account closure | AML/KYC regulatory requirement |
| Trade Records | 7 years | Financial audit trail |
| KYC Documents | 5 years | Regulatory compliance |
| Chat Messages | 2 years | Dispute resolution reference |
| Audit Logs | Permanent | Security and compliance |
| Deleted User Data | 30 days grace, then purge | GDPR right to erasure |
| Session Tokens | 24 hours inactivity | Security |

# **8\. API Requirements**

## **8.1 REST API Standards**

* All APIs follow RESTful conventions (GET, POST, PUT, DELETE, PATCH)

* JSON request and response format for all endpoints

* API versioning via URL prefix: /api/v1/

* Authentication: Bearer JWT token in Authorization header

* Rate limiting headers included in all responses (X-RateLimit-\*)

* Standardized error response format: { code, message, details }

## **8.2 Key API Endpoints**

| Method | Endpoint | Description | Auth |
| :---- | :---- | :---- | :---- |
| POST | /api/v1/auth/register | Register new user account | Public |
| POST | /api/v1/auth/login | Login and receive JWT token | Public |
| GET | /api/v1/users/me | Get current user profile | Required |
| POST | /api/v1/kyc/submit | Submit KYC documents | Required |
| GET | /api/v1/ads | List trade advertisements with filters | Public |
| POST | /api/v1/ads | Create a new trade advertisement | Required |
| PUT | /api/v1/ads/:id | Update own trade advertisement | Required |
| POST | /api/v1/trades | Initiate a trade from an ad | Required |
| POST | /api/v1/trades/:id/paid | Mark trade as paid (buyer) | Required |
| POST | /api/v1/trades/:id/release | Release escrow (seller) | Required |
| POST | /api/v1/trades/:id/cancel | Cancel an active trade | Required |
| POST | /api/v1/trades/:id/dispute | Raise a dispute on a trade | Required |
| GET | /api/v1/wallet/balance | Get user wallet balances | Required |
| POST | /api/v1/wallet/withdraw | Request a withdrawal | Required (2FA) |
| GET | /api/v1/notifications | Get user notifications | Required |
| GET | /api/v1/market/rates | Get live crypto-fiat exchange rates | Public |

# **9\. Testing Requirements**

| Test Type | Description | Coverage Target |
| :---- | :---- | :---- |
| Unit Testing | Test individual functions and components in isolation | 80% code coverage |
| Integration Testing | Test API endpoints and database interactions end-to-end | All critical paths |
| Smart Contract Testing | Test escrow contract on testnet with edge cases | 100% function coverage |
| Security Testing | OWASP Top 10 vulnerability scan; penetration test | Pre-launch mandatory |
| Load Testing | Simulate 10,000 concurrent users using k6 or Locust | Before each major release |
| UAT | User Acceptance Testing with 50-user beta group | Before public launch |
| Regression Testing | Automated regression suite run on every deployment | All core flows |
| Mobile Testing | Cross-device testing on iOS and Android (Phase 2\) | 10+ device variants |

# **10\. Constraints & Limitations**

## **10.1 Technical Constraints**

* MVP will support web only — mobile app deferred to Phase 2

* Smart contracts will initially be deployed on Ethereum testnets before mainnet

* Real-time features require WebSocket support — users on restrictive networks may have degraded experience

* Bkash and Nagad payment confirmation is manual (user uploads receipt) — automated API integration in Phase 2

## **10.2 Business Constraints**

* Bangladesh Bank does not officially recognize crypto — platform must operate via offshore entity

* Platform cannot offer fiat custody — all fiat payments are user-to-user (outside the platform)

* MVP budget ceiling: $50,000 USD

* MVP launch timeline: 6 months from project kickoff

## **10.3 Regulatory Constraints**

* Platform must not accept users from OFAC-sanctioned countries

* All KYC data must be stored in jurisdiction-compliant data centers

* GDPR compliance mandatory for EU user data

* Financial promotions must comply with local advertising laws per target market

# **11\. Appendix**

## **11.1 Revision History**

| Version | Date | Author | Changes |
| :---- | :---- | :---- | :---- |
| 0.1 | 2025-01 | Cryplio Team | Initial draft |
| 0.5 | 2025-02 | Cryplio Team | Added functional requirements sections 3.1–3.8 |
| 1.0 | 2025-03 | Cryplio Team | Full document — ready for development team review |

## **11.2 Open Issues**

| Issue ID | Description | Owner | Status |
| :---- | :---- | :---- | :---- |
| OI-001 | Finalize escrow smart contract architecture (single vs multi-sig) | CTO | Open |
| OI-002 | Confirm Bkash merchant API availability and terms | BD Team | Open |
| OI-003 | Select final KYC provider: Sumsub vs Veriff (pricing) | Founder | Open |
| OI-004 | Legal review of Terms of Service for UAE jurisdiction | Legal | Open |
| OI-005 | Decide on native token launch timing and tokenomics | Founder | Deferred |

**Cryplio SRS v1.0 — Confidential  |  cryplio.io  |  Trade Crypto. Trust the Process.**