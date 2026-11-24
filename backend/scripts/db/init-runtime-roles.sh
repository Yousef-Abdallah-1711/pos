#!/bin/sh
# Runs only when PostgreSQL initializes a new data directory. Existing databases must
# run create_runtime_roles.sql manually after their ownership/backup plan is reviewed.
set -e

(
set --
if [ -n "${POS_ACCOUNTING_PASSWORD:-}" ]; then
  set -- "$@" -v "accounting_pw=$POS_ACCOUNTING_PASSWORD"
fi
if [ -n "${POS_AUTH_LOOKUP_PASSWORD:-}" ]; then
  set -- "$@" -v "auth_lookup_pw=$POS_AUTH_LOOKUP_PASSWORD"
fi

psql "$@" -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname postgres \
  -v dbname="$POSTGRES_DB" \
  -v migrator_pw="$POS_MIGRATOR_PASSWORD" \
  -v app_pw="$POS_APP_PASSWORD" \
  -f /pos/create_runtime_roles.sql
)
