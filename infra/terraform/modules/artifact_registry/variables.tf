variable "project_id" {
  description = "Google Cloud project ID."
  type        = string
}

variable "location" {
  description = "Repository location (region)."
  type        = string
}

variable "repository_id" {
  description = "Artifact Registry repository identifier."
  type        = string
}

variable "description" {
  description = "Optional description for the repository."
  type        = string
  default     = ""
}
