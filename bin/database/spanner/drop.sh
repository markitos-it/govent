#!/usr/bin/env bash

set -euo pipefail

echo "🗑️  Deleting database 'goevents' from Spanner..."

docker run --rm --network container:go-events-spanner google/cloud-sdk:alpine /bin/sh -c "
  export SPANNER_EMULATOR_HOST=localhost:9010;
  gcloud config set auth/disable_credentials true > /dev/null 2>&1;
  gcloud config set project local-project > /dev/null 2>&1;
  
  gcloud spanner databases delete goevents --instance=local-instance --quiet 2>/dev/null || echo '   ⚠️ Database does not exist or has already been deleted.'
"
echo "✅ Database deleted successfully!"