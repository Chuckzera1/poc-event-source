# CLAUDE.md — poc-event-source

## Project Overview
Event sourcing proof-of-concept API in Go. Stores domain events in PostgreSQL and distributes them via NATS JetStream. Exposes a REST API with Gin.

- **Module:** `poc-event-source`
- **Go version:** 1.23.5
- **Web:** Gin v1.10.0
- **DB:** PostgreSQL + GORM v1.25.12 (JSONB payloads)
- **Broker:** NATS v1.41.2 with JetStream
- **Testing:** testify + testcontainers-go

---

## Architecture — Clean / Hexagonal Layers

```
domain          ← pure business models, zero imports from other internal packages
application     ← interfaces (irepository/, event/), DTOs, use cases — imports only domain
infrastructure  ← NATS & GORM implementations of application interfaces
repository      ← concrete DB repos implementing irepository interfaces
api             ← Gin handlers, route registration, messaging subscribers — wires everything
```

**Strict import direction:** outer layers depend on inner layers, never the reverse.
`api` and `infrastructure` depend on `application`; `application` depends on `domain`; `domain` depends on nothing internal.

---

## Conventions

### Constructor pattern
Every exported struct lives behind an interface and is created via a `New*()` function in its own `new.go` file.

```go
// new.go
type MyUseCase interface { Do(ctx context.Context) error }
type myUseCase struct { repo irepository.MyRepository }
func NewMyUseCase(repo irepository.MyRepository) MyUseCase { return &myUseCase{repo: repo} }
```

### One responsibility per file
- `new.go` — constructor + interface declaration
- `create.go`, `handle.go`, `find.go` — one operation each
- `create_test.go` — integration test for that operation

### Interface-first
Define the interface in `application/irepository/` or `application/usecase/<name>/new.go` **before** writing the implementation.

### Always use `context.Context`
Every method that touches the DB or broker must accept `ctx context.Context` as the first parameter.

### Error handling
Return errors, never panic. Wrap with context where useful (`fmt.Errorf("creating user: %w", err)`).

---

## Testing

- **Integration tests** use testcontainers — see `internal/repository/testutils/` for DB setup helpers.
- **Never mock the database.** Use a real containerised PostgreSQL instance.
- Tests wrap operations in a transaction and defer `tx.Rollback()` to keep the DB clean.
- Run all tests: `go test ./...`

---

## Common Commands

```bash
# Start dependencies
docker-compose up -d

# Run API
go build ./cmd/api && ./api

# Run all tests (requires Docker for testcontainers)
go test ./...

# Run a single package test
go test ./internal/repository/user/...
```

---

## What NOT to Do

- No business logic in `infrastructure/` — it only adapts external systems to interfaces.
- No direct DB calls from `api/` — always go through a use case or repository interface.
- No skipping the DTO layer — HTTP request bodies map to DTOs, not directly to domain/model structs.
- Do not add `fmt.Printf` placeholders in route handlers — implement or leave a `TODO` comment.
- Do not import `infrastructure/model` from `application/` use case implementations directly (handle.go in the event use case is a known exception to clean up).
