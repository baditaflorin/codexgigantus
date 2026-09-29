# Multi-stage Dockerfile for CodexGigantus with security hardening
# Security: Use specific version tags, not 'latest'
# Stage 1: Build the application
FROM golang:1.22-alpine3.19 AS builder

# Security: Add labels for image metadata
LABEL maintainer="CodexGigantus Security Team"
LABEL version="2.0.0"
LABEL description="Secure CodexGigantus build"

# Security: Run as non-root during build
RUN addgroup -g 1000 builder && \
    adduser -D -u 1000 -G builder builder

# Install build dependencies with specific versions
RUN apk add --no-cache \
    git=~2.43 \
    make=~4.4 \
    gcc=~13.2 \
    musl-dev=~1.2 \
    sqlite-dev=~3.44

# Set working directory
WORKDIR /build

# Security: Change ownership to builder user
RUN chown -R builder:builder /build

# Switch to builder user
USER builder

# Copy go mod files
COPY --chown=builder:builder go.mod go.sum ./

# Download dependencies with verification
RUN go mod download && go mod verify

# Copy source code
COPY --chown=builder:builder . .

# Build both CLI and Web binaries with security flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo \
    -ldflags='-w -s -extldflags "-static"' \
    -trimpath \
    -o codexgigantus-cli ./cmd/cli && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo \
    -ldflags='-w -s -extldflags "-static"' \
    -trimpath \
    -o codexgigantus-web ./cmd/web

# Security: Verify binaries are statically linked
RUN ldd codexgigantus-cli 2>&1 | grep -q "not a dynamic executable" && \
    ldd codexgigantus-web 2>&1 | grep -q "not a dynamic executable"

# Stage 2: Create minimal runtime image with security hardening
FROM alpine:3.19

# Security: Add labels
LABEL maintainer="CodexGigantus Security Team"
LABEL version="2.0.0"
LABEL description="Secure CodexGigantus runtime"

# Security: Install only essential runtime dependencies
RUN apk add --no-cache \
    ca-certificates=~20240226 \
    tzdata=~2024a && \
    update-ca-certificates

# Security: Create non-root user with fixed UID/GID
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    mkdir -p /app /app/configs /app/data /tmp/app && \
    chown -R appuser:appuser /app /tmp/app

# Set working directory
WORKDIR /app

# Security: Copy binaries with restricted permissions
COPY --from=builder --chown=root:root --chmod=555 /build/codexgigantus-cli /app/
COPY --from=builder --chown=root:root --chmod=555 /build/codexgigantus-web /app/

# Copy configuration examples
COPY --chown=root:root --chmod=444 .env.example /app/

# Security: Create necessary directories with proper permissions
RUN chown -R appuser:appuser /app/configs /app/data /tmp/app && \
    chmod 755 /app/configs /app/data && \
    chmod 1777 /tmp/app

# Expose web GUI port
EXPOSE 8080

# Security: Set environment variables
ENV APP_MODE=web \
    PATH=/app:$PATH \
    TMPDIR=/tmp/app \
    HOME=/app

# Security: Switch to non-root user
USER appuser

# Health check with reduced privileges
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider --timeout=2 http://localhost:8080/ || exit 1

# Security: Use exec form for ENTRYPOINT
ENTRYPOINT ["/app/codexgigantus-web"]

# Security: Default command (can be overridden)
CMD []
