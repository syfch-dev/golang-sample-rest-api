# Go Clean Architecture REST API — Learning Project

> Personal learning repository. Rebuilding and studying [AleksK1NG/api-mc](https://github.com/AleksK1NG/api-mc) step by step to understand Go backend development, Clean Architecture, and the Echo framework.

**This is a work-in-progress project.** Infrastructure and core packages are being built first; HTTP handlers, modules, and Docker tooling will follow.

**This is not production-ready.** The goal is to learn architecture, layering, and real-world backend patterns.

---

## Current Status

### Implemented

- [x] Project layout (`cmd/`, `config/`, `pkg/`)
- [x] YAML configuration loading ([Viper](https://github.com/spf13/viper))
- [x] Strongly-typed config structs (`ParseConfig`)
- [x] Structured logging ([Zap](https://github.com/uber-go/zap))
- [x] Jaeger / OpenTracing tracer setup
- [x] Local config file (`config/config-local.yml`)
- [x] VS Code / Cursor debug config (`.vscode/launch.json`)

### In Progress / Planned

- [ ] Echo HTTP server and middleware pipeline
- [ ] `internal/` modules (auth, news, comments, session)
- [ ] PostgreSQL + sqlx + migrations
- [ ] Redis session cache
- [ ] JWT + session + CSRF authentication
- [ ] MinIO (S3) avatar uploads
- [ ] Swagger API documentation
- [ ] Prometheus + Grafana metrics
- [ ] Docker Compose local environment
- [ ] Makefile and helper commands

---

## Quick Start (Current)

What works **today**:

### Prerequisites

- [Go 1.24+](https://go.dev/dl/)

### Run

```bash
git clone https://github.com/<your-username>/<your-repo>.git
cd <your-repo>

go mod download
go run ./cmd/api
```

Build a binary:

```bash
go build -o api ./cmd/api
./api
```

### Configuration

Config files live in `config/`. By default, `config-local.yml` is loaded.

| `config` env value | Config file |
|--------------------|-------------|
| *(empty)* | `config/config-local.yml` |
| `docker` | `config/config-docker.yml` |

```bash
# Local (default)
go run ./cmd/api

# Docker profile (when config-docker.yml is added)
config=docker go run ./cmd/api
```

Environment variables can override YAML values via Viper `AutomaticEnv()`.

### Debug (VS Code / Cursor)

1. Install the [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.go)
2. Install Delve: `go install github.com/go-delve/delve/cmd/dlv@latest`
3. Press **F5** → select **Launch API**

---

## Target Architecture

This is the **goal** — layers will be added incrementally.

```
HTTP Request
    ↓
Echo Middleware (CORS, Recover, RequestID, Metrics, ...)
    ↓
Handler (delivery/http)     ← HTTP, JSON, validation
    ↓
UseCase (usecase)           ← business logic
    ↓
Repository (repository)     ← PostgreSQL / Redis / MinIO
```

### Target Project Structure

```
.
├── cmd/api/              # Application entry point (main.go)
├── config/               # YAML configuration
├── internal/             # [planned] Private application code
│   ├── auth/             # Auth module (delivery, usecase, repository)
│   ├── news/             # News module
│   ├── comments/         # Comments module
│   ├── session/          # Redis session
│   ├── middleware/       # Auth, CSRF, logging, metrics
│   ├── models/           # Domain models
│   └── server/           # Echo server, MapHandlers
├── pkg/                  # Shared packages (logger, utils, db, ...)
├── migrations/           # [planned] PostgreSQL migration files
└── docker/               # [planned] Dockerfile, Prometheus config
```

### Current Structure

```
.
├── cmd/api/main.go
├── config/
│   ├── config.go
│   └── config-local.yml
└── pkg/
    ├── logger/zap_logger.go
    └── utils/http.go
```

---

## Tech Stack

| Category | Tool | Status |
|----------|------|--------|
| Config | [Viper](https://github.com/spf13/viper) | ✅ |
| Logging | [Zap](https://github.com/uber-go/zap) | ✅ |
| Tracing | [Jaeger](https://www.jaegertracing.io/) + OpenTracing | ✅ |
| Framework | [Echo v4](https://echo.labstack.com/) | 🔜 |
| Database | PostgreSQL, [sqlx](https://github.com/jmoiron/sqlx), [pgx](https://github.com/jackc/pgx) | 🔜 |
| Cache | [go-redis](https://github.com/go-redis/redis) | 🔜 |
| Object Storage | [MinIO](https://min.io/) | 🔜 |
| Auth | JWT, session cookie, CSRF | 🔜 |
| Validation | [validator](https://github.com/go-playground/validator) | 🔜 |
| API Docs | [Swaggo](https://github.com/swaggo/swag) | 🔜 |
| Metrics | [Prometheus](https://prometheus.io/) + [Grafana](https://grafana.com/) | 🔜 |
| Testing | [testify](https://github.com/stretchr/testify), [gomock](https://github.com/golang/mock) | 🔜 |
| DevOps | Docker Compose, Makefile | 🔜 |

---

## Roadmap

Development follows the source project in phases:

| Phase | Focus | Key deliverables |
|-------|-------|------------------|
| **1** ✅ | Foundation | Config, logger, Jaeger, project layout |
| **2** 🔜 | HTTP layer | Echo server, middleware, health check |
| **3** 🔜 | Database | PostgreSQL connection, migrations, repository pattern |
| **4** 🔜 | Auth | Register, login, JWT, session, CSRF |
| **5** 🔜 | Modules | News & comments CRUD |
| **6** 🔜 | Storage | MinIO avatar uploads |
| **7** 🔜 | Observability | Prometheus metrics, Grafana dashboards |
| **8** 🔜 | DevOps | Docker Compose, Makefile, Swagger |

---

## Quick Start (Planned)

When Docker and Makefile are added, the full local setup will look like this:

```bash
make local          # Start PostgreSQL, Redis, MinIO, Jaeger, Prometheus, Grafana
make migrate_up     # Apply database migrations
make run            # Start the API
```

### Planned Services & URLs

| Service | URL |
|---------|-----|
| API (HTTPS) | https://localhost:5000 |
| Swagger UI | https://localhost:5000/swagger/index.html |
| Jaeger UI | http://localhost:16686 |
| Prometheus | http://localhost:9090 |
| Grafana | http://localhost:3000 |
| MinIO | http://localhost:9000 |

### Planned Endpoints

| Route | Description |
|-------|-------------|
| `POST /api/v1/auth/register` | User registration |
| `POST /api/v1/auth/login` | Login |
| `GET /api/v1/auth/me` | Current authenticated user |
| `GET/POST /api/v1/news/...` | News CRUD |
| `GET/POST /api/v1/comments/...` | Comments CRUD |
| `GET /api/v1/health` | Health check |

### Planned Makefile Commands

```bash
make local          # Start Docker infrastructure
make run            # Run the API
make build          # Build binary
make test           # Run tests
make migrate_up     # Apply migration
make migrate_down   # Roll back last migration
make swaggo         # Generate Swagger docs
make run-linter     # Run golangci-lint
make develop        # Docker dev environment
make docker_delve   # Docker with Delve debugger
```

---

## Learning Checklist

Track progress while rebuilding the source project:

- [x] Understand config flow: `LoadConfig` → `ParseConfig` → typed struct
- [x] Understand `Logger` interface vs `apiLogger` implementation
- [x] Understand Jaeger tracer setup and `defer closer.Close()`
- [ ] Understand the full flow: `main.go` → config → logger → DB → server
- [ ] Understand `MapHandlers` and the Echo middleware pipeline
- [ ] Trace the auth flow: register → login → JWT/session → middleware
- [ ] Understand repository pattern and usecase separation
- [ ] Explore Redis session and CSRF protection
- [ ] Explore MinIO avatar upload flow

---

## What I'm Learning

- **Clean Architecture** — handler → usecase → repository separation
- **Echo** — routing, middleware pipeline, route groups
- **Manual dependency injection** — wiring in `MapHandlers`
- **PostgreSQL + sqlx** — database access and migrations
- **Redis** — session caching
- **JWT + session + CSRF** — authentication flow
- **MinIO (S3)** — avatar uploads
- **Zap** — structured logging
- **Swagger** — API documentation
- **Jaeger + Prometheus + Grafana** — observability
- **Docker Compose** — local development environment

---

## Source Project

This repository is based on:

- **Original:** [AleksK1NG/api-mc](https://github.com/AleksK1NG/api-mc)
- **Architecture reference:** [Clean Architecture — Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

---

## License

Subject to the source project's license. For personal learning purposes.

---

## Disclaimer

This project is for educational purposes only. Security, performance, and production best practices may be simplified or not yet implemented.
