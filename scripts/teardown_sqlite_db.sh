#!/usr/bin/env bash

# SQLite database teardown script
# Removes the SQLite database file and any related temporary files

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "${SCRIPT_DIR}")"

# Default SQLite database path (can be overridden by environment variable)
SQLITE_DB_PATH="${SQLITE_DB_PATH:-${PROJECT_DIR}/model-registry.db}"

echo "Tearing down SQLite database..."
echo "Database file: ${SQLITE_DB_PATH}"

# Remove the main database file if it exists
if [ -f "${SQLITE_DB_PATH}" ]; then
    echo "Removing database file: ${SQLITE_DB_PATH}"
    rm -f "${SQLITE_DB_PATH}"
else
    echo "Database file does not exist: ${SQLITE_DB_PATH}"
fi

# Remove any SQLite temporary files (WAL, SHM files)
SQLITE_WAL_PATH="${SQLITE_DB_PATH}-wal"
SQLITE_SHM_PATH="${SQLITE_DB_PATH}-shm"

if [ -f "${SQLITE_WAL_PATH}" ]; then
    echo "Removing WAL file: ${SQLITE_WAL_PATH}"
    rm -f "${SQLITE_WAL_PATH}"
fi

if [ -f "${SQLITE_SHM_PATH}" ]; then
    echo "Removing SHM file: ${SQLITE_SHM_PATH}"
    rm -f "${SQLITE_SHM_PATH}"
fi

echo "SQLite database teardown completed successfully."