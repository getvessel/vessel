#!/usr/bin/env bash
set -eo pipefail

BACKUP_DIR="./data/backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/vessel-backup-${TIMESTAMP}.tar.gz"

echo "📦 Starting Vessel automated backup to ${BACKUP_FILE}..."
mkdir -p "${BACKUP_DIR}"

if [ -d "./data" ]; then
  # Archive SQLite DB and Caddy configs while excluding previous backup archives
  tar --exclude='./data/backups' -czf "${BACKUP_FILE}" ./data
  echo "✅ Backup created successfully: ${BACKUP_FILE}"
else
  echo "⚠️  No ./data directory found to back up."
fi
