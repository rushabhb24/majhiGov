package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"yojana-portal/backend/internal/db"
	"yojana-portal/backend/internal/handlers"
)

func main() {
	// Load environment variables from .env file (if present)
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using system environment variables.")
	}

	// Initialize Database connection
	err = db.InitDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer db.DB.Close()

	// Setup Routes
	http.HandleFunc("/api/schemes", handlers.GetSchemesHandler)
	http.HandleFunc("/api/schemes/", handlers.GetSchemeDetailsHandler) // Handles /api/schemes/:id
	http.HandleFunc("/api/eligibility-check", handlers.CheckEligibilityHandler)

	// Fetch PORT from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Yojana Portal Backend running on http://localhost:%s", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
