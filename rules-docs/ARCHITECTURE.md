# ClearingBid — Architecture

## System Context

ClearingBid sits between retail/institutional investors and the capital markets infrastructure that manages IPO book-building. The data flow:

```
Investors (browser) → BidMan UI → REST/WS API → Go Backend → Postgres queue → RA Adapter → FIX Engine
Lead Managers (browser) → BatMan UI → REST/WS API → Go Backend → (same)
```

Inbound flow (execution reports, acks) is the reverse: FIX Engine → RA Adapter → Postgres queue → Go Backend → WebSocket → BatMan/BidMan UIs.

The system has **spiky demand**: near-zero activity between IPOs, moderate activity during pre-IPO setup, and very high concurrent load during live book-building windows.

## Deployment Topology

### Production (Target)

```
                    ┌─────────┐
                    │   ALB   │
                    └────┬────┘
                         │
              ┌──────────┼──────────┐
              │          │          │
         ┌────┴───┐ ┌───┴────┐ ┌──┴─────┐
         │ Task 1 │ │ Task 2 │ │ Task N │   ← ECS Fargate (Go binary, ~20MB image)
         └────┬───┘ └───┬────┘ └──┬─────┘
              │          │          │
              └──────┬───┼──────────┘
                     │   │
            ┌────────┘   └────────┐
            │                     │
       ┌────┴────┐          ┌─────┴─────┐
       │   RDS   │          │   Redis   │  ← ElastiCache / Valkey (cache only)
       │ (PG 16) │          │           │
       └────┬────┘          └───────────┘
            │
       ┌────┴─────────┐
       │  RA Adapter  │  ← Separate deployment (Java, own JVM)
       │  + FIX Engine│
       └──────────────┘
```

The Go binary is a static build with no OS dependencies. Docker image is `FROM scratch` or `FROM alpine` — under 30MB. Startup is instant (~100ms). Auto-scaling is target-tracking on CPU utilization and/or active connection count.

Redis is cache-only (sessions, rate limiting). All durable data lives in Postgres. If Redis is unavailable, the backend falls through to Postgres.

The RA Adapter and FIX Engine are a separate deployment unit with their own JVM. They communicate with the Go backend exclusively through Postgres queue tables — no direct network connection between Go and Java.

Static frontend assets are served from S3 + CloudFront (or equivalent CDN).

### Local Development

```
docker-compose up   →  Postgres 16 on :5432
                    →  Backend on :8080 (go run, cache=memory)
                    →  Frontend on :5173 (Vite dev server, proxies API to :8080)
```

No Redis, no RA Adapter in local dev. The cache uses an in-memory implementation. FIX message queues can be populated with test fixtures.

## Backend Architecture

### Process Structure

The Go backend is a single process with multiple goroutines:

```
main()
├── HTTP server (Fiber app, REST endpoints, WebSocket upgrade)
├── WebSocket hub (connection registry, broadcast loop)
├── Job workers (N goroutines polling jobs table via SKIP LOCKED)
├── FIX inbound consumer (LISTEN/NOTIFY + poll on inbound_fix_messages)
├── FIX outbound producer (writes to outbound_fix_messages, triggers NOTIFY)
└── Metrics server (Prometheus /metrics on separate port)
```

All goroutines share a `*sqlx.DB` for regular database queries. A separate `*pgxpool.Pool` provides native pgx connections for LISTEN/NOTIFY listeners. Shutdown is coordinated via `context.Context` cancellation — when the process receives SIGTERM, the context is cancelled, goroutines drain, and the process exits cleanly.

### Request Flow (REST)

```
HTTP Request (with session cookie)
  → Fiber middleware chain:
      → Request ID (generate or extract X-Request-ID)
      → Structured logging (attach slog fields)
      → Session validation (cookie → cache/DB lookup → set user in Locals)
      → Metrics (request count, latency histogram)
  → Route-specific handler (e.g., OrdersHandler.Submit)
  → Validation (explicit, returns typed error)
  → Repository call (sqlx query, returns (result, error))
  → Audit write (async, non-blocking)
  → JSON response
```

Handlers are methods on service/handler structs that hold their dependencies:

```go
type OrdersHandler struct {
    repo    *repo.OrdersRepo
    db      *sqlx.DB
    cache   cache.Store
    hub     *ws.Hub
    audit   *audit.Logger
    logger  *slog.Logger
}

func (h *OrdersHandler) Submit(c *fiber.Ctx) error {
    // ...
}
```

