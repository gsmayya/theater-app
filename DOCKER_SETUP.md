# 🐳 Docker Setup for Theater Booking System

This document provides comprehensive instructions for setting up and running the Theater Booking System using Docker.

## 📋 Prerequisites

- Docker Desktop 4.0+ or Docker Engine 20.10+
- Docker Compose 2.0+
- At least 4GB RAM available for Docker
- At least 2GB free disk space

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        Docker Network                           │
│                      (theater-network)                         │
│                                                                 │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │  theater-website│  │  theater-backend│  │   redis-theater │ │
│  │    (Next.js)    │  │     (Go API)    │  │     (Cache)     │ │
│  │   Port: 3000    │  │   Port: 8080    │  │   Port: 6379    │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
│           │                       │                       │     │
│           └───────────────────────┼───────────────────────┘     │
│                                   │                             │
│                          ┌─────────────────┐                   │
│                          │  mysql-theater  │                   │
│                          │   (Database)    │                   │
│                          │   Port: 3306    │                   │
│                          └─────────────────┘                   │
└─────────────────────────────────────────────────────────────────┘
```

## 🚀 Quick Start

### Option 1: Using the Management Script (Recommended)

```bash
# Make the script executable
chmod +x docker-manage.sh

# Start all services
./docker-manage.sh start

# Check status
./docker-manage.sh status

# View all available commands
./docker-manage.sh help
```

### Option 2: Using Docker Compose Directly

```bash
# Start all services in detached mode
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f
```

## 📊 Service Details

### 🌐 Frontend (theater-website)
- **Container**: `theater-website`
- **Port**: `3000`
- **URL**: http://localhost:3000
- **Technology**: Next.js 15 with React 19
- **Features**: Standalone output, optimized build, health checks

### 🔧 Backend (theater-backend)
- **Container**: `theater-backend`
- **Port**: `8080`
- **URL**: http://localhost:8080
- **Technology**: Go 1.24 with Gin framework
- **Features**: Multi-stage build, security hardening, health checks

### 🗄️ Database (mysql-theater)
- **Container**: `mysql-theater`
- **Port**: `3306`
- **Database**: `theater_booking`
- **User**: `theater_user`
- **Password**: `theater_password`
- **Features**: UTF8MB4, persistent storage, automatic schema initialization

### 🔄 Cache (redis-theater)
- **Container**: `redis-theater`
- **Port**: `6379`
- **Password**: `theater_redis_pass`
- **Features**: AOF persistence, password authentication, optimized configuration

## 💾 Data Persistence

The system uses Docker named volumes for data persistence:

- **`mysql_data`**: MySQL database files
- **`redis_data`**: Redis persistence files (RDB + AOF)

Data persists between container restarts and Docker Compose down/up cycles.

## ⚙️ Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `mysql-theater` | MySQL hostname |
| `DB_PORT` | `3306` | MySQL port |
| `DB_USER` | `theater_user` | MySQL username |
| `DB_PASSWORD` | `theater_password` | MySQL password |
| `DB_NAME` | `theater_booking` | Database name |
| `REDIS_URL` | `redis-theater:6379` | Redis connection string |
| `REDIS_PASSWORD` | `theater_redis_pass` | Redis password |
| `NEXT_PUBLIC_API_URL` | `http://localhost:8080` | Backend API URL for frontend |

### Custom Configuration Files

- **MySQL**: `./docker/mysql/my.cnf`
- **Redis**: `./docker/redis/redis.conf`

## 🔧 Common Operations

### Starting Services

```bash
# Start all services
./docker-manage.sh start

# Or with docker-compose
docker-compose up -d
```

### Viewing Logs

```bash
# All services
./docker-manage.sh logs

# Specific service
./docker-manage.sh logs-backend
./docker-manage.sh logs-frontend
./docker-manage.sh logs-mysql
./docker-manage.sh logs-redis

# With docker-compose
docker-compose logs -f [service-name]
```

### Accessing Containers

```bash
# Backend shell
./docker-manage.sh shell-backend

# MySQL CLI
./docker-manage.sh shell-mysql

# Redis CLI
./docker-manage.sh shell-redis

# Frontend shell
./docker-manage.sh shell-frontend
```

### Health Checks

```bash
# Check all services
./docker-manage.sh health

# Individual health endpoints
curl http://localhost:8080/api/v1/health  # Backend
curl http://localhost:3000                # Frontend
```

### Stopping Services

