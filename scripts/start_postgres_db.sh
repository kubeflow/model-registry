#!/bin/bash
set -e

# Check if postgres container is already running
if [ "$(docker ps -q -f name=model-registry-postgres)" ]; then
    echo "PostgreSQL container is already running"
    exit 0
fi

# Check if postgres container exists but is stopped
if [ "$(docker ps -aq -f status=exited -f name=model-registry-postgres)" ]; then
    echo "Starting existing PostgreSQL container"
    docker start model-registry-postgres
    exit 0
fi

# Create and start new postgres container
echo "Creating and starting new PostgreSQL container"
docker run --name model-registry-postgres \
    -e POSTGRES_USER=postgres \
    -e POSTGRES_PASSWORD=postgres \
    -e POSTGRES_DB=model-registry \
    -p 5432:5432 \
    -d postgres:15

# Wait for postgres to be ready
echo "Waiting for PostgreSQL to be ready..."
until docker exec model-registry-postgres pg_isready -h localhost -p 5432 -U postgres; do
    echo "PostgreSQL is unavailable - sleeping"
    sleep 1
done

echo "PostgreSQL is up and running" 