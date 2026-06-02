package services

import (
	"fmt"
	"time"

	"instaforge/internal/config"
	"instaforge/internal/logger"
	"instaforge/internal/models"
	"instaforge/internal/repositories"
	"instaforge/internal/validators"
)

// UserService handles user-related business logic.
type UserService struct {
	userRepo         *repositories.UserRepository
	apiKeyRepo       *repositories.APIKeyRepository
	analyticsRepo    *repositories.AnalyticsRepository
	popularSearchRepo *repositories.PopularSearchRepository
	instagramSvc     *InstagramService
	cfg              *config.Config
}

// NewUserService creates a new user service.
func NewUserService(cfg *config.Config, instagramSvc *InstagramService) *UserService {
	return &UserService{
		userRepo:         repositories.NewUserRepository(),
		apiKeyRepo:       repositories.NewAPIKeyRepository(),
		analyticsRepo:    repositories.NewAnalyticsRepository(),
		popularSearchRepo: repositories.NewPopularSearchRepository(),
		instagramSvc:     instagramSvc,
		cfg:              cfg,
	}
}

// SearchUsername performs the full username search flow.
func (s *UserService) SearchUsername(username string) (*models.UsernameSearchResult, error) {
	if err := validators.ValidateInstagramUsername(username); err != nil {
		return nil, fmt.Errorf("invalid username: %w", err)
	}

	if !s.cfg.Features.UsernameSearchEnabled {
		return nil, fmt.Errorf("username search feature is disabled")
	}

	start := time.Now()

	// Fetch profile (cache-first)
	profile, err := s.instagramSvc.GetUserProfile(username)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user profile: %w", err)
	}

	// Fetch posts (first page)
	feed, err := s.instagramSvc.GetUserFeed(username, "", 12)
	if err != nil {
		logger.Get().Warn().Err(err).Str("username", username).Msg("failed to fetch user feed, continuing with profile only")
		feed = &models.InstagramFeedResponse{Items: []models.InstagramPost{}}
	}

	// Fetch stories
	stories, err := s.instagramSvc.GetStories(username)
	if err != nil {
		logger.Get().Warn().Err(err).Str("username", username).Msg("failed to fetch stories, continuing without")
		stories = []models.InstagramStory{}
	}

	// Track popular search (async)
	go func() {
		if err := s.popularSearchRepo.Increment(username); err != nil {
			logger.Get().Warn().Err(err).Str("username", username).Msg("failed to track popular search")
		}
	}()

	elapsed := time.Since(start)
	logger.Get().Info().
		Str("username", username).
		Dur("elapsed", elapsed).
		Int("posts", len(feed.Items)).
		Int("stories", len(stories)).
		Msg("username search completed")

	result := &models.UsernameSearchResult{
		Profile:  *profile,
		Posts:    feed.Items,
		Stories:  stories,
	}

	// Mark reels separately (simplified — all video posts as reels)
	for _, post := range feed.Items {
		if post.IsVideo {
			result.Reels = append(result.Reels, post)
		}
	}

	return result, nil
}

// TrackUsage increments daily usage for a user.
func (s *UserService) TrackUsage(userID string) error {
	usage := &models.DailyUsage{
		UserID: userID,
		Date:   time.Now().Truncate(24 * time.Hour),
	}
	return s.analyticsRepo.UpsertDailyUsage(usage)
}

// GetPopularSearches returns the top searched usernames.
func (s *UserService) GetPopularSearches(limit int) ([]models.PopularSearch, error) {
	return s.popularSearchRepo.GetTop(limit)
}
