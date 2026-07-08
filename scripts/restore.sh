#!/usr/bin/env bash
set -eo pipefail

if [ -z "$1" ]; then
  echo "❌ Usage: ./scripts/restore.sh <path-to-backup.tar.gz>"
  exit 1
fi

BACKUP_FILE="$1"
if [ ! -f "${BACKUP_FILE}" ]; then
  echo "❌ Error: Backup file ${BACKUP_FILE} not found!"
  exit 1
fi

echo "⚠️  Restoring Vessel state from ${BACKUP_FILE}..."
tar -xzf "${BACKUP_FILE}" -C .
echo "✅ Restore completed successfully! Please restart vesseld."
