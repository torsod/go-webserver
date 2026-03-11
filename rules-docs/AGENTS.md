# ClearingBid — Agent Instructions

## What This Is

ClearingBid is a SaaS platform for democratized IPO participation. The system manages order entry from retail and institutional investors during book-building processes, and provides lead managers with real-time visibility and control over the book.

This file is the entry point for any coding agent working on this codebase. Read it fully before starting any task.

## Architecture Summary

ClearingBid has four deployable components:

1. **Frontend** — Two SPAs (React + TypeScript + Vite):
   - **BidMan**: Investor-facing order entry and management UI
   - **BatMan**: Lead manager dashboard for book management and IPO orchestration

2. **Backend** — A single Go binary serving:
   - REST API endpoints (JSON over HTTP)
   - WebSocket connections for real-time updates
   - Background job processing (Postgres-backed, SKIP LOCKED)
   - Inbound FIX message consumer (reads from Postgres queue via LISTEN/NOTIFY + polling)

3. **RA Adapter** — A small Java process (separate repository) that bridges between ClearingBid's Postgres message queues and the Rapid Addition FIX engine. This is not part of this codebase.

4. **Database** — PostgreSQL (single primary, read replicas as needed)

See `docs/ARCHITECTURE.md` for detailed component design.

## Technology Stack

| Layer | Technology | Notes |
|-------|-----------|-------|
| Frontend framework | React | 18+ |
| Frontend language | TypeScript | 5+ |
| Frontend build | Vite | 5+ |
| Frontend styling | Tailwind CSS + shadcn/ui (Radix primitives) | |
| Frontend routing | TanStack Router | Type-safe, file-based |
| Frontend server state | TanStack Query | |
| Frontend client state | Zustand | |
| Backend language | Go | 1.22+ |
| HTTP framework | Fiber | v2 (fasthttp-based) |
| Postgres driver | pgx/stdlib + sqlx | pgx v5 (LISTEN/NOTIFY, connection pool), sqlx (struct scanning) |
| WebSocket | gofiber/contrib/websocket | |
| JSON | Fiber built-in (encoding/json under the hood) | |
| Cache (production) | Redis or Valkey (via `CacheStore` interface) | |
| Cache (dev/test) | In-memory map (behind same `CacheStore` interface) | |
| Logging | log/slog (stdlib) | JSON output in production |
| Metrics | prometheus/client_golang | |
| Testing | stdlib testing + testcontainers-go | |
| Frontend testing | Vitest + Testing Library | |
| Containerization | Docker | |
| Deployment | AWS ECS Fargate (or similar; TBD) | |

See `docs/STACK.md` for rationale and constraints.

## Hard Constraints — Do Not Violate

These are non-negotiable. If a ticket asks you to do something that conflicts with these, flag it in the PR description and do not proceed.

1. **No unnecessary dependencies, but no stdlib dogma either.** Use well-maintained third-party libraries where they reduce boilerplate or solve a problem better than stdlib (Fiber, sqlx, pgx, etc.). Justify new dependencies in the PR description, but don't write 50 lines of stdlib code to avoid a proven library.
2. **No ORMs.** No GORM, no Ent, no SQLBoiler. Handwritten SQL only, via pgx.
3. **PostgreSQL only for data.** No MongoDB, no DynamoDB. PostgreSQL stores all application data, job queues, FIX message queues, and audit logs. Redis/Valkey is permitted **only** as a read-through cache for hot-path lookups (sessions, rate limiting) — never as a data store, job queue, or pub/sub bus.
4. **All external services behind interfaces.** No code should depend directly on a specific cache provider. The cache layer uses a `CacheStore` interface with swappable implementations: Redis for production, in-memory for local development and tests. The full test suite must run locally with `docker-compose up` (Postgres only) and nothing else.
5. **Idempotency on all state-mutating order endpoints.** Every write operation must accept a client-generated idempotency key (UUID). The server must enforce uniqueness. Duplicate submissions with the same key return the original result, not an error. This is non-negotiable for a financial system.
6. **No migrations in this repo.** Schema migrations are managed in a separate repository. Do not write migration files or DDL-executing code in this codebase.
7. **Errors are values.** Use Go's error handling idiom. Never panic for recoverable errors. Never silently discard errors (`_ = someFunc()`). Always wrap errors with context: `fmt.Errorf("submitting order %s: %w", orderID, err)`.
8. **Context everywhere.** Every function that does IO (database, HTTP, cache, queue) takes a `context.Context` as its first parameter. No exceptions.

## Code Conventions

