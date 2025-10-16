# Gofrik - Dead Simple Opinionated Headless CMS

A dead simple opinionated headless CMS built from scratch in Go. Gofrik uses **GraphQL exclusively** for its API and **PostgreSQL** as its only supported database.

## Philosophy

- 🎯 **Opinionated**: GraphQL-only API, PostgreSQL-only database
- 🚀 **Pure Go**: No web frameworks - complete control over HTTP handling
- 📦 **Flexible Content**: Define any content schema using JSON
- 🔐 **Built-in Auth**: Simple token-based authentication
- 💾 **PostgreSQL Native**: JSONB for flexible content storage

## Features

- ✅ **GraphQL API** - Single endpoint with full introspection
- ✅ **GraphiQL Playground** - Built-in API explorer
- ✅ **Content Type Management** - Dynamic schemas with JSON
- ✅ **User Authentication** - Token-based sessions
- ✅ **Draft/Publish Workflow** - Content status management
- ✅ **PostgreSQL JSONB** - Flexible and queryable JSON storage

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Make (optional, but recommended)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/gofrik.git
cd gofrik
```

2. Start the application with Docker:

```bash
make up
```

That's it! The server will start on `http://localhost:8080` with PostgreSQL automatically configured.

Open your browser to `http://localhost:8080/graphql` to access the GraphiQL playground!

### Common Commands

```bash
# Start all services
make up

# Start in development mode with hot reload
make dev

# View logs
make logs

# Run tests
make test

# Stop services
make down

# See all available commands
make help
```

## GraphQL API

### Endpoint

```
POST /graphql
```

### Authentication

Include the auth token in the Authorization header:

```
Authorization: Bearer YOUR_TOKEN_HERE
```

## GraphQL Examples

### Queries

#### List all content types

```graphql
query {
  contentTypes {
    id
    name
    slug
    description
    schema
    created_at
    updated_at
  }
}
```

#### Get a specific content type by ID

```graphql
query {
  contentType(id: 1) {
    id
    name
    slug
    schema
  }
}
```

#### Get content type by slug

```graphql
query {
  contentTypeBySlug(slug: "blog-post") {
    id
    name
    slug
    schema
  }
}
```

#### Get all content entries for a type

```graphql
query {
  content(typeSlug: "blog-post") {
    id
    data
    status
    created_at
    updated_at
  }
}
```

#### Get a specific content entry

```graphql
query {
  contentEntry(id: 1) {
    id
    content_type_id
    data
    status
    created_by
    created_at
    updated_at
    published_at
  }
}
```

### Mutations

#### Register a new user

```graphql
mutation {
  register(email: "user@example.com", password: "securepass123") {
    id
    email
    created_at
  }
}
```

#### Login

```graphql
mutation {
  login(email: "user@example.com", password: "securepass123") {
    token
    user {
      id
      email
      created_at
    }
  }
}
```

#### Create a content type

```graphql
mutation {
  createContentType(
    name: "Blog Post"
    slug: "blog-post"
    description: "Blog posts for the website"
    schema: "{\"type\":\"object\",\"properties\":{\"title\":{\"type\":\"string\"},\"body\":{\"type\":\"string\"},\"author\":{\"type\":\"string\"}}}"
  ) {
    id
    name
    slug
    schema
  }
}
```

#### Update a content type

```graphql
mutation {
  updateContentType(
    id: 1
    name: "Updated Blog Post"
    description: "Updated description"
  ) {
    id
    name
    description
  }
}
```

#### Delete a content type

```graphql
mutation {
  deleteContentType(id: 1)
}
```

#### Create content

```graphql
mutation {
  createContent(
    typeSlug: "blog-post"
    data: "{\"title\":\"My First Post\",\"body\":\"Hello World!\",\"author\":\"John Doe\"}"
    status: "published"
  ) {
    id
    data
    status
    created_at
  }
}
```

#### Update content

```graphql
mutation {
  updateContent(
    id: 1
    data: "{\"title\":\"Updated Post\",\"body\":\"Updated content\"}"
    status: "published"
  ) {
    id
    data
    status
    updated_at
  }
}
```

#### Delete content

```graphql
mutation {
  deleteContent(id: 1)
}
```

## Complete Example: Creating a Blog

