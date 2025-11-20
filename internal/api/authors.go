package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/listenarr/listenarr/internal/models"
)

// CreateAuthorRequest represents the request body for creating an author
type CreateAuthorRequest struct {
	Name        string `json:"name" binding:"required"`
	Biography   string `json:"biography,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
	GoodreadsID string `json:"goodreads_id,omitempty"`
}

// UpdateAuthorRequest represents the request body for updating an author
type UpdateAuthorRequest struct {
	Name        *string `json:"name,omitempty"`
	Biography   *string `json:"biography,omitempty"`
	ImageURL    *string `json:"image_url,omitempty"`
	GoodreadsID *string `json:"goodreads_id,omitempty"`
}

// AuthorResponseDetailed represents an author in API responses with timestamps as strings
type AuthorResponseDetailed struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Biography   string `json:"biography,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
	GoodreadsID string `json:"goodreads_id,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// AuthorWithBooksResponse represents an author with their books
type AuthorWithBooksResponse struct {
	AuthorResponseDetailed
	Books []BookResponse `json:"books,omitempty"`
}

// toAuthorResponseDetailed converts an Author model to API response format with string timestamps
func toAuthorResponseDetailed(author *models.Author) *AuthorResponseDetailed {
	return &AuthorResponseDetailed{
		ID:          author.ID,
		Name:        author.Name,
		Biography:   author.Biography,
		ImageURL:    author.ImageURL,
		GoodreadsID: author.GoodreadsID,
		CreatedAt:   author.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   author.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// toAuthorWithBooksResponse converts an Author model with books to API response format
func toAuthorWithBooksResponse(author *models.Author) *AuthorWithBooksResponse {
	response := &AuthorWithBooksResponse{
		AuthorResponseDetailed: *toAuthorResponseDetailed(author),
		Books:                  make([]BookResponse, len(author.Books)),
	}

	for i := range author.Books {
		response.Books[i] = *toBookResponse(&author.Books[i])
	}

	return response
}

// getAuthors handles GET /api/v1/authors
func (s *Server) getAuthors(c *gin.Context) {
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
	query := s.db.Model(&models.Author{})

	// Apply search filter
	if search := c.Query("search"); search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Apply sorting
	sortBy := c.DefaultQuery("sort", "name")
	order := c.DefaultQuery("order", "asc")
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	switch sortBy {
	case "name":
		query = query.Order("name " + order)
	case "created_at":
		query = query.Order("created_at " + order)
	default:
		query = query.Order("name " + order)
	}

	// Apply pagination
	var authors []models.Author
	err := query.Offset(offset).Limit(limit).Find(&authors).Error

	if err != nil {
		InternalErrorResponse(c, "Failed to fetch authors")
		return
	}

	// Convert to response format
	responseData := make([]*AuthorResponseDetailed, len(authors))
	for i := range authors {
		responseData[i] = toAuthorResponseDetailed(&authors[i])
	}

	PaginatedSuccessResponse(c, responseData, page, limit, int(total))
}

// getAuthor handles GET /api/v1/authors/:id
func (s *Server) getAuthor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid author ID")
		return
	}

	var author models.Author
	err = s.db.
		Preload("Books").
		Preload("Books.Series").
		First(&author, uint(id)).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFoundResponse(c, "author")
			return
		}
		InternalErrorResponse(c, "Failed to fetch author")
		return
	}

	SuccessResponse(c, StatusOK, toAuthorWithBooksResponse(&author))
}

// createAuthor handles POST /api/v1/authors
func (s *Server) createAuthor(c *gin.Context) {
	var req CreateAuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, err)
		return
	}

	// Check if author already exists
	var existingAuthor models.Author
	err := s.db.Where("name = ?", req.Name).First(&existingAuthor).Error
	if err == nil {
		ConflictResponse(c, "Author with this name already exists")
		return
	} else if err != gorm.ErrRecordNotFound {
		InternalErrorResponse(c, "Failed to check existing author")
		return
	}

	// Create author
	author := models.Author{
		Name:        req.Name,
		Biography:   req.Biography,
		ImageURL:    req.ImageURL,
		GoodreadsID: req.GoodreadsID,
	}

	if err := s.db.Create(&author).Error; err != nil {
		InternalErrorResponse(c, "Failed to create author")
		return
	}

	CreatedResponse(c, toAuthorResponseDetailed(&author))
}

// updateAuthor handles PUT /api/v1/authors/:id
func (s *Server) updateAuthor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid author ID")
		return
	}

	var req UpdateAuthorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, err)
		return
	}

	// Check if author exists
	var author models.Author
	err = s.db.First(&author, uint(id)).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFoundResponse(c, "author")
			return
		}
		InternalErrorResponse(c, "Failed to find author")
		return
	}

	// Update fields if provided
	if req.Name != nil {
		// Check for duplicate name if changing
		if *req.Name != author.Name {
			var existingAuthor models.Author
			err := s.db.Where("name = ? AND id != ?", *req.Name, uint(id)).First(&existingAuthor).Error
			if err == nil {
				ConflictResponse(c, "Author with this name already exists")
				return
			} else if err != gorm.ErrRecordNotFound {
				InternalErrorResponse(c, "Failed to check existing author")
				return
			}
		}
		author.Name = *req.Name
	}
	if req.Biography != nil {
		author.Biography = *req.Biography
	}
	if req.ImageURL != nil {
		author.ImageURL = *req.ImageURL
	}
	if req.GoodreadsID != nil {
		author.GoodreadsID = *req.GoodreadsID
	}

	if err := s.db.Save(&author).Error; err != nil {
		InternalErrorResponse(c, "Failed to update author")
		return
	}

	SuccessResponse(c, StatusOK, toAuthorResponseDetailed(&author))
}

// deleteAuthor handles DELETE /api/v1/authors/:id
func (s *Server) deleteAuthor(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid author ID")
		return
	}

	// Check if author exists
	var author models.Author
	err = s.db.First(&author, uint(id)).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFoundResponse(c, "author")
			return
		}
		InternalErrorResponse(c, "Failed to find author")
		return
	}

	// Check if author has books
	var bookCount int64
	s.db.Model(&models.Book{}).Where("author_id = ?", id).Count(&bookCount)
	if bookCount > 0 {
		ConflictResponse(c, "Cannot delete author with existing books")
		return
	}

	// Soft delete (GORM handles this automatically with DeletedAt)
	err = s.db.Delete(&author).Error
	if err != nil {
		InternalErrorResponse(c, "Failed to delete author")
		return
	}

	NoContentResponse(c)
}
