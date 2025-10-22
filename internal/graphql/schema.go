package graphql

import (
	"database/sql"

	"gofrik/internal/auth"

	"github.com/graphql-go/graphql"
)

type Schema struct {
	db     *sql.DB
	auth   *auth.Middleware
	schema graphql.Schema
}

func NewSchema(db *sql.DB, authMW *auth.Middleware) (*Schema, error) {
	s := &Schema{
		db:   db,
		auth: authMW,
	}
	
	// Define types
	userType := s.getUserType()
	contentTypeType := s.getContentTypeType()
	contentEntryType := s.getContentEntryType()
	pageInfoType := getPageInfoType()
	contentTypesResponseType := getContentTypesResponseType(contentTypeType, pageInfoType)
	contentEntriesResponseType := getContentEntriesResponseType(contentEntryType, pageInfoType)
	
	// Define root query
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"contentTypes": &graphql.Field{
				Type:        contentTypesResponseType,
				Description: "Get content types",
				Args: graphql.FieldConfigArgument{
					"limit": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 10,
						Description:  "Number of items per page (default: 10, max: 100)",
					},
					"offset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
						Description:  "Number of items to skip (default: 0)",
					},
					"orderBy": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "created_at",
						Description:  "Field to order by (id, name, slug, created_at, updated_at)",
					},
					"orderDirection": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "DESC",
						Description:  "Order direction (ASC or DESC)",
					},
				},
				Resolve: s.resolveContentTypes,
			},
			"contentType": &graphql.Field{
				Type:        contentTypeType,
				Description: "Get a content type by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: s.resolveContentType,
			},
			"contentTypeBySlug": &graphql.Field{
				Type:        contentTypeType,
				Description: "Get a content type by slug",
				Args: graphql.FieldConfigArgument{
					"slug": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: s.resolveContentTypeBySlug,
			},
			"content": &graphql.Field{
				Type:        contentEntriesResponseType,
				Description: "Get content entries by type slug",
				Args: graphql.FieldConfigArgument{
					"typeSlug": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.String),
						Description: "Content type slug",
					},
					"limit": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 10,
						Description:  "Number of items per page (default: 10, max: 100)",
					},
					"offset": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
						Description:  "Number of items to skip (default: 0)",
					},
					"orderBy": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "created_at",
						Description:  "Field to order by (id, created_at, updated_at, published_at, status)",
					},
					"orderDirection": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "DESC",
						Description:  "Order direction (ASC or DESC)",
					},
				},
				Resolve: s.resolveContent,
			},
			"contentEntry": &graphql.Field{
				Type:        contentEntryType,
				Description: "Get a content entry by ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: s.resolveContentEntry,
			},
		},
	})
	
	// Define root mutation
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"register": &graphql.Field{
				Type:        userType,
				Description: "Register a new user",
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: s.resolveRegister,
			},
			"login": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "AuthPayload",
					Fields: graphql.Fields{
						"token": &graphql.Field{Type: graphql.String},
						"user":  &graphql.Field{Type: userType},
					},
				}),
				Description: "Login user",
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: s.resolveLogin,
			},
			"createContentType": &graphql.Field{
				Type:        contentTypeType,
				Description: "Create a new content type",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"slug": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"description": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"schema": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: s.resolveCreateContentType,
			},
			"updateContentType": &graphql.Field{
				Type:        contentTypeType,
				Description: "Update a content type",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"description": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"schema": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: s.resolveUpdateContentType,
			},
			"deleteContentType": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Delete a content type",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: s.resolveDeleteContentType,
			},
			"createContent": &graphql.Field{
				Type:        contentEntryType,
				Description: "Create a new content entry",
				Args: graphql.FieldConfigArgument{
					"typeSlug": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"data": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"status": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "draft",
					},
				},
				Resolve: s.resolveCreateContent,
			},
			"updateContent": &graphql.Field{
				Type:        contentEntryType,
				Description: "Update a content entry",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"data": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"status": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: s.resolveUpdateContent,
			},
			"deleteContent": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Delete a content entry",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: s.resolveDeleteContent,
			},
		},
	})
	
	// Create schema
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
	if err != nil {
		return nil, err
	}
	
	s.schema = schema
	return s, nil
}

func (s *Schema) GetSchema() graphql.Schema {
	return s.schema
}

