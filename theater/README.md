# 🎭 Theater Booking Management System

A comprehensive Go-based theater booking management system with MySQL and Redis integration, featuring show management, booking operations, and advanced search capabilities.

## ✨ Features

### 🎪 Show Management
- **Complete Show Data**: Title, description, date, location, show number, images, videos
- **CMS Integration**: Support for CMS-hosted images and videos with ID references
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

## 🚀 Getting Started

### Prerequisites
- Go 1.24.6+
- MySQL 8.0+
- Redis 6.0+

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd theater-app/theater
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set environment variables**
   ```bash
   export DB_HOST=localhost
   export DB_PORT=3306
   export DB_USER=your_user
   export DB_PASSWORD=your_password
   export DB_NAME=theater_booking
   export REDIS_URL=localhost:6379
   ```

4. **Initialize database**
   ```bash
   chmod +x scripts/init-db.sh
   ./scripts/init-db.sh
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

The server will start at `http://localhost:8080` 🎉

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

### Database Optimization

The system includes several optimization features:

1. **Strategic Indexing**: Location, price, date, availability indexes
2. **Full-text Search**: On show names and descriptions
3. **Availability Index**: Separate table for fast availability queries
4. **Connection Pooling**: Configurable connection limits
5. **Database Triggers**: Automatic availability index maintenance

## 🧪 Testing

### Create a Test Show
```bash
curl -X POST "http://localhost:8080/api/v1/shows/create" \
  -d "name=Test Show" \
  -d "details=A test show" \
  -d "price=2500" \
  -d "total_tickets=100" \
  -d "location=Test Theater" \
  -d "show_date=2024-12-25T19:30:00Z"
```

### Create a Test Booking
```bash
curl -X POST "http://localhost:8080/api/v1/bookings/create" \
  -H "Content-Type: application/json" \
  -d '{
    "show_id": "<show_id_from_above>",
    "contact_type": "email",
    "contact_value": "test@example.com",
    "number_of_tickets": 2,
    "customer_name": "Test Customer"
  }'
```

## 🎯 Key Features Implemented

✅ **Complete Show Management**: Title, description, date, location, show number, images, videos  
✅ **Hash-based Booking IDs**: Unique IDs generated from booking information  
✅ **Multi-contact Support**: Mobile and email-based bookings  
✅ **Real-time Availability**: Automatic capacity validation and updates  
✅ **Advanced Caching**: Redis integration for optimal performance  
✅ **Comprehensive API**: RESTful endpoints for all operations  
✅ **Database Optimization**: Indexes, triggers, and connection pooling  
✅ **Error Handling**: Comprehensive validation and error responses  
✅ **Analytics**: Statistics and reporting capabilities  

## 🚀 Production Considerations

1. **Security**: Add authentication, authorization, and input sanitization
2. **Monitoring**: Implement logging, metrics, and health checks  
3. **Scaling**: Consider horizontal scaling and load balancing
4. **Backup**: Implement database backup and recovery procedures
5. **Rate Limiting**: Add API rate limiting for production use
6. **SSL/TLS**: Enable HTTPS for secure communication

## 📁 Project Structure

```
theater/
├── bookings/           # Booking domain models
├── db/                 # Database connection management
├── handlers/           # HTTP request handlers
├── repository/         # Data access layer
├── scripts/           # Database initialization scripts
├── service/           # Business logic layer
├── shows/             # Show domain models
├── utils/             # Utility functions and Redis client
├── main.go            # Application entry point
├── go.mod             # Go module definition
└── README.md          # This file
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

Built with ❤️ using Go, MySQL, and Redis 🎭
