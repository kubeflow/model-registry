#!/usr/bin/env bash

# SQLite database setup script
# Since SQLite is file-based, this script mainly ensures the database directory exists
# and provides a consistent interface with other database setup scripts

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "${SCRIPT_DIR}")"

# Default SQLite database path (can be overridden by environment variable)
SQLITE_DB_PATH="${SQLITE_DB_PATH:-${PROJECT_DIR}/model-registry.db}"

echo "Starting SQLite database setup..."
echo "Database file: ${SQLITE_DB_PATH}"

# Create the directory for the database file if it doesn't exist
DB_DIR="$(dirname "${SQLITE_DB_PATH}")"
if [ ! -d "${DB_DIR}" ]; then
    echo "Creating database directory: ${DB_DIR}"
    mkdir -p "${DB_DIR}"
fi

# SQLite will create the database file automatically when first accessed
# For consistency with other database scripts, we can touch the file
if [ ! -f "${SQLITE_DB_PATH}" ]; then
    echo "Creating empty database file: ${SQLITE_DB_PATH}"
    touch "${SQLITE_DB_PATH}"
fi

echo "SQLite database setup completed successfully."
echo "Database file: ${SQLITE_DB_PATH}"
echo ""
echo "To connect to the database manually, use:"
echo "  sqlite3 ${SQLITE_DB_PATH}"
echo ""
echo "To run migrations, use:"
echo "  make gen/gorm/sqlite"