package service

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gsmayya/theater/bookings"
	"github.com/gsmayya/theater/repository"
)

// BookingService provides business logic for theater bookings
type BookingService struct {
	bookingRepository *repository.BookingRepository
	showService       *ShowService
}

// NewBookingService creates a new booking service
func NewBookingService() *BookingService {
	return &BookingService{
		bookingRepository: repository.NewBookingRepository(),
		showService:       NewShowService(),
	}
}

// CreateBooking creates a new booking with validation and availability checks
func (s *BookingService) CreateBooking(showID uuid.UUID, contactType, contactValue string, numberOfTickets int32, customerName string) (*bookings.Booking, error) {
	// Validate input parameters
	if numberOfTickets <= 0 {
		return nil, fmt.Errorf("number of tickets must be greater than 0")
	}

	// Get show details to validate and calculate price
	show, err := s.showService.GetShow(showID.String())
	if err != nil {
		return nil, fmt.Errorf("show not found: %w", err)
	}

	// Validate ticket availability
	err = s.bookingRepository.ValidateBookingCapacity(showID, numberOfTickets, show.Total_Tickets)
	if err != nil {
		return nil, fmt.Errorf("booking validation failed: %w", err)
	}

	// Calculate total amount
	totalAmount := show.Price * numberOfTickets

	// Create the booking
	booking := bookings.NewBooking(showID, contactType, contactValue, numberOfTickets, totalAmount)
	if customerName != "" {
		booking.CustomerName = customerName
	}

	// Save to repository
	if err := s.bookingRepository.CreateBooking(booking); err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	// Update show's booked tickets count
	show.Booked_Tickets += numberOfTickets
	if err := s.showService.UpdateShow(show); err != nil {
		log.Printf("Warning: Failed to update show booked tickets count: %v", err)
		// Don't fail the booking creation for this, but log it
	}

	log.Printf("Successfully created booking: %s for show %s", booking.BookingID, showID.String())
	return booking, nil
}

// GetBooking retrieves a booking by its ID
func (s *BookingService) GetBooking(bookingID string) (*bookings.Booking, error) {
	if bookingID == "" {
		return nil, fmt.Errorf("booking ID cannot be empty")
	}

	booking, err := s.bookingRepository.GetBooking(bookingID)
	if err != nil {
		return nil, err
	}

	return booking, nil
}

// UpdateBookingStatus updates the status of a booking
func (s *BookingService) UpdateBookingStatus(bookingID, status string) error {
	if bookingID == "" {
		return fmt.Errorf("booking ID cannot be empty")
	}

	// Validate status
	validStatuses := []string{"pending", "confirmed", "cancelled"}
	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("invalid status: %s. Valid statuses are: pending, confirmed, cancelled", status)
	}

	// Get current booking to check current status
	currentBooking, err := s.bookingRepository.GetBooking(bookingID)
	if err != nil {
		return fmt.Errorf("failed to get current booking: %w", err)
	}

	// If cancelling a previously confirmed/pending booking, adjust show availability
	if status == "cancelled" && (currentBooking.Status == "confirmed" || currentBooking.Status == "pending") {
		show, err := s.showService.GetShow(currentBooking.ShowID.String())
		if err != nil {
			log.Printf("Warning: Failed to get show for booking cancellation: %v", err)
		} else {
			// Reduce booked tickets count
			show.Booked_Tickets -= currentBooking.NumberOfTickets
			if show.Booked_Tickets < 0 {
				show.Booked_Tickets = 0 // Prevent negative values
			}
			if err := s.showService.UpdateShow(show); err != nil {
				log.Printf("Warning: Failed to update show availability after cancellation: %v", err)
			}
		}
	}

	// Update booking status
	if err := s.bookingRepository.UpdateBookingStatus(bookingID, status); err != nil {
		return fmt.Errorf("failed to update booking status: %w", err)
	}

	log.Printf("Successfully updated booking %s status to %s", bookingID, status)
	return nil
}

// ConfirmBooking confirms a pending booking
func (s *BookingService) ConfirmBooking(bookingID string) error {
	return s.UpdateBookingStatus(bookingID, "confirmed")
}

// CancelBooking cancels a booking
func (s *BookingService) CancelBooking(bookingID string) error {
	return s.UpdateBookingStatus(bookingID, "cancelled")
}

// GetBookingsByShow retrieves all bookings for a specific show
func (s *BookingService) GetBookingsByShow(showID uuid.UUID) ([]*bookings.Booking, error) {
	bookingsList, err := s.bookingRepository.GetBookingsByShow(showID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookings for show: %w", err)
	}

	return bookingsList, nil
}

// GetBookingsByContact retrieves bookings by contact information
func (s *BookingService) GetBookingsByContact(contactType, contactValue string) ([]*bookings.Booking, error) {
	if contactType == "" || contactValue == "" {
		return nil, fmt.Errorf("contact type and value cannot be empty")
	}

	bookingsList, err := s.bookingRepository.GetBookingsByContact(contactType, contactValue)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookings by contact: %w", err)
	}

	return bookingsList, nil
}

