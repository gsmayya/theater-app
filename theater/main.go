package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gsmayya/theater/db"
	"github.com/gsmayya/theater/handlers"
)

const (
	DefaultPort  = "8080"
	ReadTimeout  = 15 * time.Second
	WriteTimeout = 15 * time.Second
	IdleTimeout  = 60 * time.Second
)

func main() {
	log.Println("🎭 Starting Theater Booking System...")

	// Initialize database connection
	database := db.GetDatabase()
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Test database connection
	if err := database.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("✅ Database connection established successfully")

	// Initialize services
	handlers.InitializeService()
	handlers.InitializeBookingService()
	log.Println("✅ Services initialized successfully")

	// Setup routes
	router := setupRoutes()

	// Create HTTP server with timeouts
	server := &http.Server{
		Addr:         "0.0.0.0:" + getPort(),
		Handler:      router,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
		IdleTimeout:  IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Server starting on port %s", getPort())
		logRoutes()
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited")
}

func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Legacy endpoints (for backward compatibility)
	mux.HandleFunc("/shows", handlers.ShowListHandler)
	mux.HandleFunc("/show", handlers.ShowHandler)
	mux.HandleFunc("/status", handlers.HealthCheckHandler)

	// API v1 routes
	apiV1 := "/api/v1"

	// Show management endpoints
	mux.HandleFunc(apiV1+"/search", handlers.SearchShowsHandler)
	mux.HandleFunc(apiV1+"/shows", handlers.ShowsByAllHandler)
	mux.HandleFunc(apiV1+"/shows/by-location", handlers.ShowsByLocationHandler)
	mux.HandleFunc(apiV1+"/shows/by-price-range", handlers.ShowsByPriceRangeHandler)
	mux.HandleFunc(apiV1+"/shows/create", handlers.CreateShowHandler)
	mux.HandleFunc(apiV1+"/shows/get", handlers.GetShowHandler)
	mux.HandleFunc(apiV1+"/shows/update-availability", handlers.UpdateShowAvailabilityHandler)
	mux.HandleFunc(apiV1+"/shows/booking-summary", handlers.GetShowBookingSummaryHandler)

	// Booking management endpoints
	mux.HandleFunc(apiV1+"/bookings/create", handlers.CreateBookingHandler)
	mux.HandleFunc(apiV1+"/bookings/get", handlers.GetBookingHandler)
	mux.HandleFunc(apiV1+"/bookings/update-status", handlers.UpdateBookingStatusHandler)
	mux.HandleFunc(apiV1+"/bookings/confirm", handlers.ConfirmBookingHandler)
	mux.HandleFunc(apiV1+"/bookings/cancel", handlers.CancelBookingHandler)
	mux.HandleFunc(apiV1+"/bookings/by-show", handlers.GetBookingsByShowHandler)
	mux.HandleFunc(apiV1+"/bookings/by-contact", handlers.GetBookingsByContactHandler)
	mux.HandleFunc(apiV1+"/bookings/search", handlers.SearchBookingsHandler)
	mux.HandleFunc(apiV1+"/bookings/stats", handlers.GetBookingStatsHandler)

	// System endpoints
	mux.HandleFunc(apiV1+"/stats", handlers.GetSearchStatsHandler)
	mux.HandleFunc(apiV1+"/health", handlers.HealthCheckHandler)

	return mux
}

func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return DefaultPort
}

func logRoutes() {
	log.Println("📋 Available endpoints:")
	log.Println("")
	log.Println("  🔄 Legacy endpoints:")
	log.Println("    GET  /shows                    - Get all shows")
	log.Println("    GET  /show?show_id=<id>        - Get show details")
	log.Println("    PUT  /show                     - Create new show")
	log.Println("    GET  /status                   - Health check")
	log.Println("")
	log.Println("  🎪 Show management (API v1):")
	log.Println("    GET  /api/v1/search            - Advanced show search")
	log.Println("    GET  /api/v1/shows             - Get all shows")
	log.Println("    GET  /api/v1/shows/by-location - Shows by location")
	log.Println("    GET  /api/v1/shows/by-price-range - Shows by price range")
	log.Println("    POST /api/v1/shows/create      - Create new show")
	log.Println("    GET  /api/v1/shows/get         - Get show details")
	log.Println("    PUT  /api/v1/shows/update-availability - Update show availability")
	log.Println("    GET  /api/v1/shows/booking-summary - Show booking summary")
	log.Println("")
	log.Println("  🎟️ Booking management (API v1):")
	log.Println("    POST /api/v1/bookings/create   - Create new booking")
	log.Println("    GET  /api/v1/bookings/get      - Get booking details")
	log.Println("    PUT  /api/v1/bookings/update-status - Update booking status")
	log.Println("    PUT  /api/v1/bookings/confirm  - Confirm booking")
	log.Println("    PUT  /api/v1/bookings/cancel   - Cancel booking")
	log.Println("    GET  /api/v1/bookings/by-show  - Get bookings for a show")
	log.Println("    GET  /api/v1/bookings/by-contact - Get bookings by contact")
	log.Println("    GET  /api/v1/bookings/search   - Search bookings")
	log.Println("    GET  /api/v1/bookings/stats    - Booking statistics")
	log.Println("")
	log.Println("  📊 System endpoints (API v1):")
	log.Println("    GET  /api/v1/stats             - Search statistics")
	log.Println("    GET  /api/v1/health            - Health check")
}
