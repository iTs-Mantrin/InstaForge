package models

import "time"

// DownloadStatus represents the state of a download request.
type DownloadStatus string

const (
	DownloadStatusPending    DownloadStatus = "pending"
	DownloadStatusProcessing DownloadStatus = "processing"
	DownloadStatusCompleted  DownloadStatus = "completed"
	DownloadStatusFailed     DownloadStatus = "failed"
	DownloadStatusExpired    DownloadStatus = "expired"
	DownloadStatusCancelled  DownloadStatus = "cancelled"
)

// DownloadRequest tracks each download initiated through the platform.
type DownloadRequest struct {
	ID            string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID        *string        `gorm:"type:uuid;index" json:"user_id,omitempty"`
	URL           string         `gorm:"type:text;not null" json:"url"`
	URLType       string         `gorm:"size:20;not null" json:"url_type"` // post, reel, story
	Shortcode     string         `gorm:"size:255;index" json:"shortcode,omitempty"`
	Status        DownloadStatus `gorm:"size:20;default:pending;not null;index" json:"status"`
	MediaCount    int            `gorm:"default:0" json:"media_count"`
	FileSize      int64          `gorm:"default:0" json:"file_size"`
	ErrorMessage  string         `gorm:"type:text" json:"error_message,omitempty"`
	IPAddress     string         `gorm:"size:45" json:"ip_address,omitempty"`
	UserAgent     string         `gorm:"type:text" json:"user_agent,omitempty"`
	ProcessedAt   *time.Time     `json:"processed_at,omitempty"`
	CompletedAt   *time.Time     `json:"completed_at,omitempty"`
	ExpiresAt     *time.Time     `json:"expires_at,omitempty"`
	CreatedAt     time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

func (DownloadRequest) TableName() string {
	return "download_requests"
}
