# Auth Service

> [English](README.md) · [Русский](README.ru.md)

Auth Service is responsible for user management and authentication.  
It stores user accounts, handles signup and login, issues and validates tokens, manages sessions, roles, and permissions.

Main responsibilities:
- user signup and login;
- JWT/access token issuing and validation;
- refresh tokens and user sessions;
- roles and permissions management;
- exposing user data to other services via gRPC.

---

## gRPC API

Proto source: `pkg/proto/auth/auth.proto`  
Default port: **8082**

### AuthService

| Method | Request | Response | Description |
|--------|---------|----------|-------------|
| `Login` | `LoginRequest` (login, password, RequestContext) | `LoginResponse` (User, Session, TokenPair, AuthContext) | Authenticate user, open session, issue token pair |
| `Logout` | `LogoutRequest` (refresh_token, RequestContext) | `Empty` | Revoke refresh token, close session |
| `RefreshToken` | `RefreshTokenRequest` (refresh_token, RequestContext) | `RefreshTokenResponse` (TokenPair, AuthContext) | Issue new token pair using refresh token |

### AccessManagementService

| Method | Request | Response | Description |
|--------|---------|----------|-------------|
| `CreateUser` | `CreateUserRequest` (login, email, phone, password, role_codes[]) | `CreateUserResponse` (User) | Create a new user with assigned roles |

### Key message types

| Message | Fields |
|---------|--------|
| `TokenPair` | access_token, refresh_token, token_type, access_token_expires_at, refresh_token_expires_at |
| `Session` | id, user_id, status, device, created_at, updated_at, expires_at, last_seen_at |
| `AuthContext` | user_id, session_id, role_codes[], permission_codes[] |
| `RequestContext` | trace_id, request_id, ip, user_agent, device_id, platform |

---

## Configuration

Loaded from `local.env` (local dev) or `docker.env` (Docker), with prefix `AUTH`.

| ENV | Default | Description |
|-----|---------|-------------|
| `AUTH_LOG_LEVEL` | `info` | Log verbosity (debug / info / warn / error) |
| `AUTH_GRPC_PORT` | `8082` | gRPC server listen port |

---

## Running locally

```bash
make up    # docker compose up -d  (PostgreSQL + goose migrate + app)
make down  # docker compose down
```

---

## Database

PostgreSQL 17 (port 5432). Migrations managed with [Goose](https://github.com/pressly/goose).  
Migration files: `migrations/`

| Table | Description |
|-------|-------------|
| `users` | User accounts (login, email, phone, password_hash, status, lock info) |
| `devices` | Device fingerprints (platform, user_agent) |
| `sessions` | User sessions with status, expiry, last_seen_at |
| `refresh_tokens` | Token hashes with expiration and revocation flag |
| `roles` | Named roles (code, name) |
| `permissions` | Named permissions (code, name) |
| `user_roles` | Many-to-many: users ↔ roles |
| `role_permissions` | Many-to-many: roles ↔ permissions |

---

## Dependencies

| Direction | Service | Transport |
|-----------|---------|-----------|
| Called by | Gateway Service | gRPC |
| Calls | — | — |

---

## Project Structure

```text
internal/auth-service/
├── cmd/
│   ├── application.go      # cobra command wiring
│   └── cmd.go
├── internal/
│   ├── app/
│   │   ├── application.go  # dependency wiring
│   │   └── grpc/
│   │       └── grpc.go     # gRPC server lifecycle
│   ├── config/
│   │   └── config.go       # service config struct
│   ├── jwt/
│   │   └── manager.go      # JWT sign / verify
│   ├── models/
│   │   ├── errors.go
│   │   └── user.go
│   ├── services/
│   │   └── auth/
│   │       └── service.go  # business logic
│   ├── storage/
│   │   └── user/
│   │       ├── dto.go
│   │       └── repository.go
│   └── transport/
│       └── grpc/
│           └── server.go   # gRPC handler implementations
├── migrations/
│   └── 20260330175820_init_schema.sql
├── Dockerfile
├── Makefile
├── docker-compose.yaml
├── docker.env
├── local.env
└── main.go
```