### Request Flow (WebSocket)

```
WS Upgrade Request (session cookie sent automatically by browser)
  → Session validation (same as REST, via Fiber middleware)
  → Upgrade to WebSocket (gofiber/contrib/websocket)
  → Hub registers connection with user ID and subscription channels
  → Read loop: receives subscription changes from client
  → Write loop: receives broadcast messages from Hub, writes to client
  → On disconnect: Hub unregisters connection
```

The Hub is a central goroutine that manages all WebSocket connections:

```go
type Hub struct {
    connections map[string]map[*Conn]bool  // channel → set of connections
    broadcast   chan Message               // inbound from domain events
    register    chan *Conn
    unregister  chan *Conn
}
```

Domain events (order submitted, book updated, etc.) are sent to the Hub's broadcast channel. The Hub fans out to all connections subscribed to the relevant channel.

For multi-instance deployments, cross-instance broadcast uses Postgres LISTEN/NOTIFY on a `ws_broadcast` channel. Each instance listens and re-broadcasts to its local Hub.

### Database Access Pattern

Two database clients coexist:

- **`sqlx.DB`** (via `pgx/stdlib` driver) — all regular queries. Provides struct scanning via `db` tags.
- **`pgxpool.Pool`** (native pgx) — LISTEN/NOTIFY only. A single dedicated connection from this pool is used for notification listeners.

Domain structs double as scan targets:

```go
type Order struct {
    ID             uuid.UUID  `db:"id"              json:"id"`
    IPOID          uuid.UUID  `db:"ipo_id"          json:"ipo_id"`
    InvestorID     uuid.UUID  `db:"investor_id"     json:"investor_id"`
    Quantity       float64    `db:"quantity"         json:"quantity"`
    LimitPrice     *float64   `db:"limit_price"     json:"limit_price,omitempty"`
    Status         string     `db:"status"          json:"status"`
    IdempotencyKey uuid.UUID  `db:"idempotency_key" json:"idempotency_key"`
    CreatedAt      time.Time  `db:"created_at"      json:"created_at"`
    UpdatedAt      time.Time  `db:"updated_at"      json:"updated_at"`
}
```

Repository methods use sqlx for scanning:

