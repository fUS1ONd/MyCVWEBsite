# Personal Web Platform (CV + AI Blog)

Персональный веб-сайт с публичной CV-страницей и AI блогом для авторизованных пользователей.

## 🏗 Архитектура

Проект организован как **монорепозиторий**:

```
curriculum_vitae/
├── backend/          # Go API сервер
├── frontend/         # React SPA приложение
├── docker-compose.yml
├── Makefile
└── plan.md           # Детальный план разработки
```

## 🛠 Технологический стек

### Backend
- **Go** 1.25 с Chi router
- **PostgreSQL** 15 (база данных)
- **pgx/v5** (PostgreSQL драйвер)
- **goth** (OAuth аутентификация: VK, Google, GitHub)
- **Clean Architecture** (domain, repository, service, transport layers)

### Frontend
- **React** 18+ с TypeScript
- **Vite** (build tool)
- **Tailwind CSS** (стилизация)
- **React Router** (навигация)
- **React Query** (state management)

### Infrastructure
- **Docker** + Docker Compose
- **Nginx** (reverse proxy для production)
- **Let's Encrypt** (SSL сертификаты)

## 🚀 Быстрый старт

### Предварительные требования

- Docker и Docker Compose
- Make
- Go 1.25 (опционально, для локальной разработки)
- Node.js 20+ (опционально, для локальной разработки)

### Автоматический запуск (рекомендуется)

Запустить **весь проект одной командой**:

```bash
make dev
```

Эта команда автоматически:
- ✅ Запустит Docker контейнеры (PostgreSQL, Backend, Frontend)
- ✅ Применит миграции базы данных
- ✅ Заполнит тестовыми данными (если есть)
- ✅ Настроит окружение для разработки

После запуска:
- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080
- **PostgreSQL**: localhost:5432

### Альтернативные способы

#### Вариант 1: Пошаговый запуск

```bash
# 1. Запустить Docker контейнеры
make docker-up

# 2. Применить миграции
make migrate-up

# 3. (Опционально) Заполнить тестовыми данными
make seed
```

#### Вариант 2: Локальная разработка (без Docker)

```bash
# Запустить только PostgreSQL через Docker
docker compose up db -d

# Применить миграции
make migrate-up

# Запустить backend (в одном терминале)
make backend-run

# Запустить frontend (в другом терминале)
make frontend-dev
```

### Проверка работоспособности

После запуска проверьте:

```bash
# Backend health
curl http://localhost:8080/health

# API профиля
curl http://localhost:8080/api/v1/profile

# Откройте в браузере
open http://localhost:5173
```

## 🔍 Code Quality & Linting

Проект настроен с автоматическими проверками качества кода на двух уровнях:

### Pre-commit hooks (локальная проверка)

Pre-commit hooks автоматически запускаются перед каждым коммитом и проверяют код.

**Что проверяется:**
- **Backend:** `go fmt`, `go vet`, `go mod tidy`
- **Frontend:** `prettier`, `eslint --fix`
- **Общее:** trailing whitespace, большие файлы, merge конфликты

**Установка pre-commit (выберите один способ):**

**Вариант 1: Через системный менеджер пакетов (рекомендуется)**
```bash
# Ubuntu/Debian
sudo apt install pre-commit

# macOS
brew install pre-commit

# Arch Linux
sudo pacman -S pre-commit
```

**Вариант 2: Через pip (глобально, без виртуального окружения)**
```bash
pip install pre-commit
# или
pip3 install --user pre-commit
```

**Вариант 3: Через pipx (изолированная установка CLI инструментов)**
```bash
# Установить pipx если его нет
pip install --user pipx
pipx ensurepath

# Установить pre-commit через pipx
pipx install pre-commit
```

**Активация hooks в проекте:**
```bash
# После установки pre-commit любым способом
cd /path/to/curriculum_vitae
pre-commit install

# Теперь hooks запускаются автоматически при git commit!
```

**Ручной запуск:**
```bash
# Запустить на всех файлах
pre-commit run --all-files

# Запустить на staged файлах
pre-commit run
```

### Линтеры

**Backend (golangci-lint):**
```bash
cd backend
golangci-lint run

# Конфигурация: backend/.golangci.yml
# Включено 20+ линтеров: gosec, errcheck, govet, revive, и другие
```

**Frontend (ESLint + Prettier):**
```bash
cd frontend

# Проверка кода
npm run lint              # ESLint
npm run format:check      # Prettier check

# Автоматическое исправление
npm run lint:fix          # ESLint --fix
npm run format            # Prettier --write
```

### CI/CD (GitHub Actions)

Автоматические проверки запускаются на GitHub при:
- Push в `main` или `develop`
- Создании Pull Request

**Backend CI** (`.github/workflows/backend.yml`):
- ✅ Lint (go fmt, go vet, golangci-lint)
- ✅ Test (с PostgreSQL, race detector, coverage 70%+)
- ✅ Build

**Frontend CI** (`.github/workflows/frontend.yml`):
- ✅ Lint (prettier, eslint)
- ✅ Type Check (TypeScript)
- ✅ Test (когда будут реализованы)
- ✅ Build

## 📝 Доступные команды

### Разработка
```bash
make dev                # 🚀 Запустить ВСЁ автоматически (Docker + миграции + seed)
make dev-setup          # То же самое (алиас)
```

### Docker команды
```bash
make docker-up          # Поднять все сервисы
make docker-down        # Остановить все сервисы
docker compose logs -f  # Просмотр логов
docker compose restart  # Перезапустить сервисы
```

