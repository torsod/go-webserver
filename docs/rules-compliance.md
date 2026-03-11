Context
The rules-docs/ folder contains 5 markdown files (AGENTS.md, ARCHITECTURE.md, CONVENTIONS.md, STACK.md, TICKETS.md) that define how the ClearingBid platform should be built. The current go-webserver was ported from a Meteor.js/TypeScript app and does not yet follow these rules. This plan identifies every compliance gap and the specific changes needed.

Phase 0: Hard Constraint Violations (Must Fix First)
AGENTS.md defines non-negotiable "Hard Constraints." Three are currently violated.
0.1 Remove Migrations from This Repo
Rule: "No migrations in this repo — managed in separate repository via Flyway."
Current: migrations/ directory has 14 SQL files. internal/store/db.go has a RunMigrations() function. cmd/server/main.go calls it on startup.
Changes:

Delete migrations/ directory entirely
Remove RunMigrations() from internal/store/db.go
Remove migration call from cmd/server/main.go
Save existing DDL into docs/schema-ddl.md for reference when creating the separate migrations repo

0.2 Remove FIX Protocol Engine Code
Rule: "Phase 1 (current): No FIX live. Backend reads/writes Postgres queue tables only. Do not write FIX protocol code in this repo."
Current: Full quickfixgo engine — internal/fixengine/ (3 files), internal/service/fix_service.go, internal/handler/fix_handler.go, internal/store/fix_store.go, internal/scheduler/fix_log_cleanup.go, FIX types in internal/domain/fix.go, plus wiring in main.go.
Changes:

Delete internal/fixengine/ (engine.go, real_engine.go, simulated_engine.go)
Delete internal/service/fix_service.go
Delete internal/handler/fix_handler.go
Delete internal/store/fix_store.go
Delete internal/scheduler/fix_log_cleanup.go
Remove FIX store interfaces from internal/store/interfaces.go
Strip internal/domain/fix.go down to queue-relevant types only
Remove all FIX wiring from cmd/server/main.go
Remove github.com/quickfixgo/quickfix and github.com/shopspring/decimal from go.mod
Replacement: Build internal/queue/ (Phase 3.5) with Postgres queue tables + LISTEN/NOTIFY

0.3 Add Idempotency on State-Mutating Endpoints
Rule: "Every write operation must accept a client-generated idempotency key (UUID). Duplicate submissions return the original result."
Current: No idempotency key on any write endpoint.
Changes:

Add IdempotencyKey field to domain.Order, domain.Offering, and other write entities
Add idempotency_key to INSERT SQL with ON CONFLICT (idempotency_key) DO NOTHING
When INSERT returns no rows (conflict), look up existing record by key and return it
Document DDL (ALTER TABLE ... ADD COLUMN idempotency_key UUID UNIQUE) for migrations repo


Phase 1: Framework Migration
1.1 HTTP Framework: chi → Fiber v2
Rule: AGENTS.md/STACK.md require Fiber v2 (fasthttp-based).
Current: Uses go-chi/chi/v5 with net/http handler signatures.
Changes (every handler file + main.go):

All handler signatures: func(w http.ResponseWriter, r *http.Request) → func(c *fiber.Ctx) error
Body parsing: json.NewDecoder(r.Body).Decode(&req) → c.BodyParser(&req)
Responses: writeJSON(w, status, data) → c.Status(status).JSON(data)
URL params: chi.URLParam(r, "id") → c.Params("id")
Context: r.Context() → c.UserContext()
Dependencies: remove go-chi/chi/v5, add github.com/gofiber/fiber/v2
WebSocket: gorilla/websocket → gofiber/contrib/websocket

Files: cmd/server/main.go, all 11 files in internal/handler/
1.2 DB Layer: pgxpool → sqlx + pgx/stdlib
Rule: "sqlx.DB via pgx/stdlib for regular queries. pgxpool only for LISTEN/NOTIFY."
Current: All stores use *pgxpool.Pool with manual row.Scan() and custom scanXxx() functions.
Changes:

internal/store/db.go — create *sqlx.DB via pgx/stdlib, keep *pgxpool.Pool for LISTEN/NOTIFY only
All store files — replace manual scanning with sqlx.StructScan/SelectContext
Remove all scanXxx() helper functions (sqlx uses db: struct tags automatically)
Add github.com/jmoiron/sqlx to go.mod

Files: internal/store/db.go, all 6 store implementation files
1.3 SQL as Const Strings
Rule: "SQL queries as const strings using backtick raw literals."
Current: SQL is inline within methods or built via string concatenation.
Changes: Extract all inline SQL into named const blocks at top of each store file:
goconst insertOrder = `INSERT INTO orders (...) VALUES (...) RETURNING *`
const findOrderByID = `SELECT * FROM orders WHERE id = $1`
1.4 Package Rename: handler/ → api/, store/ → repo/
Rule: CONVENTIONS.md layout: api/ for handlers, repo/ for repositories.
Changes:

Rename internal/handler/ → internal/api/
Rename internal/store/ → internal/repo/
Update package declarations and all imports in main.go, service files


Phase 2: Convention Alignment
2.1 JSON Tags: camelCase → snake_case
Rule: "JSON field names are snake_case matching Postgres column names."
Current: All domain structs use camelCase (e.g., json:"orderType", json:"userId").
Changes: Update every json: tag in all domain files:

internal/domain/order.go (~30 tags)
internal/domain/offering.go (~60 tags)
internal/domain/user.go (~20 tags)
internal/domain/trade.go (~25 tags)
internal/domain/settings.go (~20 tags)
internal/domain/orderbook.go (~10 tags)
internal/domain/settlement.go (~15 tags)

