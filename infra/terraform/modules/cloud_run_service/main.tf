locals {
  normalized_service_name = trim(regexreplace(lower(var.service_name), "[^a-z0-9-]", "-"), "-")
  service_name            = local.normalized_service_name != "" ? local.normalized_service_name : "app"
  base_service_account_id = "${local.service_name}-${var.service_account_suffix}"
  service_account_id      = substr(local.base_service_account_id, 0, min(30, length(local.base_service_account_id)))
}

resource "google_project_service" "run" {
  project            = var.project_id
  service            = "run.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "iam" {
  project            = var.project_id
  service            = "iam.googleapis.com"
  disable_on_destroy = false
}

resource "google_service_account" "service" {
  project      = var.project_id
  account_id   = local.service_account_id
  display_name = "${var.service_name} Cloud Run service account"

  depends_on = [google_project_service.iam]
}

resource "google_project_iam_member" "artifact_registry_reader" {
  project = var.project_id
  role    = "roles/artifactregistry.reader"
  member  = "serviceAccount:${google_service_account.service.email}"
}

resource "google_cloud_run_v2_service" "service" {
  name     = local.service_name
  location = var.location
  ingress  = var.ingress

  template {
    service_account = google_service_account.service.email

    scaling {
      min_instance_count = var.min_instance_count
      max_instance_count = var.max_instance_count
    }

    containers {
      image = var.image

      ports {
        container_port = var.container_port
      }

      dynamic "env" {
        for_each = var.env_vars
        content {
          name  = env.key
          value = env.value
        }
      }
    }
  }

  traffic {
    percent = 100
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
  }

  depends_on = [
    google_project_service.run,
    google_project_service.iam,
    google_project_iam_member.artifact_registry_reader,
  ]
}

resource "google_cloud_run_v2_service_iam_member" "public_invoker" {
  count = var.allow_unauthenticated ? 1 : 0

  project  = var.project_id
  location = var.location
  name     = google_cloud_run_v2_service.service.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
