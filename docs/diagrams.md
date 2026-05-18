[← C4 Диаграммы](c4.md) · [Back to README](../README.md)

# Sequence Diagrams

Диаграммы последовательностей по всем 15 HTTP endpoints `gateway-service`. Каждая диаграмма отражает путь запроса через слои: `HTTP middleware → transport → service → storage → ответ`.

**Условные обозначения:**

- **Реализовано** — поток подтверждён кодом сервиса и тестируется руками.
- **`[скелет]`** — топология компонентов уже отражена в коде (`internal/<service>/internal/...`), бизнес-логика ещё дописывается; диаграмма описывает **целевой** поток.
- `rect rgb(230, 240, 255)` — транзакционная граница (`BeginFunc(ctx, fn)`), всё внутри идёт в одной БД-транзакции.
- Все межсервисные стрелки `gRPC` подразумевают инжектирование `x-trace-id` / `x-request-id` в metadata (опускаем в нотах, чтобы не дублировать).

## Содержание

### Health
- [GET /ping](#get-ping)

### Auth (`auth-service` — реализован)
- [POST /auth/login](#post-authlogin)
- [POST /auth/refresh](#post-authrefresh)
- [POST /auth/logout](#post-authlogout)

### Users (`auth-service` — реализован)
- [POST /users](#post-users)

### Accounts (`account-service` — `[скелет]`)
- [POST /accounts](#post-accounts)
- [GET /accounts/{account_id}](#get-accountsaccount_id)
- [GET /users/{user_id}/accounts](#get-usersuser_idaccounts)
- [GET /accounts/{account_id}/balance](#get-accountsaccount_idbalance)
- [PATCH /accounts/{account_id}/status](#patch-accountsaccount_idstatus)

### Transactions (`transaction-service` — `[скелет]`)
- [POST /transactions/transfer](#post-transactionstransfer)
- [POST /transactions/replenish](#post-transactionsreplenish)
- [GET /transactions/{transaction_id}](#get-transactionstransaction_id)
- [GET /accounts/{account_id}/transactions](#get-accountsaccount_idtransactions)

### Ledger (`ledger-service` — реализован)
- [GET /accounts/{account_id}/statement](#get-accountsaccount_idstatement)

---

## Health

### GET /ping

**Статус:** реализован. **Назначение:** health-check для liveness/readiness проб.

```mermaid
sequenceDiagram
    autonumber
    actor Client as Клиент
    participant GW as Gateway Service<br/>:8081

    Client->>GW: GET /ping
    Note over GW: Trace MW: TraceID + RequestID<br/>Logging MW: метод/путь/статус
    GW-->>Client: 200 OK<br/>{"message":"pong"}
```

---

## Auth

### POST /auth/login

**Статус:** реализован. Аутентификация пользователя: проверка пароля, создание сессии, выдача JWT-пары.

```mermaid
sequenceDiagram
    autonumber
    actor Client as Клиент
    participant GW as Gateway Service<br/>:8081
    participant AuthC as Auth gRPC Client
    participant Auth as Auth Service<br/>:8082
    participant DB as Auth DB<br/>(PostgreSQL)
    participant JWT as JWT Manager

    Client->>GW: POST /auth/login<br/>{login, password, context}
    Note over GW: Trace MW: генерирует TraceID + RequestID<br/>Logging MW: логирует входящий запрос
    GW->>AuthC: Login(LoginRequest)
    Note over AuthC: Инжектирует x-trace-id / x-request-id<br/>в gRPC metadata
    AuthC->>Auth: gRPC Login
    Auth->>DB: SELECT user WHERE login = ?
    DB-->>Auth: user row
    Note over Auth: bcrypt.CompareHashAndPassword
    Auth->>DB: INSERT device (upsert by platform+user_agent)
    DB-->>Auth: device row

    rect rgb(230, 240, 255)
        Note over Auth,DB: BEGIN TRANSACTION
        Auth->>DB: INSERT session (user_id, device_id, expires_at)
        DB-->>Auth: session
        Auth->>JWT: GenerateRefreshToken(userID, sessionID)
        JWT-->>Auth: signed refresh JWT
        Auth->>DB: INSERT refresh_token (session_id, token_hash, expires_at)
        DB-->>Auth: refresh_token row
        Note over Auth,DB: COMMIT
    end

    Auth->>JWT: GenerateAccessToken(userID, roles, permissions)
    JWT-->>Auth: signed access JWT
    Auth-->>AuthC: LoginResponse{user, session, tokens, auth_context}
    AuthC-->>GW: domain LoginResult
    GW-->>Client: 200 OK<br/>{user, session, tokens, auth_context}
```

### POST /auth/refresh

**Статус:** реализован. Ротация JWT-токенов: валидация текущего refresh token, выдача новой пары.

```mermaid
sequenceDiagram
    autonumber
    actor Client as Клиент
    participant GW as Gateway Service<br/>:8081
    participant AuthC as Auth gRPC Client
    participant Auth as Auth Service<br/>:8082
    participant DB as Auth DB<br/>(PostgreSQL)
    participant JWT as JWT Manager

    Client->>GW: POST /auth/refresh<br/>{refresh_token}
    GW->>AuthC: RefreshToken(RefreshTokenRequest)
    AuthC->>Auth: gRPC RefreshToken
    Auth->>JWT: ParseToken(refresh_token)
    JWT-->>Auth: claims{userID, sessionID}
    Auth->>DB: SELECT refresh_token WHERE session_id = ?<br/>ORDER BY created_at DESC
    DB-->>Auth: refresh_token row
    Note over Auth: Проверить: не отозван, не просрочен,<br/>hash совпадает

    Auth->>DB: SELECT session WHERE id = ?
    DB-->>Auth: session row
    Note over Auth: Проверить статус сессии = ACTIVE

    rect rgb(230, 240, 255)
        Note over Auth,DB: BEGIN TRANSACTION
        Auth->>DB: UPDATE refresh_token SET revoked_at = NOW()
        Auth->>JWT: GenerateRefreshToken(userID, sessionID)
        JWT-->>Auth: новый refresh JWT
        Auth->>DB: INSERT refresh_token (session_id, new_token_hash, expires_at)
        Note over Auth,DB: COMMIT
    end

    Auth->>JWT: GenerateAccessToken(userID, roles, permissions)
    JWT-->>Auth: новый access JWT
    Auth-->>AuthC: RefreshTokenResponse{tokens}
    AuthC-->>GW: domain TokenPair
    GW-->>Client: 200 OK<br/>{access_token, refresh_token, token_type}
```

### POST /auth/logout

**Статус:** реализован. Отзыв refresh token и завершение сессии.

```mermaid
sequenceDiagram
    autonumber
    actor Client as Клиент
    participant GW as Gateway Service<br/>:8081
    participant AuthC as Auth gRPC Client
    participant Auth as Auth Service<br/>:8082
    participant DB as Auth DB<br/>(PostgreSQL)
    participant JWT as JWT Manager

    Client->>GW: POST /auth/logout<br/>Authorization: Bearer {access_token}<br/>{refresh_token}
    Note over GW: JWT Middleware: валидирует access token,<br/>извлекает user_id и role_codes
    GW->>AuthC: Logout(LogoutRequest)
    AuthC->>Auth: gRPC Logout
    Auth->>JWT: ParseToken(refresh_token)
    JWT-->>Auth: claims{sessionID}
    Auth->>DB: SELECT refresh_token WHERE session_id = ?
    DB-->>Auth: refresh_token row
    Note over Auth: Проверить: не отозван

    rect rgb(230, 240, 255)
        Note over Auth,DB: BEGIN TRANSACTION
        Auth->>DB: UPDATE refresh_token SET revoked_at = NOW()
        Auth->>DB: UPDATE session SET status = 'closed'
        Note over Auth,DB: COMMIT
    end

    Auth-->>AuthC: Empty
    AuthC-->>GW: ok
    GW-->>Client: 200 OK
```

---

## Users

### POST /users

**Статус:** реализован. Создание пользователя с назначением ролей.

> Сейчас эндпоинт **не требует** JWT-токена (см. `internal/gateway-service/api/openapi.yaml`). В целевом облике перед production-ready запуском будет добавлена проверка роли `admin` / permission `users:write`. Этот шаг отмечен на диаграмме как «`TODO`».

```mermaid
sequenceDiagram
    autonumber
    actor Client as Клиент
    participant GW as Gateway Service<br/>:8081
    participant AuthC as Auth gRPC Client
    participant Auth as Auth Service<br/>:8082<br/>(ManagementService)
    participant DB as Auth DB<br/>(PostgreSQL)

    Client->>GW: POST /users<br/>{login, email, password, role_codes}
    Note over GW: Trace MW + Logging MW<br/>(TODO: JWT MW + проверка admin)
    GW->>AuthC: CreateUser(CreateUserRequest)
    Note over AuthC: Инжектирует x-trace-id / x-request-id<br/>в gRPC metadata
    AuthC->>Auth: gRPC CreateUser
    Note over Auth: bcrypt.GenerateFromPassword(password)

    rect rgb(230, 240, 255)
        Note over Auth,DB: BEGIN TRANSACTION
        Auth->>DB: INSERT user<br/>(login, email, password_hash, status='ACTIVE')
        DB-->>Auth: user row
        Auth->>DB: SELECT roles WHERE code IN (?, ?, ...)
        DB-->>Auth: []role rows
        Note over Auth: если roles пусто → NotFoundError<br/>→ ROLLBACK
        Auth->>DB: INSERT user_roles (user_id, role_id) × N
        DB-->>Auth: ok
        Note over Auth,DB: COMMIT
    end

    Auth-->>AuthC: CreateUserResponse{user}
    AuthC-->>GW: domain User
    GW-->>Client: 201 Created<br/>{user: {id, login, email, status, role_codes}}
```

---

## Accounts

### POST /accounts

**Статус:** `[скелет]`. Открытие банковского счёта для пользователя.

```mermaid
sequenceDiagram
    autonumber
    actor Client as Клиент
    participant GW as Gateway Service<br/>:8081
    participant AccC as Account gRPC Client
    participant Acc as Account Service<br/>:8083 [скелет]
    participant DB as Account DB<br/>(PostgreSQL)

    Client->>GW: POST /accounts<br/>Authorization: Bearer {access_token}<br/>{user_id, currency, account_type}
    Note over GW: JWT Middleware: валидация токена,<br/>permission account:create<br/>проверка user_id == claims.user_id<br/>(или роль admin)
    GW->>AccC: CreateAccount(CreateAccountRequest)
    AccC->>Acc: gRPC CreateAccount

    rect rgb(230, 240, 255)
        Note over Acc,DB: BEGIN TRANSACTION
        Acc->>DB: INSERT account<br/>(user_id, currency, balance=0,<br/>status='ACTIVE')
        DB-->>Acc: account row
        Note over Acc,DB: COMMIT
    end

    Acc-->>AccC: CreateAccountResponse{account}
    AccC-->>GW: domain Account
    GW-->>Client: 201 Created<br/>{account: {id, user_id, balance, currency, status, created_at}}
```

### GET /accounts/{account_id}

**Статус:** `[скелет]`. Получение одного счёта.

```mermaid
sequenceDiagram
    autonumber
    actor Client as Клиент
    participant GW as Gateway Service<br/>:8081
    participant AccC as Account gRPC Client
    participant Acc as Account Service<br/>:8083 [скелет]
    participant DB as Account DB<br/>(PostgreSQL)

    Client->>GW: GET /accounts/{account_id}<br/>Authorization: Bearer {access_token}
    Note over GW: JWT MW: извлечь user_id из claims
    GW->>AccC: GetAccount(GetAccountRequest)
    AccC->>Acc: gRPC GetAccount
    Acc->>DB: SELECT * FROM accounts<br/>WHERE id = ?
    alt account not found
        DB-->>Acc: no rows
        Acc-->>AccC: NotFoundError
        AccC-->>GW: domain NotFoundError
        GW-->>Client: 404 Not Found
    else found
        DB-->>Acc: account row
        Note over Acc: Проверка владения:<br/>account.user_id == ctx.user_id<br/>(или роль admin)
        Acc-->>AccC: GetAccountResponse{account}
        AccC-->>GW: domain Account
        GW-->>Client: 200 OK<br/>{account}
    end
```

### GET /users/{user_id}/accounts

**Статус:** `[скелет]`. Список всех счетов пользователя.

```mermaid
sequenceDiagram
    autonumber
    actor Client as Клиент
    participant GW as Gateway Service<br/>:8081
    participant AccC as Account gRPC Client
    participant Acc as Account Service<br/>:8083 [скелет]
    participant DB as Account DB<br/>(PostgreSQL)

    Client->>GW: GET /users/{user_id}/accounts<br/>Authorization: Bearer {access_token}
    Note over GW: JWT MW: проверка user_id из path<br/>== claims.user_id (или admin)
    GW->>AccC: GetUserAccounts(GetUserAccountsRequest)
    AccC->>Acc: gRPC GetUserAccounts
    Acc->>DB: SELECT * FROM accounts<br/>WHERE user_id = ?<br/>ORDER BY created_at DESC
    DB-->>Acc: []account rows
    Acc-->>AccC: GetUserAccountsResponse{accounts[]}
    AccC-->>GW: []domain Account
    GW-->>Client: 200 OK<br/>{accounts: [{id, balance, currency, status, ...}, ...]}
```

### GET /accounts/{account_id}/balance

**Статус:** `[скелет]`. Возвращает только баланс и валюту — лёгкий вариант `GetAccount` для частых проверок.

```mermaid
sequenceDiagram
    autonumber
    actor Client as Клиент
    participant GW as Gateway Service<br/>:8081
    participant AccC as Account gRPC Client
    participant Acc as Account Service<br/>:8083 [скелет]
    participant DB as Account DB<br/>(PostgreSQL)

    Client->>GW: GET /accounts/{account_id}/balance<br/>Authorization: Bearer {access_token}
    Note over GW: JWT MW: извлечь user_id из claims
    GW->>AccC: GetBalance(GetBalanceRequest)
    AccC->>Acc: gRPC GetBalance
    Acc->>DB: SELECT balance, currency, user_id<br/>FROM accounts WHERE id = ?
    DB-->>Acc: row{balance, currency, user_id}
    Note over Acc: Проверка владения:<br/>user_id == ctx.user_id (или admin)
    Acc-->>AccC: GetBalanceResponse{amount, currency}
    AccC-->>GW: domain Balance
    GW-->>Client: 200 OK<br/>{amount: "1500.00", currency: "USD"}
```

### PATCH /accounts/{account_id}/status

**Статус:** `[скелет]`. Изменение статуса счёта (admin / специальный permission).

> **Допустимые переходы:** `ACTIVE → BLOCKED`, `ACTIVE → CLOSED`, `BLOCKED → ACTIVE`. Любой другой переход → `400 Bad Request` (ValidationError).

```mermaid
sequenceDiagram
    autonumber
    actor Client as Клиент
    participant GW as Gateway Service<br/>:8081
    participant AccC as Account gRPC Client
    participant Acc as Account Service<br/>:8083 [скелет]
    participant DB as Account DB<br/>(PostgreSQL)

    Client->>GW: PATCH /accounts/{account_id}/status<br/>Authorization: Bearer {access_token}<br/>{new_status: "BLOCKED"}
    Note over GW: JWT MW: проверка роли admin<br/>или permission account:status:write
    GW->>AccC: UpdateAccountStatus(req)
    AccC->>Acc: gRPC UpdateAccountStatus

    rect rgb(230, 240, 255)
        Note over Acc,DB: BEGIN TRANSACTION
        Acc->>DB: SELECT status FROM accounts<br/>WHERE id = ? FOR UPDATE
        DB-->>Acc: row{status: 'ACTIVE'}
        Note over Acc: Валидация перехода:<br/>'ACTIVE' → 'BLOCKED' допустим
        alt invalid transition
            Note over Acc,DB: ROLLBACK
            Acc-->>AccC: ValidationError
            AccC-->>GW: domain ValidationError
            GW-->>Client: 400 Bad Request
        else valid
            Acc->>DB: UPDATE accounts SET status = ?,<br/>updated_at = NOW() WHERE id = ?
            DB-->>Acc: updated row
            Note over Acc,DB: COMMIT
            Acc-->>AccC: UpdateAccountStatusResponse{account}
            AccC-->>GW: domain Account
            GW-->>Client: 200 OK<br/>{account: {..., status: 'BLOCKED', updated_at}}
        end
    end
```

---

## Transactions

### POST /transactions/transfer

**Статус:** `[скелет]`. Создание перевода: списание с одного счёта, зачисление на другой, публикация события в Kafka, запись в бухгалтерский журнал.

```mermaid
sequenceDiagram
    autonumber
    actor Client as Клиент
    participant GW as Gateway Service<br/>:8081
    participant TxC as Transaction gRPC Client
    participant Tx as Transaction Service<br/>:8084
    participant AccC as Account gRPC Client<br/>(внутри Tx Service)
    participant Acc as Account Service<br/>:8083
    participant AccDB as Account DB
    participant TxDB as Transaction DB
    participant Kafka as Kafka<br/>transactions.completed
    participant LedC as Ledger Consumer
    participant LedDB as Ledger DB

    Client->>GW: POST /transactions/transfer<br/>Authorization: Bearer {access_token}<br/>{from_account_id, to_account_id, amount, currency, idempotency_key}
    Note over GW: JWT Middleware: валидирует токен,<br/>проверяет роль / permission

    GW->>TxC: Transfer(TransferRequest)
    TxC->>Tx: gRPC Transfer

    rect rgb(230, 240, 255)
        Note over Tx,TxDB: BEGIN TRANSACTION
        Tx->>TxDB: INSERT transaction (pending, idempotency_key)
        TxDB-->>Tx: transaction row
        Tx->>AccC: DebitAccount(from_account_id, amount)
        AccC->>Acc: gRPC DebitAccount
        Acc->>AccDB: UPDATE accounts SET balance = balance - amount<br/>WHERE id = ? AND balance >= amount
        AccDB-->>Acc: updated row
        Acc-->>AccC: ok
        AccC-->>Tx: ok
        Tx->>AccC: CreditAccount(to_account_id, amount)
        AccC->>Acc: gRPC CreditAccount
        Acc->>AccDB: UPDATE accounts SET balance = balance + amount<br/>WHERE id = ?
        AccDB-->>Acc: updated row
        Acc-->>AccC: ok
        AccC-->>Tx: ok
        Tx->>TxDB: UPDATE transaction SET status = 'completed'
        Note over Tx,TxDB: COMMIT
    end

    Tx->>Kafka: Produce(transactions.completed)<br/>{transaction_id, from_account_id, to_account_id, amount, currency}
    Kafka-->>Tx: ack

    Tx-->>TxC: TransferResponse{transaction}
    TxC-->>GW: domain Transaction
    GW-->>Client: 200 OK {transaction}

    Note over Kafka,LedDB: Асинхронно (Ledger Consumer)
    Kafka->>LedC: Consume(transactions.completed)
    LedC->>LedDB: INSERT ledger_entry (debit, from_account, amount, balance_after)
    LedC->>LedDB: INSERT ledger_entry (credit, to_account, amount, balance_after)
    LedDB-->>LedC: ok
    LedC->>Kafka: commit offset
```

### POST /transactions/replenish

**Статус:** `[скелет]`. Пополнение счёта (одностороннее зачисление, без счёта-источника). Похоже на `transfer`, но без `Debit` и с проводкой только `credit` в ledger.

```mermaid
sequenceDiagram
    autonumber
    actor Client as Клиент
    participant GW as Gateway Service<br/>:8081
    participant TxC as Transaction gRPC Client
    participant Tx as Transaction Service<br/>:8084 [скелет]
    participant AccC as Account gRPC Client<br/>(внутри Tx Service)
    participant Acc as Account Service<br/>:8083 [скелет]
    participant AccDB as Account DB
    participant TxDB as Transaction DB
    participant Kafka as Kafka<br/>transactions.completed
    participant LedC as Ledger Consumer
    participant LedDB as Ledger DB

    Client->>GW: POST /transactions/replenish<br/>Authorization: Bearer {access_token}<br/>{account_id, amount, currency, idempotency_key, source}
    Note over GW: JWT MW: permission transaction:replenish<br/>(обычно admin или платёжный шлюз)

    GW->>TxC: Replenish(ReplenishRequest)
    TxC->>Tx: gRPC Replenish

    Note over Tx,TxDB: Idempotency: SELECT transaction<br/>WHERE idempotency_key = ?<br/>если есть → вернуть существующую
    Tx->>TxDB: SELECT transaction WHERE idempotency_key = ?
    TxDB-->>Tx: not found

    rect rgb(230, 240, 255)
        Note over Tx,TxDB: BEGIN TRANSACTION
        Tx->>TxDB: INSERT transaction<br/>(type='replenish', status='pending',<br/>to_account_id, amount, currency, idempotency_key)
        TxDB-->>Tx: transaction row
        Tx->>AccC: CreditAccount(account_id, amount)
        AccC->>Acc: gRPC CreditAccount
        Acc->>AccDB: UPDATE accounts SET balance = balance + amount<br/>WHERE id = ? AND status = 'ACTIVE'
        AccDB-->>Acc: updated row
        Acc-->>AccC: ok
        AccC-->>Tx: ok
        Tx->>TxDB: UPDATE transaction SET status = 'completed'
        Note over Tx,TxDB: COMMIT
    end

    Tx->>Kafka: Produce(transactions.completed)<br/>{transaction_id, type='replenish', to_account_id, amount, currency}
    Kafka-->>Tx: ack

    Tx-->>TxC: ReplenishResponse{transaction}
    TxC-->>GW: domain Transaction
    GW-->>Client: 200 OK {transaction}

    Note over Kafka,LedDB: Асинхронно (Ledger Consumer)
    Kafka->>LedC: Consume(transactions.completed)
    Note over LedC: type='replenish' → одна проводка (credit)
    LedC->>LedDB: INSERT ledger_entry<br/>(credit, to_account, amount, balance_after)
    LedDB-->>LedC: ok
    LedC->>Kafka: commit offset
```

### GET /transactions/{transaction_id}

**Статус:** `[скелет]`. Получение одной транзакции по ID.

```mermaid
sequenceDiagram
    autonumber
    actor Client as Клиент
    participant GW as Gateway Service<br/>:8081
    participant TxC as Transaction gRPC Client
    participant Tx as Transaction Service<br/>:8084 [скелет]
    participant TxDB as Transaction DB<br/>(PostgreSQL)

    Client->>GW: GET /transactions/{transaction_id}<br/>Authorization: Bearer {access_token}
    Note over GW: JWT MW: извлечь user_id из claims
    GW->>TxC: GetTransaction(GetTransactionRequest)
    TxC->>Tx: gRPC GetTransaction
    Tx->>TxDB: SELECT * FROM transactions WHERE id = ?
    alt not found
        TxDB-->>Tx: no rows
        Tx-->>TxC: NotFoundError
        TxC-->>GW: domain NotFoundError
        GW-->>Client: 404 Not Found
    else found
        TxDB-->>Tx: transaction row
        Note over Tx: Проверка владения:<br/>user является владельцем from_account_id<br/>или to_account_id (запрос в account-service,<br/>здесь упрощено)
        Tx-->>TxC: GetTransactionResponse{transaction}
        TxC-->>GW: domain Transaction
        GW-->>Client: 200 OK<br/>{transaction: {id, type, status, from_account_id, to_account_id, amount, currency, created_at, ...}}
    end
```

### GET /accounts/{account_id}/transactions

**Статус:** `[скелет]`. История транзакций по счёту с пагинацией.

```mermaid
sequenceDiagram
    autonumber
    actor Client as Клиент
    participant GW as Gateway Service<br/>:8081
    participant TxC as Transaction gRPC Client
    participant Tx as Transaction Service<br/>:8084 [скелет]
    participant TxDB as Transaction DB<br/>(PostgreSQL)

    Client->>GW: GET /accounts/{account_id}/transactions<br/>Authorization: Bearer {access_token}<br/>?from=2026-01-01&to=2026-03-31<br/>&limit=50&offset=0
    Note over GW: JWT MW: проверка владения<br/>account_id ↔ user_id
    GW->>TxC: GetTransactionHistory(GetTransactionHistoryRequest)
    TxC->>Tx: gRPC GetTransactionHistory
    Tx->>TxDB: SELECT * FROM transactions<br/>WHERE (from_account_id = ? OR to_account_id = ?)<br/>AND created_at BETWEEN ? AND ?<br/>ORDER BY created_at DESC<br/>LIMIT ? OFFSET ?
    TxDB-->>Tx: []transaction rows
    Tx->>TxDB: SELECT COUNT(*) FROM transactions<br/>WHERE (from_account_id = ? OR to_account_id = ?)<br/>AND created_at BETWEEN ? AND ?
    TxDB-->>Tx: total count
    Tx-->>TxC: GetTransactionHistoryResponse{transactions[], total}
    TxC-->>GW: []domain Transaction + pagination
    GW-->>Client: 200 OK<br/>{transactions: [...], pagination: {total, limit, offset}}
```

---

## Ledger

### GET /accounts/{account_id}/statement

**Статус:** реализован. Получение бухгалтерской выписки по счёту за период.

```mermaid
sequenceDiagram
    autonumber
    actor Client as Клиент
    participant GW as Gateway Service<br/>:8081
    participant LedC as Ledger gRPC Client
    participant Led as Ledger Service<br/>:8085
    participant LedDB as Ledger DB<br/>(PostgreSQL)

    Client->>GW: GET /accounts/{account_id}/statement<br/>Authorization: Bearer {access_token}<br/>?from=2026-01-01&to=2026-03-31
    Note over GW: JWT Middleware: валидирует access token,<br/>извлекает user_id и role_codes

    GW->>LedC: GetStatement(GetStatementRequest)
    Note over LedC: Инжектирует x-trace-id / x-request-id<br/>в gRPC metadata
    LedC->>Led: gRPC GetStatement

    Led->>LedDB: SELECT * FROM ledger_entries<br/>WHERE account_id = ?<br/>AND occurred_at BETWEEN ? AND ?<br/>ORDER BY occurred_at ASC

    LedDB-->>Led: []ledger_entry rows

    Led-->>LedC: GetStatementResponse{entries[]}
    LedC-->>GW: []domain LedgerEntry
    GW-->>Client: 200 OK<br/>{entries: [{type, amount, currency, balance_after, occurred_at, description}, ...]}
```

---

## See Also

- [C4 Диаграммы](c4.md) — статическая топология: Context, Container, Component
- [API Reference](api-reference.md) — HTTP и gRPC контракты
