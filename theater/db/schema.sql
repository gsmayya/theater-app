-- Theater Booking System Database Schema
-- This schema is optimized with indexes for common query patterns

CREATE DATABASE IF NOT EXISTS theater_booking;
USE theater_booking;

-- Shows table with optimized indexing
CREATE TABLE IF NOT EXISTS shows (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    details TEXT,
    price INT NOT NULL,
    total_tickets INT NOT NULL,
    booked_tickets INT DEFAULT 0,
    location VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Primary indexes for common search patterns
    INDEX idx_location (location),
    INDEX idx_price (price),
    INDEX idx_availability (total_tickets, booked_tickets),
    INDEX idx_name (name),
    
    -- Composite indexes for complex queries
    INDEX idx_location_price (location, price),
    INDEX idx_location_availability (location, total_tickets, booked_tickets),
    INDEX idx_price_availability (price, total_tickets, booked_tickets),
    
    -- Full-text search index for show names and details
    FULLTEXT INDEX ft_search (name, details)
);

-- Tickets table with optimized indexing
CREATE TABLE IF NOT EXISTS tickets (
    id VARCHAR(36) PRIMARY KEY,
    show_id VARCHAR(36) NOT NULL,
    mobile_number VARCHAR(20) NOT NULL,
    number_of_tickets INT NOT NULL,
    customer_name VARCHAR(255) NOT NULL,
    payment_mobile_number VARCHAR(20),
    date_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    unique_id VARCHAR(255) UNIQUE NOT NULL,
    booked_by VARCHAR(255),
    
    -- Foreign key constraint
    FOREIGN KEY (show_id) REFERENCES shows(id) ON DELETE CASCADE,
    
    -- Indexes for common queries
    INDEX idx_show_id (show_id),
    INDEX idx_mobile_number (mobile_number),
    INDEX idx_unique_id (unique_id),
    INDEX idx_date_time (date_time),
    INDEX idx_booked_by (booked_by),
    
    -- Composite indexes
    INDEX idx_show_mobile (show_id, mobile_number),
    INDEX idx_customer_bookings (mobile_number, date_time)
);

-- Create a view for show availability with computed available tickets
CREATE VIEW show_availability AS
SELECT 
    s.id,
    s.name,
    s.details,
    s.price,
    s.total_tickets,
    s.booked_tickets,
    (s.total_tickets - s.booked_tickets) AS available_tickets,
    s.location,
    s.created_at,
    s.updated_at
FROM shows s;

-- Create indexes on computed columns through materialized approach
-- We'll use triggers to maintain a separate availability table for fast queries
CREATE TABLE IF NOT EXISTS show_availability_index (
    show_id VARCHAR(36) PRIMARY KEY,
    available_tickets INT NOT NULL,
    is_available BOOLEAN AS (available_tickets > 0) STORED,
    
    FOREIGN KEY (show_id) REFERENCES shows(id) ON DELETE CASCADE,
    INDEX idx_available_tickets (available_tickets),
    INDEX idx_is_available (is_available)
);

-- Triggers to maintain availability index
DELIMITER $$

CREATE TRIGGER update_availability_on_show_insert
AFTER INSERT ON shows
FOR EACH ROW
BEGIN
    INSERT INTO show_availability_index (show_id, available_tickets)
    VALUES (NEW.id, NEW.total_tickets - NEW.booked_tickets);
END$$

CREATE TRIGGER update_availability_on_show_update
AFTER UPDATE ON shows
FOR EACH ROW
BEGIN
    UPDATE show_availability_index 
    SET available_tickets = NEW.total_tickets - NEW.booked_tickets
    WHERE show_id = NEW.id;
END$$

CREATE TRIGGER update_availability_on_ticket_insert
AFTER INSERT ON tickets
FOR EACH ROW
BEGIN
    UPDATE shows 
    SET booked_tickets = booked_tickets + NEW.number_of_tickets
    WHERE id = NEW.show_id;
END$$

DELIMITER ;

-- Insert some sample data for testing
INSERT INTO shows (id, name, details, price, total_tickets, location) VALUES
(UUID(), 'Hamilton', 'The revolutionary musical about Alexander Hamilton', 150, 500, 'New York'),
(UUID(), 'The Lion King', 'Disney musical featuring the circle of life', 120, 400, 'Los Angeles'),
(UUID(), 'Phantom of the Opera', 'The mysterious phantom haunts the opera house', 100, 300, 'Chicago'),
(UUID(), 'Wicked', 'The untold story of the witches of Oz', 130, 450, 'New York'),
(UUID(), 'Chicago', 'Razzle dazzle musical set in prohibition era', 110, 350, 'Las Vegas');
