.PHONY: help build run test clean fmt vet tidy docker-build docker-run docker-stop lint

# Variables
APP_NAME=recipebook-backend
DOCKER_IMAGE=$(APP_NAME)
DOCKER_TAG=latest
PORT?=8080

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the application binary
	go build -o bin/$(APP_NAME) .

run: ## Run the application locally
	go run .

run-local: ## Run the application locally with dev environment variables
	# -tags fts5: Enables FTS5 (Full-Text Search) support in SQLite driver for recipe search functionality
	DB_BUCKET_NAME=recipebook2-local-dev \
	GOOGLE_APPLICATION_CREDENTIALS=./service-account.json \
	PORT=8080 \
	go run -tags fts5 .

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage report
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

tidy: ## Tidy go modules
	go mod tidy

lint: fmt vet ## Run all linters

docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: ## Run Docker container
	docker run -p $(PORT):8080 \
		-e DB_BUCKET_NAME=$(DB_BUCKET_NAME) \
		-e PORT=8080 \
		--name $(APP_NAME) \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

docker-stop: ## Stop and remove Docker container
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true

docker-logs: ## Show Docker container logs
	docker logs -f $(APP_NAME)

dev: ## Run in development mode with auto-reload (requires air)
	@which air > /dev/null || (echo "air not installed, run: go install github.com/air-verse/air@latest" && exit 1)
	air

install-dev-tools: ## Install development tools
	go install github.com/air-verse/air@latest

all: lint test build ## Run lint, test, and build
