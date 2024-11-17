.PHONY: all build test clean docker-up docker-down run migrate

# Variables
POSTGRES_PORT=11332
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=go_api
REDIS_PORT=6379

# Docker compose file
docker-up:
	@docker run --name postgres-go-api -e POSTGRES_USER=$(POSTGRES_USER) \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		-e POSTGRES_DB=$(POSTGRES_DB) \
		-p $(POSTGRES_PORT):5432 \
		-d postgres:latest

	@docker run --name redis-go-api \
		-p $(REDIS_PORT):6379 \
		-d redis:latest

docker-down:
	@docker stop postgres-go-api redis-go-api || true
	@docker rm postgres-go-api redis-go-api || true

# Build the application
build:
	@go build -o bin/api cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	@rm -rf bin/

# Show help
help:
	@echo "Available commands:"
	@echo "  make docker-up    - Start PostgreSQL and Redis containers"
	@echo "  make docker-down  - Stop and remove containers"
	@echo "  make build       - Build the application"
	@echo "  make run         - Run the application"
	@echo "  make test        - Run tests"
	@echo "  make clean       - Clean build artifacts"

# Default target
all: clean build
