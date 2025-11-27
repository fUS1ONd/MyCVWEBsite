.PHONY: init dev down logs lint test clean

# 1. Первая настройка: копирует .env (если надо), ставит git-hooks
init:
	cp -n backend/config/local.yaml.example backend/config/local.yaml || true
	pre-commit install

# 2. Главная команда: Запуск всего в фоне + пересборка если надо
dev:
	docker compose up -d --build

# 3. Полная остановка + удаление томов (сброс базы в ноль - это проще, чем migrate down)
reset:
	docker compose down -v

# 4. Просто остановка
stop:
	docker compose down

# 5. Просмотр логов
logs:
	docker compose logs -f

# 6. Проверка качества (запуск локальных тулзов)
check:
	pre-commit run --all-files
