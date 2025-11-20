package search

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/listenarr/listenarr/internal/models"
	"github.com/listenarr/listenarr/pkg/jackett"
)

// Service handles search operations
type Service struct {
	db      *gorm.DB
	jackett *jackett.Client
}

// NewService creates a new search service
func NewService(db *gorm.DB, jackettClient *jackett.Client) *Service {
	return &Service{
		db:      db,
		jackett: jackettClient,
	}
}

// SearchResult represents a unified search result
type SearchResult struct {
	Type        string  `json:"type"` // "book", "release"
	ID          uint    `json:"id,omitempty"`
	Title       string  `json:"title"`
	Author      string  `json:"author,omitempty"`
	Description string  `json:"description,omitempty"`
	Size        int64   `json:"size,omitempty"`
	Seeders     int     `json:"seeders,omitempty"`
	Peers       int     `json:"peers,omitempty"`
	MagnetURI   string  `json:"magnet_uri,omitempty"`
	Tracker     string  `json:"tracker,omitempty"`
	MatchScore  float64 `json:"match_score,omitempty"`
}

// SearchAudiobooks searches for audiobooks using Jackett
func (s *Service) SearchAudiobooks(query string) ([]SearchResult, error) {
	results := make([]SearchResult, 0)

	// If Jackett is configured, search using it
	if s.jackett != nil {
		jackettReq := jackett.SearchRequest{
			Query:    query,
			Category: []int{3030}, // Books category
		}

		jackettResp, err := s.jackett.Search(jackettReq)
		if err != nil {
			// Log error but continue with local search
		} else {
			// Convert Jackett results to unified format
			for _, result := range jackettResp.Results {
				results = append(results, SearchResult{
					Type:        "release",
					Title:       result.Title,
					Description: result.Description,
					Size:        result.Size,
					Seeders:     result.Seeders,
					Peers:       result.Peers,
					MagnetURI:   result.MagnetURI,
					Tracker:     result.Tracker,
				})
			}
		}
	}

	// Also search local database for books
	var books []models.Book
	s.db.Where("title LIKE ?", "%"+query+"%").
		Or("isbn = ?", query).
		Or("asin = ?", query).
		Preload("Author").
		Limit(20).
		Find(&books)

	for _, book := range books {
		authorName := ""
		if book.Author.ID != 0 {
			authorName = book.Author.Name
		}
		results = append(results, SearchResult{
			Type:        "book",
			ID:          book.ID,
			Title:       book.Title,
			Author:      authorName,
			Description: book.Description,
		})
	}

	return results, nil
}

// SearchReleases searches for releases matching a book
func (s *Service) SearchReleases(bookID uint) ([]SearchResult, error) {
	var book models.Book
	if err := s.db.First(&book, bookID).Error; err != nil {
		return nil, fmt.Errorf("book not found: %w", err)
	}

	// Build search query from book title and author
	searchQuery := book.Title
	if book.Author.ID != 0 {
		var author models.Author
		s.db.First(&author, book.Author.ID)
		searchQuery = fmt.Sprintf("%s %s", author.Name, book.Title)
	}

	// Search using Jackett
	if s.jackett == nil {
		return []SearchResult{}, nil
	}

	jackettReq := jackett.SearchRequest{
		Query:    searchQuery,
		Category: []int{3030}, // Books category
	}

	jackettResp, err := s.jackett.Search(jackettReq)
	if err != nil {
		return nil, fmt.Errorf("jackett search failed: %w", err)
	}

	results := make([]SearchResult, len(jackettResp.Results))
	for i, result := range jackettResp.Results {
		results[i] = SearchResult{
			Type:        "release",
			Title:       result.Title,
			Description: result.Description,
			Size:        result.Size,
			Seeders:     result.Seeders,
			Peers:       result.Peers,
			MagnetURI:   result.MagnetURI,
			Tracker:     result.Tracker,
		}
	}

	return results, nil
}
