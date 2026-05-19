# Account Service

> [English](README.md) · [Русский](README.ru.md)

Account Service is responsible for bank accounts and current balances.
It is the source of truth for monetary balances, manages the account lifecycle and statuses, and performs idempotent debit and credit operations.

Main responsibilities:
- account creation and retrieval;
- current balance storage;
- account status management;
- idempotent debit operations;
- idempotent credit operations;
- consistency of monetary data.

---

## gRPC API

Proto source: `pkg/proto/account/account.proto`
Default port: **8083**

### AccountService

| Method | Request | Response | Description |
|--------|---------|----------|-------------|
| `CreateAccount` | `CreateAccountRequest` (user_id, currency) | `Account` | Open a new bank account |
| `GetUserAccounts` | `GetUserAccountsRequest` (user_id) | `AccountsList` | All accounts belonging to a user |
| `GetAccount` | `GetAccountRequest` (account_id) | `Account` | Single account data |
| `GetBalance` | `GetBalanceRequest` (account_id) | `Balance` | Current account balance |
| `UpdateStatus` | `UpdateStatusRequest` (account_id, status) | `UpdateStatusResponse` | Change account status (active / blocked / closed) |
| `Debit` | `DebitRequest` (account_id, amount, idempotency_key) | `DebitResponse` | Idempotent debit (called by Transaction Service) |
| `Credit` | `CreditRequest` (account_id, amount, idempotency_key) | `CreditResponse` | Idempotent credit (called by Transaction Service) |

### Key message types

| Message | Fields |
|---------|--------|
| `Account` | id (UUID), user_id (UUID), currency, balance (decimal string), status, created_at, updated_at |
| `Balance` | account_id (UUID), amount (decimal string), currency |
| `AccountStatus` | UNSPECIFIED, ACTIVE, BLOCKED, CLOSED |
| `DebitResponse` / `CreditResponse` | account_id (UUID), balance_after (decimal string) |

### Error mapping (`UnaryErrorInterceptor`)

| Domain error | gRPC code |
|--------------|-----------|
| `NotFoundError` | `NotFound` |
| `InvalidParamsError` | `InvalidArgument` |
| `BusinessError` | `FailedPrecondition` |
| `ConflictError` | `AlreadyExists` |
| anything else | `Internal` (logged once with trace_id/request_id) |

---

## Architecture

### Idempotent Debit / Credit

Both `Debit` and `Credit` are idempotent: repeated calls with the same `idempotency_key` return the originally stored `balance_after` without re-applying the operation. The contract is enforced by the `account_operations` table (`idempotency_key UNIQUE`).

Each RPC runs inside a single `tx.Manager.BeginFunc` transaction with three steps:

1. `INSERT` an `account_operations` row (`op_type`, `amount`, `idempotency_key`). On `23505` (unique violation), `SELECT` the existing operation by `idempotency_key` and return its `balance_after` — the saga short-circuits.
2. Atomic conditional `UPDATE accounts`:
   - **Debit:** `UPDATE accounts SET balance = balance - $amount WHERE public_id = $id AND balance >= $amount AND status = 'active' RETURNING balance`
   - **Credit:** `UPDATE accounts SET balance = balance + $amount WHERE public_id = $id AND status = 'active' RETURNING balance`
   - On zero rows the repository disambiguates with a follow-up `SELECT`: missing → `NotFoundError`, non-active → `BusinessError("account is not active: …")`, debit-only insufficient → `BusinessError("insufficient funds")`. The whole transaction rolls back, so the `account_operations` row inserted at step 1 disappears too — a subsequent retry is free to try again.
3. `UPDATE account_operations SET balance_after = $newBalance` on the operation row, then commit.

There is no `SELECT ... FOR UPDATE` and no application-level version column. Concurrency is handled by the conditional `UPDATE` itself, which is atomic in PostgreSQL.

### Error propagation

- Repository returns typed domain errors (`NotFoundError`, `BusinessError`, `ConflictError`, `InvalidParamsError`) and never logs.
- Service wraps errors with `fmt.Errorf("op: %w", err)` to add context and never logs.
- `UnaryErrorInterceptor` is the only place that logs unexpected errors (with `trace_id` / `request_id`) and maps them to gRPC status codes.

