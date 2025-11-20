package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	
	"github.com/listenarr/listenarr/internal/config"
	"github.com/listenarr/listenarr/internal/auth"
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
		v1.POST("/library", s.addToLibrary)
		v1.DELETE("/library/:id", s.removeFromLibrary)

		// Download routes
		v1.GET("/downloads", s.getDownloads)
		v1.POST("/downloads", s.startDownload)

		// Processing routes
		v1.GET("/processing", s.getProcessingQueue)

		// Search routes
		v1.GET("/search", s.searchAudiobooks)
	}
}

// healthCheck returns the health status of the API
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "listenarr",
	})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	return s.router.Run(addr)
}

// Placeholder handlers (to be implemented)
func (s *Server) getLibrary(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Library endpoint - to be implemented"})
}

func (s *Server) addToLibrary(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Add to library - to be implemented"})
}

func (s *Server) removeFromLibrary(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Remove from library - to be implemented"})
}

func (s *Server) getDownloads(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Downloads endpoint - to be implemented"})
}

func (s *Server) startDownload(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Start download - to be implemented"})
}

func (s *Server) getProcessingQueue(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Processing queue - to be implemented"})
}

func (s *Server) searchAudiobooks(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Search - to be implemented"})
}

