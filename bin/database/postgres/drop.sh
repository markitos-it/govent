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

log_info "Removing database and associated user"

CONTAINER_NAME="goevents-postgres"
DATABASE_MASTER_USERNAME="admin"
DATABASE_NEW_SERVICE="goevents"

docker exec -i ${CONTAINER_NAME} psql -U ${DATABASE_MASTER_USERNAME} -d postgres -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '${DATABASE_NEW_SERVICE}';"
docker exec -i ${CONTAINER_NAME} psql -U ${DATABASE_MASTER_USERNAME} -d postgres -c "DROP DATABASE IF EXISTS ${DATABASE_NEW_SERVICE};"
docker exec -i ${CONTAINER_NAME} psql -U ${DATABASE_MASTER_USERNAME} -d postgres -c "REASSIGN OWNED BY ${DATABASE_NEW_SERVICE} TO ${DATABASE_MASTER_USERNAME};"
docker exec -i ${CONTAINER_NAME} psql -U ${DATABASE_MASTER_USERNAME} -d postgres -c "DROP OWNED BY ${DATABASE_NEW_SERVICE};"
docker exec -i ${CONTAINER_NAME} psql -U ${DATABASE_MASTER_USERNAME} -d postgres -c "DROP USER IF EXISTS ${DATABASE_NEW_SERVICE};"

log_info "Removal process completed"
#:[.'.]:>-------------------------------------