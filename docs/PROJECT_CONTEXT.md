# Go-WebServer Project Context

## Overview

A Go REST API implementing the **CB-PDP Gen-2** (ClearingBid Primary Direct Purchase) IPO/Offering platform. Ported from a Meteor.js/TypeScript application (`torsod/pdp06`) to Go with PostgreSQL replacing MongoDB.

**Repository:** `github.com/torsod/go-webserver`
**Location:** `/Users/torsoderquist/dev/go-webserver`

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.25 |
| Router | chi/v5 |
| Database | PostgreSQL via pgx/v5 |
| FIX Protocol | quickfixgo/quickfix v0.9.10 |
| WebSocket | gorilla/websocket |
| Decimal Math | shopspring/decimal |

## Architecture

```
cmd/server/main.go (bootstrap + wiring)
    |
    v
Handler (HTTP) --> Service (business logic) --> Store (PostgreSQL)
    |                    |
    v                    v
internal/handler/    internal/service/    internal/store/
(11 handlers)        (11 services)        (9 store interfaces + impls)
    |
    v
internal/domain/ (9 domain model files)
internal/fixengine/ (Engine interface + simulated + real)
internal/ws/ (WebSocket hub)
internal/scheduler/ (background jobs)
```

Pattern: Clean architecture with interface-based stores, context-aware operations, explicit error returns.

## Directory Structure

```
go-webserver/
├── cmd/
│   ├── server/main.go            # Server entry point
│   ├── seed/main.go              # Seed default users (pw: "pw")
│   └── seed-data/main.go         # Seed sample offerings + orders
├── internal/
│   ├── config/config.go          # Env-based config
│   ├── domain/                   # Domain models
│   │   ├── offering.go           # Offering entity + state machine
│   │   ├── order.go              # Order entity + status lifecycle
│   │   ├── trade.go              # Trade + AllocationSession
│   │   ├── settlement.go         # Settlement UTC/DTC records
│   │   ├── user.go               # User types (LM, IS, BD, MO)
│   │   ├── orderbook.go          # Aggregated orderbook view
│   │   ├── settings.go           # System settings
│   │   └── fix.go                # FIX protocol types
│   ├── fixengine/
│   │   ├── engine.go             # Engine interface
│   │   ├── simulated_engine.go   # Dev mode (no FIX connection)
│   │   └── real_engine.go        # quickfixgo Initiator
│   ├── handler/
│   │   ├── rest_api.go           # Read-only GET endpoints
│   │   ├── public_api.go         # Unauthenticated endpoints
│   │   ├── offering_handler.go   # Offering CRUD + state
│   │   ├── order_handler.go      # Order CRUD + cancel
│   │   ├── allocation_handler.go # Allocation lifecycle
│   │   ├── fix_handler.go        # FIX session + orders
│   │   ├── user_handler.go       # User management
│   │   ├── settings_handler.go   # System settings
│   │   ├── auth_handler.go       # Login/logout
│   │   ├── health.go             # Health check
│   │   └── middleware.go         # Recovery, CORS, logging
│   ├── service/
│   │   ├── offering_service.go   # Offering CRUD + validation
│   │   ├── offering_state.go     # State transitions + change log
│   │   ├── order_service.go      # Order validation + sequence
│   │   ├── orderbook.go          # Orderbook aggregation
│   │   ├── allocation.go         # Allocation algorithms
│   │   ├── trade_service.go      # Trade creation
│   │   ├── settlement.go         # UTC/DTC file generation
│   │   ├── fix_service.go        # FIX session + order flow
│   │   ├── user_service.go       # User CRUD
│   │   ├── settings_service.go   # Settings CRUD
│   │   └── auth.go               # Authentication
│   ├── store/
│   │   ├── interfaces.go         # All store interfaces
│   │   ├── db.go                 # Pool + migrations
│   │   ├── offering_store.go     # Offering queries
│   │   ├── order_store.go        # Order queries + sequence
│   │   ├── trade_store.go        # Trade + allocation queries
│   │   ├── user_store.go         # User queries
│   │   ├── settings_store.go     # Settings queries
│   │   ├── fix_store.go          # FIX session/order/log queries
│   │   └── helpers.go            # updateFields generic helper
│   ├── scheduler/                # Background jobs
│   └── ws/                       # WebSocket hub
├── migrations/                   # 7 SQL migration versions
│   ├── 001_create_users          # Users + BD/LM fields
│   ├── 002_create_offerings      # Offerings + JSONB columns
│   ├── 003_create_orders         # Orders + order_sequence_seq
│   ├── 004_create_trades         # AllocationSessions + Trades
│   ├── 005_create_fix_tables     # FIX sessions/orders/logs
│   ├── 006_create_settings       # System settings (seeded)
│   └── 007_fix_orders_clordid    # Index on cl_ord_id
├── static/
│   ├── index.html                # Landing page with links
│   ├── manage.html               # Management UI (3 tabs)
│   ├── api-test.html             # REST API tester
│   └── api-client.html           # WebSocket client
├── context-files/                # CB-FIX spec docs
├── docs/
│   ├── PROJECT_CONTEXT.md        # This file
│   └── quickfix-engine.md        # FIX engine documentation
├── .claude/settings.json         # Claude context file config
├── Makefile                      # Build/run/seed/test targets
└── go.mod
```

