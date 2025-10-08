# Backend API

Go REST API for RecipeBook.

## Running locally

```bash
go run main.go
```

Server will start on port 8080 (or PORT env var).

## Endpoints

- `GET /health` - Health check
- `GET /recipes` - List all recipes (TODO: implement Firestore)
- `POST /recipes` - Create a new recipe (TODO: implement Firestore)

## Deployment

Deployed to Google Cloud Run via Terraform (see `../infra/`).
