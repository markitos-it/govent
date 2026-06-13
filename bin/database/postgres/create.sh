#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'
ENVIRONMENT_FILE="bin/shared/environment.sh"
source $ENVIRONMENT_FILE

function log_info() {
    echo "[INFO] $*"
}
function log_error() {
    echo "[ERROR] $*" >&2
}

setup_environment
show_config "full"

#:[.'.]:>-------------------------------------
show_banner

CONTAINER_NAME="goevents-postgres"
DB_USER="admin"
DB_NEW="goevents"

docker exec -i ${CONTAINER_NAME} psql -U ${DB_USER} -d postgres -c "CREATE DATABASE ${DB_NEW};"
docker exec -i ${CONTAINER_NAME} psql -U ${DB_USER} -d postgres -c "CREATE USER ${DB_NEW} WITH ENCRYPTED PASSWORD '${DB_NEW}';"
docker exec -i ${CONTAINER_NAME} psql -U ${DB_USER} -d ${DB_NEW} -c "GRANT ALL PRIVILEGES ON DATABASE ${DB_NEW} TO ${DB_NEW}; GRANT ALL ON SCHEMA public TO ${DB_NEW};"
#:[.'.]:>-------------------------------------