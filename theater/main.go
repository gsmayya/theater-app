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
	handlers.InitializeBookingService()
	
	// Original endpoints (backward compatibility)
	http.HandleFunc("/", handlers.DefaultHandler)
	http.HandleFunc("/status", handlers.DefaultHandler)
	http.HandleFunc("/shows", handlers.ShowListHandler)
	http.HandleFunc("/show", handlers.ShowHandler)
	
	// Show management endpoints
	http.HandleFunc("/api/v1/search", handlers.SearchShowsHandler)
	http.HandleFunc("/api/v1/shows/by-location", handlers.ShowsByLocationHandler)
	http.HandleFunc("/api/v1/shows/by-price-range", handlers.ShowsByPriceRangeHandler)
	http.HandleFunc("/api/v1/shows/create", handlers.CreateShowHandler)
	http.HandleFunc("/api/v1/shows/get", handlers.GetShowHandler)
	http.HandleFunc("/api/v1/shows/update-availability", handlers.UpdateShowAvailabilityHandler)
	
	// Booking management endpoints
	http.HandleFunc("/api/v1/bookings/create", handlers.CreateBookingHandler)
	http.HandleFunc("/api/v1/bookings/get", handlers.GetBookingHandler)
	http.HandleFunc("/api/v1/bookings/update-status", handlers.UpdateBookingStatusHandler)
	http.HandleFunc("/api/v1/bookings/confirm", handlers.ConfirmBookingHandler)
	http.HandleFunc("/api/v1/bookings/cancel", handlers.CancelBookingHandler)
	http.HandleFunc("/api/v1/bookings/by-show", handlers.GetBookingsByShowHandler)
	http.HandleFunc("/api/v1/bookings/by-contact", handlers.GetBookingsByContactHandler)
	http.HandleFunc("/api/v1/bookings/search", handlers.SearchBookingsHandler)
	http.HandleFunc("/api/v1/bookings/stats", handlers.GetBookingStatsHandler)
	http.HandleFunc("/api/v1/shows/booking-summary", handlers.GetShowBookingSummaryHandler)
	
	// System endpoints
	http.HandleFunc("/api/v1/stats", handlers.GetSearchStatsHandler)
	http.HandleFunc("/api/v1/health", handlers.HealthCheckHandler)
	
	log.Println("🎭 Theater Booking System Server starting at :8080")
	log.Println("Available endpoints:")
	log.Println("")
	log.Println("  📋 Legacy endpoints:")
	log.Println("    /shows, /show, /status")
	log.Println("")
	log.Println("  🎪 Show management:")
	log.Println("    /api/v1/search - Advanced show search")
	log.Println("    /api/v1/shows/by-location - Shows by location")
	log.Println("    /api/v1/shows/by-price-range - Shows by price range")
	log.Println("    /api/v1/shows/create - Create new show")
	log.Println("    /api/v1/shows/get - Get show details")
	log.Println("    /api/v1/shows/update-availability - Update show availability")
	log.Println("    /api/v1/shows/booking-summary - Show booking summary")
	log.Println("")
	log.Println("  🎟️ Booking management:")
	log.Println("    /api/v1/bookings/create - Create new booking")
	log.Println("    /api/v1/bookings/get - Get booking details")
	log.Println("    /api/v1/bookings/update-status - Update booking status")
	log.Println("    /api/v1/bookings/confirm - Confirm booking")
	log.Println("    /api/v1/bookings/cancel - Cancel booking")
	log.Println("    /api/v1/bookings/by-show - Get bookings for a show")
	log.Println("    /api/v1/bookings/by-contact - Get bookings by contact")
	log.Println("    /api/v1/bookings/search - Search bookings")
	log.Println("    /api/v1/bookings/stats - Booking statistics")
	log.Println("")
	log.Println("  📊 System endpoints:")
	log.Println("    /api/v1/stats - Search statistics")
	log.Println("    /api/v1/health - Health check")
	
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
