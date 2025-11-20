# m4b-tool Integration

## Overview

m4b-tool is a PHP-based command-line utility for merging, splitting, and chapterizing audiobook files. This document covers installation methods, dependencies, and integration strategies for Docker containers.

## What is m4b-tool?

- **Language**: PHP-based application
- **Purpose**: Merge multiple audio files into single M4B files with chapters
- **Dependencies**: PHP runtime, FFmpeg, and various PHP extensions
- **Official Repository**: https://github.com/sandreas/m4b-tool
- **Official Docker Image**: `sandreas/m4b-tool:latest` on Docker Hub

## Installation Methods for Docker

### Method 1: Use Official Docker Image (Recommended for Separate Container)

The official m4b-tool Docker image is maintained by the project author and includes all dependencies.

**Pros:**
- Pre-built and maintained
- Includes all dependencies (PHP, FFmpeg, etc.)
- Regularly updated
- No build time required

**Cons:**
- Requires running as separate container or sidecar
- Additional container overhead
- Need to manage inter-container communication

**Usage Example:**
```bash
docker run -it --rm \
  -u $(id -u):$(id -g) \
  -v "$(pwd)":/mnt \
  sandreas/m4b-tool:latest \
  merge input/ --output-file output.m4b
```

**In Listenarr Context:**
- Could run as sidecar container
- Communicate via shared volume
- More complex orchestration

### Method 2: Install m4b-tool in Listenarr Container (Recommended)

Install m4b-tool directly in the Listenarr container alongside the .NET application.

**Pros:**
- Single container deployment
- Simpler architecture
- Direct subprocess execution
- No inter-container communication needed

**Cons:**
- Larger image size
- Need to manage PHP and dependencies
- Build time increases

**Implementation Approaches:**

#### 2a. Install PHP + m4b-tool via Package Manager
```dockerfile
# Install PHP and dependencies
RUN apt-get update && \
    apt-get install -y \
    php-cli \
    php-mbstring \
    php-xml \
    php-zip \
    php-curl \
    ffmpeg \
    && apt-get clean

# Download and install m4b-tool
RUN curl -L -o /usr/local/bin/m4b-tool \
    https://github.com/sandreas/m4b-tool/releases/latest/download/m4b-tool.phar && \
    chmod +x /usr/local/bin/m4b-tool
```

#### 2b. Use PHP Composer (if m4b-tool is available via Composer)
```dockerfile
# Install PHP and Composer
RUN apt-get update && \
    apt-get install -y php-cli php-mbstring php-xml php-zip php-curl composer ffmpeg && \
    apt-get clean

# Install m4b-tool via Composer
RUN composer global require sandreas/m4b-tool
```

#### 2c. Download Pre-built PHAR File
```dockerfile
# Install PHP runtime and FFmpeg
RUN apt-get update && \
    apt-get install -y \
    php-cli \
    php-mbstring \
    php-xml \
    php-zip \
    php-curl \
    ffmpeg \
    curl \
    && apt-get clean

# Download m4b-tool PHAR
RUN curl -L -o /usr/local/bin/m4b-tool \
    https://github.com/sandreas/m4b-tool/releases/latest/download/m4b-tool.phar && \
    chmod +x /usr/local/bin/m4b-tool && \
    m4b-tool --version
```

### Method 3: Build from Source

Clone and build m4b-tool from source in the Docker image.

**Pros:**
- Full control over version
- Can customize build

**Cons:**
- Longest build time
- Most complex
- Requires Composer and build tools

**Implementation:**
```dockerfile
# Install build dependencies
RUN apt-get update && \
    apt-get install -y \
    git \
    php-cli \
    php-mbstring \
    php-xml \
    php-zip \
    php-curl \
    composer \
    ffmpeg \
    && apt-get clean

# Clone and build m4b-tool
RUN git clone https://github.com/sandreas/m4b-tool.git /tmp/m4b-tool && \
    cd /tmp/m4b-tool && \
    composer install --no-dev --optimize-autoloader && \
    php bin/build.php && \
    cp dist/m4b-tool.phar /usr/local/bin/m4b-tool && \
    chmod +x /usr/local/bin/m4b-tool && \
    rm -rf /tmp/m4b-tool
```

## PHP Requirements

m4b-tool requires:
- **PHP**: 7.4+ or 8.0+ (check latest requirements)
- **PHP Extensions**:
  - `mbstring` - String handling
  - `xml` - XML processing
  - `zip` - Archive handling
  - `curl` - HTTP requests
  - `json` - JSON processing (usually included)

## FFmpeg Requirements

- **FFmpeg**: Required for audio processing
- **Version**: Latest stable recommended
- **Codecs**: Must support AAC, MP3, M4A, M4B formats

## Recommended Approach for Listenarr

