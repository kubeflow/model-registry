#!/usr/bin/env bash

set -e

CONTAINER_NAME="mysql-db"

container_exists() {
    docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"
}

stop_container() {
    if container_exists; then
        echo "Stopping MySQL database..."
        docker stop ${CONTAINER_NAME} || true
    fi
}

remove_container() {
    if container_exists; then
        echo "Removing MySQL database..."
        docker rm ${CONTAINER_NAME} || true
    fi
}

if container_exists; then
    stop_container
    remove_container
    echo "MySQL database removed"
else
    echo "MySQL database container not found"
fi
