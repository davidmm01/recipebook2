variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "region" {
  description = "GCP region for resources"
  type        = string
  default     = "us-central1"
}

variable "environment" {
  description = "Environment name (dev, staging, prod)"
  type        = string
  default     = "prod"

  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "Environment must be one of: dev, staging, prod"
  }
}

variable "alert_email" {
  description = "Email address to receive monitoring alerts"
  type        = string
}

variable "cors_origins" {
  description = "Allowed CORS origins for the images bucket"
  type        = list(string)
  default     = ["*"]
}
