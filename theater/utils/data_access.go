package utils

import (
	"log"
	"os"
)

func GetEnvOrDefault(key, defaultValue string) string {
	val, exists := os.LookupEnv(key)
	if !exists || val == "" { // Check if not exists OR if exists but is empty
		return defaultValue
	}
	return val
}

/*
This will eventually will access the database and fetch new details, for now, it is dummy
*/
func GetShows() (map[string]string, error) {
	shows := map[string]string{
		"show1": "Movie 1",
		"show2": "Movie 2",
		"show3": "Movie 3",
	}
	return shows, nil
}

func PutShow(name string, details string) error {
	// This function will eventually update the show details in the database
	log.Println("Updating show:", name, "with details:", details)
	return nil // Simulating a successful update
}
