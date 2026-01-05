package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	PENDING          = "PENDING"
	CONFIRMED        = "CONFIRMED"
	PREPARING        = "PREPARING"
	READY_FOR_PICKUP = "READY_FOR_PICKUP"
	OUT_FOR_DELIVERY = "OUT_FOR_DELIVERY"
	DELIVERED        = "DELIVERED"
	CANCELLED        = "CANCELLED"      // Order cancelled (e.g., payment timeout)
	PAYMENT_FAILED   = "PAYMENT_FAILED" // Payment explicitly failed
)

type Order struct {
	ID          uuid.UUID   `gorm:"type:uuid;primarykey"`
	UserID      uuid.UUID   `gorm:"type:uuid;not null"`
	OrderItems  []OrderItem `gorm:"foreignKey:OrderID"`
	Status      string      `gorm:"type:varchar(50);default:'PENDING'"`
	TotalAmount float64     `gorm:"type:decimal(10,2);not null"`
	CreatedAt   time.Time   `gorm:"autoCreateTime"`
	UpdatedAt   time.Time   `gorm:"autoUpdateTime"`
}
