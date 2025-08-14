package utils

import (
	"testing"
)

func TestGetShows(t *testing.T) {
	shows, err := GetShows()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(shows) == 0 {
		t.Error("Expected shows to be non-empty")
	}
	if shows["show1"] != "Movie 1" {
		t.Error("Expected show1 to be 'Movie 1'")
	}
}

func TestPutShow(t *testing.T) {
	err := PutShow("show1", "Movie 1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
