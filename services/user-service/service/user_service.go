package service

import (
	"errors"
	"user-service/models"
	"user-service/repository"
	"user-service/utils"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return UserService{userRepository: userRepository}
}

func (us *UserService) CreateUser(user models.User) error {
	log.Info("Creating user")
	_, err := us.userRepository.GetUserByUsername(user.Username)

	if err == nil {
		log.Error("username already exists")
		return errors.New("username already exists")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error(err)
		return err
	}

	user.Password = utils.GeneratePassword(user.Password)
	user.Role = "user"

	log.Info("User created successfully")
	return us.userRepository.CreateUser(user)
}

func (us *UserService) GetUserByUsername(username string) (models.User, error) {
	log.Info("Getting user by username")
	return us.userRepository.GetUserByUsername(username)
}
