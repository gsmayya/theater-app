#!/bin/bash

# Theater App Data Directory Initialization Script
# This script creates the necessary data directories and sets proper permissions

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Initializing data directories...${NC}"

# Create data directories
DIRS=(
    "./data"
    "./data/mysql"
    "./data/redis"
    "./data/backups"
    "./data/logs"
)

for dir in "${DIRS[@]}"; do
    if [ ! -d "$dir" ]; then
        echo -e "${YELLOW}Creating directory: $dir${NC}"
        mkdir -p "$dir"
    else
        echo -e "${GREEN}Directory already exists: $dir${NC}"
    fi
done

# Set proper permissions
echo -e "${YELLOW}Setting permissions...${NC}"

# MySQL data directory needs specific permissions
if [ -d "./data/mysql" ]; then
    chmod 755 "./data/mysql"
    echo -e "${GREEN}✓ MySQL data directory permissions set${NC}"
fi

# Redis data directory
if [ -d "./data/redis" ]; then
    chmod 755 "./data/redis"
    echo -e "${GREEN}✓ Redis data directory permissions set${NC}"
fi

# Backups directory
if [ -d "./data/backups" ]; then
    chmod 755 "./data/backups"
    echo -e "${GREEN}✓ Backups directory permissions set${NC}"
fi

# Logs directory
if [ -d "./data/logs" ]; then
    chmod 755 "./data/logs"
    echo -e "${GREEN}✓ Logs directory permissions set${NC}"
fi

# Create .gitkeep files to ensure directories are tracked by git
for dir in "${DIRS[@]}"; do
    if [ ! -f "$dir/.gitkeep" ]; then
        touch "$dir/.gitkeep"
        echo -e "${GREEN}✓ Created .gitkeep in $dir${NC}"
    fi
done

echo -e "${GREEN}Data directory initialization completed successfully!${NC}"
echo ""
echo -e "${YELLOW}Directory structure:${NC}"
tree ./data 2>/dev/null || ls -la ./data
