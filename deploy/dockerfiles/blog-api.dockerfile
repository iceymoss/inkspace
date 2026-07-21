# Build blog frontend
FROM node:22-alpine AS frontend

WORKDIR /src/web/blog
COPY web/blog/package.json web/blog/package-lock.json ./
RUN npm ci
COPY web/blog/ ./
RUN VITE_OUT_DIR=dist npm run build

# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

RUN rm -rf internal/webassets/blog/dist
COPY --from=frontend /src/web/blog/dist ./internal/webassets/blog/dist

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

EXPOSE 8081

# Default to user service; can be overridden in docker-compose
CMD ["./server"]
