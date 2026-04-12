# Account Service

> [English](README.md) · [Русский](README.ru.md)

Account Service is responsible for bank accounts and current balances.  
It is the source of truth for monetary balances, manages the account lifecycle and statuses, and performs debit and credit operations.

Main responsibilities:
- account creation and retrieval;
- current balance storage;
- account status management;
- debit operations;
- credit operations;
- consistency of monetary data.

---

## gRPC API (planned)

| Method | Request | Response | Description |
|--------|---------|----------|-------------|
| `CreateAccount` | `CreateAccountRequest` | `Account` | Open a new bank account |
| `GetUserAccounts` | `GetUserAccountsRequest` | `AccountsList` | All accounts belonging to a user |
| `GetAccount` | `GetAccountRequest` | `Account` | Single account data |
| `GetBalance` | `GetBalanceRequest` | `Balance` | Current account balance |
| `UpdateStatus` | `UpdateStatusRequest` | `UpdateStatusResponse` | Block / unblock an account |
| `Debit` | `DebitRequest` | `DebitResponse` | Idempotent debit (called by Transaction Service) |
| `Credit` | `CreditRequest` | `CreditResponse` | Idempotent credit (called by Transaction Service) |

---

## Dependencies

| Direction | Service | Transport |
|-----------|---------|-----------|
| Called by | Gateway Service | gRPC |
| Called by | Transaction Service | gRPC |
| Calls | — | — |
