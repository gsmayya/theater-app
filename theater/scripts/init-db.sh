#!/bin/bash

# Database initialization script for theater booking system

DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-3306}
DB_USER=${DB_USER:-user}
DB_PASSWORD=${DB_PASSWORD:-password}
DB_NAME=${DB_NAME:-theater_booking}

echo "Waiting for MySQL to be ready..."

# Wait for MySQL to be available
until mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "SELECT 1" > /dev/null 2>&1; do
  echo "MySQL is unavailable - sleeping"
  sleep 2
done

echo "MySQL is ready!"

# Create database and run schema
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" < /app/db/schema.sql

echo "Database schema initialized successfully!"

# Optional: Load sample data
echo "Database initialization complete!"
