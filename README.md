# 🛰️ Vessel

**The Ultra-Lightweight, Self-Hosted PaaS for Developers.**

Turn any bare-metal Linux VPS into your own private Vercel, Railway, or Heroku in 60 seconds with zero-downtime deployments, automated SSL routing, and an ultra-responsive web control panel.

---

## ✨ Features

- **⚡ Blazing-Fast Go Daemon (`vesseld`)**: Uses native Go concurrency and official Docker SDK with `< 30MB RAM` idle overhead.
- **💻 Main Panel Dashboard (`dashboard/`)**: Built with **Vite + TanStack Router + React + Tailwind CSS**. Features live `@xterm/xterm` terminal logs, real-time CPU/RAM stats, and dark-mode glassmorphism.
- **🔒 Automated Edge Routing (`Caddy v2`)**: Zero-config Let's Encrypt SSL/TLS certificates and automatic reverse proxy configuration.
- **🔐 Encrypted `.env` Vault**: AES-256 encrypted environment variables stored inside an embedded SQLite database.
- **🛡️ Modular Middleware & Security**: Built-in JWT authentication guards, RBAC enforcement (`admin`, `member`), and global CORS middleware.
- **🔄 1-Click Zero-Downtime Self-Updates**: Upgrade the Vessel control plane with a single click (`scripts/upgrade.sh`) while your deployed user applications experience **zero seconds of downtime**.

---

## 📂 Repository Layout

```table
vessel/
├── cmd/vesseld/          # Go Daemon entrypoint (`main.go`)
├── internal/             # Core Go packages
│   ├── api/              # HTTP handlers & route registration (`server.go`, `routes.go`)
│   ├── middleware/       # Authentication guards (`auth.go`) & CORS configuration (`cors.go`)
│   ├── orchestrator/     # Container, database, storage & cron job lifecycle managers
│   ├── proxy/            # Caddy v2 reverse proxy controller
│   ├── services/         # Token, Git, Cron, and Service Linker services
│   ├── store/            # SQLite embedded database and state management
│   └── types/            # DTOs, API payloads, and internal data structures
├── dashboard/            # 💻 Main Panel Dashboard (TanStack Router + React SPA)
├── website/              # 🌐 Public Marketing Landing Page (vessel.dev)
├── get-vessel/           # 📦 Installation host (`install.sh`, `upgrade.sh`)
├── scripts/              # 🛠️ System automation (`upgrade.sh`, `backup.sh`, `restore.sh`)
├── Dockerfile            # Multi-stage container build uniting `dashboard/` and `vesseld`
├── docker-compose.yml    # Production/dev container stack with Docker socket mounting
└── Makefile              # Build, test, lint, and development automation commands
```

---

## 🚀 Quick Install (On any Linux VPS)

```bash
curl -fsSL https://get.vessel.dev | sh
```

Access your dashboard at `http://your-server-ip:3000`.

---

## ⚡ Makefile Commands & Automation

Vessel includes a comprehensive root-level `Makefile` to streamline local development, building, testing, and container deployment.

| Command                | Description                                                                                  |
| :--------------------- | :------------------------------------------------------------------------------------------- |
| `make all`             | Runs code checks (`make check`) and compiles all frontend & backend binaries (`make build`). |
| `make check`           | Formats code with `go fmt ./...` and runs static analysis with `go vet ./...`.               |
| `make test`            | Executes the complete Go unit and integration test suite (`go test ./... -v`).               |
| `make build`           | Builds both the TanStack SPA dashboard and the Go daemon (`bin/vesseld`).                    |
| `make build-daemon`    | Compiles only the Go backend daemon binary into `bin/vesseld`.                               |
| `make build-dashboard` | Bundles the Vite + TanStack Router frontend GUI into `dashboard/dist/`.                      |
| `make dev`             | Launches the backend daemon and dashboard dev servers concurrently (`npx concurrently`).     |
| `make dev-daemon`      | Runs the standalone Go backend server (`go run ./cmd/vesseld`).                              |
| `make dev-dashboard`   | Runs the standalone Vite frontend dev server on port `3000`.                                 |
| `make dev-website`     | Runs the Astro marketing landing page dev server (`website/`).                               |
| `make docker-build`    | Builds the all-in-one Vessel container image via Docker Compose.                             |
| `make docker-up`       | Starts the container stack (`vesseld` + `docker.sock` mount) in detached mode (`-d`).        |
| `make docker-down`     | Stops and removes the running Vessel container stack.                                        |
| `make clean`           | Removes compiled binaries, temporary build artifacts, and the `bin/` directory.              |

---

## 🛠️ Local Development & Docker Stack

### Option A: Using Makefile & Native Tooling

1. **Prerequisites**: Ensure Go 1.23+ and Node.js 20+ are installed.
2. **Start Dev Environment**:

   ```bash
   make dev
   ```

   This boots the backend API (`http://localhost:8080`) and frontend TanStack router UI (`http://localhost:3000`) simultaneously.

### Option B: Using Docker Compose

Run Vessel completely inside Docker with access to the host Docker daemon (`/var/run/docker.sock`):

```bash
make docker-build
make docker-up
```

---

## 📄 License

MIT License. See `LICENSE` for details.
