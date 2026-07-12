#!/usr/bin/env bash
set -eo pipefail

if [ -z "$1" ]; then
  echo "❌ Usage: ./scripts/restore.sh <path-to-backup.tar.gz>"
  exit 1
fi

BACKUP_FILE="$1"
VESSL_DIR=${VESSL_DIR:-/vessl}

if [ ! -f "${BACKUP_FILE}" ]; then
  echo "❌ Error: Backup file ${BACKUP_FILE} not found!"
  exit 1
fi

if [ ! -d "$VESSL_DIR/data" ]; then
  if [ -d "./data" ]; then
    VESSL_DIR="."
  else
    echo "❌ No Vessl data directory found at $VESSL_DIR/data."
    exit 1
  fi
fi

echo "⚠️  Restoring Vessl state from ${BACKUP_FILE} into ${VESSL_DIR}/data..."
# Extract relative to VESSL_DIR
tar -xzf "${BACKUP_FILE}" -C "${VESSL_DIR}"

# Restarting the service
echo "🔄 Restarting Vessl container to apply restored data..."
if command -v docker &> /dev/null && [ -f "$VESSL_DIR/docker-compose.yml" ]; then
  docker compose -f "$VESSL_DIR/docker-compose.yml" restart vessl
fi

echo "✅ Restore completed successfully! Vessl is back online."
