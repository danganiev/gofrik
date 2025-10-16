package graphql

import (
	"github.com/graphql-go/graphql"
)

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

