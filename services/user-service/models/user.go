package models

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000" gorm:"primaryKey"`
	Username string    `json:"username" example:"user1"`
	Password string    `json:"password" example:"password123"`
	Role     string    `json:"role" example:"user"`
}

type CreateUser struct {
	Username string `json:"username" example:"user1"`
	Password string `json:"password" example:"password123"`
}

type UserLogin struct {
	Username string `json:"username" example:"user1"`
	Password string `json:"password" example:"password123"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

type AddUser struct {
	Username string `json:"username" example:"user1"`
	Password string `json:"password" example:"password123"`
}

func (a AddUser) ValidateUser() error {
	if a.Username == "" {
		return errors.New("username is required")
	}
	if a.Password == "" {
		return errors.New("password is required")
	}
	return nil
}