---

## Configuration

Loaded from `local.env` (local dev) or `docker.env` (Docker), with prefix `ACCOUNT`.

| ENV | Default | Description |
|-----|---------|-------------|
| `ACCOUNT_LOG_LEVEL` | `info` | Log verbosity (debug / info / warn / error) |
| `ACCOUNT_DB_HOST` | `localhost` | PostgreSQL host |
| `ACCOUNT_DB_PORT` | `5432` | PostgreSQL port |
| `ACCOUNT_DB_USER` | `postgres` | PostgreSQL user |
| `ACCOUNT_DB_PASSWORD` | `postgres` | PostgreSQL password |
| `ACCOUNT_DB_NAME` | `app_db` | PostgreSQL database name |
| `ACCOUNT_GRPC_PORT` | `8083` | gRPC server listen port |

---

## Running locally

```bash
make up    # docker compose up -d  (PostgreSQL + goose migrate + app)
make down  # docker compose down
```

Run from source:

```bash
cd internal/account-service
go run ./main.go application
```

A reachable PostgreSQL is required at the configured `ACCOUNT_DB_*` address.

---

## Database

PostgreSQL 17 (port `5433` in docker compose, mapped to internal `5432`). Migrations managed with [Goose](https://github.com/pressly/goose).
Migration files: `migrations/`

| Table | Description |
|-------|-------------|
| `accounts` | Bank accounts: `public_id` (UUID, UNIQUE), `user_id` (UUID), `currency`, `balance` (NUMERIC(20,8)), `status`, timestamps |
| `account_operations` | Idempotent log of debit/credit operations: `account_id` (UUID, FK → `accounts.public_id`), `op_type` (`debit`/`credit`), `amount` (NUMERIC(20,8)), `balance_after` (NUMERIC(20,8)), `idempotency_key` (UNIQUE), timestamps |

Indexes:

| Index | Columns |
|-------|---------|
| `accounts_user_id_idx` | `user_id` |
| `account_operations_account_id_idx` | `account_id` |

---

## Dependencies

| Direction | Service | Transport |
|-----------|---------|-----------|
| Called by | Gateway Service | gRPC |
| Called by | Transaction Service (`Debit`, `Credit`) | gRPC |
| Calls | — | — |

---

## Project Structure

```text
internal/account-service/
├── cmd/
│   ├── application.go              # cobra command wiring
│   └── cmd.go
├── internal/
│   ├── app/
│   │   ├── application.go          # dependency wiring (composition root)
│   │   └── grpc/
│   │       └── grpc.go             # gRPC server lifecycle
│   ├── config/
│   │   └── config.go               # envconfig (prefix ACCOUNT)
│   ├── models/
│   │   ├── account.go              # domain types (Account, Balance, requests/results)
│   │   └── errors.go               # typed domain errors
│   ├── services/
│   │   └── account/
│   │       └── service.go          # CRUD + Debit/Credit business logic
│   ├── storage/
│   │   ├── account/
│   │   │   ├── dto.go              # row DTO + ToDomain
│   │   │   └── repository.go       # pgx repository (idempotent Debit/Credit)
│   │   └── tx/
│   │       └── tx_manager.go       # BeginFunc transaction wrapper
│   └── transport/
│       └── grpc/
│           ├── create_account.go
│           ├── credit.go
│           ├── debit.go
│           ├── get_account.go
│           ├── get_balance.go
│           ├── get_user_accounts.go
│           ├── interceptor/
│           │   └── error.go        # domain error → gRPC status mapping
│           ├── mapping.go          # domain ↔ proto helpers
│           ├── server.go           # serverAPI + AccountService interface
│           └── update_status.go
├── migrations/
│   ├── 20260517000000_init_schema.sql
│   └── 20260519000000_account_operations.sql
├── Dockerfile
├── Makefile
├── docker-compose.yaml
├── docker.env
├── local.env
└── main.go
```
