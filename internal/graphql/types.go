package graphql

import (
	"github.com/graphql-go/graphql"
)

// Pagination info type
func getPageInfoType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name:        "PageInfo",
		Description: "Information about pagination",
		Fields: graphql.Fields{
			"totalCount": &graphql.Field{
				Type:        graphql.Int,
				Description: "Total number of items available",
			},
			"hasMore": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Whether there are more items available",
			},
			"limit": &graphql.Field{
				Type:        graphql.Int,
				Description: "Number of items per page",
			},
			"offset": &graphql.Field{
				Type:        graphql.Int,
				Description: "Current offset",
			},
		},
	})
}

func (s *Schema) getUserType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name:        "User",
		Description: "A user in the system",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	})
}

func (s *Schema) getContentTypeType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name:        "ContentType",
		Description: "A content type definition",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"slug": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"schema": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	})
}

func (s *Schema) getContentEntryType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name:        "ContentEntry",
		Description: "A content entry",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"content_type_id": &graphql.Field{
				Type: graphql.Int,
			},
			"data": &graphql.Field{
				Type: graphql.String,
			},
			"status": &graphql.Field{
				Type: graphql.String,
			},
			"created_by": &graphql.Field{
				Type: graphql.Int,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"published_at": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	})
}

// Response types for list queries
func getContentTypesResponseType(contentTypeType *graphql.Object, pageInfoType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name:        "ContentTypesResponse",
		Description: "List of content types with pagination info",
		Fields: graphql.Fields{
			"items": &graphql.Field{
				Type:        graphql.NewList(contentTypeType),
				Description: "List of content types",
			},
			"pageInfo": &graphql.Field{
				Type:        pageInfoType,
				Description: "Pagination information",
			},
		},
	})
}

func getContentEntriesResponseType(contentEntryType *graphql.Object, pageInfoType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name:        "ContentEntriesResponse",
		Description: "List of content entries with pagination info",
		Fields: graphql.Fields{
			"items": &graphql.Field{
				Type:        graphql.NewList(contentEntryType),
				Description: "List of content entries",
			},
			"pageInfo": &graphql.Field{
				Type:        pageInfoType,
				Description: "Pagination information",
			},
		},
	})
}

