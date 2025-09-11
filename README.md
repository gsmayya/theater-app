# 🎭 Theater Booking System

A comprehensive, production-ready theater booking management system built with Go, Next.js, MySQL, and Redis. Features advanced search capabilities, data persistence, comprehensive testing, and Docker containerization.

## ✨ Features

### 🎪 Show Management
- **Complete Show Data**: Title, description, date, location, show number, images, videos
- **Advanced Search**: Full-text search, location-based, price range filtering
- **Real-time Availability**: Automatic ticket availability tracking
- **Caching**: Redis-based caching for optimal performance

### 🎟️ Booking System
- **Hash-based Booking IDs**: Unique, reproducible booking identifiers
- **Multi-contact Support**: Mobile number or email-based bookings
- **Status Management**: Pending, confirmed, cancelled booking states
- **Capacity Validation**: Automatic ticket availability checks
- **Real-time Updates**: Immediate show availability updates

### 📊 Analytics & Reporting
- **Booking Statistics**: Revenue, ticket sales, status breakdowns
- **Show Analytics**: Per-show booking summaries and performance metrics
- **Search Statistics**: Advanced search performance tracking

### 🚀 Performance Features
- **Dual Storage**: MySQL for persistence, Redis for caching
- **Database Indexing**: Optimized queries with strategic indexes
- **Connection Pooling**: Efficient database connection management
- **Async Operations**: Non-blocking cache operations

### 🧪 Testing & Quality
- **Comprehensive Testing**: 56+ unit tests covering all components
- **Code Quality**: Clean, maintainable, and well-documented code
- **Error Handling**: Standardized error responses and validation
- **Graceful Shutdown**: Proper server shutdown with cleanup

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Layer    │    │  Service Layer  │    │ Repository Layer│
│                 │    │                 │    │                 │
│ • Show Handlers │◄──►│ • Show Service  │◄──►│ • Show Repo     │
│ • Booking Handlers    │ • Booking Service    │ • Booking Repo  │
│ • Search Handlers     │ • Validation         │ • Caching       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
                               ┌─────────────────────────────┐
                               │                             │
                               ▼                             ▼
                    ┌─────────────────┐           ┌─────────────────┐
                    │     MySQL       │           │     Redis       │
                    │   (Persistent)  │           │    (Cache)      │
                    └─────────────────┘           └─────────────────┘
```

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.24.6+ (for local development)
- MySQL 8.0+ (for local development)
- Redis 6.0+ (for local development)

### Docker Deployment (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd theater-app
   ```

2. **Initialize data directories**
   ```bash
   ./scripts/init-data-dirs.sh
   ```

3. **Deploy the application**
   ```bash
   ./scripts/deploy.sh
   ```

4. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - MySQL: localhost:3306
   - Redis: localhost:6379

### Local Development

1. **Install dependencies**
   ```bash
   cd theater
   go mod download
   ```

2. **Set environment variables**
   ```bash
   export DB_HOST=localhost
   export DB_PORT=3306
   export DB_USER=your_user
   export DB_PASSWORD=your_password
   export DB_NAME=theater_booking
   export REDIS_URL=localhost:6379
   ```

3. **Initialize database**
   ```bash
   chmod +x scripts/init-db.sh
   ./scripts/init-db.sh
   ```

4. **Run the application**
   ```bash
   go run main.go
   ```

## 💾 Data Persistence

The application is configured with full data persistence:

- **MySQL Data**: Stored in `./data/mysql/`
- **Redis Data**: Stored in `./data/redis/`
- **Backups**: Stored in `./data/backups/`
- **Logs**: Stored in `./data/logs/`

### Backup & Restore

```bash
# Create a backup
./scripts/backup-db.sh

# Restore from backup
./scripts/restore-db.sh theater_backup_20240115_143022.sql.gz
```

## 📋 API Endpoints

### 🎪 Show Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/shows/get?id=<show_id>` | Get show details |
| `POST` | `/api/v1/shows/create` | Create new show |
| `GET` | `/api/v1/search` | Advanced show search |
| `GET` | `/api/v1/shows/by-location?location=<location>` | Shows by location |
| `GET` | `/api/v1/shows/by-price-range?min_price=<min>&max_price=<max>` | Shows by price range |
| `PUT` | `/api/v1/shows/update-availability` | Update availability |

### 🎟️ Booking Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/v1/bookings/create` | Create new booking |
| `GET` | `/api/v1/bookings/get?booking_id=<id>` | Get booking details |
| `PUT` | `/api/v1/bookings/update-status` | Update booking status |
| `POST` | `/api/v1/bookings/confirm` | Confirm booking |
| `POST` | `/api/v1/bookings/cancel` | Cancel booking |
| `GET` | `/api/v1/bookings/by-show?show_id=<id>` | Bookings for show |
| `GET` | `/api/v1/bookings/by-contact` | Bookings by contact |
| `GET` | `/api/v1/bookings/search` | Search bookings |

