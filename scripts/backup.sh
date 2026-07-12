#!/usr/bin/env bash
set -eo pipefail

VESSL_DIR=${VESSL_DIR:-/vessl}

if [ ! -d "$VESSL_DIR/data" ]; then
  if [ -d "./data" ]; then
    VESSL_DIR="."
  else
    echo "❌ No Vessl data directory found at $VESSL_DIR/data."
    exit 1
  fi
fi

BACKUP_DIR="${VESSL_DIR}/data/backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/vessl-backup-${TIMESTAMP}.tar.gz"

echo "📦 Starting Vessl automated backup to ${BACKUP_FILE}..."
mkdir -p "${BACKUP_DIR}"

# Archive data folder (SQLite DB, Traefik configs) excluding existing backups
tar --exclude="backups" -czf "${BACKUP_FILE}" -C "${VESSL_DIR}" data
echo "✅ Backup created successfully: ${BACKUP_FILE}"
