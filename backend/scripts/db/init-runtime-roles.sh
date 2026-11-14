#!/bin/sh
# Runs only when PostgreSQL initializes a new data directory. Existing databases must
# run create_runtime_roles.sql manually after their ownership/backup plan is reviewed.
set -e

(
if [ -n "${POS_ACCOUNTING_PASSWORD:-}" ]; then
  set -- -v "accounting_pw=$POS_ACCOUNTING_PASSWORD"
else
  set --
fi

psql "$@" -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname postgres \
  -v dbname="$POSTGRES_DB" \
  -v migrator_pw="$POS_MIGRATOR_PASSWORD" \
  -v app_pw="$POS_APP_PASSWORD" \
  -f /pos/create_runtime_roles.sql
)
