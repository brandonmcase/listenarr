package api

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/listenarr/listenarr/internal/models"
)

// AddToLibraryRequest represents the request body for adding a book to library
type AddToLibraryRequest struct {
	Title          string  `json:"title" binding:"required"`
	AuthorName     string  `json:"author_name" binding:"required"`
	ISBN           *string `json:"isbn,omitempty"`
	ASIN           *string `json:"asin,omitempty"`
	SeriesName     *string `json:"series_name,omitempty"`
	SeriesPosition *int    `json:"series_position,omitempty"`
}

// LibraryItemResponse represents a library item in API responses
type LibraryItemResponse struct {
	ID            uint          `json:"id"`
	BookID        uint          `json:"book_id"`
	Status        string        `json:"status"`
	FilePath      string        `json:"file_path,omitempty"`
	FileSize      int64         `json:"file_size,omitempty"`
	AddedDate     time.Time     `json:"added_date"`
	CompletedDate *time.Time    `json:"completed_date,omitempty"`
	Book          *BookResponse `json:"book,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// BookResponse represents a book in API responses
type BookResponse struct {
	ID             uint            `json:"id"`
	Title          string          `json:"title"`
	ISBN           string          `json:"isbn,omitempty"`
	ASIN           string          `json:"asin,omitempty"`
	Description    string          `json:"description,omitempty"`
	CoverArtURL    string          `json:"cover_art_url,omitempty"`
	ReleaseDate    *time.Time      `json:"release_date,omitempty"`
	Genre          string          `json:"genre,omitempty"`
	Language       string          `json:"language,omitempty"`
	Author         *AuthorResponse `json:"author,omitempty"`
	Series         *SeriesResponse `json:"series,omitempty"`
	SeriesPosition *int            `json:"series_position,omitempty"`
}

// AuthorResponse represents an author in API responses
type AuthorResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Biography   string `json:"biography,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
	GoodreadsID string `json:"goodreads_id,omitempty"`
}

// SeriesResponse represents a series in API responses
type SeriesResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	TotalBooks  int    `json:"total_books,omitempty"`
}

// toLibraryItemResponse converts a LibraryItem model to API response format
func toLibraryItemResponse(item *models.LibraryItem) *LibraryItemResponse {
	response := &LibraryItemResponse{
		ID:            item.ID,
		BookID:        item.BookID,
		Status:        string(item.Status),
		FilePath:      item.FilePath,
		FileSize:      item.FileSize,
		AddedDate:     item.AddedDate,
		CompletedDate: item.CompletedDate,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
	}

	if item.Book.ID != 0 {
		response.Book = toBookResponse(&item.Book)
	}

	return response
}

// toBookResponse converts a Book model to API response format
func toBookResponse(book *models.Book) *BookResponse {
	response := &BookResponse{
		ID:             book.ID,
		Title:          book.Title,
		ISBN:           book.ISBN,
		ASIN:           book.ASIN,
		Description:    book.Description,
		CoverArtURL:    book.CoverArtURL,
		ReleaseDate:    book.ReleaseDate,
		Genre:          book.Genre,
		Language:       book.Language,
		SeriesPosition: book.SeriesPosition,
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

	return response
}

// getLibrary handles GET /api/v1/library
func (s *Server) getLibrary(c *gin.Context) {
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
	query := s.db.Model(&models.LibraryItem{})

	// Apply filters
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if authorIDStr := c.Query("author_id"); authorIDStr != "" {
		if authorID, err := strconv.ParseUint(authorIDStr, 10, 32); err == nil {
			query = query.Joins("JOIN books ON library_items.book_id = books.id").
				Where("books.author_id = ?", uint(authorID))
		}
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Apply sorting
	sortBy := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	// Handle special sorting cases
	switch sortBy {
	case "title":
		query = query.Joins("JOIN books ON library_items.book_id = books.id").
			Order("books.title " + order)
	case "added_date":
		query = query.Order("library_items.added_date " + order)
	default:
		query = query.Order("library_items.created_at " + order)
	}

	// Apply pagination and preload relationships
	var items []models.LibraryItem
	err := query.
		Preload("Book").
		Preload("Book.Author").
		Preload("Book.Series").
		Offset(offset).
		Limit(limit).
		Find(&items).Error

	if err != nil {
		InternalErrorResponse(c, "Failed to fetch library items")
		return
	}

	// Convert to response format
	responseData := make([]*LibraryItemResponse, len(items))
	for i := range items {
		responseData[i] = toLibraryItemResponse(&items[i])
	}

	PaginatedSuccessResponse(c, responseData, page, limit, int(total))
}

// getLibraryItem handles GET /api/v1/library/:id
func (s *Server) getLibraryItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid library item ID")
		return
	}

	var item models.LibraryItem
	err = s.db.
		Preload("Book").
		Preload("Book.Author").
		Preload("Book.Series").
		Preload("Book.Audiobook").
		Preload("Downloads").
		First(&item, uint(id)).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFoundResponse(c, "library item")
			return
		}
		InternalErrorResponse(c, "Failed to fetch library item")
		return
	}

	SuccessResponse(c, StatusOK, toLibraryItemResponse(&item))
}

