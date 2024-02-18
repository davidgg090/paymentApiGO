package model

import (
	"time"
)

type AuditLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       *uint     `gorm:"index;column:user_id" json:"user_id,omitempty"`
	ActivityType string    `gorm:"not null;column:activity_type" json:"activity_type"`
	BearerToken  *string   `gorm:"column:bearer_token" json:"bearer_token,omitempty"`
	IPAddress    *string   `gorm:"column:ip_address" json:"ip_address,omitempty"`
	Path         string    `gorm:"not null;column:path" json:"path"`
	Timestamp    time.Time `gorm:"not null;column:timestamp" json:"timestamp"`
}
