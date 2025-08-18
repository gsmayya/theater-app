package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gsmayya/theater/shows"
)

// for any new just returns ok with a message.
func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status": "Ticket Backend Service is running"}`)
}

func ShowListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	show, err := shows.GetShows()
	if err != nil {
		http.Error(w, "Failed to get shows", http.StatusInternalServerError)
		return
	}
	log.Println(show)
	json.NewEncoder(w).Encode(show)
}

func addToShowsPut(w http.ResponseWriter, r *http.Request) (string, error) {
	var show_obj *shows.ShowData
	show_info := show_obj.NewShowFromPut(r)
	uuid := show_info.Show_Id.String()
	// Call the function to add the show
	err := shows.PutShow(show_info)
	if err != nil {
		http.Error(w, "Failed to add show", http.StatusInternalServerError)
		return "", err
	}
	return uuid, nil
}

func getShowDetails(w http.ResponseWriter, r *http.Request) (string, error) {
	show_id := r.URL.Query().Get("show_id")
	if show_id == "" {
		http.Error(w, "Show ID is required", http.StatusBadRequest)
		return show_id, fmt.Errorf("show_id is required")
	}
	show, err := shows.GetShow(show_id)
	if err != nil {
		http.Error(w, "Failed to get shows", http.StatusInternalServerError)
		return show_id, err
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(show)
	return show_id, nil
}

func ShowHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPut {
		uuid, err := addToShowsPut(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("uuid", uuid)
	}

	if r.Method == http.MethodGet {
		uuid, err := getShowDetails(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("uuid", uuid)
	}
}
