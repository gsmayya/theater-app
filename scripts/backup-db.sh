#!/bin/bash

# Theater App Database Backup Script
# This script creates a backup of the MySQL database

set -e

# Configuration
DB_CONTAINER="mysql-theater"
DB_NAME="theater_booking"
DB_USER="theater_user"
DB_PASSWORD="theater_password"
BACKUP_DIR="./data/backups"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="theater_backup_${TIMESTAMP}.sql"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Starting database backup...${NC}"

# Check if container is running
if ! docker ps | grep -q $DB_CONTAINER; then
    echo -e "${RED}Error: MySQL container '$DB_CONTAINER' is not running${NC}"
    echo "Please start the containers with: docker-compose up -d"
    exit 1
fi

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Create the backup
echo -e "${YELLOW}Creating backup: $BACKUP_FILE${NC}"
docker exec $DB_CONTAINER mysqldump \
    -u $DB_USER \
    -p$DB_PASSWORD \
    --single-transaction \
    --routines \
    --triggers \
    --events \
    --add-drop-database \
    --databases $DB_NAME > "$BACKUP_DIR/$BACKUP_FILE"

# Check if backup was successful
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Backup created successfully: $BACKUP_FILE${NC}"
    
    # Compress the backup
    echo -e "${YELLOW}Compressing backup...${NC}"
    gzip "$BACKUP_DIR/$BACKUP_FILE"
    echo -e "${GREEN}✓ Backup compressed: $BACKUP_FILE.gz${NC}"
    
    # Show backup size
    BACKUP_SIZE=$(du -h "$BACKUP_DIR/$BACKUP_FILE.gz" | cut -f1)
    echo -e "${GREEN}Backup size: $BACKUP_SIZE${NC}"
    
    # Keep only the last 7 days of backups
    echo -e "${YELLOW}Cleaning old backups (keeping last 7 days)...${NC}"
    find $BACKUP_DIR -name "theater_backup_*.sql.gz" -mtime +7 -delete
    echo -e "${GREEN}✓ Old backups cleaned${NC}"
    
else
    echo -e "${RED}✗ Backup failed${NC}"
    exit 1
fi

echo -e "${GREEN}Database backup completed successfully!${NC}"
