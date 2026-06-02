package models

import "time"

// AuditLog captures security-relevant events for compliance.
type AuditLog struct {
	ID         string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Action     string    `gorm:"size:100;not null;index" json:"action"`
	Entity     string    `gorm:"size:100;not null" json:"entity"`
	EntityID   string    `gorm:"size:255" json:"entity_id,omitempty"`
	UserID     *string   `gorm:"type:uuid" json:"user_id,omitempty"`
	IPAddress  string    `gorm:"size:45" json:"ip_address,omitempty"`
	UserAgent  string    `gorm:"type:text" json:"user_agent,omitempty"`
	Metadata   string    `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt  time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
