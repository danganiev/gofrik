package api

import (
	"embed"
	"html/template"
	"log"
	"net/http"
)

//go:embed templates/*.html
var templatesFS embed.FS

var homeTemplate *template.Template

func init() {
	var err error
	homeTemplate, err = template.ParseFS(templatesFS, "templates/home.html")
	if err != nil {
		log.Fatalf("Failed to parse home template: %v", err)
	}
}

// addRoutes configures all HTTP routes
func (s *Server) addRoutes(mux *http.ServeMux) error {
	// Root page
	mux.HandleFunc("/", s.handleRoot)

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

// handleRoot displays the homepage
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	// Only handle root path
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := struct {
		HasStorage     bool
		StorageStatus  string
		DatabaseStatus string
	}{
		HasStorage:     s.storage != nil,
		StorageStatus:  "Not configured",
		DatabaseStatus: "Not connected",
	}

	if s.storage != nil {
		data.StorageStatus = "Configured"
	}

	// Check database connection
	if s.db != nil {
		if err := s.db.Ping(); err == nil {
			data.DatabaseStatus = "Connected"
		} else {
			data.DatabaseStatus = "Error"
		}
	}

	if err := homeTemplate.Execute(w, data); err != nil {
		log.Printf("Error executing home template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// handleHealth is a simple health check endpoint
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

