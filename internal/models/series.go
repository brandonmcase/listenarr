package models

import (
	"time"

	"gorm.io/gorm"
)

// Series represents a book series
type Series struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Series information
	Name        string `gorm:"not null;index" json:"name"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	TotalBooks  int    `json:"total_books,omitempty"`

	// Relationships
	Books []Book `gorm:"foreignKey:SeriesID" json:"books,omitempty"`
}

// TableName specifies the table name for Series
func (Series) TableName() string {
	return "series"
}
