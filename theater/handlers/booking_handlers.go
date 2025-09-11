package handlers

import (
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

	// Handle CORS preflight requests
	if HandleCORS(w, r) {
		return
	}

	// Validate HTTP method
	if !ValidateMethod(w, r, "POST") {
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
		WriteErrorResponse(w, http.StatusBadRequest, "Invalid booking request", err)
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
		
		// Determine appropriate HTTP status code based on error type
		var statusCode int
		switch {
		case strings.Contains(err.Error(), "not found"):
			statusCode = http.StatusNotFound
		case strings.Contains(err.Error(), "insufficient tickets"):
			statusCode = http.StatusConflict
		default:
			statusCode = http.StatusInternalServerError
		}
		
		WriteErrorResponse(w, statusCode, "Failed to create booking", err)
		return
	}

	// Success response
	responseData := map[string]interface{}{
		"booking_id": createdBooking.BookingID,
		"show_id":    createdBooking.ShowID.String(),
		"booking":    createdBooking,
	}

	WriteSuccessResponse(w, http.StatusCreated, "Booking created successfully", responseData)
}

// GetBookingHandler retrieves a specific booking by ID
func GetBookingHandler(w http.ResponseWriter, r *http.Request) {
	if bookingService == nil {
		InitializeBookingService()
	}

	// Validate HTTP method
	if !ValidateMethod(w, r, "GET") {
		return
	}

	bookingID := r.URL.Query().Get("booking_id")
	if bookingID == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Missing required parameter", 
			&HTTPError{Code: http.StatusBadRequest, Message: "booking_id parameter is required"})
		return
	}

	booking, err := bookingService.GetBooking(bookingID)
	if err != nil {
		log.Printf("Error getting booking: %v", err)
		
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}
		
		WriteErrorResponse(w, statusCode, "Failed to retrieve booking", err)
		return
	}

	WriteSuccessResponse(w, http.StatusOK, "Booking retrieved successfully", booking)
}

// UpdateBookingStatusHandler updates a booking's status
func UpdateBookingStatusHandler(w http.ResponseWriter, r *http.Request) {
	if bookingService == nil {
		InitializeBookingService()
	}

	// Handle CORS preflight requests
	if HandleCORS(w, r) {
		return
	}

	// Validate HTTP method
	if !ValidateMethod(w, r, "PUT", "POST") {
		return
	}

	bookingID := r.URL.Query().Get("booking_id")
	status := r.URL.Query().Get("status")

	if bookingID == "" || status == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Missing required parameters", 
			&HTTPError{Code: http.StatusBadRequest, Message: "booking_id and status parameters are required"})
		return
	}

	err := bookingService.UpdateBookingStatus(bookingID, status)
	if err != nil {
		log.Printf("Error updating booking status: %v", err)
		
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		} else if strings.Contains(err.Error(), "invalid status") {
			statusCode = http.StatusBadRequest
		}
		
		WriteErrorResponse(w, statusCode, "Failed to update booking status", err)
		return
	}

	responseData := map[string]interface{}{
		"booking_id": bookingID,
		"status":     status,
	}

	WriteSuccessResponse(w, http.StatusOK, "Booking status updated successfully", responseData)
}

// GetBookingsByShowHandler retrieves all bookings for a specific show
func GetBookingsByShowHandler(w http.ResponseWriter, r *http.Request) {
	if bookingService == nil {
		InitializeBookingService()
	}

	// Validate HTTP method
	if !ValidateMethod(w, r, "GET") {
		return
	}

	showIDStr := r.URL.Query().Get("show_id")
	if showIDStr == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Missing required parameter", 
			&HTTPError{Code: http.StatusBadRequest, Message: "show_id parameter is required"})
		return
	}

	showID, err := uuid.Parse(showIDStr)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Invalid show ID format", err)
		return
	}

	bookingsList, err := bookingService.GetBookingsByShow(showID)
	if err != nil {
		log.Printf("Error getting bookings by show: %v", err)
		WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve bookings", err)
		return
	}

	responseData := map[string]interface{}{
		"bookings": bookingsList,
		"count":    len(bookingsList),
		"show_id":  showIDStr,
	}

	WriteSuccessResponse(w, http.StatusOK, "Bookings retrieved successfully", responseData)
}

// GetBookingsByContactHandler retrieves bookings by contact information
func GetBookingsByContactHandler(w http.ResponseWriter, r *http.Request) {
	if bookingService == nil {
		InitializeBookingService()
	}

	// Validate HTTP method
	if !ValidateMethod(w, r, "GET") {
		return
	}

	contactType := r.URL.Query().Get("contact_type")
	contactValue := r.URL.Query().Get("contact_value")

	if contactType == "" || contactValue == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Missing required parameters", 
			&HTTPError{Code: http.StatusBadRequest, Message: "contact_type and contact_value parameters are required"})
		return
	}

	bookingsList, err := bookingService.GetBookingsByContact(contactType, contactValue)
	if err != nil {
		log.Printf("Error getting bookings by contact: %v", err)
		WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve bookings", err)
		return
	}

	responseData := map[string]interface{}{
		"bookings":      bookingsList,
		"count":         len(bookingsList),
		"contact_type":  contactType,
		"contact_value": contactValue,
	}

	WriteSuccessResponse(w, http.StatusOK, "Bookings retrieved successfully", responseData)
}

