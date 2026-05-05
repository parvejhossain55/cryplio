# Cryplio

> **Trade Crypto. Trust the Process.**

P2P cryptocurrency exchange platform — Built with Clean Architecture and Domain-Driven Design.

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
- [Smart Contracts](#smart-contracts)
- [API Reference](#api-reference)
- [Testing](#testing)
- [Domain Guide](#domain-guide)
- [Kubernetes Deployment](#kubernetes-deployment)
- [Contributing](#contributing)

---

## Overview

Cryplio is a global P2P crypto exchange where buyers and sellers trade directly with each other using an escrow-based system. No platform custody of funds — trades are secured by a smart contract escrow on EVM-compatible chains.

Key capabilities:
- Escrow-backed P2P trades (USDT,USDC)
- Multi-payment method support — Bkash, Nagad, Bank Transfer, Wise, PayPal, UPI
- Real-time trade chat via WebSocket
- Admin dispute resolution system
- Referral program with commission tracking
- Merchant portal for high-volume traders

---

## Architecture

The codebase follows **Clean Architecture** with **Domain-Driven Design**. Dependencies flow strictly inward — outer layers depend on inner layers, never the reverse.

```
┌─────────────────────────────────────────────────┐
│  interfaces/       HTTP handlers, WebSocket     │
│  ┌───────────────────────────────────────────┐  │
│  │  application/  Use cases & orchestration  │  │
│  │  ┌─────────────────────────────────────┐  │  │
│  │  │  domain/   Entities, rules, ports   │  │  │
│  │  └─────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────┘  │
│  infrastructure/   DB, blockchain, 3rd-party    │
└─────────────────────────────────────────────────┘
```

**Layer responsibilities:**

| Layer | Package | Rule |
|---|---|---|
| Domain | `internal/domain/` | Pure business logic. Zero external imports. |
| Application | `internal/application/` | One file per use case. Orchestrates domain services. |
| Interfaces | `internal/interfaces/` | HTTP handlers, DTOs, middleware, WebSocket hub. |
| Infrastructure | `internal/infrastructure/` | Postgres repos, Redis, blockchain client, Sumsub, payment gateways. |
| Shared | `pkg/` | JWT, crypto, logger, config. No domain imports allowed. |

---

## Tech Stack

| Concern | Technology |
|---|---|
| Language | Go 1.22+ |
| HTTP Framework | Gin |
| Database | PostgreSQL 17 |
| Cache / Sessions | Redis 7 |
| ORM / Query | sqlx + raw SQL |
| Migrations | golang-migrate |
| Auth | JWT (golang-jwt) + TOTP (pquerna/otp) |
| WebSocket | gorilla/websocket |
| Blockchain | go-ethereum (geth) |
| Smart Contracts | Solidity + Foundry |
| Email | SendGrid |
| SMS | Twilio |
| Push Notifications | Firebase Cloud Messaging |
| File Storage | AWS S3 |
| Config | Viper |
| Logging | Zerolog |
| Containerization | Docker + Kubernetes (K8s) |

---

## Project Structure

```
cryplio-backend/
├── cmd/
│   ├── api/            # HTTP server entrypoint
│   ├── worker/         # Background job runner
│   └── migrate/        # Migration CLI entrypoint
│
├── internal/
│   ├── domain/         # Layer 1 — pure business logic
│   │   ├── user/       # entity, repository interface, service, value objects
│   │   ├── trade/      # TradeAd, Trade, matching logic, fee calc
│   │   ├── escrow/     # EscrowLock state machine
│   │   ├── wallet/     # Balance, locked_balance, transactions
│   │   ├── dispute/    # Dispute lifecycle, evidence rules
│   │   ├── merchant/   # Merchant subscription & tier rules
│   │   ├── notification/
│   │   ├── referral/
│   │   └── market/     # Rate feed aggregation
│   │
│   ├── application/    # Layer 2 — use cases
│   │   ├── user/       # register, login, session
│   │   ├── trade/      # create_ad, initiate, mark_paid, release, cancel
│   │   ├── dispute/    # raise, assign, resolve
│   │   ├── wallet/     # deposit, withdraw, balance
│   │   ├── merchant/   # apply, verify, dashboard
│   │   └── referral/   # track, payout
│   │
│   ├── interfaces/         # Layer 3 — delivery
│   │   ├── http/
│   │   │   ├── handler/    # one file per domain
│   │   │   ├── middleware/ # jwt_auth, rate_limit, cors
│   │   │   ├── dto/        # request + response structs
│   │   │   ├── validator/
│   │   │   └── router.go
│   │   └── ws/             # WebSocket hub (trade chat, notifications)
│   │
│   └── infrastructure/     # Layer 4 — external adapters
│       ├── persistence/
│       │   ├── postgres/   # repository implementations
│       │   └── redis/      # session store, rate limiter, cache
│       ├── blockchain/     # go-ethereum escrow contract client
│       ├── payment/        # bkash, nagad, wise, paypal adapters
│       ├── notification/   # sendgrid, twilio, firebase
│       ├── market/         # coingecko / binance rate feed
│       └── storage/        # AWS S3 — dispute evidence, user uploads
│
├── pkg/                # Shared utilities — no domain imports
│   ├── jwt/
│   ├── crypto/         # bcrypt, TOTP
│   ├── pagination/
│   ├── apperrors/      # typed error codes + HTTP status mapping
│   ├── logger/
│   └── config/
│
├── migrations/         # Numbered SQL files (golang-migrate)
├── contracts/
│   ├── src/            # Escrow.sol, EscrowFactory.sol
│   ├── test/           # Foundry tests (*.t.sol)
│   ├── script/         # Foundry deploy scripts (*.s.sol)
│   └── abi/            # Compiled ABI consumed by Go blockchain client
├── tests/
│   ├── unit/
│   ├── integration/    # testcontainers — real Postgres + Redis
│   └── e2e/
│
├── k8s/
│   ├── base/           # Deployments, Services, ConfigMaps
│   │   ├── api-deployment.yaml
│   │   ├── worker-deployment.yaml
│   │   ├── postgres-statefulset.yaml
│   │   └── redis-statefulset.yaml
│   └── overlays/
│       ├── staging/    # Kustomize patches for staging
│       └── production/ # Kustomize patches for production
├── docker compose.yml  # Local development only
├── Makefile
└── go.mod
```

---

## Getting Started

### Prerequisites

- Go 1.22+
- Docker & Docker Compose (local dev)
- Foundry (`foundryup`) — for smart contract development
- `kubectl` + access to a Kubernetes cluster (staging/production)
- `make`

### Clone and install

```bash
git clone https://github.com/your-org/cryplio-backend.git
cd cryplio-backend
go mod download
```

### Start infrastructure services

```bash
docker compose up -d postgres redis
```

### Copy and configure environment

```bash
cp .env.example .env
# Edit .env with your credentials — see Environment Variables below
```

### Run database migrations

```bash
make migrate-up
```

### Start the API server

```bash
make run
```

The API will be available at `http://localhost:8080`.

---

## Environment Variables

```env
# Server
APP_ENV=development
APP_PORT=8080
APP_SECRET=your-32-char-secret-key

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=cryplio
DB_USER=cryplio
DB_PASSWORD=secret
DB_MAX_CONNS=25

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=your-jwt-secret
JWT_EXPIRY_HOURS=24

# Blockchain
ETH_RPC_URL=https://mainnet.infura.io/v3/YOUR_KEY
ESCROW_CONTRACT_ADDRESS=0x...
PLATFORM_WALLET_PRIVATE_KEY=0x...

# Email — SendGrid
SENDGRID_API_KEY=SG.xxx
EMAIL_FROM=noreply@cryplio.io


# Rate limiting
RATE_LIMIT_RPS=100
RATE_LIMIT_BURST=20
```

---

## Running the App

```bash
# Development backend + frontend
make dev

# Production build
make build
./bin/api

# Run all linters
make lint

# Format code
make fmt
```

### Makefile targets

| Target | Description |
|---|---|
| `make run` | Start API server |
| `make dev` | Start backend and frontend in development mode |
| `make build` | Compile binary to `./bin/api` |
| `make test` | Run all tests |
| `make test-unit` | Unit tests only |
| `make lint` | Run golangci-lint |
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down` | Roll back last migration |
| `make migrate-create name=xxx` | Create a new migration file |

---

## Database Migrations

Migrations live in `migrations/` and are numbered sequentially.

```bash
# Apply all
make migrate-up

# Roll back one step
make migrate-down

# Create a new migration
make migrate-create name=add_merchant_tier_column
# creates: migrations/00X_add_merchant_tier_column.up.sql
#          migrations/00X_add_merchant_tier_column.down.sql
```

Migration files follow the pattern: `{version}_{description}.{up|down}.sql`

---

## Smart Contracts

Contracts are written in Solidity and developed with [Foundry](https://book.getfoundry.sh). The compiled ABI is consumed by the Go blockchain client in `internal/infrastructure/blockchain/`.

### Install Foundry

```bash
curl -L https://foundry.paradigm.xyz | bash
foundryup
```

### Workflow

```bash
cd contracts

# Build contracts
forge build

# Run all tests (with gas report)
forge test --gas-report

# Run a specific test file
forge test --match-path test/Escrow.t.sol -vvv

# Check test coverage
forge coverage

# Format Solidity files
forge fmt

# Deploy to Sepolia testnet
forge script script/DeployEscrow.s.sol \
  --rpc-url $SEPOLIA_RPC_URL \
  --private-key $DEPLOYER_PRIVATE_KEY \
  --broadcast \
  --verify

# Export ABI to Go-consumable path
cp out/Escrow.sol/Escrow.json ../contracts/abi/
```

### Foundry project layout

```
contracts/
├── src/
│   ├── Escrow.sol          # Core escrow logic
│   └── EscrowFactory.sol   # Factory for per-trade escrow instances
├── test/
│   ├── Escrow.t.sol        # Unit tests (forge test)
│   └── EscrowFactory.t.sol
├── script/
│   └── DeployEscrow.s.sol  # Deployment script (forge script)
├── abi/                    # Exported ABI JSON for Go abigen
├── foundry.toml            # Foundry config
└── remappings.txt
```

After any contract change, regenerate the Go bindings:

```bash
abigen --abi contracts/abi/Escrow.json \
       --pkg blockchain \
       --out internal/infrastructure/blockchain/escrow_binding.go
```

---

## API Reference

Base URL: `https://cryplio.io/api/v1`

All authenticated endpoints require:
```
Authorization: Bearer <jwt_token>
```

### Auth

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/auth/register` | Public | Register new user |
| `POST` | `/auth/login` | Public | Login, receive JWT |
| `POST` | `/auth/logout` | Required | Invalidate session |
| `POST` | `/auth/refresh` | Required | Refresh JWT |
| `POST` | `/auth/2fa/enable` | Required | Enable TOTP 2FA |
| `POST` | `/auth/2fa/verify` | Required | Verify TOTP code |
| `POST` | `/auth/password/reset` | Public | Request password reset |

### Users

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/users/me` | Required | Get own profile |
| `PUT` | `/users/me` | Required | Update profile |
| `GET` | `/users/me/devices` | Required | List active sessions |
| `DELETE` | `/users/me/devices/:id` | Required | Revoke a session |

### Trade Ads

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/ads` | Public | List ads (filters: crypto, fiat, payment, type) |
| `POST` | `/ads` | Required | Create ad |
| `GET` | `/ads/:id` | Public | Get single ad |
| `PUT` | `/ads/:id` | Required | Update own ad |
| `DELETE` | `/ads/:id` | Required | Delete own ad |

### Trades

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/trades` | Required | Initiate trade from ad |
| `GET` | `/trades/:id` | Required | Get trade details |
| `POST` | `/trades/:id/paid` | Required (buyer) | Mark as paid |
| `POST` | `/trades/:id/release` | Required (seller) | Release escrow |
| `POST` | `/trades/:id/cancel` | Required | Cancel trade |
| `POST` | `/trades/:id/dispute` | Required | Raise dispute |
| `GET` | `/trades/:id/messages` | Required | Get trade chat messages |

### Wallet

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/wallet/balance` | Required | Get all balances |
| `GET` | `/wallet/deposit/:crypto` | Required | Get deposit address |
| `POST` | `/wallet/withdraw` | Required (2FA) | Request withdrawal |
| `GET` | `/wallet/transactions` | Required | Transaction history |

### Market

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/market/rates` | Public | Live crypto-fiat rates |

### Dispute (Admin)

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `GET` | `/admin/disputes` | Admin | List open disputes |
| `GET` | `/admin/disputes/:id` | Admin | Get dispute detail |
| `POST` | `/admin/disputes/:id/assign` | Admin | Assign to self |
| `POST` | `/admin/disputes/:id/resolve` | Admin | Resolve with decision |

---

## Testing

```bash
# All tests
make test

# Unit tests only (no Docker needed)
make test-unit

# With coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Testing conventions

- **Unit tests** live alongside the code they test: `domain/trade/service_test.go`
- **Integration tests** can live in `tests/integration/` when database or adapter coverage is added
- **Repository interfaces** in `domain/` make every use case fully mockable without a database

---

## Domain Guide

### Trade lifecycle

```
AD_ACTIVE
    └─▶ TRADE_CREATED   (escrow locked)
            ├─▶ PAID            (buyer marks paid → 1h seller window opens)
            │       ├─▶ COMPLETED   (seller releases escrow)
            │       └─▶ DISPUTED    (buyer/seller raises dispute)
            │               └─▶ RESOLVED  (admin decides)
            └─▶ CANCELLED       (before payment, by either party)
```

### Escrow flow

1. Seller creates ad → `escrow.Lock(amount)` called immediately
2. Trade initiated → escrow remains locked
3. Buyer marks paid → 1-hour release window starts for seller
4. Seller releases → `escrow.Release(buyerAddress)` — smart contract sends USDT to buyer
5. If seller doesn't release within 1h → auto-dispute triggered
6. Admin resolves dispute → calls either `escrow.Release` or `escrow.Return`

---

## Kubernetes Deployment

Kubernetes manifests live in `k8s/` and are managed with [Kustomize](https://kustomize.io) (built into `kubectl`).

### Structure

```
k8s/
├── base/                   # Shared base resources
│   ├── api-deployment.yaml
│   ├── worker-deployment.yaml
│   ├── postgres-statefulset.yaml
│   └── redis-statefulset.yaml
└── overlays/
    ├── staging/            # Patches: replica count, image tags, env refs
    └── production/         # Patches: HPA, resource limits, PDB
```

### Deploy

```bash
# Staging
kubectl apply -k k8s/overlays/staging

# Production
kubectl apply -k k8s/overlays/production

# Check rollout status
kubectl rollout status deployment/cryplio-api -n cryplio

# Roll back one revision
kubectl rollout undo deployment/cryplio-api -n cryplio
```

### Secrets

Never commit secrets to the repo. Inject them as Kubernetes Secrets sourced from your secrets manager (AWS Secrets Manager, Vault, or Sealed Secrets):

```bash
kubectl create secret generic cryplio-secrets \
  --from-literal=DB_PASSWORD=... \
  --from-literal=JWT_SECRET=... \
  --from-literal=ETH_RPC_URL=... \
  -n cryplio
```

Reference in your deployment manifest:

```yaml
envFrom:
  - secretRef:
      name: cryplio-secrets
```

---

## Contributing

1. Branch from `main`: `git checkout -b feat/your-feature`
2. Follow the layer rules — domain code must never import infrastructure
3. Write tests for any new use case
4. Run `make lint` and `make test` before opening a PR
5. PRs require at least one reviewer approval

### Commit convention

```
feat(trade): add 30-minute payment timer auto-cancel
fix(escrow): handle reorg on release tx confirmation
chore(deps): upgrade go-ethereum to v1.13
```

---
