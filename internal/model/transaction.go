package model

import (
	"time"
)

type Transaction struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	MerchantID     uint      `gorm:"not null" json:"merchant_id"`
	CustomerID     uint      `gorm:"not null" json:"customer_id"`
	Amount         float64   `gorm:"type:numeric(10,2);not null" json:"amount"`
	Currency       string    `gorm:"type:char(3);not null" json:"currency"`
	State          string    `gorm:"size:255;default:pending;not null" json:"state"`
	HashCreditCard string    `gorm:"size:255;not null" json:"hash_credit_card"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Token          string    `gorm:"not null;unique" json:"token"`
}
