# Banking Platform

> [English](README.md) · [Русский](README.ru.md)

Go microservice monorepo modelling a banking platform. External clients talk to `gateway-service` over REST; internal services communicate over gRPC; async integration via Kafka (planned).

Project goals:
- design and implement a microservice architecture in Go;
- define service-to-service contracts (gRPC + OpenAPI);
- practice REST, gRPC, Kafka, idempotency, and audit log patterns.

---

## Quick Start

**Prerequisites:** Go 1.25.4+, Docker, Make.

```bash
# 1. Start auth-service (Postgres + migrations + gRPC server on :8082)
cd internal/auth-service && make up

# 2. Start gateway-service (HTTP server on :8081)
cd internal/gateway-service && make up
```

Health check:

```bash
curl http://localhost:8081/ping
# {"message":"pong"}
```

---

## Services

| Service | Port | Transport | Status |
|---------|------|-----------|--------|
| `gateway-service` | `:8081` | REST (ogen/OpenAPI) | Implemented |
| `auth-service` | `:8082` | gRPC | Implemented |
| `account-service` | `:8083` | gRPC | Skeleton |
| `transaction-service` | `:8084` | gRPC | Skeleton |
| `ledger-service` | `:8085` | gRPC + Kafka consumer | Implemented |

---

## Documentation

| Guide | Description |
|-------|-------------|
| [Getting Started](docs/getting-started.md) | Prerequisites, running locally, Docker |
| [Architecture](docs/architecture.md) | Service topology, internal structure, Explicit Architecture |
| [C4 Diagrams](docs/c4.md) | C4 Context, Container, Component diagrams (Mermaid) |
| [Sequence Diagrams](docs/diagrams.md) | Login, Transfer, Refresh Token, Get Statement flows |
| [Deployment](docs/deployment.md) | Docker Compose stacks, ports, build instructions |
| [API Reference](docs/api-reference.md) | REST endpoints, gRPC API |
| [Configuration](docs/configuration.md) | Environment variables for all services |
