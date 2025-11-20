package models

import (
	"time"

	"gorm.io/gorm"
)

// Audiobook represents the audiobook version of a book
type Audiobook struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationship to Book
	BookID uint `gorm:"not null;uniqueIndex;index" json:"book_id"`
	Book   Book `gorm:"foreignKey:BookID" json:"book,omitempty"`

	// Audiobook-specific information
	Narrator  string `json:"narrator,omitempty"`
	Publisher string `json:"publisher,omitempty"`
	Duration  int    `json:"duration,omitempty"` // Duration in seconds
	Format    string `json:"format,omitempty"`   // mp3, m4b, m4a, etc.
	Bitrate   int    `json:"bitrate,omitempty"`  // kbps
	Language  string `json:"language,omitempty"`
	ASIN      string `gorm:"index" json:"asin,omitempty"` // Audible ASIN
}

// TableName specifies the table name for Audiobook
func (Audiobook) TableName() string {
	return "audiobooks"
}
