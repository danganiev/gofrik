.PHONY: help up down restart logs build rebuild clean test fmt lint shell db-shell db-migrate dev stop ps

# Default target - show help
help:
	@echo "Gofrik Development (Docker-based)"
	@echo ""
	@echo "Available targets:"
	@echo "  up          - Start all services in detached mode"
	@echo "  down        - Stop and remove all containers"
	@echo "  restart     - Restart all services"
	@echo "  stop        - Stop all services without removing containers"
	@echo "  logs        - Show logs from all services (ctrl+c to exit)"
	@echo "  logs-app    - Show logs from app service only"
	@echo "  logs-db     - Show logs from postgres service only"
	@echo "  ps          - List all running containers"
	@echo ""
	@echo "  dev         - Start services in development mode with live reload"
	@echo "  build       - Build/rebuild the app Docker image"
	@echo "  rebuild     - Rebuild the app image from scratch (no cache)"
	@echo ""
	@echo "  test        - Run tests inside Docker container"
	@echo "  fmt         - Format code using gofmt inside Docker"
	@echo "  lint        - Run linter inside Docker container"
	@echo ""
	@echo "  shell       - Open a shell in the app container"
	@echo "  db-shell    - Open a psql shell in the database"
	@echo ""
	@echo "  clean       - Remove containers, volumes, and build artifacts"

# Start all services
up:
	docker-compose up -d
	@echo "Services started. App available at http://localhost:8080"

# Stop and remove all containers
down:
	docker-compose down

# Restart all services
restart:
	docker-compose restart

# Stop services without removing containers
stop:
	docker-compose stop

# Show logs
logs:
	docker-compose logs -f

# Show app logs only
logs-app:
	docker-compose logs -f app

# Show database logs only
logs-db:
	docker-compose logs -f postgres

# List running containers
ps:
	docker-compose ps

# Build/rebuild the app image
build:
	docker-compose build app

# Rebuild from scratch (no cache)
rebuild:
	docker-compose build --no-cache app

# Run tests inside Docker
test:
	docker-compose exec app go test -v ./...

# Run tests without starting services (using temporary container)
test-standalone:
	docker-compose run --rm app go test -v ./...

# Format code
fmt:
	docker-compose exec app go fmt ./...

# Format code without starting services
fmt-standalone:
	docker-compose run --rm app go fmt ./...

# Run linter (requires golangci-lint in the container)
lint:
	docker-compose exec app sh -c "if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run; else echo 'golangci-lint not installed in container'; fi"

# Open shell in app container
shell:
	docker-compose exec app sh

# Open psql shell in database
db-shell:
	docker-compose exec postgres psql -U gofrik -d gofrik

# Development mode with live reload
dev:
	@echo "Starting in development mode with live reload..."
	@echo "Code changes will trigger automatic rebuilds"
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build

# Clean everything
clean:
	docker-compose down -v
	docker system prune -f
	@echo "Cleaned up containers, volumes, and dangling images"

# Install/update Go dependencies
deps:
	docker-compose exec app go mod download
	docker-compose exec app go mod tidy

# Run a one-off command in the app container
# Usage: make exec CMD="go version"
exec:
	docker-compose exec app $(CMD)

# Run a one-off command without starting services
# Usage: make run CMD="go version"
run:
	docker-compose run --rm app $(CMD)
