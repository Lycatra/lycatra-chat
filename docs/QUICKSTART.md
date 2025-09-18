### Quickstart

#### Prerequisites
- Git, Docker Desktop (optional for later), Python 3.x for pre-commit, Node 18+ for Prettier hooks.

#### Setup
1. Clone the repo and checkout a feature branch.
2. Install hooks:
   - PowerShell: `py -m pip install pre-commit` then `pre-commit install` and `pre-commit run -a`.
3. Ensure CI passes locally before pushing.

#### Branch/PR
- Branch: `feat/<scope>-<short-desc>`
- Commit: Conventional Commits
- PR: squash merge with semantic title; 1 approval; CI green


