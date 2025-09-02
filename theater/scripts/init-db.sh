#!/bin/bash

# Theater Booking System Database Initialization Script
# This script sets up the MySQL database for the theater management system

set -e  # Exit on any error

echo "Initializing Theater Booking System Database..."

# Database configuration
DB_HOST=${DB_HOST:-"localhost"}
DB_PORT=${DB_PORT:-"3306"}
DB_USER=${DB_USER:-"user"}
DB_PASSWORD=${DB_PASSWORD:-"password"}
DB_NAME=${DB_NAME:-"theater_booking"}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Get the directory of this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Path to SQL file
SQL_FILE="$SCRIPT_DIR/init-db.sql"

print_status "Waiting for MySQL to be ready..."

# Wait for MySQL to be available
until mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" -e "SELECT 1" > /dev/null 2>&1; do
  print_warning "MySQL is unavailable - sleeping"
  sleep 2
done

print_status "MySQL is ready!"

# Check if SQL file exists
if [ ! -f "$SQL_FILE" ]; then
    print_error "SQL file not found: $SQL_FILE"
    exit 1
fi

print_status "Executing database schema script..."
if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" < "$SQL_FILE"; then
    print_status "Database schema created successfully"
else
    print_error "Failed to execute database schema script"
    exit 1
fi

# Verify tables were created
print_status "Verifying table creation..."
TABLES=$(mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "SHOW TABLES;" -s)

if [[ $TABLES == *"shows"* ]] && [[ $TABLES == *"bookings"* ]]; then
    print_status "All required tables created successfully"
else
    print_warning "Some tables may not have been created properly"
fi

print_status "Database initialization complete! 🎭"
