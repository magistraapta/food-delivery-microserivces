package service

import (
	"payment-service/messaging/events"
	"payment-service/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	stripe "github.com/stripe/stripe-go/v84"
)

// ----- Mock Repository -----
type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) CreatePayment(p *models.Payment) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockPaymentRepository) FindByOrderId(orderId uuid.UUID) (*models.Payment, error) {
	args := m.Called(orderId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Payment), args.Error(1)
}

func (m *MockPaymentRepository) UpdateStatus(orderId uuid.UUID, status models.PaymentStatus) error {
	args := m.Called(orderId, status)
	return args.Error(0)
}

func (m *MockPaymentRepository) UpdateCheckoutSession(orderId uuid.UUID, sessionId string, checkoutURL string) error {
	args := m.Called(orderId, sessionId, checkoutURL)
	return args.Error(0)
}

func (m *MockPaymentRepository) UpdatePaymentIntent(orderId uuid.UUID, paymentIntentId string, chargeId string) error {
	args := m.Called(orderId, paymentIntentId, chargeId)
	return args.Error(0)
}

func (m *MockPaymentRepository) FindByCheckoutSessionId(sessionId string) (*models.Payment, error) {
	args := m.Called(sessionId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Payment), args.Error(1)
}

func (m *MockPaymentRepository) UpdateIntentId(orderId uuid.UUID, intentId string) error {
	args := m.Called(orderId, intentId)
	return args.Error(0)
}

// ----- Mock Stripe Client -----
type MockStripeClient struct {
	mock.Mock
}

func (m *MockStripeClient) CreateCheckoutSession(orderID string, amount float64, currency string, productName string) (*stripe.CheckoutSession, error) {
	args := m.Called(orderID, amount, currency, productName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.CheckoutSession), args.Error(1)
}

func (m *MockStripeClient) GetCheckoutSession(sessionID string) (*stripe.CheckoutSession, error) {
	args := m.Called(sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.CheckoutSession), args.Error(1)
}

func (m *MockStripeClient) CreatePaymentIntent(orderID string, amount float64, currency string) (*stripe.PaymentIntent, error) {
	args := m.Called(orderID, amount, currency)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.PaymentIntent), args.Error(1)
}

func (m *MockStripeClient) ConfirmPaymentIntent(paymentIntentID string, paymentMethodID string) (*stripe.PaymentIntent, error) {
	args := m.Called(paymentIntentID, paymentMethodID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.PaymentIntent), args.Error(1)
}

func (m *MockStripeClient) GetPaymentIntent(paymentIntentID string) (*stripe.PaymentIntent, error) {
	args := m.Called(paymentIntentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.PaymentIntent), args.Error(1)
}

func (m *MockStripeClient) CreateAndConfirmPaymentIntent(orderID string, amount float64, currency string, paymentMethodID string) (*stripe.PaymentIntent, error) {
	args := m.Called(orderID, amount, currency, paymentMethodID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.PaymentIntent), args.Error(1)
}

// ----- Mock RabbitMQ Client -----
type MockRabbitMQClient struct {
	mock.Mock
}

func (m *MockRabbitMQClient) PublishPaymentSuccess(event events.PaymentSuccessEvent) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockRabbitMQClient) PublishPaymentFailed(event events.PaymentFailedEvent) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockRabbitMQClient) PublishPaymentCheckoutCreated(event events.PaymentCheckoutCreatedEvent) error {
	args := m.Called(event)
	return args.Error(0)
}

func TestProcessOrderCreatedEvent_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockPaymentRepository)
	mockStripe := new(MockStripeClient)
	mockRabbitMQ := new(MockRabbitMQClient)

	service := NewPaymentService(mockRepo, mockStripe, mockRabbitMQ)

	orderID := uuid.New()
	userID := uuid.New()

	event := events.OrderCreatedEvent{
		OrderID:  orderID,
		UserID:   userID,
		Amount:   49.99,
		Currency: "usd",
	}

	expectedSession := &stripe.CheckoutSession{
		ID:  "cs_test_123",
		URL: "https://checkout.stripe.com/pay/cs_test_123",
	}

	// Set up mock expectations
	mockRepo.On("CreatePayment", mock.AnythingOfType("*models.Payment")).Return(nil)
	mockStripe.On("CreateCheckoutSession", orderID.String(), 49.99, "usd", "Food Order").Return(expectedSession, nil)
	mockRepo.On("UpdateCheckoutSession", orderID, "cs_test_123", "https://checkout.stripe.com/pay/cs_test_123").Return(nil)
	mockRabbitMQ.On("PublishPaymentCheckoutCreated", mock.AnythingOfType("events.PaymentCheckoutCreatedEvent")).Return(nil)

	// Act
	err := service.ProcessOrderCreatedEvent(event)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockStripe.AssertExpectations(t)
	mockRabbitMQ.AssertExpectations(t)
}
