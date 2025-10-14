# Terraform Infrastructure

This directory contains Infrastructure as Code definitions for deploying the Cloud Run blog study project.

## Layout

- `main.tf` – Root module wiring Artifact Registry and Cloud Run child modules.
- `modules/artifact_registry` – Creates a regional Docker repository.
- `modules/cloud_run_service` – Provisions the Cloud Run service, service account, and IAM bindings.
- `variables.tf` – Input variables (project, region, image, etc.).
- `outputs.tf` – Convenience outputs (service URL, service account).

## Usage

1. Ensure the Google Cloud project exists and billing/APIs are enabled.
2. Create a state bucket (or choose an alternative backend) and uncomment the `backend "gcs"` block in `main.tf`.
3. Provide required variables, for example via `terraform.tfvars` or environment variables:

```hcl
project_id     = "my-gcp-project"
region         = "us-central1"
artifact_region = "us-central1"
image          = "us-central1-docker.pkg.dev/my-gcp-project/cloud-run-blog/app:latest"
env_vars = {
  "APP_ENV" = "production"
}
```

4. Run the usual Terraform workflow:

```bash
terraform init
terraform plan
terraform apply
```

## Security Notes

- No secrets are stored in this repository. Supply sensitive values (e.g., database credentials) through Terraform variables or Cloud Run secrets at deploy time.
- The Cloud Run module can be configured to disable public access by setting `allow_unauthenticated = false`.
- IAM bindings are scoped to the provided project; adjust if you use a dedicated build/deploy project.
