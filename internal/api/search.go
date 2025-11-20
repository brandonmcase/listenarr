package api

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/listenarr/listenarr/internal/models"
)

// SearchResponse represents a search result
type SearchResponse struct {
	Query   string             `json:"query"`
	Results []SearchResultItem `json:"results"`
	Total   int                `json:"total"`
}

// SearchResultItem represents a single search result item
type SearchResultItem struct {
	Type        string  `json:"type"` // "book", "author", "series"
	ID          uint    `json:"id"`
	Title       string  `json:"title"`
	Author      string  `json:"author,omitempty"`
	Description string  `json:"description,omitempty"`
	CoverArtURL string  `json:"cover_art_url,omitempty"`
	MatchScore  float64 `json:"match_score,omitempty"`
}

// searchAudiobooks handles GET /api/v1/search
func (s *Server) searchAudiobooks(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		BadRequestResponse(c, "Search query parameter 'q' is required")
		return
	}

	// Parse pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	results := make([]SearchResultItem, 0)

	// Search books
	var books []models.Book
	bookQuery := s.db.Model(&models.Book{}).
		Where("title LIKE ?", "%"+query+"%").
		Or("isbn = ?", query).
		Or("asin = ?", query).
		Preload("Author").
		Limit(limit).
		Offset(offset)

	bookQuery.Find(&books)
	for _, book := range books {
		authorName := ""
		if book.Author.ID != 0 {
			authorName = book.Author.Name
		}
		results = append(results, SearchResultItem{
			Type:        "book",
			ID:          book.ID,
			Title:       book.Title,
			Author:      authorName,
			Description: book.Description,
			CoverArtURL: book.CoverArtURL,
		})
	}

	// Search authors
	var authors []models.Author
	authorQuery := s.db.Model(&models.Author{}).
		Where("name LIKE ?", "%"+query+"%").
		Limit(limit).
		Offset(offset)

	authorQuery.Find(&authors)
	for _, author := range authors {
		results = append(results, SearchResultItem{
			Type:        "author",
			ID:          author.ID,
			Title:       author.Name,
			Description: author.Biography,
		})
	}

	// For now, return basic search results
	// TODO: Integrate with Jackett for actual audiobook search
	searchResponse := SearchResponse{
		Query:   query,
		Results: results,
		Total:   len(results),
	}

	SuccessResponse(c, StatusOK, searchResponse)
}
