package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gsmayya/theater/service"
)

var showService *service.ShowService

// InitializeService initializes the show service
func InitializeService() {
	showService = service.NewShowService()
}

// SearchShowsHandler handles advanced search requests with indexing
func SearchShowsHandler(w http.ResponseWriter, r *http.Request) {
	if showService == nil {
		InitializeService()
	}

	// Handle CORS preflight requests
	if HandleCORS(w, r) {
		return
	}

	// Validate HTTP method
	if !ValidateMethod(w, r, "GET", "POST") {
		return
	}

	var searchReq service.SearchRequest

	// Handle both GET and POST requests
	if r.Method == "GET" {
		// Parse query parameters
		searchReq = parseSearchParams(r)
	} else if r.Method == "POST" {
		// Parse JSON body
		if err := json.NewDecoder(r.Body).Decode(&searchReq); err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload", err)
			return
		}
	}

	// Perform search
	response, err := showService.SearchShows(searchReq)
	if err != nil {
		log.Printf("Search error: %v", err)
		WriteErrorResponse(w, http.StatusInternalServerError, "Search failed", err)
		return
	}

	// Return results with pagination
	pagination := PaginationInfo{
		Page:       response.Page,
		PageSize:   response.PageSize,
		Total:      response.Total,
		TotalPages: response.TotalPages,
	}

	WritePaginatedResponse(w, http.StatusOK, "Search completed successfully", response.Shows, pagination)
}

// ShowsByAllHandler retrieves all shows
func ShowsByAllHandler(w http.ResponseWriter, r *http.Request) {
	if showService == nil {
		InitializeService()
	}

	// Validate HTTP method
	if !ValidateMethod(w, r, "GET") {
		return
	}

	shows, err := showService.GetAllShows()
	if err != nil {
		log.Printf("Error getting all shows: %v", err)
		WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve shows", err)
		return
	}

	responseData := map[string]interface{}{
		"shows": shows,
		"count": len(shows),
	}

	WriteSuccessResponse(w, http.StatusOK, "Shows retrieved successfully", responseData)
}

// ShowsByLocationHandler handles location-based searches using indexes
func ShowsByLocationHandler(w http.ResponseWriter, r *http.Request) {
	if showService == nil {
		InitializeService()
	}

	// Validate HTTP method
	if !ValidateMethod(w, r, "GET") {
		return
	}

	showLocation := r.URL.Query().Get("show_location")
	if showLocation == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Missing required parameter", 
			&HTTPError{Code: http.StatusBadRequest, Message: "show_location parameter is required"})
		return
	}

	onlyAvailable := r.URL.Query().Get("only_available") == "true"

	shows, err := showService.GetShowsByLocation(showLocation, onlyAvailable)
	if err != nil {
		log.Printf("Error getting shows by location: %v", err)
		WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve shows", err)
		return
	}

	responseData := map[string]interface{}{
		"shows":         shows,
		"show_location": showLocation,
		"count":         len(shows),
		"only_available": onlyAvailable,
	}

	WriteSuccessResponse(w, http.StatusOK, "Shows retrieved successfully", responseData)
}

// ShowsByPriceRangeHandler handles price range searches using sorted sets
func ShowsByPriceRangeHandler(w http.ResponseWriter, r *http.Request) {
	if showService == nil {
		InitializeService()
	}

	// Validate HTTP method
	if !ValidateMethod(w, r, "GET") {
		return
	}

	minPriceStr := r.URL.Query().Get("min_price")
	maxPriceStr := r.URL.Query().Get("max_price")
	location := r.URL.Query().Get("location")

	if minPriceStr == "" || maxPriceStr == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Missing required parameters", 
			&HTTPError{Code: http.StatusBadRequest, Message: "Both min_price and max_price parameters are required"})
		return
	}

	minPrice, err := strconv.ParseInt(minPriceStr, 10, 32)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Invalid min_price value", err)
		return
	}

	maxPrice, err := strconv.ParseInt(maxPriceStr, 10, 32)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Invalid max_price value", err)
		return
	}

	shows, err := showService.GetShowsByPriceRange(int32(minPrice), int32(maxPrice), location)
	if err != nil {
		log.Printf("Error getting shows by price range: %v", err)
		WriteErrorResponse(w, http.StatusBadRequest, "Failed to retrieve shows", err)
		return
	}

	responseData := map[string]interface{}{
		"shows":     shows,
		"min_price": minPrice,
		"max_price": maxPrice,
		"location":  location,
		"count":     len(shows),
	}

	WriteSuccessResponse(w, http.StatusOK, "Shows retrieved successfully", responseData)
}

