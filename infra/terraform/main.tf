terraform {
  required_version = ">= 1.8, < 2.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.41"
    }
  }

  backend "gcs" {}
}

provider "google" {
  project = var.project_id
  region  = var.region
}

provider "google" {
  alias   = "artifact_registry"
  project = var.project_id
  region  = var.artifact_region
}

module "artifact_registry" {
  source = "./modules/artifact_registry"

  project_id    = var.project_id
  location      = var.artifact_region
  repository_id = var.repository_id
  description   = "Container images for the Cloud Run Blog study project"
  manage_repository = var.manage_artifact_registry
}

module "cloud_run_service" {
  source = "./modules/cloud_run_service"

  project_id             = var.project_id
  location               = var.region
  service_name           = var.service_name
  image                  = var.image
  container_port         = var.container_port
  allow_unauthenticated  = var.allow_unauthenticated
  ingress                = var.ingress
  min_instance_count     = var.min_instance_count
  max_instance_count     = var.max_instance_count
  env_vars               = var.env_vars
  service_account_suffix = var.service_account_suffix
  depends_on             = [module.artifact_registry]
}

output "cloud_run_url" {
  description = "Default URL for the deployed Cloud Run service."
  value       = module.cloud_run_service.url
}

output "service_account_email" {
  description = "Service account email used by the Cloud Run service."
  value       = module.cloud_run_service.service_account_email
}

output "artifact_registry_repository" {
  description = "Artifact Registry repository resource name."
  value       = module.artifact_registry.repository
}
