package shows

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
}

func (s *ShowData) NewShow(name string, details string, price int32, total_tickets int32, location string) *ShowData {
	s.Show_Id = uuid.New()
	s.Booked_Tickets = 0
	s.Price = price
	s.Total_Tickets = total_tickets
	s.Location = location
	s.Name = name
	s.Details = details
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

	s.Show_Id = uuid.New()
	s.Booked_Tickets = 0
	s.Price = int32(price)
	s.Total_Tickets = int32(totalTickets)
	s.Location = showLocation
	s.Name = showName
	s.Details = showDetails
	return s
}

func (s *ShowData) ShowToMap() map[string]string {
	return map[string]string{
		"show_id":        s.Show_Id.String(),
		"name":           s.Name,
		"details":        s.Details,
		"price":          fmt.Sprintf("%d", s.Price),
		"total_tickets":  fmt.Sprintf("%d", s.Total_Tickets),
		"location":       s.Location,
		"booked_tickets": fmt.Sprintf("%d", s.Booked_Tickets),
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
