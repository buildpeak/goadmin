# GoAdmin Backend

Go API server for the GoAdmin admin panel.

## Quick Start

```bash
# Install dependencies
go mod download

# Run migrations (requires golang-migrate)
make migrate-up

# Start the server
go run cmd/api/main.go
```

## Configuration

Config files live in `config/api/`. The base config is merged with an environment-specific TOML file (selected by the `ENV` environment variable, defaults to `development`). Environment variables override TOML values — use double underscore `__` for nested keys:

```bash
export API__AUTH__JWT_SECRET="your-secret"
export DATABASE_URL="postgres://user:pass@localhost/goadmin?sslmode=disable"
```

## Makefile Targets

| Target | Description |
|--------|-------------|
| `make test` | Run all tests with race detection and coverage |
| `make lint` | Run golangci-lint |
| `make migrate-up` | Apply pending database migrations |
| `make migrate-down` | Roll back the last migration |
| `make migrate-create name=<name>` | Create a new migration pair |

## Project Layout

```
cmd/api/main.go            # Entrypoint
config/api/                # TOML configuration files
database/migrations/       # SQL migration files
internal/
  auth/                    # Authentication (handlers, service, middleware, Google)
  user/                    # User handlers and service
  domain/                  # Domain models and repository interfaces
  repository/postgres/     # PostgreSQL repository implementations
  cmd/api/                 # Router, server, config, OpenAPI validator
  platform/                # Shared packages (httpjson, httperr, logging, etc.)
  rebac/                   # Relationship-based access control (WIP)
openapi.yaml               # API specification (OpenAPI 3.0.3)
```
