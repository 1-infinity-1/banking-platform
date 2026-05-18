# Banking Platform

> [English](README.md) · [Русский](README.ru.md)

Go-монорепозиторий, моделирующий банковскую платформу. Внешние клиенты работают с `gateway-service` по REST; внутренние сервисы общаются по gRPC; асинхронная интеграция через Kafka (запланировано).

Цели проекта:
- спроектировать и реализовать микросервисную архитектуру на Go;
- зафиксировать контракты между сервисами (gRPC + OpenAPI);
- отработать подходы к REST, gRPC, Kafka, idempotency и audit log.

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
| `account-service` | `:8083` | gRPC | Скелет |
| `transaction-service` | `:8084` | gRPC | Скелет |
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
