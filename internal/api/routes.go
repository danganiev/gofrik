package api

import (
	"log"
	"net/http"
)

// addRoutes configures all HTTP routes
func (s *Server) addRoutes(mux *http.ServeMux) error {
	// GraphQL endpoint
	graphQLHandler, err := NewGraphQLHandler(s.db)
	if err != nil {
		return err
	}
	mux.Handle("/graphql", graphQLHandler)
	log.Printf("GraphQL endpoint configured at /graphql")

	// Upload endpoint (if storage is configured)
	if s.storage != nil {
		uploadHandler := NewUploadHandler(s.storage)
		mux.Handle("/upload", uploadHandler)
		log.Printf("Upload endpoint configured at /upload")
	}

	// Health check endpoint
	mux.HandleFunc("/health", s.handleHealth)

	return nil
}

// handleHealth is a simple health check endpoint
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

