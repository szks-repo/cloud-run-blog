locals {
  location = lower(var.location)
}

resource "google_project_service" "artifactregistry" {
  project            = var.project_id
  service            = "artifactregistry.googleapis.com"
  disable_on_destroy = false
}

resource "google_artifact_registry_repository" "repository" {
  project       = var.project_id
  location      = local.location
  repository_id = var.repository_id
  description   = var.description
  format        = "DOCKER"

  depends_on = [google_project_service.artifactregistry]
}
