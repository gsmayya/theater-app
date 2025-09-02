# Multi-stage build for Go application
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o theater .

# Production image
FROM alpine:3.18

# Install runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    mysql-client \
    wget \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN adduser -D -s /bin/sh -u 1001 appuser

WORKDIR /app

# Copy built application
COPY --from=builder /app/theater .
COPY --from=builder /app/scripts ./scripts
COPY --from=builder /app/run.sh .

# Make scripts executable
RUN chmod +x theater \
    && chmod +x scripts/*.sh \
    && chmod +x run.sh \
    && chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/health || exit 1

EXPOSE 8080

CMD ["./run.sh"]
