# ClearingBid — Code Conventions

## Go Conventions

### Language Version

Go 1.22+. Always format with `gofmt` (non-negotiable — Go enforces this).

### Package Layout

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Entry point: config, wiring, startup, shutdown
├── internal/                    # All application code (not importable externally)
│   ├── api/                     # HTTP handlers and router setup
│   │   ├── router.go            # Fiber app setup, middleware chain, route registration
│   │   ├── middleware.go        # Request ID, logging, auth, metrics, recovery
│   │   ├── orders.go            # OrdersHandler struct and methods
│   │   ├── ipos.go              # IPOsHandler struct and methods
│   │   ├── auth_routes.go       # Magic link, OAuth, logout endpoints
│   │   └── health.go            # /health and /health/ready
│   ├── ws/                      # WebSocket hub and connection management
│   │   ├── hub.go               # Hub goroutine, connection registry, broadcast
│   │   ├── conn.go              # Connection wrapper, read/write loops
│   │   └── handler.go           # Fiber WebSocket upgrade handler
│   ├── domain/                  # Domain types — structs, enums, interfaces
│   │   ├── order.go             # Order, OrderStatus, NewOrderRequest
│   │   ├── ipo.go               # IPO, IPOState
│   │   ├── user.go              # User, Role
│   │   └── errors.go            # Domain error types
│   ├── repo/                    # Database repositories (sqlx queries)
│   │   ├── orders.go            # OrdersRepo: SQL consts + methods
│   │   └── ipos.go              # IPOsRepo
│   ├── cache/                   # CacheStore interface + implementations
│   │   ├── store.go             # Interface
│   │   ├── redis.go             # Redis/Valkey implementation
│   │   └── memory.go            # In-memory (dev/test)
│   ├── queue/                   # FIX message queue consumer/producer
│   │   ├── consumer.go          # LISTEN/NOTIFY + poll, SKIP LOCKED claim
│   │   └── producer.go          # Insert + implicit NOTIFY via trigger
│   ├── jobs/                    # Background job queue
│   │   ├── worker.go            # LISTEN/NOTIFY + poll, dispatch by type
│   │   ├── handler.go           # JobHandler interface
│   │   ├── send_email.go        # Email job implementation
│   │   └── send_notification.go # Notification job implementation
│   ├── auth/                    # Session management, magic link, OAuth
│   │   ├── middleware.go        # Session cookie → cache/DB → context
│   │   ├── context.go           # AuthContext, FromContext helper
│   │   ├── session.go           # Session CRUD (create, validate, destroy)
│   │   ├── magiclink.go         # Token generation, verification
│   │   └── oauth.go             # Google/Microsoft OAuth flows
│   ├── audit/                   # Audit trail
│   │   └── logger.go            # Buffered async writer (channel + goroutine)
│   ├── fix/                     # FIX domain types (shared vocabulary)
│   │   └── types.go             # NewOrderSingle, ExecutionReport, etc.
│   └── config/                  # Configuration from environment
│       └── config.go
├── tests/
│   ├── integration/             # Integration tests (testcontainers-go)
│   │   ├── orders_test.go
│   │   └── testhelpers.go       # Shared setup: Postgres container, sqlx.DB, seed
│   └── regression/              # Canary suite: full HTTP round-trips
│       ├── orders_test.go
│       └── testhelpers.go       # Shared setup: start server, HTTP client
```

### Naming

- Packages: lowercase, single word (`repo`, `cache`, `queue`). No `_` in package names.
- Files: `snake_case.go`.
- Exported types and functions: `PascalCase`.
- Unexported: `camelCase`.
- Interfaces: named by what they do, not `I`-prefix. `Store`, `Handler`, `Consumer` — not `IStore`.
- Receiver variables: short, usually first letter of type (`h` for handler, `r` for repo, `c` for consumer`).

### Error Handling

Always wrap errors with context using `fmt.Errorf` and `%w`:

