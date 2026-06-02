package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"instaforge/internal/cache"
	"instaforge/internal/config"
	"instaforge/internal/logger"
	"instaforge/internal/models"
)

// InstagramProvider defines the interface for fetching Instagram data.
// Implementations can use scraping, third-party APIs, or Instagram's Graph API.
type InstagramProvider interface {
	GetPost(shortcode string) (*models.InstagramPost, error)
	GetReel(shortcode string) (*models.InstagramPost, error)
	GetStory(username, storyID string) (*models.InstagramStory, error)
	GetUserProfile(username string) (*models.InstagramUser, error)
	GetUserFeed(username, cursor string, limit int) (*models.InstagramFeedResponse, error)
	GetUserStories(username string) ([]models.InstagramStory, error)
}

// InstagramService orchestrates Instagram data fetching with cache-first strategy.
type InstagramService struct {
	provider   InstagramProvider
	cache      *cache.InstagramCache
	cfg        *config.Config
}

// NewInstagramService creates a new Instagram service.
func NewInstagramService(provider InstagramProvider, cfg *config.Config) *InstagramService {
	return &InstagramService{
		provider: provider,
		cache: cache.NewInstagramCache(
			time.Duration(cfg.Cache.TTLSeconds)*time.Second,
			time.Duration(cfg.Cache.UserTTLSeconds)*time.Second,
		),
		cfg: cfg,
	}
}

// GetPostByURL fetches a post by its URL, with cache-first strategy.
func (s *InstagramService) GetPostByURL(shortcode string) (*models.InstagramPost, error) {
	ctx := context.Background()
	// Cache-first
	cached, err := s.cache.GetPost(ctx, shortcode)
	if err == nil && cached != nil {
		logger.Get().Debug().Str("shortcode", shortcode).Msg("cache hit for post")
		return cached, nil
	}

	// Provider fetch
	post, err := s.provider.GetPost(shortcode)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch post %s: %w", shortcode, err)
	}

	// Cache the result
	if err := s.cache.SetPost(ctx, post); err != nil {
		logger.Get().Warn().Err(err).Str("shortcode", shortcode).Msg("failed to cache post")
	}

	return post, nil
}

// GetReelByURL fetches a reel (identical to post in terms of response shape).
func (s *InstagramService) GetReelByURL(shortcode string) (*models.InstagramPost, error) {
	// Reels use the same underlying data as posts
	return s.GetPostByURL(shortcode)
}

// GetUserProfile fetches a user profile with caching.
func (s *InstagramService) GetUserProfile(username string) (*models.InstagramUser, error) {
	ctx := context.Background()
	// Cache-first
	cached, err := s.cache.GetUserProfile(ctx, username)
	if err == nil && cached != nil {
		return cached, nil
	}

	user, err := s.provider.GetUserProfile(username)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user profile %s: %w", username, err)
	}

	if err := s.cache.SetUserProfile(ctx, user); err != nil {
		logger.Get().Warn().Err(err).Str("username", username).Msg("failed to cache user profile")
	}

	return user, nil
}

// GetUserFeed fetches a paginated user feed with caching.
func (s *InstagramService) GetUserFeed(username, cursor string, limit int) (*models.InstagramFeedResponse, error) {
	ctx := context.Background()
	if limit <= 0 || limit > 50 {
		limit = 12
	}

	// Cache-first (skip cache if cursor is provided — pagination)
	if cursor == "" {
		cached, err := s.cache.GetUserFeed(ctx, username, cursor)
		if err == nil && cached != nil {
			return cached, nil
		}
	}

	feed, err := s.provider.GetUserFeed(username, cursor, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feed for %s: %w", username, err)
	}

	// Cache first page
	if cursor == "" {
		if err := s.cache.SetUserFeed(ctx, username, cursor, feed); err != nil {
			logger.Get().Warn().Err(err).Str("username", username).Msg("failed to cache user feed")
		}
	}

	return feed, nil
}

// GetStories fetches user stories with caching.
func (s *InstagramService) GetStories(username string) ([]models.InstagramStory, error) {
	ctx := context.Background()
	if !s.cfg.Features.StoriesEnabled {
		return nil, fmt.Errorf("stories feature is disabled")
	}

	cached, err := s.cache.GetStories(ctx, username)
	if err == nil && cached != nil {
		return cached, nil
	}

	stories, err := s.provider.GetUserStories(username)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stories for %s: %w", username, err)
	}

	if err := s.cache.SetStories(ctx, username, stories); err != nil {
		logger.Get().Warn().Err(err).Str("username", username).Msg("failed to cache stories")
	}

	return stories, nil
}

