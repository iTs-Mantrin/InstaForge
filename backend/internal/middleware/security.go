package middleware

import (
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// SecurityHeaders adds security-related HTTP headers to all responses.
func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Response().Header.Set("X-Content-Type-Options", "nosniff")
		c.Response().Header.Set("X-Frame-Options", "DENY")
		c.Response().Header.Set("X-XSS-Protection", "1; mode=block")
		c.Response().Header.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Response().Header.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		// Content Security Policy
		c.Response().Header.Set("Content-Security-Policy",
			"default-src 'self'; "+
				"connect-src 'self' *.instagram.com *.cdninstagram.com; "+
				"img-src 'self' data: *.instagram.com *.cdninstagram.com *.picsum.photos; "+
				"style-src 'self' 'unsafe-inline'; "+
				"script-src 'self' 'unsafe-inline'; "+
				"font-src 'self' data:;",
		)

		// HSTS (only in production)
		if c.Hostname() != "localhost" && c.Hostname() != "127.0.0.1" {
			c.Response().Header.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		return c.Next()
	}
}

// InputSanitizer sanitizes user input to prevent injection attacks.
func InputSanitizer() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Sanitize query params and rebuild query string
		params := c.Queries()
		if len(params) == 0 {
			return c.Next()
		}
		var pairs []string
		for key, val := range params {
			pairs = append(pairs, key+"="+sanitizeInput(val))
		}
		c.Request().URI().SetQueryString(strings.Join(pairs, "&"))
		return c.Next()
	}
}

func sanitizeInput(input string) string {
	// Remove potentially dangerous characters
	re := regexp.MustCompile(`[<>"'\x00-\x1F\x7F]`)
	return re.ReplaceAllString(input, "")
}

// BodySizeLimit limits the maximum request body size.
func BodySizeLimit(maxBytes int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if len(c.Body()) > maxBytes {
			return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"code":    "REQUEST_TOO_LARGE",
					"message": "request body exceeds maximum allowed size",
				},
			})
		}
		return c.Next()
	}
}
