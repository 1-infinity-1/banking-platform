# Account Service

> [English](README.md) · [Русский](README.ru.md)

Account Service отвечает за банковские счета и текущие балансы.  
Сервис является source of truth для денежных остатков, управляет жизненным циклом счетов, их статусами, а также выполняет операции списания и зачисления средств.

Основные задачи:
- открытие и получение счетов;
- хранение текущих балансов;
- управление статусами счетов;
- списание средств;
- зачисление средств;
- обеспечение консистентности денежных данных.

---

## gRPC API (запланировано)

| Метод | Запрос | Ответ | Описание |
|-------|--------|-------|----------|
| `CreateAccount` | `CreateAccountRequest` | `Account` | Открытие нового банковского счёта |
| `GetUserAccounts` | `GetUserAccountsRequest` | `AccountsList` | Все счета пользователя |
| `GetAccount` | `GetAccountRequest` | `Account` | Данные конкретного счёта |
| `GetBalance` | `GetBalanceRequest` | `Balance` | Текущий баланс счёта |
| `UpdateStatus` | `UpdateStatusRequest` | `UpdateStatusResponse` | Блокировка / разблокировка счёта |
| `Debit` | `DebitRequest` | `DebitResponse` | Идемпотентное списание (вызывается из Transaction Service) |
| `Credit` | `CreditRequest` | `CreditResponse` | Идемпотентное зачисление (вызывается из Transaction Service) |

---

## Зависимости

| Направление | Сервис | Транспорт |
|-------------|--------|-----------|
| Вызывается из | Gateway Service | gRPC |
| Вызывается из | Transaction Service | gRPC |
| Вызывает | — | — |
