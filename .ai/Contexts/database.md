# Database Models

## Overview

Listenarr uses SQLite with GORM for database management. All models are defined in `internal/models/` and are automatically migrated on startup.

## Models

### Core Models

#### Author
- **Purpose**: Represents book authors
- **Key Fields**: Name, Biography, ImageURL, GoodreadsID
- **Relationships**: Has many Books
- **Indexes**: Name, GoodreadsID

#### Book
- **Purpose**: Represents a written work (the book itself)
- **Key Fields**: Title, ISBN, ASIN, Description, CoverArtURL, ReleaseDate
- **Relationships**: 
  - Belongs to Author
  - Belongs to Series (optional)
  - Has one Audiobook
  - Has many Releases
  - Has many LibraryItems
- **Indexes**: Title, ISBN, ASIN, AuthorID, SeriesID, Composite (Title+AuthorID)

#### Series
- **Purpose**: Represents a book series
- **Key Fields**: Name, Description, TotalBooks
- **Relationships**: Has many Books
- **Indexes**: Name

#### Audiobook
- **Purpose**: Represents the audiobook version of a book
- **Key Fields**: Narrator, Publisher, Duration, Format, Bitrate, Language, ASIN
- **Relationships**: Belongs to Book (one-to-one)
- **Indexes**: BookID (unique), ASIN

### Release & Download Models

#### Release
- **Purpose**: Represents a specific release/edition found on indexers
- **Key Fields**: Quality, Format, Size, Indexer, IndexerID, MagnetURL, TorrentURL, TorrentHash
- **Relationships**: Belongs to Book
- **Indexes**: BookID, IndexerID, TorrentHash

#### Download
- **Purpose**: Tracks download tasks
- **Key Fields**: Status, Progress, Speed, Size, Downloaded, Error, QBittorrentHash
- **Relationships**: 
  - Belongs to LibraryItem
  - Belongs to Release
- **Indexes**: LibraryItemID, ReleaseID, Status, QBittorrentHash
- **Status Values**: queued, downloading, completed, failed, paused

#### ProcessingTask
- **Purpose**: Tracks file processing tasks
- **Key Fields**: Status, Progress, InputPath, OutputPath, Error
- **Relationships**: Belongs to Download
- **Indexes**: DownloadID, Status
- **Status Values**: pending, processing, completed, failed

### Library Model

#### LibraryItem
- **Purpose**: Represents an item in the user's library
- **Key Fields**: Status, FilePath, FileSize, AddedDate, CompletedDate
- **Relationships**: 
  - Belongs to Book
  - Has many Downloads
- **Indexes**: BookID, Status
- **Status Values**: wanted, downloading, processing, available, error

## Relationships Diagram

```
Author
  └─ has many → Book
                  ├─ belongs to → Series (optional)
                  ├─ has one → Audiobook
                  ├─ has many → Release
                  └─ has many → LibraryItem
                                    └─ has many → Download
                                                    └─ has many → ProcessingTask
```

## Database Initialization

Models are automatically migrated on startup via `database.Initialize()`:

```go
db.AutoMigrate(
    &models.Author{},
    &models.Series{},
    &models.Book{},
    &models.Audiobook{},
    &models.Release{},
    &models.LibraryItem{},
    &models.Download{},
    &models.ProcessingTask{},
)
```

## Indexes

### Automatic Indexes (from GORM tags)
- Primary keys on all models
- Foreign key indexes
- Single field indexes (Name, ISBN, ASIN, etc.)
- Soft delete indexes (deleted_at)

### Manual Indexes
- Composite index: `idx_books_title_author` on `books(title, author_id)`

## Usage Examples

### Create Author and Book
```go
author := models.Author{Name: "J.K. Rowling"}
db.Create(&author)

book := models.Book{
    Title:    "Harry Potter and the Philosopher's Stone",
    AuthorID: author.ID,
    ISBN:     "9780747532699",
}
db.Create(&book)
```

### Query with Relationships
```go
var book models.Book
db.Preload("Author").Preload("Series").First(&book, bookID)
```

### Find Library Items
```go
var items []models.LibraryItem
db.Preload("Book.Author").Where("status = ?", "available").Find(&items)
```

## Testing

All models have comprehensive tests in `internal/models/models_test.go`:
- Model creation and retrieval
- Relationship loading
- Status method testing
- Complex queries

## Migration Strategy

- **MVP**: Use GORM AutoMigrate (current approach)
- **Future**: Consider migration files for production deployments
- **Schema Changes**: AutoMigrate handles additions, but deletions require manual migration

