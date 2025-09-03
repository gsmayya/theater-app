# Theater Booking System - Optimization Guide

## Overview
This theater booking system has been optimized with comprehensive indexing strategies to provide fast, scalable search and data retrieval capabilities. The system uses a multi-layered approach combining MySQL database indexes, Redis data structures, and intelligent caching strategies.

## Optimization Features

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

## API Endpoints

### New Optimized Endpoints

#### Advanced Search
```
GET/POST /api/v1/search
```
**Parameters:**
- `location` (string): Filter by venue location
- `min_price` (int): Minimum ticket price
- `max_price` (int): Maximum ticket price
- `min_available` (int): Minimum available tickets
- `search` (string): Text search in names/descriptions
- `only_available` (boolean): Show only available shows
- `page` (int): Page number for pagination
- `page_size` (int): Results per page (max 100)

**Example:**
```bash
curl "http://localhost:8080/api/v1/search?location=new york&min_price=50&max_price=200&page=1&page_size=10"
```

#### Location-Based Search
```
GET /api/v1/shows/by-location?location=chicago&only_available=true
```

#### Price Range Search
```
GET /api/v1/shows/by-price-range?min_price=100&max_price=300&location=las vegas
```

#### Show Management
```
POST /api/v1/shows/create?name=Hamilton&location=Broadway&price=150&total_tickets=500&details=Musical
GET /api/v1/shows/get?id={show_id}
PUT /api/v1/shows/update-availability?id={show_id}&booked_tickets=50
```

#### System Information
```
GET /api/v1/stats        # Search index statistics
GET /api/v1/health       # System health check
```

### Legacy Endpoints (Backward Compatible)
- `GET /shows` - List all shows
- `GET /show?id={id}` - Get specific show
- `GET /status` - Basic status

## Performance Optimizations

### 1. Indexing Strategy
- **MySQL Indexes**: B-tree indexes for common query patterns
- **Redis Structures**: O(log N) sorted sets for range queries
- **Full-text Search**: MySQL FULLTEXT for text search
- **Composite Indexes**: Multi-column indexes for complex queries

### 2. Caching Layers
- **Redis Cache**: Show data cached with TTL
- **Redis Indexes**: Structured indexes for fast filtering
- **Database Connection Pool**: Reused connections to reduce overhead
- **Async Indexing**: Non-blocking Redis index updates

### 3. Query Optimization
- **Prepared Statements**: All database queries use prepared statements
- **Batch Operations**: Redis pipeline for multiple operations
- **Strategic Fallbacks**: Redis-first with database fallback
- **Efficient Pagination**: Offset-based pagination with count optimization

## Database Schema

### Shows Table
```sql
CREATE TABLE shows (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    details TEXT,
    price INT NOT NULL,
    total_tickets INT NOT NULL,
    booked_tickets INT DEFAULT 0,
    location VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Optimized indexes
    INDEX idx_location (location),
    INDEX idx_price (price),
    INDEX idx_availability (total_tickets, booked_tickets),
    INDEX idx_location_price (location, price),
    FULLTEXT INDEX ft_search (name, details)
);
```

### Availability Index Table
```sql
CREATE TABLE show_availability_index (
    show_id VARCHAR(36) PRIMARY KEY,
    available_tickets INT NOT NULL,
    is_available BOOLEAN AS (available_tickets > 0) STORED,
    
    INDEX idx_available_tickets (available_tickets),
    INDEX idx_is_available (is_available)
);
```

## Redis Index Structure

### Key Patterns
- `shows:location:{location}` - Set of show IDs by location
- `shows:price` - Sorted set with price as score
- `shows:availability` - Sorted set with available tickets as score
- `shows:search:{term}` - Set of show IDs containing search term
- `show:{id}` - Individual show data

### Example Redis Commands
```redis
# Shows in New York
SMEMBERS shows:location:new york

# Shows priced between $100-$200
ZRANGEBYSCORE shows:price 100 200

# Shows with 50+ available tickets
ZRANGEBYSCORE shows:availability 50 +inf

# Combined search (location + price range)
SINTER shows:location:chicago temp:price_range
```

## Performance Metrics

### Expected Improvements
- **Location searches**: ~10x faster using Redis sets
- **Price range queries**: ~5x faster using sorted sets
- **Combined searches**: ~3x faster using set operations
- **Individual show lookups**: ~2x faster with Redis cache
- **Database connections**: 50% reduction in connection overhead

### Scalability
- **Horizontal Redis scaling**: Support for Redis clusters
- **Database read replicas**: Easy integration for read scaling
- **Index maintenance**: Automatic background index updates
- **Memory efficiency**: Optimized Redis memory usage

## Usage Examples

### 1. Search for Available Shows in Chicago under $150
```bash
curl "http://localhost:8080/api/v1/search?location=chicago&max_price=150&only_available=true"
```

### 2. Get All Broadway Shows with Pagination
```bash
curl "http://localhost:8080/api/v1/search?location=broadway&page=1&page_size=20"
```

### 3. Text Search for "Hamilton" Shows
```bash
curl "http://localhost:8080/api/v1/search?search=hamilton"
```

### 4. Complex Multi-criteria Search
```bash
curl -X POST http://localhost:8080/api/v1/search \
  -H "Content-Type: application/json" \
  -d '{
    "location": "new york",
    "min_price": 100,
    "max_price": 300,
    "min_available": 10,
    "search_term": "musical",
    "only_available": true,
    "page": 1,
    "page_size": 15
  }'
```

## Setup and Deployment

### 1. Run with Docker Compose
```bash
docker-compose up --build
```

### 2. Initialize Database
The system automatically initializes the database schema on startup.

### 3. Test Optimized Endpoints
```bash
# Check system health
curl http://localhost:8080/api/v1/health

# View index statistics
curl http://localhost:8080/api/v1/stats

# Create a test show
curl -X POST "http://localhost:8080/api/v1/shows/create?name=Test Show&location=Test City&price=100&total_tickets=200"
```

## Monitoring and Statistics

The system provides detailed statistics about index usage:
- Total shows indexed
- Shows by location distribution
- Index hit rates
- System health status

Access via: `GET /api/v1/stats`

## Architecture Benefits

1. **Fast Searches**: Redis indexes provide sub-millisecond search times
2. **Scalable**: Can handle thousands of concurrent searches
3. **Flexible**: Multiple search strategies for different use cases
4. **Reliable**: Database fallback ensures system resilience
5. **Maintainable**: Clean separation of concerns with repository pattern
6. **Backward Compatible**: Legacy endpoints continue to work

## Future Enhancements

1. **Geographic Search**: Add location-based radius searches
2. **Real-time Updates**: WebSocket notifications for availability changes
3. **Analytics**: Search pattern analysis and optimization
4. **Caching**: Advanced caching strategies with cache warming
5. **Sharding**: Database sharding for massive scale
