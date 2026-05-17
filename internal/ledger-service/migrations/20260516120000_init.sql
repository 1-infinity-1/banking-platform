-- +goose Up
-- +goose StatementBegin
CREATE TABLE ledger_entries (
    id             BIGSERIAL    PRIMARY KEY,
    public_id      UUID         NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    transaction_id UUID         NOT NULL UNIQUE,
    account_id     UUID         NOT NULL,
    type           TEXT         NOT NULL,
    amount         NUMERIC(20,8) NOT NULL,
    currency       TEXT         NOT NULL,
    balance_after  NUMERIC(20,8) NOT NULL,
    description    TEXT,
    occurred_at    TIMESTAMPTZ  NOT NULL,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ledger_entries_account_period
    ON ledger_entries (account_id, occurred_at);

COMMENT ON TABLE ledger_entries IS
    'Immutable accounting records. | Неизменяемые бухгалтерские записи.';

COMMENT ON COLUMN ledger_entries.transaction_id IS
    'Idempotency key — unique per transaction. | Ключ идемпотентности — уникален для каждой транзакции.';

COMMENT ON COLUMN ledger_entries.type IS
    'Entry direction: credit or debit. | Направление записи: credit или debit.';

COMMENT ON COLUMN ledger_entries.amount IS
    'Transaction amount. | Сумма транзакции.';

COMMENT ON COLUMN ledger_entries.balance_after IS
    'Account balance after this entry. | Баланс счёта после данной записи.';

COMMENT ON COLUMN ledger_entries.occurred_at IS
    'Timestamp when the transaction occurred. | Момент, когда произошла транзакция.';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ledger_entries;
-- +goose StatementEnd
