terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

# Cloud Storage bucket for SQLite database
resource "google_storage_bucket" "database" {
  name          = "${var.project_id}-recipebook-db"
  location      = var.region
  force_destroy = false

  uniform_bucket_level_access = true

  versioning {
    enabled = true
  }

  lifecycle_rule {
    condition {
      num_newer_versions = 5
    }
    action {
      type = "Delete"
    }
  }
}

# Service account for Cloud Run
resource "google_service_account" "backend" {
  account_id   = "recipebook-backend"
  display_name = "RecipeBook Backend Service Account"
}

# Grant Cloud Run service account access to the storage bucket
resource "google_storage_bucket_iam_member" "backend_storage_access" {
  bucket = google_storage_bucket.database.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.backend.email}"
}

# Cloud Run service for backend
resource "google_cloud_run_v2_service" "backend" {
  name     = "recipebook-backend"
  location = var.region

  template {
    service_account = google_service_account.backend.email

    # Force single instance to avoid eventual consistency issues
    scaling {
      min_instance_count = 0  # Scale to zero when idle (free)
      max_instance_count = 1  # Only one container max
    }

    containers {
      # Image will be set during deployment via gcloud
      # For initial terraform apply, use a placeholder
      image = "us-docker.pkg.dev/cloudrun/container/hello"

      env {
        name  = "DB_BUCKET_NAME"
        value = google_storage_bucket.database.name
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }
    }
  }

  # Allow unauthenticated requests (auth handled by Firebase tokens in app)
  # You may want to restrict this later
}

# Allow public access to Cloud Run service
resource "google_cloud_run_v2_service_iam_member" "public_access" {
  name     = google_cloud_run_v2_service.backend.name
  location = google_cloud_run_v2_service.backend.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# Output the bucket name for use in Cloud Run environment variable
output "database_bucket_name" {
  value       = google_storage_bucket.database.name
  description = "Name of the Cloud Storage bucket storing the SQLite database"
}

output "backend_service_account_email" {
  value       = google_service_account.backend.email
  description = "Email of the service account for Cloud Run"
}

output "backend_url" {
  value       = google_cloud_run_v2_service.backend.uri
  description = "URL of the Cloud Run backend service"
}

# TODO: Add Firebase Hosting configuration
