### RFC-001: Release → Matrix approval → Docker update

#### Goal
When a new GitHub Release is published for selected services, notify a Matrix room with the release link. On 👍 reaction by an authorized user, perform a safe rolling update with fallback.

#### Flow
1. Detect release (GitHub webhook or scheduled poll).
2. Post message in Matrix room with release title, notes link, and proposed actions.
3. Wait for 👍 reaction from an allowed user within a time window.
4. On approval: run update plan:
   - Pull new image/tag
   - Health-check on canary
   - Swap traffic (or restart compose service)
   - Monitor; on failure, rollback to previous tag
5. Report status back to Matrix.

#### Scope (MVP)
- GitHub: poll releases of configured repos every N minutes.
- Matrix: single room, one approver list (user IDs).
- Docker: Docker Compose local host; later add K8s rollout.

#### Config
```toml
[matrix]
homeserver = "https://matrix.example.com"
access_token = "$MATRIX_TOKEN"
room_id = "!room:example.com"
approvers = ["@you:example.com"]

[github]
repos = ["org/service-a", "org/service-b"]
poll_interval_seconds = 300

[deploy]
compose_file = "docker-compose.yml"
services = { "service-a" = "org/service-a", "service-b" = "org/service-b" }
healthcheck_url = "http://localhost:8080/healthz"
rollback_on_failure = true
```

#### Security
- Validate reaction user is in approvers.
- Ignore duplicate approvals; idempotent deploy.
- Keep last known-good image tag for rollback.

#### Observability
- Structured logs; summary back to Matrix.


