package service

import (
	"order-service/messaging/event"
	"order-service/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) CreateOrder(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetOrderById(id uuid.UUID) (*models.Order, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) UpdateOrderStatus(id uuid.UUID, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

type MockRabbitMQClient struct {
	mock.Mock
}

func (m *MockRabbitMQClient) PublishOrderCreated(evt event.OrderCreatedEvent) error {
	args := m.Called(evt)
	return args.Error(0)
}

func (m *MockRabbitMQClient) PublishPaymentTimeout(evt event.PaymentTimeoutEvent) error {
	args := m.Called(evt)
	return args.Error(0)
}

func (m *MockRabbitMQClient) ProcessPaymentSuccess(evt event.PaymentSuccessEvent) error {
	args := m.Called(evt)
	return args.Error(0)
}

func (m *MockRabbitMQClient) ProcessPaymentFailed(evt event.PaymentFailedEvent) error {
	args := m.Called(evt)
	return args.Error(0)
}

func (m *MockRabbitMQClient) ProcessPaymentTimeout(evt event.PaymentTimeoutEvent) error {
	args := m.Called(evt)
	return args.Error(0)
}

func TestCreateOrder(t *testing.T) {
	// Arrange
	mockRepo := new(MockOrderRepository)
	mockRabbitMQ := new(MockRabbitMQClient)

	service := NewOrderServiceImpl(mockRepo, mockRabbitMQ)

	order := &models.Order{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		TotalAmount: 49.99,
		Status:      models.PENDING,
	}

	// Set up mock expectations
	mockRepo.On("CreateOrder", order).Return(nil)
	mockRabbitMQ.On("PublishOrderCreated", mock.AnythingOfType("event.OrderCreatedEvent")).Return(nil)
	mockRabbitMQ.On("PublishPaymentTimeout", mock.AnythingOfType("event.PaymentTimeoutEvent")).Return(nil)

	// Act
	err := service.CreateOrder(order)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockRabbitMQ.AssertExpectations(t)
}
