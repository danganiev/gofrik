package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type ContentEntry struct {
	ID            int             `json:"id"`
	ContentTypeID int             `json:"content_type_id"`
	Data          json.RawMessage `json:"data"`
	Status        string          `json:"status"`
	CreatedBy     *int            `json:"created_by"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	PublishedAt   *time.Time      `json:"published_at"`
}

func CreateContentEntry(db *sql.DB, contentTypeID int, data json.RawMessage, status string, createdBy *int) (*ContentEntry, error) {
	var entry ContentEntry
	err := db.QueryRow(
		`INSERT INTO content_entries (content_type_id, data, status, created_by) 
		 VALUES ($1, $2, $3, $4) 
		 RETURNING id, content_type_id, data, status, created_by, created_at, updated_at, published_at`,
		contentTypeID, data, status, createdBy,
	).Scan(&entry.ID, &entry.ContentTypeID, &entry.Data, &entry.Status, &entry.CreatedBy, &entry.CreatedAt, &entry.UpdatedAt, &entry.PublishedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create content entry: %w", err)
	}

	return &entry, nil
}

func GetContentEntry(db *sql.DB, id int) (*ContentEntry, error) {
	var entry ContentEntry
	err := db.QueryRow(
		`SELECT id, content_type_id, data, status, created_by, created_at, updated_at, published_at 
		 FROM content_entries WHERE id = $1`,
		id,
	).Scan(&entry.ID, &entry.ContentTypeID, &entry.Data, &entry.Status, &entry.CreatedBy, &entry.CreatedAt, &entry.UpdatedAt, &entry.PublishedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("content entry not found")
		}
		return nil, fmt.Errorf("failed to get content entry: %w", err)
	}

	return &entry, nil
}

func ListContentEntries(db *sql.DB, contentTypeID int, limit, offset int, orderBy, orderDirection string) ([]ContentEntry, error) {
	// Validate orderBy to prevent SQL injection
	validOrderFields := map[string]bool{
		"id":           true,
		"created_at":   true,
		"updated_at":   true,
		"published_at": true,
		"status":       true,
	}
	if !validOrderFields[orderBy] {
		orderBy = "created_at"
	}

	// Validate order direction
	if orderDirection != "ASC" && orderDirection != "DESC" {
		orderDirection = "DESC"
	}

	query := fmt.Sprintf(
		`SELECT id, content_type_id, data, status, created_by, created_at, updated_at, published_at 
		 FROM content_entries WHERE content_type_id = $1 ORDER BY %s %s LIMIT $2 OFFSET $3`,
		orderBy, orderDirection,
	)

	rows, err := db.Query(query, contentTypeID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list content entries: %w", err)
	}
	defer rows.Close()

	var entries []ContentEntry
	for rows.Next() {
		var entry ContentEntry
		if err := rows.Scan(&entry.ID, &entry.ContentTypeID, &entry.Data, &entry.Status, &entry.CreatedBy, &entry.CreatedAt, &entry.UpdatedAt, &entry.PublishedAt); err != nil {
			return nil, fmt.Errorf("failed to scan content entry: %w", err)
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// CountContentEntries returns the total number of content entries for a given content type
func CountContentEntries(db *sql.DB, contentTypeID int) (int, error) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM content_entries WHERE content_type_id = $1`, contentTypeID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count content entries: %w", err)
	}
	return count, nil
}

func UpdateContentEntry(db *sql.DB, id int, data json.RawMessage, status string) error {
	_, err := db.Exec(
		`UPDATE content_entries 
		 SET data = $1, status = $2, updated_at = CURRENT_TIMESTAMP 
		 WHERE id = $3`,
		data, status, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update content entry: %w", err)
	}
	return nil
}

func DeleteContentEntry(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM content_entries WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete content entry: %w", err)
	}
	return nil
}

