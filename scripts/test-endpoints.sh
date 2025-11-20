#!/bin/bash
# API Endpoint Verification Script
# Systematically tests all API endpoints to ensure server is running and functioning

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
BASE_URL="${LISTENARR_URL:-http://localhost:8686}"
API_KEY="${LISTENARR_API_KEY:-}"

# Counters
PASSED=0
FAILED=0
SKIPPED=0

# Function to print test result
print_result() {
    local status=$1
    local endpoint=$2
    local message=$3
    
    if [ "$status" = "PASS" ]; then
        echo -e "${GREEN}✓${NC} $endpoint - $message"
        ((PASSED++))
    elif [ "$status" = "FAIL" ]; then
        echo -e "${RED}✗${NC} $endpoint - $message"
        ((FAILED++))
    else
        echo -e "${YELLOW}⊘${NC} $endpoint - $message"
        ((SKIPPED++))
    fi
}

# Function to test an endpoint
test_endpoint() {
    local method=$1
    local endpoint=$2
    local expected_status=$3
    local description=$4
    local data=$5
    
    local url="${BASE_URL}${endpoint}"
    local headers=()
    
    # Add API key if provided
    if [ -n "$API_KEY" ]; then
        headers+=(-H "X-API-Key: $API_KEY")
    fi
    
    # Add Content-Type for POST/PUT requests
    if [ "$method" = "POST" ] || [ "$method" = "PUT" ]; then
        headers+=(-H "Content-Type: application/json")
    fi
    
    # Make request
    if [ -n "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" "${headers[@]}" -d "$data" "$url" 2>&1)
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "${headers[@]}" "$url" 2>&1)
    fi
    
    # Extract status code (last line)
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    # Check if curl failed
    if [ $? -ne 0 ]; then
        print_result "FAIL" "$method $endpoint" "Connection failed"
        return 1
    fi
    
    # Check status code
    if [ "$http_code" = "$expected_status" ]; then
        print_result "PASS" "$method $endpoint" "$description"
        return 0
    else
        print_result "FAIL" "$method $endpoint" "Expected $expected_status, got $http_code"
        echo "  Response: $body"
        return 1
    fi
}

# Function to get API key from config file
get_api_key_from_config() {
    local config_file="${CONFIG_PATH:-./config}/config.yml"
    if [ -f "$config_file" ]; then
        # Try to extract API key using grep/sed
        grep -A 1 "auth:" "$config_file" | grep "api_key:" | sed 's/.*api_key:[[:space:]]*"\(.*\)".*/\1/' | head -1
    fi
}

# Main execution
echo "=========================================="
echo "Listenarr API Endpoint Verification"
echo "=========================================="
echo "Base URL: $BASE_URL"
echo ""

# Try to get API key if not provided
if [ -z "$API_KEY" ]; then
    API_KEY=$(get_api_key_from_config)
    if [ -n "$API_KEY" ]; then
        echo "Using API key from config file"
    else
        echo -e "${YELLOW}Warning: No API key provided. Some tests will be skipped.${NC}"
        echo "Set LISTENARR_API_KEY environment variable or ensure config file exists."
        echo ""
    fi
fi

echo "Starting endpoint tests..."
echo ""

# Test 1: Health Check (no auth required)
echo "=== Public Endpoints ==="
test_endpoint "GET" "/api/health" "200" "Health check endpoint"

echo ""
echo "=== Protected Endpoints ==="

if [ -z "$API_KEY" ]; then
    echo -e "${YELLOW}Skipping protected endpoints - no API key${NC}"
    SKIPPED=$((SKIPPED + 7))
else
    # Test 2: Get Library
    test_endpoint "GET" "/api/v1/library" "200" "Get library items"
    
    # Test 3: Add to Library
    test_endpoint "POST" "/api/v1/library" "200" "Add item to library" '{"title":"Test Book","author":"Test Author"}'
    
    # Test 4: Get Downloads
    test_endpoint "GET" "/api/v1/downloads" "200" "Get download queue"
    
    # Test 5: Start Download
    test_endpoint "POST" "/api/v1/downloads" "200" "Start download" '{"libraryItemId":"test","releaseId":"test"}'
    
    # Test 6: Get Processing Queue
    test_endpoint "GET" "/api/v1/processing" "200" "Get processing queue"
    
    # Test 7: Search Audiobooks
    test_endpoint "GET" "/api/v1/search?q=test" "200" "Search audiobooks"
    
    # Test 8: Remove from Library
    test_endpoint "DELETE" "/api/v1/library/test-id" "200" "Remove item from library"
    
    # Test 9: Test authentication failure
    echo ""
    echo "=== Authentication Tests ==="
    # Test with invalid key
    invalid_url="${BASE_URL}/api/v1/library"
    invalid_response=$(curl -s -w "\n%{http_code}" -H "X-API-Key: invalid-key" "$invalid_url" 2>&1)
    invalid_code=$(echo "$invalid_response" | tail -n1)
    if [ "$invalid_code" = "401" ]; then
        print_result "PASS" "GET /api/v1/library (invalid key)" "Reject invalid API key"
    else
        print_result "FAIL" "GET /api/v1/library (invalid key)" "Expected 401, got $invalid_code"
    fi
fi

# Test invalid endpoint
echo ""
echo "=== Error Handling ==="
test_endpoint "GET" "/api/v1/nonexistent" "404" "Handle 404 for invalid endpoint"

# Summary
echo ""
echo "=========================================="
echo "Test Summary"
echo "=========================================="
echo -e "${GREEN}Passed:${NC} $PASSED"
echo -e "${RED}Failed:${NC} $FAILED"
echo -e "${YELLOW}Skipped:${NC} $SKIPPED"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed!${NC}"
    exit 1
fi

