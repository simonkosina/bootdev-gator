#!/bin/bash

DBSTRING="host=$DBHOST user=$POSTGRES_USER password=$POSTGRES_PASSWORD dbname=$POSTGRES_DB sslmode=$DBSSL"

# Wait for Postgres to be ready
echo "Waiting for Postgres..."
until pg_isready -h "$DBHOST" -U "$POSTGRES_USER" -d "$POSTGRES_DB" >/dev/null 2>&1; do
  echo "Postgres is not ready yet. Sleeping 2 seconds..."
  sleep 2
done

echo "Postgres is ready. Running migrations..."

goose postgres "$DBSTRING" up