```go
// GOOD
order, err := h.repo.InsertOrder(ctx, params)
if err != nil {
    return fmt.Errorf("inserting order for IPO %s: %w", ipoID, err)
}

// BAD — no context
order, err := h.repo.InsertOrder(ctx, params)
if err != nil {
    return err
}

// BAD — loses error chain (uses %v instead of %w)
if err != nil {
    return fmt.Errorf("insert failed: %v", err)
}

// FORBIDDEN — silently discarding errors
_ = h.repo.InsertOrder(ctx, params)
```

For domain errors, define sentinel errors or typed errors in `internal/domain/errors.go`:

```go
var (
    ErrNotFound     = errors.New("not found")
    ErrConflict     = errors.New("conflict")
    ErrUnauthorized = errors.New("unauthorized")
    ErrForbidden    = errors.New("forbidden")
)

type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation: %s: %s", e.Field, e.Message)
}
```

HTTP handlers translate domain errors to status codes in a helper:

```go
func writeError(c *fiber.Ctx, err error) error {
    switch {
    case errors.Is(err, domain.ErrNotFound):
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
    case errors.Is(err, domain.ErrConflict):
        return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "conflict"})
    case errors.Is(err, domain.ErrForbidden):
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
    default:
        var ve *domain.ValidationError
        if errors.As(err, &ve) {
            return c.Status(fiber.StatusBadRequest).JSON(ve)
        }
        slog.Error("unhandled error", "err", err)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
    }
}
```

### HTTP Handler Pattern

Handlers are methods on structs that hold dependencies. They return `error` — Fiber's convention:

```go
type OrdersHandler struct {
    repo    *repo.OrdersRepo
    db      *sqlx.DB
    cache   cache.Store
    hub     *ws.Hub
    audit   *audit.Logger
    logger  *slog.Logger
}

func NewOrdersHandler(r *repo.OrdersRepo, db *sqlx.DB, c cache.Store, hub *ws.Hub, al *audit.Logger, log *slog.Logger) *OrdersHandler {
    return &OrdersHandler{repo: r, db: db, cache: c, hub: hub, audit: al, logger: log}
}

func (h *OrdersHandler) Submit(c *fiber.Ctx) error {
    var req domain.SubmitOrderRequest
    if err := c.BodyParser(&req); err != nil {
        return writeError(c, &domain.ValidationError{Field: "body", Message: "invalid JSON"})
    }

    if err := req.Validate(); err != nil {
        return writeError(c, err)
    }

    auth, _ := auth.FromCtx(c)

    order, err := h.repo.Insert(c.Context(), &domain.Order{
        ID:             uuid.New(),
        IPOID:          req.IPOID,
        InvestorID:     auth.UserID,
        Quantity:       req.Quantity,
        LimitPrice:     req.LimitPrice,
        Status:         "pending",
        IdempotencyKey: req.IdempotencyKey,
        CreatedAt:      time.Now(),
        UpdatedAt:      time.Now(),
    })
    if err != nil {
        return writeError(c, fmt.Errorf("submit order: %w", err))
    }
    // Handle idempotent retry: Insert uses ON CONFLICT DO NOTHING + RETURNING
    if order == nil {
        existing, err := h.repo.FindByIdempotencyKey(c.Context(), req.IdempotencyKey)
        if err != nil {
            return writeError(c, fmt.Errorf("idempotent lookup: %w", err))
        }
        return c.Status(fiber.StatusOK).JSON(existing)
    }

    h.audit.Log(audit.Entry{
        UserID:     auth.UserID,
        Action:     "order.submitted",
        EntityType: "order",
        EntityID:   order.ID,
        AfterState: order,
    })

    h.hub.Broadcast(ws.Message{
        Channel: "book:" + order.IpoID.String(),
        Type:    "order.submitted",
        Payload: order,
    })

    return c.Status(fiber.StatusCreated).JSON(order)
}
```

**Critical Fiber rule:** `*fiber.Ctx` is only valid during the handler's execution. Never pass it to goroutines. Extract any needed values (user ID, request ID, params) into local variables before spawning background work.

### Dependency Injection

No framework. `main.go` is the composition root:

