package download

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/listenarr/listenarr/internal/models"
	"github.com/listenarr/listenarr/pkg/qbit"
)

// Service handles download operations
type Service struct {
	db     *gorm.DB
	qbit   *qbit.Client
	config *ServiceConfig
}

// ServiceConfig holds configuration for the download service
type ServiceConfig struct {
	Category     string
	SavePath     string
	PollInterval time.Duration
}

// NewService creates a new download service
func NewService(db *gorm.DB, qbitClient *qbit.Client, config *ServiceConfig) *Service {
	if config == nil {
		config = &ServiceConfig{
			Category:     "Listenarr",
			PollInterval: 30 * time.Second,
		}
	}
	return &Service{
		db:     db,
		qbit:   qbitClient,
		config: config,
	}
}

// StartDownload starts a download for a library item
func (s *Service) StartDownload(libraryItemID, releaseID uint) (*models.Download, error) {
	// Get release to get torrent URL
	var release models.Release
	if err := s.db.First(&release, releaseID).Error; err != nil {
		return nil, fmt.Errorf("release not found: %w", err)
	}

	// Determine torrent URL (prefer magnet, fallback to torrent URL)
	torrentURL := release.MagnetURL
	if torrentURL == "" {
		torrentURL = release.TorrentURL
	}
	if torrentURL == "" {
		return nil, fmt.Errorf("no torrent URL or magnet URL available for release")
	}

	// Create download record
	download := models.Download{
		LibraryItemID: libraryItemID,
		ReleaseID:     releaseID,
		Status:        models.DownloadStatusQueued,
		Progress:      0,
	}

	if err := s.db.Create(&download).Error; err != nil {
		return nil, fmt.Errorf("failed to create download record: %w", err)
	}

	// Add torrent to qBittorrent
	options := &qbit.AddTorrentOptions{
		Category: s.config.Category,
		SavePath: s.config.SavePath,
	}

	if err := s.qbit.AddTorrent(torrentURL, options); err != nil {
		// Update download status to failed
		download.Status = models.DownloadStatusFailed
		download.Error = fmt.Sprintf("Failed to add torrent to qBittorrent: %v", err)
		s.db.Save(&download)
		return nil, fmt.Errorf("failed to add torrent: %w", err)
	}

	// Get torrent hash from qBittorrent (we'll need to match by name or URL)
	// For now, we'll update it later when we poll
	// TODO: Get hash from qBittorrent response or match by name

	// Update library item status
	var libraryItem models.LibraryItem
	if err := s.db.First(&libraryItem, libraryItemID).Error; err == nil {
		libraryItem.Status = models.LibraryItemStatusDownloading
		s.db.Save(&libraryItem)
	}

	return &download, nil
}

// UpdateDownloadStatus updates download status from qBittorrent
func (s *Service) UpdateDownloadStatus(download *models.Download) error {
	if download.QBittorrentHash == "" {
		// Try to find torrent by matching release info
		// This is a simplified approach - in production, we'd match more reliably
		return nil
	}

	torrent, err := s.qbit.GetTorrentInfo(download.QBittorrentHash)
	if err != nil {
		return fmt.Errorf("failed to get torrent info: %w", err)
	}

	// Update download progress
	download.Progress = torrent.Progress * 100 // Convert 0-1 to 0-100
	download.Speed = torrent.DownloadSpeed
	download.Size = torrent.Size
	download.Downloaded = torrent.Downloaded

	// Update status based on qBittorrent state
	switch torrent.State {
	case "downloading", "stalledDL", "queuedDL":
		download.Status = models.DownloadStatusDownloading
	case "uploading", "stalledUP", "queuedUP":
		download.Status = models.DownloadStatusCompleted
		now := time.Now()
		download.CompletedAt = &now
	case "error":
		download.Status = models.DownloadStatusFailed
		download.Error = "qBittorrent reported error state"
	case "pausedDL", "pausedUP":
		download.Status = models.DownloadStatusPaused
	case "missingFiles":
		download.Status = models.DownloadStatusFailed
		download.Error = "Missing files"
	}

	// Update download path if available
	if torrent.ContentPath != "" {
		download.DownloadPath = torrent.ContentPath
	}

	return s.db.Save(download).Error
}

// MonitorDownloads monitors active downloads and updates their status
func (s *Service) MonitorDownloads() error {
	var downloads []models.Download
	err := s.db.Where("status IN ?", []models.DownloadStatus{
		models.DownloadStatusQueued,
		models.DownloadStatusDownloading,
	}).Find(&downloads).Error

	if err != nil {
		return fmt.Errorf("failed to fetch active downloads: %w", err)
	}

	for i := range downloads {
		if err := s.UpdateDownloadStatus(&downloads[i]); err != nil {
			// Log error but continue with other downloads
			continue
		}

		// If download completed, trigger processing
		if downloads[i].Status == models.DownloadStatusCompleted {
			s.triggerProcessing(&downloads[i])
		}
	}

	return nil
}

// triggerProcessing creates a processing task for a completed download
func (s *Service) triggerProcessing(download *models.Download) {
	// Check if processing task already exists
	var existingTask models.ProcessingTask
	err := s.db.Where("download_id = ?", download.ID).First(&existingTask).Error
	if err == nil {
		// Task already exists
		return
	}

	// Create processing task
	task := models.ProcessingTask{
		DownloadID: download.ID,
		Status:     models.ProcessingStatusPending,
		InputPath:  download.DownloadPath,
		Progress:   0,
	}

	if err := s.db.Create(&task).Error; err != nil {
		// Log error
		return
	}

	// Update library item status
	var libraryItem models.LibraryItem
	if err := s.db.First(&libraryItem, download.LibraryItemID).Error; err == nil {
		libraryItem.Status = models.LibraryItemStatusProcessing
		s.db.Save(&libraryItem)
	}
}

// CancelDownload cancels a download
func (s *Service) CancelDownload(downloadID uint) error {
	var download models.Download
	if err := s.db.First(&download, downloadID).Error; err != nil {
		return fmt.Errorf("download not found: %w", err)
	}

	// Delete from qBittorrent if hash is available
	if download.QBittorrentHash != "" {
		if err := s.qbit.DeleteTorrent([]string{download.QBittorrentHash}, false); err != nil {
			// Log error but continue
		}
	}

	// Update download status
	download.Status = models.DownloadStatusFailed
	download.Error = "Cancelled by user"
	if err := s.db.Save(&download).Error; err != nil {
		return fmt.Errorf("failed to update download: %w", err)
	}

	// Update library item status
	var libraryItem models.LibraryItem
	if err := s.db.First(&libraryItem, download.LibraryItemID).Error; err == nil {
		libraryItem.Status = models.LibraryItemStatusWanted
		s.db.Save(&libraryItem)
	}

	return nil
}
