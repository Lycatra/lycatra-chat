# Deployment and Release Strategy

Lycatra-chat prioritizes quick, reliable updates without sacrificing stability. This document expands
on the goals briefly outlined in the README and describes the approach for reaching near zero-downtime
operations as the project matures.

## Objectives

1. **Consistent rollouts.** Releases should be repeatable and automated so that every environment
   receives the same build artifacts and configuration.
2. **Fast recovery.** If a regression slips through, we must be able to restore the previous version
   in minutes.
3. **Local-first experience.** Early iterations should remain easy to run on a developer workstation
   while paving a path to production infrastructure.

## Phased roadmap

| Phase | Environment focus        | Deployment technique | Key outcomes |
| ----- | ------------------------ | -------------------- | ------------ |
| 1     | Local development        | Manual restart       | Rapid feedback while iterating on features. |
| 2     | Home lab / single host   | Rolling restart      | Reduce downtime to seconds during upgrades. |
| 3     | Talos Kubernetes cluster | Blue/green services  | Deterministic rollouts with instant rollback. |

### Phase 1 – Local development

* Run the API server and supporting services directly on a workstation or inside lightweight
  containers.
* Use `make run` (or the appropriate language-specific command) to restart the service after
  code changes. Aim for sub-5-second startup time to encourage frequent restarts.
* Capture manual steps in `docs/OPERATIONS.md` so they can be scripted later.

### Phase 2 – Home lab rollouts

* Package the service as a container image published to a local registry (for example, using
  `docker buildx bake release`).
* Maintain two running instances behind a simple reverse proxy (Caddy, Nginx, or Traefik). Restart
  them sequentially to achieve a rolling update.
* Store configuration in environment files committed to the infrastructure repository, and load them
  through Docker Compose or systemd units.

### Phase 3 – Talos Kubernetes deployment

* Use Talos as the operating system for the cluster and manage workloads with Kubernetes manifests
  stored in Git.
* Deploy via progressive delivery: first update a staging namespace, run smoke tests, then promote
  to production through GitOps tooling (Flux or Argo CD).
* Implement a blue/green strategy by running two Deployments (current and candidate) and switch
  traffic via a Service update when the candidate is healthy.
* Automate rollback by keeping previous manifests tagged and accessible through the GitOps system.

## Operational safeguards

* **Health checks:** Implement liveness and readiness probes that align with API behavior so rollouts
  only progress when the service is ready to accept traffic.
* **Observability:** Collect metrics (request rate, latency, error percentage) and log streams during
  and after a rollout to confirm stability.
* **Release checklists:** Before promotion, verify schema migrations, configuration changes, and
  third-party integrations.

## Next steps

* Prototype the container build pipeline and document the commands in `docs/OPERATIONS.md`.
* Define the smoke-test suite that validates critical API endpoints before switching traffic.
* Draft a disaster-recovery playbook that references this deployment plan and identifies responsible
  owners.
