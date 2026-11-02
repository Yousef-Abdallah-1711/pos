#!/bin/sh
set -e

echo "⏳ Waiting for PostgreSQL to be ready..."
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER"; do
  echo "Waiting for database connection..."
  sleep 2
done

echo "✅ PostgreSQL is ready!"
echo "🔄 Running database migrations..."
# Migrations may use separate schema-owner credentials (MIGRATE_DB_USER / MIGRATE_DB_PASSWORD);
# the API itself then runs with the restricted DB_USER (see backend/scripts/db/create_runtime_roles.sql).
if [ -n "${MIGRATE_DB_USER:-}" ]; then
  DB_USER="$MIGRATE_DB_USER" DB_PASSWORD="${MIGRATE_DB_PASSWORD:-}" ./migrate
  # The migration bookkeeping table must stay out of the runtime role's reach.
  if [ -n "${DB_USER:-}" ] && [ "$DB_USER" != "$MIGRATE_DB_USER" ]; then
    PGPASSWORD="${MIGRATE_DB_PASSWORD:-}" psql -h "$DB_HOST" -p "$DB_PORT" -U "$MIGRATE_DB_USER" -d "$DB_NAME" -v ON_ERROR_STOP=1 \
      -c "REVOKE ALL ON TABLE schema_migrations FROM \"$DB_USER\"" >/dev/null
  fi
else
  ./migrate
fi

echo "🚀 Starting API server..."
exec ./api
