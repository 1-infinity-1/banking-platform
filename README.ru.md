# Banking Platform

> [English](README.md) · [Русский](README.ru.md)

Go-монорепозиторий, моделирующий банковскую платформу. Внешние клиенты работают с `gateway-service` по REST; внутренние сервисы общаются по gRPC; асинхронная интеграция через Kafka (запланировано).

Цели проекта:
- спроектировать и реализовать микросервисную архитектуру на Go;
- зафиксировать контракты между сервисами (gRPC + OpenAPI);
- отработать подходы к REST, gRPC, Kafka, idempotency и audit log.

---

## О проекте

Часть сервисов спроектирована и реализована самостоятельно — чтобы выработать единую конвенцию для монорепозитория:

- **Project Layout** — структура директорий по [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
- **Uber Go Style Guide** — соглашения по стилю кода по [Uber Go Guide](https://github.com/uber-go/guide)
- **Clean Architecture** — слои и направление зависимостей внутри сервиса по [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- **Golangci-lint Golden Config** — строгая конфигурация линтера на базе [Golden Config](https://github.com/maratori/golangci-lint-config)

Получив устойчивую конвенцию и настроенную AI-платформу ([Claude Code](https://claude.ai/code) + [AI Factory](https://github.com/lee-to/ai-factory)), дальнейшая разработка ведётся с поддержкой AI без потери качества.

---

## Быстрый старт

**Требования:** Go 1.25.4+, Docker, Make.

```bash
# 1. Запустить auth-service (Postgres + миграции + gRPC-сервер на :8082)
cd internal/auth-service && make up

# 2. Запустить gateway-service (HTTP-сервер на :8081)
cd internal/gateway-service && make up
```

Проверка работоспособности:

```bash
curl http://localhost:8081/ping
# {"message":"pong"}
```

---

## Сервисы

| Сервис | Порт | Транспорт | Статус |
|--------|------|-----------|--------|
| `gateway-service` | `:8081` | REST (ogen/OpenAPI) | Реализован |
| `auth-service` | `:8082` | gRPC | Реализован |
| `account-service` | `:8083` | gRPC | Реализован |
| `transaction-service` | `:8084` | gRPC | Реализован |
| `ledger-service` | `:8085` | gRPC + Kafka consumer | Реализован |

---

## Документация

| Раздел | Описание |
|--------|----------|
| [Начало работы](docs/getting-started.md) | Требования, запуск локально, Docker |
| [Архитектура](docs/architecture.md) | Топология сервисов, внутренняя структура, Explicit Architecture |
| [C4 Диаграммы](docs/c4.md) | C4 Context, Container, Component диаграммы (Mermaid) |
| [Sequence Diagrams](docs/diagrams.md) | Потоки: Login, Transfer, Refresh Token, Get Statement |
| [Развёртывание](docs/deployment.md) | Docker Compose стеки, порты, сборка |
| [API Reference](docs/api-reference.md) | REST-эндпоинты, gRPC API |
| [Конфигурация](docs/configuration.md) | Переменные окружения всех сервисов |
