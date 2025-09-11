package bookings

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewBooking(t *testing.T) {
	showID := uuid.New()
	contactType := "mobile"
	contactValue := "1234567890"
	numberOfTickets := int32(2)
	totalAmount := int32(200)

	booking := NewBooking(showID, contactType, contactValue, numberOfTickets, totalAmount)

	if booking == nil {
		t.Fatal("NewBooking should not return nil")
	}

	if booking.ShowID != showID {
		t.Errorf("Expected ShowID %v, got %v", showID, booking.ShowID)
	}

	if booking.ContactType != contactType {
		t.Errorf("Expected ContactType %s, got %s", contactType, booking.ContactType)
	}

	if booking.ContactValue != contactValue {
		t.Errorf("Expected ContactValue %s, got %s", contactValue, booking.ContactValue)
	}

	if booking.NumberOfTickets != numberOfTickets {
		t.Errorf("Expected NumberOfTickets %d, got %d", numberOfTickets, booking.NumberOfTickets)
	}

	if booking.TotalAmount != totalAmount {
		t.Errorf("Expected TotalAmount %d, got %d", totalAmount, booking.TotalAmount)
	}

	if booking.Status != "pending" {
		t.Errorf("Expected Status 'pending', got %s", booking.Status)
	}

	if booking.BookingID == "" {
		t.Error("BookingID should not be empty")
	}

	if !strings.HasPrefix(booking.BookingID, "BK-") {
		t.Errorf("BookingID should start with 'BK-', got %s", booking.BookingID)
	}
}

func TestNewBookingFromRequest(t *testing.T) {
	showID := uuid.New()
	
	tests := []struct {
		name           string
		queryParams    map[string]string
		expectedError  bool
		errorContains  string
	}{
		{
			name: "valid request",
			queryParams: map[string]string{
				"show_id":            showID.String(),
				"contact_type":       "mobile",
				"contact_value":      "1234567890",
				"number_of_tickets":  "2",
				"customer_name":      "John Doe",
			},
			expectedError: false,
		},
		{
			name: "missing show_id",
			queryParams: map[string]string{
				"contact_type":       "mobile",
				"contact_value":      "1234567890",
				"number_of_tickets":  "2",
			},
			expectedError: true,
			errorContains: "missing required parameters",
		},
		{
			name: "invalid show_id format",
			queryParams: map[string]string{
				"show_id":            "invalid-uuid",
				"contact_type":       "mobile",
				"contact_value":      "1234567890",
				"number_of_tickets":  "2",
			},
			expectedError: true,
			errorContains: "invalid show_id format",
		},
		{
			name: "invalid number_of_tickets",
			queryParams: map[string]string{
				"show_id":            showID.String(),
				"contact_type":       "mobile",
				"contact_value":      "1234567890",
				"number_of_tickets":  "invalid",
			},
			expectedError: true,
			errorContains: "invalid number_of_tickets",
		},
		{
			name: "zero tickets",
			queryParams: map[string]string{
				"show_id":            showID.String(),
				"contact_type":       "mobile",
				"contact_value":      "1234567890",
				"number_of_tickets":  "0",
			},
			expectedError: true,
			errorContains: "number_of_tickets must be greater than 0",
		},
		{
			name: "invalid contact_type",
			queryParams: map[string]string{
				"show_id":            showID.String(),
				"contact_type":       "invalid",
				"contact_value":      "1234567890",
				"number_of_tickets":  "2",
			},
			expectedError: true,
			errorContains: "contact_type must be either 'mobile' or 'email'",
		},
		{
			name: "invalid mobile number",
			queryParams: map[string]string{
				"show_id":            showID.String(),
				"contact_type":       "mobile",
				"contact_value":      "123",
				"number_of_tickets":  "2",
			},
			expectedError: true,
			errorContains: "invalid contact_value",
		},
		{
			name: "invalid email",
			queryParams: map[string]string{
				"show_id":            showID.String(),
				"contact_type":       "email",
				"contact_value":      "invalid-email",
				"number_of_tickets":  "2",
			},
			expectedError: true,
			errorContains: "invalid contact_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()

			booking, err := NewBookingFromRequest(req)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got '%s'", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
					return
				}
				if booking == nil {
					t.Error("Expected booking but got nil")
					return
				}
				if booking.ShowID != showID {
					t.Errorf("Expected ShowID %v, got %v", showID, booking.ShowID)
				}
			}
		})
	}
}

func TestNewBookingFromJSON(t *testing.T) {
	showID := uuid.New()
	
	tests := []struct {
		name           string
		jsonBody       string
		expectedError  bool
		errorContains  string
	}{
		{
			name: "valid JSON",
			jsonBody: `{
				"show_id": "` + showID.String() + `",
				"contact_type": "mobile",
				"contact_value": "1234567890",
				"number_of_tickets": 2,
				"customer_name": "John Doe"
			}`,
			expectedError: false,
		},
		{
			name: "invalid JSON",
			jsonBody: `{
				"show_id": "` + showID.String() + `",
				"contact_type": "mobile",
				"contact_value": "1234567890",
				"number_of_tickets": 2,
				"customer_name": "John Doe"
			`, // Missing closing brace
			expectedError: true,
			errorContains: "invalid JSON payload",
		},
		{
			name: "missing required fields",
			jsonBody: `{
				"show_id": "` + showID.String() + `",
				"contact_type": "mobile"
			}`,
			expectedError: true,
			errorContains: "missing required fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/test", strings.NewReader(tt.jsonBody))
			req.Header.Set("Content-Type", "application/json")

			booking, err := NewBookingFromJSON(req)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got '%s'", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
					return
				}
				if booking == nil {
					t.Error("Expected booking but got nil")
					return
				}
			}
		})
	}
}

