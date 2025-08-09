# Multi-stage Dockerfile for GoReleaser builds
# Optimized for minimal size and security

# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata

# Create non-root user
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
# Note: GoReleaser will override this in its build process
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -extldflags '-static'" \
    -o app \
    ./cmd/*/

# Final stage
FROM scratch

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy user
COPY --from=builder /etc/passwd /etc/passwd

# Copy the binary from GoReleaser build (when used with GoReleaser)
# Otherwise copy from builder stage
COPY app /app

# Use non-root user
USER appuser

# Expose port (adjust as needed)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app", "--health"]

# Run the binary
ENTRYPOINT ["/app"]