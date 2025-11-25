# Backend - Personal Web Platform

Backend API для персонального веб-сайта (CV + AI Blog), написанный на Go.

## Технологии

- **Go** 1.22+
- **Chi** - HTTP роутер
- **pgx/v5** - PostgreSQL драйвер
- **goth** - OAuth аутентификация (VK, Google, GitHub)
- **golang-migrate** - миграции БД
- **slog** - структурированное логирование

## Структура проекта

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

## Запуск

### Локально

```bash
# Из корня проекта
make backend-run

# Или напрямую
cd backend && go run cmd/app/main.go
```

### С Docker

```bash
# Из корня проекта
make docker-up
```

## Тестирование

```bash
make backend-test
```

## Конфигурация

Конфигурация загружается из файла YAML. Путь к файлу указывается в переменной окружения `CONFIG_PATH`.

Пример конфигурации в `config/local.yaml` (для разработки).
