#!/bin/sh

MIGRATION_DIR=internal/sql/schema

export PGPASSWORD=$DB_PASS
DB_ARGS="-h $DB_HOST -U $DB_USER -d $DB_NAME -p $DB_PORT"


if ! psql $DB_ARGS -f  "$MIGRATION_DIR"/000-down.sql; then
  echo "All tables have been removed from the $DB_NAME database."
else
  echo "ERROR: Some tables could not be removed from the $DB_NAME database."
fi

unset PGPASSWORD
