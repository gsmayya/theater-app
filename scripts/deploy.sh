#!/bin/bash

# Theater App Deployment Script
# This script handles the complete deployment process with data persistence

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
COMPOSE_FILE="docker-compose.yaml"
BACKUP_BEFORE_DEPLOY=true

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --no-backup        Skip backup before deployment"
    echo "  --force            Force deployment without confirmation"
    echo "  --help             Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                 # Deploy with backup"
    echo "  $0 --no-backup     # Deploy without backup"
    echo "  $0 --force         # Force deployment"
}

# Parse command line arguments
FORCE_DEPLOY=false
while [[ $# -gt 0 ]]; do
    case $1 in
        --no-backup)
            BACKUP_BEFORE_DEPLOY=false
            shift
            ;;
        --force)
            FORCE_DEPLOY=true
            shift
            ;;
        --help)
            show_usage
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            show_usage
            exit 1
            ;;
    esac
done

echo -e "${BLUE}🎭 Theater App Deployment Script${NC}"
echo -e "${BLUE}================================${NC}"

# Check if Docker is running
if ! docker info >/dev/null 2>&1; then
    echo -e "${RED}Error: Docker is not running${NC}"
    exit 1
fi

# Check if docker-compose is available
if ! command -v docker-compose >/dev/null 2>&1; then
    echo -e "${RED}Error: docker-compose is not installed${NC}"
    exit 1
fi

# Initialize data directories
echo -e "${YELLOW}📁 Initializing data directories...${NC}"
./scripts/init-data-dirs.sh

# Backup existing data if requested
if [ "$BACKUP_BEFORE_DEPLOY" = true ]; then
    echo -e "${YELLOW}💾 Creating backup before deployment...${NC}"
    if ./scripts/backup-db.sh; then
        echo -e "${GREEN}✓ Backup completed successfully${NC}"
    else
        echo -e "${YELLOW}⚠ Backup failed, but continuing with deployment${NC}"
    fi
fi

# Confirm deployment
if [ "$FORCE_DEPLOY" = false ]; then
    echo -e "${YELLOW}⚠ This will deploy the Theater App and may affect existing data${NC}"
    read -p "Are you sure you want to continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${YELLOW}Deployment cancelled${NC}"
        exit 0
    fi
fi

# Stop existing containers
echo -e "${YELLOW}🛑 Stopping existing containers...${NC}"
docker-compose down

# Pull latest images
echo -e "${YELLOW}📥 Pulling latest images...${NC}"
docker-compose pull

# Build and start containers
echo -e "${YELLOW}🔨 Building and starting containers...${NC}"
docker-compose up -d --build

# Wait for services to be healthy
echo -e "${YELLOW}⏳ Waiting for services to be healthy...${NC}"

# Wait for MySQL
echo -e "${YELLOW}  - Waiting for MySQL...${NC}"
timeout=60
while [ $timeout -gt 0 ]; do
    if docker-compose exec -T mysql-theater mysqladmin ping -h localhost -u theater_user -ptheater_password >/dev/null 2>&1; then
        echo -e "${GREEN}  ✓ MySQL is ready${NC}"
        break
    fi
    sleep 2
    timeout=$((timeout - 2))
done

if [ $timeout -le 0 ]; then
    echo -e "${RED}  ✗ MySQL failed to start within timeout${NC}"
    exit 1
fi

# Wait for Redis
echo -e "${YELLOW}  - Waiting for Redis...${NC}"
timeout=30
while [ $timeout -gt 0 ]; do
    if docker-compose exec -T redis-theater redis-cli -a theater_redis_pass ping >/dev/null 2>&1; then
        echo -e "${GREEN}  ✓ Redis is ready${NC}"
        break
    fi
    sleep 2
    timeout=$((timeout - 2))
done

if [ $timeout -le 0 ]; then
    echo -e "${RED}  ✗ Redis failed to start within timeout${NC}"
    exit 1
fi

# Wait for Backend
echo -e "${YELLOW}  - Waiting for Backend...${NC}"
timeout=60
while [ $timeout -gt 0 ]; do
    if curl -f http://localhost:8080/api/v1/health >/dev/null 2>&1; then
        echo -e "${GREEN}  ✓ Backend is ready${NC}"
        break
    fi
    sleep 2
    timeout=$((timeout - 2))
done

if [ $timeout -le 0 ]; then
    echo -e "${RED}  ✗ Backend failed to start within timeout${NC}"
    exit 1
fi

# Wait for Frontend
echo -e "${YELLOW}  - Waiting for Frontend...${NC}"
timeout=30
while [ $timeout -gt 0 ]; do
    if curl -f http://localhost:3000 >/dev/null 2>&1; then
        echo -e "${GREEN}  ✓ Frontend is ready${NC}"
        break
    fi
    sleep 2
    timeout=$((timeout - 2))
done

if [ $timeout -le 0 ]; then
    echo -e "${YELLOW}  ⚠ Frontend may still be starting up${NC}"
fi

# Show deployment status
echo -e "${GREEN}🎉 Deployment completed successfully!${NC}"
echo ""
echo -e "${BLUE}Service URLs:${NC}"
echo -e "  Frontend: ${GREEN}http://localhost:3000${NC}"
echo -e "  Backend:  ${GREEN}http://localhost:8080${NC}"
echo -e "  MySQL:    ${GREEN}localhost:3306${NC}"
echo -e "  Redis:    ${GREEN}localhost:6379${NC}"
echo ""
echo -e "${BLUE}Useful Commands:${NC}"
echo -e "  View logs:     ${YELLOW}docker-compose logs -f${NC}"
echo -e "  Stop services: ${YELLOW}docker-compose down${NC}"
echo -e "  Backup DB:     ${YELLOW}./scripts/backup-db.sh${NC}"
echo -e "  Restore DB:    ${YELLOW}./scripts/restore-db.sh <backup_file>${NC}"
echo ""
echo -e "${GREEN}Data is persisted in the ./data directory${NC}"
