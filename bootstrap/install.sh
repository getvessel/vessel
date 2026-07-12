#!/usr/bin/env bash
# Vessl 1-Click Installer
# Usage: curl -fsSL https://get.vessl.dev | sh
set -eo pipefail

RELEASE=${VESSL_VERSION:-latest}
VESSL_DIR=/vessl
COMPOSE_URL="https://raw.githubusercontent.com/vesslhq/vessl/main/docker-compose.yml"

echo "🛰️  Vessl — Installing v${RELEASE}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# --- Root check ---
if [ "$EUID" -ne 0 ]; then
  echo "❌ Please run as root (or with sudo)."
  exit 1
fi

# --- Docker ---
if ! command -v docker &> /dev/null; then
  echo "📦 Installing Docker..."
  curl -fsSL https://get.docker.com | sh
  systemctl enable --now docker
fi

if ! docker info &> /dev/null; then
  echo "⏳ Waiting for Docker..."
  sleep 3
fi

# --- Directory setup ---
mkdir -p "$VESSL_DIR"/data/{backups,traefik,builds}

# --- Pull files ---
echo "⬇️  Fetching configuration files..."
curl -fsSL "$COMPOSE_URL" -o "$VESSL_DIR/docker-compose.yml"

if [ ! -f "$VESSL_DIR/.env" ]; then
  echo "🔑 Generating .env file..."
  ENV_URL="https://raw.githubusercontent.com/vesslhq/vessl/main/.env.example"
  curl -fsSL "$ENV_URL" -o "$VESSL_DIR/.env"
  # Generate a random 32-character string for JWT secret
  JWT_SECRET=$(head -c 24 /dev/urandom | base64)
  sed -i "s/VESSL_JWT_SECRET=.*/VESSL_JWT_SECRET=${JWT_SECRET}/" "$VESSL_DIR/.env"
fi

# --- Pull & start ---
echo "🐳 Pulling vessl:v${RELEASE}..."
docker compose -f "$VESSL_DIR/docker-compose.yml" pull
echo "🚀 Starting Vessl..."
docker compose -f "$VESSL_DIR/docker-compose.yml" up -d

# --- systemd service (optional) ---
if command -v systemctl &> /dev/null; then
  cat > /etc/systemd/system/vessl.service <<'SERVICE'
[Unit]
Description=Vessl – Self-hosted PaaS
After=docker.service
Requires=docker.service

[Service]
Restart=always
RestartSec=10
WorkingDirectory=/vessl
ExecStartPre=-/usr/bin/docker compose -f /vessl/docker-compose.yml down
ExecStart=/usr/bin/docker compose -f /vessl/docker-compose.yml up
ExecStop=/usr/bin/docker compose -f /vessl/docker-compose.yml down

[Install]
WantedBy=multi-user.target
SERVICE
  systemctl daemon-reload
  systemctl enable --now vessl.service
fi

echo "✅ Vessl is running! Access the dashboard at http://$(curl -4fsS ifconfig.me 2>/dev/null || echo 'your-server-ip'):3000"
echo "📖 Docs: https://docs.vessl.dev"
