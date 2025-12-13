# Backend - Personal Web Platform

Go-based Backend API.

## Tech Stack

- **Core:** Go 1.25, Chi v5 (Router), Slog (Logging)
- **Database:** PostgreSQL, pgx/v5 (Pool), golang-migrate
- **Auth:** Goth (OAuth2: Google, GitHub, VK ID), Gorilla Sessions
- **Config:** Cleanenv (YAML + Env support)
- **Validation:** go-playground/validator
- **Testing:** Testcontainers (Integration tests), Testify

## Features

- **Profile:** Owner information management.
- **Blog:** Posts management, slug generation, read time calculation.
- **Social:** Nested comments, likes.
- **Security:** Rate Limiting, RBAC, CORS.
- **Admin:** Media upload, content management.

## Configuration

The backend is configured via a YAML file specified by the CONFIG_PATH environment variable (defaults to config/local.yaml).

1. File Config: Main configuration resides in config/local.yaml (for dev) or config/production.yaml (for prod).
2. Credentials: Sensitive data (DB URL, OAuth secrets, Session keys) should be placed in these YAML files.
3. Env Overrides: Critical values can be overridden by environment variables defined in the config struct (e.g., DATABASE_URL, SESSION_SECRET, GOOGLE_CLIENT_ID).

See config/README.md for detailed instructions.

## Project Structure

```
backend/
├── cmd/
│   └── app/
│       └── main.go          # Точка входа приложения
├── internal/
│   ├── domain/              # Domain модели
│   ├── repository/          # Слой работы с БД
│   ├── service/             # Бизнес-логика
│   ├── transport/           # HTTP handlers
│   │   └── http/
│   └── pkg/                 # Общие утилиты
│       └── logger/
├── config/                  # Конфигурация
├── migrations/              # SQL миграции
├── go.mod
└── Dockerfile
```
