package models

import (
	"time"

	"gorm.io/gorm"
)

// Release represents a specific release/edition of an audiobook
type Release struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationship to Book
	BookID uint `gorm:"not null;index" json:"book_id"`
	Book   Book `gorm:"foreignKey:BookID" json:"book,omitempty"`

	// Release information
	Quality     string     `json:"quality,omitempty"`                 // 64kbps, 128kbps, etc.
	Format      string     `json:"format,omitempty"`                  // mp3, m4b, etc.
	Size        int64      `json:"size,omitempty"`                    // Size in bytes
	Indexer     string     `json:"indexer,omitempty"`                 // Which indexer found this
	IndexerID   string     `gorm:"index" json:"indexer_id,omitempty"` // ID from indexer
	MagnetURL   string     `gorm:"type:text" json:"magnet_url,omitempty"`
	TorrentURL  string     `gorm:"type:text" json:"torrent_url,omitempty"`
	TorrentHash string     `gorm:"index" json:"torrent_hash,omitempty"`
	Seeders     int        `json:"seeders,omitempty"`
	Leechers    int        `json:"leechers,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}

// TableName specifies the table name for Release
func (Release) TableName() string {
	return "releases"
}
