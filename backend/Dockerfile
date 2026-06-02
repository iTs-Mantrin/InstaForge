# ============================================================================
# Stage 1: Build
# ============================================================================
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w -X main.version=$(git describe --tags --always 2>/dev/null || echo dev)" \
    -o /app/instaforge ./cmd/server

# ============================================================================
# Stage 2: Runtime
# ============================================================================
FROM alpine:3.21 AS runtime

RUN apk add --no-cache ca-certificates tzdata curl

# Create non-root user
RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

# Copy binary
COPY --from=builder /app/instaforge .

# Copy .env if present (for defaults)
COPY --from=builder /app/.env .env.example

# Create temp directory for downloads
RUN mkdir -p /tmp/instaforge && chown -R app:app /tmp/instaforge

# Use non-root user
USER app

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD curl -f http://localhost:8080/health/live || exit 1

EXPOSE 8080

ENTRYPOINT ["./instaforge"]
