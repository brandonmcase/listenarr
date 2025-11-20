# API Design and Implementation

## Overview

Listenarr uses a RESTful API built with Gin framework. All API endpoints follow a consistent response format and error handling pattern.

## API Response Format

### Standard Response Structure

All API responses follow this structure:

```json
{
  "success": true,
  "data": { ... },
  "error": null
}
```

For errors:
```json
{
  "success": false,
  "data": null,
  "error": "Error message describing what went wrong"
}
```

### Response Types

- **Success Response**: `{ success: true, data: T, error: null }`
- **Error Response**: `{ success: false, data: null, error: string }`
- **Paginated Response**: `{ success: true, data: T[], pagination: { page, limit, total, totalPages }, error: null }`

## HTTP Status Codes

### Standard Status Codes

- **200 OK**: Successful GET, PUT, PATCH requests
- **201 Created**: Successful POST requests that create resources
- **204 No Content**: Successful DELETE requests
- **400 Bad Request**: Invalid request parameters or body
- **401 Unauthorized**: Missing or invalid API key
- **404 Not Found**: Resource not found
- **409 Conflict**: Resource conflict (e.g., duplicate entry)
- **422 Unprocessable Entity**: Validation errors
- **500 Internal Server Error**: Server-side errors

## API Versioning

- **Base Path**: `/api/v1`
- **Version Strategy**: URL-based versioning (`/api/v1`, `/api/v2`, etc.)
- **Future**: Consider header-based versioning for backward compatibility

## Authentication

- **Method**: API key authentication
- **Header**: `X-API-Key: <api-key>`
- **Query Parameter**: `?apikey=<api-key>`
- **Public Endpoints**: `/api/health` (no authentication required)
- **Protected Endpoints**: All `/api/v1/*` endpoints require valid API key

See `auth.md` for detailed authentication documentation.

## Endpoints

### Implemented Endpoints

- `GET /api/health` - Health check (public) ✅
- `GET /api/v1/library` - List library items (with pagination, filtering, sorting) ✅
- `GET /api/v1/library/:id` - Get single library item with full details ✅
- `POST /api/v1/library` - Add book to library (creates Author, Book, Series if needed) ✅
- `DELETE /api/v1/library/:id` - Remove from library (soft delete) ✅

### Placeholder Endpoints

- `GET /api/v1/downloads` - List downloads
- `POST /api/v1/downloads` - Start download
- `GET /api/v1/processing` - Get processing queue
- `GET /api/v1/search` - Search audiobooks

### Planned Endpoints

#### Library Management ✅
- `GET /api/v1/library` - List library items (with pagination, filtering, sorting) ✅
- `GET /api/v1/library/:id` - Get single library item with full details ✅
- `POST /api/v1/library` - Add book to library ✅
- `PUT /api/v1/library/:id` - Update library item (planned)
- `DELETE /api/v1/library/:id` - Remove from library (soft delete) ✅

#### Authors
- `GET /api/v1/authors` - List authors
- `GET /api/v1/authors/:id` - Get author with books
- `POST /api/v1/authors` - Create author
- `PUT /api/v1/authors/:id` - Update author
- `DELETE /api/v1/authors/:id` - Delete author

#### Books
- `GET /api/v1/books` - List books
- `GET /api/v1/books/:id` - Get book with full details
- `POST /api/v1/books` - Create book
- `PUT /api/v1/books/:id` - Update book
- `DELETE /api/v1/books/:id` - Delete book

#### Downloads
- `GET /api/v1/downloads` - List downloads (with filtering by status)
- `GET /api/v1/downloads/:id` - Get download details
- `POST /api/v1/downloads` - Start download
- `DELETE /api/v1/downloads/:id` - Cancel download

#### Processing
- `GET /api/v1/processing` - Get processing queue
- `GET /api/v1/processing/:id` - Get processing task details
- `POST /api/v1/processing/:id/retry` - Retry failed processing

#### Search
- `GET /api/v1/search` - Search for audiobooks (will integrate with Jackett)

## Error Handling

### Error Response Format

```json
{
  "success": false,
  "data": null,
  "error": "Human-readable error message",
  "code": "ERROR_CODE",  // Optional: machine-readable error code
  "details": { ... }      // Optional: additional error details
}
```

### Error Types

- **Validation Errors** (422): Invalid input data
- **Not Found** (404): Resource doesn't exist
- **Conflict** (409): Resource conflict (duplicate, etc.)
- **Unauthorized** (401): Authentication required
- **Internal Error** (500): Server-side errors

## Pagination

### Pagination Parameters

- `page`: Page number (default: 1)
- `limit`: Items per page (default: 20, max: 100)
- `offset`: Alternative to page (calculated from page and limit)

### Paginated Response

```json
{
  "success": true,
  "data": [ ... ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "totalPages": 8
  },
  "error": null
}
```

## Filtering and Sorting

### Filtering

- Query parameters for filtering: `?status=available&author_id=1`
- Multiple values: `?status=available,wanted`
- Date ranges: `?added_after=2024-01-01&added_before=2024-12-31`

### Sorting

- Query parameter: `?sort=title&order=asc` or `?sort=-created_at` (desc)
- Multiple fields: `?sort=author.name,title`
- Default: `created_at DESC`

## Request Validation

- Validate all input data
- Return 422 with validation errors for invalid input
- Use struct tags for validation (e.g., `validate:"required,min=1"`)

## Response Helpers (To Be Implemented)

### Success Response
```go
func SuccessResponse(c *gin.Context, statusCode int, data interface{}) {
    c.JSON(statusCode, gin.H{
        "success": true,
        "data":    data,
        "error":   nil,
    })
}
```

### Error Response
```go
func ErrorResponse(c *gin.Context, statusCode int, err error) {
    c.JSON(statusCode, gin.H{
        "success": false,
        "data":    nil,
        "error":   err.Error(),
    })
}
```

## Testing

- All endpoints must have tests
- Test success and error cases
- Test authentication requirements
- Test validation
- Use `test-endpoints.sh` or `test-endpoints.go` for integration testing

## Documentation

- API documentation will be generated using Swagger/OpenAPI
- Endpoint documentation includes:
  - Request/response examples
  - Error responses
  - Authentication requirements
  - Query parameters

## Implementation Status

- [x] API server structure
- [x] Authentication middleware
- [x] Route setup
- [x] Response format helpers ✅
- [x] Error handling helpers ✅
- [x] HTTP status code constants ✅
- [x] Pagination helpers ✅
- [x] Library endpoints implementation ✅
- [x] Request validation (Gin binding) ✅
- [ ] Authors endpoints
- [ ] Books endpoints
- [ ] API documentation (Swagger)

