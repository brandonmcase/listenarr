package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/listenarr/listenarr/internal/models"
)

// StartDownloadRequest represents the request body for starting a download
type StartDownloadRequest struct {
	LibraryItemID uint `json:"library_item_id" binding:"required"`
	ReleaseID     uint `json:"release_id" binding:"required"`
}

// DownloadResponse represents a download in API responses
type DownloadResponse struct {
	ID              uint    `json:"id"`
	LibraryItemID   uint    `json:"library_item_id"`
	ReleaseID       uint    `json:"release_id"`
	Status          string  `json:"status"`
	Progress        float64 `json:"progress"`
	Speed           int64   `json:"speed,omitempty"`
	Size            int64   `json:"size,omitempty"`
	Downloaded      int64   `json:"downloaded,omitempty"`
	Error           string  `json:"error,omitempty"`
	QBittorrentHash string  `json:"qbittorrent_hash,omitempty"`
	DownloadPath    string  `json:"download_path,omitempty"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
	CompletedAt     *string `json:"completed_at,omitempty"`
}

// toDownloadResponse converts a Download model to API response format
func toDownloadResponse(download *models.Download) *DownloadResponse {
	response := &DownloadResponse{
		ID:              download.ID,
		LibraryItemID:   download.LibraryItemID,
		ReleaseID:       download.ReleaseID,
		Status:          string(download.Status),
		Progress:        download.Progress,
		Speed:           download.Speed,
		Size:            download.Size,
		Downloaded:      download.Downloaded,
		Error:           download.Error,
		QBittorrentHash: download.QBittorrentHash,
		DownloadPath:    download.DownloadPath,
		CreatedAt:       download.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       download.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if download.CompletedAt != nil {
		completedAt := download.CompletedAt.Format("2006-01-02T15:04:05Z07:00")
		response.CompletedAt = &completedAt
	}

	return response
}

// getDownloads handles GET /api/v1/downloads
func (s *Server) getDownloads(c *gin.Context) {
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
	query := s.db.Model(&models.Download{})

	// Apply filters
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if libraryItemIDStr := c.Query("library_item_id"); libraryItemIDStr != "" {
		if libraryItemID, err := strconv.ParseUint(libraryItemIDStr, 10, 32); err == nil {
			query = query.Where("library_item_id = ?", uint(libraryItemID))
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

	switch sortBy {
	case "status":
		query = query.Order("status " + order)
	case "progress":
		query = query.Order("progress " + order)
	case "created_at":
		query = query.Order("created_at " + order)
	default:
		query = query.Order("created_at " + order)
	}

	// Apply pagination and preload relationships
	var downloads []models.Download
	err := query.
		Preload("LibraryItem").
		Preload("LibraryItem.Book").
		Preload("LibraryItem.Book.Author").
		Preload("Release").
		Offset(offset).
		Limit(limit).
		Find(&downloads).Error

	if err != nil {
		InternalErrorResponse(c, "Failed to fetch downloads")
		return
	}

	// Convert to response format
	responseData := make([]*DownloadResponse, len(downloads))
	for i := range downloads {
		responseData[i] = toDownloadResponse(&downloads[i])
	}

	PaginatedSuccessResponse(c, responseData, page, limit, int(total))
}

// getDownload handles GET /api/v1/downloads/:id
func (s *Server) getDownload(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid download ID")
		return
	}

	var download models.Download
	err = s.db.
		Preload("LibraryItem").
		Preload("LibraryItem.Book").
		Preload("LibraryItem.Book.Author").
		Preload("Release").
		First(&download, uint(id)).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFoundResponse(c, "download")
			return
		}
		InternalErrorResponse(c, "Failed to fetch download")
		return
	}

	SuccessResponse(c, StatusOK, toDownloadResponse(&download))
}

// startDownload handles POST /api/v1/downloads
// Note: This is a placeholder implementation. Full qBittorrent integration will be added later.
func (s *Server) startDownload(c *gin.Context) {
	var req StartDownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(c, err)
		return
	}

	// Verify library item exists
	var libraryItem models.LibraryItem
	err := s.db.First(&libraryItem, req.LibraryItemID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFoundResponse(c, "library item")
			return
		}
		InternalErrorResponse(c, "Failed to find library item")
		return
	}

	// Verify release exists
	var release models.Release
	err = s.db.First(&release, req.ReleaseID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFoundResponse(c, "release")
			return
		}
		InternalErrorResponse(c, "Failed to find release")
		return
	}

	// Check if there's already an active download for this library item
	var existingDownload models.Download
	err = s.db.Where("library_item_id = ? AND status IN ?", req.LibraryItemID, []models.DownloadStatus{
		models.DownloadStatusQueued,
		models.DownloadStatusDownloading,
	}).First(&existingDownload).Error
	if err == nil {
		ConflictResponse(c, "Active download already exists for this library item")
		return
	} else if err != gorm.ErrRecordNotFound {
		InternalErrorResponse(c, "Failed to check existing downloads")
		return
	}

	// Create download
	download := models.Download{
		LibraryItemID: req.LibraryItemID,
		ReleaseID:     req.ReleaseID,
		Status:        models.DownloadStatusQueued,
		Progress:      0,
	}

	if err := s.db.Create(&download).Error; err != nil {
		InternalErrorResponse(c, "Failed to create download")
		return
	}

	// Update library item status
	libraryItem.Status = models.LibraryItemStatusDownloading
	s.db.Save(&libraryItem)

	// TODO: Integrate with qBittorrent service to actually start the download
	// For now, we just create the download record

	// Reload with relationships
	err = s.db.
		Preload("LibraryItem").
		Preload("Release").
		First(&download, download.ID).Error
	if err != nil {
		InternalErrorResponse(c, "Failed to reload download")
		return
	}

	CreatedResponse(c, toDownloadResponse(&download))
}

// cancelDownload handles DELETE /api/v1/downloads/:id
func (s *Server) cancelDownload(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid download ID")
		return
	}

	// Check if download exists
	var download models.Download
	err = s.db.First(&download, uint(id)).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFoundResponse(c, "download")
			return
		}
		InternalErrorResponse(c, "Failed to find download")
		return
	}

	// Only allow canceling queued or downloading status
	if download.Status != models.DownloadStatusQueued && download.Status != models.DownloadStatusDownloading {
		BadRequestResponse(c, "Can only cancel queued or downloading downloads")
		return
	}

	// Update status to failed (or we could add a "cancelled" status)
	download.Status = models.DownloadStatusFailed
	download.Error = "Download cancelled by user"
	if err := s.db.Save(&download).Error; err != nil {
		InternalErrorResponse(c, "Failed to cancel download")
		return
	}

	// Update library item status back to wanted
	var libraryItem models.LibraryItem
	if err := s.db.First(&libraryItem, download.LibraryItemID).Error; err == nil {
		libraryItem.Status = models.LibraryItemStatusWanted
		s.db.Save(&libraryItem)
	}

	NoContentResponse(c)
}
