#!/bin/bash

# Check database connection
if ! mysql -D $MYSQL_DATABASE -u$MYSQL_USER_NAME -p$MYSQL_ROOT_PASSWORD -e 'SELECT 1' > /dev/null 2>&1; then
    exit 1
fi

# Check if schema_migrations table exists and is not dirty
DIRTY_VERSION=$(mysql -D $MYSQL_DATABASE -u$MYSQL_USER_NAME -p$MYSQL_ROOT_PASSWORD -N -e 'SELECT version FROM schema_migrations WHERE dirty = 1')
if [ -z "$DIRTY_VERSION" ]; then
    exit 0
else
    echo "Schema migration version $DIRTY_VERSION is in a dirty state"
    exit 1
fi
