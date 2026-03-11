# ClearingBid — Technology Stack

This document records stack decisions and their rationale. Agents should understand not just _what_ to use, but _why_, and make consistent choices when encountering ambiguity.

## Backend: Go

**Decision:** Go 1.22+.

**Rationale:**
- Async is the default, not an opt-in mode. Goroutines handle concurrency at the language level. There is no "blocking vs non-blocking" decision to get wrong — every function call is effectively non-blocking from the runtime's perspective.
- Agent training data on Go is strong. The language is simple enough that agents produce consistent, idiomatic code. There's one way to handle errors (`if err != nil`), one way to do concurrency (goroutines + channels), one way to write an HTTP handler. This rigidity is a feature when agents are writing the code.
- Compilation is sub-second. The binary is a single static file (~20MB). Docker images can be `FROM scratch`. Startup is instant. Memory footprint under load is a fraction of any JVM option.
- The dependency surface is intentionally small but not dogmatically stdlib-only. We use Fiber (HTTP), sqlx + pgx (Postgres), and targeted libraries where they reduce boilerplate. The goal is fewer moving parts, not fewer imports.
- The JVM is no longer needed for the main backend. The only Java in the system is the RA Adapter, which is a separate deployment with its own lifecycle.

**Tradeoffs acknowledged:**
- Go's type system is less expressive than Java's (no generics on methods, no sealed types, no algebraic data types). Domain modeling is more verbose.
- Fewer financial-domain libraries than Java. Acceptable — our domain logic is custom, and RA handles FIX protocol details.
- Error handling is verbose. Acceptable — explicit is better than implicit, and agents handle `if err != nil` correctly without deep framework knowledge.

## HTTP Framework: Fiber v2

**Decision:** Fiber for HTTP routing, middleware, and request handling.

**Rationale:**
- Fiber provides routing, middleware chaining, body parsing, and response helpers in one package. Less boilerplate than assembling these from stdlib.
- `*fiber.Ctx` provides a clean request/response API: `c.Params()`, `c.BodyParser()`, `c.JSON()`, `c.Status()`. No juggling `http.ResponseWriter` and `*http.Request` separately.
- High performance via fasthttp. For a spiky-demand system, fasthttp's zero-allocation design and connection reuse are genuine advantages.
- Excellent agent training data — Fiber is one of the most popular Go web frameworks.
- Built-in WebSocket support via `contrib/websocket`, which means we don't need a separate WebSocket library.

**Tradeoffs acknowledged:**
- Fiber uses fasthttp, not `net/http`. This means `net/http` middleware doesn't compose directly. In practice, this rarely matters — Fiber's own middleware API covers the same ground, and we're not pulling in a large ecosystem of `net/http` middleware.
- `*fiber.Ctx` is request-scoped and must not be held beyond the handler's return. If you need values from the context in a goroutine, extract them first. Agents must not pass `*fiber.Ctx` to background goroutines.

## Postgres: pgx + sqlx

**Decision:** Two Postgres clients, each for what it does best.

- **`sqlx`** (`jmoiron/sqlx`) via the `pgx/stdlib` driver — all regular queries. Provides `StructScan`, `SelectContext`, `NamedExec` for scanning rows directly into Go structs using `db` struct tags.
- **`pgx`** (native `pgxpool.Pool`) — LISTEN/NOTIFY only. A dedicated connection from this pool runs the notification listeners for the FIX message queue, job queue, and WebSocket broadcast.

**Rationale:**
- sqlx gives you struct scanning without code generation or schema synchronization. Write a struct with `db` tags, write SQL, call `StructScan`. No generated files, no build step, no running database required to write code.
- pgx provides LISTEN/NOTIFY support that `database/sql` (and therefore sqlx) cannot — `WaitForNotification` requires a native pgx connection.
- The `pgx/stdlib` driver means sqlx gets pgx's performance and type support under the hood. It's not a separate driver — sqlx talks to pgx through the `database/sql` interface.

**Why not an ORM (GORM, Ent, SQLBoiler):**
- Same argument as the Java version: ORMs hide SQL, add magic, and create agent bug factories. Handwritten SQL is explicit, testable, and debuggable.

**Why not sqlc:**
- sqlc requires maintaining a `schema.sql` file that mirrors the production schema — effectively a second source of truth alongside Flyway. When the schema changes, you'd need to update both the migration and the sqlc schema before you can write queries against the new tables. sqlx has no such requirement.
- sqlc's code generation adds a build step. sqlx is just a library — import it and use it.

**Why not raw pgx everywhere:**
- pgx's native query interface requires manual row scanning (`rows.Scan(&field1, &field2, ...)`), which is tedious and error-prone for structs with many fields. sqlx's `StructScan` eliminates this boilerplate.

## WebSocket: Fiber contrib/websocket

**Decision:** Fiber's WebSocket contrib package (`github.com/gofiber/contrib/websocket`) for WebSocket support.

**Rationale:**
- Integrates directly with Fiber's routing and middleware chain. WebSocket upgrade is a Fiber handler, so auth middleware runs before the upgrade — no separate upgrade path.
- Built on `gorilla/websocket` internally but managed through Fiber's lifecycle.
- Keeps the dependency count low — one HTTP framework handles both REST and WebSocket.

## Frontend: React + TypeScript + Vite