// addToLibrary handles POST /api/v1/library
func (s *Server) addToLibrary(c *gin.Context) {
	var req AddToLibraryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, err)
		return
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Find or create author
	var author models.Author
	err := tx.Where("name = ?", req.AuthorName).First(&author).Error
	if err == gorm.ErrRecordNotFound {
		author = models.Author{Name: req.AuthorName}
		if err := tx.Create(&author).Error; err != nil {
			tx.Rollback()
			InternalErrorResponse(c, "Failed to create author")
			return
		}
	} else if err != nil {
		tx.Rollback()
		InternalErrorResponse(c, "Failed to find author")
		return
	}

	// Find or create series if provided
	var seriesID *uint
	if req.SeriesName != nil && *req.SeriesName != "" {
		var series models.Series
		err := tx.Where("name = ?", *req.SeriesName).First(&series).Error
		if err == gorm.ErrRecordNotFound {
			series = models.Series{Name: *req.SeriesName}
			if err := tx.Create(&series).Error; err != nil {
				tx.Rollback()
				InternalErrorResponse(c, "Failed to create series")
				return
			}
		} else if err != nil {
			tx.Rollback()
			InternalErrorResponse(c, "Failed to find series")
			return
		}
		seriesID = &series.ID
	}

	// Check if book already exists
	var book models.Book
	bookQuery := tx.Where("title = ? AND author_id = ?", req.Title, author.ID)
	if req.ISBN != nil && *req.ISBN != "" {
		bookQuery = bookQuery.Or("isbn = ?", *req.ISBN)
	}
	if req.ASIN != nil && *req.ASIN != "" {
		bookQuery = bookQuery.Or("asin = ?", *req.ASIN)
	}

	err = bookQuery.First(&book).Error
	if err == gorm.ErrRecordNotFound {
		// Create new book
		book = models.Book{
			Title:          req.Title,
			AuthorID:       author.ID,
			SeriesID:       seriesID,
			SeriesPosition: req.SeriesPosition,
		}
		if req.ISBN != nil {
			book.ISBN = *req.ISBN
		}
		if req.ASIN != nil {
			book.ASIN = *req.ASIN
		}
		if err := tx.Create(&book).Error; err != nil {
			tx.Rollback()
			InternalErrorResponse(c, "Failed to create book")
			return
		}
	} else if err != nil {
		tx.Rollback()
		InternalErrorResponse(c, "Failed to find book")
		return
	}

	// Check if library item already exists for this book
	var existingItem models.LibraryItem
	err = tx.Where("book_id = ?", book.ID).First(&existingItem).Error
	if err == nil {
		tx.Rollback()
		ConflictResponse(c, "Book already exists in library")
		return
	} else if err != gorm.ErrRecordNotFound {
		tx.Rollback()
		InternalErrorResponse(c, "Failed to check existing library item")
		return
	}

	// Create library item
	libraryItem := models.LibraryItem{
		BookID:    book.ID,
		Status:    models.LibraryItemStatusWanted,
		AddedDate: time.Now(),
	}
	if err := tx.Create(&libraryItem).Error; err != nil {
		tx.Rollback()
		InternalErrorResponse(c, "Failed to create library item")
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		InternalErrorResponse(c, "Failed to save library item")
		return
	}

	// Reload with relationships
	err = s.db.
		Preload("Book").
		Preload("Book.Author").
		Preload("Book.Series").
		First(&libraryItem, libraryItem.ID).Error
	if err != nil {
		InternalErrorResponse(c, "Failed to reload library item")
		return
	}

	CreatedResponse(c, toLibraryItemResponse(&libraryItem))
}

// removeFromLibrary handles DELETE /api/v1/library/:id
func (s *Server) removeFromLibrary(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid library item ID")
		return
	}

	// Check if item exists
	var item models.LibraryItem
	err = s.db.First(&item, uint(id)).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFoundResponse(c, "library item")
			return
		}
		InternalErrorResponse(c, "Failed to find library item")
		return
	}

	// Soft delete (GORM handles this automatically with DeletedAt)
	err = s.db.Delete(&item).Error
	if err != nil {
		InternalErrorResponse(c, "Failed to delete library item")
		return
	}

	NoContentResponse(c)
}
