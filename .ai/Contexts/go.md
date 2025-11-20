# Go (Golang) Implementation

## Overview

Listenarr is built with Go (Golang), chosen for its excellent concurrency model, small Docker images, and great fit for this use case.

## Project Structure

```
listenarr/
├── cmd/
│   └── listenarr/        # Main application entry point
├── internal/              # Private application code
│   ├── api/              # HTTP handlers and routes
│   ├── config/           # Configuration management
│   ├── database/         # Database initialization and migrations
│   ├── models/           # Data models
│   └── services/         # Business logic
│       ├── download/     # Download management
│       ├── library/      # Library management
│       ├── metadata/     # Metadata fetching
│       ├── processing/   # File processing
│       ├── search/       # Search functionality
│       └── plex/         # Plex integration
├── pkg/                   # Public packages (reusable)
│   ├── qbit/            # qBittorrent client
│   ├── jackett/          # Jackett client
│   ├── plex/             # Plex client
│   └── m4b/              # m4b-tool wrapper
├── frontend/              # React frontend
└── config/                # Configuration files
```

## Key Libraries

### Web Framework
- **Gin**: Lightweight, fast HTTP web framework
- Alternative: Echo (more features) or Fiber (Express-like)

### Database
- **GORM**: Full-featured ORM for Go
- **sqlite driver**: GORM SQLite driver
- Alternative: sqlx for more SQL control

### Configuration
- **viper**: Configuration management with support for YAML, JSON, env vars

### Logging
- **logrus**: Structured logging (recommended)
- Alternative: zap (faster, more structured)

### HTTP Client
- **Built-in net/http**: Standard library HTTP client
- Alternative: resty (more convenient API)

### Background Jobs
- **Native goroutines**: For simple concurrent tasks
- **asynq**: For more advanced job queuing (if needed)

## Concurrency Model

Go's goroutines are perfect for this application:

```go
// Example: Monitor multiple downloads concurrently
func monitorDownloads(downloads []Download) {
    for _, dl := range downloads {
        go func(d Download) {
            // Monitor download progress
            for {
                status := checkDownloadStatus(d)
                if status == "completed" {
                    triggerProcessing(d)
                    break
                }
                time.Sleep(30 * time.Second)
            }
        }(dl)
    }
}
```

## API Structure

### Routes
- `/api/health` - Health check
- `/api/v1/library` - Library management
- `/api/v1/downloads` - Download management
- `/api/v1/processing` - Processing queue
- `/api/v1/search` - Search functionality

### Response Format
```go
type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}
```

## Database Models

### Core Entities
- `Author` - Author information
- `Book` - Book metadata
- `Audiobook` - Audiobook-specific information
- `Release` - Specific release/edition
- `Download` - Download task and status
- `ProcessingTask` - Processing task and status
- `LibraryItem` - Item in user's library

## Service Layer Pattern

Each service handles a specific domain:

```go
type DownloadService struct {
    db     *gorm.DB
    qbit   *qbit.Client
    config *config.Config
}

func (s *DownloadService) StartDownload(releaseID string) error {
    // Implementation
}
```

## Error Handling

Go's error handling pattern:

```go
if err != nil {
    return fmt.Errorf("context: %w", err)
}
```

## Testing

- Use `testing` package for unit tests
- Use `testify` for assertions and mocks
- Integration tests for API endpoints

## Build and Deployment

### Local Build
```bash
go build -o bin/listenarr ./cmd/listenarr
```

### Docker Build
- Multi-stage build
- Final image uses Alpine Linux
- Single binary deployment

### Cross-Compilation
```bash
GOOS=linux GOARCH=amd64 go build -o bin/listenarr-linux-amd64 ./cmd/listenarr
GOOS=linux GOARCH=arm64 go build -o bin/listenarr-linux-arm64 ./cmd/listenarr
```

## Best Practices

1. **Error Wrapping**: Use `fmt.Errorf` with `%w` for error wrapping
2. **Context**: Use `context.Context` for cancellation and timeouts
3. **Interfaces**: Define interfaces for testability
4. **Package Structure**: Keep internal code in `internal/`, reusable code in `pkg/`
5. **Concurrency**: Use goroutines for I/O-bound operations, channels for coordination
6. **Logging**: Use structured logging with context
7. **Configuration**: Use viper for flexible configuration management

## Development Workflow

1. **Format**: `go fmt ./...`
2. **Lint**: `golangci-lint run`
3. **Test**: `go test ./...`
4. **Build**: `make build`
5. **Run**: `make run`

## Resources

- [Go Documentation](https://go.dev/doc/)
- [Gin Framework](https://gin-gonic.com/)
- [GORM Documentation](https://gorm.io/)
- [Effective Go](https://go.dev/doc/effective_go)