## Configuration (Environment Variables)

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `3000` | HTTP server port |
| `PUBLIC_API_PORT` | `0` | Separate public API port (0 = disabled) |
| `DATABASE_URL` | `postgres://localhost:5432/go_webserver?sslmode=disable` | PostgreSQL connection |
| `JWT_SECRET` | `dev-secret-change-in-production` | Auth secret |
| `FIX_SIMULATED` | `true` | Use simulated FIX engine |

## Database

**PostgreSQL database:** `go_webserver`

All tables use UUID primary keys (`gen_random_uuid()`). Complex nested data stored as JSONB. Migrations auto-run on server start.

### Key Tables

| Table | PK | Key Columns | JSONB Fields |
|-------|-----|-------------|-------------|
| users | UUID | username (UNIQUE), user_type, bd_firm_id | firm_accounts, assigned_accounts |
| offerings | UUID | symbol (UNIQUE), state, asset_type | time_windows, closing_windows_config/data, listing_minimums, excluded_broker_dealers, change_log |
| orders | UUID | symbol, side, order_type, status, order_sequence | modification_history |
| allocation_sessions | UUID | symbol, status, offering_price | group_summaries |
| trades | UUID | symbol, order_id, allocation_id (FK) | - |
| fix_sessions | UUID | session_id, connected, simulated | - |
| fix_orders | UUID | cl_ord_id, session_id, main_order_id (FK) | - |
| fix_logs | UUID | session_id, direction, level | - |
| settings | UUID | market_open/close, allocation defaults | holidays, asset_type_schedules |

### Notable Schema Details

- `orders.allocation_session_id` is `UUID` type (nullable) - must pass `nil` not `""` from Go
- `order_sequence_seq` sequence provides global ordering for orders
- `trades.allocation_id` has FK to `allocation_sessions.id`

## API Endpoints

### Read-Only (GET)

| Endpoint | Description |
|----------|-------------|
| `GET /api/health` | Health check + uptime |
| `GET /api/offerings` | All offerings (`{count, offerings}`) |
| `GET /api/orders` | All orders (`{count, orders}`) |
| `GET /api/orders/symbol/{symbol}` | Orders by symbol |
| `GET /api/orderbook/{symbol}` | Aggregated bid/offer book |
| `GET /api/cprice/{symbol}` | Clearing price |
| `GET /api/trades/{symbol}` | Trades by symbol |
| `GET /api/settlement/{sessionId}` | Settlement CSV files |
| `GET /api/settings` | System settings |
| `GET /api/users` | All users |
| `GET /api` | API documentation |

### Public (No Auth)

| Endpoint | Description |
|----------|-------------|
| `GET /api/public/offerings` | Public offerings list |
| `GET /api/public/offering/{symbol}` | Single offering |
| `GET /api/public/orderbook/{symbol}` | Public orderbook |
| `GET /api/public/snapshot` | All offerings + orderbooks |

### Write Operations

