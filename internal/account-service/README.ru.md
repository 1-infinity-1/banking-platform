# Account Service

> [English](README.md) · [Русский](README.ru.md)

Account Service отвечает за банковские счета и текущие балансы.
Сервис является source of truth для денежных остатков, управляет жизненным циклом счетов, их статусами, а также выполняет идемпотентные операции списания и зачисления средств.

Основные задачи:
- открытие и получение счетов;
- хранение текущих балансов;
- управление статусами счетов;
- идемпотентное списание средств;
- идемпотентное зачисление средств;
- обеспечение консистентности денежных данных.

---

## gRPC API

Источник proto: `pkg/proto/account/account.proto`
Порт по умолчанию: **8083**

### AccountService

| Метод | Запрос | Ответ | Описание |
|-------|--------|-------|----------|
| `CreateAccount` | `CreateAccountRequest` (user_id, currency) | `Account` | Открытие нового банковского счёта |
| `GetUserAccounts` | `GetUserAccountsRequest` (user_id) | `AccountsList` | Все счета пользователя |
| `GetAccount` | `GetAccountRequest` (account_id) | `Account` | Данные конкретного счёта |
| `GetBalance` | `GetBalanceRequest` (account_id) | `Balance` | Текущий баланс счёта |
| `UpdateStatus` | `UpdateStatusRequest` (account_id, status) | `UpdateStatusResponse` | Смена статуса счёта (active / blocked / closed) |
| `Debit` | `DebitRequest` (account_id, amount, idempotency_key) | `DebitResponse` | Идемпотентное списание (вызывается из Transaction Service) |
| `Credit` | `CreditRequest` (account_id, amount, idempotency_key) | `CreditResponse` | Идемпотентное зачисление (вызывается из Transaction Service) |

### Ключевые типы сообщений

| Сообщение | Поля |
|-----------|------|
| `Account` | id (UUID), user_id (UUID), currency, balance (decimal string), status, created_at, updated_at |
| `Balance` | account_id (UUID), amount (decimal string), currency |
| `AccountStatus` | UNSPECIFIED, ACTIVE, BLOCKED, CLOSED |
| `DebitResponse` / `CreditResponse` | account_id (UUID), balance_after (decimal string) |

### Маппинг ошибок (`UnaryErrorInterceptor`)

| Доменная ошибка | gRPC code |
|-----------------|-----------|
| `NotFoundError` | `NotFound` |
| `InvalidParamsError` | `InvalidArgument` |
| `BusinessError` | `FailedPrecondition` |
| `ConflictError` | `AlreadyExists` |
| прочие | `Internal` (логируется один раз с trace_id/request_id) |

---

## Архитектура

### Идемпотентность Debit / Credit

`Debit` и `Credit` идемпотентны: повторные вызовы с тем же `idempotency_key` возвращают сохранённый `balance_after` без повторного применения операции. Контракт обеспечивается таблицей `account_operations` (`idempotency_key UNIQUE`).

Каждый RPC выполняется в одной транзакции `tx.Manager.BeginFunc` в три шага:

1. `INSERT` строки в `account_operations` (`op_type`, `amount`, `idempotency_key`). На `23505` (нарушение UNIQUE) репозиторий делает `SELECT` существующей операции по `idempotency_key` и возвращает её `balance_after` — сага замыкается.
2. Атомарный условный `UPDATE accounts`:
   - **Debit:** `UPDATE accounts SET balance = balance - $amount WHERE public_id = $id AND balance >= $amount AND status = 'active' RETURNING balance`
   - **Credit:** `UPDATE accounts SET balance = balance + $amount WHERE public_id = $id AND status = 'active' RETURNING balance`
   - На ноль строк репозиторий уточняет причину дополнительным `SELECT`: счёт отсутствует → `NotFoundError`, не активен → `BusinessError("account is not active: …")`, для Debit при нехватке средств → `BusinessError("insufficient funds")`. Вся транзакция откатывается, поэтому INSERT из шага 1 тоже исчезает — повторный вызов может попробовать снова.
3. `UPDATE account_operations SET balance_after = $newBalance` для операции, затем commit.

Без `SELECT ... FOR UPDATE` и без application-level version-колонки. Конкурентность обеспечивает сам условный `UPDATE`, который атомарен в PostgreSQL.

### Распространение ошибок

