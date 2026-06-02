package models

import "time"

// User represents an API consumer (not an Instagram user).
type User struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Email     string    `gorm:"uniqueIndex;size:255;not null" json:"email"`
	Name      string    `gorm:"size:255" json:"name"`
	Role      string    `gorm:"size:50;default:user;not null" json:"role"`
	Status    string    `gorm:"size:50;default:active;not null" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