```go
func main() {
    cfg := config.Load()

    // sqlx.DB for regular queries (via pgx/stdlib driver)
    db := mustConnectDB(cfg.DBUrl)
    defer db.Close()

    // Native pgx pool for LISTEN/NOTIFY only
    pgxPool := mustConnectPgx(cfg.DBUrl)
    defer pgxPool.Close()

    cacheStore := mustCreateCache(cfg)
    hub := ws.NewHub()
    auditLog := audit.NewLogger(db)

    ordersRepo := &repo.OrdersRepo{DB: db}
    iposRepo := &repo.IPOsRepo{DB: db}

    ordersH := api.NewOrdersHandler(ordersRepo, db, cacheStore, hub, auditLog, slog.Default())
    iposH := api.NewIPOsHandler(iposRepo, db, cacheStore, hub, auditLog, slog.Default())
    authH := api.NewAuthHandler(db, cacheStore, slog.Default())

    app := api.NewApp(ordersH, iposH, authH, hub, cacheStore)

    // Start background goroutines
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer cancel()

    go hub.Run(ctx)
    go jobs.RunWorker(ctx, db, slog.Default())
    go queue.RunInboundConsumer(ctx, pgxPool, db, hub, slog.Default())

    // Start Fiber
    go func() {
        if err := app.Listen(cfg.Addr()); err != nil {
            slog.Error("server error", "err", err)
        }
    }()

    <-ctx.Done()
    app.Shutdown()
}
```

### JSON Handling

Fiber handles JSON serialization and deserialization natively:

- **Parsing request bodies:** `c.BodyParser(&req)` decodes JSON into a struct.
- **Writing responses:** `c.JSON(v)` serializes and sets `Content-Type: application/json`.
- **Status + JSON:** `c.Status(fiber.StatusCreated).JSON(order)`.

Struct tags control serialization:

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

JSON field names are `snake_case` to match Postgres column names and frontend expectations.

### Testing Patterns

**Unit tests** (pure logic, no external dependencies):

