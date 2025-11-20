package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Migrate all models
	err = db.AutoMigrate(
		&Author{},
		&Series{},
		&Book{},
		&Audiobook{},
		&Release{},
		&LibraryItem{},
		&Download{},
		&ProcessingTask{},
	)
	assert.NoError(t, err)

	return db
}

func TestAuthor(t *testing.T) {
	db := setupTestDB(t)

	author := Author{
		Name:      "Test Author",
		Biography: "Test biography",
		ImageURL:  "https://example.com/image.jpg",
	}

	err := db.Create(&author).Error
	assert.NoError(t, err)
	assert.NotZero(t, author.ID)
	assert.NotZero(t, author.CreatedAt)

	// Retrieve
	var retrieved Author
	err = db.First(&retrieved, author.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, author.Name, retrieved.Name)
}

func TestBook_WithAuthor(t *testing.T) {
	db := setupTestDB(t)

	// Create author first
	author := Author{Name: "Test Author"}
	db.Create(&author)

	// Create book
	book := Book{
		Title:    "Test Book",
		ISBN:     "1234567890",
		AuthorID: author.ID,
	}

	err := db.Create(&book).Error
	assert.NoError(t, err)
	assert.NotZero(t, book.ID)

	// Retrieve with author
	var retrieved Book
	err = db.Preload("Author").First(&retrieved, book.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, book.Title, retrieved.Title)
	assert.Equal(t, author.Name, retrieved.Author.Name)
}

func TestBook_WithSeries(t *testing.T) {
	db := setupTestDB(t)

	author := Author{Name: "Test Author"}
	db.Create(&author)

	series := Series{Name: "Test Series"}
	db.Create(&series)

	position := 1
	book := Book{
		Title:          "Test Book",
		AuthorID:       author.ID,
		SeriesID:       &series.ID,
		SeriesPosition: &position,
	}

	err := db.Create(&book).Error
	assert.NoError(t, err)

	// Retrieve with series
	var retrieved Book
	err = db.Preload("Series").First(&retrieved, book.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, series.Name, retrieved.Series.Name)
	assert.Equal(t, 1, *retrieved.SeriesPosition)
}

func TestAudiobook_WithBook(t *testing.T) {
	db := setupTestDB(t)

	author := Author{Name: "Test Author"}
	db.Create(&author)

	book := Book{
		Title:    "Test Book",
		AuthorID: author.ID,
	}
	db.Create(&book)

	audiobook := Audiobook{
		BookID:   book.ID,
		Narrator: "Test Narrator",
		Duration: 3600, // 1 hour
		Format:   "m4b",
		Bitrate:  128,
	}

	err := db.Create(&audiobook).Error
	assert.NoError(t, err)
	assert.NotZero(t, audiobook.ID)

	// Retrieve with book
	var retrieved Audiobook
	err = db.Preload("Book").First(&retrieved, audiobook.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, book.Title, retrieved.Book.Title)
}

func TestDownload_Status(t *testing.T) {
	db := setupTestDB(t)

	author := Author{Name: "Test Author"}
	db.Create(&author)

	book := Book{
		Title:    "Test Book",
		AuthorID: author.ID,
	}
	db.Create(&book)

	libraryItem := LibraryItem{
		BookID:    book.ID,
		Status:    LibraryItemStatusWanted,
		AddedDate: time.Now(),
	}
	db.Create(&libraryItem)

	release := Release{
		BookID: book.ID,
		Format: "m4b",
		Size:   1000000,
	}
	db.Create(&release)

	download := Download{
		LibraryItemID: libraryItem.ID,
		ReleaseID:     release.ID,
		Status:        DownloadStatusQueued,
		Progress:      0,
	}

	err := db.Create(&download).Error
	assert.NoError(t, err)

	// Test status methods
	assert.True(t, download.IsActive())
	assert.False(t, download.IsComplete())
	assert.False(t, download.IsFailed())

	// Update to completed
	download.Status = DownloadStatusCompleted
	db.Save(&download)
	assert.False(t, download.IsActive())
	assert.True(t, download.IsComplete())
}

