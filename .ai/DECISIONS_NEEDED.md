# Decisions Needed Before Starting Development

This document summarizes the key decisions that need to be made before beginning active development.

## Critical Decisions (Must Decide Before Phase 1)

### 1. Authentication & Security ✅

- [x] **Authentication Method**: **API key authentication** ✅
  - **Decision**: Start with API key for MVP, add OAuth later if needed
  - **Implementation**: API key in header (`X-API-Key`) or query parameter
  - **Storage**: API key stored in config file, generated on first run
  - **Future**: OAuth support can be added later if needed

### 2. Database Schema ✅

- [x] **Model Relationships**: Finalize relationships between Author, Book, Audiobook, Series ✅
  - **Decision**: Author → Books (one-to-many), Book → Audiobook (one-to-one), Series → Books (one-to-many), Book → Releases (one-to-many), LibraryItem → Book (many-to-one)
- [x] **Migration Strategy**: How to handle schema changes? ✅
  - **Decision**: Use GORM AutoMigrate for MVP, all models auto-migrate on startup
- [x] **Indexes**: Which fields need indexes? ✅
  - **Decision**: Implemented indexes on ISBN, ASIN, Title, AuthorID, SeriesID, composite index on (Title, AuthorID), and foreign key indexes

### 3. API Design ✅

- [x] **Response Format**: Standardize API response structure ✅
  - **Decision**: `{ success: bool, data?: T, error?: string, code?: string, details?: object }`
  - **Implementation**: Response helpers in `internal/api/response.go`
  - **Status**: Complete
- [x] **Error Handling**: How to structure error responses? ✅
  - **Decision**: Custom error types (APIError, ValidationErrors) with error codes and details
  - **Implementation**: Error helpers in `internal/api/errors.go`
  - **Status**: Complete
- [x] **HTTP Status Codes**: Standardized status code constants ✅
  - **Decision**: Use standard HTTP status codes (200, 201, 204, 400, 401, 404, 409, 422, 500)
  - **Implementation**: Constants in `internal/api/response.go`
  - **Status**: Complete
- [ ] **Pagination**: For library/search endpoints? (Recommendation: Yes, start with limit/offset)

## Important Decisions (Should Decide Early)

### 4. Metadata Strategy

- [ ] **Primary Metadata Source**: Which to implement first?
  - Recommendation: Open Library (free, no API key needed)
- [ ] **Metadata Embedding Tool**:
  - Option A: m4b-tool (if it supports metadata)
  - Option B: FFmpeg (more control)
  - Option C: AtomicParsley (specialized for MP4/M4B)
  - **Recommendation**: Start with m4b-tool, add FFmpeg if needed

### 5. File Organization

- [ ] **Library Folder Structure**:

  - Option A: `Author/Book Title.m4b`
  - Option B: `Author/Series/Book Title.m4b`
  - Option C: `Author/Series/Series # - Book Title.m4b`
  - **Recommendation**: Option B with configurable naming

- [ ] **File Naming Convention**:
  - Example: `Author - Series - Book Title.m4b`
  - Need to decide: Include narrator? Include year? Format?

### 6. Processing Pipeline

- [ ] **Processing Queue**:
  - Sequential (one at a time) or Concurrent (how many?)
  - **Recommendation**: Sequential for MVP (processing is CPU/disk intensive)
- [ ] **Error Handling**:

  - Retry failed processing? How many times?
  - **Recommendation**: 1 retry, then manual intervention

- [ ] **Temporary File Cleanup**:
  - When to clean up? Immediately or after X days?
  - **Recommendation**: Clean up immediately after successful processing

## Nice-to-Have Decisions (Can Decide Later)

### 7. Quality Profiles

- [ ] **Format Preferences**: MP3, M4B, M4A? (Recommendation: M4B preferred, accept others)
- [ ] **Bitrate Preferences**: 64kbps, 128kbps, 192kbps, 320kbps? (Recommendation: 128kbps minimum)
- [ ] **Quality Profiles**: Create named profiles? (Recommendation: Single profile for MVP)

### 8. Auto-Download

- [ ] **Default Behavior**: Auto-download when added to library?
  - **Recommendation**: Opt-in for MVP, add toggle in settings

### 9. Series Management

- [ ] **Series Detection**: How to detect series from metadata?
- [ ] **Series Ordering**: How to handle series position?
  - **Recommendation**: Use series position from metadata, allow manual override

### 10. Duplicate Detection

- [ ] **Detection Method**:
  - Primary: ISBN/ASIN match
  - Secondary: Title + Author fuzzy match
  - **Recommendation**: Implement both, warn user on potential duplicates

## Integration-Specific Decisions

### 11. qBittorrent

- [ ] **Category/Tags**: Use qBittorrent categories or tags?

  - **Recommendation**: Use category "Listenarr" for organization

- [ ] **Download Path**: Where should qBittorrent download?
  - **Recommendation**: `/downloads` volume, move to library after processing

### 12. Jackett

- [ ] **Indexer Selection**: Let user choose which indexers to use?

  - **Recommendation**: Use all configured indexers, let user filter results

- [ ] **Search Parameters**: What to search for?
  - **Recommendation**: Title + Author, with ISBN/ASIN if available

### 13. Plex

- [ ] **Library Type**: Music library or separate audiobook library?

  - **Recommendation**: Separate audiobook library type (Plex supports this)

- [ ] **Refresh Strategy**:
  - Immediate refresh after each book?
  - Batch refresh on schedule?
  - **Recommendation**: Immediate for MVP, add batch option later

## UI/UX Decisions

### 14. Dashboard Content

- [ ] **What to Display**:
  - Total books, active downloads, processing queue, recent additions
  - **Recommendation**: All of the above

### 15. Library View

- [ ] **Display Format**: List, grid, or both?

  - **Recommendation**: Start with list, add grid view later

- [ ] **Filtering/Sorting**:
  - By author, series, status, date added?
  - **Recommendation**: All of the above

### 16. Search Interface

- [ ] **Search Input**: Single field or separate title/author?

  - **Recommendation**: Single field with smart parsing

- [ ] **Result Display**: How to show search results?
  - **Recommendation**: List with release details, quality, size

## Development Workflow Decisions

### 17. Testing Strategy

- [ ] **Unit Tests**: What level of coverage?

  - **Recommendation**: Critical paths (processing, metadata matching)

- [ ] **Integration Tests**: Test with real qBittorrent/Jackett?
  - **Recommendation**: Mock for CI, real for manual testing

### 18. Documentation

- [ ] **API Documentation**: OpenAPI/Swagger?

  - **Recommendation**: Yes, use Swagger/OpenAPI for Go

- [ ] **User Documentation**: Where to host?
  - **Recommendation**: GitHub Wiki or docs folder

## Priority Order

1. **Start Development With** (Critical):

   - Authentication method
   - Database schema
   - API response format
   - Basic file organization structure

2. **Decide During Phase 1** (Important):

   - Metadata source priority
   - Processing queue strategy
   - Error handling approach

3. **Decide During Phase 2+** (Can Wait):
   - Quality profiles
   - Auto-download behavior
   - Advanced UI features

## Recommendations Summary

Based on MVP approach, here are quick recommendations:

- **Auth**: API key (simple, secure enough)
- **Metadata**: Open Library first, then Google Books
- **File Org**: `Author/Series/Book Title.m4b`
- **Processing**: Sequential, 1 retry on failure
- **Auto-download**: Opt-in toggle
- **Quality**: M4B preferred, 128kbps minimum
- **Plex**: Immediate refresh after processing

These can be changed later, but having defaults helps move forward.
