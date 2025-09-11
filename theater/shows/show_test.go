package shows

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestShowDataNewShow(t *testing.T) {
	show := &ShowData{}
	
	showName := "Test Show"
	details := "Test show details"
	price := int32(100)
	totalTickets := int32(50)
	showLocation := "Test Location"

	result := show.NewShow(showName, details, price, totalTickets, showLocation)

	if result == nil {
		t.Fatal("NewShow should not return nil")
	}

	if show.ShowName != showName {
		t.Errorf("Expected ShowName %s, got %s", showName, show.ShowName)
	}

	if show.Details != details {
		t.Errorf("Expected Details %s, got %s", details, show.Details)
	}

	if show.Price != price {
		t.Errorf("Expected Price %d, got %d", price, show.Price)
	}

	if show.Total_Tickets != totalTickets {
		t.Errorf("Expected Total_Tickets %d, got %d", totalTickets, show.Total_Tickets)
	}

	if show.ShowLocation != showLocation {
		t.Errorf("Expected ShowLocation %s, got %s", showLocation, show.ShowLocation)
	}

	if show.Booked_Tickets != 0 {
		t.Errorf("Expected Booked_Tickets 0, got %d", show.Booked_Tickets)
	}

	if show.Show_Id == uuid.Nil {
		t.Error("Show_Id should not be nil")
	}

	if show.ShowNumber == "" {
		t.Error("ShowNumber should not be empty")
	}

	if !strings.HasPrefix(show.ShowNumber, "SH-") {
		t.Errorf("ShowNumber should start with 'SH-', got %s", show.ShowNumber)
	}

	if show.ShowDate.IsZero() {
		t.Error("ShowDate should not be zero")
	}

	if show.Images == nil {
		t.Error("Images should not be nil")
	}

	if show.Videos == nil {
		t.Error("Videos should not be nil")
	}
}

func TestShowDataNewShowFromPut(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		expectedError  bool
		errorContains  string
	}{
		{
			name: "valid request",
			queryParams: map[string]string{
				"show_name":      "Test Show",
				"details":        "Test details",
				"price":          "100",
				"total_tickets":  "50",
				"show_location":  "Test Location",
				"show_number":    "SH-001",
				"show_date":      "2024-12-31T23:59:59Z",
			},
			expectedError: false,
		},
		{
			name: "invalid price",
			queryParams: map[string]string{
				"show_name":      "Test Show",
				"details":        "Test details",
				"price":          "invalid",
				"total_tickets":  "50",
				"show_location":  "Test Location",
			},
			expectedError: true,
		},
		{
			name: "invalid total_tickets",
			queryParams: map[string]string{
				"show_name":      "Test Show",
				"details":        "Test details",
				"price":          "100",
				"total_tickets":  "invalid",
				"show_location":  "Test Location",
			},
			expectedError: true,
		},
		{
			name: "invalid date format",
			queryParams: map[string]string{
				"show_name":      "Test Show",
				"details":        "Test details",
				"price":          "100",
				"total_tickets":  "50",
				"show_location":  "Test Location",
				"show_date":      "invalid-date",
			},
			expectedError: false, // Should use default date
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("PUT", "/test", nil)
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()

			show := &ShowData{}
			result := show.NewShowFromPut(req)

			if tt.expectedError {
				if result != nil {
					t.Errorf("Expected nil result but got: %v", result)
				}
			} else {
				if result == nil {
					t.Error("Expected show but got nil")
					return
				}

				if show.ShowName != tt.queryParams["show_name"] {
					t.Errorf("Expected ShowName %s, got %s", tt.queryParams["show_name"], show.ShowName)
				}

				if show.Show_Id == uuid.Nil {
					t.Error("Show_Id should not be nil")
				}

				if show.Booked_Tickets != 0 {
					t.Errorf("Expected Booked_Tickets 0, got %d", show.Booked_Tickets)
				}
			}
		})
	}
}

func TestShowDataShowToMap(t *testing.T) {
	show := &ShowData{
		Show_Id:        uuid.New(),
		ShowName:       "Test Show",
		Details:        "Test details",
		Price:          100,
		Total_Tickets:  50,
		Booked_Tickets: 10,
		ShowLocation:   "Test Location",
		ShowNumber:     "SH-001",
		ShowDate:       time.Now(),
		Images:         []string{"image1.jpg", "image2.jpg"},
		Videos:         []string{"video1.mp4"},
	}

	m := show.ShowToMap()

	expectedFields := []string{
		"show_id", "show_name", "details", "price", "total_tickets",
		"show_location", "booked_tickets", "show_number", "show_date",
		"images", "videos",
	}

	for _, field := range expectedFields {
		if _, exists := m[field]; !exists {
			t.Errorf("Expected field '%s' in map", field)
		}
	}

	if m["show_name"] != show.ShowName {
		t.Errorf("Expected show_name %s, got %s", show.ShowName, m["show_name"])
	}

	if m["price"] != "100" {
		t.Errorf("Expected price '100', got %s", m["price"])
	}
}

