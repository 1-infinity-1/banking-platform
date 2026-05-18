[← Начало работы](getting-started.md) · [Back to README](../README.md) · [C4 Диаграммы →](c4.md)

# Архитектура

Платформа использует двухуровневый подход: **Microservices** на уровне репозитория и **Explicit Architecture** (синтез Clean Architecture + Hexagonal + Onion) внутри каждого сервиса.

Подробные правила с примерами кода — в [`.ai-factory/ARCHITECTURE.md`](../.ai-factory/ARCHITECTURE.md).

## Топология сервисов

```
Внешний клиент (REST/HTTP)
        │
        ▼
┌───────────────────┐
│  gateway-service  │  ← API Gateway (порт :8081)
└───────────────────┘
        │ gRPC
        ▼
┌───────────────────┐
│   auth-service    │  ← аутентификация, JWT (порт :8082)
└───────────────────┘

        │ gRPC
        ▼
┌────────────────────┐  gRPC  ┌───────────────────────┐
│  account-service   │◄──────►│  transaction-service  │
│  (порт :8083)      │        │  (порт :8084)          │
└────────────────────┘        └───────────────────────┘
                                       │ Kafka
                                       ▼ topic: transactions.completed
                              ┌───────────────────────┐
                              │    ledger-service     │
                              │    (порт :8085)        │
                              └───────────────────────┘
```

## Взаимодействие сервисов

| Отправитель | Получатель | Транспорт | Назначение | Статус |
|-------------|------------|-----------|------------|--------|
| Клиент | gateway-service | REST | Внешний API | Реализован |
| gateway-service | auth-service | gRPC | Аутентификация, управление пользователями | Реализован |
| gateway-service | account-service | gRPC | Счета и балансы | Скелет (без логики) |
| gateway-service | transaction-service | gRPC | Переводы и история | Скелет (без логики) |
| gateway-service | ledger-service | gRPC | Выписки по счетам | Реализован |
| transaction-service | account-service | gRPC | Списание / зачисление | Скелет (без логики) |
| transaction-service | Kafka | Kafka (produce) | Публикация событий транзакций | Скелет (без логики) |
| Kafka | ledger-service | Kafka (consume) | Получение событий транзакций | Реализован |

## Explicit Architecture внутри сервиса

Каждый сервис применяет **Explicit Architecture (Technical Layers)** — прагматичный синтез Clean Architecture, Hexagonal и Onion. Ключевое правило: зависимости направлены внутрь, к домену.

### Слои и их ответственность

| Слой | Пакет | Ответственность |
|------|-------|-----------------|
| **Presentation** | `transport/` | Протокол ↔ domain: gRPC handlers, HTTP handlers, middleware (трейсинг, логирование, error mapping) |
| **Application** | `services/` | Бизнес-оркестрация: транзакционные границы, вызов репозиториев, генерация токенов |
| **Domain** | `models/` | Типы данных и ошибки домена. Без внешних зависимостей |
| **Infrastructure** | `storage/` | pgx-репозитории, tx-manager. Реализует интерфейсы, объявленные в `services/` |
| **Composition Root** | `app/` | Единственное место сборки: DB → repos → services → transport |

### Структура директорий

```
internal/<service>/
├── main.go                  # cmd.Execute()
├── cmd/                     # Cobra CLI: root + application subcommand
└── internal/
    ├── app/                 # Composition root: DB → repos → services → transport
    ├── config/              # envconfig: загрузка переменных окружения
    ├── models/              # Domain layer: типы и ошибки
    ├── services/            # Application layer: бизнес-логика + consumer-side interfaces
    ├── storage/             # Infrastructure layer: pgx-репозитории + tx-manager
    └── transport/           # Presentation layer: gRPC/HTTP handlers + middleware
```

### Dependency Inversion (Hexagonal / Ports & Adapters)

`services/` не зависит от конкретных реализаций репозиториев — только от интерфейсов, которые сами же объявляют (consumer-side interfaces). `storage/` адаптируется под контракт сервиса:

```
services/auth/interfaces.go     ← объявляет userRepo, sessionRepo, ...
storage/user/repository.go      ← реализует интерфейс из services/auth/
```

Это «порты» (interfaces) и «адаптеры» (storage implementations) из Hexagonal Architecture.

## Ключевые паттерны

### Consumer-side interfaces

Сервисы объявляют интерфейсы над репозиториями у себя в `services/<usecase>/interfaces.go`, а не в `storage/`. Зависимость инвертирована — инфраструктура адаптируется к бизнес-логике.

### Transaction boundary

Многошаговые операции выполняются через `storage/tx.Manager.BeginFunc(ctx, fn)`. Репозитории принимают `pgx.Tx`; сервис компонует их внутри `BeginFunc`.

### Typed domain errors

`models/errors.go` содержит `NotFoundError`, `InvalidParamsError`, `BusinessError`, `ConflictError`, `ErrInternal`. Repository маппирует DB-ошибки; transport маппирует domain-ошибки в gRPC status / HTTP response.

### Error propagation

```
storage/    → NotFoundError / ConflictError / fmt.Errorf("op: %w", err)
services/   → fmt.Errorf("op: %w", err)   [не логирует]
transport/  → domain error → gRPC status / HTTP ответ + request-level log
```

### Trace propagation

`TraceID` и `RequestID` хранятся в `context.Context` через `pkg/trace`. Пробрасываются в gRPC через metadata headers `x-trace-id` / `x-request-id`.

## Правила зависимостей

```
transport/    → services/ (через interfaces)
services/     → models/ + consumer-side interfaces
models/       → ничего (только stdlib)
storage/      → models/ (реализует interfaces из services/)
app/          → все слои (единственное место сборки)
```

- ✅ `transport/` зависит от `services/` через интерфейсы
- ✅ Интерфейсы объявляются в `services/`, а не в `storage/`
- ❌ `transport/` не может напрямую вызывать `storage/`
- ❌ `models/` не импортирует `storage/`, `services/` или `transport/`
- ❌ Сервисы не импортируют код друг друга через `internal/` — только через gRPC-контракты

## See Also

- [C4 Диаграммы](c4.md) — визуальная топология: Context, Container, Component (включая [легенду нотации](c4.md#легенда-нотации) и компонентные разрезы всех сервисов, в том числе скелетов)
- [Sequence Diagrams](diagrams.md) — sequence-диаграммы по всем 15 endpoints, сгруппированные по доменам (Health · Auth · Users · Accounts · Transactions · Ledger)
