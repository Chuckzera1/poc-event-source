# POC Event Source

Event Sourcing REST API in Go, using NATS JetStream as the broker and PostgreSQL for persistence.

## Architecture

The project follows Clean/Hexagonal Architecture with strict layer separation:

```
domain          → pure business models (no internal imports)
application     → interfaces, use cases, DTOs
infrastructure  → external adapters (NATS, GORM)
repository      → database access implementations
api             → HTTP handlers and NATS subscribers
```

**Event Sourcing flow:**
```
POST /api/v1/user
  └─ CreateUserUseCase
       ├─ saves CREATE_USER event to EventSource table  (write model)
       └─ publishes to NATS topic "user"
                ↓
       UserBroker.Subscribe()
            └─ creates record in User table  (read model / projection)
```

For code conventions and architecture rules, see [CLAUDE.md](CLAUDE.md).

## Stack

| Layer | Technology |
|---|---|
| Language | Go 1.23.5 |
| Web | Gin v1.10.0 |
| Database | PostgreSQL + GORM v1.25.12 |
| Broker | NATS v1.41.2 with JetStream |
| Testing | testify + testcontainers-go |

## Running locally

**Prerequisites:** Docker, Go 1.23+

```bash
# 1. Start PostgreSQL and NATS
docker-compose up -d

# 2. Configure environment variables
cp .env.example .env  # adjust variables as needed

# 3. Run the API
go run ./cmd/api

# 4. Run tests (requires Docker for testcontainers)
go test ./...
```

## Endpoints

| Method | Route | Description |
|---|---|---|
| POST | `/api/v1/user` | Creates a user (event sourcing) |
| GET | `/api/v1/user/:id` | Get user by ID *(not implemented)* |
| PUT | `/api/v1/user/:id` | Update user *(not implemented)* |

**Create user example:**
```bash
curl -X POST http://localhost:8080/api/v1/user \
  -H 'Content-Type: application/json' \
  -d '{"username": "alice", "password": "secret"}'
# Response: 201 {"status": "event published"}
```