// SearchBookings searches bookings with filters and pagination
func (s *BookingService) SearchBookings(filter *bookings.BookingFilter) ([]*bookings.Booking, int, error) {
	if filter == nil {
		filter = &bookings.BookingFilter{}
	}

	// Set reasonable defaults for pagination if not specified
	if filter.Limit == 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100 // Prevent excessive load
	}

	bookingsList, totalCount, err := s.bookingRepository.GetBookingsWithFilters(filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search bookings: %w", err)
	}

	return bookingsList, totalCount, nil
}

// GetBookingStats retrieves booking statistics
func (s *BookingService) GetBookingStats() (*bookings.BookingStats, error) {
	stats, err := s.bookingRepository.GetBookingStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get booking statistics: %w", err)
	}

	return stats, nil
}

// ValidateBookingCapacity checks if tickets are available for booking
func (s *BookingService) ValidateBookingCapacity(showID uuid.UUID, requestedTickets int32) error {
	// Get show details
	show, err := s.showService.GetShow(showID.String())
	if err != nil {
		return fmt.Errorf("show not found: %w", err)
	}

	// Check capacity using repository
	return s.bookingRepository.ValidateBookingCapacity(showID, requestedTickets, show.Total_Tickets)
}

// GetShowBookingSummary provides booking summary for a show
func (s *BookingService) GetShowBookingSummary(showID uuid.UUID) (*ShowBookingSummary, error) {
	// Get show details
	show, err := s.showService.GetShow(showID.String())
	if err != nil {
		return nil, fmt.Errorf("show not found: %w", err)
	}

	// Get actual tickets sold from bookings
	ticketsSold, err := s.bookingRepository.GetTicketsSoldForShow(showID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tickets sold: %w", err)
	}

	// Get bookings for the show
	bookingsList, err := s.bookingRepository.GetBookingsByShow(showID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookings: %w", err)
	}

	// Calculate statistics
	totalRevenue := int32(0)
	bookingsByStatus := make(map[string]int32)

	for _, booking := range bookingsList {
		if booking.Status == "confirmed" || booking.Status == "pending" {
			totalRevenue += booking.TotalAmount
		}
		bookingsByStatus[booking.Status]++
	}

	summary := &ShowBookingSummary{
		ShowID:           showID,
		ShowName:         show.ShowName,
		ShowNumber:       show.ShowNumber,
		ShowDate:         show.ShowDate,
		TotalTickets:     show.Total_Tickets,
		TicketsSold:      ticketsSold,
		TicketsAvailable: show.Total_Tickets - ticketsSold,
		TotalBookings:    int32(len(bookingsList)),
		TotalRevenue:     totalRevenue,
		BookingsByStatus: bookingsByStatus,
		RecentBookings:   bookingsList[:min(len(bookingsList), 5)], // Last 5 bookings
	}

	return summary, nil
}

// DeleteBooking deletes a booking (admin function)
func (s *BookingService) DeleteBooking(bookingID string) error {
	if bookingID == "" {
		return fmt.Errorf("booking ID cannot be empty")
	}

	// Get booking details before deletion for show update
	booking, err := s.bookingRepository.GetBooking(bookingID)
	if err != nil {
		return fmt.Errorf("booking not found: %w", err)
	}

	// Delete the booking
	if err := s.bookingRepository.DeleteBooking(bookingID); err != nil {
		return fmt.Errorf("failed to delete booking: %w", err)
	}

	// Update show availability if booking was active
	if booking.Status == "confirmed" || booking.Status == "pending" {
		show, err := s.showService.GetShow(booking.ShowID.String())
		if err != nil {
			log.Printf("Warning: Failed to get show for booking deletion: %v", err)
		} else {
			// Reduce booked tickets count
			show.Booked_Tickets -= booking.NumberOfTickets
			if show.Booked_Tickets < 0 {
				show.Booked_Tickets = 0
			}
			if err := s.showService.UpdateShow(show); err != nil {
				log.Printf("Warning: Failed to update show availability after deletion: %v", err)
			}
		}
	}

	log.Printf("Successfully deleted booking: %s", bookingID)
	return nil
}

// ShowBookingSummary represents a comprehensive booking summary for a show
type ShowBookingSummary struct {
	ShowID           uuid.UUID           `json:"show_id"`
	ShowName         string              `json:"show_name"`
	ShowNumber       string              `json:"show_number"`
	ShowDate         time.Time           `json:"show_date"`
	TotalTickets     int32               `json:"total_tickets"`
	TicketsSold      int32               `json:"tickets_sold"`
	TicketsAvailable int32               `json:"tickets_available"`
	TotalBookings    int32               `json:"total_bookings"`
	TotalRevenue     int32               `json:"total_revenue"`
	BookingsByStatus map[string]int32    `json:"bookings_by_status"`
	RecentBookings   []*bookings.Booking `json:"recent_bookings"`
}

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
