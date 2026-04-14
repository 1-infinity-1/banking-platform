# Gateway Service

> [English](README.md) · [Русский](README.ru.md)

Gateway Service — единая точка входа для клиентских приложений.  
Сервис предоставляет REST API для web и mobile клиентов, выполняет маршрутизацию запросов во внутренние микросервисы, агрегирует данные, проверяет авторизацию и адаптирует ответы под клиентский формат.

Основные задачи:
- REST API для клиентов;
- аутентификация через Auth Service;
- агрегация данных пользователя и счетов;
- маршрутизация запросов в Account, Transaction и Ledger сервисы;
- единая внешняя API-граница системы.

---

## REST API

Порт по умолчанию: **8081**

REST-код генерируется через [ogen](https://github.com/ogen-go/ogen) из `api/openapi.yaml`.  
Для перегенерации после изменения спецификации — см. `Makefile`.

| Метод | Путь | Статус | Апстрим |
|-------|------|--------|---------|
| `GET` | `/ping` | реализован | — |
| `POST` | `/api/v1/auth/login` | запланирован | auth-service |
| `POST` | `/api/v1/auth/logout` | запланирован | auth-service |
| `GET` | `/api/v1/users/me` | запланирован | auth-service + account-service (агрегация) |
| `POST` | `/api/v1/accounts` | запланирован | account-service |
| `GET` | `/api/v1/accounts` | запланирован | account-service |
| `GET` | `/api/v1/accounts/{id}/balance` | запланирован | account-service |
| `POST` | `/api/v1/transactions/transfer` | запланирован | transaction-service |
| `POST` | `/api/v1/transactions/replenish` | запланирован | transaction-service |
| `GET` | `/api/v1/transactions/history` | запланирован | transaction-service |
| `GET` | `/api/v1/statements/{accountId}` | запланирован | ledger-service |

---

## gRPC-клиенты (запланировано)

| Сервис | Адрес по умолчанию |
|--------|--------------------|
| auth-service | `:8082` |
| account-service | TBD |
| transaction-service | TBD |
| ledger-service | TBD |

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
│       ├── oas_handlers_gen.go
│       ├── oas_router_gen.go
│       ├── oas_schemas_gen.go
│       └── ...
├── internal/
│   └── transport/
│       ├── handlers.go         # реализации обработчиков
│       ├── ping.go             # GET /ping
│       └── error.go            # вспомогательные функции для ошибок
├── Makefile
└── main.go
```