func TestProcessingTask_Status(t *testing.T) {
	db := setupTestDB(t)

	author := Author{Name: "Test Author"}
	db.Create(&author)

	book := Book{
		Title:    "Test Book",
		AuthorID: author.ID,
	}
	db.Create(&book)

	libraryItem := LibraryItem{
		BookID:    book.ID,
		Status:    LibraryItemStatusWanted,
		AddedDate: time.Now(),
	}
	db.Create(&libraryItem)

	release := Release{BookID: book.ID}
	db.Create(&release)

	download := Download{
		LibraryItemID: libraryItem.ID,
		ReleaseID:     release.ID,
		Status:        DownloadStatusCompleted,
	}
	db.Create(&download)

	task := ProcessingTask{
		DownloadID: download.ID,
		Status:     ProcessingStatusPending,
		InputPath:  "/tmp/download",
		Progress:   0,
	}

	err := db.Create(&task).Error
	assert.NoError(t, err)

	// Test status methods
	assert.True(t, task.IsActive())
	assert.False(t, task.IsComplete())

	// Update to completed
	task.Status = ProcessingStatusCompleted
	db.Save(&task)
	assert.False(t, task.IsActive())
	assert.True(t, task.IsComplete())
}

func TestLibraryItem_Status(t *testing.T) {
	db := setupTestDB(t)

	author := Author{Name: "Test Author"}
	db.Create(&author)

	book := Book{
		Title:    "Test Book",
		AuthorID: author.ID,
	}
	db.Create(&book)

	item := LibraryItem{
		BookID:    book.ID,
		Status:    LibraryItemStatusWanted,
		AddedDate: time.Now(),
	}

	err := db.Create(&item).Error
	assert.NoError(t, err)

	// Test status methods
	assert.False(t, item.IsAvailable())
	assert.False(t, item.IsInProgress())

	// Update to available
	item.Status = LibraryItemStatusAvailable
	item.FilePath = "/library/author/book.m4b"
	db.Save(&item)
	assert.True(t, item.IsAvailable())
}

func TestLibraryItem_GetActiveDownload(t *testing.T) {
	db := setupTestDB(t)

	author := Author{Name: "Test Author"}
	db.Create(&author)

	book := Book{
		Title:    "Test Book",
		AuthorID: author.ID,
	}
	db.Create(&book)

	libraryItem := LibraryItem{
		BookID:    book.ID,
		Status:    LibraryItemStatusDownloading,
		AddedDate: time.Now(),
	}
	db.Create(&libraryItem)

	release := Release{BookID: book.ID}
	db.Create(&release)

	download := Download{
		LibraryItemID: libraryItem.ID,
		ReleaseID:     release.ID,
		Status:        DownloadStatusDownloading,
		Progress:      50,
	}
	db.Create(&download)

	// Test GetActiveDownload
	activeDownload, err := libraryItem.GetActiveDownload(db)
	assert.NoError(t, err)
	assert.NotNil(t, activeDownload)
	assert.Equal(t, download.ID, activeDownload.ID)
}

func TestRelease(t *testing.T) {
	db := setupTestDB(t)

	author := Author{Name: "Test Author"}
	db.Create(&author)

	book := Book{
		Title:    "Test Book",
		AuthorID: author.ID,
	}
	db.Create(&book)

	release := Release{
		BookID:    book.ID,
		Quality:   "128kbps",
		Format:    "m4b",
		Size:      50000000, // 50MB
		Indexer:   "test-indexer",
		IndexerID: "12345",
		MagnetURL: "magnet:?xt=urn:btih:test",
		Seeders:   10,
		Leechers:  2,
	}

	err := db.Create(&release).Error
	assert.NoError(t, err)
	assert.NotZero(t, release.ID)

	// Retrieve with book
	var retrieved Release
	err = db.Preload("Book").First(&retrieved, release.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, release.Indexer, retrieved.Indexer)
	assert.Equal(t, book.Title, retrieved.Book.Title)
}
