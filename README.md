# Lycatra-chat Project – LLM-Powered Assistant with

# Minimal Dependencies

## Overview

**Lycatra-chat** is a local AI assistant server designed to provide intelligent, conversational responses by leveraging
a Large Language Model (LLM) and custom tools, all while keeping the technology stack lean. The goal is to
allow users (or developers) to query data and control services through natural language. If a user issues a
predefined command (for example, /list updates), the server returns the raw data from the
appropriate endpoint. If the user instead asks a free-form question or request, the query is routed to an
LLM, which will use the available **tools** (our server’s endpoints) to fetch or manipulate data and then return
a friendly, summarized answer. This approach ensures that **non-command chat messages receive a
helpful AI-generated response instead of raw JSON or database dumps**, improving usability.

**Key Objectives:**

Minimal Dependencies: The project strives to use as few external libraries and services as possible.
Both Go and Rust are strong candidates for the implementation due to their performance and rich
standard libraries, allowing us to build a web/API server without heavy frameworks. Go’s standard
library (the net/http package) provides everything needed to build a high-performance web
server without relying on external frameworks. Similarly, in Rust, one can use a low-level HTTP
library like Hyper to avoid full-stack frameworks – as one developer noted, “hyper is all you need to
build an application like this without much incidental complexity”. By sticking close to the standard
libraries (or minimal well-vetted libraries), we reduce bloat, simplify the build, and make the system
easier to maintain.
LLM Integration via Tools: Rather than the LLM having unrestricted access or being tightly coupled,
the server will expose a set of tools (endpoints) that the LLM can invoke to get data or perform
actions. This concept follows the emerging Model Context Protocol (MCP), which “allows servers to
expose tools that can be invoked by language models”. Think of it as an API designed for AI: MCP
servers can expose data through Resources (analogous to GET endpoints) and functionality
through Tools (analogous to POST endpoints) specifically for LLM use. The LLM, acting as an
agent, can query a list of available tools and call them as needed. For example, a tools/list
request would let the LLM discover that there is a tool called "list_updates", and then a
tools/call with that name could be invoked to fetch the latest updates. Our server will
implement this tool interface (using JSON over HTTP or another simple transport), so that any AI
model or agent that understands MCP or a similar JSON-RPC tool syntax can interact with it
seamlessly. This enables dynamic, AI-driven use of the endpoints – the LLM can decide which tool to
call based on the user’s query, retrieve the data, and then format a natural language answer for the
user.

Fast Service Updates with Minimal Downtime: Lycatra-chat is intended to run in a service-oriented manner
such that updates can be rolled out quickly without significant interruptions. We plan to achieve
near zero-downtime deployments by using proven DevOps strategies. In practice, this could mean
running multiple instances or containers and using a rolling update or blue-green deployment
approach. Rolling updates gradually replace instances of the service one by one, ensuring the
application remains available and allowing quick rollback if an issue is detected. Blue-green
deployments maintain two environments (blue=current, green=new) and switch traffic to the new
one only after it’s fully ready, thereby eliminating downtime during software updates – the
cutover is just a quick routing switch. If a problem is discovered, traffic can swiftly be switched
back to the old version, making rollback almost instantaneous. Initially, while running locally, we
might simply restart the server quickly for updates (the service should start up in seconds). But as
we evolve, we’ll incorporate these deployment techniques to meet the priority of minimal
downtime and easy rollback. This will be especially important once we have multiple users
relying on the service or when we deploy to a production environment.
Easy Extensibility: The system should be straightforward for developers to extend with new
capabilities. Adding a new API endpoint (which, in MCP terms, means adding a new tool or
resource) should be a clear and simple process (see “Adding New Endpoints (Tools)” below for
documentation). This ensures the project can grow to cover more functionalities as needed, without
requiring a complete overhaul. We will provide guidelines for creating new endpoints so that any
developer on the team can implement additional features in a consistent manner. Each new tool will
automatically become accessible to the LLM agent as well, meaning the AI assistant’s knowledge and
abilities grow alongside our endpoint library.
Local-First and Self-Hosted: In the initial phase, all components will run on local servers or local
containers. We are avoiding external cloud dependencies to maintain control and privacy. The
current setup already includes a local Matrix Synapse server (running in a Docker container) for chat,
and an LLM web UI (for instance, an OpenWebUI or similar interface for local models) – all on a
Windows machine. Lycatra-chat will integrate with this environment. This means developers can run and
test everything on their development machines or local network. In the future, we plan to deploy to
a dedicated Talos Linux server (which will run Kubernetes on Docker) once it’s available, to provide
more scalable and robust infrastructure (this is covered in the TODO / Future Work section). Talos is
an ultra-minimal, secure OS designed specifically for containers/K8s, which should make our
deployment consistent and reliable when we migrate to it.
In summary, Lycatra-chat aims to provide a **conversational interface to your services and data**, using a lean
backend that is easy to maintain. Users can interact through chat, getting either direct data (for explicit
commands) or AI-curated answers (for natural language queries), and developers can iterate quickly with
minimal fuss in the codebase or deployment.

## Team workflow and getting started

- **Branching**: Create a branch from `main` using `type/scope-short-description` (e.g., `feat/api-tools`, `fix/build-error`).
- **Commits**: Use Conventional Commits (e.g., `feat(chat): add message send button`).
- **Pull Requests**: Open PRs into `main`. Squash merge with a semantic title. At least one review required.
- **Local checks**: Install and run pre-commit hooks:
  - Windows PowerShell: `py -m pip install --upgrade pip pre-commit` then `pre-commit install` and `pre-commit run -a`.
  - Alternatively: `pipx install pre-commit` (if using pipx).
- **CI enforcement**: GitHub Actions run pre-commit, branch-name lint, semantic PR title check, super-linter, and stale triage.
- **Issues/PR templates**: Use the provided templates when opening issues/PRs.
- **Security**: Report vulnerabilities privately per `SECURITY.md`.


## Architecture and Components

**High-Level Architecture:** The system is composed of a few key components that work together:

