package service

import (
	"errors"
	"user-service/models"
	"user-service/repository"
	"user-service/utils"

	"gorm.io/gorm"
)

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return UserService{userRepository: userRepository}
}

func (us *UserService) CreateUser(user models.User) error {
	_, err := us.userRepository.GetUserByUsername(user.Username)

	if err == nil {
		return errors.New("username already exists")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	user.Password = utils.GeneratePassword(user.Password)
	user.Role = "user"

	return us.userRepository.CreateUser(user)
}

func (us *UserService) GetUserByUsername(username string) (models.User, error) {
	return us.userRepository.GetUserByUsername(username)
}
