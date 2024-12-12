.PHONY: dev build clean test install-deps docker-build

# Default target
all: build

# Install dependencies
install-deps:
	go mod download
	go install github.com/cosmtrek/air@latest

# Run development server with hot reload using air
dev:
	air

# Build for production
build:
	go build -o bin/server

# Clean build artifacts
clean:
	rm -rf bin tmp
	go clean

# Run tests
test:
	go test -v ./...

# Build docker image
docker-build:
	docker build -t htmx-go .

# Help target
help:
	@echo "Available targets:"
	@echo "  make dev          - Run development server with hot reload"
	@echo "  make build        - Build for production"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make test         - Run tests"
	@echo "  make install-deps - Install project dependencies"
	@echo "  make docker-build - Build Docker image"
