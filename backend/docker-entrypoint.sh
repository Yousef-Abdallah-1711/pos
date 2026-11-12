#!/bin/sh
set -eu

wait_for_database() {
  echo "Waiting for PostgreSQL readiness..."
  until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME"; do
    sleep 2
  done
}

case "${1:-}" in
  migrate)
    : "${MIGRATE_DB_USER:?MIGRATE_DB_USER is required for migration runs}"
    : "${MIGRATE_DB_PASSWORD:?MIGRATE_DB_PASSWORD is required for migration runs}"
    wait_for_database
    DB_USER="$MIGRATE_DB_USER" DB_PASSWORD="$MIGRATE_DB_PASSWORD" ./migrate
    PGPASSWORD="$MIGRATE_DB_PASSWORD" psql \
      -h "$DB_HOST" -p "$DB_PORT" -U "$MIGRATE_DB_USER" -d "$DB_NAME" \
      -v ON_ERROR_STOP=1 -c 'REVOKE ALL ON TABLE schema_migrations FROM pos_app' >/dev/null
    ;;
  api)
    if [ -r /run/secrets/redis_password ]; then
      redis_password="$(cat /run/secrets/redis_password)"
      REDIS_URL="redis://:${redis_password}@redis:6379/0"
      export REDIS_URL
      unset redis_password
    fi
    wait_for_database
    exec ./api
    ;;
  *)
    echo "usage: docker-entrypoint.sh {api|migrate}" >&2
    exit 64
    ;;
esac
