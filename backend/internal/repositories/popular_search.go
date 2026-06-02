package repositories

import (
	"time"

	"instaforge/internal/database"
	"instaforge/internal/models"

	"gorm.io/gorm"
)

// PopularSearchRepository handles popular search database operations.
type PopularSearchRepository struct {
	db *gorm.DB
}

// NewPopularSearchRepository creates a new popular search repository.
func NewPopularSearchRepository() *PopularSearchRepository {
	return &PopularSearchRepository{db: database.Get()}
}

func (r *PopularSearchRepository) Increment(username string) error {
	var search models.PopularSearch
	result := r.db.Where("username = ?", username).First(&search)
	if result.Error != nil {
		// Create new entry
		search = models.PopularSearch{
			Username:    username,
			SearchCount: 1,
		}
		return r.db.Create(&search).Error
	}

	// Increment existing
	return r.db.Model(&search).Updates(map[string]interface{}{
		"search_count":  gorm.Expr("search_count + 1"),
		"last_searched": time.Now(),
	}).Error
}

func (r *PopularSearchRepository) GetTop(limit int) ([]models.PopularSearch, error) {
	var searches []models.PopularSearch
	err := r.db.Order("search_count DESC").Limit(limit).Find(&searches).Error
	return searches, err
}

func (r *PopularSearchRepository) GetRecent(limit int) ([]models.PopularSearch, error) {
	var searches []models.PopularSearch
	err := r.db.Order("last_searched DESC").Limit(limit).Find(&searches).Error
	return searches, err
}
