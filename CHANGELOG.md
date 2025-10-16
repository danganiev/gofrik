# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- **Complete Docker-based development workflow** - All development now happens in Docker
- Revised Makefile with Docker-first commands (`make up`, `make dev`, `make test`, etc.)
- Updated README with Docker-based quick start and development instructions
- Simplified prerequisites - only Docker and Make required

### Added

- `Dockerfile.dev` - Development Dockerfile with hot reload and tooling
- `docker-compose.dev.yml` - Development overrides for docker-compose
- `.air.toml` - Air configuration for automatic hot reloading
- `DOCKER_DEV.md` - Comprehensive Docker development guide
- Development mode with hot reload (`make dev`)
- Container shell access commands (`make shell`, `make db-shell`)
- Standalone test command for running tests without starting services
- Custom command execution (`make exec CMD="..."`, `make run CMD="..."`)
- Separate log viewing commands (`make logs-app`, `make logs-db`)

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
