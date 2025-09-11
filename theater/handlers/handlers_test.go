package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDefaultHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	DefaultHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response APIResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}

	if !strings.Contains(response.Message, "Ticket Backend Service is running") {
		t.Errorf("Expected message to contain 'Ticket Backend Service is running', got %s", response.Message)
	}
}

func TestHealthCheckHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	HealthCheckHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response APIResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}

	if response.Data == nil {
		t.Error("Expected data to be present")
	}

	// Check if data contains expected health information
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Error("Expected data to be a map")
	}

	if status, exists := data["status"]; !exists || status != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", status)
	}
}

func TestHandleCORS(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		expected bool
	}{
		{"OPTIONS request", "OPTIONS", true},
		{"GET request", "GET", false},
		{"POST request", "POST", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test", nil)
			w := httptest.NewRecorder()

			result := HandleCORS(w, req)

			if result != tt.expected {
				t.Errorf("Expected HandleCORS to return %v, got %v", tt.expected, result)
			}

			if tt.method == "OPTIONS" && w.Code != http.StatusOK {
				t.Errorf("Expected status %d for OPTIONS, got %d", http.StatusOK, w.Code)
			}
		})
	}
}

func TestValidateMethod(t *testing.T) {
	tests := []struct {
		name            string
		method          string
		allowedMethods  []string
		expected        bool
		expectedStatus  int
	}{
		{"Valid GET", "GET", []string{"GET"}, true, 0},
		{"Valid POST", "POST", []string{"POST"}, true, 0},
		{"Valid PUT", "PUT", []string{"PUT"}, true, 0},
		{"Invalid method", "DELETE", []string{"GET", "POST"}, false, http.StatusMethodNotAllowed},
		{"Multiple allowed methods", "POST", []string{"GET", "POST", "PUT"}, true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test", nil)
			w := httptest.NewRecorder()

			result := ValidateMethod(w, req, tt.allowedMethods...)

			if result != tt.expected {
				t.Errorf("Expected ValidateMethod to return %v, got %v", tt.expected, result)
			}

			if !tt.expected && w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestWriteSuccessResponse(t *testing.T) {
	w := httptest.NewRecorder()
	message := "Test success"
	data := map[string]string{"key": "value"}

	WriteSuccessResponse(w, http.StatusOK, message, data)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response APIResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}

	if response.Message != message {
		t.Errorf("Expected message %s, got %s", message, response.Message)
	}

	if response.Data == nil {
		t.Error("Expected data to be present")
	}
}

func TestWriteErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()
	message := "Test error"
	err := &HTTPError{Code: http.StatusBadRequest, Message: "Bad request"}

	WriteErrorResponse(w, http.StatusBadRequest, message, err)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response APIResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if response.Success {
		t.Error("Expected success to be false")
	}

	if response.Message != message {
		t.Errorf("Expected message %s, got %s", message, response.Message)
	}

	if response.Error != err.Error() {
		t.Errorf("Expected error %s, got %s", err.Error(), response.Error)
	}
}

func TestWritePaginatedResponse(t *testing.T) {
	w := httptest.NewRecorder()
	message := "Test paginated"
	data := []string{"item1", "item2"}
	pagination := PaginationInfo{
		Page:       1,
		PageSize:   10,
		Total:      2,
		TotalPages: 1,
	}

	WritePaginatedResponse(w, http.StatusOK, message, data, pagination)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response PaginatedResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}

	if response.Message != message {
		t.Errorf("Expected message %s, got %s", message, response.Message)
	}

	if response.Pagination.Page != pagination.Page {
		t.Errorf("Expected page %d, got %d", pagination.Page, response.Pagination.Page)
	}

	if response.Pagination.Total != pagination.Total {
		t.Errorf("Expected total %d, got %d", pagination.Total, response.Pagination.Total)
	}
}

func TestHTTPError(t *testing.T) {
	err := &HTTPError{Code: http.StatusNotFound, Message: "Not found"}

	if err.Error() != "Not found" {
		t.Errorf("Expected error message 'Not found', got %s", err.Error())
	}

	if err.Code != http.StatusNotFound {
		t.Errorf("Expected code %d, got %d", http.StatusNotFound, err.Code)
	}
}

func TestCommonErrorTypes(t *testing.T) {
	tests := []struct {
		name     string
		err      *HTTPError
		expected int
	}{
		{"InvalidRequest", ErrInvalidRequest, http.StatusBadRequest},
		{"NotFound", ErrNotFound, http.StatusNotFound},
		{"InternalServer", ErrInternalServer, http.StatusInternalServerError},
		{"Unauthorized", ErrUnauthorized, http.StatusUnauthorized},
		{"Forbidden", ErrForbidden, http.StatusForbidden},
		{"Conflict", ErrConflict, http.StatusConflict},
		{"TooManyRequests", ErrTooManyRequests, http.StatusTooManyRequests},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.expected {
				t.Errorf("Expected code %d, got %d", tt.expected, tt.err.Code)
			}
		})
	}
}

func TestResponseHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	WriteSuccessResponse(w, http.StatusOK, "test", nil)

	// Check CORS headers
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("Expected Access-Control-Allow-Origin header")
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("Expected Content-Type header to be application/json")
	}
}

func TestJSONEncodingError(t *testing.T) {
	// This test is harder to implement as it requires causing a JSON encoding error
	// In practice, this would be handled by the json.NewEncoder(w).Encode() call
	// For now, we'll just test that the function doesn't panic
	w := httptest.NewRecorder()
	
	// Use a channel to cause a JSON encoding error
	ch := make(chan int)
	data := map[string]interface{}{"channel": ch}
	
	WriteSuccessResponse(w, http.StatusOK, "test", data)
	
	// Should not panic and should return some error response
	if w.Code == 0 {
		t.Error("Expected some status code to be set")
	}
}
