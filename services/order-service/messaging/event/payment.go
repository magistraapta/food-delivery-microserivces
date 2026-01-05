package event

import "time"

type PaymentSuccessEvent struct {
	OrderID               string    `json:"order_id"`
	UserID                string    `json:"user_id"`
	Amount                float64   `json:"amount"`
	Currency              string    `json:"currency"`
	StripePaymentIntentID string    `json:"stripe_payment_intent_id"`
	StripeChargeID        string    `json:"stripe_charge_id"`
	Timestamp             time.Time `json:"timestamp"`
}

type PaymentFailedEvent struct {
	OrderID       string    `json:"order_id"`
	CustomerID    string    `json:"customer_id"`
	PaymentID     string    `json:"payment_id"`
	FailureReason string    `json:"failure_reason"`
	FailureCode   string    `json:"failure_code"`
	Timestamp     time.Time `json:"timestamp"`
}
