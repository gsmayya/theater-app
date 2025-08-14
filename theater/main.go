package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gsmayya/theater/utils"
)

// for any new just returns ok with a message.
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status": "Ticket Backend Service is running"}`)
}

func showlistHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	show := utils.GetShows()
	log.Println(show)
	json.NewEncoder(w).Encode(show)

}

func showHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.PostForm.Get("name")
		details := r.PostForm.Get("details")
		utils.PutShow(name, details)
		return
	}
	if r.Method == http.MethodGet {
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Show name is required", http.StatusBadRequest)
			return
		}
		show := utils.GetShows()
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

func main() {
	utils.TestRedis()
	fmt.Println("Starting Ticket Backend Service...")
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/status", defaultHandler)
	http.HandleFunc("/shows", showlistHandler)
	http.HandleFunc("/show", showHandler)
	log.Println("Server starting at :8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
