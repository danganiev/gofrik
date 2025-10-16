package main

import (
	"log"
	"net/http"
	"os"

	"gofrik/internal/api"
	"gofrik/internal/database"
)

func main() {
	// Get configuration from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://localhost/gofrik?sslmode=disable"
	}

	// Initialize database
	db, err := database.Connect(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create GraphQL handler
	handler, err := api.NewGraphQLHandler(db)
	if err != nil {
		log.Fatalf("Failed to create GraphQL handler: %v", err)
	}
	
	// Start server
	log.Printf("Starting Gofrik GraphQL server on port %s", port)
	log.Printf("GraphQL endpoint: http://localhost:%s/graphql", port)
	log.Printf("GraphiQL playground: http://localhost:%s/graphql", port)
	
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

