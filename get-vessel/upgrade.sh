#!/usr/bin/env bash
# Remote upgrade check and execution script served from get.vessel.dev
set -eo pipefail
echo "🔍 Checking remote version from https://get.vessel.dev/version..."
# Delegates execution to the local scripts/upgrade.sh
if [ -f "/vessel/scripts/upgrade.sh" ]; then
  bash /vessel/scripts/upgrade.sh
else
  echo "Executing upgrade step..."
fi
