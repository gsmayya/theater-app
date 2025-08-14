package show

import (
	"testing"
)

func TestPutShow(t *testing.T) {
	err := PutShowDetails("show1", "Movie 1", 10, 100, "Location 1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
