-- +goose Up
-- +goose StatementBegin
CREATE TABLE account_operations (
    id              BIGSERIAL PRIMARY KEY,
    account_id      UUID NOT NULL REFERENCES accounts(public_id),
    op_type         TEXT NOT NULL CHECK (op_type IN ('debit', 'credit')),
    amount          NUMERIC(20, 8) NOT NULL,
    balance_after   NUMERIC(20, 8),
    idempotency_key VARCHAR(255) NOT NULL UNIQUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE account_operations IS
    'Idempotent log of debit/credit operations for account-service. | Идемпотентный журнал операций списания и зачисления для account-service.';

COMMENT ON COLUMN account_operations.account_id IS
    'Target account public_id. | UUID счёта, к которому относится операция.';

COMMENT ON COLUMN account_operations.op_type IS
    'Operation type: debit or credit. | Тип операции: списание или зачисление.';

COMMENT ON COLUMN account_operations.amount IS
    'Operation amount, always positive. | Сумма операции, всегда положительная.';

COMMENT ON COLUMN account_operations.balance_after IS
    'Account balance after the operation (filled after the accounts UPDATE in the same tx). | Баланс счёта после операции (заполняется после UPDATE accounts в той же транзакции).';

COMMENT ON COLUMN account_operations.idempotency_key IS
    'Caller-provided idempotency key (typically transaction PublicID + ":debit|:credit"). | Идемпотентный ключ от вызывающего сервиса.';

CREATE INDEX account_operations_account_id_idx ON account_operations (account_id);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS account_operations;
-- +goose StatementEnd
