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

func TestShowHandlerPost(t *testing.T) {
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
	req := httptest.NewRequest(http.MethodPost, "/show", nil)
	addFormData(req, "Test Show", "Test Show Details", 100, 50, "Test Location")

	ShowHandler(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Code)
	}
	uuid := resp.Header().Get("uuid")
	return uuid
}

func addFormData(req *http.Request, name, details string, price, totalTickets int32, location string) {
	req.PostForm = make(map[string][]string)
	req.PostForm["name"] = []string{name}
	req.PostForm["details"] = []string{details}
	req.PostForm["price"] = []string{fmt.Sprintf("%d", price)}
	req.PostForm["total_tickets"] = []string{fmt.Sprintf("%d", totalTickets)}
	req.PostForm["location"] = []string{location}
}
