#!/usr/bin/env bash
# Vessel 1-Click Installation Script
# Usage: curl -fsSL https://get.vessel.dev | bash -s -- --token <YOUR_TOKEN>
set -eo pipefail

echo "🛰️  Welcome to Vessel - The Ultra-Lightweight Self-Hosted PaaS"
echo "--------------------------------------------------------------"

if [ "$EUID" -ne 0 ]; then
  echo "❌ Error: Please run this script as root (or with sudo)."
  exit 1
fi

if ! command -v docker &> /dev/null; then
  echo "📦 Docker not detected. Installing Docker Engine..."
  curl -fsSL https://get.docker.com | sh
fi

echo "🚀 Provisioning Vessel directories (/vessel/data)..."
mkdir -p /vessel/data/backups /vessel/data/caddy

echo "✅ Vessel directories provisioned! Access dashboard on port 3000 once started."
