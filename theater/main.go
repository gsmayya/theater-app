package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gsmayya/theater/handlers"
	"github.com/gsmayya/theater/utils"
)

func main() {

	utils.TestRedis()
	fmt.Println("Starting Ticket Backend Service...")
	http.HandleFunc("/", handlers.DefaultHandler)
	http.HandleFunc("/status", handlers.DefaultHandler)
	http.HandleFunc("/shows", handlers.ShowListHandler)
	http.HandleFunc("/show", handlers.ShowHandler)
	log.Println("Server starting at :8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
