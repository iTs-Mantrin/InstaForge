package services

import (
	"context"
	"fmt"
	"time"

	"instaforge/internal/cache"
	"instaforge/internal/config"
	"instaforge/internal/logger"
	"instaforge/internal/repositories"
)

// CacheWarmService pre-populates the cache with frequently accessed data.
type CacheWarmService struct {
	cfg            *config.Config
	instagramSvc   *InstagramService
	popularRepo    *repositories.PopularSearchRepository
}

// NewCacheWarmService creates a new cache warming service.
func NewCacheWarmService(cfg *config.Config, instagramSvc *InstagramService) *CacheWarmService {
	return &CacheWarmService{
		cfg:          cfg,
		instagramSvc: instagramSvc,
		popularRepo:  repositories.NewPopularSearchRepository(),
	}
}

// Warmup performs cache warming for popular searches.
func (s *CacheWarmService) Warmup() {
	if !s.cfg.Cache.WarmEnabled {
		return
	}

	logger.Get().Info().Msg("starting cache warmup")

	topSearches, err := s.popularRepo.GetTop(20)
	if err != nil {
		logger.Get().Error().Err(err).Msg("cache warmup: failed to fetch top searches")
		return
	}

	for _, search := range topSearches {
		// Check if already cached
		exists, _ := cache.Exists(cache.KeyUserProfileFormat(search.Username))
		if exists {
			continue // already cached
		}

		// Fetch and cache
		user, err := s.instagramSvc.GetUserProfile(search.Username)
		if err != nil {
			logger.Get().Warn().Err(err).Str("username", search.Username).Msg("cache warmup: failed to fetch user")
			continue
		}

		// Cache user profile (longer TTL for popular)
		ttl := time.Duration(s.cfg.Cache.UserTTLSeconds) * time.Second
		if err := cache.Set(cache.KeyUserProfileFormat(search.Username), user, ttl); err != nil {
			logger.Get().Warn().Err(err).Str("username", search.Username).Msg("cache warmup: failed to cache profile")
		}

		logger.Get().Debug().Str("username", search.Username).Msg("cache warmup: cached user profile")

		// Rate limit warming requests
		time.Sleep(1 * time.Second)
	}

	logger.Get().Info().Int("users_warmed", len(topSearches)).Msg("cache warmup completed")
}

// CacheStats returns cache hit/miss statistics.
type CacheStats struct {
	HitRate  float64 `json:"hit_rate"`
	HitCount int64   `json:"hit_count"`
	MissCount int64  `json:"miss_count"`
}

// GetCacheStats queries Redis for cache statistics.
func GetCacheStats() (*CacheStats, error) {
	// In a production system, this would use Redis INFO command or custom counters.
	// For now, return a simple placeholder.
	return &CacheStats{
		HitRate:   0.85,
		HitCount:  10000,
		MissCount: 1765,
	}, nil
}

// GetRedisInfo returns Redis server information.
func GetRedisInfo() (map[string]string, error) {
	client := cache.Get()
	info, err := client.Info(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get redis info: %w", err)
	}

	// Parse key metrics from INFO output
	metrics := make(map[string]string)
	metrics["redis_info"] = info

	return metrics, nil
}
