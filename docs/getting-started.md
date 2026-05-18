[Back to README](../README.md) · [Архитектура →](architecture.md)

# Начало работы

## Требования

| Инструмент | Версия | Зачем |
|-----------|--------|-------|
| Go | 1.25.4+ | Сборка и запуск сервисов |
| Docker + Docker Compose | 24+ | Запуск PostgreSQL и сервисов |
| Make | любая | Команды сборки |
| protoc | 3.x | Регенерация gRPC-кода (опционально) |

## Запуск с Docker

Каждый реализованный сервис запускается одной командой. Docker Compose поднимает всю необходимую инфраструктуру (PostgreSQL, migrations).

### auth-service

```bash
cd internal/auth-service
make up     # docker compose up -d: Postgres + goose migrations + gRPC server
make down   # остановить
```

Сервис доступен на `localhost:8082` (gRPC).

### gateway-service

```bash
cd internal/gateway-service
make up     # docker compose up -d: HTTP server
make down   # остановить
```

Сервис доступен на `localhost:8081` (HTTP).

### account-service

```bash
cd internal/account-service
make up     # docker compose up -d: Postgres + migrations + gRPC server
make down   # остановить
```

Сервис доступен на `localhost:8083` (gRPC).

### transaction-service

```bash
cd internal/transaction-service
make up     # docker compose up -d: Postgres + migrations + gRPC server + Kafka producer
make down   # остановить
```

Сервис доступен на `localhost:8084` (gRPC). Требует запущенного Kafka для публикации событий.

### ledger-service

```bash
cd internal/ledger-service
make up     # docker compose up -d: Postgres + migrations + gRPC server + Kafka consumer
make down   # остановить
```

Сервис доступен на `localhost:8085` (gRPC) и потребляет топик `transactions.completed`. Требует запущенного Kafka.

## Запуск локально (без Docker)

Требует запущенного PostgreSQL. Переменные окружения загружаются из `local.env` через godotenv.

```bash
# auth-service
cd internal/auth-service
go run ./main.go application

# gateway-service (в отдельном терминале)
cd internal/gateway-service
go run ./main.go application

# account-service
cd internal/account-service
go run ./main.go application

# transaction-service
cd internal/transaction-service
go run ./main.go application

# ledger-service: запускает и gRPC-сервер, и Kafka-consumer
cd internal/ledger-service
go run ./main.go grpc      # только gRPC
go run ./main.go consumer  # только Kafka consumer
```

Файл `local.env` должен находиться в рабочей директории при запуске. Смотри [Конфигурацию](configuration.md) для списка переменных.

## Проверка работоспособности

```bash
# Health check
curl http://localhost:8081/ping
# Ответ: {"message":"pong"}

# Создать пользователя
curl -X POST http://localhost:8081/users \
  -H "Content-Type: application/json" \
  -d '{"login":"admin","password":"secret123","role_codes":["customer"]}'
```

## Регенерация кода

```bash
# gRPC stubs из pkg/proto/auth/auth.proto
make proto          # из корня репозитория

# HTTP handlers из internal/gateway-service/api/openapi.yaml
cd internal/gateway-service
make generate
```

Сгенерированные файлы (`pkg/proto/generated/`, `internal/gateway-service/api/ogen/`) нельзя редактировать вручную.

## Линтер

```bash
# Из корня репозитория
make install-lint   # установить golangci-lint в ./bin
make lint           # проверка
make lint-fix       # автоисправление
```

## Структура репозитория

```text
banking-platform/
├── internal/
│   ├── auth-service/        # gRPC: аутентификация, JWT, PostgreSQL
│   ├── gateway-service/     # REST: HTTP gateway → все внутренние сервисы
│   ├── account-service/     # gRPC: счета, балансы, PostgreSQL
│   ├── transaction-service/ # gRPC: переводы, пополнения, Kafka producer
│   └── ledger-service/      # gRPC + Kafka consumer: проводки, выписки
├── pkg/
│   ├── config/              # godotenv + envconfig загрузчик
│   ├── logger/              # slog-обёртка (local/dev/prod)
│   ├── trace/               # TraceID + RequestID в context
│   ├── infrastructure/postgres/ # pgx.Pool
│   ├── grpc/interceptor/    # переиспользуемые gRPC-перехватчики
│   └── proto/
│       ├── auth/auth.proto  # gRPC-контракт auth-service
│       └── generated/       # сгенерированные stubs
├── go.mod                   # единый Go-модуль
├── Makefile                 # make proto, lint
└── .golangci.yml            # строгая конфигурация линтера
```

## See Also

- [Архитектура](architecture.md) — топология сервисов и внутренняя структура
- [Конфигурация](configuration.md) — полный список переменных окружения
- [API Reference](api-reference.md) — REST и gRPC эндпоинты
