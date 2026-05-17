# Ledger Service

> [English](README.md) · [Русский](README.ru.md)

Ledger Service is responsible for accounting and audit records of operations.  
It asynchronously consumes completed transaction events, stores immutable ledger entries, and builds account statements for a requested period.

**Status:** implemented — Cobra CLI, gRPC server, Kafka consumer, PostgreSQL via pgx + goose migrations, Docker.

Main responsibilities:
- consuming completed transaction events from Kafka;
- immutable log storage (idempotent — duplicate events are safely ignored);
- operation audit trail;
- account statement generation via gRPC;
- separation of accounting model from transactional model.

---

## gRPC API

| Method | Request | Response | Description |
|--------|---------|----------|-------------|
| `GetStatement` | `GetStatementRequest` (account_id, from, to) | `Statement` | Account statement for a requested period |

Proto definition: [`pkg/proto/ledger/ledger.proto`](../../pkg/proto/ledger/ledger.proto)

---

## Kafka

| Role | Topic | Description |
|------|-------|-------------|
| Consumer | `transactions.completed` | Receives completed transaction events from Transaction Service |

---

## Dependencies

| Direction | Service | Transport |
|-----------|---------|-----------|
| Called by | Gateway Service | gRPC |
| Consumes from | `transactions.completed` | Kafka |
| Calls | — | — |

---

## Configuration

All environment variables use the `LEDGER_` prefix.

| Variable | Default | Description |
|----------|---------|-------------|
| `LEDGER_LOG_LEVEL` | `info` | Log level (`local`, `dev`, `prod`, `info`) |
| `LEDGER_DB_HOST` | `localhost` | PostgreSQL host |
| `LEDGER_DB_PORT` | `5432` | PostgreSQL port |
| `LEDGER_DB_USER` | `postgres` | PostgreSQL user |
| `LEDGER_DB_PASSWORD` | `postgres` | PostgreSQL password |
| `LEDGER_DB_NAME` | `ledger_db` | PostgreSQL database name |
| `LEDGER_GRPC_PORT` | `8083` | gRPC server port |
| `LEDGER_KAFKA_BROKERS` | `localhost:9092` | Comma-separated Kafka broker addresses |
| `LEDGER_KAFKA_TOPIC` | `transactions.completed` | Kafka topic to consume |
| `LEDGER_KAFKA_GROUP` | `ledger-service` | Kafka consumer group ID |

---

## Running locally

Requirements: PostgreSQL, Kafka.

```bash
# Copy and adjust env file
cp local.env.example local.env

# Run gRPC server
cd internal/ledger-service && go run ./main.go grpc

# Run Kafka consumer (separate terminal)
cd internal/ledger-service && go run ./main.go consumer
```

---

## Running with Docker

```bash
cd internal/ledger-service

# Start all services (gRPC, consumer, PostgreSQL, Kafka, migrations)
make up

# Stop all services
make down
```
