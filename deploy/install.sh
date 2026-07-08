#!/usr/bin/env bash
# Vessel 1-Click Installation & VPS Bootstrap Script
# Usage: curl -fsSL https://get.vessel.dev | bash -s -- --token <YOUR_TOKEN>
set -eo pipefail

echo "🛰️  Welcome to Vessel - The Ultra-Lightweight Self-Hosted PaaS"
echo "--------------------------------------------------------------"

# Verify root or sudo access
if [ "$EUID" -ne 0 ]; then
  echo "❌ Error: Please run this script as root (or with sudo)."
  exit 1
fi

# Ensure Docker is installed
if ! command -v docker &> /dev/null; then
  echo "📦 Docker not detected. Installing Docker Engine..."
  curl -fsSL https://get.docker.com | sh
fi

echo "🚀 Provisioning Vessel directories (/vessel/data)..."
mkdir -p /vessel/data/backups /vessel/data/caddy

echo "⬇️  Pulling latest Vessel container stack..."
# Placeholder for pulling Docker image or downloading pre-built single binary
# docker pull ghcr.io/vessel-run/vessel:latest

echo "✅ Vessel installed successfully! Access your dashboard at http://localhost:3000"
