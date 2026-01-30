terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 5.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

provider "google-beta" {
  project = var.project_id
  region  = var.region
}

# =============================================================================
# ENABLE REQUIRED APIS
# =============================================================================

resource "google_project_service" "required_apis" {
  for_each = toset([
    "run.googleapis.com",
    "storage.googleapis.com",
    "secretmanager.googleapis.com",
    "monitoring.googleapis.com",
    "firebase.googleapis.com",
    "firebasehosting.googleapis.com",
    "artifactregistry.googleapis.com",
  ])

  service            = each.value
  disable_on_destroy = false
}

# =============================================================================
# SERVICE ACCOUNT
# =============================================================================

resource "google_service_account" "backend" {
  account_id   = "recipebook-backend"
  display_name = "RecipeBook Backend Service Account"

  depends_on = [google_project_service.required_apis]
}

# =============================================================================
# CLOUD STORAGE - DATABASE BUCKET (PRIVATE)
# =============================================================================

resource "google_storage_bucket" "database" {
  name          = "${var.project_id}-recipebook-db"
  location      = var.region
  force_destroy = false

  # Keep bucket private - only service account can access
  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"

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

  labels = {
    app         = "recipebook"
    environment = var.environment
    purpose     = "database"
  }

  depends_on = [google_project_service.required_apis]
}

# Grant service account access to database bucket
resource "google_storage_bucket_iam_member" "backend_database_access" {
  bucket = google_storage_bucket.database.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.backend.email}"
}

# =============================================================================
# CLOUD STORAGE - BACKUPS BUCKET (PRIVATE)
# =============================================================================

resource "google_storage_bucket" "backups" {
  name          = "${var.project_id}-recipebook-backups"
  location      = var.region
  force_destroy = false

  # Keep bucket private
  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"

  # Keep backups for 90 days, then delete
  lifecycle_rule {
    condition {
      age = 90
    }
    action {
      type = "Delete"
    }
  }

  labels = {
    app         = "recipebook"
    environment = var.environment
    purpose     = "backups"
  }

  depends_on = [google_project_service.required_apis]
}

# Grant service account access to backups bucket
resource "google_storage_bucket_iam_member" "backend_backups_access" {
  bucket = google_storage_bucket.backups.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.backend.email}"
}

# =============================================================================
# CLOUD STORAGE - IMAGES BUCKET (PUBLIC READ)
# =============================================================================

resource "google_storage_bucket" "images" {
  name          = "${var.project_id}-recipebook-images"
  location      = var.region
  force_destroy = false

  uniform_bucket_level_access = true

  # CORS configuration for frontend direct uploads/access
  cors {
    origin          = var.cors_origins
    method          = ["GET", "HEAD", "OPTIONS"]
    response_header = ["Content-Type", "Content-Length", "Content-Range"]
    max_age_seconds = 3600
  }

  # Auto-delete incomplete uploads after 1 day
  lifecycle_rule {
    condition {
      age = 1
    }
    action {
      type = "AbortIncompleteMultipartUpload"
    }
  }

  labels = {
    app         = "recipebook"
    environment = var.environment
    purpose     = "images"
  }

  depends_on = [google_project_service.required_apis]
}

# Grant service account write access to images bucket
resource "google_storage_bucket_iam_member" "backend_images_access" {
  bucket = google_storage_bucket.images.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.backend.email}"
}

# Grant public read access to images bucket
resource "google_storage_bucket_iam_member" "images_public_read" {
  bucket = google_storage_bucket.images.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}

# =============================================================================
# SECRET MANAGER
# =============================================================================

# Secret for Firebase service account JSON
resource "google_secret_manager_secret" "firebase_service_account" {
  secret_id = "firebase-service-account"

  replication {
    auto {}
  }

  labels = {
    app         = "recipebook"
    environment = var.environment
  }

  depends_on = [google_project_service.required_apis]
}

# Grant service account access to read the secret
resource "google_secret_manager_secret_iam_member" "backend_secret_access" {
  secret_id = google_secret_manager_secret.firebase_service_account.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.backend.email}"
}

# =============================================================================
# CLOUD RUN SERVICE
# =============================================================================