func TestShowDataShowToJSON(t *testing.T) {
	show := &ShowData{
		Show_Id:        uuid.New(),
		ShowName:       "Test Show",
		Details:        "Test details",
		Price:          100,
		Total_Tickets:  50,
		Booked_Tickets: 10,
		ShowLocation:   "Test Location",
		ShowNumber:     "SH-001",
		ShowDate:       time.Now(),
		Images:         []string{"image1.jpg"},
		Videos:         []string{"video1.mp4"},
	}

	jsonStr, err := show.ShowToJSON()
	if err != nil {
		t.Errorf("ShowToJSON() error: %v", err)
	}

	if jsonStr == "" {
		t.Error("ShowToJSON() should not return empty string")
	}

	// Test that it's valid JSON by checking if it contains expected fields
	if !strings.Contains(jsonStr, "Test Show") {
		t.Error("JSON should contain show name")
	}

	if !strings.Contains(jsonStr, "100") {
		t.Error("JSON should contain price")
	}
}

func TestShowDataJSONToShow(t *testing.T) {
	original := &ShowData{
		Show_Id:        uuid.New(),
		ShowName:       "Test Show",
		Details:        "Test details",
		Price:          100,
		Total_Tickets:  50,
		Booked_Tickets: 10,
		ShowLocation:   "Test Location",
		ShowNumber:     "SH-001",
		ShowDate:       time.Now(),
		Images:         []string{"image1.jpg"},
		Videos:         []string{"video1.mp4"},
	}

	jsonStr, err := original.ShowToJSON()
	if err != nil {
		t.Fatalf("ShowToJSON() error: %v", err)
	}

	show := &ShowData{}
	result, err := show.JSONToShow(jsonStr)
	if err != nil {
		t.Errorf("JSONToShow() error: %v", err)
	}

	if result == nil {
		t.Error("JSONToShow() should not return nil")
	}

	if show.ShowName != original.ShowName {
		t.Errorf("Expected ShowName %s, got %s", original.ShowName, show.ShowName)
	}

	if show.Price != original.Price {
		t.Errorf("Expected Price %d, got %d", original.Price, show.Price)
	}

	if len(show.Images) != len(original.Images) {
		t.Errorf("Expected %d images, got %d", len(original.Images), len(show.Images))
	}
}

func TestShowDataJSONToShowInvalidJSON(t *testing.T) {
	show := &ShowData{}
	result, err := show.JSONToShow("invalid json")
	
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}

	if result != nil {
		t.Error("Expected nil result for invalid JSON")
	}
}

func TestShowDataDefaultValues(t *testing.T) {
	show := &ShowData{}
	show.NewShow("Test", "Details", 100, 50, "Location")

	// Test default values
	if show.Booked_Tickets != 0 {
		t.Errorf("Expected Booked_Tickets 0, got %d", show.Booked_Tickets)
	}

	if show.ShowNumber == "" {
		t.Error("ShowNumber should not be empty")
	}

	if !strings.HasPrefix(show.ShowNumber, "SH-") {
		t.Errorf("ShowNumber should start with 'SH-', got %s", show.ShowNumber)
	}

	if show.ShowDate.IsZero() {
		t.Error("ShowDate should not be zero")
	}

	// ShowDate should be approximately 30 days from now
	expectedDate := time.Now().AddDate(0, 0, 30)
	timeDiff := show.ShowDate.Sub(expectedDate)
	if timeDiff > time.Hour || timeDiff < -time.Hour {
		t.Errorf("ShowDate should be approximately 30 days from now, got %v", show.ShowDate)
	}

	if show.Images == nil {
		t.Error("Images should not be nil")
	}

	if show.Videos == nil {
		t.Error("Videos should not be nil")
	}

	if len(show.Images) != 0 {
		t.Errorf("Expected empty Images slice, got %v", show.Images)
	}

	if len(show.Videos) != 0 {
		t.Errorf("Expected empty Videos slice, got %v", show.Videos)
	}
}