**Decision:** React 18+, TypeScript 5+, Vite 5+, Tailwind CSS, shadcn/ui + Radix, TanStack Router, TanStack Query.

**Rationale:** (unchanged from previous iteration)
- React has the deepest agent training data. TypeScript catches runtime errors at compile time. Vite is the standard build tool.
- shadcn/ui provides copy-paste-and-own component primitives built on Radix accessibility primitives. Components live in the codebase, not `node_modules`.
- TanStack Router for end-to-end type-safe routing. TanStack Query for server state management.
- Zustand for client-side UI state. No Redux.
- See docs from previous iteration for detailed frontend rationale — the frontend stack is identical.

## Database: PostgreSQL 16+

**Decision:** PostgreSQL. Single source of truth for everything: application data, job queue, FIX message queues, audit log, sessions.

**Rationale:**
- Serves as both database and message broker (via LISTEN/NOTIFY + SKIP LOCKED). This eliminates the need for a separate message queue (SQS, Kafka, RabbitMQ) for the FIX integration boundary.
- JSONB for flexible payloads (job payloads, FIX message payloads, audit before/after state).
- LISTEN/NOTIFY is transactional — notifications fire on commit, so consumers never see uncommitted data.
- One datastore to operate, back up, and monitor.

**Constraints:**
- All timestamps: `TIMESTAMPTZ`. No exceptions.
- All primary keys: `UUID`.
- All tables include `created_at` and `updated_at`.
- Column naming: `snake_case`.
- No stored procedures for business logic.

## Cache: Redis/Valkey (Production) + In-Memory (Dev/Test)

**Decision:** Abstract cache behind a `Store` interface. Redis (or Valkey/ElastiCache) in production; in-memory Go map for development and tests.

**Rationale:**
- Session validation runs on every request. Under spiky load, a cache eliminates DB pressure for repeated lookups of the same session.
- Redis provides cross-instance consistency: a revoked session is immediately invisible to all instances. In-memory caches per instance create a coherence window.
- The `Store` interface ensures the full test suite runs locally without Redis. `go test ./...` works with nothing but a Postgres container.

**Hard rule:** The application must function correctly with a dead cache. Misses and errors fall through to Postgres.

## Database Migrations: Flyway (Separate Repo)

**Decision:** Hosted in a separate repository, executed as a separate devops process.

**Rationale:**
- Decouples schema changes from application deployments.
- Migrations handle advisory locking, partial failure recovery, and concurrent execution.
- The Go codebase has no migration logic, no DDL, no schema management code.

## FIX Integration: Postgres Queue + RA Adapter

**Decision:** The Go backend and the RA FIX engine communicate through Postgres queue tables. No gRPC, no REST, no message broker.

**Rationale:**
- The process boundary between the Go backend and the RA JVM already exists — RA plugins run in RA's runtime, not ours. A network boundary is unavoidable.
- Postgres as the transport means both sides need only a database connection. The RA Adapter is a thin Java process that reads/writes Postgres tables and calls RA's plugin API. No additional infrastructure.
- LISTEN/NOTIFY provides near-instant event delivery (sub-100ms). Periodic polling (every 5s) provides a safety net for missed notifications.
- IPO book-building has latency tolerance measured in seconds, not milliseconds. A few seconds of propagation through the queue is irrelevant to the business process.
- The queue tables are inspectable with SQL. Debugging is `SELECT * FROM outbound_fix_messages WHERE status = 'failed'`. No log diving through Kafka offsets or SQS receipt handles.

**Tradeoffs acknowledged:**
- Postgres LISTEN/NOTIFY has no guaranteed delivery — if nobody is listening, the notification is lost. Mitigated by SKIP LOCKED polling as a fallback.
- Queue throughput is limited by Postgres write performance. At our scale (hundreds to low thousands of messages per second during peak), this is not a concern.

## Testing

**Backend:**
- Go stdlib `testing` package. Testify is acceptable for assertions if it reduces boilerplate.
- testcontainers-go for integration tests (real Postgres in Docker).
- Table-driven tests for handler logic.
- Regression tests in `tests/regression/` — full HTTP round-trip tests against a real server.

**Frontend:**
- Vitest for unit tests.
- Testing Library for component tests.
- MSW (Mock Service Worker) for API mocking in tests.

## Observability: slog + Prometheus

**Decision:** `log/slog` for structured logging, `prometheus/client_golang` for metrics, Postgres for audit trail.

**Rationale:**
- `slog` is in the Go standard library as of Go 1.21. No external logging dependency needed. JSON handler for production, text handler for development.
- Prometheus is the most widely supported metrics format. Integrates with Grafana, CloudWatch, Datadog.
- Audit trail to Postgres keeps everything queryable with SQL.

## Deployment: Docker + ECS Fargate

**Decision:** Docker containers on AWS ECS Fargate. No Kubernetes.

**Rationale:**
- Go binaries produce tiny Docker images (~20-30MB). Startup is instant.
- ECS Fargate provides auto-scaling without Kubernetes overhead.
- Target-tracking auto-scaling on CPU handles the spiky demand profile.
- The RA Adapter is a separate ECS service (or EC2 instance, depending on RA's deployment model).

**Go build for containers:**

```dockerfile
FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

FROM scratch
COPY --from=builder /server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
```
