output "repository" {
  description = "Full resource name of the Artifact Registry repository."
  value       = google_artifact_registry_repository.repository.name
}
