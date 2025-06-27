#!/bin/bash
set -e

# Stop the postgres container if it's running
if [ "$(docker ps -q -f name=model-registry-postgres)" ]; then
    echo "Stopping PostgreSQL container"
    docker stop model-registry-postgres
fi

# Remove the postgres container if it exists
if [ "$(docker ps -aq -f name=model-registry-postgres)" ]; then
    echo "Removing PostgreSQL container"
    docker rm model-registry-postgres
fi 