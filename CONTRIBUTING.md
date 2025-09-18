### Contributing

Thank you for helping improve this project! This guide explains how we work together.

### Workflow overview
- **Branching**: Create a feature branch from `main` using the pattern `type/scope-short-description`, e.g. `feat/ui-chat-panel`.
- **Commits**: Use [Conventional Commits](https://www.conventionalcommits.org/) (e.g. `feat(chat): add message send button`).
- **Pull Requests**: Open a PR to `main`. Ensure CI is green and request review.
- **Reviews**: At least one approval before merge. Prefer small, focused PRs.
- **Merges**: Use squash merges with a semantic title.

### Getting started
1. Install Git and a recent Node or Python (as needed by the subproject).
2. Install pre-commit hooks:
   - Python: `pipx install pre-commit` or `pip install pre-commit`
   - Then run: `pre-commit install`
3. For Node formatting in hooks, ensure Node 18+ is available.

### Pre-commit checks
We use `pre-commit` to run common checks (YAML, trailing whitespace, Prettier, Ruff, Markdown lint). CI will enforce these on every PR.

### Coding standards
- Keep functions small and readable; prefer clarity over cleverness.
- Match existing code style; format with Prettier/Ruff.
- Write tests where applicable.

### PR checklist
- [ ] Follows Conventional Commits
- [ ] Descriptive title and summary
- [ ] Adds/updates docs as needed
- [ ] Passes pre-commit locally (`pre-commit run -a`)

### Security
Please report vulnerabilities privately per `SECURITY.md`.


