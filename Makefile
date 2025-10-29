.PHONY: build run dev test lint lint-fix docker-up docker-down docker-restart clean help

# Build the application
build:
	go build -o bin/app main.go

# Run the application
run:
	go run main.go

# Development mode with hot reload (requires air)
dev:
	air

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	go fmt ./...
	go vet ./...
	swag fmt -d .
	golangci-lint run -v --fast

# Format and fix code with golangci-lint
lint-fix:
	golangci-lint run --fix ./...

# Start docker-compose services
docker-up:
	docker-compose up -d

# Stop docker-compose services
docker-down:
	docker-compose down

# Restart docker-compose services
docker-restart:
	docker-compose restart

# Clean build artifacts
clean:
	rm -rf bin/

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