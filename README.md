# 🛰️ Vessl

**The Ultra-Lightweight, Self-Hosted PaaS for Developers.**

Turn any bare-metal Linux VPS into your own private Vercel, Railway, or Heroku in 60 seconds with zero-downtime deployments, automated SSL routing, and an ultra-responsive web control panel.

---

## ✨ Features

- **⚡ Blazing-Fast Go Daemon (`vessld`)**: Uses native Go concurrency and official Docker SDK with `< 30MB RAM` idle overhead.
- **💻 Self-Hosted Dashboard (`dashboard/`)**: Built with **Vite + TanStack Router + React + Tailwind CSS**. Served directly by the Go daemon. Features live `@xterm/xterm` terminal logs, real-time CPU/RAM stats, and dark-mode glassmorphism.
- **🔒 Automated Edge Routing (`Traefik v3`)**: Zero-config Let's Encrypt SSL/TLS certificates and automatic reverse proxy configuration.
- **🔐 Encrypted `.env` Vault**: AES-256 encrypted environment variables stored inside an embedded SQLite database.
- **☁️ Native Cloud Backend**: Easily compile the SaaS backend (`vesslcloud`) to act as a managed multi-tenant UI proxy with Postgres and Stripe billing.

---

## 📂 Repository Layout

```text
vessl/
├── cmd/
│   ├── vessld/           # Self-Hosted Go Daemon entrypoint
│   └── cloud/            # SaaS Cloud Control Plane entrypoint
├── internal/             # Core Go packages
│   ├── cloud/            # Cloud-only features (Billing, Multi-tenant routing, Webhooks)
│   ├── core/             # Agent and Docker engine components
│   ├── http/             # Unified HTTP handlers, routes, and middleware
│   ├── models/           # Domain models (project, database, user)
│   ├── proxy/            # Traefik v3 reverse proxy controller
│   └── services/         # Business logic services (cron, git, notifications)
├── dashboard/            # 💻 React/Vite Dashboard (compiles to static assets)
├── web/                  # 🌐 Marketing site
└── docs/                 # 📖 Documentation
```

---

## 🚀 Quick Install (On any Linux VPS)

```bash
curl -fsSL https://get.vessl.dev | sh
```

Access your dashboard at `http://your-server-ip:3000`.

---

## ⚡ Makefile Commands & Local Development

Vessl includes a comprehensive root-level `Makefile` to streamline local development for both Self-Hosted and Cloud mode.

**Self-Hosted Edition (`vessld`)**

- `make dev`: Launches the backend daemon (`cmd/vessld`) and dashboard dev servers concurrently.
- `make build`: Builds both the dashboard and the `bin/vessld` binary.

**Cloud SaaS Edition (`vesslcloud`)**

- `make dev-cloud`: Launches the Postgres SaaS backend (`cmd/cloud`) and cloud-configured dashboard concurrently.
- `make build-cloud`: Builds the SaaS backend into `bin/vesslcloud`.

### Getting Started Locally

1. **Prerequisites**: Ensure Go 1.23+ and Node.js 20+ are installed.
2. **Environment**: Copy `.env.example` to `.env`.
3. **Start Dev Environment**:

   ```bash
   make dev
   ```

---

## 📄 License

Vessl Source-Available License. You are free to view, use, and modify the code for personal or internal business use. However, redistribution, reselling, or using Vessl to provide a competing commercial managed PaaS is strictly prohibited without explicit written permission. See `LICENSE` for details.
