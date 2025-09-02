-- Theater Booking System Database Schema
-- This script creates the necessary tables for the theater management system

-- Create database if it doesn't exist
CREATE DATABASE IF NOT EXISTS theater_booking 
CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

USE theater_booking;

-- Shows table with updated structure
CREATE TABLE IF NOT EXISTS shows (
    id VARCHAR(36) PRIMARY KEY,                    -- UUID as string
    name VARCHAR(255) NOT NULL,                    -- Show name/title
    details TEXT,                                  -- Show description
    price INT NOT NULL,                            -- Ticket price in cents or smallest currency unit
    total_tickets INT NOT NULL,                    -- Total available tickets
    booked_tickets INT DEFAULT 0,                  -- Currently booked tickets
    location VARCHAR(255) NOT NULL,                -- Show location
    show_number VARCHAR(50) NOT NULL UNIQUE,       -- Unique show number (e.g., SH-123456)
    show_date DATETIME NOT NULL,                   -- Date and time of the show
    images JSON,                                   -- Array of CMS image IDs
    videos JSON,                                   -- Array of CMS video IDs
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Indexes for performance
    INDEX idx_shows_location (location),
    INDEX idx_shows_price (price),
    INDEX idx_shows_date (show_date),
    INDEX idx_shows_availability ((total_tickets - booked_tickets)),
    INDEX idx_shows_show_number (show_number),
    
    -- Full-text search index for show names and details
    FULLTEXT(name, details)
) ENGINE=InnoDB;

-- Bookings table
CREATE TABLE IF NOT EXISTS bookings (
    booking_id VARCHAR(20) PRIMARY KEY,            -- Hash-based unique ID (BK-XXXXXXXXX)
    show_id VARCHAR(36) NOT NULL,                  -- Foreign key to shows.id
    contact_type ENUM('mobile', 'email') NOT NULL, -- Type of contact information
    contact_value VARCHAR(255) NOT NULL,           -- Mobile number or email address
    number_of_tickets INT NOT NULL,                -- Number of tickets booked
    customer_name VARCHAR(255),                    -- Optional customer name
    total_amount INT NOT NULL,                     -- Total amount paid/to be paid
    booking_date DATETIME NOT NULL,                -- When the booking was made for
    status ENUM('pending', 'confirmed', 'cancelled') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Foreign key constraint
    FOREIGN KEY (show_id) REFERENCES shows(id) ON DELETE CASCADE,
    
    -- Indexes for performance
    INDEX idx_bookings_show_id (show_id),
    INDEX idx_bookings_contact (contact_type, contact_value),
    INDEX idx_bookings_status (status),
    INDEX idx_bookings_date (booking_date),
    INDEX idx_bookings_created (created_at),
    
    -- Composite indexes for common queries
    INDEX idx_bookings_show_status (show_id, status),
    INDEX idx_bookings_contact_status (contact_type, contact_value, status)
) ENGINE=InnoDB;

-- Show availability index table for optimized queries
CREATE TABLE IF NOT EXISTS show_availability_index (
    show_id VARCHAR(36) PRIMARY KEY,
    is_available BOOLEAN NOT NULL,
    available_tickets INT NOT NULL,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (show_id) REFERENCES shows(id) ON DELETE CASCADE,
    INDEX idx_availability (is_available),
    INDEX idx_available_tickets (available_tickets)
) ENGINE=InnoDB;

-- Triggers to maintain show availability index
DELIMITER //

CREATE TRIGGER IF NOT EXISTS update_show_availability_on_insert
    AFTER INSERT ON shows
    FOR EACH ROW
BEGIN
    INSERT INTO show_availability_index (show_id, is_available, available_tickets)
    VALUES (NEW.id, (NEW.total_tickets > NEW.booked_tickets), (NEW.total_tickets - NEW.booked_tickets))
    ON DUPLICATE KEY UPDATE
        is_available = (NEW.total_tickets > NEW.booked_tickets),
        available_tickets = (NEW.total_tickets - NEW.booked_tickets);
END//

CREATE TRIGGER IF NOT EXISTS update_show_availability_on_update
    AFTER UPDATE ON shows
    FOR EACH ROW
BEGIN
    UPDATE show_availability_index
    SET 
        is_available = (NEW.total_tickets > NEW.booked_tickets),
        available_tickets = (NEW.total_tickets - NEW.booked_tickets)
    WHERE show_id = NEW.id;
END//

CREATE TRIGGER IF NOT EXISTS cleanup_show_availability_on_delete
    AFTER DELETE ON shows
    FOR EACH ROW
