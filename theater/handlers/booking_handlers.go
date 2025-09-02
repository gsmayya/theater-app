package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gsmayya/theater/bookings"
	"github.com/gsmayya/theater/service"
	"github.com/google/uuid"
)

var bookingService *service.BookingService

// InitializeBookingService initializes the booking service
func InitializeBookingService() {
	bookingService = service.NewBookingService()
}

// CreateBookingHandler handles booking creation requests
func CreateBookingHandler(w http.ResponseWriter, r *http.Request) {
	if bookingService == nil {
		InitializeBookingService()
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var booking *bookings.Booking
	var err error

	// Handle both JSON and form-encoded requests
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		// Parse JSON body
		booking, err = bookings.NewBookingFromJSON(r)
	} else {
		// Parse form parameters
		booking, err = bookings.NewBookingFromRequest(r)
	}

	if err != nil {
		log.Printf("Error parsing booking request: %v", err)
		response := &bookings.BookingResponse{
			Success: false,
			Error:   err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create booking using service
	createdBooking, err := bookingService.CreateBooking(
		booking.ShowID,
		booking.ContactType,
		booking.ContactValue,
		booking.NumberOfTickets,
		booking.CustomerName,
	)

	if err != nil {
		log.Printf("Error creating booking: %v", err)
		response := &bookings.BookingResponse{
			Success: false,
			Error:   err.Error(),
		}
		
		// Determine appropriate HTTP status code
		if strings.Contains(err.Error(), "not found") {
			w.WriteHeader(http.StatusNotFound)
		} else if strings.Contains(err.Error(), "insufficient tickets") {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		
		json.NewEncoder(w).Encode(response)
		return
	}

	// Success response
	response := &bookings.BookingResponse{
		Success: true,
		Booking: createdBooking,
		Message: "Booking created successfully",
		Data: map[string]interface{}{
			"booking_id": createdBooking.BookingID,
			"show_id":    createdBooking.ShowID.String(),
		},
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetBookingHandler retrieves a specific booking by ID
func GetBookingHandler(w http.ResponseWriter, r *http.Request) {
	if bookingService == nil {
		InitializeBookingService()
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bookingID := r.URL.Query().Get("booking_id")
	if bookingID == "" {
		response := &bookings.BookingResponse{
			Success: false,
			Error:   "booking_id parameter is required",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	booking, err := bookingService.GetBooking(bookingID)
	if err != nil {
		log.Printf("Error getting booking: %v", err)
		response := &bookings.BookingResponse{
			Success: false,
			Error:   err.Error(),
		}
		
		if strings.Contains(err.Error(), "not found") {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		
		json.NewEncoder(w).Encode(response)
		return
	}

	response := &bookings.BookingResponse{
		Success: true,
		Booking: booking,
		Message: "Booking retrieved successfully",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// UpdateBookingStatusHandler updates a booking's status
func UpdateBookingStatusHandler(w http.ResponseWriter, r *http.Request) {
	if bookingService == nil {
		InitializeBookingService()
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "PUT" && r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bookingID := r.URL.Query().Get("booking_id")
	status := r.URL.Query().Get("status")

	if bookingID == "" || status == "" {
		response := &bookings.BookingResponse{
			Success: false,
			Error:   "booking_id and status parameters are required",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err := bookingService.UpdateBookingStatus(bookingID, status)
	if err != nil {
		log.Printf("Error updating booking status: %v", err)
		response := &bookings.BookingResponse{
			Success: false,
			Error:   err.Error(),
		}
		
		if strings.Contains(err.Error(), "not found") {
			w.WriteHeader(http.StatusNotFound)
		} else if strings.Contains(err.Error(), "invalid status") {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		
		json.NewEncoder(w).Encode(response)
		return
	}

	response := &bookings.BookingResponse{
		Success: true,
		Message: "Booking status updated successfully",
		Data: map[string]interface{}{
			"booking_id": bookingID,
			"status":     status,
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetBookingsByShowHandler retrieves all bookings for a specific show
func GetBookingsByShowHandler(w http.ResponseWriter, r *http.Request) {
	if bookingService == nil {
		InitializeBookingService()
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	showIDStr := r.URL.Query().Get("show_id")
	if showIDStr == "" {
		response := &bookings.BookingResponse{
			Success: false,
			Error:   "show_id parameter is required",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	showID, err := uuid.Parse(showIDStr)
	if err != nil {
		response := &bookings.BookingResponse{
			Success: false,
			Error:   "invalid show_id format",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	bookingsList, err := bookingService.GetBookingsByShow(showID)
	if err != nil {
		log.Printf("Error getting bookings by show: %v", err)
		response := &bookings.BookingResponse{
			Success: false,
			Error:   err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"bookings": bookingsList,
		"count":    len(bookingsList),
		"show_id":  showIDStr,
	})
}

// GetBookingsByContactHandler retrieves bookings by contact information
func GetBookingsByContactHandler(w http.ResponseWriter, r *http.Request) {
	if bookingService == nil {
		InitializeBookingService()
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	contactType := r.URL.Query().Get("contact_type")
	contactValue := r.URL.Query().Get("contact_value")

	if contactType == "" || contactValue == "" {
		response := &bookings.BookingResponse{
			Success: false,
			Error:   "contact_type and contact_value parameters are required",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	bookingsList, err := bookingService.GetBookingsByContact(contactType, contactValue)
	if err != nil {
		log.Printf("Error getting bookings by contact: %v", err)
		response := &bookings.BookingResponse{
			Success: false,
			Error:   err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":       true,
		"bookings":      bookingsList,
		"count":         len(bookingsList),
		"contact_type":  contactType,
		"contact_value": contactValue,
	})
}

// SearchBookingsHandler handles booking search with filters
func SearchBookingsHandler(w http.ResponseWriter, r *http.Request) {
	if bookingService == nil {
		InitializeBookingService()
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" && r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse search filters
	filter := &bookings.BookingFilter{}

	// Parse query parameters
	if showIDStr := r.URL.Query().Get("show_id"); showIDStr != "" {
		if showID, err := uuid.Parse(showIDStr); err == nil {
			filter.ShowID = &showID
		}
	}

	filter.ContactType = r.URL.Query().Get("contact_type")
	filter.Status = r.URL.Query().Get("status")

	// Parse date filters
	if dateFromStr := r.URL.Query().Get("date_from"); dateFromStr != "" {
		if dateFrom, err := time.Parse(time.RFC3339, dateFromStr); err == nil {
			filter.DateFrom = &dateFrom
		}
	}

	if dateToStr := r.URL.Query().Get("date_to"); dateToStr != "" {
		if dateTo, err := time.Parse(time.RFC3339, dateToStr); err == nil {
			filter.DateTo = &dateTo
		}
	}

	// Parse pagination
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filter.Offset = offset
		}
	}

	bookingsList, totalCount, err := bookingService.SearchBookings(filter)
	if err != nil {
		log.Printf("Error searching bookings: %v", err)
		response := &bookings.BookingResponse{
			Success: false,
			Error:   err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"bookings":    bookingsList,
		"total_count": totalCount,
		"count":       len(bookingsList),
		"filter":      filter,
	})
}

// GetBookingStatsHandler retrieves booking statistics
func GetBookingStatsHandler(w http.ResponseWriter, r *http.Request) {
	if bookingService == nil {
		InitializeBookingService()
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := bookingService.GetBookingStats()
	if err != nil {
		log.Printf("Error getting booking stats: %v", err)
		response := &bookings.BookingResponse{
			Success: false,
			Error:   err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"stats":   stats,
	})
}

// GetShowBookingSummaryHandler gets booking summary for a specific show
func GetShowBookingSummaryHandler(w http.ResponseWriter, r *http.Request) {
	if bookingService == nil {
		InitializeBookingService()
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	showIDStr := r.URL.Query().Get("show_id")
	if showIDStr == "" {
		response := &bookings.BookingResponse{
			Success: false,
			Error:   "show_id parameter is required",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	showID, err := uuid.Parse(showIDStr)
	if err != nil {
		response := &bookings.BookingResponse{
			Success: false,
			Error:   "invalid show_id format",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	summary, err := bookingService.GetShowBookingSummary(showID)
	if err != nil {
		log.Printf("Error getting show booking summary: %v", err)
		response := &bookings.BookingResponse{
			Success: false,
			Error:   err.Error(),
		}
		
		if strings.Contains(err.Error(), "not found") {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"summary": summary,
	})
}

// ConfirmBookingHandler confirms a pending booking
func ConfirmBookingHandler(w http.ResponseWriter, r *http.Request) {
	if bookingService == nil {
		InitializeBookingService()
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "PUT" && r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bookingID := r.URL.Query().Get("booking_id")
	if bookingID == "" {
		response := &bookings.BookingResponse{
			Success: false,
			Error:   "booking_id parameter is required",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err := bookingService.ConfirmBooking(bookingID)
	if err != nil {
		log.Printf("Error confirming booking: %v", err)
		response := &bookings.BookingResponse{
			Success: false,
			Error:   err.Error(),
		}
		
		if strings.Contains(err.Error(), "not found") {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		
		json.NewEncoder(w).Encode(response)
		return
	}

	response := &bookings.BookingResponse{
		Success: true,
		Message: "Booking confirmed successfully",
		Data: map[string]interface{}{
			"booking_id": bookingID,
			"status":     "confirmed",
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CancelBookingHandler cancels a booking
func CancelBookingHandler(w http.ResponseWriter, r *http.Request) {
	if bookingService == nil {
		InitializeBookingService()
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "PUT" && r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bookingID := r.URL.Query().Get("booking_id")
	if bookingID == "" {
		response := &bookings.BookingResponse{
			Success: false,
			Error:   "booking_id parameter is required",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err := bookingService.CancelBooking(bookingID)
	if err != nil {
		log.Printf("Error cancelling booking: %v", err)
		response := &bookings.BookingResponse{
			Success: false,
			Error:   err.Error(),
		}
		
		if strings.Contains(err.Error(), "not found") {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		
		json.NewEncoder(w).Encode(response)
		return
	}

	response := &bookings.BookingResponse{
		Success: true,
		Message: "Booking cancelled successfully",
		Data: map[string]interface{}{
			"booking_id": bookingID,
			"status":     "cancelled",
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
