[← API Reference](api-reference.md) · [Back to README](../README.md)

# Конфигурация

Каждый сервис загружает переменные из `local.env` (локальный запуск) или `docker.env` (Docker Compose) через библиотеки `godotenv` + `envconfig`. Файл должен находиться в рабочей директории при запуске.

**Паттерн конфигурации:** `pkg/config.NewLoaderConfig(envFilePath, prefix).Load(&cfg)`. Каждый сервис использует свой prefix (например `AUTH_`, `GATEWAY_`).

---

## auth-service

**Prefix:** `AUTH_`
**Конфиг-файлы:** `internal/auth-service/local.env` (локально), `internal/auth-service/docker.env` (Docker)

| Переменная | По умолчанию (local) | Описание |
|-----------|---------------------|----------|
| `AUTH_LOG_LEVEL` | `info` | Уровень логирования: `local`, `dev`, `prod`, `info` |
| `AUTH_DB_HOST` | `localhost` | PostgreSQL хост |
| `AUTH_DB_PORT` | `5432` | PostgreSQL порт |
| `AUTH_DB_USER` | `postgres` | Пользователь БД |
| `AUTH_DB_PASSWORD` | `postgres` | Пароль БД |
| `AUTH_DB_NAME` | `app_db` | Имя базы данных |
| `AUTH_GRPC_PORT` | `8082` | Порт gRPC-сервера |
| `AUTH_ACCESS_TOKEN_TTL` | `15m` | Время жизни access token |
| `AUTH_REFRESH_TOKEN_TTL` | `24h` | Время жизни refresh token |
| `AUTH_SECRET_KEY_FOR_TOKEN` | `secret123` | **Секрет для подписи JWT — обязательно изменить в продакшене** |

> ⚠️ `AUTH_SECRET_KEY_FOR_TOKEN` должен быть длинной случайной строкой в продакшен-среде.

---

## gateway-service

**Prefix:** `GATEWAY_`
**Конфиг-файлы:** `internal/gateway-service/local.env` (локально), `internal/gateway-service/docker.env` (Docker)

| Переменная | По умолчанию | Описание |
|-----------|--------------|----------|
| `GATEWAY_LOG_LEVEL` | `info` | Уровень логирования: `local`, `dev`, `prod`, `info` |
| `GATEWAY_HTTP_PORT` | `8081` | Порт HTTP-сервера |
| `GATEWAY_AUTH_GRPC_HOST` | `localhost` | Хост auth-service gRPC |
| `GATEWAY_AUTH_GRPC_PORT` | `8082` | Порт auth-service gRPC |
| `GATEWAY_ACCOUNT_GRPC_HOST` | `localhost` | Хост account-service gRPC |
| `GATEWAY_ACCOUNT_GRPC_PORT` | `8083` | Порт account-service gRPC |
| `GATEWAY_TRANSACTION_GRPC_HOST` | `localhost` | Хост transaction-service gRPC |
| `GATEWAY_TRANSACTION_GRPC_PORT` | `8084` | Порт transaction-service gRPC |
| `GATEWAY_LEDGER_GRPC_HOST` | `localhost` | Хост ledger-service gRPC |
| `GATEWAY_LEDGER_GRPC_PORT` | `8085` | Порт ledger-service gRPC |
| `GATEWAY_JWT_SECRET` | `change-me-in-production` | **Секрет для валидации JWT — обязательно изменить в продакшене** |

> ⚠️ `GATEWAY_JWT_SECRET` должен совпадать с `AUTH_SECRET_KEY_FOR_TOKEN` в auth-service.
> ℹ️ `GATEWAY_LEDGER_GRPC_PORT` и `LEDGER_GRPC_PORT` синхронизированы на `8085`.

---

## account-service

**Prefix:** `ACCOUNT_`
**Конфиг-файлы:** `internal/account-service/local.env`, `internal/account-service/docker.env`

