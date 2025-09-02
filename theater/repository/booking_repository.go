package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gsmayya/theater/bookings"
	"github.com/gsmayya/theater/db"
	"github.com/gsmayya/theater/utils"
	"github.com/google/uuid"
)

type BookingRepository struct {
	database    *db.Database
	redisClient *utils.RedisAccess
}

func NewBookingRepository() *BookingRepository {
	return &BookingRepository{
		database:    db.GetDatabase(),
		redisClient: utils.GetStoreAccess(),
	}
}

// CreateBooking inserts a new booking into the database and cache
func (r *BookingRepository) CreateBooking(booking *bookings.Booking) error {
	query := `
		INSERT INTO bookings (booking_id, show_id, contact_type, contact_value, number_of_tickets, 
			customer_name, total_amount, booking_date, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.database.GetDB().Exec(query,
		booking.BookingID,
		booking.ShowID.String(),
		booking.ContactType,
		booking.ContactValue,
		booking.NumberOfTickets,
		booking.CustomerName,
		booking.TotalAmount,
		booking.BookingDate,
		booking.Status,
		booking.CreatedAt,
		booking.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}

	// Cache the booking data
	r.cacheBooking(booking)

	log.Printf("Booking created successfully: %s for show %s", booking.BookingID, booking.ShowID.String())
	return nil
}

// GetBooking retrieves a booking by ID, first checking cache, then database
func (r *BookingRepository) GetBooking(bookingID string) (*bookings.Booking, error) {
	// Try to get from cache first
	if cachedBooking, err := r.getBookingFromCache(bookingID); err == nil {
		return cachedBooking, nil
	}

	// If not in cache, get from database
	query := `
		SELECT booking_id, show_id, contact_type, contact_value, number_of_tickets, 
			customer_name, total_amount, booking_date, status, created_at, updated_at
		FROM bookings 
		WHERE booking_id = ?
	`

	row := r.database.GetDB().QueryRow(query, bookingID)

	booking := &bookings.Booking{}
	var showIDStr string
	
	err := row.Scan(
		&booking.BookingID,
		&showIDStr,
		&booking.ContactType,
		&booking.ContactValue,
		&booking.NumberOfTickets,
		&booking.CustomerName,
		&booking.TotalAmount,
		&booking.BookingDate,
		&booking.Status,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("booking not found: %s", bookingID)
		}
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	// Parse show ID
	showID, err := uuid.Parse(showIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid show ID in database: %w", err)
	}
	booking.ShowID = showID

	// Cache the booking for future requests
	r.cacheBooking(booking)

	return booking, nil
}

// UpdateBooking updates an existing booking
func (r *BookingRepository) UpdateBooking(booking *bookings.Booking) error {
	query := `
		UPDATE bookings 
		SET contact_type = ?, contact_value = ?, number_of_tickets = ?, customer_name = ?, 
			total_amount = ?, booking_date = ?, status = ?, updated_at = ?
		WHERE booking_id = ?
	`

	result, err := r.database.GetDB().Exec(query,
		booking.ContactType,
		booking.ContactValue,
		booking.NumberOfTickets,
		booking.CustomerName,
		booking.TotalAmount,
		booking.BookingDate,
		booking.Status,
		time.Now(),
		booking.BookingID,
	)

	if err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("booking not found: %s", booking.BookingID)
	}

	// Update cache
	booking.UpdatedAt = time.Now()
	r.cacheBooking(booking)

	return nil
}

// UpdateBookingStatus updates only the booking status
func (r *BookingRepository) UpdateBookingStatus(bookingID, status string) error {
	query := `UPDATE bookings SET status = ?, updated_at = ? WHERE booking_id = ?`

	result, err := r.database.GetDB().Exec(query, status, time.Now(), bookingID)
	if err != nil {
		return fmt.Errorf("failed to update booking status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("booking not found: %s", bookingID)
	}

	// Update cache if booking exists in cache
	if cachedBooking, err := r.getBookingFromCache(bookingID); err == nil {
		cachedBooking.UpdateStatus(status)
		r.cacheBooking(cachedBooking)
	}

	return nil
}

// DeleteBooking deletes a booking by ID
func (r *BookingRepository) DeleteBooking(bookingID string) error {
	query := "DELETE FROM bookings WHERE booking_id = ?"

	result, err := r.database.GetDB().Exec(query, bookingID)
	if err != nil {
		return fmt.Errorf("failed to delete booking: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("booking not found: %s", bookingID)
	}

	// Remove from cache
	r.removeCachedBooking(bookingID)

	return nil
}

// GetBookingsByShow retrieves all bookings for a specific show
func (r *BookingRepository) GetBookingsByShow(showID uuid.UUID) ([]*bookings.Booking, error) {
	query := `
		SELECT booking_id, show_id, contact_type, contact_value, number_of_tickets, 
			customer_name, total_amount, booking_date, status, created_at, updated_at
		FROM bookings 
		WHERE show_id = ?
		ORDER BY created_at DESC
	`

	return r.executeBookingQuery(query, showID.String())
}

// GetBookingsByContact retrieves bookings by contact information
func (r *BookingRepository) GetBookingsByContact(contactType, contactValue string) ([]*bookings.Booking, error) {
	query := `
		SELECT booking_id, show_id, contact_type, contact_value, number_of_tickets, 
			customer_name, total_amount, booking_date, status, created_at, updated_at
		FROM bookings 
		WHERE contact_type = ? AND contact_value = ?
		ORDER BY created_at DESC
	`

	return r.executeBookingQuery(query, contactType, contactValue)
}

// GetBookingsWithFilters retrieves bookings with optional filtering and pagination
func (r *BookingRepository) GetBookingsWithFilters(filter *bookings.BookingFilter) ([]*bookings.Booking, int, error) {
	// Build the WHERE clause dynamically based on filters
	whereConditions := []string{}
	args := []interface{}{}

	if filter != nil {
		if filter.ShowID != nil {
			whereConditions = append(whereConditions, "show_id = ?")
			args = append(args, filter.ShowID.String())
		}

		if filter.ContactType != "" {
			whereConditions = append(whereConditions, "contact_type = ?")
			args = append(args, filter.ContactType)
		}

		if filter.Status != "" {
			whereConditions = append(whereConditions, "status = ?")
			args = append(args, filter.Status)
		}

		if filter.DateFrom != nil {
			whereConditions = append(whereConditions, "booking_date >= ?")
			args = append(args, *filter.DateFrom)
		}

		if filter.DateTo != nil {
			whereConditions = append(whereConditions, "booking_date <= ?")
			args = append(args, *filter.DateTo)
		}
	}

	// Build the complete query
	baseQuery := "FROM bookings"
	countQuery := "SELECT COUNT(*) " + baseQuery
	selectQuery := `
		SELECT booking_id, show_id, contact_type, contact_value, number_of_tickets, 
			customer_name, total_amount, booking_date, status, created_at, updated_at ` + baseQuery

	if len(whereConditions) > 0 {
		whereClause := " WHERE " + strings.Join(whereConditions, " AND ")
		countQuery += whereClause
		selectQuery += whereClause
	}

	// Add ordering and pagination
	selectQuery += " ORDER BY created_at DESC"
	if filter != nil {
		if filter.Limit > 0 {
			selectQuery += fmt.Sprintf(" LIMIT %d", filter.Limit)
		}
		if filter.Offset > 0 {
			selectQuery += fmt.Sprintf(" OFFSET %d", filter.Offset)
		}
	}

	// Get total count first
	var totalCount int
	err := r.database.GetDB().QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Get the bookings
	rows, err := r.database.GetDB().Query(selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query bookings: %w", err)
	}
	defer rows.Close()

	var bookingsList []*bookings.Booking
	for rows.Next() {
		booking := &bookings.Booking{}
		var showIDStr string

		err := rows.Scan(
			&booking.BookingID,
			&showIDStr,
			&booking.ContactType,
			&booking.ContactValue,
			&booking.NumberOfTickets,
			&booking.CustomerName,
			&booking.TotalAmount,
			&booking.BookingDate,
			&booking.Status,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan booking: %w", err)
		}

		// Parse show ID
		showID, err := uuid.Parse(showIDStr)
		if err != nil {
			log.Printf("Warning: invalid show ID in database: %s", showIDStr)
			continue
		}
		booking.ShowID = showID

		bookingsList = append(bookingsList, booking)

		// Cache each booking asynchronously
		go r.cacheBooking(booking)
	}

	return bookingsList, totalCount, nil
}

// GetBookingStats retrieves booking statistics
func (r *BookingRepository) GetBookingStats() (*bookings.BookingStats, error) {
	stats := &bookings.BookingStats{
		BookingsByStatus: make(map[string]int32),
		BookingsByShow:   make(map[string]int32),
	}

	// Get total bookings, tickets, and revenue
	totalQuery := `
		SELECT 
			COUNT(*) as total_bookings,
			COALESCE(SUM(number_of_tickets), 0) as total_tickets,
			COALESCE(SUM(total_amount), 0) as total_revenue
		FROM bookings
	`
	
	err := r.database.GetDB().QueryRow(totalQuery).Scan(
		&stats.TotalBookings,
		&stats.TotalTickets,
		&stats.TotalRevenue,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get booking totals: %w", err)
	}

	// Get bookings by status
	statusQuery := `SELECT status, COUNT(*) as count FROM bookings GROUP BY status`
	statusRows, err := r.database.GetDB().Query(statusQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get status stats: %w", err)
	}
	defer statusRows.Close()

	for statusRows.Next() {
		var status string
		var count int32
		if err := statusRows.Scan(&status, &count); err == nil {
			stats.BookingsByStatus[status] = count
		}
	}

	// Get bookings by show (top 10)
	showQuery := `
		SELECT show_id, COUNT(*) as count 
		FROM bookings 
		GROUP BY show_id 
		ORDER BY count DESC 
		LIMIT 10
	`
	showRows, err := r.database.GetDB().Query(showQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get show stats: %w", err)
	}
	defer showRows.Close()

	for showRows.Next() {
		var showID string
		var count int32
		if err := showRows.Scan(&showID, &count); err == nil {
			stats.BookingsByShow[showID] = count
		}
	}

	// Get recent bookings (last 10)
	recentBookings, _, err := r.GetBookingsWithFilters(&bookings.BookingFilter{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		log.Printf("Warning: failed to get recent bookings: %v", err)
	} else {
		stats.RecentBookings = recentBookings
	}

	return stats, nil
}

// Helper methods for caching
func (r *BookingRepository) cacheBooking(booking *bookings.Booking) {
	if jsonData, err := booking.ToJSON(); err == nil {
		cacheKey := fmt.Sprintf("booking:%s", booking.BookingID)
		utils.AddToCache(cacheKey, jsonData, r.redisClient)
	}
}

func (r *BookingRepository) getBookingFromCache(bookingID string) (*bookings.Booking, error) {
	cacheKey := fmt.Sprintf("booking:%s", bookingID)
	jsonData, err := utils.GetFromCache(cacheKey, r.redisClient)
	if err != nil {
		return nil, err
	}

	booking := &bookings.Booking{}
	if err := booking.FromJSON(jsonData); err != nil {
		return nil, err
	}

	return booking, nil
}

func (r *BookingRepository) removeCachedBooking(bookingID string) {
	cacheKey := fmt.Sprintf("booking:%s", bookingID)
	utils.DeleteFromCache(cacheKey, r.redisClient)
}

func (r *BookingRepository) executeBookingQuery(query string, args ...interface{}) ([]*bookings.Booking, error) {
	rows, err := r.database.GetDB().Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var bookingsList []*bookings.Booking
	for rows.Next() {
		booking := &bookings.Booking{}
		var showIDStr string

		err := rows.Scan(
			&booking.BookingID,
			&showIDStr,
			&booking.ContactType,
			&booking.ContactValue,
			&booking.NumberOfTickets,
			&booking.CustomerName,
			&booking.TotalAmount,
			&booking.BookingDate,
			&booking.Status,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}

		// Parse show ID
		showID, err := uuid.Parse(showIDStr)
		if err != nil {
			log.Printf("Warning: invalid show ID in database: %s", showIDStr)
			continue
		}
		booking.ShowID = showID

		bookingsList = append(bookingsList, booking)
	}

	return bookingsList, nil
}

// GetTicketsSoldForShow returns the total number of tickets sold for a specific show
func (r *BookingRepository) GetTicketsSoldForShow(showID uuid.UUID) (int32, error) {
	query := `
		SELECT COALESCE(SUM(number_of_tickets), 0) 
		FROM bookings 
		WHERE show_id = ? AND status IN ('confirmed', 'pending')
	`

	var ticketsSold int32
	err := r.database.GetDB().QueryRow(query, showID.String()).Scan(&ticketsSold)
	if err != nil {
		return 0, fmt.Errorf("failed to get tickets sold: %w", err)
	}

	return ticketsSold, nil
}

// ValidateBookingCapacity checks if a booking can be made without exceeding show capacity
func (r *BookingRepository) ValidateBookingCapacity(showID uuid.UUID, requestedTickets int32, showTotalTickets int32) error {
	ticketsSold, err := r.GetTicketsSoldForShow(showID)
	if err != nil {
		return err
	}

	availableTickets := showTotalTickets - ticketsSold
	if requestedTickets > availableTickets {
		return fmt.Errorf("insufficient tickets available. Requested: %d, Available: %d", requestedTickets, availableTickets)
	}

	return nil
}
