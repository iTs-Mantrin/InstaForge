package middleware

import (
	"fmt"
	"strconv"
	"time"

	"instaforge/internal/cache"
	"instaforge/internal/config"
	"instaforge/internal/logger"
	"instaforge/internal/metrics"
	pkgresponse "instaforge/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// RateLimiter implements per-IP and per-user rate limiting using Redis.
type RateLimiter struct {
	cfg     *config.RateLimitConfig
	enabled bool
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(cfg *config.RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		cfg:     cfg,
		enabled: true,
	}
}

// Limit returns a middleware that enforces rate limits.
func (rl *RateLimiter) Limit() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !rl.enabled {
			return c.Next()
		}

		identifier := c.IP()
		// If user is authenticated, use their ID for more granular limits
		if userID, ok := c.Locals("user_id").(string); ok && userID != "" {
			identifier = userID
		}

		endpoint := c.Route().Path
		key := fmt.Sprintf("ratelimit:%s:%s", identifier, endpoint)

		// Increment counter with TTL
		count, err := cache.IncrWithTTL(key, time.Duration(rl.cfg.Window)*time.Second)
		if err != nil {
			logger.Get().Error().Err(err).Msg("rate limiter: redis error — fail closed")
			return pkgresponse.ServiceUnavailable(c, "rate limiting unavailable; try again later")
		}

		// Set rate limit headers
		remaining := rl.cfg.Requests - int(count)
		if remaining < 0 {
			remaining = 0
		}
		c.Response().Header.Set("X-RateLimit-Limit", strconv.Itoa(rl.cfg.Requests))
		c.Response().Header.Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Response().Header.Set("X-RateLimit-Reset", strconv.Itoa(int(time.Now().Add(time.Duration(rl.cfg.Window)*time.Second).Unix())))

		if count > int64(rl.cfg.Requests) {
			metrics.RecordRateLimitHit()
			retryAfter := int(rl.cfg.Window.Seconds())
			c.Response().Header.Set("Retry-After", strconv.Itoa(retryAfter))
			return pkgresponse.RateLimited(c, fmt.Sprintf("rate limit exceeded. retry after %d seconds", retryAfter))
		}

		return c.Next()
	}
}

// Disable disables rate limiting (for maintenance or testing).
func (rl *RateLimiter) Disable() {
	rl.enabled = false
}

// Enable enables rate limiting.
func (rl *RateLimiter) Enable() {
	rl.enabled = true
}
