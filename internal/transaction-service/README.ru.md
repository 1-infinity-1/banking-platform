# Transaction Service

> [English](README.md) · [Русский](README.ru.md)

Transaction Service отвечает за выполнение денежных операций между счетами.
Сервис принимает команды на перевод и пополнение, координирует списание и зачисление через Account Service, обеспечивает идемпотентность операций и публикует события об успешно завершённых транзакциях.

Основные задачи:
- переводы между счетами;
- пополнение счетов;
- оркестрация debit/credit сценариев;
- хранение статусов транзакций;
- идемпотентность операций;
- публикация событий в Kafka.

---

## gRPC API

Источник proto: `pkg/proto/transaction/transaction.proto`
Порт по умолчанию: **8084**

### TransactionService

| Метод | Запрос | Ответ | Описание |
|-------|--------|-------|----------|
| `Transfer` | `TransferRequest` (from_account_id, to_account_id, amount, currency, idempotency_key) | `Transaction` | Перевод между счетами по саге Pending → Completed |
| `Replenish` | `ReplenishRequest` (to_account_id, amount, currency, idempotency_key) | `Transaction` | Пополнение счёта по саге Pending → Completed |
| `GetHistory` | `GetHistoryRequest` (account_id, limit, offset) | `TransactionsList` | История операций по счёту (DESC по created_at) |
| `GetTransaction` | `GetTransactionRequest` (transaction_id) | `Transaction` | Детали конкретной транзакции |

### Ключевые типы сообщений

| Сообщение | Поля |
|-----------|------|
| `Transaction` | id (UUID), from_account_id (UUID, optional), to_account_id (UUID), amount (decimal string), currency, status, idempotency_key, created_at, updated_at |
| `TransactionStatus` | UNSPECIFIED, PENDING, COMPLETED, FAILED, CANCELLED |
| `TransactionsList` | transactions[] |

### Маппинг ошибок (`UnaryErrorInterceptor`)

| Доменная ошибка | gRPC code |
|-----------------|-----------|
| `NotFoundError` | `NotFound` |
| `InvalidParamsError` | `InvalidArgument` |
| `BusinessError` | `FailedPrecondition` |
| `ConflictError` | `AlreadyExists` |
| прочие | `Internal` (логируется один раз с trace_id/request_id) |

---

## Сага (Pending → Completed)

`Transfer` и `Replenish` исполняют один и тот же скелет; `Transfer` добавляет шаг `Debit`:

1. Валидация запроса (положительная сумма, непустые currency / idempotency_key, разные счета для Transfer).
2. `INSERT` строки со статусом `pending`. Повторные вызовы с тем же `idempotency_key` ловит unique constraint, репозиторий делает SELECT и возвращает существующую запись — сага замыкается:
   - найден `completed` → best-effort re-publish, возвращаем существующую транзакцию;
   - найден `failed` → возвращаем `BusinessError`.
3. Только `Transfer`: вызов `account-service.Debit(from, amount, "<tx-uuid>:debit")`. На ошибке — помечаем строку `failed` и возвращаем ошибку.
4. Вызов `account-service.Credit(to, amount, "<tx-uuid>:credit")`. На ошибке — помечаем строку `failed` и возвращаем ошибку.
5. Обновление статуса строки на `completed`.
6. Публикация `TransactionEvent` в Kafka-топик `transactions.completed`.

Idempotency-ключи debit/credit формируются из UUID транзакции, поэтому ретраи внутри одной саги дедуплицируются на стороне Account Service.

---

## Kafka

| Роль | Топик | Описание |
|------|-------|----------|
| Producer | `transactions.completed` | Публикуется после успешной транзакции; потребляется Ledger Service |

Payload `TransactionEvent` (JSON):

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

## Конфигурация

Загружается из `local.env` (локальная разработка) или `docker.env` (Docker), префикс `TRANSACTION`.