**1. The API/Tool Server (MCP Server):** This is the core of Lycatra-chat – essentially a custom web server
that exposes various endpoints (tools/resources) and follows a JSON-based protocol that an AI agent
can use. We’ll implement this server in either Go or Rust (to be decided, see **Tech Stack** below),

focusing on simplicity and speed. The server’s responsibilities include: hosting HTTP routes or RPC
methods for each tool, executing the corresponding logic (e.g., querying a database or performing a
computation), and returning results in a structured format. If we adhere to the Model Context
Protocol, the server will handle requests like tools/list (return the list of available tool
definitions) and tools/call (execute a specific tool with given parameters). Under the
hood, this could be implemented via a JSON-RPC mechanism or even a simple REST pattern – the key
is that the interface is consistent and machine-readable so that an LLM can invoke it. We might
utilize an open-source library for MCP to save time (for example, MCP-Go if we use Go, or an
equivalent in Rust), but we will keep any such library usage minimal. Notably, the tools we define
can either fetch data (like reading from a file or database, calling an external API, etc.) or perform
actions (like triggering a service restart, toggling a feature, etc.), depending on project needs.
Initially, an example tool will be list_updates, which might gather information about recent
service updates or changes and return them. Over time we’ll add more tools as required.
**2. The LLM Agent:** The intelligence layer is provided by a Large Language Model. This could be an
external API (like OpenAI GPT-4) or a local model (running via something like OpenAI’s text-
generation-webui, ollama, or another local inference engine). The key is that the LLM is
configured with knowledge of how to use tools. When a user sends a free-form message, we pass it
(and possibly some context or conversation history) to the LLM. The LLM will then analyze the
request and, if needed, decide to invoke one of the tools exposed by the Lycatra-chat server. Thanks to
the standardized tool interface, the LLM knows how to, say, call the list_updates tool if the user
asked “What changed in the system recently?” The LLM effectively becomes an **agent that can
perform tool-using reasoning**. It will call the tool via the API, get back structured data, and then
incorporate that data into a natural language answer. The final answer (text) is then sent back to the
user. This design follows the retrieval-augmented generation pattern and ensures the user gets an
**LLM-curated response backed by real data from our endpoints**. We will likely start by using a
simple integration: for example, if using OpenAI, we can use the Functions API (which is analogous
to tools) or if using a local model, perhaps a framework like LangChain or an MCP-compatible
orchestrator. **Important:** The agent will only trigger on messages that are not recognized
commands (no leading slash). This mirrors approaches like the _Chaz_ Matrix bot, where _“any message
that isn’t a command causes the entire conversation (up to the last reset) to be sent to the AI for
completion”_. In our case, if a user types a normal sentence or question, the LLM agent handles it,
possibly using tools; if the user types a command (prefixed by / or another special symbol), it’s
handled directly by the server without invoking the AI.
**3. The Chat Interface (Matrix Client/Bot):** To provide a user-friendly way to interact, we leverage
our existing Matrix server (Synapse) and potentially a Matrix bot. Users can chat with the **Lycatra-chat
bot** in a Matrix room or direct message. The bot will act as a bridge: it receives messages from
Matrix and forwards them to Lycatra-chat’s API (or directly to the LLM agent as appropriate), then returns
the response as a chat message. The reason we use Matrix is because we already have it set up, and
it’s a flexible, open messaging system. By plugging our assistant into Matrix, we make it accessible
from any device using any Matrix client, and we can even invite others to interact if desired. Matrix
also gives us features like persistent history and user management out-of-the-box. In practice,
this component could be a small script or service that uses the Matrix Client-Server API or a Matrix
bot SDK to listen for messages and respond. Since our Synapse is running in a container (standard
setup), the bot can run externally and connect via HTTP API. There are existing examples of such
integrations – for instance, a simple Matrix bot can take messages from a room and send them to an

AI API and then reply, as demonstrated by projects like the matrix-llm-bot (which “redirects messages
from rooms to an LLM via API”). We will document how to set up our bot user and configure it to
point at the Lycatra-chat server’s API. This way, when a user chats with the bot in Matrix, they are
effectively interacting with the Lycatra-chat server, getting either direct data (for commands) or AI-
generated answers (for natural language), depending on their input.
**4. Data Stores / Services:** (Optional, depends on features) Lycatra-chat might interface with various data
sources. For example, if one tool needs to retrieve system metrics or database records, the server
will need access to those. Initially, we can keep things simple with local files or in-memory data for
prototypes. As the project grows, we might integrate a lightweight database or connect to external
APIs. Each integration will be encapsulated in a tool. For instance, a get_status tool could check
the status of another service via HTTP, or a search_logs tool might query a log file or database.
We will ensure these dependencies are also minimized (favoring built-in OS capabilities or simple
libraries). If any credentials or config are needed (for databases, external APIs, etc.), we’ll manage
those via environment variables or config files (keeping them out of source control).

**Workflow Example:** A typical interaction cycle might go as follows: 1. **User Input:** User types a message in
Matrix chat, e.g., “Hey, what updates were deployed this week?” (no command prefix, just a question). 2.
**Message Routing:** The Matrix bot relays this message to the Lycatra-chat server’s LLM agent interface. 3. **LLM
Analysis:** The LLM receives the message. It recognizes that answering this may require calling the
list_updates tool (which perhaps returns a summary of recent deployments). The LLM formulates a
JSON request to the Lycatra-chat API: for example, it might internally create a tools/call request for the
list_updates tool. 4. **Tool Invocation:** The Lycatra-chat server receives the call, executes the
list_updates handler (perhaps this reads a changelog file or checks a deploy log), and returns the
result, say a JSON list of updates or a text snippet describing them. 5. **LLM Response Formation:** The LLM
takes that data and composes a helpful answer like: “This week, the user service was updated to v2.3 (with
bug fixes) and a new logging system was deployed on Tuesday. No other major updates were
recorded.” 6. **Bot Reply:** The composed answer is sent back through the Matrix bot to the chat, and the user
sees the nicely formatted response. The user did not have to sift through raw JSON or logs – the LLM
interpreted it for them. 7. **Alternate Path (Command):** If the user had instead typed a command like
/list updates, the Matrix bot could detect the slash command and directly call the Lycatra-chat API for
list_updates (skipping the LLM). The server would return maybe a raw JSON or a plaintext list of
updates. The bot could either format this minimally or just post it as a pre-formatted block. This is useful for
developers who want the raw data quickly, but not necessary for end-users.

