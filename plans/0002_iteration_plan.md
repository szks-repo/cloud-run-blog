# Iteration Plan - Backend & Infrastructure Enhancements

## Goals for Next Iteration
- Strengthen backend domain layer with tests and persistence abstraction.
- Extend Terraform for networking/database options while keeping secrets external.
- Prepare deployment guardrails (staging config, manual approval steps).

## Planned Tasks
1. **Backend Hardening**
   - Add repository and handler unit tests (use in-memory repo).
   - Introduce Markdown rendering (e.g., goldmark) for post preview while validating dependency footprint.
   - Add DTO/input validation helpers to reduce duplication across handlers.
2. **Persistence Path**
   - Evaluate Cloud SQL (PostgreSQL) vs. Firestore; prototype Terraform module for chosen backend with minimal IAM permissions.
   - Implement repository adapter interface for managed storage and wire via build tag or env switch.
3. **Infrastructure Enhancements**
   - Add Terraform variables for VPC connector (optional) and custom domains.
   - Configure Terraform outputs for deploy artifacts (service account, repository URL).
4. **CI/CD Improvements**
   - Introduce lint/static analysis (e.g., `go vet`, `staticcheck` if dependency acceptable).
   - Add Terraform plan job using GitHub OIDC (requires secrets) with manual approval gate.
   - Wire Terraform apply (manual or automated) so infrastructure changes are managed via IaC rather than ad-hoc gcloud commands.
5. **Documentation**
   - Expand README with local development guide, Terraform usage, and CI/CD overview.
   - Document secret management strategy (GitHub Actions secrets, Cloud Secret Manager).

## Open Questions
- Which managed storage aligns best with study goals and budget (Cloud SQL vs. Firestore)?
- Is staging environment required, or is a single dev/prod project sufficient for now?
- Should Markdown export include front matter metadata for compatibility with other static site tools?
