#!/bin/sh

MIGRATION_DIR="internal/sql/schema"

export PGPASSWORD=$DB_PASS
COMMON_DB_ARGS="-h $DB_HOST -U $DB_USER -p $DB_PORT"

# Create the database unconditionally
echo "Creating database $DB_NAME..."
CREATEDB_CMD="CREATE DATABASE \"$DB_NAME\";"
psql $COMMON_DB_ARGS -c "$CREATEDB_CMD"

# Apply the migration scripts
DB_ARGS="$COMMON_DB_ARGS -d $DB_NAME"

for migration in "$MIGRATION_DIR"/*.sql; do
  if [ "$migration" != "internal/storage/schema/000-down.sql" ]; then

    echo "Applying migration: $migration"

    if ! psql $DB_ARGS -f "$migration"; then
      echo "Migration failed: $migration"
      exit 1
    fi
  fi;
done

unset PGPASSWORD

printf "All migrations applied successfully!\n"
