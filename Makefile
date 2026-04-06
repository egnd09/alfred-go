.PHONY: all build test run dev docker-build docker-up clean help

# Default target
all: build

# Build the binary
build:
	@echo "Building..."
	@go build -o bin/server ./cmd/server
	@echo "Build complete: bin/server"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run locally (requires MongoDB and Redis)
run:
	@go run ./cmd/server

# Start development environment (Docker)
dev:
	@./scripts/dev.sh

# Build Docker image
docker-build:
	@docker build -t alfred-go:latest .

# Start with Docker Compose
docker-up:
	@docker-compose up --build

# Stop Docker Compose
docker-down:
	@docker-compose down

# Clean build artifacts
clean:
	@rm -rf bin/
	@echo "Cleaned build artifacts"

# Run linter
lint:
	@golangci-lint run ./...

# Format code
fmt:
	@go fmt ./...

# Check code
check: fmt test lint

# Install dependencies
deps:
	@go mod download
	@go mod tidy

# Help
help:
	@echo "Available targets:"
	@echo "  make build        - Build the binary"
	@echo "  make test         - Run tests"
	@echo "  make run          - Run server locally"
	@echo "  make dev          - Start dev environment"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-up    - Start with Docker Compose"
	@echo "  make docker-down  - Stop Docker Compose"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make lint         - Run linter"
	@echo "  make fmt          - Format code"
	@echo "  make check        - Format, test, and lint"
	@echo "  make deps         - Install dependencies"