This architecture ensures flexibility: **developers and power users can use command syntax to get direct
results, while other users can ask in plain language and get AI-enhanced answers.** It also cleanly
separates concerns: the Lycatra-chat server focuses on data and actions (tools), and the LLM focuses on
understanding requests and presenting results. The interface between them (MCP/JSON tools) is
standardized, which means we could even swap out the LLM or the tool implementations independently in
the future.

## Tech Stack and Dependencies

**Programming Language:** _To Be Decided – Go or Rust._ Both Go and Rust are modern, efficient languages
with strong support for systems programming and web services, and importantly, both can produce a
single static binary for easy deployment. We outline the pros of each to help in final selection:

Go – Pros: Very simple learning curve, especially if developers are new to it. It has an excellent
standard library for HTTP servers and JSON, meaning we can avoid adding frameworks. As noted,
the net/http package in Go is sufficient for a solid web API; “Go’s standard library (net/http)
really has everything you need to build a solid, high-performance web server without relying on any
external libraries”. We can define route handlers for our JSON-RPC endpoints easily. Goroutines
and channels make concurrency straightforward if needed (e.g., handling multiple requests).
Compilation is fast, and cross-compiling for Linux (for eventual Talos deployment) from a Windows
dev machine is trivial. The team’s familiarity with Python/JavaScript will find Go’s syntax and garbage-
collected model quite approachable. Also, Go’s binary can be small (a few tens of MB) and has no
runtime dependencies – just run the binary. One potential library we might consider is the MCP-Go
library, which can save us time by providing out-of-the-box support for the Model Context
Protocol (tool registration, JSON-RPC handling, etc.). Using it would slightly increase dependencies,
but it is well-maintained and geared exactly towards our use case, allowing us to “focus on building
great tools” while it handles protocol details. Even with that, the overall dependency footprint
remains small (MCP-Go itself is primarily a wrapper around standard net/http and JSON encoding
under the hood).
Rust – Pros: High performance and memory safety. Rust’s package ecosystem (Cargo) makes it easy
to include only what we need. We can build a web server with minimal crates – for example, use
hyper (a low-level async HTTP library) instead of a full framework. This aligns with the “no heavy
framework” ethos: one Rust tutorial follows the rule “if the crate calls itself a framework, don’t use it”,
and demonstrates that using Hyper plus a few small utility crates is enough to build a web API
. Hyper gives us full control and is very fast. Rust’s strong type system can help catch errors early,
which is beneficial as the project grows. The downside is that Rust has a steeper learning curve, and
development might be a bit slower at first. However, the team can gain experience in a systems
language, which is part of the goal (we expressed interest in knowing Rust better). The compiled
binary will be similarly self-contained. We would need to manage asynchronous calls (for I/O and
possibly for integrating with the LLM or other services) – Tokio runtime is commonly used (Hyper
uses it) which is one additional dependency but pretty essential in Rust for concurrency. There may
also be an **MCP Rust** implementation or we could adapt examples from the spec. If no mature
library, implementing the JSON-RPC handlers in Rust manually is doable but a bit verbose (parsing
JSON, matching method strings, etc.). The team should weigh the learning benefits vs. the
immediate productivity.

**Current Leaning:** If rapid development and simplicity are top priority, **Go might be the better choice** to
start with, given its minimal setup and the fact that one team member can spin up a basic HTTP JSON server
in a few lines. On the other hand, if we treat this project as an opportunity to delve into Rust, we can
absolutely build it in Rust while keeping dependencies minimal (Hyper + Serde for JSON would cover most
needs). Both languages align with our **no heavy runtime** requirement (no JVM, no Node.js, etc.) and can
produce small, efficient services. Ultimately, we should choose one to avoid fragmenting the codebase. For
now, we will proceed assuming Go for examples (since it’s likely faster to implement initially), but will note

where adjustments might be needed for Rust. **Either way, the architecture and high-level design remain
the same.**

**Libraries/Dependencies:** Aside from the language runtime itself, we plan to use: - JSON handling: built-in
encoding/json (Go) or serde (Rust). - HTTP server: built-in net/http (Go) or hyper (Rust). -
(Optional) MCP support: possibly the mcp-go library for Go or a community crate for Rust if available,
to handle the tool interface. This is not strictly required; we could implement the JSON-RPC methods
ourselves if we want absolute minimal external code. However, using a well-tested library for MCP can
ensure compliance with the spec and save development time (we won’t have to write our own JSON-RPC
request router and validator). The library can also manage some advanced features (like sessions,
streaming, etc.) which we might leverage later. - Matrix bot: For the Matrix integration, we can use the
Matrix Client-Server REST API directly (via simple HTTP calls, which keeps dependencies low – just use net/
http or Rust’s Reqwest to hit the endpoints). Alternatively, there are lightweight Matrix client libraries for
both Go (mautrix-go) and Rust (ruma or the mentioned headjack used in Chaz). Using an HTTP
client might actually be simplest: we can have a small Go routine or Rust async task long-poll (sync) or
WebSocket (if using Matrix’s newer push) for new messages. Given that Synapse is local, network latency is
minimal. For sending, it’s just an HTTP POST to the _matrix/client/v3/rooms/<roomId>/send
endpoint. We will likely create a separate small program or even include it as part of the Lycatra-chat binary
(though separating concerns might be cleaner). Since this is an internal tool, using a simple approach (like
polling every second for new events) is fine and avoids adding complex dependencies. The bot will need an
access token (we’ll create a Matrix user for Lycatra-chat Bot) – we’ll document configuration for this.

Other than the above, we aim **not** to include databases or heavy services in the initial version. If we need
persistence (for example, storing user-specific preferences or logging interactions), we might use a
lightweight embedded database or just write to a file (again focusing on simplicity first). The whole system
should remain **Docker-friendly** – meaning minimal state kept on disk so that containers can be ephemeral.
For now, local development can just use the host filesystem for any needed storage.

