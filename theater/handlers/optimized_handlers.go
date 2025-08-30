package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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

	w.Header().Set("Content-Type", "application/json")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var searchReq service.SearchRequest

	// Handle both GET and POST requests
	if r.Method == "GET" {
		// Parse query parameters
		searchReq = parseSearchParams(r)
	} else if r.Method == "POST" {
		// Parse JSON body
		if err := json.NewDecoder(r.Body).Decode(&searchReq); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Perform search
	response, err := showService.SearchShows(searchReq)
	if err != nil {
		log.Printf("Search error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return results
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ShowsByLocationHandler handles location-based searches using indexes
func ShowsByLocationHandler(w http.ResponseWriter, r *http.Request) {
	if showService == nil {
		InitializeService()
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	location := r.URL.Query().Get("location")
	if location == "" {
		http.Error(w, "Location parameter is required", http.StatusBadRequest)
		return
	}

	onlyAvailable := r.URL.Query().Get("only_available") == "true"

	shows, err := showService.GetShowsByLocation(location, onlyAvailable)
	if err != nil {
		log.Printf("Error getting shows by location: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"shows":    shows,
		"location": location,
		"count":    len(shows),
	})
}

// ShowsByPriceRangeHandler handles price range searches using sorted sets
func ShowsByPriceRangeHandler(w http.ResponseWriter, r *http.Request) {
	if showService == nil {
		InitializeService()
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	minPriceStr := r.URL.Query().Get("min_price")
	maxPriceStr := r.URL.Query().Get("max_price")
	location := r.URL.Query().Get("location")

	if minPriceStr == "" || maxPriceStr == "" {
		http.Error(w, "Both min_price and max_price parameters are required", http.StatusBadRequest)
		return
	}

	minPrice, err := strconv.ParseInt(minPriceStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid min_price value", http.StatusBadRequest)
		return
	}

	maxPrice, err := strconv.ParseInt(maxPriceStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid max_price value", http.StatusBadRequest)
		return
	}

	shows, err := showService.GetShowsByPriceRange(int32(minPrice), int32(maxPrice), location)
	if err != nil {
		log.Printf("Error getting shows by price range: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"shows":     shows,
		"min_price": minPrice,
		"max_price": maxPrice,
		"location":  location,
		"count":     len(shows),
	})
}

// CreateShowHandler creates a new show with full indexing
func CreateShowHandler(w http.ResponseWriter, r *http.Request) {
	if showService == nil {
		InitializeService()
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data or query parameters
	name := r.URL.Query().Get("name")
	details := r.URL.Query().Get("details")
	location := r.URL.Query().Get("location")
	priceStr := r.URL.Query().Get("price")
	totalTicketsStr := r.URL.Query().Get("total_tickets")

	if name == "" || location == "" || priceStr == "" || totalTicketsStr == "" {
		http.Error(w, "Missing required parameters: name, location, price, total_tickets", http.StatusBadRequest)
		return
	}

	price, err := strconv.ParseInt(priceStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid price value", http.StatusBadRequest)
		return
	}

	totalTickets, err := strconv.ParseInt(totalTicketsStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid total_tickets value", http.StatusBadRequest)
		return
	}

	// Create show using service
	show, err := showService.CreateShow(name, details, location, int32(price), int32(totalTickets))
	if err != nil {
		log.Printf("Error creating show: %v", err)
		http.Error(w, "Failed to create show", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(show)
}

// GetShowHandler retrieves a specific show using optimized caching
func GetShowHandler(w http.ResponseWriter, r *http.Request) {
	if showService == nil {
		InitializeService()
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	showID := r.URL.Query().Get("id")
	if showID == "" {
		http.Error(w, "Show ID parameter is required", http.StatusBadRequest)
		return
	}

	show, err := showService.GetShow(showID)
	if err != nil {
		log.Printf("Error getting show: %v", err)
		if err.Error() == "show not found: "+showID {
			http.Error(w, "Show not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(show)
}

// UpdateShowAvailabilityHandler efficiently updates ticket availability
func UpdateShowAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	if showService == nil {
		InitializeService()
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != "PUT" && r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	showID := r.URL.Query().Get("id")
	bookedTicketsStr := r.URL.Query().Get("booked_tickets")

	if showID == "" || bookedTicketsStr == "" {
		http.Error(w, "Both id and booked_tickets parameters are required", http.StatusBadRequest)
		return
	}

	bookedTickets, err := strconv.ParseInt(bookedTicketsStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid booked_tickets value", http.StatusBadRequest)
		return
	}

	err = showService.UpdateTicketAvailability(showID, int32(bookedTickets))
	if err != nil {
		log.Printf("Error updating show availability: %v", err)
		http.Error(w, "Failed to update show availability", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"message": "Show availability updated",
	})
}

// GetSearchStatsHandler returns statistics about search indexes
func GetSearchStatsHandler(w http.ResponseWriter, r *http.Request) {
	if showService == nil {
		InitializeService()
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	stats, err := showService.GetSearchStatistics()
	if err != nil {
		log.Printf("Error getting search statistics: %v", err)
		http.Error(w, "Failed to get statistics", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

// HealthCheckHandler provides system health information
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": r.Context().Value("timestamp"),
		"services": map[string]string{
			"database": "connected",
			"redis":    "connected",
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(health)
}

// Helper function to parse search parameters from query string
func parseSearchParams(r *http.Request) service.SearchRequest {
	req := service.SearchRequest{}
	
	req.Location = r.URL.Query().Get("location")
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
