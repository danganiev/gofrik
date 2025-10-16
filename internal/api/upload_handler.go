package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gofrik/internal/storage"
)

// UploadHandler handles file uploads
type UploadHandler struct {
	storage *storage.Storage
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(storage *storage.Storage) *UploadHandler {
	return &UploadHandler{
		storage: storage,
	}
}

// ServeHTTP handles the upload request
func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only allow POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse form: %v", err), http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get file: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Upload the file
	url, err := h.storage.UploadFile(r.Context(), file, header)
	if err != nil {
		log.Printf("Upload error: %v", err)
		http.Error(w, fmt.Sprintf("Failed to upload file: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the URL as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"url":      url,
		"filename": header.Filename,
		"size":     header.Size,
		"message":  "File uploaded successfully",
	})
}