```go
func TestOrderValidation(t *testing.T) {
    tests := []struct {
        name    string
        req     domain.SubmitOrderRequest
        wantErr bool
    }{
        {name: "valid", req: validRequest(), wantErr: false},
        {name: "zero quantity", req: zeroQtyRequest(), wantErr: true},
        // ...
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.req.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

**Integration tests** (real Postgres via testcontainers-go):

```go
func TestInsertOrder(t *testing.T) {
    ctx := context.Background()
    db := setupTestDB(t) // starts Postgres container, runs schema, returns *sqlx.DB

    r := &repo.OrdersRepo{DB: db}
    order, err := r.Insert(ctx, testOrder())
    if err != nil {
        t.Fatalf("Insert: %v", err)
    }
    if order.Status != "pending" {
        t.Errorf("status = %q, want %q", order.Status, "pending")
    }
}
```

**Regression tests** (full HTTP round-trip):

```go
func TestOrderSubmissionRoundTrip(t *testing.T) {
    app := startTestApp(t) // starts full Fiber app with test DB

    req := httptest.NewRequest("POST", "/api/orders", orderBody())
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Cookie", testSessionCookie(t))

    resp, err := app.Test(req)
    if err != nil {
        t.Fatalf("POST /api/orders: %v", err)
    }
    if resp.StatusCode != fiber.StatusCreated {
        t.Errorf("status = %d, want %d", resp.StatusCode, fiber.StatusCreated)
    }
    // ... verify response body, check DB state, etc.
}
```

Key rules:
- Use stdlib `testing` package. Testify is acceptable for assertions if it reduces boilerplate.
- Table-driven tests for anything with multiple cases.
- testcontainers-go for real Postgres in integration/regression tests. No mocking the database.
- Test files live next to the code they test (`orders_test.go` next to `orders.go`) for unit tests, and in `tests/` for integration/regression.

## TypeScript / React Conventions

### File Structure

Feature-based:

```
src/
├── api/
│   ├── client.ts              # Fetch wrapper, base URL (cookies sent automatically)
│   ├── orders.ts              # Order-related API calls (typed)
│   ├── ipos.ts                # IPO-related API calls
│   └── ws.ts                  # WebSocket connection manager
├── components/
│   └── ui/                    # shadcn/ui components (owned source)
├── features/
│   ├── bidman/
│   │   ├── pages/             # Route-level components
│   │   ├── components/        # Feature-specific components
│   │   ├── hooks/             # Feature-specific hooks
│   │   └── types.ts
│   └── batman/
│       ├── pages/
│       ├── components/
│       ├── hooks/
│       └── types.ts
├── shared/
│   ├── components/            # App-level composed components (layouts, nav, guards)
│   ├── hooks/                 # useAuth, useWebSocket, etc.
│   ├── types/                 # Shared types (User, Order, IPO)
│   └── utils/                 # Formatting, validation helpers
├── routes/                    # TanStack Router route tree
│   ├── __root.tsx
│   ├── bidman/
│   └── batman/
├── App.tsx
├── main.tsx
└── routeTree.gen.ts           # Auto-generated (TanStack Router CLI)
```

### Routing (TanStack Router)

File-based routing. Route params use `$param` naming. Auth guards go in `beforeLoad` hooks. Use `loader` for data that should block navigation, TanStack Query hooks for data that can load after render. Link with `<Link to="/bidman/orders/$ipoId" params={{ ipoId }}>` — never string-interpolate.

### API Client

Typed functions in `src/api/`. Components use TanStack Query hooks:

```typescript
// src/features/bidman/hooks/useOrders.ts
export function useOrders(ipoId: string) {
  return useQuery({
    queryKey: ['orders', ipoId],
    queryFn: () => ordersApi.listByIpo(ipoId),
  });
}
```

### WebSocket

Singleton service. Session cookie sent automatically on same-origin upgrade. Reconnection with exponential backoff.

### Styling & Components

- Tailwind CSS utility classes. No CSS-in-JS.
- `cva` (class-variance-authority) for component variants.
- shadcn/ui for standard components. Radix primitives for anything shadcn doesn't cover.
- Never install MUI, Ant, Chakra, Mantine.

### Component Rules

- Functional components only.
- Props typed with explicit `interface`, not inline.
- No `any`. Use `unknown` and narrow.
- No default exports. TanStack Router route files export named `Route`.

## SQL Conventions

### Migrations

Managed via Flyway in a separate repository. If a ticket requires schema changes, document the DDL in the PR description. Follow these conventions for the DDL:

- `snake_case` for all identifiers.
- `UUID` primary keys.
- `TIMESTAMPTZ` for all time columns.
- `created_at` and `updated_at` on every table.
- Include an `idempotency_key` (UUID, UNIQUE) on any table that stores user-initiated writes.

### SQL in Repositories

SQL lives as `const` strings in the repo file, using backtick raw literals. Naming convention:

```go
const insertOrder = `INSERT INTO orders (...) VALUES (...) RETURNING *`
const findOrderByID = `SELECT * FROM orders WHERE id = $1`
const listOrdersByIPO = `SELECT * FROM orders WHERE ipo_id = $1 ORDER BY created_at DESC`
const updateOrderStatus = `UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3`
```

Use parameterized queries (`$1`, `$2`, ...). Use `RETURNING *` after INSERT/UPDATE. Prefer CTEs for complex queries.

## Git Conventions

### Branching

- `main` — always deployable, protected
- `feat/<issue>-<short-desc>` — feature branches
- `fix/<issue>-<short-desc>` — bug fix branches
- `chore/<desc>` — non-functional changes

### Commit Messages

```
<type>(<scope>): <subject>

<body — optional>

Closes #<issue>
```

Types: `feat`, `fix`, `chore`, `test`, `docs`, `refactor`
Scopes: `api`, `ws`, `repo`, `queue`, `jobs`, `cache`, `auth`, `bidman`, `batman`, `infra`

### PR Requirements

- Title matches commit convention
- Description includes: what changed, why, any tradeoffs
- Links to the GitHub issue
- All CI checks pass
- At least one human review approval
