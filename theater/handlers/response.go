package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// APIResponse represents a standardized API response structure
type APIResponse struct {
	Success   bool                   `json:"success"`
	Message   string                 `json:"message,omitempty"`
	Data      interface{}            `json:"data,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Meta      map[string]interface{} `json:"meta,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	APIResponse
	Pagination PaginationInfo `json:"pagination"`
}

// PaginationInfo contains pagination metadata
type PaginationInfo struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// WriteSuccessResponse writes a successful response
func WriteSuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	response := APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}
	writeJSONResponse(w, statusCode, response)
}

// WriteErrorResponse writes an error response
func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string, err error) {
	response := APIResponse{
		Success:   false,
		Message:   message,
		Error:     err.Error(),
		Timestamp: time.Now(),
	}
	writeJSONResponse(w, statusCode, response)
}

// WritePaginatedResponse writes a paginated response
func WritePaginatedResponse(w http.ResponseWriter, statusCode int, message string, data interface{}, pagination PaginationInfo) {
	response := PaginatedResponse{
		APIResponse: APIResponse{
			Success:   true,
			Message:   message,
			Data:      data,
			Timestamp: time.Now(),
		},
		Pagination: pagination,
	}
	writeJSONResponse(w, statusCode, response)
}

// writeJSONResponse writes a JSON response with proper headers
func writeJSONResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// If JSON encoding fails, write a simple error response
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success":false,"error":"Failed to encode response"}`))
	}
}

// HandleCORS handles preflight OPTIONS requests
func HandleCORS(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return true
	}
	return false
}

// ValidateMethod validates that the request method is allowed
func ValidateMethod(w http.ResponseWriter, r *http.Request, allowedMethods ...string) bool {
	for _, method := range allowedMethods {
		if r.Method == method {
			return true
		}
	}

	WriteErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed",
		&HTTPError{Code: http.StatusMethodNotAllowed, Message: "Method not allowed"})
	return false
}

// HTTPError represents an HTTP error
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

// Common error types
var (
	ErrInvalidRequest  = &HTTPError{Code: http.StatusBadRequest, Message: "Invalid request"}
	ErrNotFound        = &HTTPError{Code: http.StatusNotFound, Message: "Resource not found"}
	ErrInternalServer  = &HTTPError{Code: http.StatusInternalServerError, Message: "Internal server error"}
	ErrUnauthorized    = &HTTPError{Code: http.StatusUnauthorized, Message: "Unauthorized"}
	ErrForbidden       = &HTTPError{Code: http.StatusForbidden, Message: "Forbidden"}
	ErrConflict        = &HTTPError{Code: http.StatusConflict, Message: "Conflict"}
	ErrTooManyRequests = &HTTPError{Code: http.StatusTooManyRequests, Message: "Too many requests"}
)
