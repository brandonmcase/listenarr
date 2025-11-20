package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/listenarr/listenarr/internal/models"
)

// ProcessingTaskResponse represents a processing task in API responses
type ProcessingTaskResponse struct {
	ID          uint    `json:"id"`
	DownloadID  uint    `json:"download_id"`
	Status      string  `json:"status"`
	Progress    float64 `json:"progress"`
	InputPath   string  `json:"input_path"`
	OutputPath  string  `json:"output_path,omitempty"`
	Error       string  `json:"error,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	StartedAt   *string `json:"started_at,omitempty"`
	CompletedAt *string `json:"completed_at,omitempty"`
}

// toProcessingTaskResponse converts a ProcessingTask model to API response format
func toProcessingTaskResponse(task *models.ProcessingTask) *ProcessingTaskResponse {
	response := &ProcessingTaskResponse{
		ID:         task.ID,
		DownloadID: task.DownloadID,
		Status:     string(task.Status),
		Progress:   task.Progress,
		InputPath:  task.InputPath,
		OutputPath: task.OutputPath,
		Error:      task.Error,
		CreatedAt:  task.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  task.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if task.StartedAt != nil {
		startedAt := task.StartedAt.Format("2006-01-02T15:04:05Z07:00")
		response.StartedAt = &startedAt
	}

	if task.CompletedAt != nil {
		completedAt := task.CompletedAt.Format("2006-01-02T15:04:05Z07:00")
		response.CompletedAt = &completedAt
	}

	return response
}

// getProcessingQueue handles GET /api/v1/processing
func (s *Server) getProcessingQueue(c *gin.Context) {
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
	query := s.db.Model(&models.ProcessingTask{})

	// Apply filters
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Apply sorting (default: pending first, then by created_at)
	query = query.Order("CASE WHEN status = 'pending' THEN 0 WHEN status = 'processing' THEN 1 ELSE 2 END, created_at ASC")

	// Apply pagination and preload relationships
	var tasks []models.ProcessingTask
	err := query.
		Preload("Download").
		Preload("Download.LibraryItem").
		Preload("Download.LibraryItem.Book").
		Preload("Download.LibraryItem.Book.Author").
		Offset(offset).
		Limit(limit).
		Find(&tasks).Error

	if err != nil {
		InternalErrorResponse(c, "Failed to fetch processing queue")
		return
	}

	// Convert to response format
	responseData := make([]*ProcessingTaskResponse, len(tasks))
	for i := range tasks {
		responseData[i] = toProcessingTaskResponse(&tasks[i])
	}

	PaginatedSuccessResponse(c, responseData, page, limit, int(total))
}

// getProcessingTask handles GET /api/v1/processing/:id
func (s *Server) getProcessingTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid processing task ID")
		return
	}

	var task models.ProcessingTask
	err = s.db.
		Preload("Download").
		Preload("Download.LibraryItem").
		Preload("Download.LibraryItem.Book").
		Preload("Download.LibraryItem.Book.Author").
		First(&task, uint(id)).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFoundResponse(c, "processing task")
			return
		}
		InternalErrorResponse(c, "Failed to fetch processing task")
		return
	}

	SuccessResponse(c, StatusOK, toProcessingTaskResponse(&task))
}

// retryProcessingTask handles POST /api/v1/processing/:id/retry
func (s *Server) retryProcessingTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid processing task ID")
		return
	}

	// Check if task exists
	var task models.ProcessingTask
	err = s.db.First(&task, uint(id)).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFoundResponse(c, "processing task")
			return
		}
		InternalErrorResponse(c, "Failed to find processing task")
		return
	}

	// Only allow retrying failed tasks
	if task.Status != models.ProcessingStatusFailed {
		BadRequestResponse(c, "Can only retry failed processing tasks")
		return
	}

	// Reset task to pending
	task.Status = models.ProcessingStatusPending
	task.Progress = 0
	task.Error = ""
	task.StartedAt = nil
	task.CompletedAt = nil

	if err := s.db.Save(&task).Error; err != nil {
		InternalErrorResponse(c, "Failed to retry processing task")
		return
	}

	SuccessResponse(c, StatusOK, toProcessingTaskResponse(&task))
}
