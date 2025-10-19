# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files and source (needed for go mod tidy to see all imports)
COPY go.mod go.sum ./
COPY . .

# Tidy and download dependencies
RUN go mod tidy && go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gofrik .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/gofrik .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./gofrik"]

