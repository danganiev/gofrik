#!/bin/bash

# Example script to set up a blog with Gofrik GraphQL CMS
# Make sure the server is running on http://localhost:8080

BASE_URL="http://localhost:8080/graphql"

echo "=== Setting up a blog with Gofrik GraphQL CMS ==="
echo

# 1. Register a user
echo "1. Registering admin user..."
REGISTER_RESPONSE=$(curl -s -X POST ${BASE_URL} \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { register(email: \"admin@blog.com\", password: \"admin123\") { id email created_at } }"
  }')
echo "Response: $REGISTER_RESPONSE"
echo

# 2. Login
echo "2. Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST ${BASE_URL} \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { login(email: \"admin@blog.com\", password: \"admin123\") { token user { id email } } }"
  }')
echo "Response: $LOGIN_RESPONSE"

# Extract token (requires jq)
if command -v jq &> /dev/null; then
  TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.login.token')
  echo "Token: $TOKEN"
else
  echo "Note: Install 'jq' to automatically extract the token"
  TOKEN="YOUR_TOKEN_HERE"
fi
echo

# 3. Create Blog Post content type
echo "3. Creating Blog Post content type..."
curl -s -X POST ${BASE_URL} \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -d '{
    "query": "mutation { createContentType(name: \"Blog Post\", slug: \"blog-post\", description: \"Blog posts for the website\", schema: \"{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"title\\\":{\\\"type\\\":\\\"string\\\"},\\\"slug\\\":{\\\"type\\\":\\\"string\\\"},\\\"body\\\":{\\\"type\\\":\\\"string\\\"},\\\"excerpt\\\":{\\\"type\\\":\\\"string\\\"},\\\"author\\\":{\\\"type\\\":\\\"string\\\"},\\\"tags\\\":{\\\"type\\\":\\\"array\\\",\\\"items\\\":{\\\"type\\\":\\\"string\\\"}}}}\") { id name slug schema } }"
  }' | jq .
echo

# 4. Create Category content type
echo "4. Creating Category content type..."
curl -s -X POST ${BASE_URL} \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -d '{
    "query": "mutation { createContentType(name: \"Category\", slug: \"category\", description: \"Blog categories\", schema: \"{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"name\\\":{\\\"type\\\":\\\"string\\\"},\\\"slug\\\":{\\\"type\\\":\\\"string\\\"},\\\"description\\\":{\\\"type\\\":\\\"string\\\"}}}\") { id name slug } }"
  }' | jq .
echo

# 5. Create a category
echo "5. Creating a category..."
curl -s -X POST ${BASE_URL} \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -d '{
    "query": "mutation { createContent(typeSlug: \"category\", data: \"{\\\"name\\\":\\\"Technology\\\",\\\"slug\\\":\\\"technology\\\",\\\"description\\\":\\\"Posts about technology and programming\\\"}\", status: \"published\") { id data status } }"
  }' | jq .
echo

# 6. Create blog posts
echo "6. Creating sample blog posts..."

echo "   - Creating published post..."
curl -s -X POST ${BASE_URL} \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -d '{
    "query": "mutation { createContent(typeSlug: \"blog-post\", data: \"{\\\"title\\\":\\\"Getting Started with Go\\\",\\\"slug\\\":\\\"getting-started-with-go\\\",\\\"body\\\":\\\"Go is a statically typed, compiled programming language designed at Google...\\\",\\\"excerpt\\\":\\\"Learn the basics of Go programming\\\",\\\"author\\\":\\\"John Doe\\\",\\\"tags\\\":[\\\"go\\\",\\\"programming\\\",\\\"tutorial\\\"]}\", status: \"published\") { id data status created_at } }"
  }' | jq .
echo

echo "   - Creating draft post..."
curl -s -X POST ${BASE_URL} \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -d '{
    "query": "mutation { createContent(typeSlug: \"blog-post\", data: \"{\\\"title\\\":\\\"Building a Headless CMS\\\",\\\"slug\\\":\\\"building-a-headless-cms\\\",\\\"body\\\":\\\"A headless CMS provides content through an API without a built-in frontend...\\\",\\\"excerpt\\\":\\\"Learn how to build a CMS from scratch\\\",\\\"author\\\":\\\"Jane Smith\\\",\\\"tags\\\":[\\\"cms\\\",\\\"go\\\",\\\"backend\\\"]}\", status: \"draft\") { id data status } }"
  }' | jq .
echo

# 7. Query all blog posts
echo "7. Listing all blog posts..."
curl -s -X POST ${BASE_URL} \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -d '{
    "query": "query { content(typeSlug: \"blog-post\") { id data status created_at } }"
  }' | jq .
echo

echo "=== Blog setup complete! ==="
echo
echo "You can now:"
echo "  1. Visit http://localhost:8080/graphql to use the GraphiQL playground"
echo "  2. Query your content using GraphQL"
echo "  3. Try this query:"
echo
echo "     query {"
echo "       content(typeSlug: \"blog-post\") {"
echo "         id"
echo "         data"
echo "         status"
echo "       }"
echo "     }"
echo
