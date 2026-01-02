package config

import (
	"os"
	"payment-service/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {
	dbUrl := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Payment{})

	return db, nil
}
