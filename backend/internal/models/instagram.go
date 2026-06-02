package models

import "time"

// MediaType represents the type of Instagram media.
type MediaType string

const (
	MediaTypeImage    MediaType = "image"
	MediaTypeVideo    MediaType = "video"
	MediaTypeCarousel MediaType = "carousel"
)

// MediaItem represents a single media item within an Instagram post.
type MediaItem struct {
	ID          string    `json:"id"`
	Type        MediaType `json:"type"`
	URL         string    `json:"url"`
	Thumbnail   string    `json:"thumbnail,omitempty"`
	Width       int       `json:"width,omitempty"`
	Height      int       `json:"height,omitempty"`
	FileSize    int64     `json:"file_size,omitempty"`
	ContentType string    `json:"content_type,omitempty"`
}

// InstagramPost represents an Instagram post (single or carousel).
type InstagramPost struct {
	ID         string      `json:"id"`
	Shortcode  string      `json:"shortcode"`
	Username   string      `json:"username"`
	Caption    string      `json:"caption,omitempty"`
	MediaType  MediaType   `json:"media_type"`
	MediaItems []MediaItem `json:"media_items"`
	Thumbnail  string      `json:"thumbnail,omitempty"`
	Likes      int         `json:"likes,omitempty"`
	Comments   int         `json:"comments,omitempty"`
	Timestamp  time.Time   `json:"timestamp"`
	IsVideo    bool        `json:"is_video"`
	Duration   float64     `json:"duration,omitempty"`
}

// InstagramStory represents an Instagram story.
type InstagramStory struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	MediaType  MediaType `json:"media_type"`
	URL        string    `json:"url"`
	Thumbnail  string    `json:"thumbnail,omitempty"`
	Duration   float64   `json:"duration,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
	ExpiresAt  time.Time `json:"expires_at,omitempty"`
}

// InstagramUser represents a public Instagram user profile.
type InstagramUser struct {
	ID            string    `json:"id"`
	Username      string    `json:"username"`
	FullName      string    `json:"full_name,omitempty"`
	Biography     string    `json:"biography,omitempty"`
	ProfilePicURL string    `json:"profile_pic_url,omitempty"`
	FollowerCount int       `json:"follower_count"`
	FollowingCount int      `json:"following_count"`
	PostCount     int       `json:"post_count"`
	IsPrivate     bool      `json:"is_private"`
	IsVerified    bool      `json:"is_verified"`
	ExternalURL   string    `json:"external_url,omitempty"`
}

// InstagramFeedResponse wraps a paginated Instagram feed.
type InstagramFeedResponse struct {
	Items      []InstagramPost `json:"items"`
	HasMore    bool            `json:"has_more"`
	Cursor     string          `json:"cursor,omitempty"`
	User       *InstagramUser  `json:"user,omitempty"`
}

// UsernameSearchResult is the aggregate response for username searches.
type UsernameSearchResult struct {
	Profile  InstagramUser    `json:"profile"`
	Posts    []InstagramPost  `json:"posts,omitempty"`
	Reels    []InstagramPost  `json:"reels,omitempty"`
	Stories  []InstagramStory `json:"stories,omitempty"`
}

// InstagramURLInfo represents info extracted from an Instagram URL.
type InstagramURLInfo struct {
	Type      string `json:"type"`       // "post", "reel", "story"
	Shortcode string `json:"shortcode,omitempty"`
	Username  string `json:"username,omitempty"`
	StoryID   string `json:"story_id,omitempty"`
	MediaID   string `json:"media_id,omitempty"`
}
