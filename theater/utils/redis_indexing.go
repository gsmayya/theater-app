package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

const (
	// Redis key prefixes for different indexes
	ShowsByLocationPrefix     = "shows:location:"
	ShowsByPricePrefix        = "shows:price"
	ShowsByAvailabilityPrefix = "shows:availability"
	ShowsSearchPrefix         = "shows:search:"
	ShowsAllKey               = "shows:all"
)

// IndexedRedisClient extends the basic Redis functionality with indexing
type IndexedRedisClient struct {
	*RedisAccess
}

// NewIndexedRedisClient creates a new indexed Redis client
func NewIndexedRedisClient() *IndexedRedisClient {
	return &IndexedRedisClient{
		RedisAccess: GetStoreAccess(),
	}
}

// ShowIndexData represents the data structure for show indexing
type ShowIndexData struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Location         string `json:"location"`
	Price            int32  `json:"price"`
	AvailableTickets int32  `json:"available_tickets"`
	TotalTickets     int32  `json:"total_tickets"`
	Details          string `json:"details"`
}

// IndexShow adds a show to various Redis indexes for fast searching
func (irc *IndexedRedisClient) IndexShow(show ShowIndexData) error {
	ctx := *irc.context
	pipe := irc.client.Pipeline()

	// Index by location (set of show IDs for each location)
	locationKey := ShowsByLocationPrefix + strings.ToLower(show.Location)
	pipe.SAdd(ctx, locationKey, show.ID)

	// Index by price (sorted set with price as score)
	pipe.ZAdd(ctx, ShowsByPricePrefix, redis.Z{
		Score:  float64(show.Price),
		Member: show.ID,
	})

	// Index by availability (sorted set with available tickets as score)
	pipe.ZAdd(ctx, ShowsByAvailabilityPrefix, redis.Z{
		Score:  float64(show.AvailableTickets),
		Member: show.ID,
	})

	// Add to all shows set
	pipe.SAdd(ctx, ShowsAllKey, show.ID)

	// Index searchable terms (for simple text search)
	searchTerms := extractSearchTerms(show.Name + " " + show.Details)
	for _, term := range searchTerms {
		searchKey := ShowsSearchPrefix + strings.ToLower(term)
		pipe.SAdd(ctx, searchKey, show.ID)
	}

	// Store the show data as a hash for quick retrieval
	showHashKey := "show:" + show.ID
	showData, err := json.Marshal(show)
	if err != nil {
		return fmt.Errorf("failed to marshal show data: %w", err)
	}

	pipe.Set(ctx, showHashKey, showData, 0)

	// Execute all commands
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Printf("Error indexing show %s: %v", show.ID, err)
		return err
	}

	log.Printf("Successfully indexed show: %s (%s)", show.Name, show.ID)
	return nil
}

// RemoveShowFromIndexes removes a show from all indexes
func (irc *IndexedRedisClient) RemoveShowFromIndexes(showID string, location string) error {
	ctx := *irc.context
	pipe := irc.client.Pipeline()

	// Remove from location index
	locationKey := ShowsByLocationPrefix + strings.ToLower(location)
	pipe.SRem(ctx, locationKey, showID)

	// Remove from price index
	pipe.ZRem(ctx, ShowsByPricePrefix, showID)

	// Remove from availability index
	pipe.ZRem(ctx, ShowsByAvailabilityPrefix, showID)

	// Remove from all shows set
	pipe.SRem(ctx, ShowsAllKey, showID)

	// Remove show data
	showHashKey := "show:" + showID
	pipe.Del(ctx, showHashKey)

	// Note: Removing from search indexes would require knowing all terms,
	// which is expensive. In practice, these can be cleaned up periodically.

	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Printf("Error removing show %s from indexes: %v", showID, err)
		return err
	}

	return nil
}

// SearchShowsByLocation retrieves show IDs by location using Redis sets
func (irc *IndexedRedisClient) SearchShowsByLocation(location string) ([]string, error) {
	ctx := *irc.context
	locationKey := ShowsByLocationPrefix + strings.ToLower(location)

	showIDs, err := irc.client.SMembers(ctx, locationKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get shows by location: %w", err)
	}

	return showIDs, nil
}

