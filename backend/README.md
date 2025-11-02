# Backend API

Go REST API for RecipeBook using SQLite + Cloud Storage.

## Running locally

```bash
# Set required environment variables
export DB_BUCKET_NAME="your-project-recipebook-db"
export GOOGLE_APPLICATION_CREDENTIALS="./service-account.json"

# Run server
go run .
# Or use the Makefile
make run
```

Server will start on port 8080 (or PORT env var).

### Using the Makefile

The project includes a Makefile with common development tasks:

```bash
make help           # Show all available commands
make build          # Build the binary
make run            # Run the application
make test           # Run tests
make test-coverage  # Run tests with coverage report
make lint           # Run code linters (fmt + vet)
make docker-build   # Build Docker image
make docker-run     # Run in Docker container
make dev            # Run with auto-reload (requires air)
make all            # Run lint, test, and build
```

## Environment Variables

- `DB_BUCKET_NAME` - Cloud Storage bucket name for SQLite database (required)
- `GOOGLE_APPLICATION_CREDENTIALS` - Path to service account JSON (for local dev)
- `PORT` - Server port (defaults to 8080)

## API Documentation

See [RECIPE_API.md](RECIPE_API.md) for complete API documentation.

**Quick Overview:**
- `GET /health` - Health check (no auth)
- `GET /recipes` - List all recipes (auth required)
- `POST /recipes` - Create recipe (auth required)
- `GET /recipes/{id}` - Get single recipe (auth required)
- `PUT /recipes/{id}` - Update recipe (auth required)
- `DELETE /recipes/{id}` - Delete recipe (auth required)
- `GET /recipes/search?q=query` - Full-text search (auth required)

## Architecture

See [../docs/DATABASE_ARCHITECTURE.md](../docs/DATABASE_ARCHITECTURE.md) for database details.

## Deployment

Deployed to Google Cloud Run via Terraform (see `../infra/`).
