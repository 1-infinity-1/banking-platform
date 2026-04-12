# Auth Service

> [English](README.md) · [Русский](README.ru.md)

Auth Service отвечает за управление пользователями и авторизацию.  
Сервис хранит учетные записи пользователей, выполняет вход, выдает и проверяет токены, управляет сессиями, ролями и permissions.

Основные задачи:
- вход пользователей;
- выпуск и валидация JWT/access token;
- refresh token и пользовательские сессии;
- хранение ролей и permissions;
- предоставление данных пользователя другим сервисам через gRPC.

---

## gRPC API

Источник proto: `pkg/proto/auth/auth.proto`  
Порт по умолчанию: **8082**

### AuthService

| Метод | Запрос | Ответ | Описание |
|-------|--------|-------|----------|
| `Login` | `LoginRequest` (login, password, RequestContext) | `LoginResponse` (User, Session, TokenPair, AuthContext) | Аутентификация пользователя, открытие сессии, выдача пары токенов |
| `Logout` | `LogoutRequest` (refresh_token, RequestContext) | `Empty` | Отзыв refresh token, закрытие сессии |
| `RefreshToken` | `RefreshTokenRequest` (refresh_token, RequestContext) | `RefreshTokenResponse` (TokenPair, AuthContext) | Выдача новой пары токенов по refresh token |

### AccessManagementService

| Метод | Запрос | Ответ | Описание |
|-------|--------|-------|----------|
| `CreateUser` | `CreateUserRequest` (login, email, phone, password, role_codes[]) | `CreateUserResponse` (User) | Создание нового пользователя с назначением ролей |

### Ключевые типы сообщений

| Сообщение | Поля |
|-----------|------|
| `TokenPair` | access_token, refresh_token, token_type, access_token_expires_at, refresh_token_expires_at |
| `Session` | id, user_id, status, device, created_at, updated_at, expires_at, last_seen_at |
| `AuthContext` | user_id, session_id, role_codes[], permission_codes[] |
| `RequestContext` | trace_id, request_id, ip, user_agent, device_id, platform |

---

## Конфигурация

Загружается из `local.env` (локальная разработка) или `docker.env` (Docker), префикс `AUTH`.

| ENV | По умолчанию | Описание |
|-----|-------------|----------|
| `AUTH_LOG_LEVEL` | `info` | Уровень логирования (debug / info / warn / error) |
| `AUTH_GRPC_PORT` | `8082` | Порт gRPC-сервера |

---

## Запуск локально

```bash
make up    # docker compose up -d  (PostgreSQL + goose migrate + app)
make down  # docker compose down
```

---

## База данных

PostgreSQL 17 (порт 5432). Миграции управляются через [Goose](https://github.com/pressly/goose).  
Файлы миграций: `migrations/`

| Таблица | Описание |
|---------|----------|
| `users` | Учётные записи пользователей (login, email, phone, password_hash, статус, блокировка) |
| `devices` | Отпечатки устройств (platform, user_agent) |
| `sessions` | Пользовательские сессии со статусом, сроком жизни и last_seen_at |
| `refresh_tokens` | Хэши токенов с датой истечения и флагом отзыва |
| `roles` | Именованные роли (code, name) |
| `permissions` | Именованные разрешения (code, name) |
| `user_roles` | Связь многие-ко-многим: пользователи ↔ роли |
| `role_permissions` | Связь многие-ко-многим: роли ↔ разрешения |

---

## Зависимости

| Направление | Сервис | Транспорт |
|-------------|--------|-----------|
| Вызывается из | Gateway Service | gRPC |
| Вызывает | — | — |

---

## Структура проекта

```text
internal/auth-service/
├── cmd/
│   ├── application.go      # инициализация через cobra
│   └── cmd.go
├── internal/
│   ├── app/
│   │   ├── application.go  # сборка зависимостей
│   │   └── grpc/
│   │       └── grpc.go     # жизненный цикл gRPC-сервера
│   ├── config/
│   │   └── config.go       # конфигурационная структура
│   ├── jwt/
│   │   └── manager.go      # подпись и верификация JWT
│   ├── models/
│   │   ├── errors.go
│   │   └── user.go
│   ├── services/
│   │   └── auth/
│   │       └── service.go  # бизнес-логика
│   ├── storage/
│   │   └── user/
│   │       ├── dto.go
│   │       └── repository.go
│   └── transport/
│       └── grpc/
│           └── server.go   # реализация gRPC-обработчиков
├── migrations/
│   └── 20260330175820_init.sql
├── Dockerfile
├── Makefile
├── docker-compose.yaml
├── docker.env
├── local.env
└── main.go
```
