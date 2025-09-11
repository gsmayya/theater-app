#!/bin/bash

# Theater App Database Restore Script
# This script restores a MySQL database from a backup file

set -e

# Configuration
DB_CONTAINER="mysql-theater"
DB_NAME="theater_booking"
DB_USER="theater_user"
DB_PASSWORD="theater_password"
BACKUP_DIR="./data/backups"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to show usage
show_usage() {
    echo "Usage: $0 [backup_file]"
    echo ""
    echo "Options:"
    echo "  backup_file    Path to the backup file to restore"
    echo ""
    echo "Examples:"
    echo "  $0 theater_backup_20240115_143022.sql.gz"
    echo "  $0 ./data/backups/theater_backup_20240115_143022.sql.gz"
    echo ""
    echo "Available backups:"
    ls -la $BACKUP_DIR/theater_backup_*.sql.gz 2>/dev/null || echo "No backups found"
}

# Check if backup file is provided
if [ $# -eq 0 ]; then
    show_usage
    exit 1
fi

BACKUP_FILE=$1

# If backup file doesn't exist, try to find it in backup directory
if [ ! -f "$BACKUP_FILE" ]; then
    if [ -f "$BACKUP_DIR/$BACKUP_FILE" ]; then
        BACKUP_FILE="$BACKUP_DIR/$BACKUP_FILE"
    else
        echo -e "${RED}Error: Backup file '$BACKUP_FILE' not found${NC}"
        echo ""
        show_usage
        exit 1
    fi
fi

echo -e "${YELLOW}Starting database restore...${NC}"
echo -e "${BLUE}Backup file: $BACKUP_FILE${NC}"

# Check if container is running
if ! docker ps | grep -q $DB_CONTAINER; then
    echo -e "${RED}Error: MySQL container '$DB_CONTAINER' is not running${NC}"
    echo "Please start the containers with: docker-compose up -d"
    exit 1
fi

# Confirm restore operation
echo -e "${YELLOW}WARNING: This will replace all data in the database!${NC}"
read -p "Are you sure you want to continue? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}Restore cancelled${NC}"
    exit 0
fi

# Check if backup file is compressed
if [[ $BACKUP_FILE == *.gz ]]; then
    echo -e "${YELLOW}Decompressing backup file...${NC}"
    TEMP_FILE=$(mktemp)
    gunzip -c "$BACKUP_FILE" > "$TEMP_FILE"
    BACKUP_FILE="$TEMP_FILE"
    CLEANUP_TEMP=true
else
    CLEANUP_TEMP=false
fi

# Restore the database
echo -e "${YELLOW}Restoring database...${NC}"
docker exec -i $DB_CONTAINER mysql \
    -u $DB_USER \
    -p$DB_PASSWORD < "$BACKUP_FILE"

# Check if restore was successful
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Database restored successfully${NC}"
else
    echo -e "${RED}✗ Database restore failed${NC}"
    if [ "$CLEANUP_TEMP" = true ]; then
        rm -f "$TEMP_FILE"
    fi
    exit 1
fi

# Clean up temporary file if created
if [ "$CLEANUP_TEMP" = true ]; then
    rm -f "$TEMP_FILE"
fi

echo -e "${GREEN}Database restore completed successfully!${NC}"