// SearchShowsByPriceRange retrieves show IDs within a price range using sorted sets
func (irc *IndexedRedisClient) SearchShowsByPriceRange(minPrice, maxPrice int32) ([]string, error) {
	ctx := *irc.context

	result, err := irc.client.ZRangeByScore(ctx, ShowsByPricePrefix, &redis.ZRangeBy{
		Min: strconv.Itoa(int(minPrice)),
		Max: strconv.Itoa(int(maxPrice)),
	}).Result()

	if err != nil {
		return nil, fmt.Errorf("failed to get shows by price range: %w", err)
	}

	return result, nil
}

// SearchShowsByAvailability retrieves shows with minimum available tickets
func (irc *IndexedRedisClient) SearchShowsByAvailability(minAvailable int32) ([]string, error) {
	ctx := *irc.context

	result, err := irc.client.ZRangeByScore(ctx, ShowsByAvailabilityPrefix, &redis.ZRangeBy{
		Min: strconv.Itoa(int(minAvailable)),
		Max: "+inf",
	}).Result()

	if err != nil {
		return nil, fmt.Errorf("failed to get shows by availability: %w", err)
	}

	return result, nil
}

// SearchShowsByTerm performs simple text search using Redis sets intersection
func (irc *IndexedRedisClient) SearchShowsByTerm(searchTerm string) ([]string, error) {
	ctx := *irc.context
	terms := extractSearchTerms(searchTerm)

	if len(terms) == 0 {
		return []string{}, nil
	}

	// Create keys for each search term
	var keys []string
	for _, term := range terms {
		keys = append(keys, ShowsSearchPrefix+strings.ToLower(term))
	}

	// If only one term, just get members
	if len(keys) == 1 {
		return irc.client.SMembers(ctx, keys[0]).Result()
	}

	// For multiple terms, use intersection
	tempKey := "temp:search:" + fmt.Sprintf("%d", len(terms))
	err := irc.client.SInterStore(ctx, tempKey, keys...).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to intersect search terms: %w", err)
	}

	// Get results and clean up temp key
	result, err := irc.client.SMembers(ctx, tempKey).Result()
	irc.client.Del(ctx, tempKey) // Clean up temp key

	if err != nil {
		return nil, fmt.Errorf("failed to get search results: %w", err)
	}

	return result, nil
}

// GetShowsByIDs retrieves show data for multiple show IDs efficiently
func (irc *IndexedRedisClient) GetShowsByIDs(showIDs []string) ([]ShowIndexData, error) {
	if len(showIDs) == 0 {
		return []ShowIndexData{}, nil
	}

	ctx := *irc.context
	pipe := irc.client.Pipeline()

	// Pipeline all get commands
	var cmds []*redis.StringCmd
	for _, showID := range showIDs {
		showKey := "show:" + showID
		cmds = append(cmds, pipe.Get(ctx, showKey))
	}

	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, fmt.Errorf("failed to get shows by IDs: %w", err)
	}

	// Process results
	var shows []ShowIndexData
	for i, cmd := range cmds {
		result, err := cmd.Result()
		if err == redis.Nil {
			log.Printf("Show not found in cache: %s", showIDs[i])
			continue
		} else if err != nil {
			log.Printf("Error getting show %s: %v", showIDs[i], err)
			continue
		}

		var show ShowIndexData
		if err := json.Unmarshal([]byte(result), &show); err != nil {
			log.Printf("Error unmarshaling show %s: %v", showIDs[i], err)
			continue
		}

		shows = append(shows, show)
	}

	return shows, nil
}

