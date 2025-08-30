package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gsmayya/theater/db"
	"github.com/gsmayya/theater/handlers"
)

func main() {
	fmt.Println("Starting Optimized Theater Booking Service...")
	
	// Initialize database connection
	database := db.GetDatabase()
	defer database.Close()
	
	// Test database connection
	if err := database.Ping(); err != nil {
		log.Printf("Warning: Database connection failed: %v", err)
	} else {
		log.Println("Database connection established successfully")
	}
	
	// Initialize the service layer
	handlers.InitializeService()
	
	// Original endpoints (backward compatibility)
	http.HandleFunc("/", handlers.DefaultHandler)
	http.HandleFunc("/status", handlers.DefaultHandler)
	http.HandleFunc("/shows", handlers.ShowListHandler)
	http.HandleFunc("/show", handlers.ShowHandler)
	
	// New optimized endpoints with indexing
	http.HandleFunc("/api/v1/search", handlers.SearchShowsHandler)
	http.HandleFunc("/api/v1/shows/by-location", handlers.ShowsByLocationHandler)
	http.HandleFunc("/api/v1/shows/by-price-range", handlers.ShowsByPriceRangeHandler)
	http.HandleFunc("/api/v1/shows/create", handlers.CreateShowHandler)
	http.HandleFunc("/api/v1/shows/get", handlers.GetShowHandler)
	http.HandleFunc("/api/v1/shows/update-availability", handlers.UpdateShowAvailabilityHandler)
	http.HandleFunc("/api/v1/stats", handlers.GetSearchStatsHandler)
	http.HandleFunc("/api/v1/health", handlers.HealthCheckHandler)
	
	log.Println("Server starting at :8080 with optimized indexing")
	log.Println("Available endpoints:")
	log.Println("  Legacy: /shows, /show, /status")
	log.Println("  Search: /api/v1/search")
	log.Println("  Location: /api/v1/shows/by-location")
	log.Println("  Price: /api/v1/shows/by-price-range")
	log.Println("  CRUD: /api/v1/shows/create, /api/v1/shows/get")
	log.Println("  Stats: /api/v1/stats, /api/v1/health")
	
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
