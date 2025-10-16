package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"gofrik/internal/api"
	"gofrik/internal/database"
	"gofrik/internal/storage"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Create context that listens for interrupt signals
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Load configuration
	config := api.LoadConfig()

	// Initialize database
	db, err := database.Connect(config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.Migrate(db); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	log.Println("Database migrations completed")

	// Initialize storage (if configured)
	var storageClient *storage.Storage
	storageConfig := storage.LoadConfigFromEnv()
	if storageConfig.Bucket != "" {
		storageClient, err = storage.NewStorage(storageConfig)
		if err != nil {
			log.Printf("Warning: Failed to initialize storage: %v", err)
			log.Printf("File upload will not be available")
		} else {
			log.Printf("Storage initialized: %s (bucket: %s)", storageConfig.Provider, storageConfig.Bucket)
		}
	} else {
		log.Printf("Storage not configured. Set STORAGE_BUCKET to enable file uploads")
	}

	// Create server with all dependencies
	srv, err := api.NewServer(
		config,
		db,
		storageClient,
	)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: srv,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting Gofrik GraphQL server on %s", httpServer.Addr)
		log.Printf("GraphQL endpoint: http://localhost:%s/graphql", config.Port)
		log.Printf("GraphiQL playground: http://localhost:%s/graphql", config.Port)
		log.Printf("Health check: http://localhost:%s/health", config.Port)
		if storageClient != nil {
			log.Printf("Upload endpoint: http://localhost:%s/upload", config.Port)
		}

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	// Wait for interrupt signal
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		log.Println("Shutting down server...")

		// Create shutdown context with timeout
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()

	wg.Wait()
	log.Println("Server stopped gracefully")
	return nil
}