```bash
# Stop all services
./docker-manage.sh stop

# Or with docker-compose
docker-compose down
```

## 🔄 Development Workflow

### 1. Code Changes

For backend changes:
```bash
# Rebuild backend only
docker-compose build theater-backend
docker-compose up -d theater-backend
```

For frontend changes:
```bash
# Rebuild frontend only
docker-compose build theater-website  
docker-compose up -d theater-website
```

### 2. Database Changes

```bash
# Access MySQL to run manual queries
./docker-manage.sh shell-mysql

# Reset database (WARNING: loses all data)
./docker-manage.sh reset
```

### 3. Complete Rebuild

```bash
# Rebuild everything
./docker-manage.sh rebuild
```

## 🧹 Maintenance

### Cleanup Unused Resources

```bash
# Remove containers and networks (keeps volumes)
./docker-manage.sh cleanup

# Complete reset (removes volumes too)
./docker-manage.sh reset
```

### Backup Data

```bash
# Backup MySQL
docker-compose exec mysql-theater mysqldump -u theater_user -ptheater_password theater_booking > backup.sql

# Backup Redis
docker-compose exec redis-theater redis-cli -a theater_redis_pass --rdb /data/backup.rdb
```

### Restore Data

```bash
# Restore MySQL
cat backup.sql | docker-compose exec -T mysql-theater mysql -u theater_user -ptheater_password theater_booking

# Redis data is automatically restored from persistent volume
```

## 🐛 Troubleshooting

### Common Issues

#### 1. Port Already in Use
```bash
# Check what's using the port
lsof -i :3000  # Frontend
lsof -i :8080  # Backend
lsof -i :3306  # MySQL
lsof -i :6379  # Redis

# Kill the process
kill -9 <PID>
```

#### 2. Database Connection Issues
```bash
# Check MySQL logs
./docker-manage.sh logs-mysql

# Test connection
docker-compose exec mysql-theater mysql -u theater_user -ptheater_password -e "SELECT 1"
```

#### 3. Redis Connection Issues
```bash
# Check Redis logs
./docker-manage.sh logs-redis

# Test connection
docker-compose exec redis-theater redis-cli -a theater_redis_pass ping
```

#### 4. Backend Not Starting
```bash
# Check backend logs
./docker-manage.sh logs-backend

# Rebuild with no cache
docker-compose build --no-cache theater-backend
```

#### 5. Frontend Build Issues
```bash
# Check frontend logs
./docker-manage.sh logs-frontend

# Rebuild with no cache
docker-compose build --no-cache theater-website
```

### Health Check Commands

```bash
# Check container status
docker-compose ps

# Check resource usage
docker stats

# Check networks
docker network ls

# Check volumes
docker volume ls
```

## 🔒 Security Considerations

### Production Deployment

1. **Change default passwords**:
   ```bash
   # Update .env file with secure passwords
   MYSQL_ROOT_PASSWORD=your-secure-root-password
   DB_PASSWORD=your-secure-db-password
   REDIS_PASSWORD=your-secure-redis-password
   ```

2. **Use environment files**:
   ```bash
   cp .env.example .env
   # Edit .env with your values
   ```

3. **Limit exposed ports**:
   ```yaml
   # Remove port mappings for internal services in docker-compose.yaml
   # Only expose frontend (3000) and optionally backend (8080)
   ```

4. **Enable TLS**:
   - Add reverse proxy (nginx)
   - Configure SSL certificates
   - Update CORS settings

### Network Security

- All services communicate through isolated Docker network
- No direct external access to database or cache
- Non-root users in containers
- Security headers configured

## 📈 Performance Tuning

### Resource Limits

Add to docker-compose.yaml:
```yaml
services:
  theater-backend:
    deploy:
      resources:
        limits:
          memory: 512M
          cpus: '0.5'
        reservations:
          memory: 256M
          cpus: '0.25'
```

### Database Optimization

- Increase MySQL buffer pool size for production
- Enable query caching
- Configure slow query log

### Redis Optimization

- Adjust memory limits based on usage
- Configure appropriate eviction policy
- Monitor memory usage

## 📚 Additional Resources

- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [MySQL Docker Hub](https://hub.docker.com/_/mysql)
- [Redis Docker Hub](https://hub.docker.com/_/redis)
- [Next.js Docker Documentation](https://nextjs.org/docs/deployment#docker-image)
- [Go Docker Best Practices](https://docs.docker.com/language/golang/)
