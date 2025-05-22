#!/usr/bin/env bash

set -e

echo "Starting MySQL database..."

docker run --rm -d --name mysql-db -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=model-registry -p 3306:3306 mysql:8.3

echo "Waiting for MySQL to be ready..."

# First wait for MySQL to accept connections
while ! docker exec mysql-db mysqladmin ping -h"localhost" --silent; do
    sleep 1
done

# Then wait for MySQL to be fully ready by trying to execute a query
while ! docker exec mysql-db mysql -uroot -proot -e "SELECT 1" >/dev/null 2>&1; do
    sleep 1
done

# Additional wait to ensure MySQL is fully initialized
sleep 5

# Verify the database exists and is accessible
while ! docker exec mysql-db mysql -uroot -proot -e "USE \`model-registry\`; SELECT 1" >/dev/null 2>&1; do
    sleep 2
done

echo "MySQL is ready"
