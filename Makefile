.PHONY: build run seed test clean migrate docker

# Default target
all: build

# Build the server
build:
	go build -o bin/server ./cmd/server
	go build -o bin/seed ./cmd/seed

# Run the server
run:
	go run ./cmd/server

# Seed default users
seed:
	go run ./cmd/seed

# Run tests
test:
	go test ./... -v

# Clean build artifacts
clean:
	rm -rf bin/

# Run with hot reload (requires air: go install github.com/air-verse/air@latest)
dev:
	air -c .air.toml || go run ./cmd/server

# Docker build
docker:
	docker build -t go-webserver .

# Docker run (requires PostgreSQL running)
docker-run:
	docker run -p 3000:3000 -e DATABASE_URL=postgres://host.docker.internal:5432/go_webserver?sslmode=disable go-webserver

# Format code
fmt:
	go fmt ./...

# Lint (requires golangci-lint)
lint:
	golangci-lint run

# Get dependencies
deps:
	go mod tidy
	go mod download
