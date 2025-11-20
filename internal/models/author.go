package models

import (
	"time"

	"gorm.io/gorm"
)

// Author represents an author of books
type Author struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Author information
	Name        string `gorm:"not null;index" json:"name"`
	Biography   string `gorm:"type:text" json:"biography,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
	GoodreadsID string `gorm:"index" json:"goodreads_id,omitempty"`

	// Relationships
	Books []Book `gorm:"foreignKey:AuthorID" json:"books,omitempty"`
}

// TableName specifies the table name for Author
func (Author) TableName() string {
	return "authors"
}