// SearchBookingsHandler handles booking search with filters
func SearchBookingsHandler(w http.ResponseWriter, r *http.Request) {
	if bookingService == nil {
		InitializeBookingService()
	}

	// Handle CORS preflight requests
	if HandleCORS(w, r) {
		return
	}

	// Validate HTTP method
	if !ValidateMethod(w, r, "GET", "POST") {
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
		WriteErrorResponse(w, http.StatusInternalServerError, "Failed to search bookings", err)
		return
	}

	// Calculate pagination info
	page := 1
	pageSize := 20
	if filter.Limit > 0 {
		pageSize = filter.Limit
	}
	if filter.Offset > 0 {
		page = (filter.Offset / pageSize) + 1
	}

	pagination := PaginationInfo{
		Page:       page,
		PageSize:   pageSize,
		Total:      totalCount,
		TotalPages: (totalCount + pageSize - 1) / pageSize,
	}

	responseData := map[string]interface{}{
		"bookings": bookingsList,
		"filter":   filter,
	}

	WritePaginatedResponse(w, http.StatusOK, "Bookings retrieved successfully", responseData, pagination)
}

// GetBookingStatsHandler retrieves booking statistics
func GetBookingStatsHandler(w http.ResponseWriter, r *http.Request) {
	if bookingService == nil {
		InitializeBookingService()
	}

	// Validate HTTP method
	if !ValidateMethod(w, r, "GET") {
		return
	}

	stats, err := bookingService.GetBookingStats()
	if err != nil {
		log.Printf("Error getting booking stats: %v", err)
		WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve booking statistics", err)
		return
	}

	WriteSuccessResponse(w, http.StatusOK, "Booking statistics retrieved successfully", stats)
}

// GetShowBookingSummaryHandler gets booking summary for a specific show
func GetShowBookingSummaryHandler(w http.ResponseWriter, r *http.Request) {
	if bookingService == nil {
		InitializeBookingService()
	}

	// Validate HTTP method
	if !ValidateMethod(w, r, "GET") {
		return
	}

	showIDStr := r.URL.Query().Get("show_id")
	if showIDStr == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Missing required parameter", 
			&HTTPError{Code: http.StatusBadRequest, Message: "show_id parameter is required"})
		return
	}

	showID, err := uuid.Parse(showIDStr)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Invalid show ID format", err)
		return
	}

	summary, err := bookingService.GetShowBookingSummary(showID)
	if err != nil {
		log.Printf("Error getting show booking summary: %v", err)
		
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}
		
		WriteErrorResponse(w, statusCode, "Failed to retrieve booking summary", err)
		return
	}

	WriteSuccessResponse(w, http.StatusOK, "Booking summary retrieved successfully", summary)
}

// ConfirmBookingHandler confirms a pending booking
func ConfirmBookingHandler(w http.ResponseWriter, r *http.Request) {
	if bookingService == nil {
		InitializeBookingService()
	}

	// Handle CORS preflight requests
	if HandleCORS(w, r) {
		return
	}

	// Validate HTTP method
	if !ValidateMethod(w, r, "PUT", "POST") {
		return
	}

	bookingID := r.URL.Query().Get("booking_id")
	if bookingID == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Missing required parameter", 
			&HTTPError{Code: http.StatusBadRequest, Message: "booking_id parameter is required"})
		return
	}

	err := bookingService.ConfirmBooking(bookingID)
	if err != nil {
		log.Printf("Error confirming booking: %v", err)
		
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}
		
		WriteErrorResponse(w, statusCode, "Failed to confirm booking", err)
		return
	}

	responseData := map[string]interface{}{
		"booking_id": bookingID,
		"status":     "confirmed",
	}

	WriteSuccessResponse(w, http.StatusOK, "Booking confirmed successfully", responseData)
}

// CancelBookingHandler cancels a booking
func CancelBookingHandler(w http.ResponseWriter, r *http.Request) {
	if bookingService == nil {
		InitializeBookingService()
	}

	// Handle CORS preflight requests
	if HandleCORS(w, r) {
		return
	}

	// Validate HTTP method
	if !ValidateMethod(w, r, "PUT", "POST") {
		return
	}

	bookingID := r.URL.Query().Get("booking_id")
	if bookingID == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Missing required parameter", 
			&HTTPError{Code: http.StatusBadRequest, Message: "booking_id parameter is required"})
		return
	}

	err := bookingService.CancelBooking(bookingID)
	if err != nil {
		log.Printf("Error cancelling booking: %v", err)
		
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}
		
		WriteErrorResponse(w, statusCode, "Failed to cancel booking", err)
		return
	}

	responseData := map[string]interface{}{
		"booking_id": bookingID,
		"status":     "cancelled",
	}

	WriteSuccessResponse(w, http.StatusOK, "Booking cancelled successfully", responseData)
}
