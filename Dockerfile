# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/orders-api ./src/main.go

# Install goose for migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates and netcat for healthcheck
RUN apk --no-cache add ca-certificates netcat-openbsd

# Copy binary from builder
COPY --from=builder /app/orders-api .
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Copy config, migrations, and entrypoint
COPY --from=builder /app/src/config/config.yml ./src/config/
COPY --from=builder /app/src/database/migrations ./src/database/migrations
COPY entrypoint.sh /app/entrypoint.sh

RUN chmod +x /app/entrypoint.sh

EXPOSE 3000

ENTRYPOINT ["/app/entrypoint.sh"]


