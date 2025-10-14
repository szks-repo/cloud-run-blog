variable "project_id" {
  description = "Google Cloud project ID."
  type        = string
}

variable "location" {
  description = "Region where the Cloud Run service is deployed."
  type        = string
}

variable "service_name" {
  description = "Cloud Run service name."
  type        = string
}

variable "image" {
  description = "Container image to deploy."
  type        = string
}

variable "container_port" {
  description = "Port exposed by the container."
  type        = number
  default     = 8080
}

variable "env_vars" {
  description = "Environment variables for the container."
  type        = map(string)
  default     = {}
}

variable "allow_unauthenticated" {
  description = "Allow unauthenticated invocations."
  type        = bool
  default     = true
}

variable "ingress" {
  description = "Ingress traffic configuration."
  type        = string
  default     = "INGRESS_TRAFFIC_ALL"
}

variable "min_instance_count" {
  description = "Minimum number of container instances."
  type        = number
  default     = 0
}

variable "max_instance_count" {
  description = "Maximum number of container instances."
  type        = number
  default     = 3
}

variable "service_account_suffix" {
  description = "Suffix appended to the generated service account (keep <= 6 chars)."
  type        = string
  default     = "svc"
}