```graphql
# 1. Register a user
mutation {
  register(email: "admin@blog.com", password: "admin123") {
    id
    email
  }
}

# 2. Login to get a token
mutation {
  login(email: "admin@blog.com", password: "admin123") {
    token
    user {
      id
      email
    }
  }
}

# 3. Create a Blog Post content type
mutation {
  createContentType(
    name: "Blog Post"
    slug: "blog-post"
    description: "Blog posts for the website"
    schema: "{\"type\":\"object\",\"properties\":{\"title\":{\"type\":\"string\"},\"slug\":{\"type\":\"string\"},\"body\":{\"type\":\"string\"},\"excerpt\":{\"type\":\"string\"},\"author\":{\"type\":\"string\"},\"tags\":{\"type\":\"array\",\"items\":{\"type\":\"string\"}}}}"
  ) {
    id
    name
    slug
  }
}

# 4. Create a blog post
mutation {
  createContent(
    typeSlug: "blog-post"
    data: "{\"title\":\"Getting Started with Go\",\"body\":\"Go is amazing...\",\"author\":\"John Doe\",\"tags\":[\"go\",\"tutorial\"]}"
    status: "published"
  ) {
    id
    data
    status
  }
}

# 5. Query all blog posts
query {
  content(typeSlug: "blog-post") {
    id
    data
    status
    created_at
  }
}
```

## Using with cURL

```bash
# Register
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"mutation { register(email: \"admin@test.com\", password: \"admin123\") { id email } }"}'

# Login
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"mutation { login(email: \"admin@test.com\", password: \"admin123\") { token user { id email } } }"}'

# Query content types (with auth)
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"query":"query { contentTypes { id name slug } }"}'
```

## Development

All development happens in Docker to ensure consistency across environments. The Makefile provides convenient commands for all common tasks.

### Development Workflow

Start development mode with hot reload:

```bash
make dev
```

This starts the application with Air for automatic reloading when you change code. Any changes to `.go` files will trigger a rebuild.

### Running Tests

```bash
# Run tests in the running container
make test

# Or run tests in a standalone container
make test-standalone
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint
```

### Database Access

```bash
# Open a PostgreSQL shell
make db-shell

# Open a shell in the app container
make shell
```

### Viewing Logs

```bash
# View all logs
make logs

# View app logs only
make logs-app

# View database logs only
make logs-db
```

### Building for Production

```bash
# Build the Docker image
make build

# Or rebuild from scratch (no cache)
make rebuild
```

### Running Custom Commands

```bash
# Run a command in the running app container
make exec CMD="go version"

# Run a command in a new container
make run CMD="go mod tidy"
```

### Cleanup

```bash
# Stop and remove containers
make down

# Remove everything including volumes
make clean
```

## Configuration

Gofrik uses environment variables for configuration. See `.env.example` for all available options.

### Environment Variables

- `PORT` - Server port (default: 8080)
- `DATABASE_URL` - PostgreSQL connection string (default: postgres://gofrik:gofrik@postgres:5432/gofrik?sslmode=disable)
- `POSTGRES_DB` - PostgreSQL database name (default: gofrik)
- `POSTGRES_USER` - PostgreSQL username (default: gofrik)
- `POSTGRES_PASSWORD` - PostgreSQL password (default: gofrik)

To customize settings for development:

```bash
# Copy the example file
cp .env.example .env

# Edit .env with your settings
# Then restart services
make down
make up
```

## Why Go?

Go is great for web services such as this one.

## Why GraphQL?

I like GraphQL much better than REST APIs.

## Why PostgreSQL?

PostgreSQL is simply the best open source database on the market.

## Content Schema Design

Content types use JSON Schema to define structure. Example schemas:

### Blog Post

```json
{
  "type": "object",
  "properties": {
    "title": { "type": "string" },
    "slug": { "type": "string" },
    "body": { "type": "string" },
    "excerpt": { "type": "string" },
    "author": { "type": "string" },
    "featured_image": { "type": "string" },
    "tags": {
      "type": "array",
      "items": { "type": "string" }
    }
  },
  "required": ["title", "body"]
}
```

### Product

```json
{
  "type": "object",
  "properties": {
    "name": { "type": "string" },
    "sku": { "type": "string" },
    "description": { "type": "string" },
    "price": { "type": "number" },
    "currency": { "type": "string" },
    "inventory": { "type": "integer" },
    "images": {
      "type": "array",
      "items": { "type": "string" }
    }
  },
  "required": ["name", "sku", "price"]
}
```

## Future Enhancements

- [ ] GraphQL subscriptions for real-time updates
- [ ] Media/asset management with GraphQL
- [ ] Content versioning
- [ ] Webhooks on content changes
- [ ] API rate limiting
- [ ] Full-text search with PostgreSQL
- [ ] Role-based access control
- [ ] Content relationships and references
- [ ] GraphQL DataLoader for batching

## Contributing

Please don't.

## License

MIT License - see LICENSE file for details
