# Technology Stack

## Frontend

- **Dashboard**: React 19, TanStack Router, TanStack Query, TanStack Store
- **Admin Panel (cloud)**: React 19, TanStack Router, TanStack Query
- **Marketing Site**: Astro 7
- **Docs**: Astro 7, Starlight
- **Styling**: Tailwind CSS v4, Radix UI, class-variance-authority, tailwind-merge
- **Icons**: Lucide React
- **Terminal**: @xterm/xterm, @xterm/addon-fit

## Backend

- **Daemon**: Go 1.25+
- **HTTP Router**: `net/http` (Go 1.22+ enhanced ServeMux with path parameters)
- **Database**: SQLite (via `modernc.org/sqlite`, CGO-free)
- **Encryption**: AES-256-GCM (custom vault implementation)
- **Proxy**: Caddy v2 (automatic Let's Encrypt SSL reverse proxy)
- **Docker SDK**: `github.com/docker/docker/client`
- **WebSocket**: `gorilla/websocket`
- **Auth**: JWT (via `golang-jwt/jwt/v5`)
- **Billing (cloud)**: Stripe, Paddle, Paystack

## DevOps & Infrastructure

- **CI/CD**: GitHub Actions
- **Container Build**: Docker & Docker Compose
- **Process Supervisor**: Built into Go binary (goroutines, no external supervisor needed)
- **Buildpacks**: Railpack, Nixpacks (auto-detection)

## Languages

- Go 1.25+ (backend daemon)
- TypeScript (dashboard, admin panel)
- Astro (marketing & docs)
- Shell/Bash scripts (install, upgrade, backup, restore)
