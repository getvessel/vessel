#!/usr/bin/env bash
set -eo pipefail

VESSEL_DIR=${VESSEL_DIR:-/vessel}

if [ ! -d "$VESSEL_DIR/data" ]; then
  if [ -d "./data" ]; then
    VESSEL_DIR="."
  else
    echo "❌ No Vessel data directory found at $VESSEL_DIR/data."
    exit 1
  fi
fi

BACKUP_DIR="${VESSEL_DIR}/data/backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/vessel-backup-${TIMESTAMP}.tar.gz"

echo "📦 Starting Vessel automated backup to ${BACKUP_FILE}..."
mkdir -p "${BACKUP_DIR}"

# Archive data folder (SQLite DB, Traefik configs) excluding existing backups
tar --exclude="backups" -czf "${BACKUP_FILE}" -C "${VESSEL_DIR}" data
echo "✅ Backup created successfully: ${BACKUP_FILE}"