BEGIN
    DELETE FROM show_availability_index WHERE show_id = OLD.id;
END//

DELIMITER ;

-- Insert sample data for testing (optional)
INSERT IGNORE INTO shows (
    id, 
    name, 
    details, 
    price, 
    total_tickets, 
    booked_tickets, 
    location, 
    show_number, 
    show_date,
    images,
    videos
) VALUES 
(
    '550e8400-e29b-41d4-a716-446655440000',
    'The Lion King',
    'A spectacular musical adaptation of Disney''s beloved animated film',
    5000,
    200,
    25,
    'Broadway Theater, New York',
    'SH-1001',
    '2024-02-15 19:30:00',
    '["img_001", "img_002", "img_003"]',
    '["vid_001", "vid_002"]'
),
(
    '550e8400-e29b-41d4-a716-446655440001',
    'Hamilton',
    'The revolutionary musical about Alexander Hamilton and the founding fathers',
    7500,
    300,
    45,
    'Richard Rodgers Theatre, New York',
    'SH-1002',
    '2024-02-16 20:00:00',
    '["img_004", "img_005"]',
    '["vid_003"]'
),
(
    '550e8400-e29b-41d4-a716-446655440002',
    'Phantom of the Opera',
    'The classic musical about the mysterious phantom who haunts the Paris Opera House',
    6000,
    250,
    30,
    'Majestic Theatre, New York',
    'SH-1003',
    '2024-02-17 19:00:00',
    '["img_006", "img_007", "img_008", "img_009"]',
    '["vid_004", "vid_005"]'
);

-- Insert sample bookings for testing (optional)
INSERT IGNORE INTO bookings (
    booking_id,
    show_id,
    contact_type,
    contact_value,
    number_of_tickets,
    customer_name,
    total_amount,
    booking_date,
    status
) VALUES
(
    'BK-SAMPLE001',
    '550e8400-e29b-41d4-a716-446655440000',
    'email',
    'john.doe@example.com',
    2,
    'John Doe',
    10000,
    '2024-02-15 19:30:00',
    'confirmed'
),
(
    'BK-SAMPLE002',
    '550e8400-e29b-41d4-a716-446655440001',
    'mobile',
    '+1234567890',
    4,
    'Jane Smith',
    30000,
    '2024-02-16 20:00:00',
    'pending'
),
(
    'BK-SAMPLE003',
    '550e8400-e29b-41d4-a716-446655440002',
    'email',
    'mike.wilson@example.com',
    1,
    'Mike Wilson',
    6000,
    '2024-02-17 19:00:00',
    'confirmed'
);

-- Create views for common queries
CREATE VIEW IF NOT EXISTS show_booking_summary AS
SELECT 
    s.id as show_id,
    s.name as show_name,
    s.show_number,
    s.show_date,
    s.location,
    s.price,
    s.total_tickets,
    s.booked_tickets,
    (s.total_tickets - s.booked_tickets) as available_tickets,
    COUNT(b.booking_id) as total_bookings,
    COALESCE(SUM(CASE WHEN b.status = 'confirmed' THEN b.total_amount ELSE 0 END), 0) as confirmed_revenue,
    COALESCE(SUM(CASE WHEN b.status = 'pending' THEN b.total_amount ELSE 0 END), 0) as pending_revenue,
    COUNT(CASE WHEN b.status = 'confirmed' THEN 1 END) as confirmed_bookings,
    COUNT(CASE WHEN b.status = 'pending' THEN 1 END) as pending_bookings,
    COUNT(CASE WHEN b.status = 'cancelled' THEN 1 END) as cancelled_bookings
FROM shows s
LEFT JOIN bookings b ON s.id = b.show_id
GROUP BY s.id, s.name, s.show_number, s.show_date, s.location, s.price, s.total_tickets, s.booked_tickets;

CREATE VIEW IF NOT EXISTS recent_bookings AS
SELECT 
    b.booking_id,
    b.show_id,
    s.name as show_name,
    s.show_number,
    s.show_date,
    b.contact_type,
    b.contact_value,
    b.number_of_tickets,
    b.customer_name,
    b.total_amount,
    b.status,
    b.created_at,
    b.updated_at
FROM bookings b
JOIN shows s ON b.show_id = s.id
ORDER BY b.created_at DESC
LIMIT 50;

-- Verify tables were created successfully
SHOW TABLES;

-- Display table structures
DESCRIBE shows;
DESCRIBE bookings;
DESCRIBE show_availability_index;
