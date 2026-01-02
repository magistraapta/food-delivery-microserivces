package events

import (
	"time"

	"github.com/google/uuid"
)

type PaymentSuccessEvent struct {
	OrderID               uuid.UUID `json:"order_id"`
	UserID                uuid.UUID `json:"user_id"`
	Amount                float64   `json:"amount"`
	Currency              string    `json:"currency"`
	StripePaymentIntentID string    `json:"stripe_payment_intent_id"`
	StripeChargeID        string    `json:"stripe_charge_id"`
	Timestamp             time.Time `json:"timestamp"`
}

type PaymentFailedEvent struct {
	OrderID       uuid.UUID `json:"order_id"`
	CustomerID    uuid.UUID `json:"customer_id"`
	PaymentID     uuid.UUID `json:"payment_id"`
	FailureReason string    `json:"failure_reason"`
	FailureCode   string    `json:"failure_code"`
	Timestamp     time.Time `json:"timestamp"`
}
