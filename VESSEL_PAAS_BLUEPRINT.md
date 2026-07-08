# 🛰️ Vessel: The Ultra-Sleek, Lightweight Self-Hosted PaaS

> **Tagline**: _Turn any bare-metal VPS into your own private Vercel & Railway in 60 seconds._  
> **Project Name**: **Vessel** (`vessel.dev` / `github.com/solomonolatunji/vessel`)  
> **Mission**: Build an open-source, blazing-fast, developer-first self-hosted PaaS with a next-generation GUI, automated Docker container builds, zero-downtime deployments, Caddy edge SSL routing, and one-click self-updating/upgrading.

---

## 1. System Anatomy & Terminology

To ensure crystal clarity across the repository, Vessel clearly distinguishes the **Public Landing Page (`website/`)** from the **Main Panel Dashboard (`dashboard/`)**:

1. **`dashboard/` (Main Panel Dashboard — What the VPS installer sees & uses)**:
   - **What it is**: The self-hosted interactive control panel GUI built with **Vite + TanStack Router + React + Tailwind + `@xterm/xterm`**.
   - **Who sees it**: Anyone who installs Vessel on their VPS (`curl -fsSL https://get.vessel.dev | sh`). When they visit `http://their-server-ip:3000` or `https://app.their-domain.com`, this `dashboard/` is the control panel they log into to deploy apps, manage databases, view live logs, and configure `.env` variables.
2. **`website/` (Public Marketing Landing Page)**:
   - **What it is**: The public website hosted globally at `vessel.dev` (built with Astro / Vite / Tailwind).
   - **Who sees it**: Developers around the world discovering the project, reading features, documentation, and copying the one-line install command (`curl -fsSL https://get.vessel.dev | sh`). It is never downloaded or run on a user's VPS.
3. **`cmd/` & `internal/` (Orchestrator Backend & Daemon — Go / Golang)**:
   - **What it is**: High-performance backend orchestrator built in **Go (`vesseld`)**.
   - **What it does**: Talks directly to the Docker socket (`/var/run/docker.sock`), manages Caddy SSL/reverse-proxy rules, executes git webhooks, streams WebSocket terminal logs (`gorilla/websocket`), manages SQLite state (`data/vessel.db`), and handles self-upgrade commands (`scripts/upgrade.sh`).
4. **`get-vessel/` (Installation Host — `get.vessel.dev`)**:
   - **What it is**: Lightweight script delivery server serving `install.sh`, `upgrade.sh`, and system bootstrap scripts.
5. **`scripts/` (System Automation, Upgrade & Bootstrap)**:
   - **What it is**: Core shell automation (`upgrade.sh`, `backup.sh`, `restore.sh`, `bootstrap-host.sh`) allowing the user (and the `dashboard/` GUI via "Check for Updates" / "Upgrade Now" button) to self-update Vessel in place without losing containers or data.

---

## 2. Comprehensive Repository Structure

```
vessel/
├── cmd/
│   ├── vesseld/               # Main entrypoint (`main.go`) for the Go server daemon
│   └── vessel/                # CLI utility (`vessel deploy`, `vessel status`)
│
├── internal/
│   ├── types/                 # Domain structs (`ContainerHealth`, `ProjectConfig`, `SystemInfo`)
│   ├── api/                   # REST & WebSocket endpoints serving `dashboard/`
│   ├── orchestrator/          # Native Docker SDK container & volume management
│   ├── proxy/                 # Caddy v2 dynamic `Caddyfile` generator & hot-reloader (`caddy reload`)
│   ├── store/                 # Embedded SQLite (`CGO_ENABLED=0` sqlite) + AES `.env` vault
│   └── updater/               # Self-upgrade manager (triggers `scripts/upgrade.sh`)
│
├── dashboard/                 # 💻 MAIN PANEL DASHBOARD (What the VPS installer sees & uses)
│   ├── src/
│   │   ├── routes/            # TanStack Router type-safe routes (Dashboard, Projects, Databases)
│   │   ├── components/        # Glassmorphism UI cards, xterm.js live terminal, `.env` vault
│   │   └── hooks/             # WebSocket streaming hooks, TanStack Query mutations
│   ├── index.html             # Main entry point for local GUI
│   └── package.json           # Frontend dependencies (@tanstack/react-router, @xterm/xterm)
│
├── website/                   # 🌐 PUBLIC MARKETING LANDING PAGE (vessel.dev)
│   ├── src/                   # Hero section, one-click copy install banner, interactive demo screenshots
│   └── package.json           # Astro / Vite marketing site dependencies
│
├── get-vessel/                # 📦 INSTALLATION HOST (get.vessel.dev)
│   ├── install.sh             # One-click installation (`curl -fsSL https://get.vessel.dev | sh`)
│   ├── upgrade.sh             # Remote upgrade check script
│   └── server.js              # Lightweight delivery host
│
├── scripts/                   # 🛠️ SYSTEM AUTOMATION, UPGRADE & BOOTSTRAP
│   ├── upgrade.sh             # 🔄 In-place self-upgrade script (zero-downtime for user apps)
│   ├── backup.sh              # Automated SQLite state & Caddy backup (`/data/backups`)
│   ├── restore.sh             # One-click disaster recovery restore
│   └── bootstrap-host.sh      # Linux OS, Docker, and systemd provisioning checks
│
├── data/                      # 💾 PERSISTENT DATA DIRECTORY (Mounted on host)
│   ├── vessel.db              # SQLite primary database storing projects, users, and encrypted envs
│   ├── backups/               # Automated daily backups (`.tar.gz`)
│   └── caddy/                 # Caddy certificates and dynamic `Caddyfile`
│
├── docker-compose.yml         # All-in-one local dev & production container deployment
├── Dockerfile                 # Multi-stage production build uniting `dashboard/` and `vesseld` binary
├── Makefile                   # Build automation (`make build`, `make dev`)
└── README.md                  # Main open-source repository documentation
```

---

## 3. How Self-Updates & Upgrades Work (`scripts/upgrade.sh`)

One of Vessel's killer open-source features is **1-Click Self-Updating**:

1. **Check for Updates**: The `dashboard/` pings `GET /api/system/info` on the orchestrator (`vesseld`). The orchestrator checks GitHub Releases.
2. **One-Click Upgrade**: If a new version (`v1.2.0`) is found, the dashboard displays a banner with an **"Upgrade Now"** button.
3. **Safe In-Place Execution**:
   - When clicked, `vesseld` triggers `scripts/upgrade.sh` in the background.
   - `upgrade.sh` automatically creates a snapshot via `scripts/backup.sh` into `data/backups/vessel-backup-[timestamp].tar.gz`.
   - Pulls the latest Docker image (`docker compose pull vessel`).
   - Gracefully recreates only the `vessel` control plane container (`docker compose up -d --force-recreate vessel`).
   - Running user app containers (`Vite`, `NestJS`, `Postgres`, `Redis`) are **completely untouched and experience ZERO downtime during the Vessel upgrade**.
