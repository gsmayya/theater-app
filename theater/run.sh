#!/bin/sh

echo "Running tests"
REDIS_URL=redis:6379 go test -v ./...
echo "All tests complete" 
./theater