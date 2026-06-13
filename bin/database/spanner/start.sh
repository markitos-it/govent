#!/usr/bin/env bash

set -euo pipefail

echo "🚀 Starting Cloud Spanner emulator..."
docker compose -f localhost/docker-compose.yaml up -d spanner-emulator

echo "⏳ Waiting for Spanner emulator to accept connections..."
docker run --rm --network container:go-events-spanner alpine sh -c "while ! nc -z localhost 9010 2>/dev/null; do sleep 1; done"

echo "✅ Spanner emulator ready on port 9010!"
echo "--------------------------------------------------"
