package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gsmayya/theater/db"
	"github.com/gsmayya/theater/shows"
	"github.com/gsmayya/theater/utils"
)

type ShowRepository struct {
	database    *db.Database
	redisClient *utils.RedisAccess
}

type SearchFilters struct {
	Location      string
	MinPrice      *int32
	MaxPrice      *int32
	MinAvailable  *int32
	SearchTerm    string
	OnlyAvailable bool
}

type PaginationParams struct {
	Offset int
	Limit  int
}

func NewShowRepository() *ShowRepository {
	return &ShowRepository{
		database:    db.GetDatabase(),
		redisClient: utils.GetStoreAccess(),
	}
}

// CreateShow inserts a new show into the database
func (r *ShowRepository) CreateShow(show *shows.ShowData) error {
	imagesJSON, _ := json.Marshal(show.Images)
	videosJSON, _ := json.Marshal(show.Videos)
	
	query := `
		INSERT INTO shows (id, name, details, price, total_tickets, booked_tickets, location, show_number, show_date, images, videos)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.database.GetDB().Exec(query,
		show.Show_Id.String(),
		show.Name,
		show.Details,
		show.Price,
		show.Total_Tickets,
		show.Booked_Tickets,
		show.Location,
		show.ShowNumber,
		show.ShowDate,
		string(imagesJSON),
		string(videosJSON),
	)

	if err != nil {
		return fmt.Errorf("failed to create show: %w", err)
	}

	// Cache the show data
	r.cacheShow(show)

	log.Printf("Show created successfully: %s", show.Name)
	return nil
}

// GetShow retrieves a show by ID, first checking cache, then database
func (r *ShowRepository) GetShow(showID string) (*shows.ShowData, error) {
	// Try to get from cache first
	if cachedShow, err := r.getShowFromCache(showID); err == nil {
		return cachedShow, nil
	}

	// If not in cache, get from database
	query := `
		SELECT id, name, details, price, total_tickets, booked_tickets, location, 
		       show_number, show_date, images, videos, created_at, updated_at
		FROM shows 
		WHERE id = ?
	`

	row := r.database.GetDB().QueryRow(query, showID)

	show := &shows.ShowData{}
	var createdAt, updatedAt time.Time
	var imagesJSON, videosJSON string

	err := row.Scan(
		&show.Show_Id,
		&show.Name,
		&show.Details,
		&show.Price,
		&show.Total_Tickets,
		&show.Booked_Tickets,
		&show.Location,
		&show.ShowNumber,
		&show.ShowDate,
		&imagesJSON,
		&videosJSON,
		&createdAt,
		&updatedAt,
	)

	if err == nil {
		// Parse JSON arrays for images and videos
		if imagesJSON != "" {
			json.Unmarshal([]byte(imagesJSON), &show.Images)
		}
		if videosJSON != "" {
			json.Unmarshal([]byte(videosJSON), &show.Videos)
		}
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("show not found: %s", showID)
		}
		return nil, fmt.Errorf("failed to get show: %w", err)
	}

	// Cache the show for future requests
	r.cacheShow(show)

	return show, nil
}

// GetAllShows retrieves all shows with optional filtering and pagination
func (r *ShowRepository) GetAllShows(filters *SearchFilters, pagination *PaginationParams) ([]*shows.ShowData, int, error) {
	// Build the WHERE clause dynamically based on filters
	whereConditions := []string{}
	args := []interface{}{}

	if filters != nil {
		if filters.Location != "" {
			whereConditions = append(whereConditions, "location = ?")
			args = append(args, filters.Location)
		}

		if filters.MinPrice != nil {
			whereConditions = append(whereConditions, "price >= ?")
			args = append(args, *filters.MinPrice)
		}

		if filters.MaxPrice != nil {
			whereConditions = append(whereConditions, "price <= ?")
			args = append(args, *filters.MaxPrice)
		}

		if filters.MinAvailable != nil {
			whereConditions = append(whereConditions, "(total_tickets - booked_tickets) >= ?")
			args = append(args, *filters.MinAvailable)
		}

		if filters.OnlyAvailable {
			whereConditions = append(whereConditions, "total_tickets > booked_tickets")
		}

		if filters.SearchTerm != "" {
			whereConditions = append(whereConditions, "MATCH(name, details) AGAINST(? IN NATURAL LANGUAGE MODE)")
			args = append(args, filters.SearchTerm)
		}
	}

	// Build the complete query
	baseQuery := "FROM shows"
	countQuery := "SELECT COUNT(*) " + baseQuery
	selectQuery := `
		SELECT id, name, details, price, total_tickets, booked_tickets, location, 
		       show_number, show_date, images, videos, created_at, updated_at 
		` + baseQuery

	if len(whereConditions) > 0 {
		whereClause := " WHERE " + strings.Join(whereConditions, " AND ")
		countQuery += whereClause
		selectQuery += whereClause
	}

	// Add ordering and pagination
	selectQuery += " ORDER BY created_at DESC"
	if pagination != nil {
		selectQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", pagination.Limit, pagination.Offset)
	}

	// Get total count first
	var totalCount int
	err := r.database.GetDB().QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Get the shows
	rows, err := r.database.GetDB().Query(selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query shows: %w", err)
	}
	defer rows.Close()

	var showsList []*shows.ShowData
	for rows.Next() {
		show := &shows.ShowData{}
		var createdAt, updatedAt time.Time
		var imagesJSON, videosJSON string

		err := rows.Scan(
			&show.Show_Id,
			&show.Name,
			&show.Details,
			&show.Price,
			&show.Total_Tickets,
			&show.Booked_Tickets,
			&show.Location,
			&show.ShowNumber,
			&show.ShowDate,
			&imagesJSON,
			&videosJSON,
			&createdAt,
			&updatedAt,
		)

		if err == nil {
			// Parse JSON arrays for images and videos
			if imagesJSON != "" {
				json.Unmarshal([]byte(imagesJSON), &show.Images)
			}
			if videosJSON != "" {
				json.Unmarshal([]byte(videosJSON), &show.Videos)
			}
		}

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan show: %w", err)
		}

		showsList = append(showsList, show)

		// Cache each show
		go r.cacheShow(show) // Cache asynchronously
	}

	return showsList, totalCount, nil
}

// GetShowsByLocation uses indexed query for location-based searches
func (r *ShowRepository) GetShowsByLocation(location string, onlyAvailable bool) ([]*shows.ShowData, error) {
	var query string
	args := []interface{}{location}

	if onlyAvailable {
		query = `
			SELECT s.id, s.name, s.details, s.price, s.total_tickets, s.booked_tickets, s.location, s.created_at, s.updated_at
			FROM shows s
			INNER JOIN show_availability_index sai ON s.id = sai.show_id
			WHERE s.location = ? AND sai.is_available = true
			ORDER BY s.name
		`
	} else {
		query = `
			SELECT id, name, details, price, total_tickets, booked_tickets, location, created_at, updated_at
			FROM shows 
			WHERE location = ?
			ORDER BY name
		`
	}

	return r.executeShowQuery(query, args...)
}

// GetShowsByPriceRange uses indexed query for price range searches
func (r *ShowRepository) GetShowsByPriceRange(minPrice, maxPrice int32, location string) ([]*shows.ShowData, error) {
	var query string
	var args []interface{}

	if location != "" {
		query = `
			SELECT id, name, details, price, total_tickets, booked_tickets, location, created_at, updated_at
			FROM shows 
			WHERE price BETWEEN ? AND ? AND location = ?
			ORDER BY price, name
		`
		args = []interface{}{minPrice, maxPrice, location}
	} else {
		query = `
			SELECT id, name, details, price, total_tickets, booked_tickets, location, created_at, updated_at
			FROM shows 
			WHERE price BETWEEN ? AND ?
			ORDER BY price, name
		`
		args = []interface{}{minPrice, maxPrice}
	}

	return r.executeShowQuery(query, args...)
}

// UpdateShow updates an existing show
func (r *ShowRepository) UpdateShow(show *shows.ShowData) error {
	imagesJSON, _ := json.Marshal(show.Images)
	videosJSON, _ := json.Marshal(show.Videos)
	
	query := `
		UPDATE shows 
		SET name = ?, details = ?, price = ?, total_tickets = ?, booked_tickets = ?, location = ?, 
		    show_number = ?, show_date = ?, images = ?, videos = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	result, err := r.database.GetDB().Exec(query,
		show.Name,
		show.Details,
		show.Price,
		show.Total_Tickets,
		show.Booked_Tickets,
		show.Location,
		show.ShowNumber,
		show.ShowDate,
		string(imagesJSON),
		string(videosJSON),
		show.Show_Id.String(),
	)

	if err != nil {
		return fmt.Errorf("failed to update show: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("show not found: %s", show.Show_Id.String())
	}

	// Update cache
	r.cacheShow(show)

	return nil
}

// DeleteShow deletes a show by ID
func (r *ShowRepository) DeleteShow(showID string) error {
	query := "DELETE FROM shows WHERE id = ?"

	result, err := r.database.GetDB().Exec(query, showID)
	if err != nil {
		return fmt.Errorf("failed to delete show: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("show not found: %s", showID)
	}

	// Remove from cache
	r.removeCachedShow(showID)

	return nil
}

// Helper methods for caching
func (r *ShowRepository) cacheShow(show *shows.ShowData) {
	if jsonData, err := show.ShowToJSON(); err == nil {
		utils.AddToCache(show.Show_Id.String(), jsonData, r.redisClient)
	}
}

func (r *ShowRepository) getShowFromCache(showID string) (*shows.ShowData, error) {
	jsonData, err := utils.GetFromCache(showID, r.redisClient)
	if err != nil {
		return nil, err
	}

	show := &shows.ShowData{}
	if _, err := show.JSONToShow(jsonData); err != nil {
		return nil, err
	}

	return show, nil
}

func (r *ShowRepository) removeCachedShow(showID string) {
	utils.DeleteFromCache(showID, r.redisClient)
}

func (r *ShowRepository) executeShowQuery(query string, args ...interface{}) ([]*shows.ShowData, error) {
	rows, err := r.database.GetDB().Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var showsList []*shows.ShowData
	for rows.Next() {
		show := &shows.ShowData{}
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&show.Show_Id,
			&show.Name,
			&show.Details,
			&show.Price,
			&show.Total_Tickets,
			&show.Booked_Tickets,
			&show.Location,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan show: %w", err)
		}

		showsList = append(showsList, show)
	}

	return showsList, nil
}
