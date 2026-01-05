package models

import "github.com/google/uuid"

type OrderItem struct {
	ID       uuid.UUID `gorm:"type:uuid;primarykey"`
	OrderID  uuid.UUID `gorm:"type:uuid;not null"` // Foreign key
	FoodID   uuid.UUID `gorm:"type:uuid;not null"`
	Quantity int       `gorm:"type:int;not null"`
	Price    float64   `gorm:"type:decimal(10,2);not null"` // Price at order time
}
