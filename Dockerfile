# Build stage
FROM golang:1.25.1-alpine AS builder

# Install git and ca-certificates (needed for fetching dependencies and HTTPS)
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 for static binary
# -ldflags="-w -s" to reduce binary size (strip debug info)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o athena \
    ./cmd/api-server

# Final stage - minimal runtime image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 athena && \
    adduser -D -u 1000 -G athena athena

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/athena .

# Change ownership to non-root user
RUN chown -R athena:athena /app

# Switch to non-root user
USER athena

# Expose port
EXPOSE 1323

# Set environment variables (can be overridden at runtime)
ENV JWT_SECRET="your-secret-key-change-this-in-production"

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:1323/ping || exit 1

# Run the application
CMD ["./athena"]
