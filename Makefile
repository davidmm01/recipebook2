.PHONY: help dev devstack backend frontend db-setup db-reset install clean test

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

dev: ## Start both backend and frontend in parallel
	@echo "Starting RecipeBook development environment..."
	@echo "Backend will run on http://localhost:8080"
	@echo "Frontend will run on http://localhost:3000"
	@echo ""
	@$(MAKE) -j2 backend frontend

devstack: db-reset dev ## Reset database and start development environment
	@echo "Development stack started with fresh database!"

backend: ## Start backend server
	@cd backend && make run-local

frontend: ## Start frontend dev server
	@cd frontend && npm start

db-setup: ## Setup database with recipes
	@cd backend && make import-and-load
	@echo ""
	@echo "Database setup complete!"

db-reset: ## Delete old database and recreate with fresh data
	@echo "Resetting database..."
	@rm -f /tmp/recipes.db
	@cd backend && make import-and-load
	@echo ""
	@echo "Database reset complete!"

install: ## Install all dependencies
	@echo "Installing backend dependencies..."
	@cd backend && go mod download
	@echo "Installing frontend dependencies..."
	@cd frontend && npm install
	@echo ""
	@echo "All dependencies installed!"

test: ## Run all tests
	@echo "Running backend tests..."
	@cd backend && make test
	@echo "Running frontend tests..."
	@cd frontend && npm test -- --watchAll=false

clean: ## Clean all build artifacts
	@echo "Cleaning backend..."
	@cd backend && make clean
	@echo "Cleaning frontend..."
	@cd frontend && rm -rf build node_modules/.cache
	@echo "Clean complete!"

build: ## Build both frontend and backend for production
	@echo "Building backend..."
	@cd backend && make build
	@echo "Building frontend..."
	@cd frontend && npm run build
	@echo ""
	@echo "Build complete!"
	@echo "Backend binary: backend/bin/recipebook-backend"
	@echo "Frontend build: frontend/build/"

lint: ## Run linters for both projects
	@echo "Linting backend..."
	@cd backend && make lint
	@echo "Backend linting complete!"