func TestBookingUpdateStatus(t *testing.T) {
	booking := &Booking{
		Status:    "pending",
		UpdatedAt: time.Now().Add(-time.Hour),
	}

	newStatus := "confirmed"
	booking.UpdateStatus(newStatus)

	if booking.Status != newStatus {
		t.Errorf("Expected status %s, got %s", newStatus, booking.Status)
	}

	if booking.UpdatedAt.Before(time.Now().Add(-time.Minute)) {
		t.Error("UpdatedAt should be recent")
	}
}

func TestBookingIsValidStatus(t *testing.T) {
	tests := []struct {
		status string
		valid  bool
	}{
		{"pending", true},
		{"confirmed", true},
		{"cancelled", true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			booking := &Booking{Status: tt.status}
			if booking.IsValidStatus() != tt.valid {
				t.Errorf("Expected IsValidStatus() to return %v for status '%s'", tt.valid, tt.status)
			}
		})
	}
}

func TestBookingToJSON(t *testing.T) {
	booking := &Booking{
		BookingID:       "BK-1234567890ABCDEF",
		ShowID:          uuid.New(),
		ContactType:     "mobile",
		ContactValue:    "1234567890",
		NumberOfTickets: 2,
		CustomerName:    "John Doe",
		TotalAmount:     200,
		BookingDate:     time.Now(),
		Status:          "pending",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	jsonStr, err := booking.ToJSON()
	if err != nil {
		t.Errorf("ToJSON() error: %v", err)
	}

	if jsonStr == "" {
		t.Error("ToJSON() should not return empty string")
	}

	// Test that it's valid JSON by unmarshaling it back
	var unmarshaled Booking
	if err := json.Unmarshal([]byte(jsonStr), &unmarshaled); err != nil {
		t.Errorf("ToJSON() produced invalid JSON: %v", err)
	}
}

func TestBookingFromJSON(t *testing.T) {
	original := &Booking{
		BookingID:       "BK-1234567890ABCDEF",
		ShowID:          uuid.New(),
		ContactType:     "mobile",
		ContactValue:    "1234567890",
		NumberOfTickets: 2,
		CustomerName:    "John Doe",
		TotalAmount:     200,
		BookingDate:     time.Now(),
		Status:          "pending",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	jsonStr, err := original.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() error: %v", err)
	}

	booking := &Booking{}
	if err := booking.FromJSON(jsonStr); err != nil {
		t.Errorf("FromJSON() error: %v", err)
	}

	if booking.BookingID != original.BookingID {
		t.Errorf("Expected BookingID %s, got %s", original.BookingID, booking.BookingID)
	}
}

func TestBookingToMap(t *testing.T) {
	booking := &Booking{
		BookingID:       "BK-1234567890ABCDEF",
		ShowID:          uuid.New(),
		ContactType:     "mobile",
		ContactValue:    "1234567890",
		NumberOfTickets: 2,
		CustomerName:    "John Doe",
		TotalAmount:     200,
		BookingDate:     time.Now(),
		Status:          "pending",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	m := booking.ToMap()

	if m["booking_id"] != booking.BookingID {
		t.Errorf("Expected booking_id %s, got %s", booking.BookingID, m["booking_id"])
	}

	if m["show_id"] != booking.ShowID.String() {
		t.Errorf("Expected show_id %s, got %s", booking.ShowID.String(), m["show_id"])
	}

	if m["contact_type"] != booking.ContactType {
		t.Errorf("Expected contact_type %s, got %s", booking.ContactType, m["contact_type"])
	}
}

func TestIsValidContactType(t *testing.T) {
	tests := []struct {
		contactType string
		valid       bool
	}{
		{"mobile", true},
		{"email", true},
		{"phone", false},
		{"", false},
		{"invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.contactType, func(t *testing.T) {
			if isValidContactType(tt.contactType) != tt.valid {
				t.Errorf("Expected isValidContactType('%s') to return %v", tt.contactType, tt.valid)
			}
		})
	}
}

func TestIsValidContactValue(t *testing.T) {
	tests := []struct {
		contactType  string
		contactValue string
		valid        bool
	}{
		{"mobile", "1234567890", true},
		{"mobile", "123456789012345", true},
		{"mobile", "123", false},
		{"mobile", "abc", false},
		{"email", "test@example.com", true},
		{"email", "test@example", false},
		{"email", "test", false},
		{"invalid", "anything", false},
	}

	for _, tt := range tests {
		t.Run(tt.contactType+"_"+tt.contactValue, func(t *testing.T) {
			if isValidContactValue(tt.contactType, tt.contactValue) != tt.valid {
				t.Errorf("Expected isValidContactValue('%s', '%s') to return %v", 
					tt.contactType, tt.contactValue, tt.valid)
			}
		})
	}
}
