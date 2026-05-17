-- +goose Up
-- +goose StatementBegin
CREATE TABLE transactions (
    id               BIGSERIAL PRIMARY KEY,
    public_id        UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    from_account_id  UUID,
    to_account_id    UUID NOT NULL,
    amount           NUMERIC(20, 8) NOT NULL,
    currency         VARCHAR(10) NOT NULL,
    status           TEXT NOT NULL DEFAULT 'pending',
    idempotency_key  VARCHAR(255) NOT NULL UNIQUE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE transactions IS
    'Transaction service transactions. | Транзакции (переводы и пополнения).';

COMMENT ON COLUMN transactions.from_account_id IS
    'Source account public_id. NULL for replenishments. | UUID счёта-отправителя. NULL для пополнений.';

COMMENT ON COLUMN transactions.to_account_id IS
    'Destination account public_id. | UUID счёта-получателя.';

COMMENT ON COLUMN transactions.amount IS
    'Transaction amount. | Сумма транзакции.';

COMMENT ON COLUMN transactions.status IS
    'Transaction status: unspecified, pending, completed, failed, cancelled. | Статус транзакции.';

COMMENT ON COLUMN transactions.idempotency_key IS
    'Client-provided idempotency key. | Ключ идемпотентности от клиента.';

CREATE INDEX transactions_to_account_idx ON transactions (to_account_id);
CREATE INDEX transactions_from_account_idx ON transactions (from_account_id) WHERE from_account_id IS NOT NULL;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
-- +goose StatementEnd
