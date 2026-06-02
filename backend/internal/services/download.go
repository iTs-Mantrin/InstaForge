package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"instaforge/internal/cache"
	"instaforge/internal/config"
	"instaforge/internal/logger"
	"instaforge/internal/models"
	"instaforge/internal/queue"
	"instaforge/internal/repositories"
	"instaforge/internal/validators"

	"github.com/google/uuid"
)

// DownloadService handles the download pipeline.
type DownloadService struct {
	cfg         *config.Config
	repo        *repositories.DownloadRequestRepository
	cache       *cache.InstagramCache
	instaSvc    *InstagramService
	httpClient  *http.Client
	mu          sync.RWMutex
	activeJobs  map[string]chan struct{}
}

// NewDownloadService creates a new download service.
func NewDownloadService(cfg *config.Config, instaSvc *InstagramService) *DownloadService {
	return &DownloadService{
		cfg:        cfg,
		repo:       repositories.NewDownloadRequestRepository(),
		cache: cache.NewInstagramCache(
			time.Duration(cfg.Cache.TTLSeconds)*time.Second,
			time.Duration(cfg.Cache.UserTTLSeconds)*time.Second,
		),
		instaSvc: instaSvc,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:    50,
				IdleConnTimeout: 90 * time.Second,
			},
		},
		activeJobs: make(map[string]chan struct{}),
	}
}

// QueueDownload creates a download request and queues it for processing.
func (s *DownloadService) QueueDownload(url, ipAddress, userAgent string, userID *string) (*models.DownloadRequest, error) {
	// Validate URL
	urlType, shortcode, err := validators.ValidateInstagramURL(url)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Create download request record
	req := &models.DownloadRequest{
		ID:        uuid.New().String(),
		URL:       url,
		URLType:   string(urlType),
		Shortcode: shortcode,
		Status:    models.DownloadStatusPending,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		UserID:    userID,
	}

	if err := s.repo.Create(req); err != nil {
		return nil, fmt.Errorf("failed to create download request: %w", err)
	}

	// Queue to Kafka for async processing
	kafkaMsg := &queue.Message{
		ID:        req.ID,
		Type:      "download.requested",
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"request_id": req.ID,
			"url":        url,
			"type":       string(urlType),
			"shortcode":  shortcode,
		},
	}

	if err := queue.Publish("download_requests", kafkaMsg); err != nil {
		logger.Get().Error().Err(err).Str("request_id", req.ID).Msg("failed to queue download request")
		// Don't fail — the request is already in the database
	}

	return req, nil
}

// ProcessDownload processes a queued download request (called by workers).
func (s *DownloadService) ProcessDownload(requestID string) error {
	req, err := s.repo.FindByID(requestID)
	if err != nil {
		return fmt.Errorf("download request not found: %s", requestID)
	}

	// Update status to processing
	if err := s.repo.UpdateStatus(req.ID, models.DownloadStatusProcessing, ""); err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	// Fetch post data based on URL type
	var post *models.InstagramPost
	switch req.URLType {
	case "post":
		post, err = s.instaSvc.GetPostByURL(req.Shortcode)
	case "reel":
		post, err = s.instaSvc.GetReelByURL(req.Shortcode)
	default:
		err = fmt.Errorf("unsupported URL type: %s", req.URLType)
	}
	if err != nil {
		if updateErr := s.repo.UpdateStatus(req.ID, models.DownloadStatusFailed, err.Error()); updateErr != nil {
			logger.Get().Error().Err(updateErr).Str("request_id", req.ID).Msg("failed to update download status to failed")
		}
		return fmt.Errorf("failed to fetch media: %w", err)
	}

	// For carousel, get all media
	// For single, get the first item
	// Cache the download URL for later retrieval
	downloadData := map[string]interface{}{
		"media_items": post.MediaItems,
		"username":    post.Username,
		"caption":     post.Caption,
		"shortcode":   post.Shortcode,
	}

	cacheKey := cache.KeyDownloadURLFormat(req.ID)
	if err := cache.Set(cacheKey, downloadData, 30*time.Minute); err != nil {
		logger.Get().Warn().Err(err).Msg("failed to cache download data")
	}

	// Update request as completed
	if err := s.repo.UpdateStatus(req.ID, models.DownloadStatusCompleted, ""); err != nil {
		return fmt.Errorf("failed to mark download completed: %w", err)
	}

	// Update media count
	req.MediaCount = len(post.MediaItems)
	if updateErr := s.repo.Update(req); updateErr != nil {
		logger.Get().Error().Err(updateErr).Str("request_id", req.ID).Msg("failed to update download request metadata")
	}

	logger.Get().Info().
		Str("request_id", requestID).
		Str("shortcode", req.Shortcode).
		Int("media_count", len(post.MediaItems)).
		Msg("download processed successfully")

	return nil
}

