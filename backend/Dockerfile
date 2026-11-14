# Build stage
FROM golang:1.27.4-alpine@sha256:0a5f7fe796599ff487605987432e7619ae45beeb2fba4c2ce0b907879932aebf AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -p 1 -o /app/migrate ./cmd/migrate
RUN CGO_ENABLED=0 GOOS=linux go build -p 1 -o /app/api ./cmd/api

# Runtime stage
FROM alpine:3.23@sha256:469c9ea410b9215eaee7976ba6c0639039b444be8e5d638f7abbada189971eb2 AS runtime
RUN apk add --no-cache ca-certificates postgresql-client curl \
 && addgroup -S -g 10001 pos \
 && adduser -S -D -H -u 10001 -G pos pos \
 && mkdir -p /app/uploads \
 && chown -R pos:pos /app

WORKDIR /app

# Copy binaries from builder
COPY --from=builder --chown=pos:pos /app/migrate ./migrate
COPY --from=builder --chown=pos:pos /app/api ./api

# Copy migrations
COPY --from=builder --chown=pos:pos /app/migrations ./migrations

# Copy startup script
COPY --chown=pos:pos docker-entrypoint.sh ./docker-entrypoint.sh
RUN chmod 0555 ./docker-entrypoint.sh

USER 10001:10001
EXPOSE 8080

ENTRYPOINT ["/app/docker-entrypoint.sh"]
