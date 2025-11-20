package api

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/listenarr/listenarr/internal/models"
)

// CreateBookRequest represents the request body for creating a book
type CreateBookRequest struct {
	Title          string     `json:"title" binding:"required"`
	AuthorID       uint       `json:"author_id" binding:"required"`
	ISBN           *string    `json:"isbn,omitempty"`
	ASIN           *string    `json:"asin,omitempty"`
	Description    *string    `json:"description,omitempty"`
	CoverArtURL    *string    `json:"cover_art_url,omitempty"`
	ReleaseDate    *time.Time `json:"release_date,omitempty"`
	Genre          *string    `json:"genre,omitempty"`
	Language       *string    `json:"language,omitempty"`
	SeriesID       *uint      `json:"series_id,omitempty"`
	SeriesPosition *int       `json:"series_position,omitempty"`
}

// UpdateBookRequest represents the request body for updating a book
type UpdateBookRequest struct {
	Title          *string    `json:"title,omitempty"`
	AuthorID       *uint      `json:"author_id,omitempty"`
	ISBN           *string    `json:"isbn,omitempty"`
	ASIN           *string    `json:"asin,omitempty"`
	Description    *string    `json:"description,omitempty"`
	CoverArtURL    *string    `json:"cover_art_url,omitempty"`
	ReleaseDate    *time.Time `json:"release_date,omitempty"`
	Genre          *string    `json:"genre,omitempty"`
	Language       *string    `json:"language,omitempty"`
	SeriesID       *uint      `json:"series_id,omitempty"`
	SeriesPosition *int       `json:"series_position,omitempty"`
}

