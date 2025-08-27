#!/bin/bash
set -e

# This script removes redundant `default:NULL` from GORM model files
# and fixes SQLite-specific GORM generation issues.

# The directory containing the generated GORM models
SCHEMA_DIR="internal/db/schema"

# Check if the directory exists
if [ ! -d "$SCHEMA_DIR" ]; then
  echo "Directory $SCHEMA_DIR not found."
  exit 1
fi

# Remove `;default:NULL` from all .gen.go files in the schema directory
find "$SCHEMA_DIR" -type f -name "*.gen.go" -exec sed -i 's/;default:NULL//g' {} +

# Fix SQLite primary key issues:
# 1. Convert `*int32 primaryKey` to `int32 primaryKey;autoIncrement:true` for id fields  
# 2. Fix spacing to match MySQL/PostgreSQL formatting exactly
# 3. This ensures SQLite generates identical Go struct types and formatting

# Fix ID fields with proper spacing (23 spaces between ID and int32, 3 spaces before backtick)
find "$SCHEMA_DIR" -type f -name "*.gen.go" -exec sed -i -E 's/^([[:space:]]*ID)[[:space:]]+\*int32[[:space:]]+(`gorm:"column:id;primaryKey[^`]*`)/\1                       int32   \2;autoIncrement:true`/g' {} +

# Also fix the closing backtick issue - remove extra backtick if present
find "$SCHEMA_DIR" -type f -name "*.gen.go" -exec sed -i 's/;autoIncrement:true``;/;autoIncrement:true`/g' {} +

echo "Cleaned up GORM models: removed 'default:NULL' and fixed SQLite primary keys." 