#!/usr/bin/env bash
# Remote upgrade script served from get.vessel.dev
# Called by the local scripts/upgrade.sh or by users manually.
set -eo pipefail

VESSEL_DIR=${VESSEL_DIR:-/vessel}

echo "⬆️  Vessel Remote Upgrade Check"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ -f "$VESSEL_DIR/docker-compose.yml" ]; then
  echo "🐳 Pulling latest image..."
  docker compose -f "$VESSEL_DIR/docker-compose.yml" pull
  echo "🚀 Recreating container..."
  docker compose -f "$VESSEL_DIR/docker-compose.yml" up -d --force-recreate
  echo "✅ Upgrade complete."
else
  echo "❌ No Vessel installation found at $VESSEL_DIR."
  echo "   Run: curl -fsSL https://get.vessel.dev | sh"
  exit 1
fi
