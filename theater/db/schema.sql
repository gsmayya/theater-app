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
    show_number VARCHAR(100),
    show_date TIMESTAMP,
    images JSON,
    videos JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Primary indexes for common search patterns
    INDEX idx_location (location),
    INDEX idx_price (price),
    INDEX idx_availability (total_tickets, booked_tickets),
    INDEX idx_name (name),
    INDEX idx_show_date (show_date),

    -- Composite indexes for complex queries
    INDEX idx_location_price (location, price),
    INDEX idx_location_availability (location, total_tickets, booked_tickets),
    INDEX idx_price_availability (price, total_tickets, booked_tickets),
    
    -- Full-text search index for show names and details
    FULLTEXT INDEX ft_search (name, details)
);

-- Bookings table with optimized indexing
CREATE TABLE IF NOT EXISTS bookings (
    booking_id VARCHAR(50) PRIMARY KEY,
    show_id VARCHAR(36) NOT NULL,
    contact_type ENUM('mobile', 'email') NOT NULL,
    contact_value VARCHAR(255) NOT NULL,
    number_of_tickets INT NOT NULL,
    customer_name VARCHAR(255),
    total_amount INT NOT NULL,
    booking_date TIMESTAMP NOT NULL,
    status ENUM('pending', 'confirmed', 'cancelled') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Foreign key constraint
    FOREIGN KEY (show_id) REFERENCES shows(id) ON DELETE CASCADE,
    
    -- Indexes for common queries
    INDEX idx_show_id (show_id),
    INDEX idx_contact (contact_type, contact_value),
    INDEX idx_status (status),
    INDEX idx_booking_date (booking_date),
    INDEX idx_created_at (created_at),
    
    -- Composite indexes
    INDEX idx_show_status (show_id, status),
    INDEX idx_contact_booking (contact_type, contact_value, booking_date)
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
    s.show_number,
    s.show_date,
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

CREATE TRIGGER update_availability_on_booking_insert
AFTER INSERT ON bookings
FOR EACH ROW
BEGIN
    UPDATE shows 
    SET booked_tickets = booked_tickets + NEW.number_of_tickets
    WHERE id = NEW.show_id;
END$$

CREATE TRIGGER update_availability_on_booking_update
AFTER UPDATE ON bookings
FOR EACH ROW
BEGIN
    IF OLD.status != NEW.status THEN
        IF NEW.status = 'cancelled' AND OLD.status != 'cancelled' THEN
            UPDATE shows 
            SET booked_tickets = booked_tickets - NEW.number_of_tickets
            WHERE id = NEW.show_id;
        ELSEIF OLD.status = 'cancelled' AND NEW.status != 'cancelled' THEN
            UPDATE shows 
            SET booked_tickets = booked_tickets + NEW.number_of_tickets
            WHERE id = NEW.show_id;
        END IF;
    END IF;
END$$

DELIMITER ;

-- Insert some sample data for testing
INSERT INTO shows (id, name, details, price, total_tickets, location, show_number, show_date) VALUES
(UUID(), 'Hamilton', 'The revolutionary musical about Alexander Hamilton', 150, 500, 'New York', 'SH-001', DATE_ADD(NOW(), INTERVAL 30 DAY)),
(UUID(), 'The Lion King', 'Disney musical featuring the circle of life', 120, 400, 'Los Angeles', 'SH-002', DATE_ADD(NOW(), INTERVAL 35 DAY)),
(UUID(), 'Phantom of the Opera', 'The mysterious phantom haunts the opera house', 100, 300, 'Chicago', 'SH-003', DATE_ADD(NOW(), INTERVAL 40 DAY)),
(UUID(), 'Wicked', 'The untold story of the witches of Oz', 130, 450, 'New York', 'SH-004', DATE_ADD(NOW(), INTERVAL 45 DAY)),
(UUID(), 'Chicago', 'Razzle dazzle musical set in prohibition era', 110, 350, 'Las Vegas', 'SH-005', DATE_ADD(NOW(), INTERVAL 50 DAY));
