# Ledger Service

> [English](README.md) · [Русский](README.ru.md)

Ledger Service отвечает за бухгалтерский и аудиторский учёт операций.  
Сервис асинхронно получает события о завершённых транзакциях, сохраняет неизменяемые записи в журнал проводок и формирует выписки по счетам за период.

**Статус:** реализован — Cobra CLI, gRPC-сервер, Kafka consumer, PostgreSQL via pgx + goose-миграции, Docker.

Основные задачи:
- потребление событий о завершённых транзакциях из Kafka;
- хранение immutable log (идемпотентно — дублирующиеся события игнорируются);
- аудит операций;
- построение выписок по счетам через gRPC;
- отделение учётной модели от transactional модели.

---

## gRPC API

| Метод | Запрос | Ответ | Описание |
|-------|--------|-------|----------|
| `GetStatement` | `GetStatementRequest` (account_id, from, to) | `Statement` | Выписка по счёту за указанный период |

Proto-определение: [`pkg/proto/ledger/ledger.proto`](../../pkg/proto/ledger/ledger.proto)

---

## Kafka

| Роль | Топик | Описание |
|------|-------|----------|
| Consumer | `transactions.completed` | Получение событий о завершённых транзакциях от Transaction Service |

---

## Зависимости

| Направление | Сервис | Транспорт |
|-------------|--------|-----------|
| Вызывается из | Gateway Service | gRPC |
| Потребляет из | `transactions.completed` | Kafka |
| Вызывает | — | — |

---

## Конфигурация

Все переменные окружения используют префикс `LEDGER_`.

| Переменная | По умолчанию | Описание |
|-----------|--------------|----------|
| `LEDGER_LOG_LEVEL` | `info` | Уровень логирования (`local`, `dev`, `prod`, `info`) |
| `LEDGER_DB_HOST` | `localhost` | Хост PostgreSQL |
| `LEDGER_DB_PORT` | `5432` | Порт PostgreSQL |
| `LEDGER_DB_USER` | `postgres` | Пользователь PostgreSQL |
| `LEDGER_DB_PASSWORD` | `postgres` | Пароль PostgreSQL |
| `LEDGER_DB_NAME` | `ledger_db` | Имя базы данных PostgreSQL |
| `LEDGER_GRPC_PORT` | `8083` | Порт gRPC-сервера |
| `LEDGER_KAFKA_BROKERS` | `localhost:9092` | Адреса Kafka-брокеров (через запятую) |
| `LEDGER_KAFKA_TOPIC` | `transactions.completed` | Kafka-топик для потребления |
| `LEDGER_KAFKA_GROUP` | `ledger-service` | ID группы Kafka-консьюмера |

---

## Локальный запуск

Требования: PostgreSQL, Kafka.

```bash
# Скопировать и настроить env-файл
cp local.env.example local.env

# Запуск gRPC-сервера
cd internal/ledger-service && go run ./main.go grpc

# Запуск Kafka consumer (в отдельном терминале)
cd internal/ledger-service && go run ./main.go consumer
```

---

## Запуск через Docker

```bash
cd internal/ledger-service

# Запустить все сервисы (gRPC, consumer, PostgreSQL, Kafka, миграции)
make up

# Остановить все сервисы
make down
```
