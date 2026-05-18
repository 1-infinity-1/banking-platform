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

#### Health

| Метод | Путь | Описание | Статус |
|-------|------|----------|--------|
| `GET` | `/ping` | Health check | Реализован |

#### Auth

| Метод | Путь | Описание | Статус |
|-------|------|----------|--------|
| `POST` | `/auth/login` | Аутентификация по логину и паролю | Реализован |
| `POST` | `/auth/logout` | Отзыв refresh token и закрытие сессии | Реализован |
| `POST` | `/auth/refresh` | Ротация токенов | Реализован |

#### Users

| Метод | Путь | Описание | Статус |
|-------|------|----------|--------|
| `POST` | `/users` | Создание пользователя | Реализован |

#### Accounts

| Метод | Путь | Описание | Статус |
|-------|------|----------|--------|
| `GET` | `/users/{user_id}/accounts` | Счета пользователя | Скелет |
| `POST` | `/accounts` | Открыть счёт | Скелет |
| `GET` | `/accounts/{account_id}` | Получить счёт | Скелет |
| `GET` | `/accounts/{account_id}/balance` | Баланс счёта | Скелет |
| `PATCH` | `/accounts/{account_id}/status` | Изменить статус счёта | Скелет |

#### Transactions

| Метод | Путь | Описание | Статус |
|-------|------|----------|--------|
| `POST` | `/transactions/transfer` | Перевод между счетами | Скелет |
| `POST` | `/transactions/replenish` | Пополнение счёта | Скелет |
| `GET` | `/transactions/{transaction_id}` | Получить транзакцию | Скелет |
| `GET` | `/accounts/{account_id}/transactions` | История транзакций | Скелет |

#### Ledger

| Метод | Путь | Описание | Статус |
|-------|------|----------|--------|
| `GET` | `/accounts/{account_id}/statement` | Выписка по счёту | Реализован |

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

## Trace Headers

Для трассировки запросов через сервисы используются заголовки:

| Заголовок | Тип | Описание |
|-----------|-----|----------|
| `x-trace-id` | UUID | ID трассировки (генерируется на HTTP-входе) |
| `x-request-id` | UUID | ID конкретного запроса |

В HTTP-запросах передаются как заголовки; в gRPC — через metadata.

## See Also

- [Архитектура](architecture.md) — топология сервисов и паттерны
- [Конфигурация](configuration.md) — порты и переменные окружения
- [Начало работы](getting-started.md) — запуск и проверка
