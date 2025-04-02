package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	envPath := filepath.Join(".", ".env")
	if _, err := os.Stat(envPath); err == nil {
		if err := godotenv.Load(envPath); err != nil {
			log.Printf("Warning: Error loading .env file: %v", err)
		} else {
			log.Println("Loaded environment variables from .env file")
		}
	}

	if err := StartServer(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
