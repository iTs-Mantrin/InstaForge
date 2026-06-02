package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORSConfig returns a CORS middleware configuration.
func CORSConfig() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS, PATCH",
		AllowHeaders: "Accept, Authorization, Content-Type, X-API-Key, X-Request-ID",
		ExposeHeaders: "X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset, Retry-After, X-Request-ID",
		MaxAge:       3600,
	})
}
