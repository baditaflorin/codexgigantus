# Multi-stage Dockerfile for CodexGigantus
# Stage 1: Build the application
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make gcc musl-dev sqlite-dev

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build both CLI and Web binaries
RUN make build-cli build-web

# Stage 2: Create minimal runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates sqlite-libs

# Create app directory
WORKDIR /app

# Copy binaries from builder
COPY --from=builder /build/codexgigantus-cli /app/
COPY --from=builder /build/codexgigantus-web /app/

# Copy .env.example
COPY .env.example /app/

# Create configs directory
RUN mkdir -p /app/configs /app/data

# Expose web GUI port
EXPOSE 8080

# Default to web mode
ENV APP_MODE=web

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

# Entrypoint script to choose between CLI and Web mode
COPY docker-entrypoint.sh /app/
RUN chmod +x /app/docker-entrypoint.sh

ENTRYPOINT ["/app/docker-entrypoint.sh"]