### Primary Recommendation: Method 2c (PHAR Download)

**Rationale:**
1. **Simplicity**: Single container, no sidecar complexity
2. **Performance**: Direct subprocess execution, no network overhead
3. **Reliability**: No dependency on external container availability
4. **Size**: Reasonable image size increase (~100-200MB for PHP + FFmpeg)
5. **Maintenance**: Easy to update by changing download URL

**Dockerfile Example:**
```dockerfile
FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS runtime

# Install PHP, FFmpeg, and dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    php-cli \
    php-mbstring \
    php-xml \
    php-zip \
    php-curl \
    ffmpeg \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Download and install m4b-tool
RUN curl -L -o /usr/local/bin/m4b-tool \
    https://github.com/sandreas/m4b-tool/releases/latest/download/m4b-tool.phar && \
    chmod +x /usr/local/bin/m4b-tool && \
    m4b-tool --version

# Verify installation
RUN m4b-tool --version

# ... rest of Dockerfile
```

## Execution from .NET

### Subprocess Execution
```csharp
var processInfo = new ProcessStartInfo
{
    FileName = "m4b-tool",
    Arguments = "merge input/ --output-file output.m4b",
    UseShellExecute = false,
    RedirectStandardOutput = true,
    RedirectStandardError = true,
    WorkingDirectory = workingDir
};

using var process = Process.Start(processInfo);
await process.WaitForExitAsync();
```

### Path Considerations
- Ensure `/usr/local/bin` is in PATH
- Or use full path: `/usr/local/bin/m4b-tool`
- Test execution in container during build

## Version Management

### Pinning Versions
- Use specific release URL instead of `latest`
- Example: `https://github.com/sandreas/m4b-tool/releases/download/v0.5.1/m4b-tool.phar`
- Update version in Dockerfile or via build argument

### Build Arguments
```dockerfile
ARG M4B_TOOL_VERSION=latest
ARG M4B_TOOL_URL=https://github.com/sandreas/m4b-tool/releases/${M4B_TOOL_VERSION}/download/m4b-tool.phar

RUN curl -L -o /usr/local/bin/m4b-tool ${M4B_TOOL_URL} && \
    chmod +x /usr/local/bin/m4b-tool
```

## Testing Installation

### Verify in Container
```bash
docker run --rm listenarr/listenarr:latest m4b-tool --version
```

### Test Merge Command
```bash
docker run --rm \
  -v /path/to/audiobook:/data \
  listenarr/listenarr:latest \
  m4b-tool merge /data/input/ --output-file /data/output.m4b
```

## Troubleshooting

### Common Issues

1. **m4b-tool not found**
   - Check PATH includes `/usr/local/bin`
   - Verify file exists: `ls -la /usr/local/bin/m4b-tool`
   - Check file permissions: `chmod +x /usr/local/bin/m4b-tool`

2. **PHP errors**
   - Verify PHP extensions installed
   - Check PHP version compatibility
   - Review m4b-tool requirements

3. **FFmpeg errors**
   - Verify FFmpeg installed and in PATH
   - Check FFmpeg version supports required codecs
   - Test: `ffmpeg -version`

4. **Permission errors**
   - Ensure proper file permissions on volumes
   - Check PUID/PGID settings
   - Verify umask configuration

## Image Size Impact

### Estimated Size Increases
- **PHP CLI**: ~50-100MB
- **FFmpeg**: ~50-100MB
- **m4b-tool PHAR**: ~5-10MB
- **Total**: ~100-200MB additional

### Optimization Strategies
- Use Alpine-based PHP image (smaller)
- Remove unnecessary PHP extensions
- Use multi-stage build to minimize final image
- Consider using official m4b-tool image as base (if acceptable)

## Alternative: Alpine-based Installation

For smaller image size:
```dockerfile
FROM mcr.microsoft.com/dotnet/aspnet:8.0-alpine AS runtime

# Install PHP and FFmpeg on Alpine
RUN apk add --no-cache \
    php-cli \
    php-mbstring \
    php-xml \
    php-zip \
    php-curl \
    ffmpeg \
    curl

# Download m4b-tool
RUN curl -L -o /usr/local/bin/m4b-tool \
    https://github.com/sandreas/m4b-tool/releases/latest/download/m4b-tool.phar && \
    chmod +x /usr/local/bin/m4b-tool
```

**Note**: Alpine may have compatibility issues with some .NET applications. Test thoroughly.

## References

- [m4b-tool GitHub](https://github.com/sandreas/m4b-tool)
- [m4b-tool Docker Hub](https://hub.docker.com/r/sandreas/m4b-tool)
- [PHP Docker Images](https://hub.docker.com/_/php)
- [FFmpeg Documentation](https://ffmpeg.org/documentation.html)

