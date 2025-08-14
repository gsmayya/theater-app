package utils

import (
	"log"
	"os"
	"strconv"
)

func GetEnvOrDefault(key, defaultValue string) string {
	val, exists := os.LookupEnv(key)
	if !exists || val == "" { // Check if not exists OR if exists but is empty
		return defaultValue
	}
	return val
}

func GetInt32(str string) int32 {
	val, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		log.Fatalf("Error converting string to int32: %v", err)
	}
	return int32(val)
}
