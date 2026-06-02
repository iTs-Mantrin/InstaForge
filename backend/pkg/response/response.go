package response

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// APIResponse is the standard API response envelope.
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// APIError represents an error in the response.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Meta holds pagination and request metadata.
type Meta struct {
	RequestID string `json:"request_id,omitempty"`
	Cursor    string `json:"cursor,omitempty"`
	HasMore   bool   `json:"has_more,omitempty"`
	Total     int64  `json:"total,omitempty"`
	TookMs    int64  `json:"took_ms,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

// Standard error codes.
const (
	ErrBadRequest       = "BAD_REQUEST"
	ErrUnauthorized     = "UNAUTHORIZED"
	ErrForbidden        = "FORBIDDEN"
	ErrNotFound         = "NOT_FOUND"
	ErrRateLimited      = "RATE_LIMITED"
	ErrConflict         = "CONFLICT"
	ErrValidation       = "VALIDATION_ERROR"
	ErrInternal         = "INTERNAL_ERROR"
	ErrServiceUnhealthy = "SERVICE_UNAVAILABLE"
)

// Success sends a 200 OK response.
func Success(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(APIResponse{
		Success: true,
		Data:    data,
		Meta:    buildMeta(c),
	})
}

// SuccessWithMeta sends a success response with pagination metadata.
func SuccessWithMeta(c *fiber.Ctx, data interface{}, cursor string, hasMore bool, total int64) error {
	meta := buildMeta(c)
	meta.Cursor = cursor
	meta.HasMore = hasMore
	meta.Total = total
	return c.JSON(APIResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// Created sends a 201 Created response.
func Created(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(APIResponse{
		Success: true,
		Data:    data,
		Meta:    buildMeta(c),
	})
}

// Accepted sends a 202 Accepted response (for async operations).
func Accepted(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusAccepted).JSON(APIResponse{
		Success: true,
		Data:    data,
		Meta:    buildMeta(c),
	})
}

// NoContent sends a 204 No Content response.
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// Error sends an error response with the given status code.
func Error(c *fiber.Ctx, status int, code, message string) error {
	return c.Status(status).JSON(APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
		Meta: buildMeta(c),
	})
}

// ErrorWithDetails sends an error response with additional details.
func ErrorWithDetails(c *fiber.Ctx, status int, code, message, details string) error {
	return c.Status(status).JSON(APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
		Meta: buildMeta(c),
	})
}

// BadRequest sends a 400 response.
func BadRequest(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusBadRequest, ErrBadRequest, message)
}

// ValidationError sends a 400 with validation error code.
func ValidationError(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusBadRequest, ErrValidation, message)
}

// Unauthorized sends a 401 response.
func Unauthorized(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnauthorized, ErrUnauthorized, message)
}

// Forbidden sends a 403 response.
func Forbidden(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusForbidden, ErrForbidden, message)
}

// NotFound sends a 404 response.
func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, ErrNotFound, message)
}

// RateLimited sends a 429 response.
func RateLimited(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusTooManyRequests, ErrRateLimited, message)
}

// InternalError sends a 500 response.
func InternalError(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusInternalServerError, ErrInternal, message)
}

// ServiceUnavailable sends a 503 response.
func ServiceUnavailable(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusServiceUnavailable, ErrServiceUnhealthy, message)
}

func buildMeta(c *fiber.Ctx) *Meta {
	requestID, _ := c.Locals("request_id").(string)
	tookMs, _ := c.Locals("took_ms").(int64)
	return &Meta{
		RequestID: requestID,
		TookMs:    tookMs,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}
