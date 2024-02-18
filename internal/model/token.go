package model

import (
	"time"
)

type Token struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	AccessToken string    `gorm:"size:255;not null;unique" json:"access_token"`
	TokenType   string    `gorm:"size:255;not null" json:"token_type"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	UserID      *uint     `gorm:"index" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"-"`
	ExpiresAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"expires_at"`
}
