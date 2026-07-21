# Build admin frontend
FROM node:22-alpine AS frontend

WORKDIR /src/web/admin
COPY web/admin/package.json web/admin/package-lock.json ./
RUN npm ci
COPY web/admin/ ./
RUN VITE_OUT_DIR=dist npm run build

# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

RUN rm -rf internal/webassets/admin/dist
COPY --from=frontend /src/web/admin/dist ./internal/webassets/admin/dist

# Build admin service (cmd/admin)
RUN CGO_ENABLED=0 GOOS=linux go build -o admin ./cmd/admin

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binaries and config
COPY --from=builder /app/admin ./admin
COPY --from=builder /app/config ./config

# Create uploads directory
RUN mkdir -p uploads

EXPOSE 8083

# Default to user service; can be overridden in docker-compose
CMD ["./admin"]
