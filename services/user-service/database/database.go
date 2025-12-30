package database

import (
	"user-service/common"
	"user-service/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(config *common.Config) (*gorm.DB, error) {
	dbUrl := config.DatabaseUrl

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.User{})

	return db, nil
}
