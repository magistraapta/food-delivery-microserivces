package repository

import (
	"payment-service/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	CreatePayment(p *models.Payment) error
	FindByOrderId(orderId uuid.UUID) (*models.Payment, error)
	FindByCheckoutSessionId(sessionId string) (*models.Payment, error)
	UpdateStatus(orderId uuid.UUID, status models.PaymentStatus) error
	UpdateIntentId(orderId uuid.UUID, intentId string) error
	UpdateCheckoutSession(orderId uuid.UUID, sessionId string, checkoutURL string) error
	UpdatePaymentIntent(orderId uuid.UUID, paymentIntentId string, chargeId string) error
}

type PaymentRepositoryImpl struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &PaymentRepositoryImpl{db: db}
}

func (r *PaymentRepositoryImpl) CreatePayment(p *models.Payment) error {
	return r.db.Create(p).Error
}

func (r *PaymentRepositoryImpl) FindByOrderId(orderId uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.Where("order_id = ?", orderId).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepositoryImpl) FindByCheckoutSessionId(sessionId string) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.Where("stripe_checkout_session_id = ?", sessionId).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepositoryImpl) UpdateStatus(orderId uuid.UUID, status models.PaymentStatus) error {
	return r.db.Model(&models.Payment{}).Where("order_id = ?", orderId).Update("status", status).Error
}

func (r *PaymentRepositoryImpl) UpdateIntentId(orderId uuid.UUID, intentId string) error {
	return r.db.Model(&models.Payment{}).Where("order_id = ?", orderId).Update("stripe_payment_intent_id", intentId).Error
}

func (r *PaymentRepositoryImpl) UpdateCheckoutSession(orderId uuid.UUID, sessionId string, checkoutURL string) error {
	return r.db.Model(&models.Payment{}).Where("order_id = ?", orderId).Updates(map[string]interface{}{
		"stripe_checkout_session_id": sessionId,
		"checkout_url":               checkoutURL,
	}).Error
}

func (r *PaymentRepositoryImpl) UpdatePaymentIntent(orderId uuid.UUID, paymentIntentId string, chargeId string) error {
	return r.db.Model(&models.Payment{}).Where("order_id = ?", orderId).Updates(map[string]interface{}{
		"stripe_payment_intent_id": paymentIntentId,
		"stripe_charge_id":         chargeId,
	}).Error
}
