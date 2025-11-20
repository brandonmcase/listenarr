# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.21-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o listenarr ./cmd/listenarr

# Stage 3: Runtime
FROM alpine:latest
WORKDIR /app

# Install PHP, FFmpeg, and m4b-tool
RUN apk add --no-cache \
    php-cli \
    php-mbstring \
    php-xml \
    php-zip \
    php-curl \
    ffmpeg \
    curl \
    ca-certificates && \
    curl -L -o /usr/local/bin/m4b-tool \
    https://github.com/sandreas/m4b-tool/releases/latest/download/m4b-tool.phar && \
    chmod +x /usr/local/bin/m4b-tool && \
    m4b-tool --version

# Copy built artifacts
COPY --from=frontend-builder /app/frontend/dist ./wwwroot
COPY --from=backend-builder /app/listenarr /usr/local/bin/

# Set up user for proper permissions
RUN addgroup -S listenarr && adduser -S listenarr -G listenarr
RUN chown -R listenarr:listenarr /app
USER listenarr

EXPOSE 8686

HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD wget -q -O- http://localhost:8686/api/health || exit 1

CMD ["listenarr"]

