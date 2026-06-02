package repositories

import (
	"time"

	"instaforge/internal/database"
	"instaforge/internal/models"

	"gorm.io/gorm"
)

// APIKeyRepository handles API key database operations.
type APIKeyRepository struct {
	db *gorm.DB
}

// NewAPIKeyRepository creates a new API key repository.
func NewAPIKeyRepository() *APIKeyRepository {
	return &APIKeyRepository{db: database.Get()}
}

func (r *APIKeyRepository) Create(key *models.APIKey) error {
	return r.db.Create(key).Error
}

func (r *APIKeyRepository) FindByKey(key string) (*models.APIKey, error) {
	var apiKey models.APIKey
	err := r.db.Where("key = ? AND status = ?", key, "active").First(&apiKey).Error
	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

func (r *APIKeyRepository) FindByUserID(userID string) ([]models.APIKey, error) {
	var keys []models.APIKey
	err := r.db.Where("user_id = ?", userID).Find(&keys).Error
	return keys, err
}

func (r *APIKeyRepository) UpdateLastUsed(id string) error {
	return r.db.Model(&models.APIKey{}).Where("id = ?", id).Update("last_used", time.Now()).Error
}

func (r *APIKeyRepository) Revoke(id string) error {
	return r.db.Model(&models.APIKey{}).Where("id = ?", id).Update("status", "revoked").Error
}
