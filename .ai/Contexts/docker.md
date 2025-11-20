# Docker Containerization

## Overview

Listenarr is designed to run easily in Docker containers, following the same patterns as Lidarr and Radarr. This document tracks Docker-specific architecture, configuration, and implementation details.

## Container Architecture

### Base Image
- **Runtime**: Microsoft .NET 8.0 ASP.NET Runtime image
- **Build**: Multi-stage build with .NET SDK for compilation
- **Size Optimization**: Use Alpine-based images where possible

### Dependencies in Container
- **.NET Runtime**: Provided by base image
- **m4b-tool**: PHP-based tool, needs PHP runtime (see `m4b-tool.md` for installation details)
- **FFmpeg**: Required for audio processing
- **PHP**: PHP CLI with extensions (mbstring, xml, zip, curl)
- **System Libraries**: All dependencies for m4b-tool and FFmpeg

### Volume Mounts
- `/config`: Application configuration and SQLite database
- `/library`: Audiobook library directory (read/write)
- `/downloads`: Temporary download location (optional, can use library)
- `/processing`: Temporary processing workspace (optional, can use library)

### Ports
- **8686**: Web UI and API (default, configurable)

### Environment Variables
- `PUID`: User ID for file permissions
- `PGID`: Group ID for file permissions
- `TZ`: Timezone (e.g., `America/New_York`)
- `UMASK`: File permission mask (default: `002`)

## Build Process

### Multi-Stage Build
1. **Frontend Build Stage**: Build React frontend
2. **Backend Build Stage**: Compile .NET application
3. **Runtime Stage**: Combine artifacts with runtime dependencies

### Build Considerations
- Minimize image size
- Cache dependencies effectively
- Include health check
- Set proper file permissions
- Use non-root user when possible

## Runtime Considerations

### File Permissions
- Must handle file permissions correctly for mounted volumes
- Use PUID/PGID to match host user
- Set appropriate umask for file creation

### Resource Management
- Set CPU/memory limits if needed
- Handle large file processing (audiobooks can be GB in size)
- Monitor disk space for processing workspace

### Health Checks
- Implement `/api/health` endpoint
- Check database connectivity
- Verify critical services are running

### Logging
- Log to stdout/stderr for Docker log collection
- Structured logging (JSON format preferred)
- Log levels configurable via environment or config

## Development Workflow

### Local Development
- Use `docker-compose.yml` for local development
- Mount source code for hot-reload (development only)
- Separate development and production Dockerfiles if needed

### Testing
- Test container builds in CI/CD
- Test on multiple architectures (AMD64, ARM64)
- Verify volume mounts work correctly
- Test with different permission configurations

## Deployment

### Docker Hub / Container Registry
- Automated builds on release
- Tag versions appropriately
- Provide `latest` tag for latest stable release
- Multi-architecture builds (AMD64, ARM64)

### Docker Compose Examples
- Basic setup
- With qBittorrent and Jackett
- Full stack (Listenarr + qBittorrent + Jackett + Plex)

## Platform Support

### Architectures
- **AMD64**: Primary architecture
- **ARM64**: Support for Raspberry Pi, Apple Silicon, etc.

### Operating Systems
- Linux (primary)
- Windows (via WSL2)
- macOS (via Docker Desktop)

## Integration with Other Containers

### qBittorrent
- Can run in separate container
- Connect via network (bridge or host)
- Share download directory via volume

### Jackett
- Typically runs in separate container
- Connect via HTTP API
- No shared volumes needed

### Plex
- Typically runs in separate container or host
- Listenarr updates Plex library via API
- Library directory may be shared or separate

## Troubleshooting

### Common Issues
- **Permission Errors**: Check PUID/PGID match host user
- **m4b-tool Not Found**: Verify installation in Dockerfile
- **Large File Processing**: Ensure sufficient disk space
- **Network Issues**: Check network mode and port mappings

## Future Considerations

- [ ] Support for Docker Swarm
- [ ] Kubernetes manifests
- [ ] Helm charts
- [ ] Podman support
- [ ] Rootless container support

