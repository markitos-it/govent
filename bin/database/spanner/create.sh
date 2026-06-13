#!/usr/bin/env bash

set -euo pipefail

echo "📦 Creating instance and database in Spanner..."

docker run --rm --network container:go-events-spanner google/cloud-sdk:alpine /bin/sh -c "
  export SPANNER_EMULATOR_HOST=localhost:9010;
  gcloud config set auth/disable_credentials true > /dev/null 2>&1;
  gcloud config set project local-project > /dev/null 2>&1;
  
  echo '➡️  Creating instance local-instance...';
  gcloud spanner instances create local-instance --config=emulator-config --description='Local Instance' --nodes=1 2>/dev/null || echo '   ℹ️ Instance already exists, skipping...';
  
  echo '➡️  Creating database goevents...';
  gcloud spanner databases create goevents --instance=local-instance 2>/dev/null || echo '   ℹ️ Database already exists, skipping...';
"
echo "✅ Database 'goevents' is ready and operational!"