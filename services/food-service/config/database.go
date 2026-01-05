package config

import (
	"food-service/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadDatabase() (*gorm.DB, error) {
	dbUrl := os.Getenv("DB_URL")

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Restaurant{}, &models.Food{})

	return db, nil
}
