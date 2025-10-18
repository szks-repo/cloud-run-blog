# cloud-run-blog

Study-oriented blog API deployed on Google Cloud Run. Infrastructure is managed via Terraform and deployed through GitHub Actions.

## CI/CD Workflows

- `Deploy` – builds a container image with Cloud Build and deploys it directly to Cloud Run via `gcloud run deploy`.
- `Terraform Deploy` – builds the container image, generates a Terraform plan, and (optionally) applies infrastructure changes through Terraform.

## Terraform Deploy workflow

Runs are triggered manually (`workflow_dispatch`) to keep full control over infrastructure changes.

### Required GitHub secrets

Configure the following repository secrets before running the workflow:

- `WORKLOAD_IDENTITY_PROVIDER` – OIDC provider resource used for GitHub Actions authentication.
- `GOOGLE_CLOUD_SERVICE_ACCOUNT` – Service account email with permissions for Cloud Build, Artifact Registry, and Terraform-managed resources.
- `GOOGLE_CLOUD_PROJECT_ID` – Target Google Cloud project ID.
- `ARTIFACT_REPOSITORY` – Artifact Registry repository path (e.g. `asia-northeast1-docker.pkg.dev/my-project/cloud-run-blog/app`).
- `TERRAFORM_STATE_BUCKET` – GCS bucket name for Terraform remote state.
- `TERRAFORM_STATE_PREFIX` – Key prefix within the state bucket (e.g. `cloud-run-blog/prod`).

Create a protected GitHub environment (e.g. `production`) with required reviewers to gate the `terraform-apply` job if manual approval is desired.

### Running the workflow

1. Ensure the Terraform state bucket exists and the service account can read/write to it.
2. Trigger **Deploy** from the Actions tab.
3. The `build` job submits a Cloud Build, waits for completion, and publishes the built image digest.
4. The `terraform-plan` job initialises Terraform (`infra/terraform`), formats/validates the configuration, and creates a plan artefact plus a summary in the job log.
5. When the plan reports changes, the `terraform-apply` job becomes available and applies the uploaded plan (subject to environment approval).

## Local Terraform usage

```bash
cd infra/terraform
terraform init \
  -backend-config="bucket=YOUR_STATE_BUCKET" \
  -backend-config="prefix=cloud-run-blog/dev"
terraform plan \
  -var "project_id=YOUR_PROJECT" \
  -var "image=asia-northeast1-docker.pkg.dev/YOUR_PROJECT/cloud-run-blog/app:latest"
```
