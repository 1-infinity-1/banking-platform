# Transaction Service

> [English](README.md) · [Русский](README.ru.md)

Transaction Service is responsible for processing money movements between accounts.
It accepts transfer and replenish commands, coordinates debit and credit operations through Account Service, ensures idempotency, and publishes events for completed transactions.

Main responsibilities:
- transfers between accounts;
- account replenishment;
- orchestration of debit/credit flows;
- transaction status storage;
- operation idempotency;
- Kafka event publishing.

---

## gRPC API

Proto source: `pkg/proto/transaction/transaction.proto`
Default port: **8084**

### TransactionService

| Method | Request | Response | Description |
|--------|---------|----------|-------------|
| `Transfer` | `TransferRequest` (from_account_id, to_account_id, amount, currency, idempotency_key) | `Transaction` | Transfer between accounts via Pending → Completed saga |
| `Replenish` | `ReplenishRequest` (to_account_id, amount, currency, idempotency_key) | `Transaction` | Account replenishment via Pending → Completed saga |
| `GetHistory` | `GetHistoryRequest` (account_id, limit, offset) | `TransactionsList` | Operation history for an account (DESC by created_at) |
| `GetTransaction` | `GetTransactionRequest` (transaction_id) | `Transaction` | Single transaction details |

### Key message types

| Message | Fields |
|---------|--------|
| `Transaction` | id (UUID), from_account_id (UUID, optional), to_account_id (UUID), amount (decimal string), currency, status, idempotency_key, created_at, updated_at |
| `TransactionStatus` | UNSPECIFIED, PENDING, COMPLETED, FAILED, CANCELLED |
| `TransactionsList` | transactions[] |

### Error mapping (`UnaryErrorInterceptor`)

| Domain error | gRPC code |
|--------------|-----------|
| `NotFoundError` | `NotFound` |
| `InvalidParamsError` | `InvalidArgument` |
| `BusinessError` | `FailedPrecondition` |
| `ConflictError` | `AlreadyExists` |
| anything else | `Internal` (logged once with trace_id/request_id) |

---

## Saga (Pending → Completed)

Both `Transfer` and `Replenish` execute the same skeleton, with `Transfer` adding a `Debit` step:

1. Validate the request (positive amount, non-empty currency / idempotency_key, distinct accounts for Transfer).
2. `INSERT` a row with status `pending`. Repeated calls with the same `idempotency_key` hit the unique constraint, the repository SELECTs the existing row and returns it — the saga short-circuits:
   - existing `completed` → best-effort re-publish, return the existing transaction;
   - existing `failed` → return `BusinessError`.
3. `Transfer` only: call `account-service.Debit(from, amount, "<tx-uuid>:debit")`. On error mark the row `failed` and return.
4. Call `account-service.Credit(to, amount, "<tx-uuid>:credit")`. On error mark the row `failed` and return.
5. Update the row status to `completed`.
6. Publish a `TransactionEvent` to Kafka topic `transactions.completed`.

The debit / credit idempotency keys are derived from the transaction public UUID so retries from the same saga are deduplicated by Account Service.

---

## Kafka

| Role | Topic | Description |
|------|-------|-------------|
| Producer | `transactions.completed` | Published after a successful transaction; consumed by Ledger Service |

`TransactionEvent` payload (JSON):

```json
{
  "transaction_id": "uuid",
  "from_account_id": "uuid|null",
  "to_account_id": "uuid",
  "amount": "decimal-as-string",
  "currency": "RUB",
  "status": "completed",
  "occurred_at": "RFC3339"
}
```

---

## Configuration

Loaded from `local.env` (local dev) or `docker.env` (Docker), with prefix `TRANSACTION`.