### База данных
```bash
make migrate-up         # Применить миграции
make migrate-down       # Откатить последнюю миграцию
make migrate-create     # Создать новую миграцию
make seed               # Заполнить тестовыми данными
```

### Backend команды
```bash
make backend-build      # Собрать backend в bin/
make backend-run        # Запустить backend локально
make backend-test       # Запустить backend тесты
make backend-lint       # Запустить golangci-lint
```

### Frontend команды
```bash
make frontend-install   # Установить npm зависимости
make frontend-dev       # Запустить dev server
make frontend-build     # Собрать production build
make frontend-lint      # Запустить ESLint и Prettier проверки
make frontend-format    # Отформатировать код
```

### Общие команды
```bash
make build              # Собрать backend и frontend
make test               # Запустить все тесты
make lint               # Запустить все линтеры
make clean              # Очистить build артефакты
```

## 🏗 Структура проекта

### Backend
```
backend/
├── cmd/
│   └── app/
│       └── main.go          # Точка входа
├── internal/
│   ├── domain/              # Domain модели
│   ├── repository/          # Работа с БД
│   ├── service/             # Бизнес-логика
│   ├── transport/http/      # HTTP handlers
│   └── pkg/                 # Общие утилиты
├── config/                  # Конфигурация
├── migrations/              # SQL миграции
└── go.mod
```

### Frontend
```
frontend/
├── src/
│   ├── components/          # React компоненты
│   ├── pages/               # Страницы
│   ├── hooks/               # Custom hooks
│   ├── services/            # API клиенты
│   ├── types/               # TypeScript типы
│   └── utils/               # Утилиты
├── public/                  # Статические файлы
└── package.json
```

## 🔧 Конфигурация

### Backend конфигурация

Создайте файл `backend/config/local.yaml` для локальной разработки:

```yaml
env: "local"
http_server:
  address: "0.0.0.0:8080"
  timeout: "4s"
  idle_timeout: "60s"
database:
  url: "postgres://postgres:postgres@localhost:5432/pwp_db?sslmode=disable"
profile:
  name: "Your Name"
  description: "Your description"
  photo_url: "/static/photo.jpg"
  activity: "Your activity"
  contacts:
    email: "your@email.com"
    github: "yourusername"
    linkedin: "yourprofile"
```

### Environment Variables

Создайте файл `.env` в корне проекта:
```bash
CONFIG_PATH=backend/config/local.yaml
DATABASE_URL=postgres://postgres:postgres@localhost:5432/pwp_db?sslmode=disable
```

## 🔧 Troubleshooting

### PostgreSQL не запускается
```bash
docker compose down -v  # Удалить volumes
make dev                # Запустить заново
```

### Миграции не применяются

**Если `make migrate-up` выдает ошибку "migrate: not found":**

```bash
# Установить migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Добавить в PATH
export PATH="$HOME/go/bin:$PATH"

# ИЛИ использовать через Docker (автоматически при make dev)
```

### Frontend не собирается
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
npm run dev
```

### Pre-commit hooks ошибка с nodeenv

**Проблема:** `IndexError: list index out of range` при запуске pre-commit

**Решение:** Используйте локальные hooks вместо mirrors. Конфигурация в проекте уже обновлена для использования локальных node_modules.

### golangci-lint не найден

**Проблема:** `golangci-lint: not found` при запуске `make lint`

**Решение:** Установите golangci-lint:
```bash
# Через go install (рекомендуется)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Добавьте в PATH (если еще нет)
export PATH=$PATH:$(go env GOPATH)/bin

# Или добавьте в ~/.bashrc или ~/.zshrc
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
```

### Порты уже заняты

**Проблема:** Ошибка `port is already allocated`

**Решение:**
```bash
# Проверить какие порты используются
docker compose ps
lsof -i :8080  # Backend
lsof -i :5173  # Frontend
lsof -i :5432  # PostgreSQL

# Остановить конфликтующие сервисы
make docker-down
```

## 🧪 Тестирование

```bash
# Backend тесты
cd backend && go test -v -cover ./...

# Frontend тесты (когда будут реализованы)
cd frontend && npm test
```

## 📚 Дополнительная документация

- [Backend README](./backend/README.md) - детальная документация по backend
- [План разработки](./plan.md) - пошаговый план с чекбоксами

## 🎯 Основные функции

### Для неавторизованных пользователей
- ✅ Публичная CV-страница с информацией о навыках и контактах

### Для авторизованных пользователей
- 🔐 OAuth авторизация (VK, Google, GitHub)
- 📝 Просмотр постов об AI
- 💬 Комментирование постов
- 🔔 Уведомления о новых постах (email/push)

### Для администратора
- ✏️ Создание и редактирование постов
- 📷 Загрузка изображений и видео
- 🗑️ Модерация комментариев
- 📊 Просмотр аналитики

## 🛣 Roadmap

См. [план разработки](./plan.md) для детального roadmap с 13 этапами разработки.

**Текущий этап:** Этап 1 - Инфраструктура и базовая настройка (9/11 задач завершено - 81%)

**Недавно завершено:**
- ✅ Pre-commit hooks для автоматической проверки кода
- ✅ Конфигурация golangci-lint для backend
- ✅ Конфигурация ESLint и Prettier для frontend
- ✅ CI/CD workflows (GitHub Actions)

## 📄 Лицензия

MIT

## 👤 Автор

[Ваше имя]

---

**Дата создания:** 2025-11-25
**Версия:** 1.0