See `docs/CONVENTIONS.md` for full details. Key highlights:

### Go

- Package structure follows the standard Go project layout. See Project Structure below.
- SQL queries are `const` strings in the repository file, using backtick raw literals.
- HTTP handlers are `func(c *fiber.Ctx) error`. Fiber provides routing, middleware, and context handling.
- All database access returns `(result, error)` — no wrapper types, no monads.
- SQL queries are `const` strings in the repository file, using backtick raw literals. Row scanning uses `sqlx` struct tags (`db:"column_name"`) — no manual row scanning.
- Dependency injection is explicit: constructors take interfaces, `main()` wires everything.
- Concurrency via goroutines and channels. No external concurrency frameworks.

### TypeScript / React

- Functional components only. No class components.
- State management: TanStack Query for server state, Zustand for client state. No Redux.
- Routing: TanStack Router. Type-safe route definitions, file-based — not React Router.
- API calls go through a typed client layer (`src/api/`). Never call `fetch` directly from components.
- UI components: shadcn/ui built on Radix primitives with Tailwind styling. Do not install alternative component libraries (MUI, Ant, Chakra, Mantine).
- WebSocket connections managed via a singleton hook/service, not per-component.
- File structure: feature-based (`src/features/bidman/`, `src/features/batman/`) not type-based.
- Components: Emphasise reusable components. Do not inline duplicate Tailwind styles -- extract reusable components instead.

### SQL

- Use `snake_case` for all database identifiers.
- All tables have `id` (UUID, PK), `created_at` (timestamptz), `updated_at` (timestamptz).
- Use `timestamptz` everywhere, never `timestamp`.
- **Migrations are managed in a separate repository.** If a ticket requires schema changes, document the required DDL in the PR description and flag it as needing a companion migration.

## Project Structure

```
clearingbid/
├── AGENTS.md                    # This file
├── docs/
│   ├── ARCHITECTURE.md          # Detailed architecture and component design
│   ├── STACK.md                 # Technology choices and rationale
│   ├── CONVENTIONS.md           # Full code conventions
│   ├── TICKETS.md               # How tickets work, PR workflow
│   └── spec/                    # Product specification (the domain bible)
├── backend/
│   ├── go.mod
│   ├── go.sum
│   ├── cmd/
│   │   └── server/
│   │       └── main.go          # Entry point, wiring, startup
│   ├── internal/
│   │   ├── api/                 # HTTP handlers and route definitions
│   │   │   ├── router.go        # Fiber app setup, middleware chain
│   │   │   ├── orders.go        # Order-related handlers
│   │   │   └── ipos.go          # IPO-related handlers
│   │   ├── ws/                  # WebSocket handlers and hub
│   │   │   ├── hub.go           # Connection registry, broadcast
│   │   │   └── handler.go       # Upgrade handler, read/write loops
│   │   ├── domain/              # Domain types (structs, enums, interfaces)
│   │   │   ├── order.go
│   │   │   ├── ipo.go
│   │   │   └── user.go
│   │   ├── repo/                # Database repositories (pgx queries)
│   │   │   ├── orders.go
│   │   │   └── ipos.go
│   │   ├── cache/               # CacheStore interface + implementations
│   │   │   ├── store.go         # Interface definition
│   │   │   ├── redis.go         # Redis implementation
│   │   │   └── memory.go        # In-memory implementation (dev/test)
│   │   ├── queue/               # FIX message queue (Postgres + LISTEN/NOTIFY)
│   │   │   ├── consumer.go      # Inbound message consumer
│   │   │   └── producer.go      # Outbound message producer
│   │   ├── jobs/                # Background job queue (SKIP LOCKED)
│   │   │   ├── worker.go        # Job polling and dispatch
│   │   │   └── handlers.go      # Job type implementations
│   │   ├── auth/                # Session, magic link, OAuth handlers
│   │   │   ├── middleware.go     # Session validation middleware
│   │   │   ├── magiclink.go
│   │   │   └── oauth.go
│   │   ├── audit/               # Audit trail writer
│   │   │   └── logger.go
│   │   ├── fix/                 # FIX domain types (shared vocabulary with RA adapter)
│   │   │   └── types.go         # NewOrderSingle, ExecutionReport, etc.
│   │   └── config/              # Configuration loading
│   │       └── config.go
│   └── tests/
│       ├── integration/         # Integration tests (Testcontainers)
│       └── regression/          # Canary regression suite (full stack)
├── frontend/
│   ├── package.json
│   ├── vite.config.ts
│   ├── components.json          # shadcn/ui configuration
│   ├── src/
│   │   ├── api/                 # Typed API client, WebSocket service
│   │   ├── components/
│   │   │   └── ui/              # shadcn/ui components (owned source)
│   │   ├── features/
│   │   │   ├── bidman/          # BidMan UI components, pages, hooks
│   │   │   └── batman/          # BatMan UI components, pages, hooks
│   │   ├── routes/              # TanStack Router route tree
│   │   ├── shared/              # Shared components, types, utilities
│   │   └── App.tsx
│   └── tsconfig.json
├── docker/
│   ├── Dockerfile.backend
│   ├── Dockerfile.frontend
│   └── docker-compose.yml       # Local dev stack (app + Postgres)
└── .github/
    └── workflows/               # CI pipeline
```

