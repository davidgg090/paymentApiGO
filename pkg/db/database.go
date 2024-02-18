package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"paymentApi/internal/model"
)

var DB *gorm.DB

func InitDB(dsn string) {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("failed to connect database")
	}

	err = DB.AutoMigrate(&model.Customer{}, &model.Merchant{}, &model.Transaction{}, &model.User{}, &model.Token{}, &model.AuditLog{})
	if err != nil {
		log.Panic("failed to migrate database:", err)
	}
}
