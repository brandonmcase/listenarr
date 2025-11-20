package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/listenarr/listenarr/internal/models"
)

// Initialize creates and returns a database connection
func Initialize(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate all models
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// Create additional indexes
	if err := CreateIndexes(db); err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	return db, nil
}

// migrate runs database migrations for all models
func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Author{},
		&models.Series{},
		&models.Book{},
		&models.Audiobook{},
		&models.Release{},
		&models.LibraryItem{},
		&models.Download{},
		&models.ProcessingTask{},
	)
}

// CreateIndexes creates additional indexes for performance
func CreateIndexes(db *gorm.DB) error {
	// Composite index for book searches (title + author)
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_books_title_author ON books(title, author_id)").Error; err != nil {
		return fmt.Errorf("failed to create composite index: %w", err)
	}

	// Index for ISBN/ASIN lookups (if not already created by GORM)
	// GORM should handle these from the model tags, but we can add more if needed

	return nil
}
