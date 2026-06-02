package repositories

import (
	"instaforge/internal/database"
	"instaforge/internal/models"

	"gorm.io/gorm"
)

// AuditLogRepository handles audit log database operations.
type AuditLogRepository struct {
	db *gorm.DB
}

// NewAuditLogRepository creates a new audit log repository.
func NewAuditLogRepository() *AuditLogRepository {
	return &AuditLogRepository{db: database.Get()}
}

func (r *AuditLogRepository) Create(log *models.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *AuditLogRepository) List(limit, offset int, filters map[string]interface{}) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64
	query := r.db.Model(&models.AuditLog{})
	for k, v := range filters {
		query = query.Where(k, v)
	}
	query.Count(&total)
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}
