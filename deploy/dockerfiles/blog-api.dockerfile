# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build user service (cmd/server)
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binaries and config
COPY --from=builder /app/server ./server
COPY --from=builder /app/config ./config

# Create uploads directory
RUN mkdir -p uploads

EXPOSE 8080

# Default to user service; can be overridden in docker-compose
CMD ["./server"]