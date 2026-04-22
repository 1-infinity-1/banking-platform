-- +goose Up
-- +goose StatementBegin
INSERT INTO roles (code, name)
VALUES
    ('admin', 'Administrator'),
    ('user_base', 'Base User')
ON CONFLICT (code) DO NOTHING;

INSERT INTO permissions (code, name)
VALUES
    ('user.read', 'Read user'),
    ('user.write', 'Write user'),
    ('session.read', 'Read session'),
    ('session.write', 'Write session')
ON CONFLICT (code) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
    'user.read',
    'user.write',
    'session.read',
    'session.write'
)
WHERE r.code = 'admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
    'user.read',
    'session.read'
)
WHERE r.code = 'user_base'
ON CONFLICT (role_id, permission_id) DO NOTHING;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM role_permissions
WHERE role_id IN (
    SELECT id FROM roles WHERE code IN ('admin', 'user_base')
);

DELETE FROM permissions
WHERE code IN (
    'user.read',
    'user.write',
    'session.read',
    'session.write'
);

DELETE FROM roles
WHERE code IN ('admin', 'user_base');
-- +goose StatementEnd
