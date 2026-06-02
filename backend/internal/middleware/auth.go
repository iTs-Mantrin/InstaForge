package middleware

import (
	"strings"
	"time"

	"instaforge/internal/logger"
	"instaforge/internal/repositories"
	"instaforge/internal/security"
	pkgresponse "instaforge/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// JWTAuthMiddleware validates JWT tokens on protected routes.
func JWTAuthMiddleware(jwtManager *security.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return pkgresponse.Unauthorized(c, "missing authorization header")
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return pkgresponse.Unauthorized(c, "invalid authorization format, use Bearer <token>")
		}

		token := parts[1]
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			logger.Get().Warn().Err(err).Str("token_prefix", token[:min(10, len(token))]+"...").Msg("jwt validation failed")
			return pkgresponse.Unauthorized(c, "invalid or expired token")
		}

		// Store claims in context
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("user_role", claims.Role)
		c.Locals("user_scopes", claims.Scopes)

		return c.Next()
	}
}

// APIKeyMiddleware validates API key authentication and enriches context.
func APIKeyMiddleware(repo *repositories.APIKeyRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}
		if apiKey == "" {
			// Optional — skip validation, let rate limiter handle anonymous
			return c.Next()
		}

		// Hash the raw key for database lookup (DB stores SHA-256 hex)
		hasher := security.NewAPIKeyGenerator()
		hashedKey := hasher.Hash(apiKey)

		// Validate against database
		keyRecord, err := repo.FindByKey(hashedKey)
		if err != nil {
			logger.Get().Warn().Str("api_key_prefix", apiKey[:min(8, len(apiKey))]+"...").Msg("invalid API key attempt")
			return pkgresponse.Unauthorized(c, "invalid API key")
		}

		// Check if key has expired
		if keyRecord.ExpiresAt != nil && keyRecord.ExpiresAt.Before(time.Now()) {
			logger.Get().Warn().Str("key_id", keyRecord.ID).Msg("expired API key attempt")
			return pkgresponse.Unauthorized(c, "API key has expired")
		}

		// Update last used timestamp (async, non-blocking)
		go func() {
			if updateErr := repo.UpdateLastUsed(keyRecord.ID); updateErr != nil {
				logger.Get().Error().Err(updateErr).Str("key_id", keyRecord.ID).Msg("failed to update API key last_used")
			}
		}()

		// Enrich context with user info for rate limiting and audit
		c.Locals("user_id", keyRecord.UserID)
		c.Locals("api_key_id", keyRecord.ID)
		c.Locals("api_key_name", keyRecord.Name)
		c.Locals("api_key_scopes", keyRecord.Scopes)

		return c.Next()
	}
}

// OptionalAuthMiddleware extracts JWT if present but doesn't require it.
func OptionalAuthMiddleware(jwtManager *security.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Next()
		}

		claims, err := jwtManager.ValidateToken(parts[1])
		if err == nil {
			c.Locals("user_id", claims.UserID)
			c.Locals("user_email", claims.Email)
			c.Locals("user_role", claims.Role)
		}

		return c.Next()
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
