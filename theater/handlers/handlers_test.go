package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultHandler(t *testing.T) {
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	DefaultHandler(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Code)
	}
}

func TestShowListHandler(t *testing.T) {
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/shows", nil)

	ShowListHandler(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Code)
	}
}

func TestShowHandlerPut(t *testing.T) {
	uuid := insertValue(t)
	log.Println("Value obtained is ", uuid)
}

func TestShowHandlerGet(t *testing.T) {
	uuid := insertValue(t)
	log.Println("Value obtained is ", uuid)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/show?show_id="+uuid, nil)

	ShowHandler(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Code)
	}
	uuid2 := resp.Header().Get("uuid")
	log.Println("Value obtained after get ", uuid2)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(b))

}

func insertValue(t *testing.T) string {
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/show", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	addQueryData(req, "Test Show", "Test Show Details", 100, 50, "Test Location")

	ShowHandler(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Code)
	}
	uuid := resp.Header().Get("uuid")
	return uuid
}

func addQueryData(req *http.Request, name, details string, price, totalTickets int32, location string) {
	query := req.URL.Query()
	query.Add("name", name)
	query.Add("details", details)
	query.Add("price", fmt.Sprintf("%d", price))
	query.Add("total_tickets", fmt.Sprintf("%d", totalTickets))
	query.Add("location", location)
	req.URL.RawQuery = query.Encode()
}
