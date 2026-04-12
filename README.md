# Banking Platform

> [English](README.md) · [Русский](README.ru.md)

Banking Platform is a backend project in Go implemented as a microservice monorepo.  
The project models a basic banking platform with a single entry point for client applications, a user and authentication service, an account service, a transaction service, and a ledger service.

Project goals:
- design and implement a microservice architecture in Go;
- define service-to-service contracts;
- design data schemas and business flows;
- demonstrate clear responsibility boundaries between services;
- practice REST, gRPC, Kafka, idempotency, and audit log patterns.

---

# Services

## `gateway-service`
Single REST/BFF entry point for web and mobile clients. Routes requests to internal services, aggregates data, and provides the external API boundary of the system.

## `auth-service`
Manages users, authentication, sessions, roles, and permissions. Issues and validates tokens and provides user data to other services.

## `account-service`
Responsible for bank accounts, their statuses, and current balances. Acts as the source of truth for monetary balances.

## `transaction-service`
Orchestrates money operations such as transfers and replenishments. Coordinates debit/credit flows, stores transaction statuses, and publishes completed transaction events.

## `ledger-service`
Maintains accounting and audit records of operations. Asynchronously consumes transaction events and builds account statements.

---

# Architecture

The system is built around the following principles:
- client applications communicate only with `gateway-service` over REST;
- internal synchronous service-to-service communication uses gRPC;
- asynchronous integration uses Kafka;
- `account-service` stores current balances;
- `ledger-service` stores immutable ledger history;
- `transaction-service` manages money operation workflows;
- `auth-service` owns authentication and user access context.

---

# Service Interaction

| From | To | Transport | Purpose |
|---|---|---|---|
| Client Apps | Gateway Service | REST | External API |
| Gateway Service | Auth Service | gRPC | Auth and user data |
| Gateway Service | Account Service | gRPC | Accounts and balances |
| Gateway Service | Transaction Service | gRPC | Transfers and history |
| Gateway Service | Ledger Service | gRPC | Statements |
| Transaction Service | Account Service | gRPC | Debit / credit |
| Transaction Service | Kafka | Kafka | Publish completed transaction events |
| Kafka | Ledger Service | Kafka | Consume transaction events |

---

# Repository Structure

```text
.
├── internal/
│   ├── auth-service/
│   ├── gateway-service/
│   ├── account-service/
│   ├── transaction-service/
│   └── ledger-service/
├── pkg/
│   ├── config/
│   ├── logger/
│   └── proto/
│       ├── auth/
│       └── generated/
├── go.mod
└── Makefile
```
