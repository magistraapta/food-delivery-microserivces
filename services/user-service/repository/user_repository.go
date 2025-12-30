package repository

import (
	"user-service/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user models.User) error
	GetUserByUsername(username string) (models.User, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (ur *UserRepositoryImpl) CreateUser(user models.User) error {

	return ur.db.Create(&user).Error
}

func (ur *UserRepositoryImpl) GetUserByUsername(username string) (models.User, error) {
	var user models.User

	return user, ur.db.Where("username = ?", username).First(&user).Error
}
