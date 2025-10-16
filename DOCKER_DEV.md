# Docker-Based Development Guide

This project uses Docker for all development to ensure consistency across environments.

## Quick Reference

### Essential Commands

```bash
# Start services (production mode)
make up

# Start services (development mode with hot reload)
make dev

# View logs
make logs

# Stop services
make down

# Run tests
make test

# Get help
make help
```

## Development Modes

### Production Mode (`make up`)

- Uses the multi-stage `Dockerfile`
- Optimized for production deployment
- No hot reload
- Smaller image size

### Development Mode (`make dev`)

- Uses `Dockerfile.dev`
- Includes development tools (Air, golangci-lint)
- Source code mounted as volume
- Automatic hot reload with Air
- Changes to `.go` files trigger automatic rebuild

## File Structure

```
.
├── Dockerfile              # Production build (multi-stage)
├── Dockerfile.dev          # Development build (with tools)
├── docker-compose.yml      # Base compose configuration
├── docker-compose.dev.yml  # Development overrides
├── .air.toml              # Air hot reload configuration
├── Makefile               # Development commands
└── .env.example           # Environment variables template
```

## How It Works

### Production Build

The production `Dockerfile` uses a multi-stage build:

1. **Builder stage**: Compiles the Go binary
2. **Final stage**: Minimal Alpine image with just the binary

### Development Build

The development setup:

1. Uses `Dockerfile.dev` with Go tools installed
2. Mounts source code as a volume (changes reflected immediately)
3. Runs Air for automatic reloading
4. Includes linting and formatting tools

### Docker Compose

- **docker-compose.yml**: Base configuration for PostgreSQL and app
- **docker-compose.dev.yml**: Development overrides (volumes, command, etc.)

When you run `make dev`, it merges both files:

```bash
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```

## Common Workflows

### Starting a New Feature

```bash
# Start development environment
make dev

# In another terminal, view logs
make logs-app

# Make your changes in your editor
# Air will automatically reload on save
```

### Running Tests

```bash
# With services already running
make test

# Or without starting services first
make test-standalone
```

### Database Operations

```bash
# Open PostgreSQL shell
make db-shell

# Then in psql:
# \dt              List tables
# \d+ users        Describe users table
# SELECT * FROM users;
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint
```

### Debugging

```bash
# Open a shell in the app container
make shell

# Check Go version
make exec CMD="go version"

# View environment variables
make exec CMD="env"

# Run database migrations manually
make exec CMD="go run main.go migrate"
```

### Cleanup

```bash
# Stop containers (keep volumes)
make down

# Remove everything including database data
make clean
```

## Environment Variables

Create a `.env` file from the example:

```bash
cp .env.example .env
```

Edit `.env` to customize:

- Database credentials
- Server port
- Other configuration

Docker Compose will automatically load this file.

## Troubleshooting

### Services won't start

```bash
# Check what's running
make ps

# View logs for errors
make logs

# Clean everything and start fresh
make clean
make up
```

### Code changes not reloading

- Check that you're running `make dev` (not `make up`)
- Verify Air is running: `make logs-app`
- Check `.air.toml` configuration

### Database connection issues

```bash
# Check database is healthy
make ps

# View database logs
make logs-db

# Connect to database directly
make db-shell
```

### Port already in use

If port 8080 or 5432 is already in use:

```bash
# Edit .env file
PORT=8081  # or any other port

# Restart
make down
make up
```

## Advanced Usage

### Running arbitrary commands

```bash
# In running container
make exec CMD="go mod tidy"

# In new container
make run CMD="go test -bench ."
```

### Building production image

```bash
# Build the production image
make build

# Test the production image
docker-compose up
```

### Accessing container directly

```bash
# Shell into app container
make shell

# Shell into database container
docker-compose exec postgres sh
```

## Tips

1. **Always use `make dev`** for active development (hot reload)
2. **Use `make up`** only for testing production builds
3. **Keep services running** and use `make exec` for commands
4. **Check logs frequently** with `make logs-app`
5. **Use `make help`** to see all available commands

## Integration with IDEs

### VS Code

With services running (`make up`), you can still use VS Code's Go extension for:

- IntelliSense
- Go to definition
- Debugging (attach to process)

### Terminal

Keep multiple terminals open:

1. `make dev` - Running services
2. `make logs-app` - Viewing logs
3. Your regular shell - Running commands

## Further Reading

- [Air Documentation](https://github.com/cosmtrek/air)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Go in Docker Best Practices](https://docs.docker.com/language/golang/)
