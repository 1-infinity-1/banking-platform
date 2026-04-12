# Ledger Service

> [English](README.md) · [Русский](README.ru.md)

Ledger Service is responsible for accounting and audit records of operations.  
It asynchronously consumes completed transaction events, stores immutable ledger entries, and builds account statements for a requested period.

Main responsibilities:
- consuming completed transaction events;
- immutable log storage;
- operation audit trail;
- account statement generation;
- separation of accounting model from transactional model.

---

## gRPC API (planned)

| Method | Request | Response | Description |
|--------|---------|----------|-------------|
| `GetStatement` | `GetStatementRequest` (account_id, from, to) | `Statement` | Account statement for a requested period |

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