// InvalidateCache clears cached data for a user.
func (s *InstagramService) InvalidateCache(username string) {
	ctx := context.Background()
	s.cache.InvalidateUser(ctx, username)
}

// === Default Provider Implementation (Third-Party API) ===

// ThirdPartyInstagramProvider fetches data from a third-party Instagram API.
type ThirdPartyInstagramProvider struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewThirdPartyInstagramProvider creates a third-party API provider.
func NewThirdPartyInstagramProvider(cfg config.InstagramConfig) *ThirdPartyInstagramProvider {
	return &ThirdPartyInstagramProvider{
		apiKey:  cfg.APIKey,
		baseURL: cfg.APIBaseURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 20,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

func (p *ThirdPartyInstagramProvider) doRequest(path string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", p.baseURL, path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("provider request to %s failed: %w", path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read provider response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("provider returned status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Provider response structures (adapt to your specific provider)
type providerPostResponse struct {
	Data struct {
		ID        string `json:"id"`
		Shortcode string `json:"shortcode"`
		Owner     struct {
			Username string `json:"username"`
		} `json:"owner"`
		Caption   string `json:"caption"`
		MediaType int    `json:"media_type"` // 1=image, 2=video, 8=carousel
		MediaURL  string `json:"media_url"`
		Thumbnail string `json:"thumbnail"`
		Likes     int    `json:"likes"`
		Comments  int    `json:"comments"`
		Timestamp string `json:"timestamp"`
		IsVideo   bool   `json:"is_video"`
		Duration  float64 `json:"duration"`
		Children  []struct {
			ID       string `json:"id"`
			MediaURL string `json:"media_url"`
			MediaType int   `json:"media_type"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
		} `json:"children"`
	} `json:"data"`
}

func (p *ThirdPartyInstagramProvider) GetPost(shortcode string) (*models.InstagramPost, error) {
	data, err := p.doRequest(fmt.Sprintf("/api/post/%s", shortcode))
	if err != nil {
		return nil, err
	}

	var resp providerPostResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse provider response: %w", err)
	}

	d := resp.Data
	ts, err := time.Parse(time.RFC3339, d.Timestamp)
	if err != nil {
		logger.Get().Warn().Err(err).Str("timestamp", d.Timestamp).Str("shortcode", d.Shortcode).Msg("invalid timestamp format from provider, using zero value")
	}

	post := &models.InstagramPost{
		ID:        d.ID,
		Shortcode: d.Shortcode,
		Username:  d.Owner.Username,
		Caption:   d.Caption,
		Thumbnail: d.Thumbnail,
		Likes:     d.Likes,
		Comments:  d.Comments,
		Timestamp: ts,
		IsVideo:   d.IsVideo,
		Duration:  d.Duration,
	}

	if d.MediaType == 8 || len(d.Children) > 0 {
		post.MediaType = models.MediaTypeCarousel
		for _, child := range d.Children {
			mt := models.MediaTypeImage
			if child.MediaType == 2 {
				mt = models.MediaTypeVideo
			}
			post.MediaItems = append(post.MediaItems, models.MediaItem{
				ID:     child.ID,
				Type:   mt,
				URL:    child.MediaURL,
				Width:  child.Width,
				Height: child.Height,
			})
		}
	} else {
		post.MediaItems = []models.MediaItem{{
			ID:   d.ID,
			URL:  d.MediaURL,
			Type: models.MediaTypeImage,
		}}
		if d.IsVideo {
			post.MediaItems[0].Type = models.MediaTypeVideo
		}
	}

	return post, nil
}

func (p *ThirdPartyInstagramProvider) GetReel(shortcode string) (*models.InstagramPost, error) {
	return p.GetPost(shortcode)
}

func (p *ThirdPartyInstagramProvider) GetStory(username, storyID string) (*models.InstagramStory, error) {
	data, err := p.doRequest(fmt.Sprintf("/api/story/%s/%s", username, storyID))
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data models.InstagramStory `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (p *ThirdPartyInstagramProvider) GetUserProfile(username string) (*models.InstagramUser, error) {
	data, err := p.doRequest(fmt.Sprintf("/api/user/%s", username))
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data models.InstagramUser `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (p *ThirdPartyInstagramProvider) GetUserFeed(username, cursor string, limit int) (*models.InstagramFeedResponse, error) {
	path := fmt.Sprintf("/api/user/%s/feed?limit=%d", username, limit)
	if cursor != "" {
		path += "&cursor=" + cursor
	}

	data, err := p.doRequest(path)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data models.InstagramFeedResponse `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (p *ThirdPartyInstagramProvider) GetUserStories(username string) ([]models.InstagramStory, error) {
	data, err := p.doRequest(fmt.Sprintf("/api/user/%s/stories", username))
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data []models.InstagramStory `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

// MockInstagramProvider provides mock data for development/testing.
type MockInstagramProvider struct{}

func NewMockInstagramProvider() *MockInstagramProvider {
	return &MockInstagramProvider{}
}

func (m *MockInstagramProvider) GetPost(shortcode string) (*models.InstagramPost, error) {
	return &models.InstagramPost{
		ID:        "12345",
		Shortcode: shortcode,
		Username:  "test_user",
		Caption:   "Test caption for Instagram post",
		MediaType: models.MediaTypeImage,
		MediaItems: []models.MediaItem{{
			ID:   "media_1",
			Type: models.MediaTypeImage,
			URL:  "https://picsum.photos/800/800",
		}},
		Likes:     150,
		Comments:  25,
		Timestamp: time.Now(),
	}, nil
}

func (m *MockInstagramProvider) GetReel(shortcode string) (*models.InstagramPost, error) {
	post, err := m.GetPost(shortcode)
	if err != nil {
		return nil, fmt.Errorf("mock reel fetch failed: %w", err)
	}
	post.IsVideo = true
	post.MediaType = models.MediaTypeVideo
	if len(post.MediaItems) > 0 {
		post.MediaItems[0].Type = models.MediaTypeVideo
	}
	post.Duration = 15.5
	return post, nil
}

func (m *MockInstagramProvider) GetStory(username, storyID string) (*models.InstagramStory, error) {
	return &models.InstagramStory{
		ID:        storyID,
		Username:  username,
		MediaType: models.MediaTypeImage,
		URL:       "https://picsum.photos/1080/1920",
		Duration:  10.0,
		Timestamp: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, nil
}

func (m *MockInstagramProvider) GetUserProfile(username string) (*models.InstagramUser, error) {
	return &models.InstagramUser{
		ID:            "user_123",
		Username:      username,
		FullName:      "Test User",
		Biography:     "This is a test Instagram user profile.",
		ProfilePicURL: "https://picsum.photos/200/200",
		FollowerCount: 1000,
		FollowingCount: 500,
		PostCount:     42,
		IsPrivate:     false,
		IsVerified:    false,
	}, nil
}

func (m *MockInstagramProvider) GetUserFeed(username, cursor string, limit int) (*models.InstagramFeedResponse, error) {
	var posts []models.InstagramPost
	for i := 0; i < limit; i++ {
		posts = append(posts, models.InstagramPost{
			ID:        fmt.Sprintf("post_%d", i),
			Shortcode: fmt.Sprintf("ABC%d", i),
			Username:  username,
			Caption:   fmt.Sprintf("Mock post %d", i+1),
			MediaType: models.MediaTypeImage,
			MediaItems: []models.MediaItem{{
				ID:   fmt.Sprintf("media_%d", i),
				Type: models.MediaTypeImage,
				URL:  fmt.Sprintf("https://picsum.photos/800/800?random=%d", i),
			}},
			Likes:     100 + i*10,
			Comments:  10 + i,
			Timestamp: time.Now().Add(-time.Duration(i) * time.Hour),
		})
	}
	return &models.InstagramFeedResponse{
		Items:   posts,
		HasMore: true,
		Cursor:  fmt.Sprintf("cursor_after_%s", username),
		User: &models.InstagramUser{
			Username: username,
			FullName: "Test User",
		},
	}, nil
}

func (m *MockInstagramProvider) GetUserStories(username string) ([]models.InstagramStory, error) {
	return []models.InstagramStory{
		{
			ID:        "story_1",
			Username:  username,
			MediaType: models.MediaTypeImage,
			URL:       "https://picsum.photos/1080/1920?random=1",
			Duration:  5.0,
			Timestamp: time.Now(),
		},
		{
			ID:        "story_2",
			Username:  username,
			MediaType: models.MediaTypeVideo,
			URL:       "https://picsum.photos/1080/1920?random=2",
			Thumbnail: "https://picsum.photos/200/200?random=2",
			Duration:  15.0,
			Timestamp: time.Now(),
		},
	}, nil
}
