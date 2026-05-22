# GoAdmin

A production-ready admin panel starter built with [Go](https://golang.org/) and [React](https://reactjs.org/). Batteries included: JWT authentication, Google OAuth, PostgreSQL persistence, OpenAPI-validated endpoints, and an Ant Design frontend.

## Features

- **Authentication** — username/password login, signup, and Google OAuth sign-in
- **JWT tokens** — access + refresh tokens with revocation on logout
- **User management** — list, view, and update users via REST API
- **OpenAPI 3.0** — request validation middleware and auto-generated API docs
- **RBAC-ready** — database schema for roles, permissions, and ReBAC primitives
- **Observability** — structured logging (slog), OpenTelemetry hooks (optional)

## Project Structure

```
goadmin/
  backend/          # Go API server
    cmd/api/        # Application entrypoint
    config/api/     # TOML configs (base, development, production, test)
    database/       # SQL migrations and seed data
    internal/
      auth/         # Auth handlers, service, middleware, Google ID token
      user/         # User handlers and service
      domain/       # Domain models and repository interfaces
      repository/   # PostgreSQL implementations (pgx v5)
      cmd/api/      # Router, server, config, OpenAPI validator
      platform/     # Reusable packages (httpjson, httperr, httproute, logging)
    openapi.yaml    # API specification (OpenAPI 3.0.3)
  frontend/
    ant-design/     # Primary React frontend (CRA + Ant Design 5 + Pro Components)
    react-admin/    # Alternative frontend scaffold (Vite + react-admin)
```

## Quick Start

### Prerequisites

- Go 1.22+
- Node.js 18+
- PostgreSQL 15+
- [golang-migrate](https://github.com/golang-migrate/migrate) (for database migrations)

### 1. Backend

```bash
cd backend

# Copy and fill in your config
cp config/api/development.toml config/api/local.toml

# Set environment variables for secrets
export API__AUTH__JWT_SECRET="your-secret-key"
export GOOGLE__CLIENT_ID="your-google-client-id"
export GOOGLE__CLIENT_SECRET="your-google-client-secret"

# Run database migrations
make migrate-up

# Start the server (default: :3600)
go run cmd/api/main.go
```

### 2. Frontend

```bash
cd frontend/ant-design

# Install dependencies
npm install

# Start dev server (default: :3000)
npm start
```

The frontend proxies API calls to `http://localhost:3600`. Set `REACT_APP_BACKEND_URL` in `.env` to override.

### 3. Configuration

Backend config is loaded from `config/api/base.toml` merged with `config/api/{env}.toml`, then overridden by environment variables (double-underscore `__` maps to nested keys). Key settings:

| Env Variable | Config Path | Description |
|---|---|---|
| `API__AUTH__JWT_SECRET` | `api.auth.jwt_secret` | HS256 signing key |
| `API__PORT` | `api.port` | Server port (default 3600) |
| `DATABASE_URL` | `database_url` | PostgreSQL connection string |
| `GOOGLE__CLIENT_ID` | `google.client_id` | Google OAuth client ID |

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/health` | Public | Health check |
| POST | `/auth/login` | Public | Username/password login |
| POST | `/auth/signup` | Public | Register new user |
| POST | `/auth/signin-with-google` | Public | Google OAuth sign-in |
| POST | `/auth/logout` | Bearer | Invalidate current token |
| GET | `/auth/profile` | Bearer | Get current user profile |
| GET | `/v1/users` | Bearer | List all users |
| GET | `/v1/users/{id}` | Bearer | Get user by ID |
| PATCH | `/v1/users/{id}` | Bearer | Update user |
| GET | `/v1/users/{id}/roles` | Bearer | List user roles |

Full API spec at `backend/openapi.yaml`.

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go 1.22, Chi router, pgx v5 |
| Auth | golang-jwt (HS256), bcrypt, Google ID Token |
| Frontend | React 18, TypeScript, Ant Design 5, Pro Components |
| Database | PostgreSQL with golang-migrate |
| Spec | OpenAPI 3.0.3 (libopenapi-validator) |
| Logging | slog + tint (pretty console) |
| Telemetry | OpenTelemetry (optional, Honeycomb-ready) |

## License

MIT
