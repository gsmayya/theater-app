package bookings

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Booking represents a theater booking
type Booking struct {
	BookingID       string    `json:"booking_id"`       // Hash-generated unique ID
	ShowID          uuid.UUID `json:"show_id"`          // Reference to the show
	ContactType     string    `json:"contact_type"`     // "mobile" or "email"
	ContactValue    string    `json:"contact_value"`    // Mobile number or email address
	NumberOfTickets int32     `json:"number_of_tickets"`
	BookingDate     time.Time `json:"booking_date"`
	Status          string    `json:"status"`           // "confirmed", "pending", "cancelled"
	CustomerName    string    `json:"customer_name,omitempty"`
	TotalAmount     int32     `json:"total_amount"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// BookingRequest represents the request payload for creating a booking
type BookingRequest struct {
	ShowID          string `json:"show_id"`
	ContactType     string `json:"contact_type"`     // "mobile" or "email"
	ContactValue    string `json:"contact_value"`
	NumberOfTickets int32  `json:"number_of_tickets"`
	CustomerName    string `json:"customer_name,omitempty"`
}

// NewBooking creates a new booking with generated hash ID
func NewBooking(showID uuid.UUID, contactType, contactValue string, numberOfTickets int32, totalAmount int32) *Booking {
	now := time.Now()
	
	booking := &Booking{
		ShowID:          showID,
		ContactType:     contactType,
		ContactValue:    contactValue,
		NumberOfTickets: numberOfTickets,
		TotalAmount:     totalAmount,
		BookingDate:     now,
		Status:          "pending",
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	
	// Generate hash-based unique ID
	booking.BookingID = booking.generateHashID()
	return booking
}

// NewBookingFromRequest creates a booking from HTTP request
func NewBookingFromRequest(r *http.Request) (*Booking, error) {
	showIDStr := r.URL.Query().Get("show_id")
	contactType := r.URL.Query().Get("contact_type")
	contactValue := r.URL.Query().Get("contact_value")
	numberOfTicketsStr := r.URL.Query().Get("number_of_tickets")
	customerName := r.URL.Query().Get("customer_name")
	
	if showIDStr == "" || contactType == "" || contactValue == "" || numberOfTicketsStr == "" {
		return nil, fmt.Errorf("missing required parameters: show_id, contact_type, contact_value, number_of_tickets")
	}
	
	showID, err := uuid.Parse(showIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid show_id format: %w", err)
	}
	
	numberOfTickets, err := strconv.ParseInt(numberOfTicketsStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid number_of_tickets: %w", err)
	}
	
	if numberOfTickets <= 0 {
		return nil, fmt.Errorf("number_of_tickets must be greater than 0")
	}
	
	if !isValidContactType(contactType) {
		return nil, fmt.Errorf("contact_type must be either 'mobile' or 'email'")
	}
	
	if !isValidContactValue(contactType, contactValue) {
		return nil, fmt.Errorf("invalid contact_value for contact_type %s", contactType)
	}
	
	now := time.Now()
	booking := &Booking{
		ShowID:          showID,
		ContactType:     contactType,
		ContactValue:    contactValue,
		NumberOfTickets: int32(numberOfTickets),
		CustomerName:    customerName,
		BookingDate:     now,
		Status:          "pending",
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	
	booking.BookingID = booking.generateHashID()
	return booking, nil
}

// NewBookingFromJSON creates a booking from JSON request body
func NewBookingFromJSON(r *http.Request) (*Booking, error) {
	var req BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid JSON payload: %w", err)
	}
	
	if req.ShowID == "" || req.ContactType == "" || req.ContactValue == "" || req.NumberOfTickets <= 0 {
		return nil, fmt.Errorf("missing required fields: show_id, contact_type, contact_value, number_of_tickets")
	}
	
	showID, err := uuid.Parse(req.ShowID)
	if err != nil {
		return nil, fmt.Errorf("invalid show_id format: %w", err)
	}
	
	if !isValidContactType(req.ContactType) {
		return nil, fmt.Errorf("contact_type must be either 'mobile' or 'email'")
	}
	
	if !isValidContactValue(req.ContactType, req.ContactValue) {
		return nil, fmt.Errorf("invalid contact_value for contact_type %s", req.ContactType)
	}
	
	now := time.Now()
	booking := &Booking{
		ShowID:          showID,
		ContactType:     req.ContactType,
		ContactValue:    req.ContactValue,
		NumberOfTickets: req.NumberOfTickets,
		CustomerName:    req.CustomerName,
		BookingDate:     now,
		Status:          "pending",
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	
	booking.BookingID = booking.generateHashID()
	return booking, nil
}

// generateHashID creates a unique hash-based ID from booking information
func (b *Booking) generateHashID() string {
	// Create hash from show_id + contact_type + contact_value + booking_date + number_of_tickets
	hashInput := fmt.Sprintf("%s:%s:%s:%d:%d",
		b.ShowID.String(),
		b.ContactType,
		b.ContactValue,
		b.BookingDate.Unix(),
		b.NumberOfTickets,
	)
	
	hash := sha256.Sum256([]byte(hashInput))
	// Return first 16 characters of hex representation for a shorter, readable ID
	return fmt.Sprintf("BK-%X", hash)[:18] // BK- prefix + 16 hex chars
}

// UpdateStatus updates the booking status and timestamp
func (b *Booking) UpdateStatus(status string) {
	b.Status = status
	b.UpdatedAt = time.Now()
}

// IsValidStatus checks if the status is valid
func (b *Booking) IsValidStatus() bool {
	validStatuses := []string{"pending", "confirmed", "cancelled"}
	for _, valid := range validStatuses {
		if b.Status == valid {
			return true
		}
	}
	return false
}

// ToJSON converts booking to JSON string
func (b *Booking) ToJSON() (string, error) {
	jsonData, err := json.Marshal(b)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// FromJSON populates booking from JSON string
func (b *Booking) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), b)
}

// ToMap converts booking to map for easier handling
func (b *Booking) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"booking_id":        b.BookingID,
		"show_id":           b.ShowID.String(),
		"contact_type":      b.ContactType,
		"contact_value":     b.ContactValue,
		"number_of_tickets": b.NumberOfTickets,
		"customer_name":     b.CustomerName,
		"total_amount":      b.TotalAmount,
		"booking_date":      b.BookingDate.Format(time.RFC3339),
		"status":            b.Status,
		"created_at":        b.CreatedAt.Format(time.RFC3339),
		"updated_at":        b.UpdatedAt.Format(time.RFC3339),
	}
}

// BookingResponse represents the response structure for booking operations
type BookingResponse struct {
	Success bool                   `json:"success"`
	Booking *Booking               `json:"booking,omitempty"`
	Message string                 `json:"message,omitempty"`
	Error   string                 `json:"error,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// Helper functions for validation

func isValidContactType(contactType string) bool {
	return contactType == "mobile" || contactType == "email"
}

func isValidContactValue(contactType, contactValue string) bool {
	switch contactType {
	case "mobile":
		// Basic mobile number validation (10-15 digits, may include + and -)
		return len(contactValue) >= 10 && len(contactValue) <= 15 && 
			   strings.ContainsAny(contactValue, "0123456789")
	case "email":
		// Basic email validation
		return strings.Contains(contactValue, "@") && strings.Contains(contactValue, ".") &&
			   len(contactValue) > 5
	default:
		return false
	}
}

// GetBookingsByShow returns bookings for a specific show (to be used by repository)
type BookingFilter struct {
	ShowID      *uuid.UUID `json:"show_id,omitempty"`
	ContactType string     `json:"contact_type,omitempty"`
	Status      string     `json:"status,omitempty"`
	DateFrom    *time.Time `json:"date_from,omitempty"`
	DateTo      *time.Time `json:"date_to,omitempty"`
	Limit       int        `json:"limit,omitempty"`
	Offset      int        `json:"offset,omitempty"`
}

// BookingStats represents booking statistics
type BookingStats struct {
	TotalBookings     int32                  `json:"total_bookings"`
	TotalTickets      int32                  `json:"total_tickets"`
	TotalRevenue      int32                  `json:"total_revenue"`
	BookingsByStatus  map[string]int32       `json:"bookings_by_status"`
	BookingsByShow    map[string]int32       `json:"bookings_by_show"`
	RecentBookings    []*Booking             `json:"recent_bookings"`
}
