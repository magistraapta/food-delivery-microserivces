package events

import (
	"time"

	"github.com/google/uuid"
)

type OrderCreatedEvent struct {
	OrderID         uuid.UUID `json:"order_id"`
	UserID          uuid.UUID `json:"user_id"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	PaymentMethodID string    `json:"payment_method_id"`
	Timestamp       time.Time `json:"timestamp"`
}
