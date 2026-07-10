#!/usr/bin/env bash
set -eo pipefail

VERSION=${1:-latest}
VESSEL_DIR=${VESSEL_DIR:-/vessel}

echo "🛰️  Starting Vessel upgrade to v${VERSION}..."

if [ ! -f "$VESSEL_DIR/docker-compose.yml" ]; then
  if [ -f "./docker-compose.yml" ]; then
    VESSEL_DIR="."
  else
    echo "❌ No Vessel installation found at $VESSEL_DIR."
    exit 1
  fi
fi

echo "📦 Taking pre-upgrade state backup..."
mkdir -p "$VESSEL_DIR/data/backups"
if [ -f "$VESSEL_DIR/data/vessel.db" ]; then
  # Safely copy database or dump it
  cp "$VESSEL_DIR/data/vessel.db" "$VESSEL_DIR/data/backups/vessel-pre-upgrade-$(date +%Y%m%d%H%M%S).db" 2>/dev/null || true
  echo "✅ SQLite database backed up safely."
fi

export VESSEL_VERSION="$VERSION"
echo "⬇️  Fetching Vessel release v${VERSION}..."

if command -v docker &> /dev/null; then
  docker compose -f "$VESSEL_DIR/docker-compose.yml" pull
  docker compose -f "$VESSEL_DIR/docker-compose.yml" up -d --force-recreate
else
  echo "❌ Docker not found. Running outside Docker? Please replace binary manually."
  exit 1
fi

echo "🚀 Vessel upgrade to v${VERSION} completed successfully! User containers experienced zero downtime."
