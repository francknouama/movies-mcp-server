# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o movies-server cmd/server/main.go

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/movies-server .
COPY --from=builder /app/migrations ./migrations

# Create non-root user
RUN addgroup -g 1000 -S movies && \
    adduser -u 1000 -S movies -G movies

# Change ownership
RUN chown -R movies:movies /root

# Switch to non-root user
USER movies

# Run the binary
CMD ["./movies-server"]