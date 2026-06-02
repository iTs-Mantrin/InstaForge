package api

import (
	"instaforge/internal/api/controllers"
	"instaforge/internal/health"
	"instaforge/internal/metrics"
	"instaforge/internal/middleware"
	"instaforge/internal/security"

	"github.com/gofiber/fiber/v2"
)

// RouterConfig holds dependencies for route setup.
type RouterConfig struct {
	InstagramCtrl *controllers.InstagramController
	DownloadCtrl  *controllers.DownloadController
	HealthCtrl    *controllers.HealthController
	JWTManager    *security.JWTManager
	RateLimiter   *middleware.RateLimiter
	Checker       *health.Checker
}

// SetupRoutes configures all application routes.
func SetupRoutes(app *fiber.App, cfg RouterConfig) {
	api := app.Group("/api", middleware.RequestID(), middleware.RequestLogger())

	// Health & Metrics (no auth)
	app.Get("/health/live", middleware.RequestID(), cfg.HealthCtrl.Liveness)
	app.Get("/health/ready", middleware.RequestID(), cfg.HealthCtrl.Readiness)

	// Frontend health probe (HEAD for Socket.IO pre-flight)
	app.Head("/api/health", cfg.HealthCtrl.Liveness)

	// Public routes (rate limited)
	v1 := api.Group("/v1", middleware.SecurityHeaders(), middleware.CORSConfig(), cfg.RateLimiter.Limit())

	// Ping
	v1.Get("/ping", cfg.HealthCtrl.Ping)

	// Instagram preview (no auth needed)
	v1.Post("/preview", cfg.InstagramCtrl.PreviewURL)

	// User search (no auth needed)
	v1.Get("/user/:username", cfg.InstagramCtrl.SearchUser)
	v1.Get("/user/:username/feed", cfg.InstagramCtrl.GetUserFeed)
	v1.Get("/user/:username/stories", cfg.InstagramCtrl.GetUserStories)

	// Download (no auth needed)
	v1.Post("/download", cfg.DownloadCtrl.QueueDownload)
	v1.Get("/download/:id", cfg.DownloadCtrl.GetDownloadStatus)
	v1.Get("/download/stream", cfg.DownloadCtrl.StreamMedia)

	// Authenticated routes (JWT required)
	auth := v1.Group("/auth", middleware.JWTAuthMiddleware(cfg.JWTManager))
	auth.Get("/profile", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"user_id": c.Locals("user_id"),
			"email":   c.Locals("user_email"),
			"role":    c.Locals("user_role"),
		})
	})

	// ============================================================
	// Frontend-compatible routes (/instagram/ prefix, camelCase JSON)
	// Matches the frontend's api.ts endpoint patterns.
	// ============================================================
	ig := app.Group("/instagram", middleware.RequestID(), middleware.RequestLogger(), middleware.SecurityHeaders(), middleware.CORSConfig(), cfg.RateLimiter.Limit())

	ig.Post("/preview", cfg.InstagramCtrl.PreviewFe)
	ig.Get("/user/:username/info", cfg.InstagramCtrl.UserInfo)
	ig.Get("/user/:username/feed", cfg.InstagramCtrl.UserFeedFe)
	ig.Get("/user/:username/stories", cfg.InstagramCtrl.StoriesFe)
	ig.Post("/download", cfg.DownloadCtrl.QueueDownloadFe)
	ig.Get("/progress/:taskId", cfg.DownloadCtrl.GetProgressResult)
	ig.Get("/file/:taskId", cfg.DownloadCtrl.GetFileResult)
	ig.Delete("/:taskId", cfg.DownloadCtrl.CancelDownload)

	// Prometheus metrics (mounted on a separate path for scraping)
	app.Get("/metrics", middleware.RequestID(), metrics.Handler())

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"code":    "NOT_FOUND",
				"message": "endpoint not found",
			},
		})
	})
}
