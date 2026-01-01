# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make gcc musl-dev

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binaries
ARG BUILD=unknown
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags "-X main.build=${BUILD} -s -w" -o /build/bin/indexer cmd/indexer/main.go
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags "-X main.build=${BUILD} -s -w" -o /build/bin/consumer cmd/consumer/main.go

# Indexer runtime stage
FROM alpine:3.19 AS indexer

RUN apk add --no-cache ca-certificates tzdata wget

WORKDIR /app

# Copy indexer binary
COPY --from=builder /build/bin/indexer /app/indexer

# Create data directory
RUN mkdir -p /app/data

# Run as non-root user
RUN addgroup -g 1000 indexer && \
    adduser -D -u 1000 -G indexer indexer && \
    chown -R indexer:indexer /app

USER indexer

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
  CMD wget --spider -q http://localhost:8080/health || exit 1

ENTRYPOINT ["/app/indexer"]
CMD ["-config", "/app/config.toml"]

# Consumer runtime stage
FROM alpine:3.19 AS consumer

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy consumer binary
COPY --from=builder /build/bin/consumer /app/consumer

# Run as non-root user
RUN addgroup -g 1000 consumer && \
    adduser -D -u 1000 -G consumer consumer && \
    chown -R consumer:consumer /app

USER consumer

ENTRYPOINT ["/app/consumer"]
CMD ["-config", "/app/config.toml"]
