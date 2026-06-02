package database

import (
	"fmt"

	"instaforge/internal/logger"
	"instaforge/internal/models"

	"gorm.io/gorm"
)

// Migration represents a database migration.
type Migration struct {
	ID      string
	Migrate func(tx *gorm.DB) error
}

// RunMigrations executes all pending migrations.
func RunMigrations(db *gorm.DB) error {
	log := logger.Get()
	log.Info().Msg("running database migrations")

	// Enable UUID extension
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		return fmt.Errorf("failed to enable uuid-ossp extension: %w", err)
	}

	// Enable pg_trgm for text search
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "pg_trgm"`).Error; err != nil {
		return fmt.Errorf("failed to enable pg_trgm extension: %w", err)
	}

	// Auto-migrate models
	modelsToMigrate := []interface{}{
		&models.User{},
		&models.APIKey{},
		&models.DownloadRequest{},
		&models.AnalyticsEvent{},
		&models.AnalyticsSummary{},
		&models.DailyUsage{},
		&models.AuditLog{},
		&models.PopularSearch{},
	}

	for _, model := range modelsToMigrate {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %w", model, err)
		}
	}

	// Create indexes
	if err := createIndexes(db); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	log.Info().Msg("database migrations completed successfully")
	return nil
}

func createIndexes(db *gorm.DB) error {
	indexes := []string{
		// Analytics
		`CREATE INDEX IF NOT EXISTS idx_analytics_events_created_at_brin ON analytics_events USING BRIN(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_analytics_events_endpoint_created ON analytics_events(endpoint, created_at)`,

		// Download requests
		`CREATE INDEX IF NOT EXISTS idx_download_requests_status_created ON download_requests(status, created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_download_requests_ip_created ON download_requests(ip_address, created_at)`,

		// Audit logs
		`CREATE INDEX IF NOT EXISTS idx_audit_logs_action_created ON audit_logs(action, created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id)`,

		// Popular searches
		`CREATE INDEX IF NOT EXISTS idx_popular_searches_count ON popular_searches(search_count DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_popular_searches_last_searched ON popular_searches(last_searched DESC)`,

		// Text search with pg_trgm
		`CREATE INDEX IF NOT EXISTS idx_popular_searches_username_trgm ON popular_searches USING GIN(username gin_trgm_ops)`,

		// API keys
		`CREATE INDEX IF NOT EXISTS idx_api_keys_user_status ON api_keys(user_id, status)`,
	}

	for _, idx := range indexes {
		if err := db.Exec(idx).Error; err != nil {
			logger.Get().Warn().Err(err).Str("index", idx).Msg("index creation warning (may already exist)")
		}
	}

	return nil
}
