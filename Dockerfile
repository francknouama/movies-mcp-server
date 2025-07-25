# ==============================================================================
# Production Dockerfile for Movies MCP Server
# ==============================================================================

# Build stage - use official Go image for building
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set build arguments for metadata
ARG VERSION=unknown
ARG BUILD_TIME=unknown
ARG GIT_COMMIT=unknown

# Create non-root user for build
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /build

# Copy go mod files first for better Docker layer caching
COPY go.mod go.sum ./

# Download dependencies (this layer is cached unless go.mod/go.sum changes)
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
    -a \
    -installsuffix cgo \
    -ldflags="-w -s -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}" \
    -o movies-mcp-server \
    ./cmd/server

# ==============================================================================
# Runtime stage - use minimal distroless image for security
FROM gcr.io/distroless/static-debian12:nonroot

# Set metadata labels
LABEL org.opencontainers.image.title="Movies MCP Server"
LABEL org.opencontainers.image.description="Model Context Protocol server for movie database operations"
LABEL org.opencontainers.image.url="https://github.com/your-org/movies-mcp-server"
LABEL org.opencontainers.image.documentation="https://github.com/your-org/movies-mcp-server/blob/main/README.md"
LABEL org.opencontainers.image.source="https://github.com/your-org/movies-mcp-server"
LABEL org.opencontainers.image.version="${VERSION:-unknown}"
LABEL org.opencontainers.image.created="${BUILD_TIME:-unknown}"
LABEL org.opencontainers.image.revision="${GIT_COMMIT:-unknown}"
LABEL org.opencontainers.image.vendor="Your Organization"
LABEL org.opencontainers.image.licenses="MIT"

# Copy timezone data and CA certificates from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary from builder stage
COPY --from=builder /build/movies-mcp-server /usr/local/bin/movies-mcp-server

# Copy migration files (needed for database setup)
COPY --from=builder /build/migrations /migrations

# Set environment variables
ENV TZ=UTC
ENV GO_ENV=production

# Health check configuration
ENV HEALTH_CHECK_PORT=8080
ENV HEALTH_CHECK_PATH=/health

# Database configuration (will be overridden by environment)
ENV DB_HOST=localhost
ENV DB_PORT=5432
ENV DB_NAME=movies_mcp
ENV DB_USER=movies_user
ENV DB_SSLMODE=require
ENV DB_MAX_OPEN_CONNS=25
ENV DB_MAX_IDLE_CONNS=5
ENV DB_CONN_MAX_LIFETIME=1h

# Server configuration
ENV LOG_LEVEL=info
ENV SERVER_TIMEOUT=30s

# Image processing configuration
ENV MAX_IMAGE_SIZE=5242880
ENV ALLOWED_IMAGE_TYPES=image/jpeg,image/png,image/webp
ENV ENABLE_THUMBNAILS=true
ENV THUMBNAIL_SIZE=200x200

# Metrics and monitoring
ENV METRICS_ENABLED=true
ENV METRICS_INTERVAL=30s

# Use non-root user (distroless default is 'nonroot' with UID 65532)
USER 65532:65532

# Expose health check port (MCP protocol uses stdin/stdout, no network port needed)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ["/usr/local/bin/movies-mcp-server", "-health-check"]

# Set the entrypoint
ENTRYPOINT ["/usr/local/bin/movies-mcp-server"]

# Default command (can be overridden)
CMD []