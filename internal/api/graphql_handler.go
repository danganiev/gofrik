package api

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"gofrik/internal/auth"
	gofrikGraphQL "gofrik/internal/graphql"

	"github.com/graphql-go/handler"
)

type GraphQLHandler struct {
	db      *sql.DB
	auth    *auth.Middleware
	handler *handler.Handler
}

func NewGraphQLHandler(db *sql.DB) (*GraphQLHandler, error) {
	authMW := auth.NewMiddleware()

	// Create GraphQL schema
	schema, err := gofrikGraphQL.NewSchema(db, authMW)
	if err != nil {
		return nil, err
	}

	// Create GraphQL handler with GraphiQL enabled
	gqlSchema := schema.GetSchema()
	h := handler.New(&handler.Config{
		Schema:     &gqlSchema,
		Pretty:     true,
		GraphiQL:   true,
		Playground: true,
	})

	return &GraphQLHandler{
		db:      db,
		auth:    authMW,
		handler: h,
	}, nil
}

func (h *GraphQLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Start with request context
	ctx := r.Context()

	// Check for authentication token
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			token := parts[1]
			if session, ok := h.auth.GetSession(token); ok {
				// Add session to context
				ctx = context.WithValue(ctx, "session", session)
			}
		}
	}

	// Serve GraphQL
	h.handler.ContextHandler(ctx, w, r)
}

