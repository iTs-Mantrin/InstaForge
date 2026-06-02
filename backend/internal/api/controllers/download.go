package controllers

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"instaforge/internal/dto"
	"instaforge/internal/metrics"
	"instaforge/internal/models"
	"instaforge/internal/services"
	"instaforge/internal/validators"
	pkgresponse "instaforge/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// DownloadController handles media download endpoints.
type DownloadController struct {
	downloadService *services.DownloadService
}

// NewDownloadController creates a new download controller.
func NewDownloadController(downloadService *services.DownloadService) *DownloadController {
	return &DownloadController{
		downloadService: downloadService,
	}
}

// QueueDownload handles POST /api/download (queue a download).
func (ctrl *DownloadController) QueueDownload(c *fiber.Ctx) error {
	var req dto.DownloadRequest
	if err := c.BodyParser(&req); err != nil {
		return pkgresponse.BadRequest(c, "invalid request body")
	}

	if req.URL == "" {
		return pkgresponse.ValidationError(c, "url is required")
	}

	// Validate URL
	_, _, err := validators.ValidateInstagramURL(req.URL)
	if err != nil {
		return pkgresponse.ValidationError(c, "invalid Instagram URL: "+err.Error())
	}

	// Extract user identity from context if authenticated
	var userID *string
	if uid, ok := c.Locals("user_id").(string); ok && uid != "" {
		userID = &uid
	}

	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")

	downloadReq, err := ctrl.downloadService.QueueDownload(req.URL, ipAddress, userAgent, userID)
	if err != nil {
		metrics.RecordDownload("unknown", false)
		return pkgresponse.InternalError(c, "failed to queue download, please try again later")
	}

	metrics.RecordDownload(downloadReq.URLType, true)

	return pkgresponse.Accepted(c, dto.DownloadResponse{
		RequestID: downloadReq.ID,
		Status:    string(downloadReq.Status),
	})
}

// GetDownloadStatus handles GET /api/download/:id (check download status).
func (ctrl *DownloadController) GetDownloadStatus(c *fiber.Ctx) error {
	requestID := c.Params("id")
	if requestID == "" {
		return pkgresponse.ValidationError(c, "download request ID is required")
	}

	result, err := ctrl.downloadService.GetDownloadResult(requestID)
	if err != nil {
		return pkgresponse.NotFound(c, "download result not found or expired")
	}

	return pkgresponse.Success(c, result)
}

// StreamMedia handles GET /api/download/:id/stream (stream media file).
func (ctrl *DownloadController) StreamMedia(c *fiber.Ctx) error {
	mediaURL := c.Query("url")
	if mediaURL == "" {
		return pkgresponse.ValidationError(c, "media URL parameter is required")
	}

	// Validate the media URL is from an allowed domain
	if !strings.HasPrefix(mediaURL, "https://") {
		return pkgresponse.BadRequest(c, "invalid media URL")
	}

	reader, contentType, contentLength, err := ctrl.downloadService.StreamMedia(mediaURL)
	if err != nil {
		return pkgresponse.InternalError(c, "failed to stream media, please try again later")
	}
	defer reader.Close()

	// Set headers for download
	ext := filepath.Ext(mediaURL)
	if ext == "" {
		ext = ".mp4"
	}
	c.Response().Header.Set("Content-Type", contentType)
	c.Response().Header.Set("Content-Disposition", "attachment; filename=\"instaforge_download"+ext+"\"")
	if contentLength > 0 {
		c.Response().Header.Set("Content-Length", fmt.Sprintf("%d", contentLength))
	}

	// Stream the file
	_, err = io.Copy(c.Response().BodyWriter(), reader)
	if err != nil {
		return pkgresponse.InternalError(c, "failed to stream media response")
	}

	return nil
}

// ============================================================
// Frontend-compatible handlers (camelCase JSON, /instagram/ prefix)
// ============================================================

// statusToFrontend maps backend status to frontend-compatible status string.
func statusToFrontend(status models.DownloadStatus) string {
	switch status {
	case models.DownloadStatusPending:
		return "queued"
	case models.DownloadStatusProcessing:
		return "downloading"
	case models.DownloadStatusCompleted:
		return "done"
	case models.DownloadStatusFailed:
		return "failed"
	case models.DownloadStatusCancelled:
		return "cancelled"
	default:
		return "queued"
	}
}

