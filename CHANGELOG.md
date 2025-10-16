# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **S3-Compatible Storage Support** - File upload system supporting AWS S3, Google Cloud Storage, MinIO, DigitalOcean Spaces, Backblaze B2, and any S3-compatible service
- `/upload` HTTP endpoint for multipart file uploads
- `internal/storage` package with S3 client integration
- Image upload with automatic validation (type, size limits)
- Unique filename generation to prevent collisions
- Support for custom S3 endpoints (MinIO, DigitalOcean, etc.)
- CDN integration via `STORAGE_PUBLIC_URL` configuration
- Storage configuration via environment variables
- `STORAGE.md` - Comprehensive storage configuration guide
- `examples/storage_setup.sh` - Interactive storage setup script
- AWS SDK Go v2 dependencies for S3 operations
- File type validation (JPEG, PNG, GIF, WebP, SVG)
- 10MB file size limit for uploads

### Changed

- **Complete Docker-based development workflow** - All development now happens in Docker
- Revised Makefile with Docker-first commands (`make up`, `make dev`, `make test`, etc.)
- Updated README with Docker-based quick start and development instructions
- Updated README with storage configuration section
- Simplified prerequisites - only Docker and Make required
- Enhanced `docker-compose.yml` with storage environment variables
- Enhanced `docker-compose.dev.yml` with storage configuration
- Updated main.go to use HTTP mux for multiple endpoints (`/graphql` and `/upload`)
- Storage is now optional - gracefully disables if not configured

## [0.1.0] - 2025-10-16

### Added

- GraphQL API with single endpoint and full introspection
- Built-in GraphiQL playground for API exploration
- Dynamic content type management using JSON schemas
- Token-based user authentication system
- User registration and login mutations
- Draft/publish workflow for content entries
- Content type CRUD operations (create, read, update, delete)
- Content entry CRUD operations
- PostgreSQL JSONB storage for flexible content data
- Database connection and migrations
- Docker and Docker Compose support
- Environment variable configuration (PORT, DATABASE_URL)
- Example scripts and GraphQL queries
- Query content by type slug
- Query content types by ID or slug
- HTTP handler with GraphQL endpoint at `/graphql`

[Unreleased]: https://github.com/yourusername/gofrik/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/yourusername/gofrik/releases/tag/v0.1.0
