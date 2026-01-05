package dto

import (
	"time"

	"github.com/google/uuid"
)

type OrderItemRequest struct {
	FoodID   uuid.UUID `json:"food_id"`
	Quantity int       `json:"quantity"`
}

type CreateOrderRequest struct {
	UserID     uuid.UUID          `json:"user_id"`
	OrderItems []OrderItemRequest `json:"order_items"`
}

type OrderItemResponse struct {
	ID       uuid.UUID `json:"id"`
	FoodID   uuid.UUID `json:"food_id"`
	Quantity int       `json:"quantity"`
	Price    float64   `json:"price"`
}

type OrderResponse struct {
	ID          uuid.UUID           `json:"id"`
	UserID      uuid.UUID           `json:"user_id"`
	Status      string              `json:"status"`
	TotalAmount float64             `json:"total_amount"`
	OrderItems  []OrderItemResponse `json:"order_items"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}
