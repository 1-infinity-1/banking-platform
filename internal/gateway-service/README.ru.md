# Gateway Service

> [English](README.md) · [Русский](README.ru.md)

Gateway Service — единая точка входа для клиентских приложений.  
Сервис предоставляет REST API для web и mobile клиентов, выполняет маршрутизацию запросов во внутренние микросервисы, проверяет JWT-токены и адаптирует ответы под клиентский формат.

Основные задачи:
- REST API для клиентов (генерируется ogen из `api/openapi.yaml`);
- JWT Bearer-аутентификация через `SecurityHandler`;
- управление аутентификацией и сессиями через Auth Service;
- управление счетами через Account Service;
- переводы и пополнения через Transaction Service;
- получение выписки по счёту через Ledger Service.

---

## REST API

Порт по умолчанию: **8081**

REST-код генерируется через [ogen](https://github.com/ogen-go/ogen) из `api/openapi.yaml`.  
Для перегенерации после изменения спецификации — `make generate` внутри `internal/gateway-service/`.

Защищённые эндпоинты требуют заголовок `Authorization: Bearer <jwt>`.

| Метод | Путь | Авторизация | Апстрим |
|-------|------|-------------|---------|
| `GET` | `/ping` | — | — |
| `POST` | `/api/v1/auth/register` | — | auth-service |
| `POST` | `/api/v1/auth/login` | — | auth-service |
| `POST` | `/api/v1/auth/logout` | JWT | auth-service |
| `POST` | `/api/v1/auth/refresh` | — | auth-service |
| `POST` | `/api/v1/accounts` | JWT | account-service |
| `GET` | `/api/v1/users/{user_id}/accounts` | JWT | account-service |
| `GET` | `/api/v1/accounts/{account_id}` | JWT | account-service |
| `GET` | `/api/v1/accounts/{account_id}/balance` | JWT | account-service |
| `PATCH` | `/api/v1/accounts/{account_id}/status` | JWT | account-service |
| `POST` | `/api/v1/transactions/transfer` | JWT | transaction-service |
| `POST` | `/api/v1/transactions/replenish` | JWT | transaction-service |
| `GET` | `/api/v1/accounts/{account_id}/transactions` | JWT | transaction-service |
| `GET` | `/api/v1/transactions/{transaction_id}` | JWT | transaction-service |
| `GET` | `/api/v1/accounts/{account_id}/statement` | JWT | ledger-service |

---

## gRPC-клиенты

| Сервис | Префикс env | Адрес по умолчанию |
|--------|-------------|-------------------|
| auth-service | `GATEWAY_AUTH_GRPC_*` | `localhost:8082` |
| account-service | `GATEWAY_ACCOUNT_GRPC_*` | `localhost:8083` |
| transaction-service | `GATEWAY_TRANSACTION_GRPC_*` | `localhost:8084` |
| ledger-service | `GATEWAY_LEDGER_GRPC_*` | `localhost:8085` |

---

## Конфигурация

Загружается из `local.env` (локальный запуск) или `docker.env` (Docker). Префикс env: `GATEWAY`.

| Переменная | Умолчание | Описание |
|------------|-----------|----------|
| `GATEWAY_LOG_LEVEL` | `info` | Уровень логирования |
| `GATEWAY_HTTP_PORT` | `8081` | HTTP-порт |
| `GATEWAY_AUTH_GRPC_HOST` | `localhost` | хост gRPC auth-service |
| `GATEWAY_AUTH_GRPC_PORT` | `8082` | порт gRPC auth-service |
| `GATEWAY_ACCOUNT_GRPC_HOST` | `localhost` | хост gRPC account-service |
| `GATEWAY_ACCOUNT_GRPC_PORT` | `8083` | порт gRPC account-service |
| `GATEWAY_TRANSACTION_GRPC_HOST` | `localhost` | хост gRPC transaction-service |
| `GATEWAY_TRANSACTION_GRPC_PORT` | `8084` | порт gRPC transaction-service |
| `GATEWAY_LEDGER_GRPC_HOST` | `localhost` | хост gRPC ledger-service |
| `GATEWAY_LEDGER_GRPC_PORT` | `8085` | порт gRPC ledger-service |
| `GATEWAY_JWT_SECRET` | **обязательно** | HS256-секрет для валидации JWT |

---

## Зависимости

| Направление | Сервис | Транспорт |
|-------------|--------|-----------|
| Вызывается из | Web / mobile клиенты | REST |
| Вызывает | auth-service | gRPC |
| Вызывает | account-service | gRPC |
| Вызывает | transaction-service | gRPC |
| Вызывает | ledger-service | gRPC |

---

## Структура проекта

```text
internal/gateway-service/
├── api/
│   ├── openapi.yaml            # OpenAPI 3.0 спецификация (источник истины для REST)
│   └── ogen/                   # сгенерировано ogen — не редактировать вручную
├── internal/
│   ├── app/
│   │   ├── application.go      # корень композиции: соединяет клиентов, сервисы, транспорт
│   │   └── srv/                # обёртка HTTP-сервера
│   ├── clients/
│   │   ├── auth/               # gRPC-клиент → auth-service
│   │   ├── account/            # gRPC-клиент → account-service
│   │   ├── transaction/        # gRPC-клиент → transaction-service
│   │   └── ledger/             # gRPC-клиент → ledger-service
│   ├── config/                 # envconfig-структура (префикс GATEWAY)
│   ├── models/                 # доменные типы и типизированные ошибки
│   ├── services/
│   │   ├── auth/               # обёртка бизнес-логики аутентификации
│   │   ├── account/            # обёртка бизнес-логики счетов
│   │   ├── transaction/        # обёртка бизнес-логики транзакций
│   │   ├── ledger/             # обёртка бизнес-логики выписок
│   │   └── management/         # управление пользователями (CreateUser)
│   └── transport/
│       ├── handlers.go         # структура GatewayHandler и интерфейсы сервисов
│       ├── middleware/
│       │   └── auth.go         # JWTSecurityHandler (ogen SecurityHandler)
│       ├── error.go            # NewError: доменные ошибки → HTTP-ответы
│       └── *.go                # по одному файлу на эндпоинт
├── Makefile
└── main.go
```
