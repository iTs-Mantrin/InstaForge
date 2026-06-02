package controllers

import (
	"instaforge/internal/health"
	pkgresponse "instaforge/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// HealthController handles health check endpoints.
type HealthController struct {
	checker *health.Checker
}

// NewHealthController creates a new health controller.
func NewHealthController(checker *health.Checker) *HealthController {
	return &HealthController{
		checker: checker,
	}
}

// Liveness handles GET /health/live
func (ctrl *HealthController) Liveness(c *fiber.Ctx) error {
	return ctrl.checker.LiveHandler(c)
}

// Readiness handles GET /health/ready
func (ctrl *HealthController) Readiness(c *fiber.Ctx) error {
	return ctrl.checker.ReadyHandler(c)
}

// Ping handles GET /api/ping
func (ctrl *HealthController) Ping(c *fiber.Ctx) error {
	return pkgresponse.Success(c, map[string]string{
		"message": "pong",
	})
}
