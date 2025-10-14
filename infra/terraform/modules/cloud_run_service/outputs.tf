output "url" {
  description = "Cloud Run service URL."
  value       = google_cloud_run_v2_service.service.uri
}

output "service_account_email" {
  description = "Service account email for the Cloud Run service."
  value       = google_service_account.service.email
}
