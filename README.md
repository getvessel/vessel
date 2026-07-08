# 🛰️ Vessel

**The Ultra-Lightweight, Self-Hosted PaaS for Developers.**

Turn any bare-metal Linux VPS into your own private Vercel, Railway, or Heroku in 60 seconds with zero-downtime deployments, automated SSL routing, and an ultra-responsive web control panel.

---

## ✨ Features

- **⚡ Blazing-Fast Go Daemon (`vesseld`)**: Uses native Go concurrency and official Docker SDK with `< 30MB RAM` idle overhead.
- **💻 Main Panel Dashboard (`dashboard/`)**: Built with **Vite + TanStack Router + React + Tailwind CSS**. Features live `@xterm/xterm` terminal logs, real-time CPU/RAM stats, and dark-mode glassmorphism.
- **🔒 Automated Edge Routing (`Caddy v2`)**: Zero-config Let's Encrypt SSL/TLS certificates and automatic reverse proxy configuration.
- **🔐 Encrypted `.env` Vault**: AES-256 encrypted environment variables stored inside an embedded SQLite database.
- **🔄 1-Click Zero-Downtime Self-Updates**: Upgrade the Vessel control plane with a single click (`scripts/upgrade.sh`) while your deployed user applications experience **zero seconds of downtime**.

---

## 📂 Repository Layout

```
vessel/
├── cmd/vesseld/          # Go Daemon entrypoint (`main.go`)
├── internal/             # Go packages (`types`, `orchestrator`, `proxy`, `store`, `api`, `updater`)
├── dashboard/            # 💻 Main Panel Dashboard (Vite + TanStack Router GUI that VPS installers see)
├── website/              # 🌐 Public Marketing Landing Page (vessel.dev)
├── get-vessel/           # 📦 Installation host (`install.sh`, `upgrade.sh`, `server.js`)
├── scripts/              # 🛠️ System automation (`upgrade.sh`, `backup.sh`, `restore.sh`, `bootstrap-host.sh`)
└── data/                 # 💾 Persistent data storage (SQLite DB, backups, Caddyfile)
```

---

## 🚀 Quick Install (On any Linux VPS)

```bash
curl -fsSL https://get.vessel.dev | sh
```

Access your dashboard at `http://your-server-ip:3000`.

---

## 🛠️ Local Development

1. **Verify Go & Docker**:
   ```bash
   go version
   docker --version
   ```

2. **Run the Backend Daemon (`vesseld`)**:
   ```bash
   go run ./cmd/vesseld
   ```

3. **Run the Dashboard GUI (`dashboard/`)**:
   ```bash
   cd dashboard && npm install && npm run dev
   ```

---

## 📄 License

MIT License. See `LICENSE` for details.
