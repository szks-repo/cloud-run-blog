output "repository" {
  description = "Full resource name of the Artifact Registry repository."
  value = coalesce(
    try(google_artifact_registry_repository.repository[0].name, null),
    try(data.google_artifact_registry_repository.existing[0].name, null)
  )
}
