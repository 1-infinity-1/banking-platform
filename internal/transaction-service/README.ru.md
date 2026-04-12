# Transaction Service

> [English](README.md) · [Русский](README.ru.md)

Transaction Service отвечает за выполнение денежных операций между счетами.  
Сервис принимает команды на перевод и пополнение, координирует списание и зачисление через Account Service, обеспечивает идемпотентность операций и публикует события об успешно завершенных транзакциях.

Основные задачи:
- переводы между счетами;
- пополнение счетов;
- оркестрация debit/credit сценариев;
- хранение статусов транзакций;
- идемпотентность операций;
- публикация событий в Kafka.

---

## gRPC API (запланировано)

| Метод | Запрос | Ответ | Описание |
|-------|--------|-------|----------|
| `Transfer` | `TransferRequest` | `Transaction` | Перевод между счетами (паттерн Saga) |
| `Replenish` | `ReplenishRequest` | `Transaction` | Пополнение счёта |
| `GetHistory` | `GetHistoryRequest` | `TransactionsList` | История операций по счёту |
| `GetTransaction` | `GetTransactionRequest` | `Transaction` | Детали конкретной транзакции |

---

## Kafka

| Роль | Топик | Описание |
|------|-------|----------|
| Producer | `transactions.completed` | Публикуется после успешной транзакции; потребляется Ledger Service |

---

## Зависимости

| Направление | Сервис | Транспорт |
|-------------|--------|-----------|
| Вызывается из | Gateway Service | gRPC |
| Вызывает | account-service (`GetAccount`, `Debit`, `Credit`) | gRPC |
| Публикует в | `transactions.completed` | Kafka |
