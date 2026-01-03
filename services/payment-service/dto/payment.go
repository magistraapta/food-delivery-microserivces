package dto

import (
	"time"

	"github.com/google/uuid"
)

// CheckoutResponse is returned when requesting a checkout URL
type CheckoutResponse struct {
	OrderID     uuid.UUID `json:"order_id"`
	CheckoutURL string    `json:"checkout_url"`
	Status      string    `json:"status"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
}

// PaymentStatusResponse is returned when requesting payment status
type PaymentStatusResponse struct {
	OrderID   uuid.UUID `json:"order_id"`
	Status    string    `json:"status"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
