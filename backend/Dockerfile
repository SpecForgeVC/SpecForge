# Build stage
FROM golang:1.25.6-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod tidy

# Copy source code
COPY . .

# Build the application with optimizations and CGO disabled for portability
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server

# Final stage
FROM alpine:3.19

# Security: Run as a non-root user
RUN addgroup -S specforge && adduser -S specforge -G specforge
USER specforge

WORKDIR /home/specforge

# Copy the binary from the builder stage
COPY --from=builder --chown=specforge:specforge /app/server .
# Copy migrations for startup
COPY --from=builder --chown=specforge:specforge /app/migrations ./migrations

# Expose port
EXPOSE 8080

# Healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the server
ENTRYPOINT ["./server"]
