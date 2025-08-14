package show

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gsmayya/theater/utils"
)

type Show struct {
	Show_Id        uuid.UUID
	Name           string
	Details        string
	Price          int32
	Total_Tickets  int32
	Location       string
	Booked_Tickets int32
}

func NewShow(name string, details string, price int32, total_tickets int32, location string) *Show {
	return &Show{
		Show_Id:        uuid.New(),
		Name:           name,
		Details:        details,
		Price:          price,
		Total_Tickets:  total_tickets,
		Location:       location,
		Booked_Tickets: 0,
	}
}

func NewShowFromPost(r *http.Request) *Show {
	name := r.PostForm.Get("name")
	details := r.PostForm.Get("details")
	price := utils.GetInt32(r.PostForm.Get("price"))
	totalTickets := utils.GetInt32(r.PostForm.Get("total_tickets"))
	location := r.PostForm.Get("location")

	return &Show{
		Show_Id:        uuid.New(),
		Name:           name,
		Details:        details,
		Price:          price,
		Total_Tickets:  totalTickets,
		Location:       location,
		Booked_Tickets: 0,
	}
}

func ShowToMap(s *Show) map[string]string {
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

func ShowToJSON(s *Show) (string, error) {
	jsonData, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func JSONToShow(data string) (*Show, error) {
	var show Show
	err := json.Unmarshal([]byte(data), &show)
	if err != nil {
		return nil, err
	}
	return &show, nil
}
