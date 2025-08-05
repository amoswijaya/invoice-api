package db

import (
	"invoice-api/internal/model"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	dsn := os.Getenv("POSTGRES_DSN")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}
	DB.AutoMigrate(&model.User{}, &model.Invoice{}, &model.InvoiceItem{})
}