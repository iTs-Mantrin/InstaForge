package models

import "time"

// APIKey represents an API key for programmatic access.
type APIKey struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Key       string    `gorm:"uniqueIndex;size:255;not null" json:"-"`
	Name      string    `gorm:"size:255" json:"name"`
	Scopes    string    `gorm:"type:text" json:"scopes,omitempty"`
	LastUsed  *time.Time `json:"last_used,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	Status    string    `gorm:"size:50;default:active;not null" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
}

func (APIKey) TableName() string {
	return "api_keys"
}
