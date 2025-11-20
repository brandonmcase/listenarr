# Testing Standards

## Overview

Listenarr follows strict testing requirements. All code must have tests, and all tests must pass before code can be committed or merged.

## Test Requirements

### Mandatory Testing
- **All packages must have tests** - Every `.go` file should have a corresponding `*_test.go` file
- **All functions must be testable** - Write testable code with clear interfaces
- **Tests must pass** - No exceptions, tests are a quality gate
- **Test coverage** - Aim for >80% coverage on critical paths

### Test Types

#### Unit Tests
- Test individual functions and methods in isolation
- Mock external dependencies
- Fast execution (< 1 second per test)
- Located in `*_test.go` files alongside code

#### Integration Tests
- Test API endpoints and external integrations
- May require running services (database, external APIs)
- Can be in `*_integration_test.go` files
- Use `-tags=integration` build tag if needed

#### End-to-End Tests
- Test complete workflows
- Use API endpoint testing scripts
- May require full server setup

## Test Structure

### Go Test Files
```go
package packagename

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestFunctionName(t *testing.T) {
    // Arrange
    input := "test"
    
    // Act
    result := FunctionName(input)
    
    // Assert
    assert.Equal(t, "expected", result)
}
```

### Table-Driven Tests
```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"test case 1", "input1", "output1"},
        {"test case 2", "input2", "output2"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := FunctionName(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

## Running Tests

### Basic Commands
```bash
# Run all tests
go test ./...

# Run tests in specific package
go test ./internal/auth

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Using Make
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests with verbose output
make test-verbose
```

## Test Best Practices

### 1. Test Naming
- Use descriptive test names: `TestFunctionName_Scenario_ExpectedResult`
- Example: `TestAPIKeyMiddleware_InvalidKey_Returns401`

### 2. Test Organization
- One test per scenario
- Use subtests with `t.Run()` for related scenarios
- Group related tests together

### 3. Assertions
- Use `require` for critical assertions (stops test on failure)
- Use `assert` for non-critical checks (continues test)
- Provide clear error messages

### 4. Test Data
- Use test fixtures for complex data
- Create helper functions for common setup
- Clean up test data after tests

### 5. Mocking
- Mock external dependencies (databases, APIs, file system)
- Use interfaces to enable mocking
- Keep mocks simple and focused

### 6. Test Isolation
- Tests should not depend on each other
- Tests should not depend on execution order
- Clean up state between tests

## API Endpoint Testing

### Manual Testing
```bash
# Test health endpoint
curl http://localhost:8686/api/health

# Test protected endpoint
curl -H "X-API-Key: your-key" http://localhost:8686/api/v1/library
```

### Automated Testing
```bash
# Run endpoint test script
make test-endpoints
# or
bash scripts/test-endpoints.sh
```

## Test Coverage Goals

- **Critical paths**: >90% coverage
  - Authentication
  - Configuration loading
  - Database operations
  - API handlers

- **Business logic**: >80% coverage
  - Service layer
  - Processing logic
  - Metadata matching

- **Utilities**: >70% coverage
  - Helper functions
  - Utility packages

## Continuous Integration

Tests should run:
- Before every commit (pre-commit hook recommended)
- In CI/CD pipeline
- As part of `build-all.sh`

## Common Test Patterns

### Testing HTTP Handlers
```go
func TestHandler(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := gin.New()
    router.GET("/test", handler)
    
    req, _ := http.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
}
```

### Testing with Database
```go
func TestWithDB(t *testing.T) {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)
    defer db.Close()
    
    // Test database operations
}
```

### Testing Configuration
```go
func TestConfig(t *testing.T) {
    os.Setenv("TEST_VAR", "value")
    defer os.Unsetenv("TEST_VAR")
    
    cfg, err := config.Load()
    require.NoError(t, err)
    assert.Equal(t, "value", cfg.TestVar)
}
```

## Resources

- [Go Testing Documentation](https://go.dev/doc/tutorial/add-a-test)
- [testify Documentation](https://github.com/stretchr/testify)
- [Testing Best Practices](https://golang.org/doc/effective_go#testing)

