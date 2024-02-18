package model

import (
	"time"
)

type Customer struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"size:255;not null" json:"name"`
	Email          string    `gorm:"size:255;not null;unique" json:"email"`
	Address        string    `gorm:"size:255" json:"address"`
	HashCreditCard string    `gorm:"size:255;not null" json:"hash_credit_card"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
