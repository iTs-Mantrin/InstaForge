package health

import (
	"fmt"
	"sync"
	"time"

	"instaforge/internal/cache"
	"instaforge/internal/database"
	"instaforge/internal/dto"
	"instaforge/internal/queue"

	"github.com/gofiber/fiber/v2"
)

// Check represents a named health check.
type Check struct {
	Name   string
	Check  func() error
}

// Checker runs registered health checks.
type Checker struct {
	checks []Check
	mu     sync.RWMutex
}

// NewChecker creates a new health checker.
func NewChecker() *Checker {
	return &Checker{
		checks: make([]Check, 0),
	}
}

// Register adds a health check.
func (c *Checker) Register(name string, check func() error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.checks = append(c.checks, Check{Name: name, Check: check})
}

// LiveHandler returns a simple liveness probe (always 200 if process is running).
func (c *Checker) LiveHandler(ctx *fiber.Ctx) error {
	return ctx.JSON(dto.HealthResponse{
		Status:    "alive",
		Version:   "1.0.0",
		Timestamp: time.Now().UTC(),
	})
}

// ReadyHandler runs all readiness checks.
func (c *Checker) ReadyHandler(ctx *fiber.Ctx) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	results := make(map[string]string)
	status := "ready"

	for _, check := range c.checks {
		if err := check.Check(); err != nil {
			results[check.Name] = fmt.Sprintf("unhealthy: %v", err)
			status = "not_ready"
		} else {
			results[check.Name] = "healthy"
		}
	}

	httpStatus := fiber.StatusOK
	if status == "not_ready" {
		httpStatus = fiber.StatusServiceUnavailable
	}

	return ctx.Status(httpStatus).JSON(dto.HealthResponse{
		Status:    status,
		Version:   "1.0.0",
		Timestamp: time.Now().UTC(),
		Checks:    results,
	})
}

// DefaultChecks registers standard health checks for database, cache, and queue.
func (c *Checker) RegisterDefaults() {
	c.Register("database", func() error {
		return database.HealthCheck()
	})

	c.Register("cache", func() error {
		return cache.HealthCheck()
	})

	c.Register("queue", func() error {
		return queue.HealthCheck()
	})
}
