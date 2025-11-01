.PHONY: build run dev docs test lint lint-fix docker-up docker-down docker-restart clean migrate help

# Build the application
build:
	go build -o bin/app main.go

# Run the application
run:
	go run main.go

# Development mode with hot reload (requires air)
dev:
	air

# Generate docs api with swagger
docs:
	swag init	

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	go fmt ./...
	go vet ./...
	swag fmt -d .
	golangci-lint run -v --fast

# Format and fix code with golangci-lint
lint-fix:
	golangci-lint run --fix ./...

# Tidy dependencies
tidy:
	go mod tidy
	go mod download

# Install tools
install-tools:
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/air-verse/air@latest

# Start docker-compose services
docker-up:
	docker-compose -f docker-compose-dev.yaml up -d

# Stop docker-compose services
docker-down:
	docker-compose -f docker-compose-dev.yaml down

# Restart docker-compose services
docker-restart:
	docker-compose -f docker-compose-dev.yaml restart

# View docker logs
docker-logs:
	docker-compose -f docker-compose-dev.yaml logs -f

# Run database migrations
migrate:
	@echo "Running migrations..."
	@if [ -f .env ]; then \
		export $$(cat .env | xargs) && \
		psql -h $$DB_HOST -U $$DB_USER -d $$DB_NAME -f migrations/001_add_auth_fields.sql; \
	else \
		echo "Error: .env file not found"; \
		exit 1; \
	fi

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Show help
help:
	@echo "  make test           - Run tests"
	@echo "  make lint           - Run linter"
	@echo "  make lint-fix       - Format and fix code with golangci-lint"
	@echo "  make docker-up      - Start docker-compose services"
	@echo "  make dev            - Run in development mode"
	@echo "  make test           - Run tests"
	@echo "  make lint           - Run linter"
	@echo "  make docker-up      - Start docker-compose services"
	@echo "  make docker-down    - Stop docker-compose services"
	@echo "  make docker-restart - Restart docker-compose services"
	@echo "  make clean          - Clean build artifacts"