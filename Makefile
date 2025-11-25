.PHONY: build run test lint docker-up docker-down clean backend-build backend-run backend-test frontend-install frontend-dev

APP_NAME=personal-web-platform
BACKEND_DIR=backend
FRONTEND_DIR=frontend
CMD_PATH=$(BACKEND_DIR)/cmd/app/main.go

# Backend commands
backend-build:
	cd $(BACKEND_DIR) && go build -o ../bin/$(APP_NAME) cmd/app/main.go

backend-run:
	cd $(BACKEND_DIR) && go run cmd/app/main.go

backend-test:
	cd $(BACKEND_DIR) && go test -v -cover ./...

backend-lint:
	cd $(BACKEND_DIR) && golangci-lint run

# Frontend commands
frontend-install:
	cd $(FRONTEND_DIR) && npm install

frontend-dev:
	cd $(FRONTEND_DIR) && npm run dev

frontend-build:
	cd $(FRONTEND_DIR) && npm run build

# Docker commands
docker-up:
	docker-compose up -d --build

docker-down:
	docker-compose down

# Shortcut commands
build: backend-build frontend-build

run: backend-run

test: backend-test

lint: backend-lint

clean:
	rm -rf bin/
	rm -rf $(FRONTEND_DIR)/dist
	rm -rf $(FRONTEND_DIR)/node_modules