// QueueDownloadFe handles POST /instagram/download (frontend-facing, returns taskId).
func (ctrl *DownloadController) QueueDownloadFe(c *fiber.Ctx) error {
	var req dto.DownloadRequest
	if err := c.BodyParser(&req); err != nil {
		return pkgresponse.BadRequest(c, "invalid request body")
	}

	if req.URL == "" {
		return pkgresponse.ValidationError(c, "url is required")
	}

	_, _, err := validators.ValidateInstagramURL(req.URL)
	if err != nil {
		return pkgresponse.ValidationError(c, "invalid Instagram URL: "+err.Error())
	}

	var userID *string
	if uid, ok := c.Locals("user_id").(string); ok && uid != "" {
		userID = &uid
	}

	downloadReq, err := ctrl.downloadService.QueueDownload(req.URL, c.IP(), c.Get("User-Agent"), userID)
	if err != nil {
		metrics.RecordDownload("unknown", false)
		return pkgresponse.InternalError(c, "failed to queue download, please try again later")
	}

	metrics.RecordDownload(downloadReq.URLType, true)

	return pkgresponse.Accepted(c, dto.DownloadJobResponse{
		TaskId:  downloadReq.ID,
		Status:  statusToFrontend(downloadReq.Status),
		Source:  req.URL,
	})
}

// GetProgressResult handles GET /instagram/progress/:taskId (frontend-facing).
func (ctrl *DownloadController) GetProgressResult(c *fiber.Ctx) error {
	taskID := c.Params("taskId")
	if taskID == "" {
		return pkgresponse.ValidationError(c, "task id is required")
	}

	req, err := ctrl.downloadService.GetDownloadProgress(taskID)
	if err != nil {
		return pkgresponse.NotFound(c, "download not found")
	}

	percent := 0
	switch req.Status {
	case models.DownloadStatusProcessing:
		percent = 50
	case models.DownloadStatusCompleted:
		percent = 100
	case models.DownloadStatusFailed, models.DownloadStatusCancelled:
		percent = 0
	default:
		percent = 0
	}

	errorMsg := ""
	if req.Status == models.DownloadStatusFailed && req.ErrorMessage != "" {
		errorMsg = req.ErrorMessage
	}

	downloadURL := ""
	if req.Status == models.DownloadStatusCompleted {
		imgURL, err := ctrl.downloadService.GetDownloadFileResult(taskID)
		if err == nil {
			downloadURL = imgURL
		}
	}

	return pkgresponse.Success(c, dto.ProgressResponse{
		TaskId:      taskID,
		Percent:     percent,
		Speed:       "N/A",
		Eta:         "N/A",
		Filename:    "",
		Status:      statusToFrontend(req.Status),
		ErrorMsg:    errorMsg,
		DownloadUrl: downloadURL,
	})
}

// GetFileResult handles GET /instagram/file/:taskId (frontend-facing).
func (ctrl *DownloadController) GetFileResult(c *fiber.Ctx) error {
	taskID := c.Params("taskId")
	if taskID == "" {
		return pkgresponse.ValidationError(c, "task id is required")
	}

	downloadURL, err := ctrl.downloadService.GetDownloadFileResult(taskID)
	if err != nil {
		return pkgresponse.NotFound(c, "download result not found or expired")
	}

	return pkgresponse.Success(c, dto.FileResponse{
		DownloadUrl: downloadURL,
	})
}

// CancelDownload handles DELETE /instagram/:taskId (frontend-facing).
func (ctrl *DownloadController) CancelDownload(c *fiber.Ctx) error {
	taskID := c.Params("taskId")
	if taskID == "" {
		return pkgresponse.ValidationError(c, "task id is required")
	}

	if err := ctrl.downloadService.CancelJob(taskID); err != nil {
		return pkgresponse.BadRequest(c, err.Error())
	}

	return pkgresponse.Success(c, fiber.Map{
		"message": "download cancelled",
	})
}
