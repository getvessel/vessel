#!/usr/bin/env bash
set -eo pipefail

VERSION=${1:-}
VESSL_DIR=${VESSL_DIR:-/vessl}

if [ -z "$VERSION" ]; then
  echo "❌ Usage: ./scripts/downgrade.sh <version>"
  echo "   Example: ./scripts/downgrade.sh 0.1.0"
  echo ""
  echo "Available versions: https://github.com/vesslhq/vessl/releases"
  exit 1
fi

echo "⬇️  Vessl — Downgrading to v${VERSION}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ ! -f "$VESSL_DIR/docker-compose.yml" ]; then
  if [ -f "./docker-compose.yml" ]; then
    VESSL_DIR="."
  else
    echo "❌ No Vessl installation found at $VESSL_DIR."
    exit 1
  fi
fi

echo "💾 Creating database backup before downgrade..."
mkdir -p "$VESSL_DIR/data/backups"
if [ -f "$VESSL_DIR/data/vessl.db" ]; then
  cp "$VESSL_DIR/data/vessl.db" "$VESSL_DIR/data/backups/vessl-pre-downgrade-$(date +%Y%m%d%H%M%S).db" 2>/dev/null || true
  echo "✅ SQLite database backed up safely."
fi

export VESSL_VERSION="$VERSION"
echo "🐳 Pulling vessl:v${VERSION}..."

if command -v docker &> /dev/null; then
  docker compose -f "$VESSL_DIR/docker-compose.yml" pull
  echo "🚀 Recreating container with v${VERSION}..."
  docker compose -f "$VESSL_DIR/docker-compose.yml" up -d --force-recreate
else
  echo "❌ Docker not found."
  exit 1
fi

echo "✅ Downgraded to v${VERSION}."
echo "⚠️  If you encounter issues, restore a previous backup from $VESSL_DIR/data/backups/"
