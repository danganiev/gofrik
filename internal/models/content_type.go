package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type ContentType struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	Slug        string          `json:"slug"`
	Description string          `json:"description"`
	Schema      json.RawMessage `json:"schema"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func CreateContentType(db *sql.DB, name, slug, description string, schema json.RawMessage) (*ContentType, error) {
	var ct ContentType
	err := db.QueryRow(
		`INSERT INTO content_types (name, slug, description, schema) 
		 VALUES ($1, $2, $3, $4) 
		 RETURNING id, name, slug, description, schema, created_at, updated_at`,
		name, slug, description, schema,
	).Scan(&ct.ID, &ct.Name, &ct.Slug, &ct.Description, &ct.Schema, &ct.CreatedAt, &ct.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create content type: %w", err)
	}

	return &ct, nil
}

func GetContentType(db *sql.DB, id int) (*ContentType, error) {
	var ct ContentType
	err := db.QueryRow(
		`SELECT id, name, slug, description, schema, created_at, updated_at 
		 FROM content_types WHERE id = $1`,
		id,
	).Scan(&ct.ID, &ct.Name, &ct.Slug, &ct.Description, &ct.Schema, &ct.CreatedAt, &ct.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("content type not found")
		}
		return nil, fmt.Errorf("failed to get content type: %w", err)
	}

	return &ct, nil
}

func GetContentTypeBySlug(db *sql.DB, slug string) (*ContentType, error) {
	var ct ContentType
	err := db.QueryRow(
		`SELECT id, name, slug, description, schema, created_at, updated_at 
		 FROM content_types WHERE slug = $1`,
		slug,
	).Scan(&ct.ID, &ct.Name, &ct.Slug, &ct.Description, &ct.Schema, &ct.CreatedAt, &ct.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("content type not found")
		}
		return nil, fmt.Errorf("failed to get content type: %w", err)
	}

	return &ct, nil
}

func ListContentTypes(db *sql.DB) ([]ContentType, error) {
	rows, err := db.Query(
		`SELECT id, name, slug, description, schema, created_at, updated_at 
		 FROM content_types ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list content types: %w", err)
	}
	defer rows.Close()

	var types []ContentType
	for rows.Next() {
		var ct ContentType
		if err := rows.Scan(&ct.ID, &ct.Name, &ct.Slug, &ct.Description, &ct.Schema, &ct.CreatedAt, &ct.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan content type: %w", err)
		}
		types = append(types, ct)
	}

	return types, nil
}

func UpdateContentType(db *sql.DB, id int, name, description string, schema json.RawMessage) error {
	_, err := db.Exec(
		`UPDATE content_types 
		 SET name = $1, description = $2, schema = $3, updated_at = CURRENT_TIMESTAMP 
		 WHERE id = $4`,
		name, description, schema, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update content type: %w", err)
	}
	return nil
}

func DeleteContentType(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM content_types WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete content type: %w", err)
	}
	return nil
}

