#!/bin/bash

# Development setup script
# Automatically sets up the entire development environment

set -e

echo "🚀 Starting development environment..."

# Check if docker compose is running
if ! docker compose ps | grep -q "Up"; then
    echo "📦 Starting Docker containers..."
    docker compose up -d --build
else
    echo "✅ Docker containers already running"
fi

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL..."
timeout 30 bash -c 'until docker compose exec -T db pg_isready -U postgres > /dev/null 2>&1; do sleep 1; done' || {
    echo "❌ PostgreSQL failed to start"
    exit 1
}

echo "✅ PostgreSQL is ready"

# Apply migrations
echo "📊 Applying database migrations..."
if command -v migrate &> /dev/null; then
    # Use migrate CLI if available
    migrate -path backend/migrations -database "postgres://postgres:postgres@localhost:5432/pwp_db?sslmode=disable" up
else
    # Fallback: use docker exec with golang-migrate in container
    docker compose exec -T backend sh -c "
        if ! command -v migrate &> /dev/null; then
            echo 'Installing migrate...'
            go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
        fi
        /root/go/bin/migrate -path /app/migrations -database 'postgres://postgres:postgres@db:5432/pwp_db?sslmode=disable' up
    "
fi

echo "✅ Migrations applied"

# Optional: Apply seed data
if [ -f "backend/migrations/seeds.sql" ]; then
    echo "🌱 Applying seed data..."
    docker compose exec -T db psql -U postgres -d pwp_db < backend/migrations/seeds.sql 2>/dev/null || echo "⚠️  Seed data already applied or not available"
fi

echo ""
echo "✨ Development environment is ready!"
echo ""
echo "📋 Services:"
echo "   - Backend:  http://localhost:8080"
echo "   - Frontend: http://localhost:5173"
echo "   - Database: postgresql://localhost:5432/pwp_db"
echo ""
echo "📝 Useful commands:"
echo "   - View logs:    docker compose logs -f"
echo "   - Stop all:     docker compose down"
echo "   - Restart:      docker compose restart"
echo ""
