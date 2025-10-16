package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// Config holds the storage configuration
type Config struct {
	Provider        string // "s3", "gcs", "minio", etc. (all use S3 API)
	Bucket          string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string // Custom endpoint for S3-compatible services (MinIO, DigitalOcean, etc.)
	PublicURL       string // Public URL base for accessing files (e.g., CDN URL)
}

// Storage handles file uploads to S3-compatible storage
type Storage struct {
	config *Config
	client *s3.Client
}

// NewStorage creates a new storage instance
func NewStorage(cfg *Config) (*Storage, error) {
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("storage bucket is required")
	}

	// Create AWS config
	var awsCfg aws.Config
	var err error

	if cfg.AccessKeyID != "" && cfg.SecretAccessKey != "" {
		// Use static credentials
		awsCfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(cfg.Region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				cfg.AccessKeyID,
				cfg.SecretAccessKey,
				"",
			)),
		)
	} else {
		// Use default credential chain (IAM roles, env vars, etc.)
		awsCfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(cfg.Region),
		)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client with optional custom endpoint
	var s3Client *s3.Client
	if cfg.Endpoint != "" {
		// Custom endpoint for S3-compatible services
		s3Client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			o.UsePathStyle = true // Required for MinIO and some other S3-compatible services
		})
	} else {
		s3Client = s3.NewFromConfig(awsCfg)
	}

	return &Storage{
		config: cfg,
		client: s3Client,
	}, nil
}

// UploadFile uploads a file to S3-compatible storage and returns the public URL
func (s *Storage) UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	// Validate file size (max 10MB)
	const maxFileSize = 10 * 1024 * 1024
	if header.Size > maxFileSize {
		return "", fmt.Errorf("file too large: maximum size is 10MB")
	}

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	if !isValidImageType(contentType) {
		return "", fmt.Errorf("invalid file type: only images are allowed (jpg, jpeg, png, gif, webp, svg)")
	}

	// Generate unique filename
	filename := generateUniqueFilename(header.Filename)

	// Create the upload input
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.config.Bucket),
		Key:         aws.String(filename),
		Body:        file,
		ContentType: aws.String(contentType),
		ACL:         "public-read", // Make file publicly accessible
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Generate public URL
	url := s.getPublicURL(filename)

	return url, nil
}

// DeleteFile deletes a file from storage
func (s *Storage) DeleteFile(ctx context.Context, url string) error {
	// Extract filename from URL
	filename := s.extractFilenameFromURL(url)
	if filename == "" {
		return fmt.Errorf("invalid URL")
	}

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(filename),
	})

	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// getPublicURL generates the public URL for a file
func (s *Storage) getPublicURL(filename string) string {
	// If custom public URL is provided, use it
	if s.config.PublicURL != "" {
		return fmt.Sprintf("%s/%s", strings.TrimRight(s.config.PublicURL, "/"), filename)
	}

	// If custom endpoint is provided (MinIO, DigitalOcean, etc.)
	if s.config.Endpoint != "" {
		return fmt.Sprintf("%s/%s/%s", strings.TrimRight(s.config.Endpoint, "/"), s.config.Bucket, filename)
	}

	// Default S3 URL format
	if s.config.Region == "us-east-1" {
		return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.config.Bucket, filename)
	}
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.config.Bucket, s.config.Region, filename)
}

// extractFilenameFromURL extracts the filename from a public URL
func (s *Storage) extractFilenameFromURL(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

// isValidImageType checks if the content type is a valid image type
func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
		"image/svg+xml",
	}

	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}
	return false
}

// generateUniqueFilename generates a unique filename with timestamp and UUID
func generateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	timestamp := time.Now().Unix()
	uniqueID := uuid.New().String()[:8]
	
	// Clean the original filename (remove extension and special chars)
	baseName := strings.TrimSuffix(originalFilename, ext)
	baseName = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		return '-'
	}, baseName)
	
	// Limit base name length
	if len(baseName) > 50 {
		baseName = baseName[:50]
	}
	
	return fmt.Sprintf("%d-%s-%s%s", timestamp, uniqueID, baseName, ext)
}

// LoadConfigFromEnv loads storage configuration from environment variables
func LoadConfigFromEnv() *Config {
	provider := os.Getenv("STORAGE_PROVIDER")
	if provider == "" {
		provider = "s3" // Default to AWS S3
	}

	return &Config{
		Provider:        provider,
		Bucket:          os.Getenv("STORAGE_BUCKET"),
		Region:          getEnvOrDefault("STORAGE_REGION", "us-east-1"),
		AccessKeyID:     os.Getenv("STORAGE_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("STORAGE_SECRET_ACCESS_KEY"),
		Endpoint:        os.Getenv("STORAGE_ENDPOINT"),
		PublicURL:       os.Getenv("STORAGE_PUBLIC_URL"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

