package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"instaforge/internal/models"

	"github.com/redis/go-redis/v9"
)

// InstagramCache provides cache-first access to Instagram data.
type InstagramCache struct {
	client     redis.UniversalClient
	defaultTTL time.Duration
	userTTL    time.Duration
}

// NewInstagramCache creates a new Instagram cache service.
func NewInstagramCache(defaultTTL, userTTL time.Duration) *InstagramCache {
	return &InstagramCache{
		client:     Get(),
		defaultTTL: defaultTTL,
		userTTL:    userTTL,
	}
}

// GetPost retrieves a cached post.
func (c *InstagramCache) GetPost(ctx context.Context, shortcode string) (*models.InstagramPost, error) {
	data, err := c.client.Get(ctx, KeyPostFormat(shortcode)).Bytes()
	if err != nil {
		return nil, err
	}
	var post models.InstagramPost
	err = json.Unmarshal(data, &post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// SetPost caches a post.
func (c *InstagramCache) SetPost(ctx context.Context, post *models.InstagramPost) error {
	data, err := json.Marshal(post)
	if err != nil {
		return fmt.Errorf("cache marshal error: %w", err)
	}
	return c.client.Set(ctx, KeyPostFormat(post.Shortcode), data, c.defaultTTL).Err()
}

// GetUserProfile retrieves a cached user profile.
func (c *InstagramCache) GetUserProfile(ctx context.Context, username string) (*models.InstagramUser, error) {
	data, err := c.client.Get(ctx, KeyUserProfileFormat(username)).Bytes()
	if err != nil {
		return nil, err
	}
	var user models.InstagramUser
	err = json.Unmarshal(data, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// SetUserProfile caches a user profile (longer TTL).
func (c *InstagramCache) SetUserProfile(ctx context.Context, user *models.InstagramUser) error {
	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("cache marshal error: %w", err)
	}
	return c.client.Set(ctx, KeyUserProfileFormat(user.Username), data, c.userTTL).Err()
}

// GetUserFeed retrieves a cached feed page.
func (c *InstagramCache) GetUserFeed(ctx context.Context, username, cursor string) (*models.InstagramFeedResponse, error) {
	data, err := c.client.Get(ctx, KeyUserFeedFormat(username, cursor)).Bytes()
	if err != nil {
		return nil, err
	}
	var feed models.InstagramFeedResponse
	err = json.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}
	return &feed, nil
}

// SetUserFeed caches a feed page.
func (c *InstagramCache) SetUserFeed(ctx context.Context, username, cursor string, feed *models.InstagramFeedResponse) error {
	data, err := json.Marshal(feed)
	if err != nil {
		return fmt.Errorf("cache marshal error: %w", err)
	}
	return c.client.Set(ctx, KeyUserFeedFormat(username, cursor), data, c.defaultTTL).Err()
}

// GetStories retrieves cached stories for a user.
func (c *InstagramCache) GetStories(ctx context.Context, username string) ([]models.InstagramStory, error) {
	data, err := c.client.Get(ctx, KeyUserStoriesFormat(username)).Bytes()
	if err != nil {
		return nil, err
	}
	var stories []models.InstagramStory
	err = json.Unmarshal(data, &stories)
	if err != nil {
		return nil, err
	}
	return stories, nil
}

// SetStories caches stories for a user (short TTL — stories expire fast).
func (c *InstagramCache) SetStories(ctx context.Context, username string, stories []models.InstagramStory) error {
	data, err := json.Marshal(stories)
	if err != nil {
		return fmt.Errorf("cache marshal error: %w", err)
	}
	ttl := 5 * time.Minute // stories change frequently
	return c.client.Set(ctx, KeyUserStoriesFormat(username), data, ttl).Err()
}

// InvalidateUser clears all cache entries for a user.
func (c *InstagramCache) InvalidateUser(ctx context.Context, username string) {
	if err := c.client.Del(ctx, KeyUserProfileFormat(username)).Err(); err != nil {
		log.Printf("failed to invalidate profile cache for %s: %v", username, err)
	}
	if err := c.client.Del(ctx, KeyUserStoriesFormat(username)).Err(); err != nil {
		log.Printf("failed to invalidate stories cache for %s: %v", username, err)
	}
}

// WarmUserProfile proactively caches a user profile.
func (c *InstagramCache) WarmUserProfile(ctx context.Context, user *models.InstagramUser) {
	if err := c.SetUserProfile(ctx, user); err != nil {
		log.Printf("cache warm failed for user %s: %v", user.Username, err)
	}
}