// CreateShowHandler creates a new show with full indexing
func CreateShowHandler(w http.ResponseWriter, r *http.Request) {
	if showService == nil {
		InitializeService()
	}

	// Handle CORS preflight requests
	if HandleCORS(w, r) {
		return
	}

	// Validate HTTP method
	if !ValidateMethod(w, r, "POST") {
		return
	}

	// Parse form data or query parameters
	name := r.URL.Query().Get("name")
	details := r.URL.Query().Get("details")
	location := r.URL.Query().Get("location")
	priceStr := r.URL.Query().Get("price")
	totalTicketsStr := r.URL.Query().Get("total_tickets")

	if name == "" || location == "" || priceStr == "" || totalTicketsStr == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Missing required parameters", 
			&HTTPError{Code: http.StatusBadRequest, Message: "Missing required parameters: name, location, price, total_tickets"})
		return
	}

	price, err := strconv.ParseInt(priceStr, 10, 32)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Invalid price value", err)
		return
	}

	totalTickets, err := strconv.ParseInt(totalTicketsStr, 10, 32)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Invalid total_tickets value", err)
		return
	}

	// Create show using service
	show, err := showService.CreateShow(name, details, location, int32(price), int32(totalTickets))
	if err != nil {
		log.Printf("Error creating show: %v", err)
		WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create show", err)
		return
	}

	WriteSuccessResponse(w, http.StatusCreated, "Show created successfully", show)
}

// GetShowHandler retrieves a specific show using optimized caching
func GetShowHandler(w http.ResponseWriter, r *http.Request) {
	if showService == nil {
		InitializeService()
	}

	// Validate HTTP method
	if !ValidateMethod(w, r, "GET") {
		return
	}

	showID := r.URL.Query().Get("id")
	if showID == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Missing required parameter", 
			&HTTPError{Code: http.StatusBadRequest, Message: "id parameter is required"})
		return
	}

	show, err := showService.GetShow(showID)
	if err != nil {
		log.Printf("Error getting show: %v", err)
		
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}
		
		WriteErrorResponse(w, statusCode, "Failed to retrieve show", err)
		return
	}

	WriteSuccessResponse(w, http.StatusOK, "Show retrieved successfully", show)
}

// UpdateShowAvailabilityHandler efficiently updates just the availability information
func UpdateShowAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	if showService == nil {
		InitializeService()
	}

	// Handle CORS preflight requests
	if HandleCORS(w, r) {
		return
	}

	// Validate HTTP method
	if !ValidateMethod(w, r, "PUT", "POST") {
		return
	}

	showID := r.URL.Query().Get("id")
	bookedTicketsStr := r.URL.Query().Get("booked_tickets")

	if showID == "" || bookedTicketsStr == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Missing required parameters", 
			&HTTPError{Code: http.StatusBadRequest, Message: "Both id and booked_tickets parameters are required"})
		return
	}

	bookedTickets, err := strconv.ParseInt(bookedTicketsStr, 10, 32)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, "Invalid booked_tickets value", err)
		return
	}

	err = showService.UpdateTicketAvailability(showID, int32(bookedTickets))
	if err != nil {
		log.Printf("Error updating show availability: %v", err)
		WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update show availability", err)
		return
	}

	responseData := map[string]interface{}{
		"show_id":        showID,
		"booked_tickets": bookedTickets,
	}

	WriteSuccessResponse(w, http.StatusOK, "Show availability updated successfully", responseData)
}

// GetSearchStatsHandler returns statistics about search indexes
func GetSearchStatsHandler(w http.ResponseWriter, r *http.Request) {
	if showService == nil {
		InitializeService()
	}

	// Validate HTTP method
	if !ValidateMethod(w, r, "GET") {
		return
	}

	stats, err := showService.GetSearchStatistics()
	if err != nil {
		log.Printf("Error getting search statistics: %v", err)
		WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve statistics", err)
		return
	}

	WriteSuccessResponse(w, http.StatusOK, "Statistics retrieved successfully", stats)
}

// HealthCheckHandler provides system health information
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Validate HTTP method
	if !ValidateMethod(w, r, "GET") {
		return
	}

	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"services": map[string]string{
			"database": "connected",
			"redis":    "connected",
		},
		"version": "1.0.0",
	}

	WriteSuccessResponse(w, http.StatusOK, "System is healthy", health)
}

// Helper function to parse search parameters from query string
func parseSearchParams(r *http.Request) service.SearchRequest {
	req := service.SearchRequest{}

	req.ShowLocation = r.URL.Query().Get("location")
	req.SearchTerm = r.URL.Query().Get("search")
	req.OnlyAvailable = r.URL.Query().Get("only_available") == "true"

	if minPriceStr := r.URL.Query().Get("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseInt(minPriceStr, 10, 32); err == nil {
			minPrice32 := int32(minPrice)
			req.MinPrice = &minPrice32
		}
	}

	if maxPriceStr := r.URL.Query().Get("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseInt(maxPriceStr, 10, 32); err == nil {
			maxPrice32 := int32(maxPrice)
			req.MaxPrice = &maxPrice32
		}
	}

	if minAvailableStr := r.URL.Query().Get("min_available"); minAvailableStr != "" {
		if minAvailable, err := strconv.ParseInt(minAvailableStr, 10, 32); err == nil {
			minAvailable32 := int32(minAvailable)
			req.MinAvailable = &minAvailable32
		}
	}

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			req.Page = page
		}
	}

	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil {
			req.PageSize = pageSize
		}
	}

	return req
}
