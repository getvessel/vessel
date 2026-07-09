# Contributing to Vessel 🛰️

Thank you for your interest in contributing to **Vessel**! We are building the most lightweight, developer-friendly, open-source self-hosted PaaS.

---

## 🛠️ Development Environment Setup

### Prerequisites

- **Go**: `v1.22+`
- **Node.js**: `v22.12+`
- **Docker**: `v24+`

### Quick Start

1. **Fork & Clone**:

   ```bash
   git clone https://github.com/solomonolatunji/vessel.git
   cd vessel
   ```

2. **Set up environment**:

   ```bash
   cp .env.example .env
   ```

3. **Run the Go Orchestrator (`cmd/vesseld`)**:

   ```bash
   go run ./cmd/vesseld
   ```

4. **Run the Dashboard (`dashboard/`)**:

   ```bash
   cd dashboard
   npm install
   npm run dev
   ```

5. **Run the Marketing Landing Page (`web/`)**:

   ```bash
   cd web
   npm install
   npm run dev
   ```

---

## 📂 Architecture Overview

- **`cmd/vesseld` & `internal/`**: Written in **Go**. Handles Docker socket execution, Caddy proxy configurations, `.env` encryption, and WebSocket streaming.
- **`dashboard/`**: Written in **TypeScript (TanStack + Vite + Tailwind v4)**. The self-hosted web UI where users deploy and manage their projects.
- **`web/`**: Written in **Astro 7**. The public marketing and documentation landing page (`vessel.dev`).

---

## 🔒 Code Style & Standards

- **Go**: Ensure `go fmt ./...` and `go vet ./...` pass cleanly before submitting PRs.
- **TypeScript**: Ensure `npm run check` or `tsc --noEmit` runs without type errors in `dashboard/` and `web/`.
- **Commits**: Follow [Conventional Commits](https://www.conventionalcommits.org/) (e.g., `feat(daemon): add container CPU usage websocket stream`).

---

## 🐛 Submitting Issues & Pull Requests

- If reporting a bug, please include your Linux kernel version (`uname -r`), Docker version, and `vesseld` logs (`journalctl -u vesseld -e`).
- All pull requests must include clear descriptions of changes and verify that `make test` or `go build ./cmd/vesseld` compiles successfully without errors.
