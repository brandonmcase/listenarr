package models

import (
	"time"

	"gorm.io/gorm"
)

// Book represents a book (the written work)
type Book struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Book information
	Title       string     `gorm:"not null;index" json:"title"`
	ISBN        string     `gorm:"index" json:"isbn,omitempty"`
	ASIN        string     `gorm:"index" json:"asin,omitempty"`
	Description string     `gorm:"type:text" json:"description,omitempty"`
	CoverArtURL string     `json:"cover_art_url,omitempty"`
	ReleaseDate *time.Time `json:"release_date,omitempty"`
	Genre       string     `json:"genre,omitempty"`
	Language    string     `json:"language,omitempty"`

	// Relationships
	AuthorID uint   `gorm:"not null;index" json:"author_id"`
	Author   Author `gorm:"foreignKey:AuthorID" json:"author,omitempty"`

	SeriesID       *uint   `gorm:"index" json:"series_id,omitempty"`
	Series         *Series `gorm:"foreignKey:SeriesID" json:"series,omitempty"`
	SeriesPosition *int    `json:"series_position,omitempty"`

	// Related models
	Audiobook    *Audiobook    `gorm:"foreignKey:BookID" json:"audiobook,omitempty"`
	Releases     []Release     `gorm:"foreignKey:BookID" json:"releases,omitempty"`
	LibraryItems []LibraryItem `gorm:"foreignKey:BookID" json:"library_items,omitempty"`
}

// TableName specifies the table name for Book
func (Book) TableName() string {
	return "books"
}

// CompositeIndex creates an index on title and author for faster searches
func (Book) CompositeIndex() string {
	return "idx_books_title_author"
}
