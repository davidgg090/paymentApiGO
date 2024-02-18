package model

import (
	"time"
)

type Merchant struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Name              string    `gorm:"size:255;not null" json:"name"`
	Email             string    `gorm:"size:255;not null;unique" json:"email"`
	IsActive          bool      `gorm:"default:true" json:"is_active"`
	AmountAccount     float64   `gorm:"default:0;not null" json:"amount_account"`
	CreatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	AuthenticationKey string    `gorm:"not null" json:"authentication_key"`
}
