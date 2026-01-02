package models

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusSuccess PaymentStatus = "success"
	PaymentStatusFailed  PaymentStatus = "failed"
)

type Payment struct {
	ID                    uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	OrderID               uuid.UUID     `gorm:"type:uuid;not null;index" json:"order_id"`
	UserID                uuid.UUID     `gorm:"type:uuid;not null" json:"user_id"`
	Amount                float64       `gorm:"type:decimal(10,2);not null" json:"amount"`
	Currency              string        `gorm:"type:varchar(3);default:'usd'" json:"currency"`
	Status                PaymentStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	StripePaymentIntentID string        `gorm:"type:varchar(255)" json:"stripe_payment_intent_id"`
	StripeChargeID        string        `gorm:"type:varchar(255)" json:"stripe_charge_id"`
	FailureReason         string        `gorm:"type:text" json:"failure_reason,omitempty"`
	CreatedAt             time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}