| Endpoint | Description |
|----------|-------------|
| `POST /api/auth/login` | Login (`{username, password}`) |
| `POST /api/auth/logout` | Logout |
| `POST /api/offerings` | Create offering |
| `PUT /api/offerings/{id}` | Update offering |
| `POST /api/offerings/{id}/state` | Change state (`{targetState, reason, userId}`) |
| `POST /api/orders` | Create order (`{symbol, side, orderType, quantity, price, priorityGroup, userId}`) |
| `PUT /api/orders/{id}` | Update order (`{quantity, price, minQty, userId}`) |
| `DELETE /api/orders/{id}` | Cancel order (`{userId, reason}`) |
| `POST /api/orders/cancel-all` | Bulk cancel (`{userId, symbol}`) |
| `POST /api/users` | Create user |
| `PUT /api/users/{id}` | Update user |
| `PUT /api/settings` | Update settings |

### Allocation Lifecycle

| Endpoint | Description |
|----------|-------------|
| `POST /api/allocation/temp-close` | Create temporary allocation |
| `POST /api/allocation/bust/{sessionId}` | Bust allocation |
| `POST /api/allocation/confirm/{sessionId}` | Confirm (creates trades) |
| `POST /api/allocation/cancel-leaves/{sessionId}` | Cancel remaining qty |
| `GET /api/allocation/sessions` | List sessions |
| `GET /api/allocation/sessions/{sessionId}/trades` | Session trades |

### FIX Protocol

| Endpoint | Description |
|----------|-------------|
| `POST /api/fix/session/start` | Start FIX session |
| `POST /api/fix/session/stop` | Stop session |
| `GET /api/fix/session/status` | Session status |
| `POST /api/fix/orders` | Send order (`{symbol, quantity, price, priorityGroup, bdUser}`) |
| `POST /api/fix/orders/cancel` | Cancel by ClOrdID |
| `POST /api/fix/orders/csv` | Bulk CSV upload |
| `GET /api/fix/logs` | FIX message logs |
| `GET /api/fix/orders/{sessionId}` | FIX orders in session |
| `DELETE /api/fix/logs` | Clear logs |

### WebSocket

| Endpoint | Description |
|----------|-------------|
| `GET /ws` | Real-time update hub |

## Domain Models - Key Details

### Offering States (State Machine)

```
NEW --> UPCOMING --> OPEN --> CLOSE_PENDING --> CLOSING --> CLEARING --> CLOSED
                      |
                      +--> HALTED --> OPEN (resume)
                      |           --> CANCELED
                      +--> FROZEN --> OPEN (resume)
                      +--> CANCELED
```

### Offering Key Types

- **AssetType:** `STOCK`, `BOND`, `FUND_ETF`
- **SecurityType:** `CS` (Common Stock), `CORP` (Corporate Bond), `PS` (Preferred), `NT` (Note)
- **AllocationMethod:** `PRICE_TIME`, `PRO_RATA`, `PRIORITY_GROUP_PRO_RATA`
- **CpricePublishMode:** `MIN_IPO_PRICE`, `PUBLISH_NOW`, `NOT_PUBLISHED_AUTO_ON`, `NOT_PUBLISHED_AUTO_OFF`
- **PreferentialBids:** `YES`, `NO`, `LEAD_MANAGER_BD_ONLY`

### Order Key Types

- **Side:** `BID`, `OFFER`
- **OrderType:** `COMPETITIVE`, `PREFERENTIAL`, `SECONDARY_SELL`
- **Status:** `ACTIVE`, `CANCELED`, `HALTED`, `FILLED`, `PARTIAL_FILL`, `REJECTED`
- **ExecInst:** `""` (none), `"p"` (preferential), `"t"` (DTC tracking)

### User Types

- **LM** (Lead Manager) - manages offerings, has `lm_firm_id`
- **IS** (Issuer) - issuer of securities
- **BD** (Broker-Dealer) - enters orders, has `bd_firm_id`, `bd_role`, `qsr_active`
- **MO** (Market Operations) - platform operations
- **BD Roles:** `ADMIN`, `BROKER`, `OMS_USER`, `GROUP_MANAGER`, `MASTER_BROKER`

## FIX Engine

The FIX 4.2 engine uses an **Engine interface** pattern:

- **SimulatedEngine** (default, `FIX_SIMULATED=true`): No external connection, auto-generates execution reports after 100-300ms. Used for development.
- **RealEngine** (`FIX_SIMULATED=false`): Uses quickfixgo Initiator, connects to actual FIX counterparty. MD5 password hash on Logon (tag 554).

