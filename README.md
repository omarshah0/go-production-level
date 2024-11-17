# Go Production Level

This is a simple REST API example, written in Go and featuring:

- User authentication with JSON Web Tokens
- User management (CRUD)
- Swagger documentation
- Docker support
- PostgreSQL database
- Redis caching
- Environment variables configuration
- Dependency management with Go modules
- Unit tests and integration tests

## Getting started

### Prerequisites

- Docker
- Docker Compose
- Go

### Build and run

1.  Clone the repository
2.  Run `docker-compose up` to build and start the containers
3.  Run `make run` to start the API

## Testing

- Run `go test` to run unit tests
- Run `go test -tags=integration` to run integration tests
