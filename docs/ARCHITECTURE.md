### Architecture (Go)

- **API/Tool Server (MCP-ish)**: Go HTTP server exposing tools as JSON endpoints.
- **LLM Agent**: Calls tools when needed; initial MVP can be simple function-calling via API.
- **Matrix Bot**: Bridges chat to server; listens for releases/commands; posts messages and reacts.
- **Update Orchestrator**: Watches releases, asks for approval in Matrix, runs Docker updates with rollback.

Key principles: minimal deps, small binaries, stateless HTTP, clear tool boundaries.


