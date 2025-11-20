package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/listenarr/listenarr/internal/auth"
	"github.com/listenarr/listenarr/internal/config"
)

// Server represents the API server
type Server struct {
	config *config.Config
	db     *gorm.DB
	router *gin.Engine
}

// NewServer creates a new API server instance
func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	// Set Gin mode based on environment
	if cfg.Server.Host == "0.0.0.0" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	server := &Server{
		config: cfg,
		db:     db,
		router: router,
	}

	server.setupRoutes()

	return server
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// Health check endpoint (no auth required)
	s.router.GET("/api/health", s.healthCheck)

	// Apply authentication middleware if enabled
	if s.config.Auth.Enabled && s.config.Auth.APIKey != "" {
		s.router.Use(auth.APIKeyMiddleware(s.config.Auth.APIKey))
	}

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Library routes
		v1.GET("/library", s.getLibrary)
		v1.GET("/library/:id", s.getLibraryItem)
		v1.POST("/library", s.addToLibrary)
		v1.DELETE("/library/:id", s.removeFromLibrary)

		// Author routes
		v1.GET("/authors", s.getAuthors)
		v1.GET("/authors/:id", s.getAuthor)
		v1.POST("/authors", s.createAuthor)
		v1.PUT("/authors/:id", s.updateAuthor)
		v1.DELETE("/authors/:id", s.deleteAuthor)

		// Book routes
		v1.GET("/books", s.getBooks)
		v1.GET("/books/:id", s.getBook)
		v1.POST("/books", s.createBook)
		v1.PUT("/books/:id", s.updateBook)
		v1.DELETE("/books/:id", s.deleteBook)

		// Download routes
		v1.GET("/downloads", s.getDownloads)
		v1.GET("/downloads/:id", s.getDownload)
		v1.POST("/downloads", s.startDownload)
		v1.DELETE("/downloads/:id", s.cancelDownload)

		// Processing routes
		v1.GET("/processing", s.getProcessingQueue)
		v1.GET("/processing/:id", s.getProcessingTask)
		v1.POST("/processing/:id/retry", s.retryProcessingTask)

		// Search routes
		v1.GET("/search", s.searchAudiobooks)
	}
}

// healthCheck returns the health status of the API
func (s *Server) healthCheck(c *gin.Context) {
	// Health check uses simple format (not standard API response format)
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "listenarr",
	})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	return s.router.Run(addr)
}

// All handlers are implemented in separate files:
// - Library handlers: library.go
// - Author handlers: authors.go
// - Book handlers: books.go
// - Download handlers: downloads.go
// - Processing handlers: processing.go
// - Search handler: search.go
