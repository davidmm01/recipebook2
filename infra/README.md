# Infrastructure

Terraform configuration for RecipeBook GCP infrastructure.

## Resources

- Firestore database
- Cloud Run service (backend API)
- Firebase Hosting (frontend)

## Usage

```bash
# Initialize Terraform
terraform init

# Plan changes
terraform plan

# Apply changes
terraform apply
```

## Configuration

Set your GCP project ID in `terraform.tfvars`:

```hcl
project_id = "your-gcp-project-id"
region     = "us-central1"
```
