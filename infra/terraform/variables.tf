variable "project_id" {
  description = "Google Cloud project ID where resources will be created."
  type        = string
}

variable "region" {
  description = "Region for Cloud Run service deployment."
  type        = string
  default     = "asia-northeast1"
}

variable "artifact_region" {
  description = "Region for Artifact Registry repository (defaults to the Cloud Run region)."
  type        = string
  default     = "asia-northeast1"
}

variable "repository_id" {
  description = "Artifact Registry repository identifier."
  type        = string
  default     = "cloud-run-blog"
}

variable "manage_artifact_registry" {
  description = "Whether Terraform should create/manage the Artifact Registry repository. Set to false to use an existing repository."
  type        = bool
  default     = true
}

variable "service_name" {
  description = "Name of the Cloud Run service."
  type        = string
  default     = "cloud-run-blog"
}

variable "image" {
  description = "Container image reference (e.g. asia-northeast1-docker.pkg.dev/PROJECT/REPO/app:tag)."
  type        = string
}

variable "container_port" {
  description = "Container port exposed by the Cloud Run service."
  type        = number
  default     = 8080
}

variable "allow_unauthenticated" {
  description = "Whether to allow unauthenticated invocations of the Cloud Run service."
  type        = bool
  default     = true
}

variable "ingress" {
  description = "Ingress setting for the Cloud Run service."
  type        = string
  default     = "INGRESS_TRAFFIC_ALL"
}

variable "min_instance_count" {
  description = "Minimum number of instances for the Cloud Run service."
  type        = number
  default     = 0
}

variable "max_instance_count" {
  description = "Maximum number of instances for the Cloud Run service."
  type        = number
  default     = 3
}

variable "env_vars" {
  description = "Environment variables passed to the Cloud Run container."
  type        = map(string)
  default     = {}
}

variable "service_account_suffix" {
  description = "Optional suffix appended to the generated Cloud Run service account ID."
  type        = string
  default     = "svc"
}
