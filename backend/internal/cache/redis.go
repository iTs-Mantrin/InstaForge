package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"instaforge/internal/config"

	"github.com/redis/go-redis/v9"
)

var (
	client redis.UniversalClient
	once   sync.Once
	ctx    = context.Background()
)

// Connect initializes the Redis cluster client.
func Connect(cfg config.RedisConfig) (redis.UniversalClient, error) {
	var err error
	once.Do(func() {
		opts := &redis.ClusterOptions{
			Addrs:         cfg.Addrs,
			Password:      cfg.Password,
			PoolSize:      cfg.PoolSize,
			MinIdleConns:  cfg.MinIdleConns,
			DialTimeout:   cfg.DialTimeout,
			ReadTimeout:   cfg.ReadTimeout,
			WriteTimeout:  cfg.WriteTimeout,
			ReadOnly:      true,
			RouteRandomly: true,
		}

		client = redis.NewClusterClient(opts)

		// Verify connectivity
		if pingErr := client.Ping(ctx).Err(); pingErr != nil {
			err = fmt.Errorf("redis cluster connection failed: %w", pingErr)
			return
		}
	})

	return client, err
}

// ConnectSingle initializes a single-node Redis client (for development).
func ConnectSingle(cfg config.RedisConfig) (redis.UniversalClient, error) {
	var err error
	once.Do(func() {
		opts := &redis.Options{
			Addr:         cfg.Addrs[0],
			Password:     cfg.Password,
			DB:           cfg.DB,
			PoolSize:     cfg.PoolSize,
			MinIdleConns: cfg.MinIdleConns,
			DialTimeout:  cfg.DialTimeout,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		}

		client = redis.NewClient(opts)

		if pingErr := client.Ping(ctx).Err(); pingErr != nil {
			err = fmt.Errorf("redis connection failed: %w", pingErr)
			return
		}
	})

	return client, err
}

// Get returns the Redis client instance.
func Get() redis.UniversalClient {
	if client == nil {
		panic("redis not initialized — call cache.Connect() first")
	}
	return client
}

// Set stores a value in the cache with TTL.
func Set(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache marshal error: %w", err)
	}
	return client.Set(ctx, key, data, ttl).Err()
}

// GetJSON retrieves and deserializes a cached value.
func GetJSON(key string, dest interface{}) error {
	data, err := client.Get(ctx, key).Bytes()
	if err != nil {
		return err // includes redis.Nil
	}
	return json.Unmarshal(data, dest)
}

// Delete removes a key from cache.
func Delete(key string) error {
	return client.Del(ctx, key).Err()
}

// Exists checks if a key exists in cache.
func Exists(key string) (bool, error) {
	n, err := client.Exists(ctx, key).Result()
	return n > 0, err
}

// Keys returns all keys matching a pattern (use with caution on clusters).
func Keys(pattern string) ([]string, error) {
	return client.Keys(ctx, pattern).Result()
}

// Incr increments a counter and returns the new value.
func Incr(key string) (int64, error) {
	return client.Incr(ctx, key).Result()
}

// IncrWithTTL increments a counter and sets TTL on first creation.
func IncrWithTTL(key string, ttl time.Duration) (int64, error) {
	v, err := client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if v == 1 {
		if err := client.Expire(ctx, key, ttl).Err(); err != nil {
			fmt.Printf("failed to set TTL on rate limit key %s: %v", key, err)
		}
	}
	return v, nil
}

// Expire sets TTL on a key.
func Expire(key string, ttl time.Duration) error {
	return client.Expire(ctx, key, ttl).Err()
}

// TTL returns the remaining TTL for a key.
func TTL(key string) (time.Duration, error) {
	return client.TTL(ctx, key).Result()
}

// Pipeline creates a pipeline for batch operations.
func Pipeline() redis.Pipeliner {
	return client.Pipeline()
}

// HealthCheck verifies Redis connectivity.
func HealthCheck() error {
	return client.Ping(ctx).Err()
}

// Close gracefully closes the Redis connection.
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}

// Cache keys constants
const (
	KeyPost           = "post:%s"           // post:{shortcode}
	KeyUserProfile    = "user:%s"           // user:{username}
	KeyUserFeed       = "feed:%s:%s"        // feed:{username}:{cursor}
	KeyUserStories    = "stories:%s"        // stories:{username}
	KeyDownloadURL    = "download:%s"       // download:{request_id}
	KeyRateLimit      = "ratelimit:%s:%s"   // ratelimit:{identifier}:{endpoint}
	KeyAnalyticsDaily = "analytics:daily:%s" // analytics:daily:{date}
)

// KeyPostFormat returns the cache key for a post.
func KeyPostFormat(shortcode string) string {
	return fmt.Sprintf(KeyPost, shortcode)
}

// KeyUserProfileFormat returns the cache key for a user profile.
func KeyUserProfileFormat(username string) string {
	return fmt.Sprintf(KeyUserProfile, username)
}

// KeyUserFeedFormat returns the cache key for a user feed.
func KeyUserFeedFormat(username, cursor string) string {
	return fmt.Sprintf(KeyUserFeed, username, cursor)
}

// KeyUserStoriesFormat returns the cache key for user stories.
func KeyUserStoriesFormat(username string) string {
	return fmt.Sprintf(KeyUserStories, username)
}

// KeyDownloadURLFormat returns the cache key for a download request.
func KeyDownloadURLFormat(requestID string) string {
	return fmt.Sprintf(KeyDownloadURL, requestID)
}
