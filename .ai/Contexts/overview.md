# Project Overview

## Project Description

**Listenarr** is an audiobook collection manager similar to Lidarr (music) and Radarr (movies), designed to automate the process of searching, downloading, processing, and organizing audiobooks. It integrates with qBittorrent, Jackett, and Plex, and uses m4b-tool to merge audiobook files into single m4b files with embedded metadata.

## Project Areas

This document outlines the main areas of the codebase and maps contexts to corresponding parts of the project.

### Context Mapping

- **Overview** (`overview.md`): This file - high-level project structure and context mapping
- **Planning** (`../PLANNING.md`): Comprehensive project planning document
- **Decisions Needed** (`../DECISIONS_NEEDED.md`): Key decisions that need to be made before starting development
- **Language Comparison** (`../LANGUAGE_COMPARISON.md`): Detailed comparison of programming language options
- **Go** (`go.md`): Go implementation details, project structure, and best practices
- **Auth** (`auth.md`): Authentication system (API key authentication)
- **Database** (`database.md`): Database models, relationships, and schema
- **Testing** (`testing.md`): Testing standards, requirements, and best practices
- **Docker** (`docker.md`): Docker containerization architecture and configuration
- **m4b-tool** (`m4b-tool.md`): m4b-tool installation methods, dependencies, and Docker integration
- Additional context files will be added as the project evolves:
  - `search.md`: Search module and Jackett integration
  - `download.md`: Download module and qBittorrent integration
  - `processing.md`: File processing module and m4b-tool integration
  - `metadata.md`: Metadata fetching and embedding
  - `library.md`: Library organization and management
  - `plex.md`: Plex integration module
  - `api.md`: API endpoints and routes
  - `ui.md`: User interface components and views
  - `database.md`: Database schema and models
  - `config.md`: Configuration and environment setup

## Technology Stack

- **Backend**: **Go (Golang)** ✅
  - **Framework**: Gin (web framework)
  - **Database**: SQLite with GORM or sqlx
  - **Background Jobs**: Native goroutines/channels or asynq
  - **HTTP Client**: Built-in net/http or resty
  - **Config**: viper
  - **Logging**: logrus or zap
- **Frontend**: React, Material-UI
- **Database**: SQLite (simple) or PostgreSQL (if needed)
- **Integrations**: qBittorrent API, Jackett API, Plex API, m4b-tool

## Project Structure

_To be populated as the project develops_

## Key Components

### Core Modules
1. **Search Module**: Interfaces with Jackett to search for audiobooks
2. **Download Module**: Manages downloads via qBittorrent
3. **Processing Module**: Uses m4b-tool to merge and process files
4. **Metadata Module**: Fetches and embeds audiobook metadata
5. **Library Module**: Organizes and manages audiobook library
6. **Plex Integration Module**: Integrates with Plex Media Server

## Architecture

See `PLANNING.md` for detailed architecture and workflow design.