## Domain Specification

The product specification lives in `docs/spec/`. This is the authoritative source for business rules, entity definitions, order lifecycle, IPO states, and all domain logic.

**Before implementing any ticket that involves business logic, read the relevant sections of the spec.** The spec takes precedence over any assumptions in these architectural documents.

## FIX Integration (Phased)

FIX connectivity is **not being integrated yet.** The current phase focuses on building the core platform.

**Phase 1 (current):** No FIX. The Go backend reads from and writes to Postgres queue tables (`inbound_orders`, `outbound_orders`). In development, these tables are populated by test fixtures or a stub script. The queue consumer and producer are built and tested against real Postgres.

**Phase 2 (planned):** RA Adapter — a small Java process (separate repo) that reads outbound orders from Postgres, translates to FIX via RA's plugin API, and writes inbound FIX messages back to Postgres. Communication is entirely through the Postgres queue tables using LISTEN/NOTIFY + SKIP LOCKED polling.

**Phase 3 (production):** RA Adapter connected to live FIX counterparties.

Agents should implement against the queue tables. Do not write FIX protocol handling code, RA-specific code, or any Java code in this repo.

## Observability

All backend code must be observable from day one. This is part of every feature, not optional infrastructure.

### Structured Logging

- Use `log/slog` (Go stdlib). JSON output in production, text in development.
- Every log line must include: request ID, user ID (if authenticated), and relevant entity IDs.
- Log at appropriate levels: `slog.Error` for unexpected failures, `slog.Warn` for recoverable issues, `slog.Info` for significant business events, `slog.Debug` for diagnostics.
- Never use `fmt.Println` or `log.Println`. Always use `slog`.
- Add structured fields, don't interpolate into message strings: `slog.Info("order submitted", "order_id", orderID, "ipo_id", ipoID)`.

### Metrics

- Use `prometheus/client_golang`. Expose `/metrics` endpoint.
- Mandatory metrics for every new feature: request count/latency by endpoint, active WebSocket connections, job queue depth, error rates, DB query latency.

### Audit Trail

- All state-mutating operations produce an audit record (append-only `audit_log` table).
- Audit writes are async (goroutine + channel buffer) and must not fail the parent operation.

### Health Checks

- `GET /health` — liveness (process is running).
- `GET /health/ready` — readiness (DB connected, job workers running).

## Working With Tickets

See `docs/TICKETS.md` for the full workflow. Summary:

1. Each GitHub issue contains a clear description, acceptance criteria, and affected areas.
2. Agent reads the ticket + this AGENTS.md + relevant docs.
3. Agent creates a feature branch: `feat/<issue-number>-<short-description>`.
4. Agent implements, writes tests, runs them locally.
5. Agent pushes and opens a PR referencing the issue (`Closes #<N>`).
6. PR is reviewed by a human developer. Feedback is given as PR comments.
7. Agent addresses review comments on the same branch and pushes again.

## What Agents Must Do Before Submitting a PR

- [ ] Code compiles with no errors (`go build ./...`).
- [ ] All existing tests pass (`go test ./...`) — no regressions.
- [ ] `go vet ./...` reports no issues.
- [ ] Frontend builds cleanly (`npm run build`, `npm run lint`, `npm run typecheck`).
- [ ] New code has unit tests covering the new functionality.
- [ ] DB-touching code has integration tests using testcontainers-go.
- [ ] State-mutating endpoints have tests verifying idempotency behavior.
- [ ] New/modified endpoints have a corresponding regression test in `tests/regression/`.
- [ ] Observable: logging, metrics, and audit trail per Observability section.
- [ ] If schema changes are needed, document the required DDL in the PR description.
- [ ] No new dependencies added without justification in the PR description.
- [ ] No violations of the Hard Constraints listed above.
- [ ] PR description includes: what changed, why, and any open questions or tradeoffs.
