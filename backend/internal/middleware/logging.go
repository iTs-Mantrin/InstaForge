package middleware

import (
	"time"

	"instaforge/internal/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
)

// RequestID adds a unique request ID to each request.
func RequestID() fiber.Handler {
	return requestid.New(requestid.Config{
		Header: "X-Request-ID",
		Generator: func() string {
			return uuid.New().String()
		},
	})
}

// RequestLogger logs each HTTP request with structured details.
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)
		durationMs := duration.Milliseconds()

		// Store for response metadata
		c.Locals("took_ms", durationMs)

		// Enrich with request context
		requestID, _ := c.Locals("request_id").(string)

		// Structured logging
		log := logger.WithRequest(
			c.Method(),
			c.Path(),
			requestID,
			c.IP(),
		)

		log.Info().
			Int("status", c.Response().StatusCode()).
			Int64("duration_ms", durationMs).
			Int("body_size", len(c.Response().Body())).
			Str("user_agent", c.Get("User-Agent")).
			Str("referer", c.Get("Referer")).
			Msg("request completed")
		return nil
	}
}

// RecoverMiddleware handles panics gracefully and logs them.
func RecoverMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
		requestID, _ := c.Locals("request_id").(string)
				log := logger.WithRequest(c.Method(), c.Path(), requestID, c.IP())
				log.Error().
					Interface("panic", r).
					Msg("panic recovered")

				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"success": false,
					"error": fiber.Map{
						"code":    "INTERNAL_ERROR",
						"message": "an unexpected error occurred",
					},
				})
			}
		}()
		return c.Next()
	}
}
