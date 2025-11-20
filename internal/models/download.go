package models

import (
	"time"

	"gorm.io/gorm"
)

// DownloadStatus represents the status of a download
type DownloadStatus string

const (
	DownloadStatusQueued      DownloadStatus = "queued"
	DownloadStatusDownloading DownloadStatus = "downloading"
	DownloadStatusCompleted   DownloadStatus = "completed"
	DownloadStatusFailed      DownloadStatus = "failed"
	DownloadStatusPaused      DownloadStatus = "paused"
)

// Download represents a download task
type Download struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	LibraryItemID uint        `gorm:"not null;index" json:"library_item_id"`
	LibraryItem   LibraryItem `gorm:"foreignKey:LibraryItemID" json:"library_item,omitempty"`
	ReleaseID     uint        `gorm:"not null;index" json:"release_id"`
	Release       Release     `gorm:"foreignKey:ReleaseID" json:"release,omitempty"`

	// Download information
	Status          DownloadStatus `gorm:"not null;index;default:'queued'" json:"status"`
	Progress        float64        `gorm:"default:0" json:"progress"` // 0-100
	Speed           int64          `json:"speed,omitempty"`           // bytes per second
	Size            int64          `json:"size,omitempty"`            // total size in bytes
	Downloaded      int64          `json:"downloaded,omitempty"`      // bytes downloaded
	Error           string         `gorm:"type:text" json:"error,omitempty"`
	QBittorrentHash string         `gorm:"index" json:"qbittorrent_hash,omitempty"`  // qBittorrent torrent hash
	DownloadPath    string         `gorm:"type:text" json:"download_path,omitempty"` // Path where files are downloaded
	CompletedAt     *time.Time     `json:"completed_at,omitempty"`
}

// TableName specifies the table name for Download
func (Download) TableName() string {
	return "downloads"
}

// IsActive returns true if download is in progress
func (d *Download) IsActive() bool {
	return d.Status == DownloadStatusDownloading || d.Status == DownloadStatusQueued
}

// IsComplete returns true if download is completed
func (d *Download) IsComplete() bool {
	return d.Status == DownloadStatusCompleted
}

// IsFailed returns true if download failed
func (d *Download) IsFailed() bool {
	return d.Status == DownloadStatusFailed
}
