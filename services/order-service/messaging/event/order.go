package event

import "time"

type OrderCreatedEvent struct {
	OrderID         string  `json:"order_id"`
	UserID          string  `json:"user_id"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	PaymentMethodID string  `json:"payment_method_id"`
}

// PaymentTimeoutEvent is published after a delay to check if payment was completed
type PaymentTimeoutEvent struct {
	OrderID   string    `json:"order_id"`
	CreatedAt time.Time `json:"created_at"`
}
