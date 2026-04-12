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

## gRPC API (planned)

| Method | Request | Response | Description |
|--------|---------|----------|-------------|
| `Transfer` | `TransferRequest` | `Transaction` | Transfer between accounts (Saga pattern) |
| `Replenish` | `ReplenishRequest` | `Transaction` | Account replenishment |
| `GetHistory` | `GetHistoryRequest` | `TransactionsList` | Operation history for an account |
| `GetTransaction` | `GetTransactionRequest` | `Transaction` | Single transaction details |

---

## Kafka

| Role | Topic | Description |
|------|-------|-------------|
| Producer | `transactions.completed` | Published after a successful transaction; consumed by Ledger Service |

---

## Dependencies

| Direction | Service | Transport |
|-----------|---------|-----------|
| Called by | Gateway Service | gRPC |
| Calls | account-service (`GetAccount`, `Debit`, `Credit`) | gRPC |
| Publishes to | `transactions.completed` | Kafka |
