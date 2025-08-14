package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gsmayya/theater/utils"
)

// for any new just returns ok with a message.
func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status": "Ticket Backend Service is running"}`)
}

func ShowListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	show, err := utils.GetShows()
	if err != nil {
		http.Error(w, "Failed to get shows", http.StatusInternalServerError)
		return
	}
	log.Println(show)
	json.NewEncoder(w).Encode(show)

}

func ShowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.PostForm.Get("name")
		details := r.PostForm.Get("details")
		err := utils.PutShow(name, details)
		if err != nil {
			http.Error(w, "Failed to update show", http.StatusInternalServerError)
			return
		}
		return
	}
	if r.Method == http.MethodGet {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Show name is required", http.StatusBadRequest)
			return
		}
		show, err := utils.GetShows()
		if err != nil {
			http.Error(w, "Failed to get shows", http.StatusInternalServerError)
			return
		}
		details, exists := show[name]
		if !exists {
			http.Error(w, "Show not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{name: details})
		return
	}
}
