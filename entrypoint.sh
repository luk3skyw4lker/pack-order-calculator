#!/bin/sh
set -e

# Wait for postgres to be ready using pg_isready
echo "Waiting for postgres..."
timeout=30
counter=0
until nc -z postgres 5432 2>/dev/null || [ $counter -eq $timeout ]; do
  counter=$((counter + 1))
  echo "Attempt $counter/$timeout: Waiting for postgres..."
  sleep 1
done

if [ $counter -eq $timeout ]; then
  echo "Failed to connect to postgres after $timeout seconds"
  exit 1
fi

echo "PostgreSQL is ready"

# Run migrations
echo "Running database migrations..."
cd /app/src/database/migrations
goose postgres "host=postgres port=5432 user=${DATABASE_USER} password=${DATABASE_PASSWORD} dbname=${DATABASE_NAME} sslmode=${DATABASE_SSL_MODE}" up
cd /app

# Start the application
echo "Starting application..."
exec ./orders-api