// CombinedSearch performs a complex search combining multiple criteria
func (irc *IndexedRedisClient) CombinedSearch(location string, minPrice, maxPrice int32, minAvailable int32, searchTerm string) ([]string, error) {
	ctx := *irc.context
	var keys []string
	tempKeys := []string{}

	// Collect all criteria keys
	if location != "" {
		keys = append(keys, ShowsByLocationPrefix+strings.ToLower(location))
	}

	// Handle price range
	if minPrice > 0 || maxPrice > 0 {
		priceKey := "temp:price"
		actualMinPrice := minPrice
		actualMaxPrice := maxPrice
		if actualMaxPrice == 0 {
			actualMaxPrice = 999999 // Large number for no upper limit
		}

		err := irc.client.ZRangeByScore(ctx, priceKey, &redis.ZRangeBy{
			Min: strconv.Itoa(int(actualMinPrice)),
			Max: strconv.Itoa(int(actualMaxPrice)),
		}).Err()

		if err != nil {
			return nil, fmt.Errorf("failed to filter by price range: %w", err)
		}

		keys = append(keys, priceKey)
		tempKeys = append(tempKeys, priceKey)
	}

	// Handle availability
	if minAvailable > 0 {
		availKey := "temp:availability"
		err := irc.client.ZRangeByScore(ctx, availKey, &redis.ZRangeBy{
			Min: strconv.Itoa(int(minAvailable)),
			Max: "+inf",
		}).Err()

		if err != nil {
			return nil, fmt.Errorf("failed to filter by availability: %w", err)
		}

		keys = append(keys, availKey)
		tempKeys = append(tempKeys, availKey)
	}

	// Handle search terms
	if searchTerm != "" {
		terms := extractSearchTerms(searchTerm)
		for _, term := range terms {
			keys = append(keys, ShowsSearchPrefix+strings.ToLower(term))
		}
	}

	// If no criteria, return all shows
	if len(keys) == 0 {
		result, err := irc.client.SMembers(ctx, ShowsAllKey).Result()
		return result, err
	}

	// Perform intersection of all criteria
	resultKey := "temp:final_result"
	tempKeys = append(tempKeys, resultKey)

	err := irc.client.SInterStore(ctx, resultKey, keys...).Err()
	if err != nil {
		// Clean up temp keys
		if len(tempKeys) > 0 {
			irc.client.Del(ctx, tempKeys...)
		}
		return nil, fmt.Errorf("failed to perform combined search: %w", err)
	}

	// Get final results
	result, err := irc.client.SMembers(ctx, resultKey).Result()

	// Clean up all temp keys
	if len(tempKeys) > 0 {
		irc.client.Del(ctx, tempKeys...)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get combined search results: %w", err)
	}

	return result, nil
}

// Helper function to extract search terms from text
func extractSearchTerms(text string) []string {
	// Simple tokenization - split by spaces and remove short words
	words := strings.Fields(strings.ToLower(text))
	var terms []string

	for _, word := range words {
		// Remove punctuation and filter short words
		cleaned := strings.Trim(word, ".,!?;:")
		if len(cleaned) >= 3 { // Only index words with 3+ characters
			terms = append(terms, cleaned)
		}
	}

	return terms
}

// UpdateShowAvailability efficiently updates just the availability index
func (irc *IndexedRedisClient) UpdateShowAvailability(showID string, availableTickets int32) error {
	ctx := *irc.context

	// Update availability index
	err := irc.client.ZAdd(ctx, ShowsByAvailabilityPrefix, redis.Z{
		Score:  float64(availableTickets),
		Member: showID,
	}).Err()

	if err != nil {
		return fmt.Errorf("failed to update show availability: %w", err)
	}

	return nil
}

// GetShowStatistics returns statistics about indexed shows
func (irc *IndexedRedisClient) GetShowStatistics() (map[string]interface{}, error) {
	ctx := *irc.context
	pipe := irc.client.Pipeline()

	totalShowsCmd := pipe.SCard(ctx, ShowsAllKey)
	priceStatsCmd := pipe.ZCard(ctx, ShowsByPricePrefix)
	availabilityStatsCmd := pipe.ZCard(ctx, ShowsByAvailabilityPrefix)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}

	stats := map[string]interface{}{
		"total_shows":                 totalShowsCmd.Val(),
		"shows_in_price_index":        priceStatsCmd.Val(),
		"shows_in_availability_index": availabilityStatsCmd.Val(),
	}

	// Get location statistics
	locationPattern := ShowsByLocationPrefix + "*"
	locationKeys, err := irc.client.Keys(ctx, locationPattern).Result()
	if err == nil {
		locationCounts := make(map[string]int64)
		for _, key := range locationKeys {
			location := strings.TrimPrefix(key, ShowsByLocationPrefix)
			count := irc.client.SCard(ctx, key).Val()
			locationCounts[location] = count
		}
		stats["shows_by_location"] = locationCounts
	}

	return stats, nil
}
