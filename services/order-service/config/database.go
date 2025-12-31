package config

import (
	"order-service/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {
	dbUrl := "postgres://postgres:postgres@localhost:5432/order-service-db"

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Order{}, &models.OrderItem{})

	return db, nil
}
