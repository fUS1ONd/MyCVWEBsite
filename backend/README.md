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