// BookResponseDetailed represents a book in API responses with full details
type BookResponseDetailed struct {
	ID             uint            `json:"id"`
	Title          string          `json:"title"`
	ISBN           string          `json:"isbn,omitempty"`
	ASIN           string          `json:"asin,omitempty"`
	Description    string          `json:"description,omitempty"`
	CoverArtURL    string          `json:"cover_art_url,omitempty"`
	ReleaseDate    *time.Time      `json:"release_date,omitempty"`
	Genre          string          `json:"genre,omitempty"`
	Language       string          `json:"language,omitempty"`
	AuthorID       uint            `json:"author_id"`
	Author         *AuthorResponse `json:"author,omitempty"`
	SeriesID       *uint           `json:"series_id,omitempty"`
	Series         *SeriesResponse `json:"series,omitempty"`
	SeriesPosition *int            `json:"series_position,omitempty"`
	Audiobook      interface{}     `json:"audiobook,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// toBookResponseDetailed converts a Book model to detailed API response format
func toBookResponseDetailed(book *models.Book) *BookResponseDetailed {
	response := &BookResponseDetailed{
		ID:             book.ID,
		Title:          book.Title,
		ISBN:           book.ISBN,
		ASIN:           book.ASIN,
		Description:    book.Description,
		CoverArtURL:    book.CoverArtURL,
		ReleaseDate:    book.ReleaseDate,
		Genre:          book.Genre,
		Language:       book.Language,
		AuthorID:       book.AuthorID,
		SeriesID:       book.SeriesID,
		SeriesPosition: book.SeriesPosition,
		CreatedAt:      book.CreatedAt,
		UpdatedAt:      book.UpdatedAt,
	}

	if book.Author.ID != 0 {
		response.Author = &AuthorResponse{
			ID:          book.Author.ID,
			Name:        book.Author.Name,
			Biography:   book.Author.Biography,
			ImageURL:    book.Author.ImageURL,
			GoodreadsID: book.Author.GoodreadsID,
		}
	}

	if book.Series != nil && book.Series.ID != 0 {
		response.Series = &SeriesResponse{
			ID:          book.Series.ID,
			Name:        book.Series.Name,
			Description: book.Series.Description,
			TotalBooks:  book.Series.TotalBooks,
		}
	}

	if book.Audiobook != nil {
		response.Audiobook = map[string]interface{}{
			"id":       book.Audiobook.ID,
			"narrator": book.Audiobook.Narrator,
			"duration": book.Audiobook.Duration,
			"format":   book.Audiobook.Format,
			"bitrate":  book.Audiobook.Bitrate,
			"language": book.Audiobook.Language,
			"asin":     book.Audiobook.ASIN,
		}
	}

	return response
}

// getBooks handles GET /api/v1/books
func (s *Server) getBooks(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Validate pagination
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

	// Build query
	query := s.db.Model(&models.Book{})

	// Apply filters
	if search := c.Query("search"); search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}
	if authorIDStr := c.Query("author_id"); authorIDStr != "" {
		if authorID, err := strconv.ParseUint(authorIDStr, 10, 32); err == nil {
			query = query.Where("author_id = ?", uint(authorID))
		}
	}
	if seriesIDStr := c.Query("series_id"); seriesIDStr != "" {
		if seriesID, err := strconv.ParseUint(seriesIDStr, 10, 32); err == nil {
			query = query.Where("series_id = ?", uint(seriesID))
		}
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Apply sorting
	sortBy := c.DefaultQuery("sort", "title")
	order := c.DefaultQuery("order", "asc")
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	switch sortBy {
	case "title":
		query = query.Order("title " + order)
	case "created_at":
		query = query.Order("created_at " + order)
	case "release_date":
		query = query.Order("release_date " + order)
	default:
		query = query.Order("title " + order)
	}

	// Apply pagination and preload relationships
	var books []models.Book
	err := query.
		Preload("Author").
		Preload("Series").
		Offset(offset).
		Limit(limit).
		Find(&books).Error

	if err != nil {
		InternalErrorResponse(c, "Failed to fetch books")
		return
	}

	// Convert to response format
	responseData := make([]*BookResponseDetailed, len(books))
	for i := range books {
		responseData[i] = toBookResponseDetailed(&books[i])
	}

	PaginatedSuccessResponse(c, responseData, page, limit, int(total))
}

// getBook handles GET /api/v1/books/:id
func (s *Server) getBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid book ID")
		return
	}

	var book models.Book
	err = s.db.
		Preload("Author").
		Preload("Series").
		Preload("Audiobook").
		Preload("Releases").
		Preload("LibraryItems").
		First(&book, uint(id)).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFoundResponse(c, "book")
			return
		}
		InternalErrorResponse(c, "Failed to fetch book")
		return
	}

	SuccessResponse(c, StatusOK, toBookResponseDetailed(&book))
}

// createBook handles POST /api/v1/books
func (s *Server) createBook(c *gin.Context) {
	var req CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, err)
		return
	}

	// Verify author exists
	var author models.Author
	err := s.db.First(&author, req.AuthorID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFoundResponse(c, "author")
			return
		}
		InternalErrorResponse(c, "Failed to find author")
		return
	}

	// Verify series exists if provided
	if req.SeriesID != nil {
		var series models.Series
		err := s.db.First(&series, *req.SeriesID).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				NotFoundResponse(c, "series")
				return
			}
			InternalErrorResponse(c, "Failed to find series")
			return
		}
	}

	// Check for duplicate book (by title + author or ISBN/ASIN)
	var existingBook models.Book
	bookQuery := s.db.Where("title = ? AND author_id = ?", req.Title, req.AuthorID)
	if req.ISBN != nil && *req.ISBN != "" {
		bookQuery = bookQuery.Or("isbn = ?", *req.ISBN)
	}
	if req.ASIN != nil && *req.ASIN != "" {
		bookQuery = bookQuery.Or("asin = ?", *req.ASIN)
	}

	err = bookQuery.First(&existingBook).Error
	if err == nil {
		ConflictResponse(c, "Book already exists")
		return
	} else if err != gorm.ErrRecordNotFound {
		InternalErrorResponse(c, "Failed to check existing book")
		return
	}

	// Create book
	book := models.Book{
		Title:          req.Title,
		AuthorID:       req.AuthorID,
		SeriesID:       req.SeriesID,
		SeriesPosition: req.SeriesPosition,
	}
	if req.ISBN != nil {
		book.ISBN = *req.ISBN
	}
	if req.ASIN != nil {
		book.ASIN = *req.ASIN
	}
	if req.Description != nil {
		book.Description = *req.Description
	}
	if req.CoverArtURL != nil {
		book.CoverArtURL = *req.CoverArtURL
	}
	if req.ReleaseDate != nil {
		book.ReleaseDate = req.ReleaseDate
	}
	if req.Genre != nil {
		book.Genre = *req.Genre
	}
	if req.Language != nil {
		book.Language = *req.Language
	}

	if err := s.db.Create(&book).Error; err != nil {
		InternalErrorResponse(c, "Failed to create book")
		return
	}

	// Reload with relationships
	err = s.db.
		Preload("Author").
		Preload("Series").
		First(&book, book.ID).Error
	if err != nil {
		InternalErrorResponse(c, "Failed to reload book")
		return
	}

	CreatedResponse(c, toBookResponseDetailed(&book))
}

// updateBook handles PUT /api/v1/books/:id
func (s *Server) updateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid book ID")
		return
	}

	var req UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, err)
		return
	}

	// Check if book exists
	var book models.Book
	err = s.db.First(&book, uint(id)).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFoundResponse(c, "book")
			return
		}
		InternalErrorResponse(c, "Failed to find book")
		return
	}

	// Update fields if provided
	if req.Title != nil {
		book.Title = *req.Title
	}
	if req.AuthorID != nil {
		// Verify author exists
		var author models.Author
		err := s.db.First(&author, *req.AuthorID).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				NotFoundResponse(c, "author")
				return
			}
			InternalErrorResponse(c, "Failed to find author")
			return
		}
		book.AuthorID = *req.AuthorID
	}
	if req.ISBN != nil {
		book.ISBN = *req.ISBN
	}
	if req.ASIN != nil {
		book.ASIN = *req.ASIN
	}
	if req.Description != nil {
		book.Description = *req.Description
	}
	if req.CoverArtURL != nil {
		book.CoverArtURL = *req.CoverArtURL
	}
	if req.ReleaseDate != nil {
		book.ReleaseDate = req.ReleaseDate
	}
	if req.Genre != nil {
		book.Genre = *req.Genre
	}
	if req.Language != nil {
		book.Language = *req.Language
	}
	if req.SeriesID != nil {
		if *req.SeriesID == 0 {
			book.SeriesID = nil
		} else {
			// Verify series exists
			var series models.Series
			err := s.db.First(&series, *req.SeriesID).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					NotFoundResponse(c, "series")
					return
				}
				InternalErrorResponse(c, "Failed to find series")
				return
			}
			book.SeriesID = req.SeriesID
		}
	}
	if req.SeriesPosition != nil {
		book.SeriesPosition = req.SeriesPosition
	}

	if err := s.db.Save(&book).Error; err != nil {
		InternalErrorResponse(c, "Failed to update book")
		return
	}

	// Reload with relationships
	err = s.db.
		Preload("Author").
		Preload("Series").
		Preload("Audiobook").
		First(&book, book.ID).Error
	if err != nil {
		InternalErrorResponse(c, "Failed to reload book")
		return
	}

	SuccessResponse(c, StatusOK, toBookResponseDetailed(&book))
}

// deleteBook handles DELETE /api/v1/books/:id
func (s *Server) deleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid book ID")
		return
	}

	// Check if book exists
	var book models.Book
	err = s.db.First(&book, uint(id)).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFoundResponse(c, "book")
			return
		}
		InternalErrorResponse(c, "Failed to find book")
		return
	}

	// Check if book has library items
	var libraryItemCount int64
	s.db.Model(&models.LibraryItem{}).Where("book_id = ?", id).Count(&libraryItemCount)
	if libraryItemCount > 0 {
		ConflictResponse(c, "Cannot delete book with existing library items")
		return
	}

	// Soft delete (GORM handles this automatically with DeletedAt)
	err = s.db.Delete(&book).Error
	if err != nil {
		InternalErrorResponse(c, "Failed to delete book")
		return
	}

	NoContentResponse(c)
}
