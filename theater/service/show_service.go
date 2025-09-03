package service

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/gsmayya/theater/repository"
	"github.com/gsmayya/theater/shows"
	"github.com/gsmayya/theater/utils"
)

// ShowService provides business logic for theater shows with optimized caching
type ShowService struct {
	repository *repository.ShowRepository
	redisIndex *utils.IndexedRedisClient
}

// SearchRequest represents a search query with all possible filters
type SearchRequest struct {
	ShowLocation  string `json:"show_location,omitempty"`
	MinPrice      *int32 `json:"min_price,omitempty"`
	MaxPrice      *int32 `json:"max_price,omitempty"`
	MinAvailable  *int32 `json:"min_available,omitempty"`
	SearchTerm    string `json:"search_term,omitempty"`
	OnlyAvailable bool   `json:"only_available,omitempty"`
	Page          int    `json:"page,omitempty"`
	PageSize      int    `json:"page_size,omitempty"`
}

// SearchResponse represents the response from a search query
type SearchResponse struct {
	Shows      []*shows.ShowData `json:"shows"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

// NewShowService creates a new show service with optimized caching
func NewShowService() *ShowService {
	return &ShowService{
		repository: repository.NewShowRepository(),
		redisIndex: utils.NewIndexedRedisClient(),
	}
}

// CreateShow creates a new show with full indexing
func (s *ShowService) CreateShow(show_name, details, show_location string, price, totalTickets int32) (*shows.ShowData, error) {
	// Create show data
	show := &shows.ShowData{}
	show.NewShow(show_name, details, price, totalTickets, show_location)

	// Save to database
	if err := s.repository.CreateShow(show); err != nil {
		return nil, fmt.Errorf("failed to create show in database: %w", err)
	}

	// Index in Redis for fast searches
	indexData := utils.ShowIndexData{
		ID:               show.Show_Id.String(),
		ShowName:         show.ShowName,
		ShowLocation:     show.ShowLocation,
		Price:            show.Price,
		AvailableTickets: show.Total_Tickets - show.Booked_Tickets,
		TotalTickets:     show.Total_Tickets,
		Details:          show.Details,
	}

	if err := s.redisIndex.IndexShow(indexData); err != nil {
		log.Printf("Warning: Failed to index show in Redis: %v", err)
		// Don't fail the entire operation for indexing errors
	}

	log.Printf("Successfully created show: %s (ID: %s)", show.ShowName, show.Show_Id.String())
	return show, nil
}

// GetShow retrieves a show by ID using cache-first strategy
func (s *ShowService) GetShow(showID string) (*shows.ShowData, error) {
	// Validate UUID format
	if _, err := uuid.Parse(showID); err != nil {
		return nil, fmt.Errorf("invalid show ID format: %s", showID)
	}

	// Try repository (which checks cache first, then database)
	show, err := s.repository.GetShow(showID)
	if err != nil {
		return nil, err
	}

	return show, nil
}

// SearchShows performs optimized searching using multiple strategies
func (s *ShowService) SearchShows(req SearchRequest) (*SearchResponse, error) {
	// Set default pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100 // Limit max page size
	}

	// Strategy 1: Use Redis indexes for fast filtering if we have specific criteria
	if s.shouldUseRedisIndex(req) {
		return s.searchWithRedisIndex(req)
	}

	// Strategy 2: Use database queries with indexes for complex searches
	return s.searchWithDatabase(req)
}

func (s *ShowService) GetAllShows() ([]*shows.ShowData, error) {
	return s.GetShowsByPriceRange(0, 100000, "")
}

// GetShowsByLocation uses Redis indexing for location-based searches
func (s *ShowService) GetShowsByLocation(show_location string, onlyAvailable bool) ([]*shows.ShowData, error) {
	if show_location == "" {
		return nil, fmt.Errorf("location cannot be empty")
	}

	// Try Redis first for fast results
	showIDs, err := s.redisIndex.SearchShowsByLocation(show_location)
	if err == nil && len(showIDs) > 0 {
		// Get show data from Redis
		indexedShows, err := s.redisIndex.GetShowsByIDs(showIDs)
		if err == nil {
			// Convert to ShowData format
			var result []*shows.ShowData
			for _, indexed := range indexedShows {
				if onlyAvailable && indexed.AvailableTickets <= 0 {
					continue
				}
				show := s.convertFromIndexed(indexed)
				result = append(result, show)
			}
			return result, nil
		}
	}

	// Fall back to database query
	return s.repository.GetShowsByLocation(show_location, onlyAvailable)
}

// GetShowsByPriceRange uses Redis sorted sets for efficient price range queries
func (s *ShowService) GetShowsByPriceRange(minPrice, maxPrice int32, show_location string) ([]*shows.ShowData, error) {
	if minPrice < 0 || maxPrice < 0 {
		return nil, fmt.Errorf("price values cannot be negative")
	}
	if minPrice > maxPrice {
		return nil, fmt.Errorf("minimum price cannot be greater than maximum price")
	}

	// Try Redis first
	showIDs, err := s.redisIndex.SearchShowsByPriceRange(minPrice, maxPrice)
	if err == nil {
		// Filter by location if specified
		if show_location != "" {
			locationIDs, err := s.redisIndex.SearchShowsByLocation(show_location)
			if err == nil {
				showIDs = intersectSlices(showIDs, locationIDs)
			}
		}

		if len(showIDs) > 0 {
			indexedShows, err := s.redisIndex.GetShowsByIDs(showIDs)
			if err == nil {
				var result []*shows.ShowData
				for _, indexed := range indexedShows {
					result = append(result, s.convertFromIndexed(indexed))
				}
				return result, nil
			}
		}
	}

	// Fall back to database
	return s.repository.GetShowsByPriceRange(minPrice, maxPrice, show_location)
}

// UpdateShow updates a show and maintains all indexes
func (s *ShowService) UpdateShow(show *shows.ShowData) error {
	// Update in database
	if err := s.repository.UpdateShow(show); err != nil {
		return err
	}

	// Update Redis indexes
	indexData := utils.ShowIndexData{
		ID:               show.Show_Id.String(),
		ShowName:         show.ShowName,
		ShowLocation:     show.ShowLocation,
		Price:            show.Price,
		AvailableTickets: show.Total_Tickets - show.Booked_Tickets,
		TotalTickets:     show.Total_Tickets,
		Details:          show.Details,
	}

	if err := s.redisIndex.IndexShow(indexData); err != nil {
		log.Printf("Warning: Failed to update show in Redis index: %v", err)
	}

	return nil
}

// UpdateTicketAvailability efficiently updates just the availability information
func (s *ShowService) UpdateTicketAvailability(showID string, bookedTickets int32) error {
	// Get current show data
	show, err := s.GetShow(showID)
	if err != nil {
		return err
	}

	// Update booked tickets
	show.Booked_Tickets = bookedTickets

	// Update in database
	if err := s.repository.UpdateShow(show); err != nil {
		return err
	}

	// Update Redis availability index efficiently
	availableTickets := show.Total_Tickets - show.Booked_Tickets
	if err := s.redisIndex.UpdateShowAvailability(showID, availableTickets); err != nil {
		log.Printf("Warning: Failed to update availability in Redis: %v", err)
	}

	return nil
}

// DeleteShow removes a show from all systems
func (s *ShowService) DeleteShow(showID string) error {
	// Get show data for cleanup
	show, err := s.GetShow(showID)
	if err != nil {
		return err
	}

	// Delete from database
	if err := s.repository.DeleteShow(showID); err != nil {
		return err
	}

	// Remove from Redis indexes
	if err := s.redisIndex.RemoveShowFromIndexes(showID, show.ShowLocation); err != nil {
		log.Printf("Warning: Failed to remove show from Redis indexes: %v", err)
	}

	return nil
}

// GetSearchStatistics returns statistics about the search indexes
func (s *ShowService) GetSearchStatistics() (map[string]interface{}, error) {
	return s.redisIndex.GetShowStatistics()
}

// Helper methods

func (s *ShowService) shouldUseRedisIndex(req SearchRequest) bool {
	// Use Redis when we have specific criteria that benefit from indexing
	return req.ShowLocation != "" ||
		req.MinPrice != nil ||
		req.MaxPrice != nil ||
		req.MinAvailable != nil ||
		req.SearchTerm != ""
}

func (s *ShowService) searchWithRedisIndex(req SearchRequest) (*SearchResponse, error) {
	var minPrice, maxPrice, minAvailable int32

	if req.MinPrice != nil {
		minPrice = *req.MinPrice
	}
	if req.MaxPrice != nil {
		maxPrice = *req.MaxPrice
	}
	if req.MinAvailable != nil {
		minAvailable = *req.MinAvailable
	}

	// Perform combined search using Redis
	showIDs, err := s.redisIndex.CombinedSearch(req.ShowLocation, minPrice, maxPrice, minAvailable, req.SearchTerm)
	if err != nil {
		log.Printf("Redis search failed, falling back to database: %v", err)
		return s.searchWithDatabase(req)
	}

	// Get show data
	indexedShows, err := s.redisIndex.GetShowsByIDs(showIDs)
	if err != nil {
		return s.searchWithDatabase(req)
	}

	// Apply availability filter if needed
	var filteredShows []utils.ShowIndexData
	for _, show := range indexedShows {
		if req.OnlyAvailable && show.AvailableTickets <= 0 {
			continue
		}
		filteredShows = append(filteredShows, show)
	}

	// Apply pagination
	total := len(filteredShows)
	offset := (req.Page - 1) * req.PageSize
	end := offset + req.PageSize

	if offset >= total {
		return &SearchResponse{
			Shows:      []*shows.ShowData{},
			Total:      total,
			Page:       req.Page,
			PageSize:   req.PageSize,
			TotalPages: (total + req.PageSize - 1) / req.PageSize,
		}, nil
	}

	if end > total {
		end = total
	}

	// Convert to ShowData
	var result []*shows.ShowData
	for i := offset; i < end; i++ {
		result = append(result, s.convertFromIndexed(filteredShows[i]))
	}

	return &SearchResponse{
		Shows:      result,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: (total + req.PageSize - 1) / req.PageSize,
	}, nil
}

func (s *ShowService) searchWithDatabase(req SearchRequest) (*SearchResponse, error) {
	// Convert request to repository filters
	filters := &repository.SearchFilters{
		ShowLocation:  req.ShowLocation,
		MinPrice:      req.MinPrice,
		MaxPrice:      req.MaxPrice,
		MinAvailable:  req.MinAvailable,
		SearchTerm:    req.SearchTerm,
		OnlyAvailable: req.OnlyAvailable,
	}

	pagination := &repository.PaginationParams{
		Offset: (req.Page - 1) * req.PageSize,
		Limit:  req.PageSize,
	}

	showsList, total, err := s.repository.GetAllShows(filters, pagination)
	if err != nil {
		return nil, err
	}

	return &SearchResponse{
		Shows:      showsList,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: (total + req.PageSize - 1) / req.PageSize,
	}, nil
}

func (s *ShowService) convertFromIndexed(indexed utils.ShowIndexData) *shows.ShowData {
	showID, _ := uuid.Parse(indexed.ID)
	return &shows.ShowData{
		Show_Id:        showID,
		ShowName:       indexed.ShowName,
		Details:        indexed.Details,
		Price:          indexed.Price,
		Total_Tickets:  indexed.TotalTickets,
		Booked_Tickets: indexed.TotalTickets - indexed.AvailableTickets,
		ShowLocation:   indexed.ShowLocation,
	}
}

// Utility function to find intersection of two string slices
func intersectSlices(a, b []string) []string {
	set := make(map[string]bool)
	var result []string

	// Add all elements from first slice to set
	for _, item := range a {
		set[item] = true
	}

	// Check elements from second slice
	for _, item := range b {
		if set[item] {
			result = append(result, item)
			delete(set, item) // Avoid duplicates
		}
	}

	return result
}
