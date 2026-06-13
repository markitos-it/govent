#!/usr/bin/env bash

set -euo pipefail

echo "🛑 Stopping Cloud Spanner emulator..."
docker compose -f localhost/docker-compose.yaml stop spanner-emulator spanner-init

echo "✅ Spanner stopped successfully!"
