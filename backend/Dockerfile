# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies with retry logic for package mirrors
RUN apk update && \
    apk add --no-cache git make || \
    (apk update && apk add --no-cache git make) || \
    (sleep 5 && apk update && apk add --no-cache git make)

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/migrate ./cmd/migrate
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/api ./cmd/api

# Runtime stage
FROM alpine:3.20

# Install runtime dependencies with retry logic for package mirrors
RUN apk update && \
    apk add --no-cache ca-certificates postgresql-client curl || \
    (apk update && apk add --no-cache ca-certificates postgresql-client curl) || \
    (sleep 5 && apk update && apk add --no-cache ca-certificates postgresql-client curl)

WORKDIR /root/

# Copy binaries from builder
COPY --from=builder /app/migrate .
COPY --from=builder /app/api .

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Copy startup script
COPY docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["./docker-entrypoint.sh"]
