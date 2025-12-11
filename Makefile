.PHONY: init dev down logs lint test coverage

# 1. Initial setup: copies .env (if needed), installs git-hooks
init:
	cp -n backend/config/local.yaml.example backend/config/local.yaml || true
	pre-commit install

# 2. Main command: Start everything in background + rebuild if needed
dev:
	docker compose up -d --build

# 3. Full stop + remove volumes (database reset - easier than migrate down)
reset:
	docker compose down -v

# 4. Just stop
stop:
	docker compose down

# 5. View logs
logs:
	docker compose logs -f

# 6. Quality check (run local tools)
check:
	pre-commit run --all-files

# Run all tests
test:
	cd backend && go test -v -tags=integration ./...

# Generate coverage report
coverage:
	cd backend && go test -tags=integration -coverprofile=coverage.out ./...
	cd backend && go tool cover -html=coverage.out -o coverage.html
