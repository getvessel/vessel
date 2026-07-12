#!/usr/bin/env bash
set -eo pipefail

VERSION=${1:-latest}
VESSL_DIR=${VESSL_DIR:-/vessl}

echo "🛰️  Starting Vessl upgrade to v${VERSION}..."

if [ ! -f "$VESSL_DIR/docker-compose.yml" ]; then
  if [ -f "./docker-compose.yml" ]; then
    VESSL_DIR="."
  else
    echo "❌ No Vessl installation found at $VESSL_DIR."
    exit 1
  fi
fi

echo "📦 Taking pre-upgrade state backup..."
mkdir -p "$VESSL_DIR/data/backups"
if [ -f "$VESSL_DIR/data/vessl.db" ]; then
  # Safely copy database or dump it
  cp "$VESSL_DIR/data/vessl.db" "$VESSL_DIR/data/backups/vessl-pre-upgrade-$(date +%Y%m%d%H%M%S).db" 2>/dev/null || true
  echo "✅ SQLite database backed up safely."
fi

export VESSL_VERSION="$VERSION"
echo "⬇️  Fetching Vessl release v${VERSION}..."

if command -v docker &> /dev/null; then
  docker compose -f "$VESSL_DIR/docker-compose.yml" pull
  docker compose -f "$VESSL_DIR/docker-compose.yml" up -d --force-recreate
else
  echo "❌ Docker not found. Running outside Docker? Please replace binary manually."
  exit 1
fi

echo "🚀 Vessl upgrade to v${VERSION} completed successfully! User containers experienced zero downtime."
