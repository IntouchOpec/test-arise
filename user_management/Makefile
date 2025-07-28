.PHONY: help build run test clean docker-build docker-run docker-stop lint format deps migrate

# Default target
help: ## Display this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Development commands
deps: ## Download dependencies
	go mod download
	go mod tidy

build: ## Build the application
	go build -o bin/main .

run: ## Run the application locally
	go run main.go

test: ## Run unit tests
	go test -v ./tests/...

test-coverage: ## Run tests with coverage
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html

lint: ## Run linting
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

format: ## Format code
	go fmt ./...

# Database commands
migrate: ## Run database migrations
	@echo "Running database migrations..."
	go run main.go --migrate-only || true

# Docker commands
docker-build: ## Build Docker image
	docker build -t user-management-api .

docker-run: ## Run with Docker Compose
	docker-compose up --build -d

docker-run-proxy: ## Run with Docker Compose including Nginx proxy
	docker-compose --profile proxy up --build -d

docker-stop: ## Stop Docker containers
	docker-compose down

docker-clean: ## Clean Docker containers and volumes
	docker-compose down -v
	docker system prune -f

docker-logs: ## View Docker logs
	docker-compose logs -f api

# Development setup
dev-setup: deps ## Set up development environment
	@echo "Setting up development environment..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@echo "Development environment ready!"

# Clean up
clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean

# API testing
api-test: ## Test API endpoints (requires running server)
	@echo "Testing API endpoints..."
	@echo "Health check:"
	curl -s http://localhost:8080/health | jq .
	@echo "\nCreating test user:"
	curl -s -X POST http://localhost:8080/api/v1/users \
		-H "Content-Type: application/json" \
		-d '{"name":"Test User","email":"test@example.com","age":25}' | jq .

# Production commands
prod-build: ## Build for production
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o bin/main .

# Documentation
docs: ## Generate documentation
	@echo "API documentation available at: http://localhost:8080/health"
	@echo "Database schema:"
	@echo "- users table with fields: id, name, email, age, phone, address, is_active, created_at, updated_at, deleted_at"
