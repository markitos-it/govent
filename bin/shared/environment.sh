#!/bin/bash

# Centralized environment configuration

DEFAULT_DATABASE_HOST="localhost"
DEFAULT_DATABASE_DRIVER="postgres"
DEFAULT_DATABASE_USER="goevents"
DEFAULT_DATABASE_PASSWORD="goevents"
DEFAULT_DATABASE_NAME="goevents"
DEFAULT_DATABASE_SSL_MODE="disable"
DEFAULT_POSTGRES_DSN="host=${DEFAULT_DATABASE_HOST} user=${DEFAULT_DATABASE_USER} password=${DEFAULT_DATABASE_PASSWORD} dbname=${DEFAULT_DATABASE_NAME} sslmode=${DEFAULT_DATABASE_SSL_MODE}"
DEFAULT_MARIADB_DSN="${DEFAULT_DATABASE_USER}:${DEFAULT_DATABASE_PASSWORD}@tcp(${DEFAULT_DATABASE_HOST}:3306)/${DEFAULT_DATABASE_NAME}?parseTime=true"
DEFAULT_GRPC_SERVER_ADDRESS=":30000"
DEFAULT_POSTGRES_CONTAINER_NAME="goevents-postgres"
DEFAULT_SPANNER_PROJECT="local-project"
DEFAULT_SPANNER_INSTANCE="local-instance"
DEFAULT_SPANNER_DATABASE="goevents"
DEFAULT_SPANNER_DSN="projects/${DEFAULT_SPANNER_PROJECT}/instances/${DEFAULT_SPANNER_INSTANCE}/databases/${DEFAULT_SPANNER_DATABASE}"

function setup_environment() {
    if [[ "$DEFAULT_DATABASE_DRIVER" == "spanner" ]]; then
        DEFAULT_DATABASE_DSN="${DEFAULT_SPANNER_DSN}"
    elif [[ "$DEFAULT_DATABASE_DRIVER" == "mariadb" ]]; then
        DEFAULT_DATABASE_DSN="${DEFAULT_MARIADB_DSN}"
    else
        DEFAULT_DATABASE_DSN="${DEFAULT_POSTGRES_DSN}"
    fi

    : ${DATABASE_DSN:="${DEFAULT_DATABASE_DSN}"}
    : ${GRPC_SERVER_ADDRESS:="${DEFAULT_GRPC_SERVER_ADDRESS}"}
    : ${POSTGRES_CONTAINER_NAME:="${DEFAULT_POSTGRES_CONTAINER_NAME}"}
    : ${SPANNER_PROJECT:="${DEFAULT_SPANNER_PROJECT}"}
    : ${SPANNER_INSTANCE:="${DEFAULT_SPANNER_INSTANCE}"}
    : ${SPANNER_DATABASE:="${DEFAULT_SPANNER_DATABASE}"}
    : ${SPANNER_DSN:="${DEFAULT_SPANNER_DSN}"}
    : ${DATABASE_DRIVER:="${DEFAULT_DATABASE_DRIVER}"}

    export DATABASE_DRIVER
    export DATABASE_DSN
    export GRPC_SERVER_ADDRESS
    export POSTGRES_CONTAINER_NAME
    export SPANNER_PROJECT
    export SPANNER_INSTANCE
    export SPANNER_DATABASE
    export SPANNER_DSN
    export SPANNER_EMULATOR_HOST="localhost:9010"
}

function show_config() {
    echo "⚙️  Configuration Profile:"
    echo "   🔹 DRIVER   : $DATABASE_DRIVER"
    echo "   🔹 DSN      : $DATABASE_DSN"
    echo "   🔹 GRPC     : $GRPC_SERVER_ADDRESS"
    echo "   🔹 PG_CONT  : $POSTGRES_CONTAINER_NAME"
    echo "   🔹 SP_PROJ  : $SPANNER_PROJECT"
    echo "   🔹 SP_INST  : $SPANNER_INSTANCE"
    echo "   🔹 SP_DB    : $SPANNER_DATABASE"
    echo "   🔹 SP_DSN   : $SPANNER_DSN"
    
    echo "─────────────────────────────────────"
}

function show_banner() {
    echo "============================================="
    echo " __  __  ____  _  __"
    echo "|  \/  |  _ \| |/ /"
    echo "| \  / | | | | ' / "
    echo "| |\/| | | | |  <  "
    echo "| |  | | |_| | . \ "
    echo "|_|  |_|____/|_|\\_\\"
    echo ""
    echo "Creator: Marco Antonio - markitos"
    echo "============================================="
    echo " > (mArKit0sDevSecOpsKit)"
    echo " > Markitos DevSecOps Kulture"
    echo ""
}