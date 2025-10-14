# Cloud Run Blog - Initial Plan

## Objectives
- Build a study-oriented blog platform deployed on Google Cloud Run.
- Keep the stack simple and transparent to minimize supply chain risks.
- Manage infrastructure as code with Terraform and publish the project on GitHub without exposing secrets.

## High-Level Architecture
- **Backend:** Go 1.25+ HTTP API exposing endpoints for post CRUD and Markdown export.
- **Storage:** Cloud SQL (PostgreSQL) or Firestore (evaluate trade-offs); start with in-memory or file-based store for local dev, add managed option for Cloud Run.
- **Frontend:** Minimal server-rendered HTML templates bundled with Go backend; no heavy frontend frameworks.
- **CI/CD:** GitHub Actions workflow building/tests, Terraform fmt/validate, build and deploy container to Cloud Run using deploy user-provided secrets (via GitHub Actions Encrypted Secrets).
- **IaC:** Terraform modules for Cloud Run service, container registry (Artifact Registry), state bucket, optional database, and IAM bindings.

## Workstreams & Tasks
1. **Project Foundation**
   - Initialize Go module and folder structure (`cmd/api`, `internal/...`).
   - Configure basic HTTP server with health endpoint.
   - Create minimal HTML template rendering list of posts.
2. **Content Management**
   - Define post model and repository interface.
   - Implement in-memory repository for local dev; design interface for future managed DB.
   - Add CRUD handlers (list, create, update, delete) with simple form submissions.
   - Implement Markdown export endpoint producing `.md` payload.
3. **Infrastructure**
   - Set up Terraform project (`infra/terraform`) with modules for:
     - Artifact Registry
     - Cloud Run service
     - (Optional) Cloud SQL/Firestore
     - Service account & IAM roles
   - Configure remote state (Google Cloud Storage bucket) with placeholders.
4. **CI/CD**
   - GitHub Actions workflow for Go tests, static analysis (`go vet`), and build.
   - Workflow for Terraform fmt/validate and plan (requires OIDC or workload identity later).
   - Deployment workflow building container image and deploying to Cloud Run (manual approval or on release).
5. **Documentation**
   - Update README with architecture overview, local dev instructions, deployment steps.
   - Document secret management approach (env vars, GitHub secrets, Terraform variables).

## Next Steps
- Decide on managed storage (Cloud SQL vs Firestore) before Terraform implementation.
- Draft detailed task issues or sub-plans per workstream.
- Create roadmap for incrementally shipping MVP features.

