package models

import (
	"time"

	"gorm.io/gorm"
)

// LibraryItemStatus represents the status of a library item
type LibraryItemStatus string

const (
	LibraryItemStatusWanted      LibraryItemStatus = "wanted"
	LibraryItemStatusDownloading LibraryItemStatus = "downloading"
	LibraryItemStatusProcessing  LibraryItemStatus = "processing"
	LibraryItemStatusAvailable   LibraryItemStatus = "available"
	LibraryItemStatusError       LibraryItemStatus = "error"
)

// LibraryItem represents an item in the user's library
type LibraryItem struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationship to Book
	BookID uint `gorm:"not null;index" json:"book_id"`
	Book   Book `gorm:"foreignKey:BookID" json:"book,omitempty"`

	// Library item information
	Status        LibraryItemStatus `gorm:"not null;index;default:'wanted'" json:"status"`
	FilePath      string            `gorm:"type:text" json:"file_path,omitempty"` // Path to final m4b file
	FileSize      int64             `json:"file_size,omitempty"`                  // Size in bytes
	AddedDate     time.Time         `gorm:"not null" json:"added_date"`
	CompletedDate *time.Time        `json:"completed_date,omitempty"`

	// Relationships
	Downloads       []Download       `gorm:"foreignKey:LibraryItemID" json:"downloads,omitempty"`
	ProcessingTasks []ProcessingTask `gorm:"foreignKey:DownloadID" json:"processing_tasks,omitempty"` // Through Download
}

// TableName specifies the table name for LibraryItem
func (LibraryItem) TableName() string {
	return "library_items"
}

// IsAvailable returns true if item is available in library
func (l *LibraryItem) IsAvailable() bool {
	return l.Status == LibraryItemStatusAvailable
}

// IsInProgress returns true if item is being downloaded or processed
func (l *LibraryItem) IsInProgress() bool {
	return l.Status == LibraryItemStatusDownloading || l.Status == LibraryItemStatusProcessing
}

// GetActiveDownload returns the active download if any
func (l *LibraryItem) GetActiveDownload(db *gorm.DB) (*Download, error) {
	var download Download
	err := db.Where("library_item_id = ? AND status IN ?", l.ID, []DownloadStatus{
		DownloadStatusQueued,
		DownloadStatusDownloading,
	}).First(&download).Error
	if err != nil {
		return nil, err
	}
	return &download, nil
}
