# 🛰️ Vessl

**The Ultra-Lightweight, Self-Hosted PaaS for Developers.**

Turn any bare-metal Linux VPS into your own private Vercel, Railway, or Heroku in 60 seconds with zero-downtime deployments, automated SSL routing, and an ultra-responsive web control panel.

---

## ✨ Features

- **⚡ Blazing-Fast Go Daemon (`vessld`)**: Uses native Go concurrency and official Docker SDK with `< 30MB RAM` idle overhead.
- **💻 Self-Hosted Dashboard (`dashboard/`)**: Built with **Vite + TanStack Router + React + Tailwind CSS**. Served directly by the Go daemon. Features live `@xterm/xterm` terminal logs, real-time CPU/RAM stats, and dark-mode glassmorphism.
- **🔒 Automated Edge Routing (`Traefik v3`)**: Zero-config Let's Encrypt SSL/TLS certificates and automatic reverse proxy configuration.
- **🔐 Encrypted `.env` Vault**: AES-256 encrypted environment variables stored inside an embedded SQLite database.
- **🛡️ Modular Middleware & Security**: Built-in JWT authentication guards, RBAC enforcement (`admin`, `member`), and global CORS middleware.
- **🔄 1-Click Zero-Downtime Self-Updates**: Upgrade the Vessl control plane with a single click (`scripts/upgrade.sh`) while your deployed user applications experience **zero seconds of downtime**.

---

## 📂 Repository Layout

```text
vessl/
├── cmd/          # Go Daemon entrypoint (`main.go`)
├── internal/             # Core Go packages (horizontal layers)
│   ├── agent/            # Agent mode for remote cloud connectivity
│   ├── engine/           # Docker engine lifecycle, build strategies, deployer
│   ├── handlers/         # HTTP handlers (auth, backup, canvas, database, deployment, etc.)
│   ├── http/             # HTTP server setup, routes, CORS, auth middleware
│   ├── middleware/       # Authentication guards & CORS configuration
│   ├── models/           # Domain model structs (project, database, user, etc.)
│   ├── proxy/            # Traefik v3 reverse proxy controller
│   ├── repositories/     # SQLite data access layer (per-domain repositories)
│   ├── services/         # Business logic services (auth, cron, deploy, git, etc.)
│   └── vault/            # AES-256-GCM encryption vault for secrets
├── dashboard/            # 💻 Self-hosted dashboard — served by daemon binary

---

## 📢 Notifications Architecture (Self-Hosted)

Vessl (Self-Hosted) is completely standalone and fully owned by the administrator.
- **Database-Driven Channels:** There are NO environment variables for configuring notification channels (like SMTP, Slack, Discord). All channels are configured directly in the Dashboard UI and saved securely in the SQLite database (`ServerSettings`).
- **No Hardcoded Fallbacks:** If a channel is toggled off or credentials are not provided in the UI, Vessl simply does not dispatch messages to that channel.
- **Isolated Dispatchers:** Each notification channel (Mailer, Discord, Telegram, Pushover, Slack) is separated cleanly inside `internal/notifications/`, making it incredibly easy to contribute new channels without touching core dispatch logic.
├── web/                  # 🌐 Marketing site — `vessl.dev`
├── docs/                 # 📖 Documentation — `docs.vessl.dev`
├── bootstrap/            # 📦 One-line install server (`install.sh`)
├── scripts/              # 🛠️ System automation (`upgrade.sh`, `backup.sh`, `restore.sh`, `downgrade.sh`)
├── Dockerfile            # Multi-stage container build uniting `dashboard/` and `vessld`
├── docker-compose.yml    # Production/dev container stack with Docker socket mounting
└── Makefile              # Build, test, lint, and development automation commands
```

---

## 🚀 Quick Install (On any Linux VPS)

```bash
curl -fsSL https://get.vessl.dev | sh
```

Access your dashboard at `http://your-server-ip:3000`.

---

## ⚡ Makefile Commands & Automation

Vessl includes a comprehensive root-level `Makefile` to streamline local development, building, testing, and container deployment.

| Command                | Description                                                                                  |
| :--------------------- | :------------------------------------------------------------------------------------------- |
| `make all`             | Runs code checks (`make check`) and compiles all frontend & backend binaries (`make build`). |
| `make check`           | Formats code with `go fmt ./...` and runs static analysis with `go vet ./...`.               |
| `make test`            | Executes the complete Go unit and integration test suite (`go test ./... -v`).               |
| `make build`           | Builds both the TanStack SPA dashboard and the Go daemon (`bin/vessld`).                     |
| `make build-daemon`    | Compiles only the Go backend daemon binary into `bin/vessld`.                                |
| `make build-dashboard` | Bundles the Vite + TanStack Router frontend GUI into `dashboard/dist/`.                      |
| `make dev`             | Launches the backend daemon and dashboard dev servers concurrently (`npx concurrently`).     |
| `make dev-daemon`      | Runs the standalone Go backend server (`go run ./cmd`).                                      |
| `make dev-dashboard`   | Runs the standalone Vite frontend dev server on port `3000`.                                 |
| `make dev-web`         | Runs the Astro marketing landing page dev server (`web/`).                                   |
| `make docker-build`    | Builds the all-in-one Vessl container image via Docker Compose.                              |
| `make docker-up`       | Starts the container stack (`vessld` + `docker.sock` mount) in detached mode (`-d`).         |
| `make docker-down`     | Stops and removes the running Vessl container stack.                                         |
| `make clean`           | Removes compiled binaries, temporary build artifacts, and the `bin/` directory.              |

---

## 🛠️ Local Development & Docker Stack

### Option A: Using Makefile & Native Tooling

1. **Prerequisites**: Ensure Go 1.23+ and Node.js 20+ are installed.
2. **Environment**: Copy `.env.example` to `.env` and adjust if needed:

   ```bash
   cp .env.example .env
   ```

   > **Note on Port Conflicts (Local Dev):** If you already have Apache (`httpd`) or Nginx running on your machine, Traefik will fail to start on ports 80/443. You can easily resolve this by uncommenting and modifying the `VESSL_TRAEFIK_*` ports in your `.env` file (e.g., `VESSL_TRAEFIK_HTTP_PORT=8081`).

3. **Start Dev Environment**:

   ```bash
   make dev
   ```

   This boots the backend API (`http://localhost:8080`) and frontend TanStack router UI (`http://localhost:3000`) simultaneously.

### Option B: Using Docker Compose

Run Vessl completely inside Docker with access to the host Docker daemon (`/var/run/docker.sock`):

```bash
make docker-build
make docker-up
```

---

## 📄 License

MIT License. See `LICENSE` for details.

## 📢 Notifications Architecture (Self-Hosted)

Vessl (Self-Hosted) uses an event-driven notification dispatcher located at `internal/services/notifications`.
Unlike managed cloud environments, the self-hosted daemon **does not rely on environment variables** for configuring SMTP, Slack, Discord, or Telegram. Instead, all notification channel settings are stored persistently in the local SQLite database.

- **Dynamic Configuration:** Administrators configure their SMTP credentials and webhook URLs directly through the Vessl dashboard settings.
- **Enabled Toggles:** A channel (e.g., Email, Telegram) is considered "enabled" if its respective settings exist in the database and the `Enabled` flag is true.
- **No Global Fallbacks:** Because each self-hosted instance is fully independent, there is no fallback to a "managed" or "global" provider. If the user hasn't configured SMTP, emails simply aren't sent.

To contribute a new channel:

1. Create a new file (e.g., `slack.go`) inside `internal/services/notifications/`.
2. Implement the sending logic.
3. Hook it into the central event dispatcher.
