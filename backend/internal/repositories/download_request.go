package repositories

import (
	"time"

	"instaforge/internal/database"
	"instaforge/internal/models"

	"gorm.io/gorm"
)

// DownloadRequestRepository handles download request database operations.
type DownloadRequestRepository struct {
	db *gorm.DB
}

// NewDownloadRequestRepository creates a new download request repository.
func NewDownloadRequestRepository() *DownloadRequestRepository {
	return &DownloadRequestRepository{db: database.Get()}
}

func (r *DownloadRequestRepository) Create(req *models.DownloadRequest) error {
	return r.db.Create(req).Error
}

func (r *DownloadRequestRepository) FindByID(id string) (*models.DownloadRequest, error) {
	var req models.DownloadRequest
	err := r.db.Where("id = ?", id).First(&req).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *DownloadRequestRepository) Update(req *models.DownloadRequest) error {
	return r.db.Save(req).Error
}

func (r *DownloadRequestRepository) UpdateStatus(id string, status models.DownloadStatus, errMsg string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if status == models.DownloadStatusProcessing {
		now := time.Now()
		updates["processed_at"] = &now
	}
	if status == models.DownloadStatusCompleted {
		now := time.Now()
		updates["completed_at"] = &now
	}
	if errMsg != "" {
		updates["error_message"] = errMsg
	}
	return r.db.Model(&models.DownloadRequest{}).Where("id = ?", id).Updates(updates).Error
}

func (r *DownloadRequestRepository) ListByStatus(status models.DownloadStatus, limit int) ([]models.DownloadRequest, error) {
	var reqs []models.DownloadRequest
	err := r.db.Where("status = ?", status).Order("created_at ASC").Limit(limit).Find(&reqs).Error
	return reqs, err
}

func (r *DownloadRequestRepository) CountByIP(ip string, since time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&models.DownloadRequest{}).Where("ip_address = ? AND created_at > ?", ip, since).Count(&count).Error
	return count, err
}

func (r *DownloadRequestRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.DownloadRequest{}).Error
}