resource "google_cloud_run_v2_service" "backend" {
  name     = "recipebook-backend"
  location = var.region

  template {
    service_account = google_service_account.backend.email

    # Force single instance to avoid SQLite consistency issues
    scaling {
      min_instance_count = 0
      max_instance_count = 1
    }

    containers {
      # Image will be set during deployment via gcloud
      image = "us-docker.pkg.dev/cloudrun/container/hello"

      env {
        name  = "DB_BUCKET_NAME"
        value = google_storage_bucket.database.name
      }

      env {
        name  = "IMAGES_BUCKET_NAME"
        value = google_storage_bucket.images.name
      }

      env {
        name  = "ENVIRONMENT"
        value = var.environment
      }

      # Reference secret from Secret Manager
      env {
        name = "FIREBASE_SERVICE_ACCOUNT"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.firebase_service_account.secret_id
            version = "latest"
          }
        }
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }
    }
  }

  depends_on = [
    google_project_service.required_apis,
    google_secret_manager_secret.firebase_service_account,
  ]
}

# Allow public access to Cloud Run service (auth handled by Firebase tokens)
resource "google_cloud_run_v2_service_iam_member" "public_access" {
  name     = google_cloud_run_v2_service.backend.name
  location = google_cloud_run_v2_service.backend.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# =============================================================================
# FIREBASE HOSTING
# =============================================================================

resource "google_firebase_project" "default" {
  provider = google-beta
  project  = var.project_id

  depends_on = [google_project_service.required_apis]
}

resource "google_firebase_hosting_site" "default" {
  provider = google-beta
  project  = var.project_id
  site_id  = "${var.project_id}-web"

  depends_on = [google_firebase_project.default]
}

# =============================================================================
# CLOUD MONITORING
# =============================================================================

# Email notification channel for alerts
resource "google_monitoring_notification_channel" "email" {
  display_name = "RecipeBook Alert Email"
  type         = "email"

  labels = {
    email_address = var.alert_email
  }

  depends_on = [google_project_service.required_apis]
}

# Alert: Cloud Run high error rate
resource "google_monitoring_alert_policy" "cloud_run_errors" {
  display_name = "RecipeBook - High Error Rate"
  combiner     = "OR"

  conditions {
    display_name = "Cloud Run 5xx error rate > 5%"

    condition_threshold {
      filter          = "resource.type = \"cloud_run_revision\" AND resource.labels.service_name = \"recipebook-backend\" AND metric.type = \"run.googleapis.com/request_count\" AND metric.labels.response_code_class = \"5xx\""
      duration        = "300s"
      comparison      = "COMPARISON_GT"
      threshold_value = 0.05

      aggregations {
        alignment_period     = "60s"
        per_series_aligner   = "ALIGN_RATE"
        cross_series_reducer = "REDUCE_SUM"
      }
    }
  }

  notification_channels = [google_monitoring_notification_channel.email.id]

  alert_strategy {
    auto_close = "1800s"
  }

  depends_on = [google_cloud_run_v2_service.backend]
}

# Alert: Cloud Run high latency
resource "google_monitoring_alert_policy" "cloud_run_latency" {
  display_name = "RecipeBook - High Latency"
  combiner     = "OR"

  conditions {
    display_name = "Cloud Run p99 latency > 5s"

    condition_threshold {
      filter          = "resource.type = \"cloud_run_revision\" AND resource.labels.service_name = \"recipebook-backend\" AND metric.type = \"run.googleapis.com/request_latencies\""
      duration        = "300s"
      comparison      = "COMPARISON_GT"
      threshold_value = 5000

      aggregations {
        alignment_period     = "60s"
        per_series_aligner   = "ALIGN_PERCENTILE_99"
        cross_series_reducer = "REDUCE_SUM"
      }
    }
  }

  notification_channels = [google_monitoring_notification_channel.email.id]

  alert_strategy {
    auto_close = "1800s"
  }

  depends_on = [google_cloud_run_v2_service.backend]
}

# Alert: Cloud Run instance count at maximum
resource "google_monitoring_alert_policy" "cloud_run_scaling" {
  display_name = "RecipeBook - Instance at Max Capacity"
  combiner     = "OR"

  conditions {
    display_name = "Cloud Run instance count = 1 (max)"

    condition_threshold {
      filter          = "resource.type = \"cloud_run_revision\" AND resource.labels.service_name = \"recipebook-backend\" AND metric.type = \"run.googleapis.com/container/instance_count\""
      duration        = "600s"
      comparison      = "COMPARISON_GT"
      threshold_value = 0.9

      aggregations {
        alignment_period     = "60s"
        per_series_aligner   = "ALIGN_MAX"
        cross_series_reducer = "REDUCE_SUM"
      }
    }
  }

  notification_channels = [google_monitoring_notification_channel.email.id]

  alert_strategy {
    auto_close = "1800s"
  }

  depends_on = [google_cloud_run_v2_service.backend]
}

# =============================================================================
# ARTIFACT REGISTRY REPOSITORY
# =============================================================================

# Docker repository for Cloud Run source deployments
resource "google_artifact_registry_repository" "cloud_run_source" {
  repository_id = "cloud-run-source-deploy"
  location      = var.region
  description   = "Docker repository for Cloud Run deployments from source"
  format        = "DOCKER"

  labels = {
    app         = "recipebook"
    environment = var.environment
  }

  depends_on = [google_project_service.required_apis]
}

# =============================================================================
# GITHUB ACTIONS SERVICE ACCOUNT
# =============================================================================

# Service account for GitHub Actions to deploy to Cloud Run
resource "google_service_account" "github_actions" {
  account_id   = "github-actions"
  display_name = "GitHub Actions Deployment"

  depends_on = [google_project_service.required_apis]
}

# Grant Cloud Run admin permission
resource "google_project_iam_member" "github_actions_run_admin" {
  project = var.project_id
  role    = "roles/run.admin"
  member  = "serviceAccount:${google_service_account.github_actions.email}"
}

# Grant permission to act as service accounts (needed for Cloud Run deployment)
resource "google_project_iam_member" "github_actions_sa_user" {
  project = var.project_id
  role    = "roles/iam.serviceAccountUser"
  member  = "serviceAccount:${google_service_account.github_actions.email}"
}

# Grant Storage Admin for build artifacts
resource "google_project_iam_member" "github_actions_storage_admin" {
  project = var.project_id
  role    = "roles/storage.admin"
  member  = "serviceAccount:${google_service_account.github_actions.email}"
}

# Grant Artifact Registry Writer for container images
resource "google_project_iam_member" "github_actions_artifact_registry" {
  project = var.project_id
  role    = "roles/artifactregistry.writer"
  member  = "serviceAccount:${google_service_account.github_actions.email}"
}

# Grant Firebase Hosting Admin for frontend deployments
resource "google_project_iam_member" "github_actions_firebase_hosting" {
  project = var.project_id
  role    = "roles/firebasehosting.admin"
  member  = "serviceAccount:${google_service_account.github_actions.email}"
}

# =============================================================================
# OUTPUTS
# =============================================================================

output "database_bucket_name" {
  value       = google_storage_bucket.database.name
  description = "Name of the Cloud Storage bucket storing the SQLite database"
}

output "images_bucket_name" {
  value       = google_storage_bucket.images.name
  description = "Name of the Cloud Storage bucket storing recipe images"
}

output "backups_bucket_name" {
  value       = google_storage_bucket.backups.name
  description = "Name of the Cloud Storage bucket storing database backups"
}

output "backend_service_account_email" {
  value       = google_service_account.backend.email
  description = "Email of the service account for Cloud Run"
}

output "backend_url" {
  value       = google_cloud_run_v2_service.backend.uri
  description = "URL of the Cloud Run backend service"
}

output "firebase_hosting_site" {
  value       = google_firebase_hosting_site.default.default_url
  description = "Default URL for Firebase Hosting"
}

output "secret_manager_secret_id" {
  value       = google_secret_manager_secret.firebase_service_account.secret_id
  description = "Secret Manager secret ID for Firebase service account"
}

output "github_actions_service_account_email" {
  value       = google_service_account.github_actions.email
  description = "Email of the GitHub Actions service account"
}
