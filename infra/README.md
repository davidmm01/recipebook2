# Infrastructure

Terraform configuration for RecipeBook GCP infrastructure.

## Resources

- **Cloud Storage (Database)** - Private bucket for SQLite database with versioning
- **Cloud Storage (Images)** - Public bucket for recipe images with CORS
- **Service Account** - For Cloud Run to access Storage and Secrets
- **Secret Manager** - Stores Firebase service account credentials
- **Cloud Run** - Backend API deployment with `maxScale=1`
- **Firebase Hosting** - Frontend React app deployment
- **Cloud Monitoring** - Alert policies for errors, latency, and scaling

## Quick Start

### 1. Configure Variables

Copy the example tfvars file and fill in your values:

```bash
cp terraform.tfvars.example terraform.tfvars
```

Edit `terraform.tfvars`:

```hcl
project_id  = "your-gcp-project-id"
region      = "us-central1"
environment = "prod"
alert_email = "your-email@example.com"

# Optional: Restrict CORS to your domains
# cors_origins = ["https://your-project-web.web.app", "http://localhost:3000"]
```

### 2. Initialize and Apply Terraform

```bash
terraform init
terraform plan
terraform apply
```

### 3. Upload Firebase Service Account to Secret Manager

After Terraform creates the secret, upload your service account JSON:

```bash
# Get the secret ID from Terraform output
terraform output secret_manager_secret_id

# Upload the secret value
gcloud secrets versions add firebase-service-account \
  --data-file=path/to/your/service-account.json
```

### 4. Deploy Backend to Cloud Run

```bash
cd ../backend
gcloud run deploy recipebook-backend \
  --source . \
  --region us-central1 \
  --allow-unauthenticated
```

### 5. Deploy Frontend to Firebase Hosting

```bash
cd ../frontend
npm run build
firebase deploy --only hosting
```

## Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `project_id` | Yes | - | GCP Project ID |
| `region` | No | `us-central1` | GCP region for resources |
| `environment` | No | `prod` | Environment name (dev/staging/prod) |
| `alert_email` | Yes | - | Email for monitoring alerts |
| `cors_origins` | No | `["*"]` | Allowed CORS origins for images bucket |

## Outputs

| Output | Description |
|--------|-------------|
| `database_bucket_name` | Private bucket for SQLite database |
| `images_bucket_name` | Public bucket for recipe images |
| `backend_service_account_email` | Service account email for Cloud Run |
| `backend_url` | Cloud Run service URL |
| `firebase_hosting_site` | Firebase Hosting URL |
| `secret_manager_secret_id` | Secret ID for Firebase credentials |

## Architecture

### Storage Buckets

| Bucket | Access | Purpose |
|--------|--------|---------|
| `{project}-recipebook-db` | Private (enforced) | SQLite database, versioned (5 versions) |
| `{project}-recipebook-images` | Public read | Recipe images and icons |

### Cloud Run Configuration

- **Single instance** (`maxScale=1`) to avoid SQLite consistency issues
- **Scales to zero** when idle (cost-effective)
- **512MB memory, 1 CPU**
- Environment variables:
  - `DB_BUCKET_NAME` - Database bucket
  - `IMAGES_BUCKET_NAME` - Images bucket
  - `ENVIRONMENT` - Current environment
  - `FIREBASE_SERVICE_ACCOUNT` - Injected from Secret Manager

### Monitoring Alerts

| Alert | Condition | Duration |
|-------|-----------|----------|
| High Error Rate | 5xx errors > 5% | 5 minutes |
| High Latency | p99 > 5 seconds | 5 minutes |
| Max Capacity | Instance count = 1 | 10 minutes |

## Useful Commands

```bash
# View all outputs
terraform output

# View specific output
terraform output backend_url

# Plan changes without applying
terraform plan

# Destroy all resources (careful!)
terraform destroy

# Format Terraform files
terraform fmt

# Validate configuration
terraform validate
```

## Related Documentation

- [Deployment Guide](../docs/DEPLOYMENT.md)
- [Database Architecture](../docs/DATABASE_ARCHITECTURE.md)
- [Cloud Run Database Behavior](../docs/CLOUD_RUN_DATABASE_BEHAVIOR.md)
- [API Documentation](../docs/API.md)
