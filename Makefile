.PHONY: help dev devstack backend frontend db-setup db-reset install clean test auto-promote-admin

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

devstack: db-reset ## Reset database and start development environment
	@echo "Starting RecipeBook development environment..."
	@echo "Backend will run on http://localhost:8080"
	@echo "Frontend will run on http://localhost:3000"
	@echo ""
	@$(MAKE) -j3 backend frontend auto-promote-admin

backend: ## Start backend server
	@cd backend && make run-local

frontend: ## Start frontend dev server
	@cd frontend && npm start

auto-promote-admin: ## Background: auto-promote all local dev users to admin
	@while true; do \
		sqlite3 /tmp/recipes.db "PRAGMA journal_mode=WAL; PRAGMA busy_timeout=5000; UPDATE users SET role = 'admin' WHERE role != 'admin';" 2>/dev/null; \
		sleep 2; \
	done

db-reset: ## Pull latest database backup from GCP
	@echo "Pulling latest database backup from GCP..."
	@rm -f /tmp/recipes.db /tmp/recipes.db-wal /tmp/recipes.db-shm
	@LATEST_BACKUP=$$(gsutil ls -l gs://recipebook2-d0440-recipebook-backups/*.db 2>/dev/null | grep -v TOTAL | sort -k2 | tail -1 | awk '{print $$3}'); \
	if [ -z "$$LATEST_BACKUP" ]; then \
		echo "No backups found, downloading from primary DB bucket..."; \
		gsutil cp gs://recipebook2-d0440-recipebook-db/recipes.db /tmp/recipes.db; \
	else \
		echo "Downloading: $$LATEST_BACKUP"; \
		gsutil cp "$$LATEST_BACKUP" /tmp/recipes.db; \
	fi
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