| ENV | Default | Description |
|-----|---------|-------------|
| `TRANSACTION_LOG_LEVEL` | `info` | Log verbosity (debug / info / warn / error) |
| `TRANSACTION_DB_HOST` | `localhost` | PostgreSQL host |
| `TRANSACTION_DB_PORT` | `5432` | PostgreSQL port |
| `TRANSACTION_DB_USER` | `postgres` | PostgreSQL user |
| `TRANSACTION_DB_PASSWORD` | `postgres` | PostgreSQL password |
| `TRANSACTION_DB_NAME` | `app_db` | PostgreSQL database name |
| `TRANSACTION_GRPC_PORT` | `8084` | gRPC server listen port |
| `TRANSACTION_ACCOUNT_SERVICE_HOST` | `localhost` | Account Service gRPC host |
| `TRANSACTION_ACCOUNT_SERVICE_PORT` | `8083` | Account Service gRPC port |
| `TRANSACTION_KAFKA_BROKERS` | `localhost:9092` | Kafka brokers (comma-separated) |
| `TRANSACTION_KAFKA_TOPIC` | `transactions.completed` | Topic for completed-transaction events |

---

## Running locally

```bash
make up    # docker compose up -d  (PostgreSQL + goose migrate + app)
make down  # docker compose down
```

Run from source:

```bash
cd internal/transaction-service
go run ./main.go application
```

`account-service` and a Kafka broker must be reachable at the configured addresses for transfers and replenishments to succeed end-to-end.

---

## Database

PostgreSQL 17 (port `5434` in docker compose, mapped to internal `5432`). Migrations managed with [Goose](https://github.com/pressly/goose).
Migration files: `migrations/`

| Table | Description |
|-------|-------------|
| `transactions` | Transfer and replenish operations: `from_account_id` (UUID, nullable for replenishments), `to_account_id` (UUID), `amount` (NUMERIC(20,8)), `currency`, `status`, `idempotency_key` (UNIQUE), timestamps |

Indexes:

| Index | Columns |
|-------|---------|
| `transactions_to_account_idx` | `to_account_id` |
| `transactions_from_account_idx` | `from_account_id` (partial, non-NULL) |

---

## Dependencies

| Direction | Service | Transport |
|-----------|---------|-----------|
| Called by | Gateway Service | gRPC |
| Calls | Account Service (`Debit`, `Credit`) | gRPC |
| Publishes to | Ledger Service via `transactions.completed` | Kafka |

---

## Project Structure

```text
internal/transaction-service/
├── cmd/
│   ├── application.go              # cobra command wiring
│   └── cmd.go
├── internal/
│   ├── app/
│   │   ├── application.go          # dependency wiring (composition root)
│   │   └── grpc/
│   │       └── grpc.go             # gRPC server lifecycle
│   ├── clients/
│   │   └── account/
│   │       └── client.go           # gRPC client to account-service
│   ├── config/
│   │   └── config.go               # envconfig (prefix TRANSACTION)
│   ├── kafka/
│   │   └── producer.go             # segmentio/kafka-go producer
│   ├── models/
│   │   ├── errors.go               # typed domain errors
│   │   └── transaction.go          # domain types + TransactionEvent
│   ├── services/
│   │   └── transaction/
│   │       └── service.go          # Transfer/Replenish saga, GetHistory, GetTransaction
│   ├── storage/
│   │   ├── transaction/
│   │   │   ├── dto.go              # row DTO + ToDomain
│   │   │   └── repository.go       # pgx repository (idempotent INSERT)
│   │   └── tx/
│   │       └── tx_manager.go       # BeginFunc transaction wrapper
│   └── transport/
│       └── grpc/
│           ├── get_history.go
│           ├── get_transaction.go
│           ├── interceptor/
│           │   └── error.go        # domain error → gRPC status mapping
│           ├── mapping.go          # domain → proto helpers
│           ├── replenish.go
│           ├── server.go           # serverAPI + TransactionService interface
│           └── transfer.go
├── migrations/
│   └── 20260517000000_init_schema.sql
├── Dockerfile
├── Makefile
├── docker-compose.yaml
├── docker.env
├── local.env
└── main.go
```
