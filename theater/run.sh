#!/bin/sh

echo "Starting optimized theater booking service..."

# Initialize database schema
echo "Initializing database..."
./scripts/init-db.sh

echo "Running tests"
REDIS_URL=redis:6379 go test -v ./...
echo "All tests complete"

echo "Starting theater booking service with indexing..."
./theater