| ENV | По умолчанию | Описание |
|-----|--------------|----------|
| `TRANSACTION_LOG_LEVEL` | `info` | Уровень логирования (debug / info / warn / error) |
| `TRANSACTION_DB_HOST` | `localhost` | Хост PostgreSQL |
| `TRANSACTION_DB_PORT` | `5432` | Порт PostgreSQL |
| `TRANSACTION_DB_USER` | `postgres` | Пользователь PostgreSQL |
| `TRANSACTION_DB_PASSWORD` | `postgres` | Пароль PostgreSQL |
| `TRANSACTION_DB_NAME` | `app_db` | Имя базы данных PostgreSQL |
| `TRANSACTION_GRPC_PORT` | `8084` | Порт gRPC-сервера |
| `TRANSACTION_ACCOUNT_SERVICE_HOST` | `localhost` | gRPC-хост Account Service |
| `TRANSACTION_ACCOUNT_SERVICE_PORT` | `8083` | gRPC-порт Account Service |
| `TRANSACTION_KAFKA_BROKERS` | `localhost:9092` | Брокеры Kafka (через запятую) |
| `TRANSACTION_KAFKA_TOPIC` | `transactions.completed` | Топик для событий о завершённых транзакциях |

---

## Запуск локально

```bash
make up    # docker compose up -d  (PostgreSQL + goose migrate + app)
make down  # docker compose down
```

Запуск из исходников:

```bash
cd internal/transaction-service
go run ./main.go application
```

Чтобы перевод/пополнение прошли end-to-end, account-service и брокер Kafka должны быть доступны по указанным адресам.

---

## База данных

PostgreSQL 17 (порт `5434` в docker compose, проброшен на внутренний `5432`). Миграции через [Goose](https://github.com/pressly/goose).
Файлы миграций: `migrations/`

| Таблица | Описание |
|---------|----------|
| `transactions` | Операции перевода и пополнения: `from_account_id` (UUID, NULL для пополнений), `to_account_id` (UUID), `amount` (NUMERIC(20,8)), `currency`, `status`, `idempotency_key` (UNIQUE), таймстемпы |

Индексы:

| Индекс | Колонки |
|--------|---------|
| `transactions_to_account_idx` | `to_account_id` |
| `transactions_from_account_idx` | `from_account_id` (partial, non-NULL) |

---

## Зависимости

| Направление | Сервис | Транспорт |
|-------------|--------|-----------|
| Вызывается из | Gateway Service | gRPC |
| Вызывает | Account Service (`Debit`, `Credit`) | gRPC |
| Публикует в | Ledger Service через `transactions.completed` | Kafka |

---

## Структура проекта

```text
internal/transaction-service/
├── cmd/
│   ├── application.go              # связка cobra-команд
│   └── cmd.go
├── internal/
│   ├── app/
│   │   ├── application.go          # сборка зависимостей (composition root)
│   │   └── grpc/
│   │       └── grpc.go             # жизненный цикл gRPC-сервера
│   ├── clients/
│   │   └── account/
│   │       └── client.go           # gRPC-клиент к account-service
│   ├── config/
│   │   └── config.go               # envconfig (префикс TRANSACTION)
│   ├── kafka/
│   │   └── producer.go             # producer на segmentio/kafka-go
│   ├── models/
│   │   ├── errors.go               # типизированные доменные ошибки
│   │   └── transaction.go          # доменные типы + TransactionEvent
│   ├── services/
│   │   └── transaction/
│   │       └── service.go          # сага Transfer/Replenish, GetHistory, GetTransaction
│   ├── storage/
│   │   ├── transaction/
│   │   │   ├── dto.go              # row DTO + ToDomain
│   │   │   └── repository.go       # pgx-репозиторий (идемпотентный INSERT)
│   │   └── tx/
│   │       └── tx_manager.go       # обёртка BeginFunc для транзакций
│   └── transport/
│       └── grpc/
│           ├── get_history.go
│           ├── get_transaction.go
│           ├── interceptor/
│           │   └── error.go        # маппинг доменных ошибок в gRPC status
│           ├── mapping.go          # helpers domain → proto
│           ├── replenish.go
│           ├── server.go           # serverAPI + интерфейс TransactionService
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
