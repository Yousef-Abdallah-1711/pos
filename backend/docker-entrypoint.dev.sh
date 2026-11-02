#!/bin/sh
set -e

echo "Waiting for PostgreSQL..."
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER"; do sleep 2; done

echo "Running database migrations..."
if [ -n "${MIGRATE_DB_USER:-}" ]; then
  DB_USER="$MIGRATE_DB_USER" DB_PASSWORD="${MIGRATE_DB_PASSWORD:-}" go run ./cmd/migrate
  if [ -n "${DB_USER:-}" ] && [ "$DB_USER" != "$MIGRATE_DB_USER" ]; then
    PGPASSWORD="${MIGRATE_DB_PASSWORD:-}" psql -h "$DB_HOST" -p "$DB_PORT" -U "$MIGRATE_DB_USER" -d "$DB_NAME" -v ON_ERROR_STOP=1 \
      -c "REVOKE ALL ON TABLE schema_migrations FROM \"$DB_USER\"" >/dev/null
  fi
else
  go run ./cmd/migrate
fi

echo "Starting API with hot reload..."
exec air --build.cmd "go build -o /tmp/api ./cmd/api" --build.bin "/tmp/api" --build.include_ext "go" --build.exclude_dir "tmp,vendor,uploads,tests" --build.poll true
