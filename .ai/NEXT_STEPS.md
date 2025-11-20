# Next Steps for Listenarr Development

## Current Status

✅ **Completed:**
- Project structure (Go backend + React frontend)
- API key authentication
- Testing infrastructure
- Docker setup
- Build system
- Documentation
- **Database models and schema** ✅ (Author, Book, Audiobook, Series, Release, Download, ProcessingTask, LibraryItem)
- **Database migrations** ✅ (GORM AutoMigrate with indexes)

## Immediate Next Steps (Priority Order)

### 1. Database Models & Schema ✅ **COMPLETE**

**Status**: All models created, relationships defined, migrations set up, tests passing.

**Completed:**
- ✅ All 8 models created with GORM tags
- ✅ All relationships defined and tested
- ✅ GORM AutoMigrate configured
- ✅ Indexes added (single-field, composite, foreign keys)
- ✅ Comprehensive unit tests (9 test cases, all passing)

### 2. Standardize API Response Format ✅ **COMPLETE**

**Status**: All response helpers, error handling, and status codes implemented and tested.

**Completed:**
- ✅ Response helper functions (`SuccessResponse`, `ErrorResponse`, `CreatedResponse`, `NoContentResponse`)
- ✅ Pagination helper (`PaginatedSuccessResponse`)
- ✅ Error types (`APIError`, `ValidationErrors`) with error codes
- ✅ HTTP status code constants
- ✅ Comprehensive tests (all passing)
- ✅ Updated placeholder endpoints to use new format

### 3. Implement Core API Endpoints

**Why**: Need working endpoints to build frontend and integrations.

**Priority Order:**
1. **Library Endpoints** ✅ **COMPLETE**:
   - [x] `GET /api/v1/library` - List all library items (with pagination, filtering, sorting) ✅
   - [x] `GET /api/v1/library/:id` - Get single library item ✅
   - [x] `POST /api/v1/library` - Add book to library (creates Author, Book, Series if needed) ✅
   - [x] `DELETE /api/v1/library/:id` - Remove from library (soft delete) ✅
   - [x] Write tests for all endpoints ✅

2. **Search Endpoints**:
   - [ ] `GET /api/v1/search` - Search for audiobooks
   - [ ] Implement basic search (will integrate Jackett later)

3. **Download Endpoints**:
   - [ ] `GET /api/v1/downloads` - List downloads
   - [ ] `POST /api/v1/downloads` - Start download (placeholder for now)

**Estimated Time**: 4-6 hours

### 4. qBittorrent Integration

**Why**: Core functionality - need to download audiobooks.

**Tasks:**
- [ ] Create `pkg/qbit/client.go` - qBittorrent API client
- [ ] Implement authentication
- [ ] Implement: add torrent, get status, monitor progress
- [ ] Create download service in `internal/services/download/`
- [ ] Write tests (with mocks)
- [ ] Integrate with download endpoints

**Estimated Time**: 4-6 hours

### 5. Jackett Integration

**Why**: Need to search for audiobooks.

**Tasks:**
- [ ] Create `pkg/jackett/client.go` - Jackett API client
- [ ] Implement search functionality
- [ ] Parse and filter results
- [ ] Create search service in `internal/services/search/`
- [ ] Write tests
- [ ] Integrate with search endpoints

**Estimated Time**: 3-4 hours

## Recommended Development Order

### Week 1: Foundation
1. **Day 1-2**: Database models + API response format
2. **Day 3-4**: Core library API endpoints
3. **Day 5**: Testing and refinement

### Week 2: Integrations
1. **Day 1-2**: qBittorrent integration
2. **Day 3-4**: Jackett integration
3. **Day 5**: Testing and integration

### Week 3: Processing
1. **Day 1-2**: m4b-tool integration
2. **Day 3-4**: File processing pipeline
3. **Day 5**: Testing

## Quick Wins (Can Do Anytime)

- [ ] Add pagination to library endpoint
- [ ] Add filtering/sorting to library endpoint
- [ ] Improve error messages
- [ ] Add request validation
- [ ] Create API documentation (Swagger/OpenAPI)

## Decisions Still Needed

Before implementing, decide on:
1. **File Organization**: `Author/Series/Book Title.m4b` format
2. **Metadata Source**: Start with Open Library
3. **Processing Queue**: Sequential processing for MVP

See `DECISIONS_NEEDED.md` for full list.

## Getting Started Right Now

**Immediate action items:**

1. **Create database models**:
   ```bash
   # Create model files
   touch internal/models/{author,book,audiobook,series,release,download,processing_task,library_item}.go
   ```

2. **Start with Author model** (simplest):
   - Fields: ID, Name, Biography, ImageURL, CreatedAt, UpdatedAt
   - GORM tags for database
   - Basic CRUD operations

3. **Then Book model**:
   - Fields: ID, Title, ISBN, ASIN, Description, CoverArtURL, ReleaseDate
   - Relationship to Author
   - Relationship to Series (optional)

4. **Set up AutoMigrate**:
   - Update `internal/database/database.go` to migrate all models

5. **Write tests** for models

## Testing Checklist

For each new feature:
- [ ] Unit tests written
- [ ] Tests pass (`make test`)
- [ ] Linting passes (`make lint`)
- [ ] Code formatted (`make fmt`)
- [ ] Integration tests (if applicable)
- [ ] API endpoint tests (if applicable)

## Resources

- [GORM Documentation](https://gorm.io/docs/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [qBittorrent API](https://github.com/qbittorrent/qBittorrent/wiki/Web-API-Reference)
- [Jackett API](https://github.com/Jackett/Jackett)

---

**Recommendation**: Start with **Database Models** - it's the foundation everything else builds on.

