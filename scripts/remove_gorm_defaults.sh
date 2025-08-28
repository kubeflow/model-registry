#!/bin/bash
set -e

# This script removes redundant `default:NULL` from GORM model files.
# The gorm/gen tool incorrectly adds this tag for nullable columns in PostgreSQL.

# The directory containing the generated GORM models
SCHEMA_DIR="internal/db/schema"

# Check if the directory exists
if [ ! -d "$SCHEMA_DIR" ]; then
  echo "Directory $SCHEMA_DIR not found."
  exit 1
fi

# Use sed to remove `;default:NULL` from all .gen.go files in the schema directory
find "$SCHEMA_DIR" -type f -name "*.gen.go" -exec sed -i 's/;default:NULL//g' {} +

echo "Cleaned up 'default:NULL' from GORM models." 