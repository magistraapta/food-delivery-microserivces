package config

import (
	"food-service/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadDatabase() (*gorm.DB, error) {
	dbUrl := "postgres://postgres:postgres@localhost:5432/food-service-db?sslmode=disable"

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Restaurant{}, &models.Food{})

	return db, nil
}
