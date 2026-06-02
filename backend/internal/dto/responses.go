package dto

import (
	"time"
	"instaforge/internal/models"
)

// APIResponse is the standard API response wrapper.
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *APIMeta    `json:"meta,omitempty"`
}

// APIError represents an error response.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// APIMeta holds pagination and metadata.
type APIMeta struct {
	RequestID string `json:"request_id,omitempty"`
	Cursor    string `json:"cursor,omitempty"`
	HasMore   bool   `json:"has_more"`
	Total     int64  `json:"total,omitempty"`
	TookMs    int64  `json:"took_ms"`
}

// PreviewResponse is the response for a URL preview.
type PreviewResponse struct {
	Type       string          `json:"type"`       // "post", "reel", "story"
	Username   string          `json:"username"`
	Caption    string          `json:"caption,omitempty"`
	MediaCount int             `json:"media_count"`
	MediaItems []MediaItemDTO  `json:"media_items"`
	Thumbnail  string          `json:"thumbnail,omitempty"`
	Likes      int             `json:"likes,omitempty"`
	Comments   int             `json:"comments,omitempty"`
	Timestamp  time.Time       `json:"timestamp"`
}

// MediaItemDTO is a serializable media item for API responses.
type MediaItemDTO struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	URL         string `json:"url"`
	Thumbnail   string `json:"thumbnail,omitempty"`
	Width       int    `json:"width,omitempty"`
	Height      int    `json:"height,omitempty"`
	FileSize    int64  `json:"file_size,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Duration    float64 `json:"duration,omitempty"`
}

// DownloadResponse is returned after a download request is queued.
type DownloadResponse struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
	MediaURL  string `json:"media_url,omitempty"`
	ExpiresIn int    `json:"expires_in,omitempty"` // seconds
}

// UsernameSearchResponse is the full user search result.
type UsernameSearchResponse struct {
	Profile ProfileDTO         `json:"profile"`
	Posts   []*models.InstagramPost `json:"posts"`
	Reels   []*models.InstagramPost `json:"reels"`
	Stories []*models.InstagramStory `json:"stories"`
}

// ProfileDTO is a serializable user profile.
type ProfileDTO struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	FullName      string `json:"full_name"`
	Biography     string `json:"biography"`
	ProfilePicURL string `json:"profile_pic_url"`
	FollowerCount int    `json:"follower_count"`
	FollowingCount int   `json:"following_count"`
	PostCount     int    `json:"post_count"`
	IsVerified    bool   `json:"is_verified"`
}

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Timestamp time.Time         `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}

// MetricsSnapshot is a summary of key metrics for dashboard.
type MetricsSnapshot struct {
	RequestsTotal   int64   `json:"requests_total"`
	RequestsPerSec  float64 `json:"requests_per_sec"`
	AvgResponseMs   float64 `json:"avg_response_ms"`
	P99ResponseMs   float64 `json:"p99_response_ms"`
	ActiveWorkers   int     `json:"active_workers"`
	QueueDepth      int64   `json:"queue_depth"`
	CacheHitRate    float64 `json:"cache_hit_rate"`
	ErrorRate       float64 `json:"error_rate"`
}

// ============================================================
// Frontend-compatible response types (camelCase JSON, matches frontend TS types)
// ============================================================

// PostPreviewResponse matches frontend InstagramPost type.
type PostPreviewResponse struct {
	ID             string            `json:"id"`
	Shortcode      string            `json:"shortcode"`
	MediaType      string            `json:"mediaType"`
	MediaItems     []MediaItemFEDTO  `json:"mediaItems"`
	Caption        string            `json:"caption,omitempty"`
	LikesCount     int               `json:"likesCount"`
	CommentsCount  int               `json:"commentsCount"`
	Timestamp      int64             `json:"timestamp"` // Unix milliseconds
	OwnerUsername  string            `json:"ownerUsername"`
	OwnerFullName  string            `json:"ownerFullName,omitempty"`
	OwnerProfilePic string           `json:"ownerProfilePic,omitempty"`
	DisplayUrl     string            `json:"displayUrl"`
	IsSponsored    bool              `json:"isSponsored,omitempty"`
}

// MediaItemFEDTO matches frontend InstagramMediaItem.
type MediaItemFEDTO struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	URL         string `json:"url"`
	Thumbnail   string `json:"thumbnail,omitempty"`
	Width       int    `json:"width,omitempty"`
	Height      int    `json:"height,omitempty"`
	HasAudio    bool   `json:"hasAudio,omitempty"`
	DownloadUrl string `json:"downloadUrl,omitempty"`
}

// UserInfoResponse matches frontend InstagramUser.
type UserInfoResponse struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	FullName      string `json:"fullName"`
	ProfilePicUrl string `json:"profilePicUrl"`
	FollowerCount int    `json:"followerCount"`
	FollowingCount int   `json:"followingCount"`
	PostCount     int    `json:"postCount"`
	IsVerified    bool   `json:"isVerified"`
	Biography     string `json:"biography,omitempty"`
	ExternalUrl   string `json:"externalUrl,omitempty"`
	IsPrivate     bool   `json:"isPrivate,omitempty"`
}

// FeedResponse matches frontend InstagramFeedResponse.
type FeedResponse struct {
	Items   []PostPreviewResponse `json:"items"`
	HasMore bool                  `json:"hasMore"`
	Cursor  string                `json:"cursor,omitempty"`
}

// DownloadJobResponse matches frontend DownloadJob.
type DownloadJobResponse struct {
	TaskId  string `json:"taskId"`
	Status  string `json:"status"`
	Source  string `json:"source"`
	Message string `json:"message"`
}

// ProgressResponse matches frontend ProgressData.
type ProgressResponse struct {
	TaskId      string `json:"taskId"`
	Percent     int    `json:"percent"`
	Speed       string `json:"speed"`
	Eta         string `json:"eta"`
	Filename    string `json:"filename"`
	Status      string `json:"status"`
	ErrorMsg    string `json:"errorMsg"`
	DownloadUrl string `json:"downloadUrl"`
}

// FileResponse matches frontend DownloadResult.
type FileResponse struct {
	DownloadUrl string `json:"downloadUrl"`
}
