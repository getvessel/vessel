#!/usr/bin/env bash
set -eo pipefail

if [ -z "$1" ]; then
  echo "❌ Usage: ./scripts/restore.sh <path-to-backup.tar.gz>"
  exit 1
fi

BACKUP_FILE="$1"
VESSEL_DIR=${VESSEL_DIR:-/vessel}

if [ ! -f "${BACKUP_FILE}" ]; then
  echo "❌ Error: Backup file ${BACKUP_FILE} not found!"
  exit 1
fi

if [ ! -d "$VESSEL_DIR/data" ]; then
  if [ -d "./data" ]; then
    VESSEL_DIR="."
  else
    echo "❌ No Vessel data directory found at $VESSEL_DIR/data."
    exit 1
  fi
fi

echo "⚠️  Restoring Vessel state from ${BACKUP_FILE} into ${VESSEL_DIR}/data..."
# Extract relative to VESSEL_DIR
tar -xzf "${BACKUP_FILE}" -C "${VESSEL_DIR}"

# Restarting the service
echo "🔄 Restarting Vessel container to apply restored data..."
if command -v docker &> /dev/null && [ -f "$VESSEL_DIR/docker-compose.yml" ]; then
  docker compose -f "$VESSEL_DIR/docker-compose.yml" restart vessel
fi

echo "✅ Restore completed successfully! Vessel is back online."