- Репозиторий возвращает типизированные доменные ошибки (`NotFoundError`, `BusinessError`, `ConflictError`, `InvalidParamsError`) и не логирует.
- Service оборачивает ошибки через `fmt.Errorf("op: %w", err)` для контекста и не логирует.
- `UnaryErrorInterceptor` — единственное место, которое логирует unexpected errors (с `trace_id` / `request_id`) и мапит их на gRPC status.

---

## Конфигурация

Загружается из `local.env` (локальная разработка) или `docker.env` (Docker), префикс `ACCOUNT`.

| ENV | По умолчанию | Описание |
|-----|--------------|----------|
| `ACCOUNT_LOG_LEVEL` | `info` | Уровень логирования (debug / info / warn / error) |
| `ACCOUNT_DB_HOST` | `localhost` | Хост PostgreSQL |
| `ACCOUNT_DB_PORT` | `5432` | Порт PostgreSQL |
| `ACCOUNT_DB_USER` | `postgres` | Пользователь PostgreSQL |
| `ACCOUNT_DB_PASSWORD` | `postgres` | Пароль PostgreSQL |
| `ACCOUNT_DB_NAME` | `app_db` | Имя базы данных PostgreSQL |
| `ACCOUNT_GRPC_PORT` | `8083` | Порт gRPC-сервера |

---

## Запуск локально

```bash
make up    # docker compose up -d  (PostgreSQL + goose migrate + app)
make down  # docker compose down
```

Запуск из исходников:

```bash
cd internal/account-service
go run ./main.go application
```

Для работы сервиса требуется доступный PostgreSQL по адресу из `ACCOUNT_DB_*`.

---

## База данных

PostgreSQL 17 (порт `5433` в docker compose, проброшен на внутренний `5432`). Миграции через [Goose](https://github.com/pressly/goose).
Файлы миграций: `migrations/`

| Таблица | Описание |
|---------|----------|
| `accounts` | Банковские счета: `public_id` (UUID, UNIQUE), `user_id` (UUID), `currency`, `balance` (NUMERIC(20,8)), `status`, таймстемпы |
| `account_operations` | Идемпотентный журнал операций списания/зачисления: `account_id` (UUID, FK → `accounts.public_id`), `op_type` (`debit`/`credit`), `amount` (NUMERIC(20,8)), `balance_after` (NUMERIC(20,8)), `idempotency_key` (UNIQUE), таймстемпы |

Индексы:

| Индекс | Колонки |
|--------|---------|
| `accounts_user_id_idx` | `user_id` |
| `account_operations_account_id_idx` | `account_id` |

---

## Зависимости

| Направление | Сервис | Транспорт |
|-------------|--------|-----------|
| Вызывается из | Gateway Service | gRPC |
| Вызывается из | Transaction Service (`Debit`, `Credit`) | gRPC |
| Вызывает | — | — |

---

## Структура проекта

```text
internal/account-service/
├── cmd/
│   ├── application.go              # связка cobra-команд
│   └── cmd.go
├── internal/
│   ├── app/
│   │   ├── application.go          # сборка зависимостей (composition root)
│   │   └── grpc/
│   │       └── grpc.go             # жизненный цикл gRPC-сервера
│   ├── config/
│   │   └── config.go               # envconfig (префикс ACCOUNT)
│   ├── models/
│   │   ├── account.go              # доменные типы (Account, Balance, requests/results)
│   │   └── errors.go               # типизированные доменные ошибки
│   ├── services/
│   │   └── account/
│   │       └── service.go          # CRUD + Debit/Credit бизнес-логика
│   ├── storage/
│   │   ├── account/
│   │   │   ├── dto.go              # row DTO + ToDomain
│   │   │   └── repository.go       # pgx-репозиторий (идемпотентные Debit/Credit)
│   │   └── tx/
│   │       └── tx_manager.go       # обёртка BeginFunc для транзакций
│   └── transport/
│       └── grpc/
│           ├── create_account.go
│           ├── credit.go
│           ├── debit.go
│           ├── get_account.go
│           ├── get_balance.go
│           ├── get_user_accounts.go
│           ├── interceptor/
│           │   └── error.go        # маппинг доменных ошибок в gRPC status
│           ├── mapping.go          # хелперы domain ↔ proto
│           ├── server.go           # serverAPI + интерфейс AccountService
│           └── update_status.go
├── migrations/
│   ├── 20260517000000_init_schema.sql
│   └── 20260519000000_account_operations.sql
├── Dockerfile
├── Makefile
├── docker-compose.yaml
├── docker.env
├── local.env
└── main.go
```
