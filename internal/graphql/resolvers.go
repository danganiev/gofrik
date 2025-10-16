package graphql

import (
	"encoding/json"
	"fmt"

	"gofrik/internal/auth"
	"gofrik/internal/models"

	"github.com/graphql-go/graphql"
)

// Query Resolvers

func (s *Schema) resolveContentTypes(p graphql.ResolveParams) (interface{}, error) {
	types, err := models.ListContentTypes(s.db)
	if err != nil {
		return nil, err
	}
	
	// Convert to map format for GraphQL
	var result []map[string]interface{}
	for _, ct := range types {
		result = append(result, map[string]interface{}{
			"id":          ct.ID,
			"name":        ct.Name,
			"slug":        ct.Slug,
			"description": ct.Description,
			"schema":      string(ct.Schema),
			"created_at":  ct.CreatedAt,
			"updated_at":  ct.UpdatedAt,
		})
	}
	
	return result, nil
}

func (s *Schema) resolveContentType(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(int)
	if !ok {
		return nil, fmt.Errorf("invalid id")
	}
	
	ct, err := models.GetContentType(s.db, id)
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"id":          ct.ID,
		"name":        ct.Name,
		"slug":        ct.Slug,
		"description": ct.Description,
		"schema":      string(ct.Schema),
		"created_at":  ct.CreatedAt,
		"updated_at":  ct.UpdatedAt,
	}, nil
}

func (s *Schema) resolveContentTypeBySlug(p graphql.ResolveParams) (interface{}, error) {
	slug, ok := p.Args["slug"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid slug")
	}
	
	ct, err := models.GetContentTypeBySlug(s.db, slug)
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"id":          ct.ID,
		"name":        ct.Name,
		"slug":        ct.Slug,
		"description": ct.Description,
		"schema":      string(ct.Schema),
		"created_at":  ct.CreatedAt,
		"updated_at":  ct.UpdatedAt,
	}, nil
}

func (s *Schema) resolveContent(p graphql.ResolveParams) (interface{}, error) {
	typeSlug, ok := p.Args["typeSlug"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid typeSlug")
	}
	
	ct, err := models.GetContentTypeBySlug(s.db, typeSlug)
	if err != nil {
		return nil, err
	}
	
	entries, err := models.ListContentEntries(s.db, ct.ID)
	if err != nil {
		return nil, err
	}
	
	// Convert to map format for GraphQL
	var result []map[string]interface{}
	for _, entry := range entries {
		item := map[string]interface{}{
			"id":              entry.ID,
			"content_type_id": entry.ContentTypeID,
			"data":            string(entry.Data),
			"status":          entry.Status,
			"created_at":      entry.CreatedAt,
			"updated_at":      entry.UpdatedAt,
		}
		if entry.CreatedBy != nil {
			item["created_by"] = *entry.CreatedBy
		}
		if entry.PublishedAt != nil {
			item["published_at"] = *entry.PublishedAt
		}
		result = append(result, item)
	}
	
	return result, nil
}

func (s *Schema) resolveContentEntry(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(int)
	if !ok {
		return nil, fmt.Errorf("invalid id")
	}
	
	entry, err := models.GetContentEntry(s.db, id)
	if err != nil {
		return nil, err
	}
	
	result := map[string]interface{}{
		"id":              entry.ID,
		"content_type_id": entry.ContentTypeID,
		"data":            string(entry.Data),
		"status":          entry.Status,
		"created_at":      entry.CreatedAt,
		"updated_at":      entry.UpdatedAt,
	}
	if entry.CreatedBy != nil {
		result["created_by"] = *entry.CreatedBy
	}
	if entry.PublishedAt != nil {
		result["published_at"] = *entry.PublishedAt
	}
	
	return result, nil
}

// Mutation Resolvers

func (s *Schema) resolveRegister(p graphql.ResolveParams) (interface{}, error) {
	email, _ := p.Args["email"].(string)
	password, _ := p.Args["password"].(string)
	
	if email == "" || password == "" {
		return nil, fmt.Errorf("email and password are required")
	}
	
	// Check if a user already exists (only allow 1 user total)
	count, err := models.CountUsers(s.db)
	if err != nil {
		return nil, err
	}
	
	if count >= 1 {
		return nil, fmt.Errorf("registration is not allowed: a user is already registered")
	}
	
	user, err := models.CreateUser(s.db, email, password)
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"id":         user.ID,
		"email":      user.Email,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}, nil
}

func (s *Schema) resolveLogin(p graphql.ResolveParams) (interface{}, error) {
	email, _ := p.Args["email"].(string)
	password, _ := p.Args["password"].(string)
	
	user, err := models.GetUserByEmail(s.db, email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	
	if !user.CheckPassword(password) {
		return nil, fmt.Errorf("invalid credentials")
	}
	
	token, err := s.auth.CreateSession(user.ID, user.Email)
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":         user.ID,
			"email":      user.Email,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	}, nil
}

func (s *Schema) resolveCreateContentType(p graphql.ResolveParams) (interface{}, error) {
	name, _ := p.Args["name"].(string)
	slug, _ := p.Args["slug"].(string)
	description, _ := p.Args["description"].(string)
	schemaStr, _ := p.Args["schema"].(string)
	
	if name == "" || slug == "" || schemaStr == "" {
		return nil, fmt.Errorf("name, slug, and schema are required")
	}
	
	// Validate JSON schema
	var schemaData interface{}
	if err := json.Unmarshal([]byte(schemaStr), &schemaData); err != nil {
		return nil, fmt.Errorf("invalid schema JSON: %w", err)
	}
	
	ct, err := models.CreateContentType(s.db, name, slug, description, json.RawMessage(schemaStr))
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"id":          ct.ID,
		"name":        ct.Name,
		"slug":        ct.Slug,
		"description": ct.Description,
		"schema":      string(ct.Schema),
		"created_at":  ct.CreatedAt,
		"updated_at":  ct.UpdatedAt,
	}, nil
}

