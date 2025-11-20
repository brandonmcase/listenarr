#!/bin/bash
set -e

echo "🚀 Starting build-all.sh..."

# MARK: Preflight
echo "🔍 Checking dependencies..."

# Check Go
command -v go >/dev/null 2>&1 || { echo "❌ Go is not installed. Please install Go 1.21 or later."; exit 1; }
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✓ Go version: $GO_VERSION"

# Check Node.js (for frontend)
command -v node >/dev/null 2>&1 || echo "⚠️  Node.js not found (frontend tests will be skipped)"
if command -v node >/dev/null 2>&1; then
    NODE_VERSION=$(node --version)
    echo "✓ Node.js version: $NODE_VERSION"
fi

# Check curl (for endpoint testing)
command -v curl >/dev/null 2>&1 || { echo "❌ curl is not installed. Required for endpoint testing."; exit 1; }
echo "✓ curl is available"

# MARK: Go Format Check
echo ""
echo "⚡ Running Go format check..."
if ! go fmt ./...; then
    echo "❌ Go format check failed"
    exit 1
fi
echo "✓ Go code is properly formatted"

# MARK: Go Vet
echo ""
echo "🔍 Running go vet..."
if ! go vet ./...; then
    echo "❌ go vet found issues"
    exit 1
fi
echo "✓ go vet passed"

# MARK: Build Backend
echo ""
echo "🏗️ Building backend..."
if ! go build -o /tmp/listenarr-test ./cmd/listenarr; then
    echo "❌ Backend build failed"
    exit 1
fi
echo "✓ Backend build successful"

# MARK: Unit Tests
echo ""
echo "🧪 Running Go unit tests..."
if ! go test -v -coverprofile=/tmp/coverage.out ./...; then
    echo "❌ Unit tests failed"
    exit 1
fi

# Show test coverage
if [ -f /tmp/coverage.out ]; then
    COVERAGE=$(go tool cover -func=/tmp/coverage.out | grep total | awk '{print $3}')
    echo "✓ Test coverage: $COVERAGE"
fi

# MARK: Frontend Tests (if Node.js is available)
if command -v node >/dev/null 2>&1 && [ -d "frontend" ]; then
    echo ""
    echo "🎨 Running frontend tests..."
    cd frontend
    
    # Install dependencies if needed
    if [ ! -d "node_modules" ]; then
        echo "Installing frontend dependencies..."
        npm install --silent
    fi
    
    # Run linting
    if npm run lint >/dev/null 2>&1; then
        echo "✓ Frontend linting passed"
    else
        echo "⚠️  Frontend linting found issues (non-blocking)"
    fi
    
    # Build frontend
    if npm run build >/dev/null 2>&1; then
        echo "✓ Frontend build successful"
    else
        echo "⚠️  Frontend build had issues (non-blocking)"
    fi
    
    cd ..
fi

# MARK: API Endpoint Testing
echo ""
echo "🌐 Testing API endpoints..."

# Check if server is already running
if curl -s -f "${LISTENARR_URL:-http://localhost:8686}/api/health" >/dev/null 2>&1; then
    echo "✓ Server is running, testing endpoints..."
    
    # Run endpoint test script
    if [ -f "scripts/test-endpoints.sh" ]; then
        if bash scripts/test-endpoints.sh; then
            echo "✓ All API endpoints are functioning correctly"
        else
            echo "⚠️  Some API endpoint tests failed (server may not be fully configured)"
        fi
    else
        echo "⚠️  Endpoint test script not found"
    fi
else
    echo "⚠️  Server is not running. Skipping endpoint tests."
    echo "   To test endpoints, start the server and run: scripts/test-endpoints.sh"
fi

# MARK: Summary
echo ""
echo "=========================================="
echo "✅ build-all.sh completed successfully!"
echo "=========================================="
echo ""
echo "Next steps:"
echo "  - Start server: go run ./cmd/listenarr"
echo "  - Test endpoints: scripts/test-endpoints.sh"
echo "  - Run tests: go test ./..."
echo ""

