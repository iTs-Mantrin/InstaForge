package repositories

import (
	"time"

	"instaforge/internal/database"
	"instaforge/internal/models"

	"gorm.io/gorm"
)

// AnalyticsRepository handles analytics database operations.
type AnalyticsRepository struct {
	db *gorm.DB
}

// NewAnalyticsRepository creates a new analytics repository.
func NewAnalyticsRepository() *AnalyticsRepository {
	return &AnalyticsRepository{db: database.Get()}
}

func (r *AnalyticsRepository) RecordEvent(event *models.AnalyticsEvent) error {
	return r.db.Create(event).Error
}

func (r *AnalyticsRepository) RecordEventBatch(events []models.AnalyticsEvent) error {
	if len(events) == 0 {
		return nil
	}
	return r.db.CreateInBatches(events, 100).Error
}

func (r *AnalyticsRepository) GetCountByEndpoint(since time.Time) ([]struct {
	Endpoint string
	Count    int64
}, error) {
	var results []struct {
		Endpoint string
		Count    int64
	}
	err := r.db.Model(&models.AnalyticsEvent{}).
		Select("endpoint, COUNT(*) as count").
		Where("created_at > ?", since).
		Group("endpoint").
		Order("count DESC").
		Find(&results).Error
	return results, err
}

func (r *AnalyticsRepository) GetErrorRate(since time.Time) (float64, error) {
	var total int64
	var errors int64
	r.db.Model(&models.AnalyticsEvent{}).Where("created_at > ?", since).Count(&total)
	r.db.Model(&models.AnalyticsEvent{}).Where("created_at > ? AND status_code >= 400", since).Count(&errors)
	if total == 0 {
		return 0, nil
	}
	return float64(errors) / float64(total), nil
}

func (r *AnalyticsRepository) GetDailyUsage(userID string, date time.Time) (*models.DailyUsage, error) {
	var usage models.DailyUsage
	err := r.db.Where("user_id = ? AND date = ?", userID, date.Format("2006-01-02")).First(&usage).Error
	if err != nil {
		return nil, err
	}
	return &usage, nil
}

func (r *AnalyticsRepository) UpsertDailyUsage(usage *models.DailyUsage) error {
	return r.db.Where("user_id = ? AND date = ?", usage.UserID, usage.Date.Format("2006-01-02")).
		Assign(usage).
		FirstOrCreate(usage).Error
}