func (s *Schema) resolveUpdateContentType(p graphql.ResolveParams) (interface{}, error) {
	id, _ := p.Args["id"].(int)
	
	// Get existing content type
	ct, err := models.GetContentType(s.db, id)
	if err != nil {
		return nil, err
	}
	
	// Update fields if provided
	name := ct.Name
	if n, ok := p.Args["name"].(string); ok && n != "" {
		name = n
	}
	
	description := ct.Description
	if d, ok := p.Args["description"].(string); ok {
		description = d
	}
	
	schema := ct.Schema
	if s, ok := p.Args["schema"].(string); ok && s != "" {
		// Validate JSON
		var schemaData interface{}
		if err := json.Unmarshal([]byte(s), &schemaData); err != nil {
			return nil, fmt.Errorf("invalid schema JSON: %w", err)
		}
		schema = json.RawMessage(s)
	}
	
	if err := models.UpdateContentType(s.db, id, name, description, schema); err != nil {
		return nil, err
	}
	
	// Fetch updated content type
	updated, err := models.GetContentType(s.db, id)
	if err != nil {
		return nil, err
	}
	
	return map[string]interface{}{
		"id":          updated.ID,
		"name":        updated.Name,
		"slug":        updated.Slug,
		"description": updated.Description,
		"schema":      string(updated.Schema),
		"created_at":  updated.CreatedAt,
		"updated_at":  updated.UpdatedAt,
	}, nil
}

func (s *Schema) resolveDeleteContentType(p graphql.ResolveParams) (interface{}, error) {
	id, _ := p.Args["id"].(int)
	
	if err := models.DeleteContentType(s.db, id); err != nil {
		return false, err
	}
	
	return true, nil
}

func (s *Schema) resolveCreateContent(p graphql.ResolveParams) (interface{}, error) {
	typeSlug, _ := p.Args["typeSlug"].(string)
	dataStr, _ := p.Args["data"].(string)
	status, _ := p.Args["status"].(string)
	
	if typeSlug == "" || dataStr == "" {
		return nil, fmt.Errorf("typeSlug and data are required")
	}
	
	if status == "" {
		status = "draft"
	}
	
	// Get content type
	ct, err := models.GetContentTypeBySlug(s.db, typeSlug)
	if err != nil {
		return nil, err
	}
	
	// Validate JSON data
	var data interface{}
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		return nil, fmt.Errorf("invalid data JSON: %w", err)
	}
	
	// Get user from context if available
	var createdBy *int
	if session, ok := p.Context.Value("session").(*auth.Session); ok && session != nil {
		createdBy = &session.UserID
	}
	
	entry, err := models.CreateContentEntry(s.db, ct.ID, json.RawMessage(dataStr), status, createdBy)
	if err != nil {
		return nil, err
	}
	
	result := map[string]interface{}{
		"id":              entry.ID,
		"content_type_id": entry.ContentTypeID,
		"data":            string(entry.Data),
		"status":          entry.Status,
		"created_at":      entry.CreatedAt,
		"updated_at":      entry.UpdatedAt,
	}
	if entry.CreatedBy != nil {
		result["created_by"] = *entry.CreatedBy
	}
	if entry.PublishedAt != nil {
		result["published_at"] = *entry.PublishedAt
	}
	
	return result, nil
}

func (s *Schema) resolveUpdateContent(p graphql.ResolveParams) (interface{}, error) {
	id, _ := p.Args["id"].(int)
	
	// Get existing entry
	entry, err := models.GetContentEntry(s.db, id)
	if err != nil {
		return nil, err
	}
	
	// Update fields if provided
	data := entry.Data
	if d, ok := p.Args["data"].(string); ok && d != "" {
		// Validate JSON
		var jsonData interface{}
		if err := json.Unmarshal([]byte(d), &jsonData); err != nil {
			return nil, fmt.Errorf("invalid data JSON: %w", err)
		}
		data = json.RawMessage(d)
	}
	
	status := entry.Status
	if s, ok := p.Args["status"].(string); ok && s != "" {
		status = s
	}
	
	if err := models.UpdateContentEntry(s.db, id, data, status); err != nil {
		return nil, err
	}
	
	// Fetch updated entry
	updated, err := models.GetContentEntry(s.db, id)
	if err != nil {
		return nil, err
	}
	
	result := map[string]interface{}{
		"id":              updated.ID,
		"content_type_id": updated.ContentTypeID,
		"data":            string(updated.Data),
		"status":          updated.Status,
		"created_at":      updated.CreatedAt,
		"updated_at":      updated.UpdatedAt,
	}
	if updated.CreatedBy != nil {
		result["created_by"] = *updated.CreatedBy
	}
	if updated.PublishedAt != nil {
		result["published_at"] = *updated.PublishedAt
	}
	
	return result, nil
}

func (s *Schema) resolveDeleteContent(p graphql.ResolveParams) (interface{}, error) {
	id, _ := p.Args["id"].(int)
	
	if err := models.DeleteContentEntry(s.db, id); err != nil {
		return false, err
	}
	
	return true, nil
}