### FIX Message Flow

1. **New Order Single (D):** ClOrdID(11), Symbol(55), Side(54), OrderQty(38), Price(44), OrdType(40)=Limit, TimeInForce(59)=GTC
2. **Execution Report (8):** Parsed by engine, updates FIXOrder status via callback
3. **Order Cancel Request (F):** OrigClOrdID, ClOrdID, Symbol, Side, OrderQty

### CSV Bulk Upload Format

```
Symbol,Quantity,Price,Priority,BD User[,Side,OrderType,Account,ExecInst,MinQty]
LNRG,1000,20.50,1,bd1
SML,500,14.00,2,bd2,BID,COMPETITIVE,,p,100
```

## Makefile Targets

| Target | Command |
|--------|---------|
| `make build` | Build server + seed + seed-data binaries |
| `make run` | Run server (`go run ./cmd/server`) |
| `make dev` | Hot reload with air (fallback to `go run`) |
| `make seed` | Seed default users (lm1, is1, bd1-bd3, mo1) |
| `make seed-data` | Seed sample offerings (LNRG, SML, GRNB) + orders |
| `make test` | Run all tests |
| `make fmt` | Format code |
| `make lint` | golangci-lint |
| `make deps` | go mod tidy + download |
| `make clean` | Remove bin/ |

## Default Seed Data

### Users (`make seed`, password: `pw`)

| Username | Type | Details |
|----------|------|---------|
| lm1 | Lead Manager | LM_FIRM_1 |
| is1 | Issuer | - |
| bd1 | Broker-Dealer | MASTER_BROKER, BD_FIRM_1, QSR active |
| bd2 | Broker-Dealer | BROKER, BD_FIRM_2, QSR active |
| bd3 | Broker-Dealer | ADMIN, BD_FIRM_3, QSR active |
| mo1 | Market Operations | - |

### Offerings (`make seed-data`)

| Symbol | Name | Asset | State | Price Range | Primary Qty |
|--------|------|-------|-------|-------------|-------------|
| LNRG | LunarEnergy Inc. | STOCK/CS | OPEN | $18-$22 | 1,000,000 |
| SML | SmallCap Technologies | STOCK/CS | OPEN | $12-$16 | 500,000 |
| GRNB | GreenBond Corp 5Y Note | BOND/CORP | UPCOMING | $99-$101 | 200,000 |

Also seeds 15 bid + 4 offer orders for LNRG, 9 bid orders for SML.

## API Response Formats

List endpoints return wrapped objects, not raw arrays:

```json
// GET /api/offerings
{"count": 3, "offerings": [...]}

// GET /api/orders
{"count": 28, "orders": [...]}
```

## Known Patterns & Gotchas

1. **UUID columns:** PostgreSQL UUID columns reject empty strings. Use `*string` (pointer) in Go structs for nullable UUID fields like `allocation_session_id`. The `pgx` driver sends `nil` pointers as SQL NULL.

2. **JSONB marshaling:** Complex fields (TimeWindows, ClosingWindowsConfig, etc.) are marshaled to JSON before INSERT/UPDATE and unmarshaled after SELECT. Store layer handles this transparently.

3. **Order sequence:** Global ordering via `order_sequence_seq` PostgreSQL sequence. Each order gets a unique monotonic sequence number.

4. **State transitions:** Offering state changes are validated in `offering_state.go` and recorded in the `change_log` JSONB array with timestamps, user, and reason.

5. **Modification history:** Order changes are tracked in `modification_history` JSONB with field name, old/new values, user, timestamp, and whether time priority is impacted.

6. **Static files:** Served at both `/static/*` and root paths (`/`, `/manage.html`, `/api-test.html`, `/api-client.html`).

7. **Graceful shutdown:** Server handles SIGINT/SIGTERM, stops FIX session, cancels context, allows 10s drain.

## Quick Start

```bash
# Prerequisites: PostgreSQL running, database created
createdb go_webserver

# Install dependencies
make deps

# Seed users and sample data
make seed
make seed-data

# Run server
make run
# or with hot reload:
make dev

# Browse to:
# http://localhost:3000           - Landing page
# http://localhost:3000/manage.html - Management UI
```
