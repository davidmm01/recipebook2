# Infrastructure

Terraform configuration for RecipeBook GCP infrastructure.

## Resources

- **Cloud Storage bucket** - Stores SQLite database file with versioning
- **Service Account** - For Cloud Run to access Cloud Storage
- **Cloud Run service** - Backend API deployment with `maxScale=1`
- **IAM permissions** - Public access to Cloud Run service
- **Firebase Hosting** (TODO) - Frontend deployment

## Current Setup

This creates:
1. Cloud Storage bucket: `{project-id}-recipebook-db`
   - Versioning enabled (keeps last 5 versions)
   - Used to store `recipes.db` SQLite file
2. Service account: `recipebook-backend@{project-id}.iam.gserviceaccount.com`
   - Has `roles/storage.objectAdmin` on the bucket
3. Cloud Run service: `recipebook-backend`
   - **Single instance** (`maxScale=1`) to avoid consistency issues
   - Scales to zero when idle (free)
   - 512MB memory, 1 CPU
   - Environment: `DB_BUCKET_NAME` set automatically

## Usage

```bash
# Initialize Terraform
terraform init

# Plan changes
terraform plan

# Apply changes
terraform apply

# View outputs (bucket name and service account)
terraform output
```

## Configuration

Create `terraform.tfvars` with your GCP project ID:

```hcl
project_id = "your-gcp-project-id"
region     = "us-central1"  # Optional, defaults to us-central1
```

## Outputs

After applying, Terraform outputs:
- `database_bucket_name` - Cloud Storage bucket name
- `backend_service_account_email` - Service account email
- `backend_url` - Cloud Run service URL

## Deployment

See [../docs/DEPLOYMENT.md](../docs/DEPLOYMENT.md) for complete deployment instructions.

## Architecture

See [../docs/DATABASE_ARCHITECTURE.md](../docs/DATABASE_ARCHITECTURE.md) for details on the SQLite + Cloud Storage architecture.

See [../docs/CLOUD_RUN_DATABASE_BEHAVIOR.md](../docs/CLOUD_RUN_DATABASE_BEHAVIOR.md) for details on how the database works in Cloud Run.
