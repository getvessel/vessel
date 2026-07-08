# 🛰️ Vessel: The Ultra-Sleek, Lightweight Self-Hosted PaaS

> **Tagline**: _Turn any bare-metal VPS into your own private Vercel & Railway in 60 seconds._  
> **Project Name**: **Vessel** (`vessel.dev` / `github.com/vessel-run/vessel`)  
> **Mission**: Build an open-source, blazing-fast, developer-first self-hosted PaaS with a next-generation GUI, automated Docker container builds, zero-downtime deployments, Caddy edge SSL routing, and one-click self-updating/upgrading.

---

## 1. System Anatomy & Terminology

To ensure crystal clarity across the repository, Vessel follows the **Standard Go Cloud-Native Infrastructure Architecture** (modeled after standard cloud tools like Traefik, Caddy, Dokku, Grafana, ArgoCD):

1. **`www/` (Public Marketing Landing Page)**:
   - **What it is**: The public website hosted at `vessel.dev` (built with Astro / Vite / Tailwind).
   - **Who sees it**: Developers globally discovering the project, reading features, documentation, and copying the one-line install command (`curl -fsSL https://get.vessel.dev | bash`).
2. **`ui/` (Self-Hosted Web GUI Dashboard)**:
   - **What it is**: The interactive management dashboard built with **Vite + TanStack Router + React + Tailwind + `@xterm/xterm`**.
   - **Who sees it**: Anyone who installs Vessel on their VPS. When they visit `https://their-server-ip:3000` or `https://app.their-domain.com`, this is the GUI control panel they log into to deploy apps, manage databases, view live logs, and configure `.env` variables.
3. **`cmd/vesseld` & `internal/` (Orchestrator Backend & Daemon — Go / Golang)**:
   - **What it is**: High-performance backend orchestrator built in **Go**.
   - **What it does**: Talks directly to the Docker socket (`/var/run/docker.sock`), manages Caddy SSL/reverse-proxy rules, executes git webhooks, streams WebSocket terminal logs (`gorilla/websocket`), manages SQLite state, and handles self-upgrade/update commands (`deploy/upgrade.sh`).
4. **`deploy/` (Installation, Bootstrap & Self-Upgrade Engine)**:
   - **What it is**: Consolidated shell automation allowing one-click install (`install.sh`), pre-flight checks (`bootstrap.sh`), and zero-downtime self-updating (`upgrade.sh`).

---

## 2. Comprehensive Repository Structure

Vessel uses the **Standard Go Project Layout**:

```
vessel/
├── cmd/
│   ├── vesseld/               # Main entrypoint (`main.go`) for the Go server daemon
│   └── vessel/                # CLI utility (`vessel deploy`, `vessel status`)
│
├── internal/
│   ├── types/                 # Domain structs (`ContainerHealth`, `ProjectConfig`, `SystemInfo`)
│   ├── api/                   # REST & WebSocket endpoints serving `ui/` dashboard
│   ├── orchestrator/          # Native Docker SDK container & volume management
│   ├── proxy/                 # Caddy v2 dynamic `Caddyfile` generator & hot-reloader (`caddy reload`)
│   ├── store/                 # Embedded SQLite (`CGO_ENABLED=0` sqlite) + AES `.env` vault
│   └── updater/               # Self-upgrade manager (triggers `deploy/upgrade.sh`)
│
├── ui/                        # 💻 SELF-HOSTED WEB GUI (What the VPS installer sees & uses)
│   ├── src/
│   │   ├── routes/            # TanStack Router type-safe routes (Dashboard, Projects, Databases)
│   │   ├── components/        # Glassmorphism UI cards, xterm.js live terminal, `.env` vault
│   │   └── hooks/             # WebSocket streaming hooks, TanStack Query mutations
│   ├── index.html             # Main entry point for local GUI
│   └── package.json           # Frontend dependencies (@tanstack/react-router, @xterm/xterm)
│
├── www/                       # 🌐 PUBLIC MARKETING LANDING PAGE (vessel.dev)
│   ├── src/                   # Hero section, one-click copy install banner, interactive demo screenshots
│   └── package.json           # Astro / Vite landing page dependencies
│
├── deploy/                    # 📦 INSTALLATION & SELF-UPGRADE SCRIPTS
│   ├── install.sh             # One-click curl installation (`curl -fsSL https://get.vessel.dev | bash`)
│   ├── upgrade.sh             # In-place self-upgrade script executed by GUI/backend during updates
│   └── bootstrap.sh           # Linux OS check, Docker checking, systemd unit provisioning
│
├── data/                      # 💾 PERSISTENT DATA DIRECTORY (Mounted in Docker/Systemd)
│   ├── vessel.db              # SQLite primary database storing projects, users, and encrypted envs
│   ├── backups/               # Automated daily backups (`.tar.gz`)
│   └── caddy/                 # Caddy certificates and dynamic `Caddyfile`
│
├── docker-compose.yml         # All-in-one local dev & production container deployment
├── Dockerfile                 # Multi-stage production build uniting `ui/` and `vesseld` binary
├── Makefile                   # Build automation (`make build-daemon`, `make dev`)
└── README.md                  # Main open-source repository documentation
```

---

## 3. How Self-Updates & Upgrades Work (`deploy/upgrade.sh`)

One of Vessel's killer open-source features is **1-Click Self-Updating**:

1. **Check for Updates**: The GUI (`ui/`) pings `GET /api/system/info` on the orchestrator (`vesseld`). The orchestrator checks GitHub Releases (`github.com/vessel-run/vessel/releases/latest`).
2. **One-Click Upgrade**: If a new version (`v1.2.0`) is found, the GUI displays a banner with an **"Upgrade Now"** button.
3. **Safe In-Place Execution**:
   - When clicked, `vesseld` triggers `deploy/upgrade.sh` in the background.
   - `upgrade.sh` automatically creates a snapshot (`/data/backups/vessel-pre-upgrade-[timestamp].db`).
   - Pulls the latest Docker image (`docker compose pull vessel`).
   - Gracefully recreates only the `vessel` control plane container (`docker compose up -d --force-recreate vessel`).
   - Running user app containers (`Vite`, `NestJS`, `Postgres`, `Redis`) are **completely untouched and experience ZERO downtime during the Vessel upgrade**.
