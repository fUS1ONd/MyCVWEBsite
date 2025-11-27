.PHONY: build run test lint format docker-up docker-down clean backend-build backend-run backend-test backend-lint frontend-install frontend-dev frontend-build frontend-lint frontend-format migrate-up migrate-down migrate-create seed

APP_NAME=personal-web-platform
BACKEND_DIR=backend
FRONTEND_DIR=frontend
CMD_PATH=$(BACKEND_DIR)/cmd/app/main.go
MIGRATE_DIR=$(BACKEND_DIR)/migrations
DATABASE_URL?=postgres://postgres:postgres@localhost:5432/pwp_db?sslmode=disable

# Backend commands
backend-build:
	cd $(BACKEND_DIR) && go build -o ../bin/$(APP_NAME) cmd/app/main.go

backend-run:
	cd $(BACKEND_DIR) && go run cmd/app/main.go

backend-test:
	cd $(BACKEND_DIR) && go test -v -cover ./...

backend-lint:
	cd $(BACKEND_DIR) && golangci-lint run

# Database migration commands
migrate-up:
	migrate -path $(MIGRATE_DIR) -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path $(MIGRATE_DIR) -database "$(DATABASE_URL)" down

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATE_DIR) -seq $$name

seed:
	psql "$(DATABASE_URL)" -f $(MIGRATE_DIR)/seeds.sql

# Frontend commands
frontend-install:
	cd $(FRONTEND_DIR) && npm install

frontend-dev:
	cd $(FRONTEND_DIR) && npm run dev

frontend-build:
	cd $(FRONTEND_DIR) && npm run build

frontend-lint:
	cd $(FRONTEND_DIR) && npm run lint && npm run format:check

frontend-format:
	cd $(FRONTEND_DIR) && npm run format

# Docker commands
docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

# Shortcut commands
build: backend-build frontend-build

run: backend-run

test: backend-test

lint: backend-lint frontend-lint

format: frontend-format

clean:
	rm -rf bin/
	rm -rf $(FRONTEND_DIR)/dist
	rm -rf $(FRONTEND_DIR)/node_modules
