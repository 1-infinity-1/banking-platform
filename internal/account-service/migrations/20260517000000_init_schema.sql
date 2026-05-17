-- +goose Up
-- +goose StatementBegin
CREATE TABLE accounts (
    id         BIGSERIAL PRIMARY KEY,
    public_id  UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL,
    currency   VARCHAR(10) NOT NULL,
    balance    NUMERIC(20, 8) NOT NULL DEFAULT 0,
    status     TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE accounts IS
    'Account service accounts. | Счета пользователей.';

COMMENT ON COLUMN accounts.user_id IS
    'Reference to auth-service user public_id. | Ссылка на public_id пользователя из auth-service.';

COMMENT ON COLUMN accounts.currency IS
    'ISO 4217 currency code. | Код валюты по ISO 4217.';

COMMENT ON COLUMN accounts.balance IS
    'Current account balance. | Текущий баланс счёта.';

COMMENT ON COLUMN accounts.status IS
    'Account status: unspecified, active, blocked, closed. | Статус счёта.';

CREATE INDEX accounts_user_id_idx ON accounts (user_id);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS accounts;
-- +goose StatementEnd