// GetDownloadResult retrieves the result of a completed download.
func (s *DownloadService) GetDownloadResult(requestID string) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := cache.GetJSON(cache.KeyDownloadURLFormat(requestID), &data)
	if err != nil {
		return nil, fmt.Errorf("download result not found or expired")
	}
	return data, nil
}

// StreamMedia streams a media file to the response writer.
func (s *DownloadService) StreamMedia(mediaURL string) (io.ReadCloser, string, int64, error) {
	resp, err := s.httpClient.Get(mediaURL)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to fetch media: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, "", 0, fmt.Errorf("media server returned status %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	contentLength := resp.ContentLength

	return resp.Body, contentType, contentLength, nil
}

// GetDownloadFileResult returns the download URL for a completed download.
func (s *DownloadService) GetDownloadFileResult(requestID string) (string, error) {
	data, err := s.GetDownloadResult(requestID)
	if err != nil {
		return "", err
	}

	// Try to extract a media URL from the cached download data
	if items, ok := data["media_items"].([]interface{}); ok && len(items) > 0 {
		if first, ok := items[0].(map[string]interface{}); ok {
			if url, ok := first["url"].(string); ok && url != "" {
				return url, nil
			}
		}
	}

	return "", fmt.Errorf("no media URL available for download: %s", requestID)
}

// GetDownloadProgress returns the current status/progress of a download for the frontend.
func (s *DownloadService) GetDownloadProgress(requestID string) (*models.DownloadRequest, error) {
	req, err := s.repo.FindByID(requestID)
	if err != nil {
		return nil, fmt.Errorf("download request not found: %s", requestID)
	}
	return req, nil
}

// CancelJob cancels a queued or in-progress download.
func (s *DownloadService) CancelJob(requestID string) error {
	req, err := s.repo.FindByID(requestID)
	if err != nil {
		return fmt.Errorf("download request not found: %s", requestID)
	}
	if req.Status == models.DownloadStatusCompleted || req.Status == models.DownloadStatusExpired {
		return fmt.Errorf("download already completed, cannot cancel")
	}
	if err := s.repo.UpdateStatus(requestID, models.DownloadStatusCancelled, "cancelled by user"); err != nil {
		return fmt.Errorf("failed to cancel download: %w", err)
	}

	// Signal active job if being processed
	s.mu.RLock()
	cancelCh, exists := s.activeJobs[requestID]
	s.mu.RUnlock()
	if exists {
		select {
		case cancelCh <- struct{}{}:
		default:
		}
	}

	return nil
}

// CleanupExpired removes expired download records and files.
func (s *DownloadService) CleanupExpired() error {
	return s.repo.DeleteExpired()
}

// SaveMediaToTemp stores media data to a temporary file.
func (s *DownloadService) SaveMediaToTemp(data io.Reader, ext string) (string, error) {
	if err := os.MkdirAll(s.cfg.Storage.TempDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	filename := fmt.Sprintf("%s_%s%s", uuid.New().String(), time.Now().Format("20060102150405"), ext)
	filepath := filepath.Join(s.cfg.Storage.TempDir, filename)

	file, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, data)
	if err != nil {
		os.Remove(filepath)
		return "", fmt.Errorf("failed to write temp file: %w", err)
	}

	return filepath, nil
}
