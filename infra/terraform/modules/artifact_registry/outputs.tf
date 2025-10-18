output "repository" {
  description = "Full resource name of the Artifact Registry repository."
  value = var.manage_repository ?
    google_artifact_registry_repository.repository[0].name :
    data.google_artifact_registry_repository.existing[0].name
}
