package dto

// PreviewRequest is the request body for URL preview.
type PreviewRequest struct {
	URL string `json:"url" validate:"required,instagram_url"`
}

// PreviewBulkRequest handles batch URL preview.
type PreviewBulkRequest struct {
	URLs []string `json:"urls" validate:"required,min=1,max=10,dive,instagram_url"`
}

// UsernameSearchParams holds query parameters for username search.
type UsernameSearchParams struct {
	Username string `query:"username" validate:"required,instagram_username"`
	Cursor   string `query:"cursor,omitempty"`
	Limit    int    `query:"limit,omitempty" validate:"omitempty,min=1,max=50"`
}

// UserFeedParams holds pagination parameters for user feed.
type UserFeedParams struct {
	Username string `query:"username" validate:"required,instagram_username"`
	Cursor   string `query:"cursor,omitempty"`
	Limit    int    `query:"limit,omitempty" validate:"omitempty,min=1,max=50"`
	MediaType string `query:"media_type,omitempty" validate:"omitempty,oneof=post reel story"`
}

// DownloadRequest is the body for initiating a download.
type DownloadRequest struct {
	URL      string `json:"url" validate:"required,instagram_url"`
	MediaID  string `json:"media_id,omitempty"`
	Quality  string `json:"quality,omitempty" validate:"omitempty,oneof=best medium low"`
}

// APIKeyCreateRequest is the body for creating a new API key.
type APIKeyCreateRequest struct {
	Name   string   `json:"name" validate:"required,min=3,max=100"`
	Scopes []string `json:"scopes,omitempty"`
}

// PaginationParams are common pagination query params.
type PaginationParams struct {
	Cursor string `query:"cursor,omitempty"`
	Limit  int    `query:"limit,omitempty" validate:"omitempty,min=1,max=100"`
}
