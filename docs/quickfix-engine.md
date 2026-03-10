# QuickFIX/Go FIX 4.2 Engine Integration

## Overview

Added a FIX 4.2 protocol engine to the go-webserver using the `github.com/quickfixgo/quickfix` library. The engine supports both a real quickfixgo initiator for production FIX connections and a simulated mode for development and testing. All FIX operations are exposed via REST endpoints under `/api/fix/`.

## Architecture

```
REST Handler (/api/fix/*)
    │
    ▼
FIXService (business logic, CSV parsing, order linkage)
    │
    ├── Engine interface
    │       ├── RealEngine    (quickfixgo Initiator + Application)
    │       └── SimulatedEngine (dev mode, fake responses)
    │
    ├── FIXSessionStore  (PostgreSQL)
    ├── FIXOrderStore    (PostgreSQL)
    ├── FIXLogStore      (PostgreSQL)
    └── OrderStore       (main orders table linkage)
```

## Files Created

| File | Purpose |
|------|---------|
| `internal/fixengine/engine.go` | Engine interface, ExecutionReport struct, SendOrderParams struct |
| `internal/fixengine/real_engine.go` | quickfixgo Application implementation with FIX 4.2 Initiator |
| `internal/fixengine/simulated_engine.go` | Development-mode engine with simulated FIX responses |
| `internal/service/fix_service.go` | Business logic: session management, order flow, CSV parsing, execution report handling |
| `internal/handler/fix_handler.go` | REST API endpoints for all FIX operations |
| `migrations/007_fix_orders_clordid_index.up.sql` | Index on fix_orders.cl_ord_id for order lookups |
| `migrations/007_fix_orders_clordid_index.down.sql` | Rollback for the index |

## Files Modified

| File | Changes |
|------|---------|
| `internal/domain/fix.go` | Added FIXConnectionSettings, CSVOrder, DefaultFIXSettings(), FIX order status constants |
| `internal/store/interfaces.go` | Extended FIXOrderStore (FindByClOrdID, UpdateFields) and FIXLogStore (FindBySessionID, DeleteBySessionID) |
| `internal/store/fix_store.go` | Implemented the new store interface methods |
| `internal/config/config.go` | Added FIXSimulated config field (env: FIX_SIMULATED, default: true) |
| `cmd/server/main.go` | Wired FIX stores, engine, service, handler, routes, orphan cleanup, and graceful shutdown |
| `go.mod` | Added github.com/quickfixgo/quickfix and github.com/shopspring/decimal dependencies |

## REST API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/fix/session/start` | Start a FIX session (accepts optional connection settings JSON) |
| POST | `/api/fix/session/stop` | Stop the active FIX session |
| GET | `/api/fix/session/status` | Get current session status and connection info |
| POST | `/api/fix/orders` | Send a single order via FIX |
| POST | `/api/fix/orders/cancel` | Cancel an order by clOrdId |
| POST | `/api/fix/orders/csv` | Bulk upload orders from CSV content with 50ms throttling |
| GET | `/api/fix/logs?sessionId=&limit=` | Get FIX message logs |
| GET | `/api/fix/orders/{sessionId}` | Get FIX orders for a session |
| DELETE | `/api/fix/logs?sessionId=` | Clear FIX logs |

## FIX 4.2 Message Types

### Outgoing

**New Order Single (MsgType=D)**
- Tag 11: ClOrdID (auto-generated: ORD-{timestamp}-{random})
- Tag 21: HandlInst = 2 (Automated execution, public)
- Tag 55: Symbol
- Tag 54: Side (1=Buy, 2=Sell)
- Tag 60: TransactTime (UTC)
- Tag 38: OrderQty
- Tag 40: OrdType = 2 (Limit)
- Tag 44: Price
- Tag 1: Account (optional)
- Tag 59: TimeInForce = 1 (GTC)
- Tag 58: Text (encodes priority group)

**Order Cancel Request (MsgType=F)**
- Tag 41: OrigClOrdID
- Tag 11: ClOrdID (auto-generated: CXL-{timestamp}-{random})
- Tag 55: Symbol
- Tag 54: Side
- Tag 38: OrderQty

### Incoming

**Execution Report (MsgType=8)**
- ExecType 0 = New → FIX order status NEW
- ExecType 1 = Partial Fill → PARTIAL
- ExecType 2 = Fill → FILLED
- ExecType 4 = Canceled → CANCELED
- ExecType 8 = Rejected → REJECTED

## Engine Modes

### Simulated Mode (default)

Enabled by `FIX_SIMULATED=true` (the default). The simulated engine:
- Simulates connection with a 200ms delay
- Auto-generates execution reports after 100-300ms random delay
- Returns ExecType=0 (New) for orders, ExecType=4 (Canceled) for cancels
- Runs a heartbeat goroutine at the configured interval
- No external FIX connection required

### Real Mode

Enabled by `FIX_SIMULATED=false`. The real engine:
- Creates a quickfixgo Initiator with programmatic settings (no .cfg file)
- Implements the quickfix.Application interface (OnCreate, OnLogon, OnLogout, ToAdmin, FromAdmin, ToApp, FromApp)
- Injects MD5-hashed password on Logon (tag 554) per ClearingBid spec
- Routes incoming Execution Reports (MsgType=8) through a callback to the service layer
- Uses in-memory message store
- Supports configurable host, port, SenderCompID, TargetCompID, heartbeat interval

## Order Flow

1. Client sends POST `/api/fix/orders` with order details (symbol, quantity, price, priorityGroup, bdUser)
2. FIXService generates a unique ClOrdID
3. A main Order record is created in the orders table (status=ACTIVE)
4. A FIXOrder record is created in the fix_orders table (status=PENDING, linked to main order via main_order_id)
5. The FIX message is logged and sent via the engine
6. When the engine receives an Execution Report, the callback updates the FIXOrder status
7. Session message counters are incremented

## CSV Bulk Upload

POST `/api/fix/orders/csv` accepts CSV content:

```
Symbol,Quantity,Price,Priority,BD User
LNRG,1000,25.50,1,bd1
LNRG,2500,26.00,2,bd1
SML,5000,21.00,1,bd2
```

Orders are submitted with a 50ms delay between each to stay under the 1000 msg/sec throttle limit.

## Session Lifecycle

- **Startup**: Server cleans up any orphaned sessions (still marked connected=true from a previous crash)
- **Start Session**: Creates DB record, starts engine, registers execution report callback
- **Stop Session**: Sends Logout, stops engine, marks session disconnected
- **Shutdown**: Graceful FIX session stop before server exit

## Configuration

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| FIX_SIMULATED | true | Use simulated engine (true) or real quickfixgo (false) |

Connection settings are passed via the POST `/api/fix/session/start` request body:

```json
{
  "host": "localhost",
  "port": 9878,
  "senderCompId": "CLIENT_COMP",
  "targetCompId": "CBID",
  "username": "",
  "password": "",
  "heartbeatInterval": 30,
  "simulated": true
}
```

## Dependencies Added

- `github.com/quickfixgo/quickfix` v0.9.10 — FIX protocol engine
- `github.com/shopspring/decimal` v1.4.0 — Required by quickfixgo for price/quantity fields
