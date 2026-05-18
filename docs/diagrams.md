[← C4 Диаграммы](c4.md) · [Back to README](../README.md) · [База данных →](database.md)

# Sequence Diagrams

Диаграммы последовательностей ключевых use cases. Каждая диаграмма отражает реальный путь запроса через слои: HTTP middleware → transport → service → storage → ответ.

---

## Login

Аутентификация пользователя: проверка пароля, создание сессии, выдача JWT-пары.

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

---

## Refresh Token

Ротация JWT-токенов: валидация текущего refresh token, выдача новой пары.

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

---

## Logout

Отзыв refresh token и завершение сессии.

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

## Transfer (перевод между счетами)

Создание перевода: списание с одного счёта, зачисление на другой, публикация события в Kafka, запись в бухгалтерский журнал.

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

---

## Get Statement (выписка по счёту)

Получение бухгалтерской выписки по счёту за период.

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
- [База данных](database.md) — ER-схемы всех четырёх баз данных
- [API Reference](api-reference.md) — HTTP и gRPC контракты
