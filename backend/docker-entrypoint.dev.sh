#!/bin/sh
set -e

echo "Waiting for PostgreSQL..."
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER"; do sleep 2; done

echo "Running database migrations..."
go run ./cmd/migrate

echo "Starting API with hot reload..."
exec air --build.cmd "go build -o /tmp/api ./cmd/api" --build.bin "/tmp/api" --build.include_ext "go" --build.exclude_dir "tmp,vendor,uploads,tests" --build.poll true
