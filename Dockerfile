# Build blog frontend
FROM node:22-alpine AS blog-builder

WORKDIR /src/web/blog
COPY web/blog/package.json web/blog/package-lock.json ./
RUN npm ci
COPY web/blog/ ./
RUN VITE_OUT_DIR=dist npm run build

# Build admin frontend
FROM node:22-alpine AS admin-builder

WORKDIR /src/web/admin
COPY web/admin/package.json web/admin/package-lock.json ./
RUN npm ci
COPY web/admin/ ./
RUN VITE_OUT_DIR=dist npm run build

# Shared Go source and dependency cache
FROM golang:1.26-alpine AS go-base

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the user service with only the blog SPA.
FROM go-base AS server-builder
RUN rm -rf internal/webassets/blog/dist
COPY --from=blog-builder /src/web/blog/dist ./internal/webassets/blog/dist
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Build the admin service with only the admin SPA.
FROM go-base AS admin-go-builder
RUN rm -rf internal/webassets/admin/dist
COPY --from=admin-builder /src/web/admin/dist ./internal/webassets/admin/dist
RUN CGO_ENABLED=0 GOOS=linux go build -o admin ./cmd/admin

# Build scheduler without any frontend build dependency.
FROM go-base AS scheduler-builder
RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler ./cmd/scheduler

# Shared runtime base
FROM alpine:latest AS runtime-base

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

COPY --from=go-base /app/config ./config

# Create uploads directory
RUN mkdir -p uploads

# Admin runtime contains only the admin binary and admin SPA.
FROM runtime-base AS admin-runtime
COPY --from=admin-go-builder /app/admin ./admin
EXPOSE 8083
CMD ["./admin"]

# Scheduler runtime contains no frontend assets.
FROM runtime-base AS scheduler-runtime
COPY --from=scheduler-builder /app/scheduler ./scheduler
CMD ["./scheduler"]

# Server is the default target and contains only the blog SPA.
FROM runtime-base AS server-runtime
COPY --from=server-builder /app/server ./server
EXPOSE 8081
CMD ["./server"]
