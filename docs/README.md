# RecipeBook Documentation

Complete documentation for the RecipeBook application.

## Quick Links

- **[API.md](API.md)** - Complete API reference and endpoints
- **[DEPLOYMENT.md](DEPLOYMENT.md)** - Step-by-step deployment guide
- **[IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md)** - What's built and next steps
- **[DATABASE_ARCHITECTURE.md](DATABASE_ARCHITECTURE.md)** - SQLite + Cloud Storage design
- **[CLOUD_RUN_DATABASE_BEHAVIOR.md](CLOUD_RUN_DATABASE_BEHAVIOR.md)** - How database works in containers
- **[PUBLIC_READ_MODEL.md](PUBLIC_READ_MODEL.md)** - Public read, authenticated write model
- **[SECURITY_OPTIONS.md](SECURITY_OPTIONS.md)** - Security considerations and alternatives

## Getting Started

1. Read [IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md) for an overview
2. Follow [DEPLOYMENT.md](DEPLOYMENT.md) to deploy to GCP
3. Understand [PUBLIC_READ_MODEL.md](PUBLIC_READ_MODEL.md) for the auth model

## Architecture Documents

### Database
- [DATABASE_ARCHITECTURE.md](DATABASE_ARCHITECTURE.md) - Why SQLite + Cloud Storage, schema, costs
- [CLOUD_RUN_DATABASE_BEHAVIOR.md](CLOUD_RUN_DATABASE_BEHAVIOR.md) - Container lifecycle, consistency model

### Security
- [PUBLIC_READ_MODEL.md](PUBLIC_READ_MODEL.md) - Public reads, authenticated writes
- [SECURITY_OPTIONS.md](SECURITY_OPTIONS.md) - IAM, Cloud Armor, Firebase App Check options

### Operations
- [DEPLOYMENT.md](DEPLOYMENT.md) - Deployment steps, monitoring, troubleshooting
- [IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md) - Current status and roadmap

## API Documentation

See [API.md](API.md) for complete API reference including all endpoints, request/response formats, and authentication details.

## Infrastructure

See [../infra/README.md](../infra/README.md) for Terraform setup.
