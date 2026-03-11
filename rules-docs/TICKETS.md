# ClearingBid — Ticket Workflow

## Ticket Structure

Every GitHub issue follows this template:

```markdown
## Summary
[One-paragraph description of what needs to be done]

## Affected Areas
- [ ] Backend (API/handlers)
- [ ] Backend (DB/repos)
- [ ] Backend (Jobs)
- [ ] Backend (WebSocket)
- [ ] Backend (FIX queue)
- [ ] Frontend (BidMan)
- [ ] Frontend (BatMan)
- [ ] Frontend (shared)
- [ ] Infrastructure (Docker/CI)
- [ ] Schema change needed (requires migration in migrations repo)

## Acceptance Criteria
- [ ] [Specific, testable criterion]
- [ ] [Another criterion]

## Technical Notes
[Optional: hints about implementation approach, references to existing code, relevant docs sections]

## Subtasks
- [ ] [If the ticket is large, break it into ordered subtasks]
```

## Agent Workflow

When an agent receives a ticket, it follows this sequence:

### 1. Read Context

Before writing any code:
- Read `AGENTS.md` (this is the primary context document)
- Read the relevant `docs/` files based on the affected areas
- Read any files referenced in the ticket's Technical Notes
- Look at existing code in the affected packages to understand current patterns

### 2. Plan

For non-trivial tickets, write a brief implementation plan as a comment on the issue before starting. This is especially important for:
- Database schema changes (which cascade into repo queries, domain types, handlers, and frontend)
- New API endpoints (which require handler, route registration, repo methods, and frontend integration)
- Cross-cutting changes (auth, WebSocket protocol)

### 3. Branch

```bash
git checkout main
git pull origin main
git checkout -b feat/<issue-number>-<short-description>
```

### 4. Implement

Follow the order that makes the code testable at each step:

**For backend-only changes:**
1. Domain types (if new entities or request/response types)
2. Repository (SQL consts + sqlx methods in `internal/repo/`)
3. Handler methods
4. Route registration
5. Unit tests
6. Integration tests (testcontainers-go)
7. Regression test if new/modified endpoint

**For frontend-only changes:**
1. Types (if new API response shapes)
2. API client functions
3. Components / hooks
4. Page-level integration
5. Tests

**For full-stack changes:**
1. Backend: domain types → repo methods → handler → route
2. Backend tests
3. Frontend: types → API client → components → pages
4. Frontend tests

If the change requires schema modifications, document the required DDL in the PR description under a "Required Schema Changes" heading.

### 5. Test Locally

```bash
# Backend
cd backend
go build ./...        # Compile
go vet ./...          # Static analysis
go test ./...         # All tests (unit + integration via testcontainers)

# Frontend
cd frontend
npm run lint          # ESLint
npm run typecheck     # tsc --noEmit
npm run test          # Vitest
npm run build         # Production build
```

All of the above must pass before pushing.

### 6. Commit and Push

```bash
git add -A
git commit -m "feat(api): implement order submission endpoint

- Add Order domain types and SubmitOrderRequest
- Add OrdersRepo with sqlx queries
- Add OrdersHandler.Submit with idempotency
- Integration tests with testcontainers-go
- Regression test for POST /api/orders

Closes #42"

git push -u origin feat/42-order-submission
```

### 7. Open PR

PR description template:

```markdown
## What
[Brief description of the change]

## Why
[Link to issue, business context]

## How
[Key implementation decisions, any tradeoffs]

## Required Schema Changes
[If applicable: DDL needed in the migrations repo. If none, delete this section.]

## Testing
[What was tested, how to verify]

## Checklist
- [ ] `go build ./...` succeeds
- [ ] `go vet ./...` clean
- [ ] `go test ./...` passes (no regressions)
- [ ] Frontend builds cleanly (if applicable)
- [ ] New functionality has unit tests
- [ ] DB-touching code has integration tests (testcontainers-go)
- [ ] New/modified endpoints have regression tests
- [ ] State-mutating endpoints enforce idempotency
- [ ] Logging, metrics, and audit trail included
- [ ] No new dependencies without justification
- [ ] No Hard Constraint violations
```

### 8. Address Review Feedback

When a reviewer leaves comments:
- Read all comments before making changes
- Address each comment (fix it, or explain why not)
- Push new commits to the same branch (do not force-push)
- Reply to each review comment indicating how it was addressed

## Testing Requirements

### Unit Tests

Every ticket that adds or modifies functionality must include unit tests. For Go, use table-driven tests for any function with multiple input/output cases. Tests go next to the code: `orders_test.go` alongside `orders.go`.

### Integration Tests

Any code that touches the database must have integration tests using testcontainers-go (real Postgres, not mocks). These tests verify SQL queries, transactions, constraints, and the LISTEN/NOTIFY + SKIP LOCKED patterns.

### Regression Suite

The `tests/regression/` directory contains end-to-end tests that exercise API endpoints through the full stack (HTTP request → handler → repo → DB → response). These tests:

- Start the full Go server with a testcontainers Postgres instance.
- Exercise each endpoint's happy path and critical error paths.
- Are run as a **canary suite during deployments** — a failing regression test blocks the deployment.

When a ticket adds a new endpoint or modifies an existing one, the agent must add or update the corresponding regression test.

## Ticket Sizing

Tickets should be scoped so an agent can complete them in a single session:

- **Small** (1-3 files changed): Bug fix, add a field, simple UI change
- **Medium** (4-8 files changed): New API endpoint, new UI feature, new repo methods
- **Large** (9+ files changed): Should be broken into subtasks

If a ticket feels too large, the agent should comment on the issue suggesting a breakdown before starting.

## Dependencies Between Tickets

Tickets may depend on other tickets. Dependencies should be noted in the ticket description. Agents should not start a ticket whose dependencies are not yet merged to `main`.

## Two-Developer Model

- **Frontend Developer**: Owns `frontend/`, reviews frontend PRs, writes BidMan/BatMan tickets.
- **Backend Developer**: Owns `backend/`, `docker/`, reviews backend/infra PRs, writes API/DB/jobs tickets.

Full-stack tickets should be split into a backend subtask and a frontend subtask. The backend subtask should be completed first.
