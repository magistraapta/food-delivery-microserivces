package models

import (
	"time"

	"github.com/google/uuid"
)

type Food struct {
	ID           uuid.UUID `gorm:"type:uuid;primarykey"`
	RestaurantID uuid.UUID `gorm:"type:uuid;not null"`
	Name         string    `gorm:"type:varchar(255);not null"`
	Price        float64   `gorm:"type:int;not null"`
	Description  string    `gorm:"type:text"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}
