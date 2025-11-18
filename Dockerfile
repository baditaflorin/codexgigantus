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

# Install build dependencies
# Note: Alpine apk doesn't support flexible version constraints, using latest available in Alpine 3.19
RUN apk add --no-cache \
    git \
    make \
    gcc \
    musl-dev \
    sqlite-dev

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
# Note: CGO_ENABLED=1 is required for SQLite support (go-sqlite3)
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -a \
    -ldflags='-w -s' \
    -trimpath \
    -o codexgigantus-cli ./cmd/cli && \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -a \
    -ldflags='-w -s' \
    -trimpath \
    -o codexgigantus-web ./cmd/web

# Security: Verify binaries exist and are executable
RUN test -x codexgigantus-cli && test -x codexgigantus-web

# Stage 2: Create minimal runtime image with security hardening
FROM alpine:3.19

# Security: Add labels
LABEL maintainer="CodexGigantus Security Team"
LABEL version="2.0.0"
LABEL description="Secure CodexGigantus runtime"

# Security: Install only essential runtime dependencies
# Note: SQLite requires libc and other runtime libraries
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    sqlite-libs && \
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
