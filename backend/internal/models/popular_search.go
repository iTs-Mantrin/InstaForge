package models

import "time"

// PopularSearch tracks frequently searched usernames for cache warming.
type PopularSearch struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Username    string    `gorm:"uniqueIndex;size:255;not null" json:"username"`
	SearchCount int64     `gorm:"default:1;not null" json:"search_count"`
	LastSearched time.Time `gorm:"autoCreateTime" json:"last_searched"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (PopularSearch) TableName() string {
	return "popular_searches"
}
