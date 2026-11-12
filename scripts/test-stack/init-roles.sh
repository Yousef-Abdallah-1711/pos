#!/bin/sh
# Runs inside the postgres container on first start (docker-entrypoint-initdb.d).
set -e
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname postgres \
  -v dbname="$POSTGRES_DB" \
  -v migrator_pw="$POS_MIGRATOR_PASSWORD" \
  -v app_pw="$POS_APP_PASSWORD" \
  -f /pos/create_runtime_roles.sql
