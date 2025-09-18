### Operations

#### Branch protection (GitHub)
- Protect `main`: require PRs, 1 approval, and checks (pre-commit, super-linter, semantic-pr, branch-name).

#### Secrets
- Store tokens in repo secrets or `.env` (local only). Never commit secrets.

#### CI
- Pre-commit and Super-Linter run on PRs/push to `main`.

#### Docker updates
- Standard deploy: build image, update compose or K8s manifests, rollout with health checks.
- Rollback: revert image tag and redeploy.