Also update static/manage.html JavaScript and any client code referencing camelCase API fields.
Warning: This is a breaking API change for any existing clients.
2.2 Domain Error Types
Rule: CONVENTIONS.md requires domain/errors.go with sentinel errors and typed errors.
Current: No custom error types. All errors are plain fmt.Errorf strings.
Changes:

Create internal/domain/errors.go with sentinel errors: ErrNotFound, ErrConflict, ErrUnauthorized, ErrForbidden
Add ValidationError type with Error() method
Update service layer to return domain errors
Update handler writeError to translate domain errors → HTTP status codes via errors.Is/errors.As

2.3 Config Expansion and Validation
Current: Minimal config: Port, DatabaseURL, JWTSecret, FIXSimulated. No validation.
Changes:

Add: CacheProvider, RedisURL, LogFormat, MetricsPort, SessionTTL
Remove: JWTSecret (sessions replace JWT), FIXSimulated (no FIX engine)
Add validation (required fields, valid enum values)


Phase 3: New Subsystems
3.1 Auth: Session-Based with Cookies (Replace JWT)
Rule: Session-based with HTTP-only cookies, no JWT. Cache-first session lookup.
Current: Bearer token auth where token = username.
New files:

internal/auth/middleware.go — session cookie validation
internal/auth/context.go — AuthContext{UserID, Role, OrgID} + FromCtx()
internal/auth/session.go — session CRUD, cache-first DB-fallback lookup

DDL to document: sessions table (id, user_id, expires_at, last_active_at, ip_address, user_agent)
3.2 Cache Layer (NEW)
Rule: Store interface with Get/Set/Delete. MemoryStore for dev, RedisStore for prod. App must work with dead cache.
New files:

internal/cache/store.go — interface + ErrNotFound
internal/cache/memory.go — sync.RWMutex + map[string]entry with TTL
internal/cache/redis.go — go-redis/v9 wrapper

3.3 Audit Trail (NEW)
Rule: Append-only audit_log table. Async writes via channel + goroutine. Never blocks request path.
New files:

internal/audit/logger.go — buffered channel writer, batch insert goroutine

DDL to document: audit_log table (id, user_id, action, entity_type, entity_id, before_state JSONB, after_state JSONB, ip_address, created_at)
3.4 Job Queue (NEW)
Rule: jobs table with SKIP LOCKED + LISTEN/NOTIFY. JobHandler interface. Exponential backoff. Dead letter.
New files:

internal/jobs/worker.go — LISTEN + poll loop, claim via SKIP LOCKED
internal/jobs/handler.go — JobHandler interface: Type() + Handle(ctx, payload)

3.5 FIX Message Queue via Postgres (Replaces fixengine)
Rule: outbound_fix_messages and inbound_fix_messages tables. LISTEN/NOTIFY + SKIP LOCKED.
New files:

internal/queue/consumer.go — inbound message consumer
internal/queue/producer.go — outbound message writer
internal/fix/types.go — FIX domain types only (no protocol code)


Phase 4: Observability
4.1 Request ID Middleware
Rule: Generate or extract X-Request-ID on every request. Propagate through context.
4.2 Request-Scoped Structured Logger
Rule: Every log line includes request_id, user_id. JSON in production, text in dev.
Current: Global slog with text handler only.
Changes:

main.go selects JSON/text handler based on Config.LogFormat
Middleware creates request-scoped *slog.Logger with request_id + user_id
Pass logger through to service/repo layers

4.3 Prometheus Metrics
Rule: Prometheus on :9090. Required: request count/latency, WS connections, job queue depth, DB query latency.
Changes:

Add github.com/prometheus/client_golang
Fiber middleware for HTTP metrics
Separate HTTP server on :9090 serving /metrics

4.4 Health Check Redesign
Rule: GET /health (liveness) + GET /health/ready (readiness: DB pool, workers).
Current: Single GET /api/health with uptime.

Phase 5: WebSocket Redesign
5.1 Channel-Based Hub with Cross-Instance Support
Rule: Hub subscriptions by channel ("book:ipoID"). Cross-instance via pg_notify('ws_events').
Current: Hub broadcasts to ALL clients. No channels. No cross-instance.
Changes:

Redesign Hub: connections map[string]map[*Conn]bool (channel → connections)
Add subscribe/unsubscribe from client messages
Broadcast fans out only to matching channel subscribers
Add pg_notify for cross-instance delivery
Hub.Run takes context.Context for shutdown
Switch to gofiber/contrib/websocket

Files: internal/ws/hub.go, internal/ws/client.go

Phase 6: Testing (Currently Zero Tests)
6.1 Unit Tests — table-driven, *_test.go next to code

Domain validation (order price/qty, state transitions)
Service business logic (allocation, orderbook)
Config loading, cache implementations

6.2 Integration Tests — testcontainers-go, real Postgres

Add github.com/testcontainers/testcontainers-go
Test helpers for Postgres container lifecycle
Tests for each repo method

6.3 Regression Tests — tests/regression/, full HTTP round-trips

End-to-end API tests using Fiber's app.Test()
Happy path + critical error paths for all endpoints


Gap Summary
PhaseDescriptionGap Count0Hard Constraint Violations31Framework Migration42Convention Alignment33New Subsystems54Observability45WebSocket Redesign16Testing3Total23
Verification
After each phase:

go build ./... — compiles
go vet ./... — no issues
go test ./... — all tests pass (Phase 6+)
Manual test via manage.html and curl — API still functional
go mod tidy — no unused dependencies