| Переменная | По умолчанию | Описание |
|-----------|--------------|----------|
| `ACCOUNT_LOG_LEVEL` | `info` | Уровень логирования |
| `ACCOUNT_DB_HOST` | `localhost` | PostgreSQL хост |
| `ACCOUNT_DB_PORT` | `5432` | PostgreSQL порт |
| `ACCOUNT_DB_USER` | `postgres` | Пользователь БД |
| `ACCOUNT_DB_PASSWORD` | `postgres` | Пароль БД |
| `ACCOUNT_DB_NAME` | `app_db` | Имя базы данных |
| `ACCOUNT_GRPC_PORT` | `8083` | Порт gRPC-сервера |

---

## transaction-service

**Prefix:** `TRANSACTION_`
**Конфиг-файлы:** `internal/transaction-service/local.env`, `internal/transaction-service/docker.env`

| Переменная | По умолчанию | Описание |
|-----------|--------------|----------|
| `TRANSACTION_LOG_LEVEL` | `info` | Уровень логирования |
| `TRANSACTION_DB_HOST` | `localhost` | PostgreSQL хост |
| `TRANSACTION_DB_PORT` | `5432` | PostgreSQL порт |
| `TRANSACTION_DB_USER` | `postgres` | Пользователь БД |
| `TRANSACTION_DB_PASSWORD` | `postgres` | Пароль БД |
| `TRANSACTION_DB_NAME` | `app_db` | Имя базы данных |
| `TRANSACTION_GRPC_PORT` | `8084` | Порт gRPC-сервера |
| `TRANSACTION_ACCOUNT_SERVICE_HOST` | `localhost` | Хост account-service |
| `TRANSACTION_ACCOUNT_SERVICE_PORT` | `8083` | Порт account-service gRPC |
| `TRANSACTION_KAFKA_BROKERS` | `localhost:9092` | Kafka broker адрес |
| `TRANSACTION_KAFKA_TOPIC` | `transactions.completed` | Топик для публикации событий |

---

## ledger-service

**Prefix:** `LEDGER_`
**Конфиг-файлы:** `internal/ledger-service/local.env`, `internal/ledger-service/docker.env`

| Переменная | По умолчанию | Описание |
|-----------|--------------|----------|
| `LEDGER_LOG_LEVEL` | `info` | Уровень логирования |
| `LEDGER_DB_HOST` | `localhost` | PostgreSQL хост |
| `LEDGER_DB_PORT` | `5432` | PostgreSQL порт |
| `LEDGER_DB_USER` | `postgres` | Пользователь БД |
| `LEDGER_DB_PASSWORD` | `postgres` | Пароль БД |
| `LEDGER_DB_NAME` | `ledger_db` | Имя базы данных (отдельная от остальных) |
| `LEDGER_GRPC_PORT` | `8085` | Порт gRPC-сервера |
| `LEDGER_KAFKA_BROKERS` | `localhost:9092` | Kafka broker адрес |
| `LEDGER_KAFKA_TOPIC` | `transactions.completed` | Топик для потребления событий |
| `LEDGER_KAFKA_GROUP` | `ledger-service` | Consumer group ID |

> ℹ️ `ledger-service` использует отдельную базу данных `ledger_db` (не `app_db`).

---

## Уровни логирования

Все сервисы используют `pkg/logger` на базе `slog`:

| Значение `LOG_LEVEL` | Формат | Debug-записи |
|---------------------|--------|--------------|
| `local` | text | ✅ |
| `dev` | JSON | ✅ |
| `prod` / `info` | JSON | ❌ |

---

## Пример local.env (auth-service)

```env
AUTH_LOG_LEVEL=local

AUTH_DB_HOST=localhost
AUTH_DB_PORT=5432
AUTH_DB_USER=postgres
AUTH_DB_PASSWORD=postgres
AUTH_DB_NAME=app_db

AUTH_GRPC_PORT=8082

AUTH_ACCESS_TOKEN_TTL=15m
AUTH_REFRESH_TOKEN_TTL=24h

AUTH_SECRET_KEY_FOR_TOKEN=my-local-dev-secret
```

## See Also

- [Начало работы](getting-started.md) — как запустить сервисы с этими переменными
- [Архитектура](architecture.md) — как конфигурация загружается через composition root
- [API Reference](api-reference.md) — порты и эндпоинты сервисов
