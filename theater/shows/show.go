package shows

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// Show represents a show in theater
type ShowData struct {
	Show_Id        uuid.UUID `json:"show_id"`
	Name           string    `json:"name"`
	Details        string    `json:"details"`
	Price          int32     `json:"price"`
	Total_Tickets  int32     `json:"total_tickets"`
	Location       string    `json:"location"`
	Booked_Tickets int32     `json:"booked_tickets"`
	ShowNumber     string    `json:"show_number"`
	ShowDate       time.Time `json:"show_date"`
	Images         []string  `json:"images,omitempty"`        // CMS image IDs
	Videos         []string  `json:"videos,omitempty"`        // CMS video IDs
}

func (s *ShowData) NewShow(name string, details string, price int32, total_tickets int32, location string) *ShowData {
	s.Show_Id = uuid.New()
	s.Booked_Tickets = 0
	s.Price = price
	s.Total_Tickets = total_tickets
	s.Location = location
	s.Name = name
	s.Details = details
	s.ShowNumber = fmt.Sprintf("SH-%d", time.Now().Unix()) // Generate show number
	s.ShowDate = time.Now().AddDate(0, 0, 30) // Default 30 days from now
	s.Images = []string{}
	s.Videos = []string{}
	return s
}

func (s *ShowData) NewShowFromPut(r *http.Request) *ShowData {
	showName := r.URL.Query().Get("name")
	showDetails := r.URL.Query().Get("details")
	priceStr := r.URL.Query().Get("price")
	price, err := strconv.ParseInt(priceStr, 10, 32)
	if err != nil {
		return nil
	}

	totalTicketsStr := r.URL.Query().Get("total_tickets")
	totalTickets, err := strconv.ParseInt(totalTicketsStr, 10, 32)
	if err != nil {
		return nil
	}

	showLocation := r.URL.Query().Get("location")
	showNumber := r.URL.Query().Get("show_number")
	showDateStr := r.URL.Query().Get("show_date")
	
	// Parse show date
	var showDate time.Time
	if showDateStr != "" {
		if parsedDate, err := time.Parse(time.RFC3339, showDateStr); err == nil {
			showDate = parsedDate
		} else {
			showDate = time.Now().AddDate(0, 0, 30) // Default 30 days from now
		}
	} else {
		showDate = time.Now().AddDate(0, 0, 30)
	}

	s.Show_Id = uuid.New()
	s.Booked_Tickets = 0
	s.Price = int32(price)
	s.Total_Tickets = int32(totalTickets)
	s.Location = showLocation
	s.Name = showName
	s.Details = showDetails
	s.ShowNumber = showNumber
	if s.ShowNumber == "" {
		s.ShowNumber = fmt.Sprintf("SH-%d", time.Now().Unix())
	}
	s.ShowDate = showDate
	s.Images = []string{}
	s.Videos = []string{}
	return s
}

func (s *ShowData) ShowToMap() map[string]string {
	imagesJson, _ := json.Marshal(s.Images)
	videosJson, _ := json.Marshal(s.Videos)
	return map[string]string{
		"show_id":        s.Show_Id.String(),
		"name":           s.Name,
		"details":        s.Details,
		"price":          fmt.Sprintf("%d", s.Price),
		"total_tickets":  fmt.Sprintf("%d", s.Total_Tickets),
		"location":       s.Location,
		"booked_tickets": fmt.Sprintf("%d", s.Booked_Tickets),
		"show_number":    s.ShowNumber,
		"show_date":      s.ShowDate.Format(time.RFC3339),
		"images":         string(imagesJson),
		"videos":         string(videosJson),
	}
}

func (s *ShowData) ShowToJSON() (string, error) {
	jsonData, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (s *ShowData) JSONToShow(data string) (*ShowData, error) {

	err := json.Unmarshal([]byte(data), &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}
