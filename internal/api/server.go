package api

import (
	"database/sql"
	"log"
	"net/http"

	"gofrik/internal/storage"
)

// Server represents the HTTP server with all dependencies
type Server struct {
	config  *Config
	db      *sql.DB
	storage *storage.Storage
}

// NewServer creates a new HTTP server with all dependencies
func NewServer(
	config *Config,
	db *sql.DB,
	storageClient *storage.Storage,
) (http.Handler, error) {
	srv := &Server{
		config:  config,
		db:      db,
		storage: storageClient,
	}

	// Create mux and add routes
	mux := http.NewServeMux()
	if err := srv.addRoutes(mux); err != nil {
		return nil, err
	}

	// Wrap with middleware
	var handler http.Handler = mux
	handler = loggingMiddleware(handler)
	handler = corsMiddleware(handler)

	return handler, nil
}

// loggingMiddleware logs all HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// corsMiddleware adds CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

