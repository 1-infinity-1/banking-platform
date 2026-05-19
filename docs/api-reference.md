[← Развёртывание](deployment.md) · [Back to README](../README.md) · [Конфигурация →](configuration.md)

# API Reference

## REST API (gateway-service)

**Base URL:** `http://localhost:8081`

Документация генерируется из [`api/openapi.yaml`](../internal/gateway-service/api/openapi.yaml).

### Аутентификация

Защищённые эндпоинты требуют заголовок:

```
Authorization: Bearer <access_token>
```

Access token выдаётся при успешном `POST /auth/login`.

### Эндпоинты

> **Поток** — ссылка на sequence-диаграмму в [`diagrams.md`](diagrams.md). Все 15 endpoints имеют диаграмму потока.

#### Health

| Метод | Путь | Описание | Статус | Поток |
|-------|------|----------|--------|-------|
| `GET` | `/ping` | Health check | Реализован | [→](diagrams.md#get-ping) |

#### Auth

| Метод | Путь | Описание | Статус | Поток |
|-------|------|----------|--------|-------|
| `POST` | `/auth/login` | Аутентификация по логину и паролю | Реализован | [→](diagrams.md#post-authlogin) |
| `POST` | `/auth/logout` | Отзыв refresh token и закрытие сессии | Реализован | [→](diagrams.md#post-authlogout) |
| `POST` | `/auth/refresh` | Ротация токенов | Реализован | [→](diagrams.md#post-authrefresh) |

#### Users

| Метод | Путь | Описание | Статус | Поток |
|-------|------|----------|--------|-------|
| `POST` | `/users` | Создание пользователя | Реализован | [→](diagrams.md#post-users) |

#### Accounts

| Метод | Путь | Описание | Статус | Поток |
|-------|------|----------|--------|-------|
| `GET` | `/users/{user_id}/accounts` | Счета пользователя | Реализован | [→](diagrams.md#get-usersuser_idaccounts) |
| `POST` | `/accounts` | Открыть счёт | Реализован | [→](diagrams.md#post-accounts) |
| `GET` | `/accounts/{account_id}` | Получить счёт | Реализован | [→](diagrams.md#get-accountsaccount_id) |
| `GET` | `/accounts/{account_id}/balance` | Баланс счёта | Реализован | [→](diagrams.md#get-accountsaccount_idbalance) |
| `PATCH` | `/accounts/{account_id}/status` | Изменить статус счёта | Реализован | [→](diagrams.md#patch-accountsaccount_idstatus) |

#### Transactions

| Метод | Путь | Описание | Статус | Поток |
|-------|------|----------|--------|-------|
| `POST` | `/transactions/transfer` | Перевод между счетами | Реализован | [→](diagrams.md#post-transactionstransfer) |
| `POST` | `/transactions/replenish` | Пополнение счёта | Реализован | [→](diagrams.md#post-transactionsreplenish) |
| `GET` | `/transactions/{transaction_id}` | Получить транзакцию | Реализован | [→](diagrams.md#get-transactionstransaction_id) |
| `GET` | `/accounts/{account_id}/transactions` | История транзакций | Реализован | [→](diagrams.md#get-accountsaccount_idtransactions) |

#### Ledger

| Метод | Путь | Описание | Статус | Поток |
|-------|------|----------|--------|-------|
| `GET` | `/accounts/{account_id}/statement` | Выписка по счёту | Реализован | [→](diagrams.md#get-accountsaccount_idstatement) |

### Примеры запросов

**POST /auth/login**

```bash
curl -X POST http://localhost:8081/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "login": "user@example.com",
    "password": "secret123",
    "context": {"platform": "web", "user_agent": "Mozilla/5.0"}
  }'
```

Ответ `200`:
```json
{
  "user": {"id": "...", "login": "user@example.com", "status": "USER_STATUS_ACTIVE"},
  "session": {"id": "...", "status": "SESSION_STATUS_ACTIVE"},
  "tokens": {
    "access_token": "eyJ...",
    "refresh_token": "eyJ...",
    "token_type": "Bearer"
  },
  "auth_context": {"user_id": "...", "role_codes": ["customer"]}
}
```

**POST /users** (создание пользователя)

```bash
curl -X POST http://localhost:8081/users \
  -H "Content-Type: application/json" \
  -d '{
    "login": "newuser",
    "email": "newuser@example.com",
    "password": "strongpassword",
    "role_codes": ["customer"]
  }'
```

**POST /transactions/transfer**

```bash
curl -X POST http://localhost:8081/transactions/transfer \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "from_account_id": "uuid-...",
    "to_account_id": "uuid-...",
    "amount": "100.00",
    "currency": "USD",
    "idempotency_key": "unique-key-123"
  }'
```

---

## gRPC API (auth-service)

**Адрес:** `localhost:8082`
**Proto-файл:** [`pkg/proto/auth/auth.proto`](../pkg/proto/auth/auth.proto)
**Package:** `auth.v1`

### AuthService

| RPC | Описание |
|-----|----------|
| `Login(LoginRequest) → LoginResponse` | Аутентификация, создание сессии, выдача токенов |
| `Logout(LogoutRequest) → Empty` | Отзыв refresh token, закрытие сессии |
| `RefreshToken(RefreshTokenRequest) → RefreshTokenResponse` | Ротация токенов |

### AccessManagementService

| RPC | Описание |
|-----|----------|
| `CreateUser(CreateUserRequest) → CreateUserResponse` | Создание пользователя |

### Ключевые типы

```protobuf
enum UserStatus {
  USER_STATUS_ACTIVE   = 1;
  USER_STATUS_BLOCKED  = 2;
  USER_STATUS_LOCKED   = 3;
  USER_STATUS_DISABLED = 4;
}

message TokenPair {
  string access_token  = 1;
  string refresh_token = 2;
  string token_type    = 3; // Bearer
}

message AuthContext {
  string   user_id          = 1;
  string   session_id       = 2;
  repeated string role_codes       = 3;
  repeated string permission_codes = 4;
}
```

### Коды ошибок gRPC

| gRPC Status | Причина |
|-------------|---------|
| `NOT_FOUND` | Пользователь или ресурс не найден |
| `ALREADY_EXISTS` | Конфликт (например, логин уже занят) |
| `INVALID_ARGUMENT` | Невалидный запрос |
| `UNAUTHENTICATED` | Неверный или просроченный токен |
| `INTERNAL` | Внутренняя ошибка сервера |

---

## gRPC API (account-service)

**Адрес:** `localhost:8083`
**Proto-файл:** [`pkg/proto/account/account.proto`](../pkg/proto/account/account.proto)
**Package:** `account.v1`

### AccountService

| RPC | Запрос | Ответ | Описание |
|-----|--------|-------|----------|
| `CreateAccount` | `user_id`, `currency` | `Account` | Открыть новый счёт |
| `GetUserAccounts` | `user_id` | `AccountsList` | Все счета пользователя |
| `GetAccount` | `account_id` | `Account` | Данные одного счёта |
| `GetBalance` | `account_id` | `Balance` | Текущий баланс |
| `UpdateStatus` | `account_id`, `status` | `UpdateStatusResponse` | Изменить статус счёта (active / blocked / closed) |
| `Debit` | `account_id`, `amount`, `idempotency_key` | `DebitResponse` (balance_after) | Идемпотентное списание |
| `Credit` | `account_id`, `amount`, `idempotency_key` | `CreditResponse` (balance_after) | Идемпотентное зачисление |

### Коды ошибок gRPC (account-service)

| gRPC Status | Причина |
|-------------|---------|
| `NOT_FOUND` | Счёт не найден |
| `ALREADY_EXISTS` | Конфликт (duplicate idempotency_key) |
| `INVALID_ARGUMENT` | Невалидный запрос |
| `FAILED_PRECONDITION` | Счёт неактивен, недостаточно средств |
| `INTERNAL` | Внутренняя ошибка сервера |

---

## gRPC API (transaction-service)

**Адрес:** `localhost:8084`
**Proto-файл:** [`pkg/proto/transaction/transaction.proto`](../pkg/proto/transaction/transaction.proto)
**Package:** `transaction.v1`

### TransactionService

| RPC | Запрос | Ответ | Описание |
|-----|--------|-------|----------|
| `Transfer` | `from_account_id`, `to_account_id`, `amount`, `currency`, `idempotency_key` | `Transaction` | Перевод между счетами через Pending → Completed saga |
| `Replenish` | `to_account_id`, `amount`, `currency`, `idempotency_key` | `Transaction` | Пополнение счёта через Pending → Completed saga |
| `GetHistory` | `account_id`, `limit`, `offset` | `TransactionsList` | История операций по счёту (DESC by created_at) |
| `GetTransaction` | `transaction_id` | `Transaction` | Данные одной транзакции |

### Коды ошибок gRPC (transaction-service)

| gRPC Status | Причина |
|-------------|---------|
| `NOT_FOUND` | Транзакция не найдена |
| `ALREADY_EXISTS` | Конфликт (duplicate idempotency_key при статусе failed) |
| `INVALID_ARGUMENT` | Невалидный запрос |
| `FAILED_PRECONDITION` | Бизнес-ошибка (недостаточно средств, неактивный счёт) |
| `INTERNAL` | Внутренняя ошибка сервера |

---

## Trace Headers

Для трассировки запросов через сервисы используются заголовки:

| Заголовок | Тип | Описание |
|-----------|-----|----------|
| `x-trace-id` | UUID | ID трассировки (генерируется на HTTP-входе) |
| `x-request-id` | UUID | ID конкретного запроса |

В HTTP-запросах передаются как заголовки; в gRPC — через metadata.

## See Also

- [Sequence Diagrams](diagrams.md) — sequence-диаграммы по всем 15 endpoints (колонка «Поток» в таблицах выше)
- [C4 Диаграммы](c4.md) — статическая топология и компонентный разрез сервисов
- [Архитектура](architecture.md) — топология сервисов и паттерны
- [Конфигурация](configuration.md) — порты и переменные окружения
- [Начало работы](getting-started.md) — запуск и проверка
