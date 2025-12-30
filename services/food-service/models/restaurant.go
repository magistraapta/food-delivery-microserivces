package models

import (
	"time"

	"github.com/google/uuid"
)

type Restaurant struct {
	ID        uuid.UUID `gorm:"type:uuid;primarykey"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Address   string    `gorm:"type:text;not null"`
	Foods     []Food    `gorm:"foreignKey:RestaurantID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
