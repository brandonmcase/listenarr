# Listenarr

Audiobook collection manager similar to Lidarr and Radarr, designed to automate the process of searching, downloading, processing, and organizing audiobooks.

## Features

- 🔍 Search for audiobooks across multiple torrent indexers (via Jackett)
- ⬇️ Download audiobooks using qBittorrent
- 🔧 Process downloaded files (merge into single m4b files using m4b-tool)
- 📝 Fetch and embed metadata
- 📚 Organize files for Plex integration
- 🐳 Docker-first design

## Technology Stack

- **Backend**: Go (Golang) with Gin framework
- **Frontend**: React with Material-UI
- **Database**: SQLite
- **Integrations**: qBittorrent, Jackett, Plex, m4b-tool

## Prerequisites

- Go 1.21 or later
- Node.js 18+ (for frontend development)
- Docker (for containerized deployment)

## Development Setup

### Backend

```bash
# Install dependencies
go mod download

# Run development server
go run cmd/listenarr/main.go

# Build binary
go build -o bin/listenarr cmd/listenarr/main.go
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

## Docker

```bash
# Build image
docker build -t listenarr:latest .

# Run container
docker run -d \
  --name listenarr \
  -p 8686:8686 \
  -v $(pwd)/config:/config \
  -v $(pwd)/library:/library \
  listenarr:latest
```

Or use Docker Compose:

```bash
docker-compose up -d
```

## Configuration

See `config/config.example.yml` for configuration options.

## Documentation

- [Planning Document](.ai/PLANNING.md)
- [Language Comparison](.ai/LANGUAGE_COMPARISON.md)
- [Docker Setup](.ai/Contexts/docker.md)
- [m4b-tool Integration](.ai/Contexts/m4b-tool.md)

## License

GPL-3.0 (to match Lidarr/Radarr)