**Summary of Chosen Stack (to be confirmed):**

Language: Go 1.21+ (preferred for quick start) or Rust 1.72+ (if we commit to the Rust route).
Platform: Linux containers for deployment (but develop on Windows/Mac is fine). We will
containerize the service to run on Docker since eventual target is Talos/K8s. During development,
running natively is fine too.
Interface Protocol: HTTP + JSON (simple REST or JSON-RPC endpoints). Conforms to MCP syntax for
tool calls so that LLMs can interface easily.
LLM Service: Initially OpenAI API (GPT-4) for convenience (no setup needed, just an API key) – this
gives us a quick way to test the agent behavior. In parallel, we are exploring local models via the
existing OpenWebUI/Ollama setup. The architecture allows either to be used. (We might
structure the LLM agent as a separate module that can call either an external API or a local model
process).
Chat Interface: Matrix (Synapse) with a bot user. Already running Synapse in Docker; we will just
add our bot. No additional chat server needed – reuse what’s in place for simplicity.

## Development Setup and Running the Project

One of our goals is to keep the development environment straightforward so any team member can quickly
get the project up and running. Below are instructions for setting up and running Lycatra-chat in a dev setting.

### Prerequisites

Go or Rust Toolchain: Depending on the chosen language, install either Go (https://go.dev/dl/) or
Rust (https://rustup.rs/). Both are cross-platform and easy to install. Ensure your PATH is set for the
toolchain. For Go, verify by running go version; for Rust, rustc --version and cargo --
version.
Git: The project will be managed in a Git repository. Clone the repository to your local machine:
git clone <repo_url> (repo URL to be provided once initialized).
Docker (optional but recommended): Since our other services (Matrix Synapse, etc.) are running in
Docker, having Docker allows you to run Lycatra-chat in a container as well, or at least to coordinate with
other services using something like Docker Compose. It’s not strictly required for coding or unit
testing, but it will be for integration testing (especially before deploying to Talos). Ensure Docker
Desktop or Docker Engine is installed if you plan to use it.
### Project Structure

After cloning, the repository structure will look something like:

Lycatra-chat-project/
├── cmd/ # Application entry points
│ └── Lycatra-chat/ # Main Lycatra-chat server binary source
│ └── main.go # Or main.rs (depending on language)
├── internal/ # (if Go) or src/ (if Rust) – internal library code
│ ├── tools/ # Implementation of each tool/endpoint
│ ├── llm/ # LLM integration logic (tool invocation agent)
│ └── matrix_bot/ # Matrix bot integration (listener & sender)
├── pkg/ # (optional) for shared packages or library code
├── Dockerfile # Dockerfile to containerize the Lycatra-chat server
├── docker-compose.yml # Compose file to orchestrate Lycatra-chat + Synapse +
others (dev/test only)
├── README.md # (This document) project documentation
├── DESIGN.md # Additional design notes (if needed)
└──... (other docs or config files, e.g.,.env.example for environment
variables)
_(Note: This is a proposed structure. In Go, we typically separate cmd/ for binaries and pkg/ or internal/
for libraries. In Rust, everything is under src/ by default, but we can achieve a similar organization with
modules. We will adapt to language conventions accordingly.)_

Key files initially: - **main.go / main.rs:** The entry point that starts the HTTP server, registers tools, and
begins listening for requests. - **Tools implementation:** Could be individual files for each tool (e.g., tools/
list_updates.go or tools/list_updates.rs) containing the logic for that endpoint. - **LLM
integration:** Possibly a module that handles sending queries to the LLM (OpenAI or local) and processing
the response (including hooking into the tools if needed – though if using MCP, the LLM agent might handle
tool calling logic externally; see design note below). - **Matrix bot:** A module or service that polls for
messages and relays them. We might also include a simple script or separate binary for the bot (to keep the
Lycatra-chat server independent). Alternatively, the Lycatra-chat server could itself connect to Matrix and act as the
bot. This is a design choice: running them as separate processes might be cleaner (one process for the API,
one for the chat interface), but it’s also possible to integrate (just spawn a goroutine or tokio task for the bot
loop inside the main program). We’ll decide based on simplicity – a separate small program might actually
reduce complexity of each component.

### Running in Development

**1. Configuration:** First, copy or create a configuration file if needed (e.g.,.env or a config.yaml). At
minimum, you may need to provide: - **Matrix Bot Credentials:** The homeserver URL (e.g., [http://](http://)
localhost:8008 for local Synapse), the bot user ID (@Lycatra-chat:yourdomain), and access token or
login/password. This allows the bot to connect to Matrix. - **LLM API Key or Endpoint:** If using OpenAI, set
    OPENAI_API_KEY. If using a local LLM, specify the endpoint (e.g., [http://localhost:4000/api](http://localhost:4000/api) for
OpenWebUI, etc.). We will not commit secrets to the repo; instead, use environment variables or a.env
file that is gitignored. - Optionally, any other service credentials if a tool needs it (not in initial version). - For
now, defaults can be coded (e.g., assume Matrix on localhost, etc.), so config may be optional.
**2. Launch Supporting Services:** Ensure the Matrix Synapse server is running (if you have Docker, docker
ps to see if the synapse container is up, or start it via docker-compose up synapse if using a compose
file). Also launch any LLM local server if needed (for example, if using an Ollama or other model server, start
it). If using OpenAI API, no local process needed, just internet access.
**3. Run Lycatra-chat Server (Development Mode):** - **Using Go:** Navigate to the project directory and run go
run./cmd/Lycatra-chat (this compiles and runs the server). You should see logs indicating the server has
started, e.g., “Server starting on [http://localhost:8080”](http://localhost:8080”) (the default port can be 8080 or configurable). - If
you want to test the MCP interface manually, you can use curl or a tool like httpie to hit the
endpoints. For example, to list tools: curl -X POST -H "Content-Type: application/json" -d
'{"jsonrpc":"2.0","id":1,"method":"tools/list"}' [http://localhost:8080/mcp.](http://localhost:8080/mcp.) The server
should respond with a JSON listing of available tools and their schemas. Or to test a specific tool directly:
    curl -X POST -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","id":
2,"method":"tools/call","params":{"name":"list_updates","arguments":{}}}' [http://](http://)
localhost:8080/mcp and you should get back a result with either the data or an error if something is
wrong. This can confirm that the backend logic works. - **Using Rust:** Run cargo run in the project
directory (or cargo run --bin Lycatra-chat if we have multiple binaries). This will compile and run. The rest
is analogous to Go – the server should start and listen on a port (likely 8080 as well). Use similar curl
commands as above (the endpoints might be at e.g. [http://localhost:8080/tools/call](http://localhost:8080/tools/call) depending
on how we implement routing – details will be in the documentation/comments of the code). - In either
case, the server by default runs in a single-process, single-instance mode. It’s fine for dev. If you make
changes to code, just stop (Ctrl+C) and rerun the command to see changes. Because we’re not using hot-

reload frameworks (keeping things simple), you need to rebuild to apply changes. Both Go and Rust
compile quickly for our project size (fractions of a second to a couple seconds).

**4. Run Matrix Bot (if separate):** If we implement the Matrix bot as a separate program, you would run it
similarly: - Possibly go run./cmd/matrixbot or cargo run --bin matrixbot. This should start a
process that connects to the Matrix homeserver and begins listening for messages (likely using the Matrix
sync API). - Ensure the bot is invited to the room or knows which room to monitor. In a simple setup, we
might have the bot auto-join a specific room or listen to its DM. - Then try sending messages in Matrix: e.g.,
open Element (or any Matrix client) as a normal user, start a chat with the Lycatra-chat bot user, and send a
greeting or ask a question. The bot process logs should show it received a message and forwarded it to the
Lycatra-chat server, and you should soon see a reply. For a command, try /list updates in the chat – the bot
should catch the command (perhaps by pattern or because Matrix clients don’t send it to the bot if it’s a
command? We might need to configure how to capture commands; possibly just treat everything as a
message and decide in our code if it starts with /). - In case of issues, check the bot logs for errors (auth
issues, etc.), and also the Lycatra-chat server logs (to see if it received the request or if the LLM call failed, etc.).
**5. Testing the LLM Calls:** If the LLM integration is configured, test a query that triggers it. For instance, ask
the bot in chat: “Could you list the recent updates in a summary?” This should result in the LLM being
invoked. If using OpenAI, ensure the API key is correct and observe if the assistant responds sensibly with a
summary containing real data. If using a local LLM, ensure that service is running and that Lycatra-chat’s LLM
module can reach it (maybe check by calling the local LLM API directly once).

Throughout development, you can run unit tests if available (go test./... or cargo test) which we
will write for critical components (especially the tool handlers logic). We will also create some integration
tests for the JSON interface.

**Developer Convenience:** The environment is kept simple intentionally: - No database that you need to set
up for basic functionality (unless your feature requires one – and if so, we’ll likely use SQLite or similar for
ease). - The services can all run on localhost with default configs. We will supply example config values for
things like the Matrix bot in a.env.example file. - A Docker Compose configuration will be provided to
launch the whole stack (Synapse, Lycatra-chat, possibly a local LLM container if we have one). This is useful for
end-to-end testing. For daily dev, you might just run Lycatra-chat directly and rely on an already-running
Synapse container. - We avoid any proprietary IDE requirements or heavy build scripts – just standard Go or
Rust commands.

**Note on MCP Tools and LLM Agent:** If we use an **MCP-compliant client** (for example, some devs might
use a tool or an IDE plugin that speaks MCP to interact with Lycatra-chat), we can directly connect the LLM to our
server via the MCP protocol. Otherwise, our integration as described (Matrix + custom code) works fine.
MCP is basically JSON-RPC as shown, so even our Matrix bot could technically send MCP-formatted requests
on behalf of the LLM. The separation in development is: implement all tools and confirm they work via the
API; then separately implement how the LLM will call them (either through our code or by relying on the
LLM’s abilities if it has built-in support for tool usage).

## Adding New Endpoints (Tools)

One of the design requirements is making it easy to extend Lycatra-chat with new functionality. This section
serves as a developer guide for creating new endpoints – referred to as **tools** in the MCP context – and
exposing them to both direct users and the LLM.

When asking for new features, we (or the product team) will often describe them in terms of high-level
actions or queries. As a developer, your task is to map those requests to one or more tools. A tool typically
corresponds to a single well-defined operation or query. For example, if the new feature is _"Allow the
assistant to fetch user account details by email"_, you might create a tool called get_user_by_email that
takes an email address and returns the user’s info from the database.

**Steps to Add a New Tool:**

Design the Tool Interface: Decide on the tool’s name, description, and input/output schema.
The name should be a short identifier (e.g., "get_user_by_email", "restart_service", etc.).
Use snake_case and keep it alphanumeric. This name is what the LLM will use to invoke the tool.
Write a brief description of what the tool does. This helps both documentation and the LLM (some
LLM agents use the description to decide when to use a tool).
Define the input parameters the tool needs. Each parameter should have a type (string, number,
boolean, etc.) and possibly a description. If certain parameters are required, note that. If using an
MCP library, you’ll use functions to specify these (e.g., mcp.WithString("email",
mcp.Required(), mcp.Description("User email address")) in Go). If writing manually,
document the expected JSON structure for the call (e.g., { "name": "get_user_by_email",
"arguments": { "email": "<string>" } }).
Determine the output format. In many cases, our tools will return text or JSON data that the LLM
can parse. According to MCP, a tool result is typically either some text content or structured data.
Often, returning a human-readable text (maybe a short summary or list) is useful because the LLM
can directly forward that as an answer if needed. However, if the LLM is expected to interpret the
result, structured data (JSON) could be returned and the LLM will incorporate it. For now, returning a
text blob (or markdown) is simplest unless the use case demands structured output.
Implement the Tool Logic: This involves writing a function that performs the desired action. For
example, in Go, you might create a function func handleGetUserByEmail(ctx
context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) (following
the signature expected by the MCP-Go library), or if not using the library, you’d implement it inside
the HTTP handler case for that method. In Rust, similarly, you’d write an async function to handle it.
Within this function, do the following:
Parse and validate input parameters from the request. (The library provides helpers like
req.RequireString("email") which returns the string or an error if missing. Without the
library, you’d manually extract from JSON.)
Perform the operation: e.g., query the database for the given email. (We may have a utility module
for database interactions if applicable.)
Formulate the result. Using MCP-Go, you might use mcp.NewToolResultText(resultString)
to create a text result. Or manually, you’d craft a JSON-RPC response object with the result

content. Handle errors by either returning an error (which the library might convert into a proper
error response) or by creating an error result (e.g., mcp.NewToolResultError("message")).
Ensure thread-safety and performance considerations: If the tool might be called concurrently, make
sure to avoid global state mutations without locks, etc. Usually, reading operations are fine. If a tool
writes to something, consider synchronizing or queueing as needed (out of scope for many read-
only tools).
Example (Pseudo-code in Go):
tool:= mcp.NewTool("get_user_by_email",
mcp.WithDescription("Lookup user details by email"),
mcp.WithString("email", mcp.Required(),mcp.Description("Email address
of the user")),
)
s.AddTool(tool, func(ctxcontext.Context, reqmcp.CallToolRequest)
(*mcp.CallToolResult, error) {
email, err := req.RequireString("email")
if err!=nil {
return mcp.NewToolResultError("Email is required"),nil
}
user, err:= database.LookupUserByEmail(email)
if err!=nil {
return mcp.NewToolResultError("Lookup failed: "+ err.Error()),nil
}
if user==nil{
return mcp.NewToolResultText("No user found with that email"),nil
}
// Format user info as text
info:= fmt.Sprintf("User %s: %s, plan=%s", user.ID, user.Name,
user.Plan)
return mcp.NewToolResultText(info),nil
})
This snippet (illustrative) registers the tool on the server s. The handler fetches a user and returns
a text summary. The MCP library handles wiring this into the JSON-RPC response. If we weren’t using
the library, we’d achieve something similar by manually checking the RPC method name in a handler
function.
Register the Tool with the Server: If using the library, as shown above, you call AddTool on the
server instance to add your new tool and its handler. If not using a library, you’ll need to add a case
in your request routing logic. For instance, if we have an HTTP endpoint for /mcp, inside it we
parse the JSON and see if method == "tools/call" and params.name ==
"get_user_by_email", then call our new function. We might maintain a map of tool names to
handler functions to simplify this. It’s important that after adding a tool, it also appears in the output
of tools/list. So if doing manually, add an entry to the list output. If using the library, it updates
the list for us. The tools/list response includes each tool’s name, description, and input schema


, so verify that your new tool’s info is correctly reflected (the library will include it
automatically if registered).
Document the Tool (for the team): Update the project documentation (could be this README or a
separate TOOLS.md file) with the details of the new endpoint. Include:
Name and description (so others know it exists).
Example usage (maybe a sample JSON request and response).
Any special considerations (e.g., “this tool requires an API key set in config” or “this tool will modify
database state, use with caution”). This helps other developers and also helps when writing prompts
for the LLM (the LLM’s system prompt or instructions might include a brief of available tools).
Testing the New Tool: Write unit tests for the tool’s logic if possible. For instance, test that
get_user_by_email returns the expected result for a known user, and an appropriate message
for an unknown email. You can call the handler function directly in tests (bypassing full server). Also
test error conditions (invalid input). Then test integration:
Start the server and try a tools/list to ensure your tool is listed.
Try a direct tools/call via curl or HTTP client to see it working end-to-end.
If the Matrix bot/LLM is running, try asking the assistant to use it (for example: “Do we have an
account under the email alice@example.com?”) and see if it triggers the tool (you might check server
logs to confirm the tool was called) and returns a correct answer. Adjust as needed (sometimes you
might need to tweak the tool’s output format to be more LLM-friendly).
By following the above steps, adding new capabilities is relatively straightforward. The architecture’s
separation means you usually don’t have to touch the LLM code or the bot code at all – just add the new
tool on the server side. As soon as it’s available, the LLM can discover it and use it. (If using OpenAI function
calling, you might need to update the functions JSON you pass, but if using an MCP-aware agent, it will
fetch tools automatically via tools/list due to our declared capabilities.)

**Important:** Maintain **consistency** and **security** : - Keep naming and parameter conventions uniform (e.g.,
use lower_snake_case for JSON fields, provide clear descriptions). - Do not expose dangerous system calls as
tools without proper safeguards. Always consider what could happen if the LLM misuses a tool. We have a
“human in the loop” approach now (since only our team is using it), but as we expand, we might incorporate
confirmation steps for certain actions. - If a tool accesses sensitive data, ensure proper authentication or
permission checks if needed. In our initial local environment, this isn’t a big concern, but in a multi-user
scenario it would be. - Clean up on removal: If a tool becomes deprecated, remove its handler and any
associated resources, and update the list output.

## Deployment Strategy for Updates

As mentioned, one of Lycatra-chat’s priorities is enabling quick updates with minimal downtime. This section
outlines how we manage deployments and rollbacks in practice.


**Continuous Integration (CI):** We will set up a simple CI workflow (e.g., GitHub Actions or GitLab CI
depending on our repo) that runs tests and builds the binary (for both Windows and Linux targets,
possibly). This ensures that any new commit is verifiably building and passing tests.

**Versioning:** We’ll use semantic versioning (X.Y.Z) for the Lycatra-chat server. Small changes and new tools
increment the minor or patch version. Major changes (breaking ones) increment major. This will be
referenced in logs or status commands.

**Local Deployment (Current, Manual):** During the initial phase (all local servers on a Windows machine
with Docker), deploying a new version might be as simple as building the new Docker image and restarting
the container: - We will have a Dockerfile that produces an image of the Lycatra-chat service. For example, a
multi-stage Dockerfile that compiles the Go binary and then copies it into a scratch or alpine base for
minimal image size. - To update, we build the image (docker build -t Lycatra-chat:<tag>.) and then
update the Docker Compose file or Docker CLI to use the new image. If using Compose, we might just do
docker-compose down && docker-compose up -d to recreate the container with the new image. This
will incur a few seconds of downtime (while the container restarts). We can minimize this by instead doing a
rolling update: - Start a new container Lycatra-chat_v2 on the new image, confirm it’s running well
(healthchecks pass). - Update routing (if we had a reverse proxy or if the Matrix bot knows two endpoints,
etc.) to point to the new one, or simply instruct the Matrix bot to switch to the new instance if it’s connecting
directly (not applicable if the bot and server are in one process – but in that case, just one container restart
is fine). - Stop the old container. - This is a manual blue-green deployment on the local machine. Since the
user base is small at this stage, a few seconds of downtime might be acceptable, but we aim to practice for
zero-downtime as we’ll need it later. - Rollback in this scenario is equally manual: if the new image has
issues, re-launch the old image (we should keep the last known-good image tagged, e.g., Lycatra-chat:
1.0.0).

**Talos (Kubernetes) Deployment (Future):** Once the Talos Linux server is ready (running Docker/K8s on the
Windows machine or separate), we will move to a more robust deployment pipeline: - **Kubernetes
Deployment:** We’ll create K8s manifests or Helm charts for Lycatra-chat. The deployment will specify a replica
count (likely 2 for HA) and use a RollingUpdate strategy by default. This means when we push a new
container image (through CI/CD), Kubernetes will spin up a new pod with the new version, wait for it to
become healthy, then take an old one out, one by one. This achieves minimal downtime deployment
automatically. - **Blue-Green via Kubernetes:** Alternatively, we might implement Blue-Green by
deploying the new version as a separate Deployment or using tools like Argo Rollouts. But even without
fancy tools, a simple way is: - Deploy new version as Lycatra-chat-green while Lycatra-chat-blue is live. - Test Lycatra-chat-
green (perhaps have it join a testing Matrix room or use port-forward to hit it). - Switch a service or ingress
to point to Lycatra-chat-green (this is the cutover; since both sets of pods are running, the switch is instant for
users). - Keep Lycatra-chat-blue pods around for a short while. If an issue is reported, switch back.
Otherwise, retire Lycatra-chat-blue. - Talos will help here because it’s designed for running Kubernetes clusters
smoothly; it’s minimal and should keep overhead low, aligning with our dependency minimalism (Talos itself
provides the OS, we don’t have to manage a full Linux distro). - **Downtime goal:** With K8s, we aim for
essentially zero downtime. Users might not even notice an upgrade beyond maybe a slight delay if their
request hit during a pod restart (which K8s mitigates via readiness checks and not routing traffic to a pod
until it’s ready). - **Rollback:** Kubernetes makes rollback easy – we can keep an older ReplicaSet and just scale
it up or use kubectl rollout undo to go back to a previous version if needed. Our CI could also
maintain images with tags like prod-previous to quickly re-deploy the last version.

**State Management:** Lycatra-chat is largely stateless (especially in the HTTP layer). It should not keep significant
in-memory state that can’t be rederived, except maybe caches. This statelessness is key to quick restarts
and multiple instances. The only stateful component in our architecture is Synapse (which has its DB) and
any future databases we connect to. Those are outside Lycatra-chat and handled separately. We should ensure
that if two instances of Lycatra-chat run, they won’t conflict (they shouldn’t, as long as they don’t both try to do
something like write to the same file – which we will avoid; if needed, centralize such actions or use DB
transactions).

**Monitoring and Logging:** During deployment, it’s important to monitor if the new version is functioning
correctly: - We will implement basic **health endpoints** (like /healthz) that Kubernetes or our docker-
compose can ping. This can simply check that the server is responsive and maybe that it can reach the LLM
service. - Logging should be sufficient to spot errors. We’ll use structured logging with timestamps and
levels so that if something goes wrong after an update, we can quickly identify what and where. - If we
integrate any metrics (later on), we could track things like number of tool calls, response times, etc., to
ensure performance hasn’t regressed after an update.

**Client impact:** Because clients (Matrix bot, etc.) connect over network and likely will reconnect on failure, a
short blip in Lycatra-chat availability will cause the bot to retry. For example, if Lycatra-chat restarts, the Matrix bot’s
API call might fail; we can code the bot to handle that gracefully (retry after a few seconds). In K8s, we’d
update the bot’s configuration to use a service that load-balances between Lycatra-chat pods, so even during
rollouts the bot always finds a live instance.

**Conclusion on Updates:** We are building the operational aspect early: from the get-go, we will practice
doing fast, low-downtime updates to catch any issues in our procedure. This means even in development,
using techniques like running a second instance for testing upgrades, etc., whenever feasible. The
combination of **containerization and orchestration** gives us the tools for minimal downtime deployments

. We just have to use them in a simple, automated way. Documentation (this section) will be kept up-to-
date so that any team member can deploy a new version by following steps, and know how to rollback if
needed. As the system matures, we might automate this completely with a CI/CD pipeline.

## TODO / Future Enhancements

As the project is in its early stage, we have a roadmap of enhancements and tasks to tackle. This section
lists those to keep track and ensure we plan for them:

**1. Finalize Language Choice:** Decide between Go vs Rust for the implementation (see **Tech Stack**
above for comparison). Currently leaning towards Go for faster initial progress, but Rust remains an
option. _Todo:_ Make a decision and update the documentation to reflect the chosen stack, removing
the ambiguity so the team can focus on one. If Rust is chosen, allocate time for team ramp-up/
training.
**2. Implement Initial Tool Set:** Begin with a minimal but useful set of tools. For example:

list_updates – as discussed, fetches recent changes. Implementation: maybe read from a
CHANGELOG file or a Wiki API.

get_service_status – checks if dependent services are up (could ping an endpoint or check
docker container status).
help – returns a help message or list of commands (the LLM likely doesn’t need this, but for direct
command users).
Todo: Define and implement these initial tools. Ensure they are well-tested.
**3. LLM Integration Layer:** Right now, the design assumes the LLM can call our tools. We need to
implement that connection:

For OpenAI: Use the function calling API. We’d provide the OpenAI model a list of functions (tools)
with their schemas on each conversation. There’s some development to format our tool list into the
OpenAI function format. Todo: Write a module that given our tool definitions produces the JSON to
send to OpenAI, and handles the function call responses (calls the appropriate tool and returns the
result to the model).
For local model: Possibly set up an agent using LangChain or similar. Or use an MCP client – if there’s
an open-source agent that speaks MCP, we could use that to interface with local models. This might
require more R&D. In the interim, using OpenAI as a reliable baseline is fine.
Todo: Implement a simple LLMAgent class (or Go equivalent) that can take a user message and
return a response by coordinating with the LLM. Start with OpenAI for simplicity (since our focus is
not to train a model but to use one). Later, swap or add local model support.
**4. Matrix Bot Integration:** Set up the Matrix bot user for Lycatra-chat:

Register a new user on Synapse (or use an existing one, but better to have a dedicated bot account,
e.g., username "Lycatra-chat").
Develop the bot logic as per Architecture section. Could be in Go or Rust or even Python, but to
keep dependencies minimal, doing it in the same language as the server (and possibly as part of the
server binary) is preferred.
Test in a private room. Then potentially expose in a public room for internal users.
Todo: Write the Matrix bot connector. Ensure it handles basic commands and messages correctly.
Possibly use a small existing library if it saves time (but if not, raw HTTP calls are fine).
In the future, consider bridging to other chat platforms as well (but Matrix covers a lot and can
bridge to Slack, etc., if needed).
**5. Deploy on Talos/K8s:** Once the new Talos Linux server is ready (this will provide our Kubernetes
cluster on a Docker setup), we will containerize and deploy Lycatra-chat to it:

Write Kubernetes manifests: Deployment, Service, ConfigMap/Secrets (for API keys).
Possibly use Helm or Kustomize for easier maintenance.
Ensure the Matrix bot can run in-cluster or access the Matrix server (which might remain on the
same machine or another).
This will significantly reduce downtime on updates due to the rolling update features described. It
also sets the stage for scaling out if needed.
Todo: Setup CI to build/push Docker images. Deploy to Talos cluster in a dev namespace for testing.
Eventually cut over production usage to the Talos-hosted instance.

We will treat the current Docker on Windows as a dev environment and Talos/K8s as staging/prod
environment.
**6. Performance Tuning & Load Testing:** As we add features, we should ensure the system remains
fast. The use of Go/Rust is already a plus for performance. We should test the latency of tool calls
and LLM responses.

If local LLMs are too slow, consider adding caching for repeated queries or using faster model
inference (maybe switch to a smaller model for quick answers).
Monitor resource usage. The MCP approach allows parallel calls, but we might restrict concurrency if
the LLM or database can’t handle too many at once.
Todo: Create a few benchmark scenarios (like 10 concurrent users asking for something) and
measure response times. Optimize if needed (for example, use persistent connections, adjust
threadpool, etc.).
**7. Additional Tools & Features:** After the initial version is stable, we can expand capabilities:

More data-oriented tools (integration with project management APIs, status pages, etc. – whatever is
useful for our context).
Tools that modify state (e.g., /restart service X tool to restart a Docker container). These need
careful security and likely confirmation from the user (to ensure the LLM doesn’t restart things
without explicit ask). Perhaps require commands for destructive actions.
Multi-step workflows: The LLM can orchestrate calling multiple tools if needed (MCP allows multiple
calls in a conversation). We should test and possibly build prompt templates for that.
User-specific contexts: In the future, if multiple users use the assistant, we might give different
permissions. MCP has a notion of sessions and per-session tools. We might later utilize that
if needed (e.g., an admin user has a deploy_tool, others do not).
Todo: Gather feedback from initial users (our team) on what other tasks would be helpful to
automate via this assistant, and implement accordingly.
**8. Documentation & Knowledge Base:** Keep improving documentation.

This README should be kept up-to-date with any architectural changes.
Possibly maintain a changelog of new tools or changes for quick reference.
Write usage guides for end users if we ever expose it beyond dev team.
If we integrate a lot of domain knowledge, consider adding a resources/ tool for LLM (MCP’s
resource concept) to load reference docs into context. For example, the assistant could have
access to a “docs://README” resource (which contains this README content) to answer questions
about the system itself.
**9. Resilience and Error Handling:** As a future improvement, implement more robust error handling:

If the LLM produces an invalid tool call, handle gracefully (perhaps send an apology message or a
clarification request).

If a tool fails (exception, timeout), the system should catch it and inform the user in a friendly way,
rather than crash or hang.
Use timeouts for tool calls so an unresponsive service doesn’t hang the whole chat response
indefinitely.
Todo: Implement middleware or wrappers for tool calls to enforce timeouts and catch panics (if in
Go). The MCP spec and library encourage such safety measures.
**10. Security Audit:** Before any broader deployment, review security:

Ensure the API is not exposed to the public internet without auth (for now, our API might listen on
localhost or a private network. Eventually, if exposing, use authentication/authorization).
Lock down the Matrix bot’s capabilities (it should probably ignore messages not meant for it, etc.,
and not join arbitrary rooms).
Since Talos is hardened, that helps on OS level. But also consider network policies between services
in K8s.
Todo: Possibly implement an auth token for the MCP API if needed (so only the bot/LLM can call it,
not random users). However, since everything is internal, it’s low risk at the moment.
By addressing these TODOs over time, we will move from a basic working prototype to a robust,
production-ready AI assistant platform. The overarching theme is **keep it simple** at each step – add
functionality gradually without overcomplicating the design or introducing unnecessary dependencies.

_This README serves as the foundation for our development efforts. It contains the knowledge and guidelines
needed for the team to start building Lycatra-chat. As we progress, we’ll update this document to reflect the current
state of the project. Every team member should feel free to propose improvements or clarifications to this README
so that it remains an authoritative resource. Let’s build this minimal, powerful LLM assistant to make our lives
easier!_
