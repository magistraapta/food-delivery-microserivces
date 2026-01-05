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

// PaymentCheckoutCreatedEvent is published when a Stripe Checkout session is created
// This contains the URL the user should be redirected to for payment
type PaymentCheckoutCreatedEvent struct {
	OrderID     uuid.UUID `json:"order_id"`
	UserID      uuid.UUID `json:"user_id"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	CheckoutURL string    `json:"checkout_url"`
	SessionID   string    `json:"session_id"`
	ExpiresAt   time.Time `json:"expires_at"`
	Timestamp   time.Time `json:"timestamp"`
}
