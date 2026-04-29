-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    login TEXT NOT NULL UNIQUE,
    email TEXT UNIQUE,
    phone TEXT UNIQUE,
    password_hash TEXT NOT NULL,
    status TEXT NOT NULL,
    failed_login_attempts INTEGER NOT NULL DEFAULT 0,
    locked_until TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE users IS
    'Auth service users. | Пользователи auth-сервиса.';

COMMENT ON COLUMN users.login IS
    'User login used for authentication. | Логин пользователя для аутентификации.';

COMMENT ON COLUMN users.email IS
    'User email address. | Email пользователя.';

COMMENT ON COLUMN users.phone IS
    'User phone number. | Телефон пользователя.';

COMMENT ON COLUMN users.password_hash IS
    'Hashed user password. | Хеш пароля пользователя.';

COMMENT ON COLUMN users.status IS
    'Current user status. | Текущий статус пользователя.';

COMMENT ON COLUMN users.failed_login_attempts IS
    'Number of consecutive failed login attempts. | Количество подряд неудачных попыток входа.';

COMMENT ON COLUMN users.locked_until IS
    'Time until the user is temporarily locked. | Момент времени, до которого пользователь временно заблокирован.';

CREATE TABLE devices (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    platform TEXT NOT NULL,
    user_agent TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE devices IS
    'User devices used for authentication. | Устройства, с которых пользователь проходит аутентификацию.';

COMMENT ON COLUMN devices.platform IS
    'Device or client platform. | Платформа устройства или клиента.';

COMMENT ON COLUMN devices.user_agent IS
    'User-Agent string of the client. | User-Agent клиента, использованный при входе.';

CREATE TABLE sessions (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id BIGINT NOT NULL REFERENCES devices(id) ON DELETE RESTRICT,
    status TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON COLUMN sessions.status IS
    'Current session status. | Текущий статус сессии.';

COMMENT ON COLUMN sessions.expires_at IS
    'Session expiration timestamp. | Момент истечения срока действия сессии.';

COMMENT ON COLUMN sessions.last_seen_at IS
    'Timestamp of last activity in the session. | Момент последней активности в рамках сессии.';

CREATE TABLE refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    session_id BIGINT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON COLUMN refresh_tokens.token_hash IS
    'Hash of the refresh token. | Хеш refresh токена.';

COMMENT ON COLUMN refresh_tokens.expires_at IS
    'Refresh token expiration timestamp. | Момент истечения срока действия refresh токена.';

COMMENT ON COLUMN refresh_tokens.revoked_at IS
    'Timestamp when the token was revoked. | Момент отзыва refresh токена.';

CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON COLUMN roles.code IS
    'Unique role code. | Уникальный код роли.';

COMMENT ON COLUMN roles.name IS
    'Human-readable role name. | Человекочитаемое имя роли.';

CREATE TABLE permissions (
    id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON COLUMN permissions.code IS
    'Unique permission code. | Уникальный код разрешения.';

COMMENT ON COLUMN permissions.name IS
    'Human-readable permission name. | Человекочитаемое имя разрешения.';

CREATE TABLE user_roles (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE role_permissions (
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id BIGINT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (role_id, permission_id)
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS devices;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd