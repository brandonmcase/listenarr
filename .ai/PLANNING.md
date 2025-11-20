# Listenarr - Audiobook Management System

## Project Planning Document

**Goal**: Create an audiobook collection manager similar to [Lidarr](https://github.com/Lidarr/Lidarr) and [Radarr](https://github.com/Radarr/Radarr), but specifically designed for audiobooks with automated file processing and metadata management.

---

## 1. Project Overview

### 1.1 Core Concept

Listenarr will automate the process of:

- Searching for audiobooks across multiple torrent indexers (via Jackett)
- Downloading audiobooks using qBittorrent
- Processing downloaded files (merging into single m4b files using [m4b-tool](https://github.com/sandreas/m4b-tool))
- Fetching and embedding metadata
- Organizing files for Plex integration
- Managing your audiobook library

### 1.2 Key Differentiators from Lidarr/Radarr

- **File Processing**: Automatically merge multi-file audiobooks into single m4b files
- **Chapter Management**: Extract and embed chapter information
- **Audiobook-Specific Metadata**: Focus on book metadata (author, narrator, series, etc.) rather than music/movie metadata
- **Plex Audiobook Library**: Optimized for Plex's audiobook library structure
- **Docker-First**: Designed from the ground up to run easily in Docker containers

---

## 2. Technology Stack

### 2.1 Backend ✅

- **Language**: **Go (Golang)** ✅ - Chosen for excellent concurrency, small Docker images, and great fit for this use case
- **Framework**: **Gin** ✅ - Lightweight, fast, well-documented (chosen)
- **Database**: **SQLite** ✅ - For simplicity, matching Lidarr/Radarr approach
- **ORM/Database Library**: **GORM** ✅ - Full-featured ORM (chosen for productivity)
- **Background Jobs**: **Native goroutines with channels** ✅ - Start simple, add asynq if needed
- **HTTP Client**: **Built-in `net/http`** ✅ - Standard library is sufficient
- **Configuration**: **viper** ✅ - Config management with YAML/env support
- **Logging**: **logrus** ✅ - Structured logging (chosen for simplicity)

### 2.2 Frontend ✅

- **Framework**: **React 18** ✅ - With TypeScript
- **Build Tool**: **Vite** ✅ - Modern, fast build tool
- **UI Library**: **Material-UI (MUI)** ✅ - Component library with dark theme
- **State Management**: **React Query + Zustand** ✅ - React Query for server state, Zustand for client state (if needed)
- **Routing**: **React Router v6** ✅ - Navigation
- **HTTP Client**: **Axios** ✅ - For API communication
- **Styling**: **Emotion** ✅ - CSS-in-JS (comes with MUI)

### 2.3 External Tools & Integrations

- **qBittorrent**: Web API for download management
- **Jackett**: API for torrent indexer aggregation
- **m4b-tool**: Command-line tool for audiobook file processing
- **Plex**: API for library management
- **FFmpeg**: Audio processing (may be required by m4b-tool)

### 2.4 Containerization ✅

- **Docker**: Primary deployment method - must run easily in Docker containers
- **Docker Compose**: ✅ Created for local development and multi-container setups
- **Base Image**: **Alpine Linux** ✅ - Small, efficient base for Go binary
- **Multi-stage Builds**: ✅ 3-stage build (frontend, backend, runtime)
- **Volume Mounts**: ✅ Configured for persistent data (database, config, library)
- **Dockerfile**: ✅ Created with Go + React + m4b-tool

### 2.5 Metadata Sources

- **Audible API**: Primary source for audiobook metadata (may require scraping or third-party APIs)
- **Open Library API**: Free metadata source for books
- **Goodreads API**: Additional metadata and ratings
- **Google Books API**: Fallback metadata source
- **MusicBrainz**: For narrator/artist information (if available)

---

## 3. Docker Architecture

### 3.1 Container Design

- **Single Container Approach**: Main application runs in one container
- **Dependencies**: Container must include:
  - Runtime (depends on language: Go binary, Python runtime, or .NET runtime)
  - m4b-tool binary (or install during build)
  - FFmpeg (if required by m4b-tool)
  - PHP runtime (for m4b-tool)
  - All necessary system libraries

### 3.2 Volume Mounts

- **Config Volume**: `/config` - Application configuration, database
- **Library Volume**: `/library` - Audiobook library (read/write)
- **Downloads Volume**: `/downloads` - Temporary download location (optional)
- **Processing Volume**: `/processing` - Temporary processing workspace (optional)

### 3.3 Environment Variables

- **PUID/PGID**: User/group IDs for file permissions (matching host user)
- **TZ**: Timezone setting
- **UMASK**: File permission mask
- **API Keys**: Optional environment-based configuration

### 3.4 Network Requirements

- **Host Network**: May need host network mode for qBittorrent/Jackett discovery
- **Bridge Network**: Default Docker bridge for container communication
- **Ports**:
  - Web UI port (default: 8686 or configurable)
  - API port (if separate from UI)

### 3.5 Dockerfile Structure

**Note**: Dockerfile structure depends on chosen language. See `LANGUAGE_COMPARISON.md` for detailed language-specific examples and recommendations.

#### Example Dockerfile Structures

**Go Example** (Recommended):

```dockerfile
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
RUN CGO_ENABLED=0 GOOS=linux go build -o listenarr

# Stage 3: Runtime
FROM alpine:latest
WORKDIR /app

# Install PHP, FFmpeg, and m4b-tool
RUN apk add --no-cache \
    php-cli php-mbstring php-xml php-zip php-curl \
    ffmpeg curl && \
    curl -L -o /usr/local/bin/m4b-tool \
    https://github.com/sandreas/m4b-tool/releases/latest/download/m4b-tool.phar && \
    chmod +x /usr/local/bin/m4b-tool

# Copy built artifacts
COPY --from=frontend-builder /app/frontend/dist ./wwwroot
COPY --from=backend-builder /app/listenarr /usr/local/bin/

RUN addgroup -S listenarr && adduser -S listenarr -G listenarr
USER listenarr

EXPOSE 8686
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD wget -q -O- http://localhost:8686/api/health || exit 1

CMD ["listenarr"]
```

**Python Example** (Alternative):

```dockerfile
# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: Build Backend
FROM python:3.11-slim AS backend-builder
WORKDIR /app
COPY requirements.txt .
RUN pip install --user -r requirements.txt

# Stage 3: Runtime
FROM python:3.11-slim
WORKDIR /app

# Install PHP, FFmpeg, and m4b-tool
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    php-cli php-mbstring php-xml php-zip php-curl \
    ffmpeg curl && \
    curl -L -o /usr/local/bin/m4b-tool \
    https://github.com/sandreas/m4b-tool/releases/latest/download/m4b-tool.phar && \
    chmod +x /usr/local/bin/m4b-tool && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=frontend-builder /app/frontend/dist ./wwwroot
COPY --from=backend-builder /root/.local /root/.local
COPY . .
ENV PATH=/root/.local/bin:$PATH

RUN groupadd -r listenarr && useradd -r -g listenarr listenarr
USER listenarr

EXPOSE 8686
CMD ["uvicorn", "listenarr.main:app", "--host", "0.0.0.0", "--port", "8686"]
```

**Note**: See `/ai/LANGUAGE_COMPARISON.md` for detailed language comparison and `/ai/Contexts/m4b-tool.md` for m4b-tool installation details.

#### Example docker-compose.yml

```yaml
version: "3.8"

services:
  listenarr:
    image: listenarr/listenarr:latest
    container_name: listenarr
    environment:
      - PUID=1000
      - PGID=1000
      - TZ=America/New_York
      - UMASK=002
    volumes:
      - ./config:/config
      - ./library:/library
      - ./downloads:/downloads
      - ./processing:/processing
    ports:
      - "8686:8686"
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8686/api/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
    networks:
      - listenarr-network

networks:
  listenarr-network:
    driver: bridge
```

### 3.6 Container Considerations

- **File Permissions**: Ensure proper permissions for mounted volumes
- **Resource Limits**: CPU/memory limits for processing tasks
- **Health Checks**: Implement health check endpoint
- **Signal Handling**: Graceful shutdown on SIGTERM
- **Logging**: Structured logging to stdout/stderr for Docker log collection

---

## 4. System Architecture

### 4.1 Core Modules

#### 4.1.1 Search Module

- **Purpose**: Interface with Jackett to search for audiobooks
- **Responsibilities**:
  - Query multiple indexers via Jackett
  - Filter results by quality, format, and relevance
  - Rank results based on user preferences
  - Handle search result caching

#### 4.1.2 Download Module

- **Purpose**: Manage downloads via qBittorrent
- **Responsibilities**:
  - Add torrents to qBittorrent
  - Monitor download progress
  - Handle download failures and retries
  - Manage download queue priorities
  - Track download history

#### 4.1.3 Processing Module

- **Purpose**: Process downloaded audiobook files using m4b-tool
- **Responsibilities**:
  - Detect when downloads complete
  - Identify audiobook file formats (mp3, m4a, m4b, etc.)
  - Merge multiple files into single m4b file
  - Extract and embed chapter information
  - Handle processing errors and retries
  - Clean up temporary files

#### 4.1.4 Metadata Module

- **Purpose**: Fetch and embed audiobook metadata
- **Responsibilities**:
  - Search metadata sources (Audible, Open Library, etc.)
  - Match downloaded files to metadata records
  - Extract metadata from existing files (if available)
  - Embed metadata into processed m4b files
  - Download and embed cover art
  - Handle metadata conflicts and user preferences

#### 4.1.5 Library Module

- **Purpose**: Organize and manage audiobook library
- **Responsibilities**:
  - Organize files into library structure (Author/Series/Book format)
  - Track library contents and status
  - Handle file renaming and organization
  - Manage library scanning and updates
  - Integration with Plex library updates

#### 4.1.6 Plex Integration Module

- **Purpose**: Integrate with Plex Media Server
- **Responsibilities**:
  - Update Plex library when new audiobooks are added
  - Send notifications to Plex
  - Handle Plex library refresh requests
  - Support Plex metadata preferences

### 4.2 Data Models

#### Core Entities (To Be Implemented)

- **Author**: Author information and metadata
  - Fields: ID, Name, Biography, Image URL
- **Book**: Book information (title, ISBN, etc.)
  - Fields: ID, Title, ISBN, ASIN, Description, Cover Art URL, Release Date
- **Audiobook**: Audiobook-specific information (narrator, duration, format)
  - Fields: ID, Book ID, Narrator, Duration, Format, Bitrate, Language
- **Series**: Book series information
  - Fields: ID, Name, Description, Total Books
- **Release**: Specific release/edition of an audiobook
  - Fields: ID, Audiobook ID, Quality, Format, Size, Indexer, Magnet/Torrent URL
- **Download**: Download task and status
  - Fields: ID, Library Item ID, Release ID, Status, Progress, Speed, Error
- **ProcessingTask**: Processing task and status
  - Fields: ID, Download ID, Status, Progress, Input Path, Output Path, Error
- **LibraryItem**: Item in user's library with status
  - Fields: ID, Book ID, Audiobook ID, Status, File Path, Added Date, Completed Date

---

## 5. Workflow Design

### 5.1 Complete Download & Processing Workflow

```
1. User adds audiobook to library (manual or automatic)
   ↓
2. Search Module queries Jackett for available releases
   ↓
3. User selects release (or auto-selects based on quality preferences)
   ↓
4. Download Module adds torrent to qBittorrent
   ↓
5. Monitor download progress
   ↓
6. Download completes → Trigger Processing Module
   ↓
7. Processing Module:
   a. Identifies file types and structure
   b. Calls m4b-tool to merge files into single m4b
   c. Extracts chapter information
   ↓
8. Metadata Module:
   a. Searches metadata sources for book information
   b. Matches downloaded book to metadata record
   c. Downloads cover art
   ↓
9. Embed metadata and cover art into m4b file
   ↓
10. Library Module:
    a. Organizes file into library structure
    b. Renames file according to naming convention
    ↓
11. Plex Integration Module:
    a. Updates Plex library
    b. Triggers library refresh
    ↓
12. Mark audiobook as "Available" in library
```

### 5.2 File Processing Details

#### m4b-tool Integration

- **Command Structure**: Execute m4b-tool as subprocess
- **Input**: Downloaded audiobook files (mp3, m4a, etc.)
- **Output**: Single m4b file with chapters
- **Error Handling**: Log errors, retry on failures, notify user
- **Progress Tracking**: Parse m4b-tool output for progress updates

#### Metadata Embedding

- Use m4b-tool's metadata embedding capabilities
- Or use FFmpeg/AtomicParsley for metadata embedding
- Embed: Title, Author, Narrator, Series, Series Position, Description, Cover Art, Chapters

---

## 6. Integration Specifications

### 6.1 qBittorrent Integration

- **API Endpoints**:
  - Authentication
  - Add torrent (from URL or file)
  - Get torrent list
  - Get torrent info
  - Pause/Resume torrents
  - Delete torrents
- **Polling**: Check download status every 30 seconds
- **Event Handling**: Detect completion, handle errors

### 6.2 Jackett Integration

- **API Endpoints**:
  - Search across all indexers
  - Get indexer capabilities
  - Test indexer connectivity
- **Search Parameters**: Title, Author, ISBN
- **Result Filtering**: Filter by format, quality, size

### 6.3 m4b-tool Integration

- **Execution**: Run as subprocess with appropriate flags
- **Key Commands**:
  - Merge: `m4b-tool merge input/ --output-file output.m4b`
  - Chapters: `m4b-tool chapters --merge input/`
  - Metadata: Embed via m4b-tool or separate tool
- **Error Handling**: Parse stderr for errors, handle gracefully

### 6.4 Plex Integration

- **API Endpoints**:
  - Update library section
  - Refresh metadata
  - Get library sections
- **Library Structure**: Organize files in Plex-compatible structure
- **Notifications**: Send webhook notifications when new items added

---

## 7. Metadata Strategy

### 7.1 Metadata Sources Priority

1. **Audible** (if accessible): Most comprehensive audiobook metadata
2. **Open Library**: Free, open metadata source
3. **Goodreads**: Ratings and additional metadata
4. **Google Books**: Fallback option
5. **User Input**: Manual entry when automatic sources fail

### 7.2 Metadata Fields

- **Core**: Title, Author, Narrator, Publisher, Release Date
- **Audiobook-Specific**: Duration, Format, Bitrate, Language
- **Organization**: Series, Series Position, Genre, Tags
- **Media**: Cover Art, Description, ISBN/ASIN
- **Technical**: File Path, File Size, Processing Date

### 7.3 Metadata Matching

- **Matching Strategies**:
  - Title + Author matching
  - ISBN/ASIN matching (most reliable)
  - Fuzzy matching for variations
  - User confirmation for ambiguous matches

---

## 8. User Interface Design

### 8.1 Main Views ✅

- **Dashboard**: Overview of library, recent activity, download queue
- **Library**: Browse audiobook library, filter by author/series/genre
- **Add Audiobook**: Search and add new audiobooks
- **Download Queue**: Monitor active downloads
- **Processing Queue**: Monitor file processing status
- **Settings**: Configuration for all integrations

### 8.2 Key Features

- **Search**: Search across Jackett indexers
- **Quality Profiles**: Define preferred formats and quality
- **Naming Conventions**: Customize file/folder naming
- **Automatic Management**: Auto-download, auto-process, auto-organize
- **Manual Override**: Manual search, selection, and processing

---

## 9. Development Phases

### Phase 1: Foundation (Weeks 1-4)

- [x] Set up project structure (Go backend, React frontend) ✅
- [x] Create basic API structure ✅
- [x] Set up development environment and tooling ✅
- [x] Create initial Dockerfile and docker-compose.yml ✅
- [x] Implement database schema and models ✅
- [x] Set up database migrations (GORM AutoMigrate) ✅
- [x] Implement basic authentication/authorization (API key) ✅
- [x] Standardize API response format and error handling ✅
- [x] Implement Library API endpoints ✅
- [x] Implement Authors API endpoints ✅
- [x] Implement Books API endpoints ✅
- [x] Implement Downloads API endpoints ✅
- [x] Implement Processing API endpoints ✅
- [x] Implement Search endpoint (basic) ✅
- [ ] Test container deployment locally
- [ ] Set up CI/CD pipeline (optional for MVP)

### Phase 2: Core Integrations (Weeks 5-8)

- [ ] Implement Jackett integration (search functionality)
- [ ] Implement qBittorrent integration (download management)
- [ ] Implement basic download monitoring
- [ ] Create download queue management
- [ ] Test integrations end-to-end

### Phase 3: File Processing (Weeks 9-12)

- [ ] Integrate m4b-tool subprocess execution
- [ ] Implement file detection and format identification
- [ ] Create merge/processing pipeline
- [ ] Implement chapter extraction
- [ ] Add processing queue and status tracking
- [ ] Error handling and retry logic

### Phase 4: Metadata Management (Weeks 13-16)

- [ ] Research and implement metadata source APIs
- [ ] Create metadata matching algorithms
- [ ] Implement metadata embedding
- [ ] Add cover art download and embedding
- [ ] Create metadata editing interface

### Phase 5: Library Management (Weeks 17-20)

- [ ] Implement library organization logic
- [ ] Create file naming and organization system
- [ ] Implement library scanning
- [ ] Add library browsing and filtering
- [ ] Create library management UI

### Phase 6: Plex Integration (Weeks 21-22)

- [ ] Implement Plex API integration
- [ ] Add library update triggers
- [ ] Test Plex library refresh
- [ ] Handle Plex-specific requirements

### Phase 7: Polish & Testing (Weeks 23-26)

- [ ] UI/UX improvements
- [ ] Comprehensive testing
- [ ] Performance optimization
- [ ] Documentation
- [ ] Error handling improvements
- [ ] User feedback integration

---

## 10. Technical Considerations

### 10.1 File Processing Challenges

- **Large Files**: Audiobooks can be very large (several GB)
- **Processing Time**: Merging can take significant time
- **Disk Space**: Need temporary space for processing
- **Error Recovery**: Handle partial processing failures

### 10.2 Metadata Challenges

- **API Limitations**: Some metadata sources have rate limits
- **Matching Accuracy**: Improve matching algorithms over time
- **Missing Data**: Handle cases where metadata is incomplete
- **User Override**: Allow manual metadata correction

### 10.3 Performance Considerations

- **Background Processing**: Process files asynchronously
- **Queue Management**: Prioritize processing queue
- **Caching**: Cache metadata and search results
- **Database Optimization**: Index frequently queried fields

### 10.4 Security Considerations

- **API Keys**: Secure storage of API credentials
- **File System**: Secure file operations
- **User Data**: Protect user library and preferences

### 10.5 Docker-Specific Considerations

- **m4b-tool Installation**: Bundle in image or install at runtime?
- **FFmpeg Dependencies**: Include in base image or install separately?
- **Volume Permissions**: Handle file ownership/permissions correctly
- **Resource Management**: Set appropriate CPU/memory limits
- **Multi-Architecture**: Support ARM64 (Raspberry Pi, Apple Silicon) and AMD64
- **Image Size**: Optimize to keep image size reasonable
- **Build Time**: Optimize Docker build process for faster iterations
- **Development Workflow**: Hot-reload support in development containers

---

## 11. Open Questions & Decisions Needed

### 11.1 Technical Decisions

- [x] **Programming Language**: **Go (Golang)** ✅
- [x] **Web Framework**: **Gin** ✅
- [x] **Database**: **SQLite** ✅ (Start with SQLite, migrate to PostgreSQL if needed)
- [x] **ORM**: **GORM** ✅
- [x] **Background Jobs**: **Native goroutines** ✅ (Start simple, add asynq if needed)
- [x] **Frontend Framework**: **React 18 + TypeScript** ✅
- [x] **Build Tool**: **Vite** ✅
- [x] **UI Library**: **Material-UI** ✅
- [x] **State Management**: **React Query + Zustand** ✅
- [x] **m4b-tool in Docker**: **Install PHAR file directly** ✅
- [x] **Base Image**: **Alpine Linux** ✅
- [x] **Multi-stage Build**: **3 stages** ✅ (frontend, backend, runtime)
- [ ] **Metadata Embedding**: Use m4b-tool or separate tool (FFmpeg/AtomicParsley)?
- [ ] **Chapter Information**: Extract from existing files or generate?
- [x] **Authentication**: **API key authentication** ✅ - Start with API key for MVP, add OAuth later if needed
- [ ] **API Versioning**: How to handle API versioning?

### 11.2 Feature Decisions

- [ ] **Auto-download**: Enable by default or opt-in? (Recommendation: opt-in for MVP)
- [ ] **Quality Profiles**: How granular? (Format only or bitrate too?) (Recommendation: Start with format only)
- [ ] **Series Management**: How to handle multi-book series? (Recommendation: Track series and position)
- [ ] **Duplicate Handling**: How to detect and handle duplicates? (Recommendation: By ISBN/ASIN, then title+author)
- [ ] **Library Organization**: Folder structure? (Recommendation: Author/Series/Book or Author/Book)
- [ ] **File Naming**: Naming convention? (Recommendation: Author - Series - Book Title.m4b)
- [ ] **Processing Queue**: Concurrent processing limit? (Recommendation: 1-2 at a time due to resource usage)
- [ ] **Download Queue**: Priority system? (Recommendation: FIFO for MVP, add priorities later)

### 11.3 Integration Decisions

- [ ] **Audible API**: Official API or scraping? (Legal considerations) (Recommendation: Start with Open Library, add Audible later if needed)
- [ ] **Multiple Download Clients**: Support only qBittorrent or others? (Recommendation: qBittorrent only for MVP, add others later)
- [ ] **Plex Alternatives**: Support other media servers? (Recommendation: Plex only for MVP, add Jellyfin/Emby later)
- [ ] **Metadata Sources Priority**: Which sources to implement first? (Recommendation: Open Library → Google Books → Audible)
- [ ] **Jackett Indexers**: Which indexers to support/test? (Recommendation: Let user configure in Jackett)

---

## 12. Success Criteria

### 12.1 MVP (Minimum Viable Product)

- [ ] Search for audiobooks via Jackett
- [ ] Download via qBittorrent
- [ ] Process files into single m4b with m4b-tool
- [ ] Embed basic metadata
- [ ] Organize files for Plex
- [ ] Basic UI for library management
- [ ] **Docker container runs successfully**
- [ ] **Docker Compose setup for easy deployment**

### 12.2 Full Release

- [ ] All MVP features working reliably
- [ ] Comprehensive metadata management
- [ ] Automatic processing pipeline
- [ ] Polished UI matching Lidarr/Radarr quality
- [ ] Comprehensive documentation
- [ ] Error handling and recovery
- [ ] Performance optimization
- [ ] **Production-ready Docker image**
- [ ] **Multi-architecture support (AMD64, ARM64)**
- [ ] **Docker Hub/GitHub Container Registry publishing**
- [ ] **Docker Compose examples for common setups**

---

## 13. Resources & References

### 13.1 Similar Projects

- [Lidarr](https://github.com/Lidarr/Lidarr) - Music collection manager
- [Radarr](https://github.com/Radarr/Radarr) - Movie collection manager
- [Sonarr](https://github.com/Sonarr/Sonarr) - TV show collection manager

### 13.2 Tools & Libraries

- [m4b-tool](https://github.com/sandreas/m4b-tool) - Audiobook file processing
- [qBittorrent Web API](<https://github.com/qbittorrent/qBittorrent/wiki/Web-API-Reference-(qBittorrent-4.1)>)
- [Jackett API](https://github.com/Jackett/Jackett)
- [Plex API](https://www.plex.tv/api/)

### 13.3 Metadata Sources

- [Open Library API](https://openlibrary.org/developers/api)
- [Google Books API](https://developers.google.com/books)
- [Goodreads API](https://www.goodreads.com/api) (deprecated, may need alternatives)

### 13.4 Docker Resources

- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [Multi-stage Builds](https://docs.docker.com/build/building/multi-stage/)
- [.NET Docker Images](https://hub.docker.com/_/microsoft-dotnet)
- [Lidarr Docker Image](https://hub.docker.com/r/lidarr/lidarr) - Reference implementation
- [Radarr Docker Image](https://hub.docker.com/r/radarr/radarr) - Reference implementation

---

## 14. Next Steps

1. **Review this planning document** and refine requirements
2. **Make key technical decisions** (database, background jobs, Docker base image, etc.)
3. **Set up development environment** and project structure
4. **Create initial Dockerfile and docker-compose.yml** for development
5. **Begin Phase 1** implementation
6. **Create detailed technical specifications** for each module
7. **Set up project management** (issues, milestones, etc.)

---

**Document Version**: 2.0  
**Last Updated**: After Go/React Setup  
**Status**: Ready for Development - Key Decisions Made

## 15. Project Status Summary

### ✅ Completed

- [x] Technology stack decisions (Go, React, Gin, GORM, Material-UI, Vite)
- [x] Project structure created (backend + frontend)
- [x] Docker setup complete (Dockerfile + docker-compose.yml)
- [x] Basic API server structure
- [x] Frontend layout and routing
- [x] Configuration system (viper)
- [x] Database initialization setup

### 🚧 In Progress / Next Steps

- [x] Database models and schema ✅
- [x] API response format standardization ✅
- [ ] API endpoint implementations (Library, Authors, Books)
- [ ] Integration clients (qBittorrent, Jackett, Plex)
- [ ] File processing pipeline
- [ ] Metadata fetching and embedding

### ❓ Decisions Needed Before Starting Development

See Section 11 for detailed decisions needed.
