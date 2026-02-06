# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build scheduler service (cmd/scheduler)
RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler ./cmd/scheduler

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/scheduler ./scheduler
COPY --from=builder /app/config ./config

# Create uploads directory
RUN mkdir -p uploads

# Default to user service; can be overridden in docker-compose
CMD ["./scheduler"]
