package models

import (
	"time"

	"gorm.io/gorm"
)

// ProcessingStatus represents the status of a processing task
type ProcessingStatus string

const (
	ProcessingStatusPending    ProcessingStatus = "pending"
	ProcessingStatusProcessing ProcessingStatus = "processing"
	ProcessingStatusCompleted  ProcessingStatus = "completed"
	ProcessingStatusFailed     ProcessingStatus = "failed"
)

// ProcessingTask represents a file processing task
type ProcessingTask struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationship to Download
	DownloadID uint     `gorm:"not null;index" json:"download_id"`
	Download   Download `gorm:"foreignKey:DownloadID" json:"download,omitempty"`

	// Processing information
	Status      ProcessingStatus `gorm:"not null;index;default:'pending'" json:"status"`
	Progress    float64          `gorm:"default:0" json:"progress"`              // 0-100
	InputPath   string           `gorm:"type:text;not null" json:"input_path"`   // Path to downloaded files
	OutputPath  string           `gorm:"type:text" json:"output_path,omitempty"` // Path to processed m4b file
	Error       string           `gorm:"type:text" json:"error,omitempty"`
	StartedAt   *time.Time       `json:"started_at,omitempty"`
	CompletedAt *time.Time       `json:"completed_at,omitempty"`
}

// TableName specifies the table name for ProcessingTask
func (ProcessingTask) TableName() string {
	return "processing_tasks"
}

// IsActive returns true if processing is in progress
func (p *ProcessingTask) IsActive() bool {
	return p.Status == ProcessingStatusProcessing || p.Status == ProcessingStatusPending
}

// IsComplete returns true if processing is completed
func (p *ProcessingTask) IsComplete() bool {
	return p.Status == ProcessingStatusCompleted
}

// IsFailed returns true if processing failed
func (p *ProcessingTask) IsFailed() bool {
	return p.Status == ProcessingStatusFailed
}
