# 🛰️ Vessel

**The Ultra-Lightweight, Self-Hosted PaaS for Developers.**

Turn any bare-metal Linux VPS into your own private Vercel, Railway, or Heroku in 60 seconds with zero-downtime deployments, automated SSL routing, and an ultra-responsive web control panel.

---

## ✨ Features

- **⚡ Blazing-Fast Go Daemon (`vesseld`)**: Uses native Go concurrency and official Docker SDK with `< 30MB RAM` idle overhead.
- **🖥️ Next-Gen Web Dashboard (`ui/`)**: Built with **Vite + TanStack Router + React + Tailwind CSS**. Features live `@xterm/xterm` terminal logs, real-time CPU/RAM stats, and dark-mode glassmorphism.
- **🔒 Automated Edge Routing (`Caddy v2`)**: Zero-config Let's Encrypt SSL/TLS certificates and automatic reverse proxy configuration.
- **🔐 Encrypted `.env` Vault**: AES-256 encrypted environment variables stored inside an embedded SQLite database.
- **🔄 1-Click Zero-Downtime Self-Updates**: Upgrade the Vessel control plane with a single click (`deploy/upgrade.sh`) while your deployed user applications experience **zero seconds of downtime**.

---

## 📂 Repository Layout (Standard Go Project Layout)

```
vessel/
├── cmd/vesseld/          # Go Daemon entrypoint (`main.go`)
├── internal/             # Go packages (`types`, `orchestrator`, `proxy`, `store`, `api`, `updater`)
├── ui/                   # 💻 Self-Hosted Web GUI Dashboard (Vite + TanStack Router)
├── www/                  # 🌐 Public Marketing Landing Page (vessel.dev)
├── deploy/               # 📦 Installation (`install.sh`), Upgrade (`upgrade.sh`), and Bootstrap scripts
└── data/                 # 💾 Persistent data storage (SQLite DB, backups, Caddyfile)
```

---

## 🚀 Quick Install (On any Linux VPS)

```bash
curl -fsSL https://get.vessel.dev | bash -s -- --token YOUR_SECRET_TOKEN
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

3. **Run the Dashboard GUI (`ui/`)**:
   ```bash
   cd ui && npm install && npm run dev
   ```

---

## 📄 License

MIT License. See `LICENSE` for details.
