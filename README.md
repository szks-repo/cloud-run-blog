# cloud-run-blog

Study-oriented blog API deployed on Google Cloud Run. Infrastructure is managed via Terraform and deployed through GitHub Actions.

## CI/CD Workflows

- `Deploy` – builds a container image with Cloud Build and deploys it directly to Cloud Run via `gcloud run deploy`.
- `Terraform Deploy` – builds the container image, generates a Terraform plan, and (optionally) applies infrastructure changes through Terraform.

## Terraform Deploy workflow

Runs are triggered manually (`workflow_dispatch`) to keep full control over infrastructure changes.

### Required GitHub secrets

Before running the workflow, ensure that the required GitHub Secrets referenced inside `.github/workflows/terraform-deploy.yml` and `.github/workflows/deploy.yml` are populated (OIDC provider, service account, project ID, Artifact Registry reference, and Terraform state configuration). Avoid storing raw credential values in the repository—use Secrets only.

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
  -var "image=asia-northeast1-docker.pkg.dev/YOUR_PROJECT/YOUR_REPOSITORY/cloud-run-blog:latest" \
  -var "repository_id=YOUR_REPOSITORY" \
  -var "manage_artifact_registry=false" # Set to false if the repository already exists
```
