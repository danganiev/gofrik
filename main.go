package main

import (
	"context"
	"fmt"
	"io"
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
	ctx := context.Background()
	if err := run(ctx, os.Getenv, os.Stdout, os.Stderr, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(
	ctx context.Context,
	getenv func(string) string,
	stdout, stderr io.Writer,
	args []string,
) error {
	// Create context that listens for interrupt signals
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Create logger that writes to stdout
	logger := log.New(stdout, "", log.LstdFlags)

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
	logger.Println("Database migrations completed")

	// Initialize storage (if configured)
	var storageClient *storage.Storage
	storageConfig := storage.LoadConfigFromEnv()
	if storageConfig.Bucket != "" {
		storageClient, err = storage.NewStorage(storageConfig)
		if err != nil {
			logger.Printf("Warning: Failed to initialize storage: %v", err)
			logger.Printf("File upload will not be available")
		} else {
			logger.Printf("Storage initialized: %s (bucket: %s)", storageConfig.Provider, storageConfig.Bucket)
		}
	} else {
		logger.Printf("Storage not configured. Set STORAGE_BUCKET to enable file uploads")
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
		logger.Printf("Starting Gofrik GraphQL server on %s", httpServer.Addr)
		logger.Printf("GraphQL endpoint: http://localhost:%s/graphql", config.Port)
		logger.Printf("GraphiQL playground: http://localhost:%s/graphql", config.Port)
		logger.Printf("Health check: http://localhost:%s/health", config.Port)
		if storageClient != nil {
			logger.Printf("Upload endpoint: http://localhost:%s/upload", config.Port)
		}

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(stderr, "error listening and serving: %s\n", err)
		}
	}()

	// Wait for interrupt signal
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		logger.Println("Shutting down server...")

		// Create shutdown context with timeout
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(stderr, "error shutting down http server: %s\n", err)
		}
	}()

	wg.Wait()
	logger.Println("Server stopped gracefully")
	return nil
}

