package models

import "time"

// AnalyticsEvent stores aggregated request analytics.
type AnalyticsEvent struct {
	ID         string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Endpoint   string    `gorm:"size:255;not null;index" json:"endpoint"`
	Method     string    `gorm:"size:10;not null" json:"method"`
	StatusCode int       `gorm:"not null" json:"status_code"`
	DurationMs int64     `gorm:"not null" json:"duration_ms"`
	IPAddress  string    `gorm:"size:45" json:"ip_address,omitempty"`
	UserAgent  string    `gorm:"type:text" json:"user_agent,omitempty"`
	UserID     *string   `gorm:"type:uuid" json:"user_id,omitempty"`
	Country    string    `gorm:"size:5" json:"country,omitempty"`
	CreatedAt  time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

func (AnalyticsEvent) TableName() string {
	return "analytics_events"
}

// AnalyticsSummary holds pre-aggregated metrics for dashboards.
type AnalyticsSummary struct {
	ID           string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Date         time.Time `gorm:"uniqueIndex:idx_date_endpoint;not null" json:"date"`
	Endpoint     string    `gorm:"uniqueIndex:idx_date_endpoint;size:255;not null" json:"endpoint"`
	TotalReq     int64     `gorm:"default:0" json:"total_requests"`
	SuccessReq   int64     `gorm:"default:0" json:"success_requests"`
	ErrorReq     int64     `gorm:"default:0" json:"error_requests"`
	AvgDurationMs float64  `gorm:"default:0" json:"avg_duration_ms"`
	P99DurationMs float64  `gorm:"default:0" json:"p99_duration_ms"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (AnalyticsSummary) TableName() string {
	return "analytics_summaries"
}

// DailyUsage tracks per-user daily usage for rate limiting.
type DailyUsage struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    string    `gorm:"type:uuid;uniqueIndex:idx_user_date;not null" json:"user_id"`
	Date      time.Time `gorm:"uniqueIndex:idx_user_date;not null" json:"date"`
	ReqCount  int64     `gorm:"default:0" json:"request_count"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (DailyUsage) TableName() string {
	return "daily_usage"
}
