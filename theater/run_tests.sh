#!/bin/bash

# Theater Booking System Test Runner
# This script runs all tests with proper setup and reporting

set -e

echo "🎭 Theater Booking System - Test Runner"
echo "========================================"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed or not in PATH"
    exit 1
fi

print_status "Go version: $(go version)"

# Set test environment variables
export DB_HOST=localhost
export DB_USER=test_user
export DB_PASSWORD=test_password
export DB_NAME=test_theater_booking
export DB_PORT=3306
export REDIS_URL=localhost:6379
export REDIS_PASSWORD=test_password

print_status "Test environment variables set"

# Run tests with verbose output and coverage
print_status "Running tests with coverage..."

# Create coverage directory
mkdir -p coverage

# Run tests for each package
packages=(
    "./db"
    "./bookings"
    "./shows"
    "./handlers"
    "./utils"
    "./service"
    "./repository"
)

total_tests=0
failed_tests=0

for package in "${packages[@]}"; do
    if [ -d "$package" ]; then
        print_status "Testing package: $package"
        
        # Run tests for this package
        if go test -v -coverprofile="coverage/${package##*/}.out" -covermode=atomic "$package"; then
            print_success "Package $package passed"
        else
            print_error "Package $package failed"
            ((failed_tests++))
        fi
        ((total_tests++))
    else
        print_warning "Package $package not found, skipping"
    fi
done

# Generate overall coverage report
print_status "Generating coverage report..."

# Combine all coverage files
echo "mode: atomic" > coverage/combined.out
for package in "${packages[@]}"; do
    if [ -f "coverage/${package##*/}.out" ]; then
        tail -n +2 "coverage/${package##*/}.out" >> coverage/combined.out
    fi
done

# Generate HTML coverage report
if command -v go &> /dev/null; then
    go tool cover -html=coverage/combined.out -o coverage/coverage.html
    print_success "Coverage report generated: coverage/coverage.html"
fi

# Show coverage summary
print_status "Coverage summary:"
go tool cover -func=coverage/combined.out | tail -1

# Run integration tests if they exist
if [ -d "tests/integration" ]; then
    print_status "Running integration tests..."
    if go test -v ./tests/integration/...; then
        print_success "Integration tests passed"
    else
        print_error "Integration tests failed"
        ((failed_tests++))
    fi
    ((total_tests++))
fi

# Run benchmark tests
print_status "Running benchmark tests..."
if go test -bench=. -benchmem ./...; then
    print_success "Benchmark tests completed"
else
    print_warning "Some benchmark tests failed"
fi

# Summary
echo ""
echo "========================================"
echo "Test Summary"
echo "========================================"
print_status "Total test suites: $total_tests"
if [ $failed_tests -eq 0 ]; then
    print_success "All tests passed! 🎉"
    exit 0
else
    print_error "$failed_tests test suite(s) failed"
    exit 1
fi