```go
type OrdersRepo struct {
    db *sqlx.DB
}

const insertOrder = `
    INSERT INTO orders (
        id, ipo_id, investor_id, quantity, limit_price,
        status, idempotency_key, created_at, updated_at
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    ON CONFLICT (idempotency_key) DO NOTHING
    RETURNING *`

func (r *OrdersRepo) Insert(ctx context.Context, o *Order) (*Order, error) {
    var result Order
    err := r.db.QueryRowxContext(ctx, insertOrder,
        o.ID, o.IPOID, o.InvestorID, o.Quantity, o.LimitPrice,
        o.Status, o.IdempotencyKey, o.CreatedAt, o.UpdatedAt,
    ).StructScan(&result)
    if err != nil {
        return nil, fmt.Errorf("insert order: %w", err)
    }
    return &result, nil
}

const listOrdersByIPO = `SELECT * FROM orders WHERE ipo_id = $1 ORDER BY created_at DESC`

func (r *OrdersRepo) ListByIPO(ctx context.Context, ipoID uuid.UUID) ([]Order, error) {
    var orders []Order
    if err := r.db.SelectContext(ctx, &orders, listOrdersByIPO, ipoID); err != nil {
        return nil, fmt.Errorf("list orders by IPO %s: %w", ipoID, err)
    }
    return orders, nil
}
```

Key rules:
- SQL as `const` strings in the repo file. Backtick raw literals.
- One repo struct per aggregate root. Dependency: `*sqlx.DB`.
- Use `StructScan` for single rows, `SelectContext` for lists.
- All methods take `context.Context` as first parameter.
- `db` struct tags match Postgres column names (`snake_case`).

Transactions use `sqlx.Tx`:

```go
tx, err := r.db.BeginTxx(ctx, nil)
if err != nil {
    return fmt.Errorf("begin tx: %w", err)
}
defer tx.Rollback() // no-op if already committed

var order Order
err = tx.QueryRowxContext(ctx, insertOrder, args...).StructScan(&order)
if err != nil {
    return fmt.Errorf("insert order: %w", err)
}
// ... more work in the same transaction ...
return tx.Commit()
```

### Configuration

Configuration is loaded from environment variables. No config files.

```go
type Config struct {
    Port        int    `env:"PORT" default:"8080"`
    DBUrl       string `env:"DATABASE_URL" required:"true"`
    CacheProvider string `env:"CACHE_PROVIDER" default:"memory"` // "memory" or "redis"
    RedisURL    string `env:"REDIS_URL"`
    LogFormat   string `env:"LOG_FORMAT" default:"text"` // "text" or "json"
    // ...
}
```

Use a lightweight env parser (e.g., `caarlos0/env`) or just `os.Getenv` with manual parsing. No Viper, no complex config frameworks.

## FIX Message Queue

### Design

The Go backend and the RA Adapter communicate exclusively through two Postgres tables. There is no direct network connection between them.

```
Go Backend ←── inbound_fix_messages ←── RA Adapter ←── FIX Engine
Go Backend ──→ outbound_fix_messages ──→ RA Adapter ──→ FIX Engine
```

Both directions use the same pattern: writer inserts a row, Postgres trigger fires `pg_notify`, reader wakes up via LISTEN and claims the row with SKIP LOCKED. A periodic poll (every 5s) acts as a safety net for missed notifications.

### Queue Tables

```sql
CREATE TABLE outbound_fix_messages (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_type    TEXT NOT NULL,
    payload         JSONB NOT NULL,
    status          TEXT NOT NULL DEFAULT 'pending'
                    CHECK (status IN ('pending', 'processing', 'sent', 'failed', 'dead')),
    attempts        INT NOT NULL DEFAULT 0,
    max_attempts    INT NOT NULL DEFAULT 5,
    last_error      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    processed_at    TIMESTAMPTZ
);

CREATE INDEX idx_outbound_fix_poll ON outbound_fix_messages (created_at)
    WHERE status = 'pending';

CREATE TABLE inbound_fix_messages (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_type    TEXT NOT NULL,
    payload         JSONB NOT NULL,
    status          TEXT NOT NULL DEFAULT 'pending'
                    CHECK (status IN ('pending', 'processing', 'applied', 'failed', 'dead')),
    attempts        INT NOT NULL DEFAULT 0,
    max_attempts    INT NOT NULL DEFAULT 5,
    last_error      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    processed_at    TIMESTAMPTZ
);

CREATE INDEX idx_inbound_fix_poll ON inbound_fix_messages (created_at)
    WHERE status = 'pending';
```

### NOTIFY Triggers

```sql
CREATE OR REPLACE FUNCTION notify_queue() RETURNS trigger AS $$
BEGIN
    PERFORM pg_notify(TG_TABLE_NAME, NEW.id::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER outbound_fix_notify
    AFTER INSERT ON outbound_fix_messages
    FOR EACH ROW EXECUTE FUNCTION notify_queue();

CREATE TRIGGER inbound_fix_notify
    AFTER INSERT ON inbound_fix_messages
    FOR EACH ROW EXECUTE FUNCTION notify_queue();
```

### Consumer Pattern (Go side, inbound)

```go
func (c *InboundConsumer) Run(ctx context.Context) error {
    conn, err := c.pool.Acquire(ctx)
    if err != nil {
        return fmt.Errorf("acquire conn for LISTEN: %w", err)
    }
    defer conn.Release()

    _, err = conn.Exec(ctx, "LISTEN inbound_fix_messages")
    if err != nil {
        return fmt.Errorf("LISTEN: %w", err)
    }

    // Poll immediately on startup to catch anything queued while we were down
    c.processBatch(ctx)

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()

        default:
            // Wait for notification with timeout (poll interval as fallback)
            waitCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
            _, err := conn.Conn().WaitForNotification(waitCtx)
            cancel()

            if err != nil && !errors.Is(err, context.DeadlineExceeded) {
                c.logger.Error("notification error", "err", err)
                continue
            }
            // Either notification received or timeout — process either way
            c.processBatch(ctx)
        }
    }
}

func (c *InboundConsumer) processBatch(ctx context.Context) {
    for {
        msg, err := c.claimNext(ctx) // SELECT ... FOR UPDATE SKIP LOCKED
        if err != nil {
            if errors.Is(err, pgx.ErrNoRows) {
                return // queue drained
            }
            c.logger.Error("claim failed", "err", err)
            return
        }
        if err := c.handle(ctx, msg); err != nil {
            c.logger.Error("handle failed", "msg_id", msg.ID, "err", err)
            c.markFailed(ctx, msg, err)
        } else {
            c.markApplied(ctx, msg)
        }
    }
}
```

The key design points:
- `WaitForNotification` with a timeout doubles as the poll interval. If a NOTIFY fires, it returns immediately. If not, it returns after 5s and we poll anyway.
- On startup, `processBatch` runs before the first `WaitForNotification`, catching anything queued during downtime.
- Each notification triggers processing of ALL pending messages (batch drain), not just the one that was notified. This handles the case where multiple inserts happen in rapid succession.
- NOTIFY is transactional — it fires on commit, so the consumer never sees an uncommitted row.

## Job Queue

Background jobs (notifications, emails, webhooks, reconciliation) use the same SKIP LOCKED + LISTEN/NOTIFY pattern as the FIX queue.

### Schema

```sql
CREATE TABLE jobs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type            TEXT NOT NULL,
    payload         JSONB NOT NULL DEFAULT '{}',
    status          TEXT NOT NULL DEFAULT 'pending'
                    CHECK (status IN ('pending', 'running', 'completed', 'failed', 'dead')),
    attempts        INT NOT NULL DEFAULT 0,
    max_attempts    INT NOT NULL DEFAULT 5,
    last_error      TEXT,
    run_after       TIMESTAMPTZ NOT NULL DEFAULT now(),
    started_at      TIMESTAMPTZ,
    completed_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_jobs_poll ON jobs (run_after) WHERE status = 'pending';

CREATE TRIGGER jobs_notify
    AFTER INSERT ON jobs
    FOR EACH ROW EXECUTE FUNCTION notify_queue();
```

### Job Dispatch

Each job type implements a handler:

```go
type JobHandler interface {
    Type() string
    Handle(ctx context.Context, payload json.RawMessage) error
}
```

The worker registers handlers by type string. On claiming a job, it dispatches to the matching handler. On failure, it reschedules with exponential backoff: `run_after = now() + 2^attempts * base_delay`. After `max_attempts`, status moves to `dead`.

## Cache Layer

### Interface

```go
type Store interface {
    Get(ctx context.Context, key string) (string, error)  // returns ErrNotFound on miss
    Set(ctx context.Context, key string, value string, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
}
```

### Implementations

**`RedisStore`**: Uses `github.com/redis/go-redis/v9`. Production.

**`MemoryStore`**: `sync.RWMutex` + `map[string]entry` with TTL expiry checked on access and periodic cleanup. Dev and test. ~50 lines.

Selection is config-driven: `CACHE_PROVIDER=memory` (default) or `CACHE_PROVIDER=redis`.

### Fallback

If the cache returns an error (Redis down), callers fall through to Postgres. Cache failures are logged as warnings, never returned to the HTTP caller.

### What Gets Cached

| Key pattern | Value | TTL | Invalidation |
|-------------|-------|-----|-------------|
| `session:<id>` | Session JSON | 5 min | Explicit DELETE on logout |
| `ratelimit:<user_id>:<window>` | Counter | Window duration | Expires naturally |

Nothing else. Do not cache domain entities or query results without explicit architectural discussion.

## WebSocket Event System

### Internal Broadcast

Domain events (order submitted, book updated) are published to the Hub's broadcast channel:

```go
h.hub.Broadcast(ws.Message{
    Channel: "book:" + ipoID,
    Type:    "order.submitted",
    Payload: orderJSON,
})
```

The Hub fans out to all locally-connected clients subscribed to the matching channel.

### Cross-Instance Broadcast

For multi-instance deployments, Postgres LISTEN/NOTIFY distributes events across instances:

1. When a domain event occurs, the backend writes to a `ws_events` NOTIFY channel: `pg_notify('ws_events', json_payload)`.
2. Each backend instance listens on `ws_events` and forwards received messages to its local Hub.
3. The Hub broadcasts to local connections as normal.

This is the same LISTEN/NOTIFY infrastructure used for the FIX message queue and job queue — one pattern, reused everywhere.

## Authentication & Authorization

### Overview

Session-based authentication with HTTP-only cookies. No JWT.

Login methods:
- **Email magic link**: User enters email → server creates token, enqueues email job → user clicks link → server validates token, creates session, sets cookie.
- **OAuth (Google, Microsoft)**: Standard OAuth 2.0 authorization code flow → server handles callback, upserts user, creates session, sets cookie.

MFA is required for production but currently **out of scope**. The session model accommodates MFA state.

### Session Model

```sql
CREATE TABLE sessions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at      TIMESTAMPTZ NOT NULL,
    last_active_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    ip_address      TEXT,
    user_agent      TEXT
);

CREATE INDEX idx_sessions_user ON sessions (user_id);
CREATE INDEX idx_sessions_expires ON sessions (expires_at);
```

Session lookup uses the cache layer: check Redis/memory first, fall through to Postgres on miss, populate cache on hit. Logout deletes from both Postgres and cache — immediate invalidation across all instances.

### Cookie Configuration

```
Set-Cookie: cb_session=<session_id>; HttpOnly; Secure; SameSite=Lax; Path=/; Max-Age=<session_ttl>
```

WebSocket upgrades receive the cookie automatically (same-origin browser request).

### Authorization

Role-based, enforced in handler code. No framework.

The session middleware sets user info in Fiber's `Locals`:

```go
type AuthContext struct {
    UserID uuid.UUID
    Role   string   // "investor", "manager", "admin"
    OrgID  uuid.UUID
}

func FromCtx(c *fiber.Ctx) (AuthContext, bool) {
    ac, ok := c.Locals("auth").(AuthContext)
    return ac, ok
}
```

Handlers check authorization explicitly:

```go
auth, ok := auth.FromCtx(c)
if !ok || auth.Role != "manager" && auth.Role != "admin" {
    return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
}
```

## Observability

### Logging

Structured JSON logging via `log/slog`. A middleware attaches a request-scoped logger with request ID, user ID, and request metadata:

```go
slog.Info("order submitted",
    "request_id", requestID,
    "user_id", auth.UserID,
    "order_id", order.ID,
    "ipo_id", order.IPOID,
)
```

### Metrics

Prometheus client library. Exposed on a separate port (`:9090/metrics`) to keep it off the public ALB.

| Metric | Type | Labels |
|--------|------|--------|
| `cb_http_requests_total` | Counter | method, path, status |
| `cb_http_request_duration_seconds` | Histogram | method, path |
| `cb_ws_connections_active` | Gauge | channel_prefix |
| `cb_orders_submitted_total` | Counter | ipo_id |
| `cb_job_queue_depth` | Gauge | status |
| `cb_job_execution_duration_seconds` | Histogram | job_type, outcome |
| `cb_db_query_duration_seconds` | Histogram | query_name |
| `cb_fix_queue_depth` | Gauge | direction, status |

### Audit Trail

```sql
CREATE TABLE audit_log (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID,
    action          TEXT NOT NULL,
    entity_type     TEXT NOT NULL,
    entity_id       UUID NOT NULL,
    before_state    JSONB,
    after_state     JSONB NOT NULL,
    ip_address      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_entity ON audit_log (entity_type, entity_id);
CREATE INDEX idx_audit_user ON audit_log (user_id);
```

Audit writes are buffered through a channel and written in batches by a dedicated goroutine. The buffer has a capacity (e.g., 1000). If the buffer is full, audit writes are dropped and a metric is incremented — never block the request path.

### Health Checks

- `GET /health` — 200 if process is running.
- `GET /health/ready` — 200 if DB pool has active connections, job workers are running, and FIX consumer is connected.

## Performance Budget

| Metric | Target |
|--------|--------|
| REST API p99 latency | < 50ms |
| Order submission throughput | > 2,000 orders/sec sustained |
| WebSocket event delivery | < 200ms from DB write to client receipt |
| Concurrent WebSocket connections | > 50,000 per instance |
| Job processing throughput | > 100 jobs/sec (not on critical path) |
| FIX queue latency (NOTIFY path) | < 100ms from insert to consumer processing |
| Binary size | < 30MB |
| Container startup | < 500ms |
| Memory under load | < 256MB per instance |

These targets reflect Go's lighter resource footprint compared to JVM options. The goroutine model and pgx connection pooling should handle these comfortably on modest Fargate task sizes.