### 📊 Analytics

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/bookings/stats` | Booking statistics |
| `GET` | `/api/v1/shows/booking-summary?show_id=<id>` | Show booking summary |
| `GET` | `/api/v1/stats` | System statistics |
| `GET` | `/api/v1/health` | Health check |

## 💾 Data Models

### Show Data Structure
```json
{
  "show_id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "The Lion King",
  "details": "A spectacular musical adaptation...",
  "price": 5000,
  "total_tickets": 200,
  "booked_tickets": 25,
  "location": "Broadway Theater, New York",
  "show_number": "SH-1001",
  "show_date": "2024-02-15T19:30:00Z",
  "images": ["img_001", "img_002", "img_003"],
  "videos": ["vid_001", "vid_002"]
}
```

### Booking Data Structure
```json
{
  "booking_id": "BK-A1B2C3D4E5F6G7H8",
  "show_id": "550e8400-e29b-41d4-a716-446655440000",
  "contact_type": "email",
  "contact_value": "customer@example.com",
  "number_of_tickets": 2,
  "customer_name": "John Doe",
  "total_amount": 10000,
  "booking_date": "2024-02-15T19:30:00Z",
  "status": "confirmed",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T11:00:00Z"
}
```

## 🗄️ Database Schema

### Shows Table
```sql
CREATE TABLE shows (
    id VARCHAR(36) PRIMARY KEY,           -- UUID
    name VARCHAR(255) NOT NULL,           -- Show title
    details TEXT,                         -- Description
    price INT NOT NULL,                   -- Price in cents
    total_tickets INT NOT NULL,           -- Total capacity
    booked_tickets INT DEFAULT 0,         -- Currently booked
    location VARCHAR(255) NOT NULL,       -- Venue location
    show_number VARCHAR(50) UNIQUE,       -- Show identifier
    show_date DATETIME NOT NULL,          -- Show date/time
    images JSON,                          -- CMS image IDs
    videos JSON,                          -- CMS video IDs
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### Bookings Table
```sql
CREATE TABLE bookings (
    booking_id VARCHAR(20) PRIMARY KEY,   -- Hash-generated ID
    show_id VARCHAR(36) NOT NULL,         -- Foreign key to shows
    contact_type ENUM('mobile', 'email'), -- Contact method
    contact_value VARCHAR(255) NOT NULL,  -- Phone/email
    number_of_tickets INT NOT NULL,       -- Tickets count
    customer_name VARCHAR(255),           -- Optional name
    total_amount INT NOT NULL,            -- Total cost
    booking_date DATETIME NOT NULL,       -- Booking timestamp
    status ENUM('pending', 'confirmed', 'cancelled') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## 🔧 Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `localhost` | MySQL host |
| `DB_PORT` | `3306` | MySQL port |
| `DB_USER` | `user` | MySQL username |
| `DB_PASSWORD` | `password` | MySQL password |
| `DB_NAME` | `theater_booking` | Database name |
| `REDIS_URL` | `localhost:6379` | Redis connection string |

### Docker Services

#### Frontend (theater-website)
- **Container**: `theater-website`
- **Port**: `3000`
- **Technology**: Next.js 15 with React 19
- **Features**: Standalone output, optimized build, health checks

#### Backend (theater-backend)
- **Container**: `theater-backend`
- **Port**: `8080`
- **Technology**: Go 1.24 with Gin framework
- **Features**: Multi-stage build, security hardening, health checks

#### Database (mysql-theater)
- **Container**: `mysql-theater`
- **Port**: `3306`
- **Database**: `theater_booking`
- **User**: `theater_user`
- **Password**: `theater_password`
- **Features**: UTF8MB4, persistent storage, automatic schema initialization

#### Cache (redis-theater)
- **Container**: `redis-theater`
- **Port**: `6379`
- **Password**: `theater_redis_pass`
- **Features**: AOF persistence, password authentication, optimized configuration

## 🧪 Testing

### Running Tests

The application includes comprehensive unit tests for all major components:

```bash
# Run all tests
go test -v ./...

# Run tests for specific packages
go test -v ./bookings
go test -v ./shows
go test -v ./handlers
go test -v ./utils

# Run tests with coverage
go test -v -cover ./...
```

### Test Coverage

- **Bookings Package**: 15 tests covering all booking operations
- **Shows Package**: 6 tests covering show data management
- **Handlers Package**: 12 tests covering HTTP handlers and responses
- **Utils Package**: 8 tests covering utility functions
- **Database Package**: 3 tests (requires running database)

## 🚀 Performance Optimizations

### 1. MySQL Database Indexing
- **Primary Indexes**: Location, price, availability, show name
- **Composite Indexes**: Location+price, location+availability, price+availability
- **Full-text Search**: MySQL FULLTEXT index on show names and details
- **Availability Triggers**: Automatic maintenance of availability indexes
- **Connection Pooling**: Optimized database connections (25 max, 5-minute lifetime)

### 2. Redis Advanced Indexing
- **Location Sets**: Fast location-based lookups using Redis Sets
- **Price Sorted Sets**: Range queries using Redis ZSets with price as score
- **Availability Sorted Sets**: Availability filtering using sorted sets
- **Text Search**: Simple keyword search using set intersections
- **Combined Search**: Multi-criteria searches using set operations

### 3. Multi-Strategy Search System
- **Cache-First Strategy**: Redis indexes for fast common queries
- **Database Fallback**: Complex queries fall back to MySQL with indexes
- **Intelligent Routing**: Automatic selection of optimal search strategy
- **Pagination Support**: Efficient pagination for large result sets

## 🔄 Docker Management

### Common Operations

```bash
# Start all services
docker-compose up -d

# Stop all services
docker-compose down

# View logs
docker-compose logs -f

# Check status
docker-compose ps

# Rebuild services
docker-compose build --no-cache
```

### Data Management

```bash
# Initialize data directories
./scripts/init-data-dirs.sh

# Create backup
./scripts/backup-db.sh

# Restore from backup
./scripts/restore-db.sh <backup_file>

# Deploy with backup safety
./scripts/deploy.sh
```

### Container Access

```bash
# Backend shell
docker-compose exec theater-backend sh

# MySQL CLI
docker-compose exec mysql-theater mysql -u theater_user -ptheater_password theater_booking

# Redis CLI
docker-compose exec redis-theater redis-cli -a theater_redis_pass
```

## 🛠️ Development

### Project Structure

```
theater-app/
├── theater/                 # Go backend API
│   ├── bookings/           # Booking domain models
│   ├── db/                 # Database connection management
│   ├── handlers/           # HTTP request handlers
│   ├── repository/         # Data access layer
│   ├── service/           # Business logic layer
│   ├── shows/             # Show domain models
│   ├── utils/             # Utility functions and Redis client
│   └── main.go            # Application entry point
├── theater-website/        # Next.js frontend
├── scripts/               # Deployment and maintenance scripts
├── docker/                # Docker configuration files
├── data/                  # Persistent data storage
└── docker-compose.yaml    # Docker Compose configuration
```

### Code Quality Features

- **Comprehensive Testing**: 56+ unit tests with full coverage
- **Error Handling**: Standardized error responses and validation
- **Code Organization**: Clean separation of concerns
- **Documentation**: Well-documented code with inline comments
- **Graceful Shutdown**: Proper server shutdown with cleanup
- **Standardized Responses**: Consistent API response format

## 🔒 Security

### Data Protection
- **File Permissions**: Proper directory permissions
- **Backup Security**: Secure backup storage
- **Access Control**: Limited access to data directories
- **Credential Management**: Secure password handling

### Network Security
- **Internal Networks**: Services communicate over private network
- **Port Exposure**: Only necessary ports exposed
- **Authentication**: Password-protected database access
- **CORS Handling**: Proper cross-origin request handling

## 📈 Production Considerations

### Scalability
- **Resource Limits**: Memory and CPU limits configured
- **Connection Pooling**: Optimized database connections
- **Caching Strategy**: Redis for performance
- **Load Balancing**: Ready for load balancer integration

### Reliability
- **Data Persistence**: Full data retention between deployments
- **Backup Strategy**: Comprehensive backup and restore
- **Health Monitoring**: Service health checks
- **Error Handling**: Graceful error handling and recovery

### Monitoring
- **Health Checks**: Docker health checks for all services
- **Logging**: Centralized logging system
- **Metrics**: Performance and usage metrics
- **Alerting**: Health status monitoring

## 🚨 Troubleshooting

### Common Issues

#### Port Already in Use
```bash
# Check what's using the port
lsof -i :3000  # Frontend
lsof -i :8080  # Backend
lsof -i :3306  # MySQL
lsof -i :6379  # Redis
```

#### Database Connection Issues
```bash
# Check MySQL logs
docker-compose logs mysql-theater

# Test connection
docker-compose exec mysql-theater mysql -u theater_user -ptheater_password -e "SELECT 1"
```

#### Data Not Persisting
```bash
# Check volume mounts
docker inspect mysql-theater | grep -A 10 "Mounts"

# Verify directory permissions
ls -la data/
```

## 📚 Additional Resources

- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [MySQL Docker Hub](https://hub.docker.com/_/mysql)
- [Redis Docker Hub](https://hub.docker.com/_/redis)
- [Next.js Docker Documentation](https://nextjs.org/docs/deployment#docker-image)
- [Go Docker Best Practices](https://docs.docker.com/language/golang/)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

Built with ❤️ using Go, Next.js, MySQL, and Redis 🎭