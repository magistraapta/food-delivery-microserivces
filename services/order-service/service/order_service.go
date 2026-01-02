package service

import (
	"log"
	"order-service/messaging"
	"order-service/messaging/event"
	"order-service/models"
	"order-service/repository"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(order *models.Order) error
	GetOrderById(id uuid.UUID) (*models.Order, error)
	ProcessPaymentSuccess(evt event.PaymentSuccessEvent) error
	ProcessPaymentFailed(evt event.PaymentFailedEvent) error
}

type OrderServiceImpl struct {
	orderRepository repository.OrderRepository
	rabbitMQClient  *messaging.RabbitMQClient
}

func NewOrderServiceImpl(orderRepository repository.OrderRepository, rabbitMQClient *messaging.RabbitMQClient) OrderService {
	return &OrderServiceImpl{
		orderRepository: orderRepository,
		rabbitMQClient:  rabbitMQClient,
	}
}

// CreateOrder creates a new order and publishes an order.created event
func (s *OrderServiceImpl) CreateOrder(order *models.Order) error {
	// Save order to database
	if err := s.orderRepository.CreateOrder(order); err != nil {
		return err
	}

	// Publish order.created event to RabbitMQ
	evt := event.OrderCreatedEvent{
		OrderID:         order.ID.String(),
		UserID:          order.UserID.String(),
		Amount:          order.TotalAmount,
		Currency:        "usd", // TODO: Add currency to Order model
		PaymentMethodID: "",    // TODO: Add payment method to Order model or request
	}

	if err := s.rabbitMQClient.PublishOrderCreated(evt); err != nil {
		log.Printf("Failed to publish order.created event for OrderID %s: %v", order.ID, err)
		// Note: Order is already created; consider implementing saga or compensation logic
		return err
	}

	log.Printf("Order created and event published: OrderID=%s", order.ID)
	return nil
}

func (s *OrderServiceImpl) GetOrderById(id uuid.UUID) (*models.Order, error) {
	return s.orderRepository.GetOrderById(id)
}

// ProcessPaymentSuccess handles payment.success events from Payment Service
func (s *OrderServiceImpl) ProcessPaymentSuccess(evt event.PaymentSuccessEvent) error {
	log.Printf("Processing payment success for OrderID: %s", evt.OrderID)

	orderID, err := uuid.Parse(evt.OrderID)
	if err != nil {
		log.Printf("Invalid OrderID format: %s", evt.OrderID)
		return err
	}

	// Update order status to CONFIRMED
	if err := s.orderRepository.UpdateOrderStatus(orderID, models.CONFIRMED); err != nil {
		log.Printf("Failed to update order status to CONFIRMED: %v", err)
		return err
	}

	log.Printf("Order %s status updated to CONFIRMED", evt.OrderID)
	return nil
}

// ProcessPaymentFailed handles payment.failed events from Payment Service
func (s *OrderServiceImpl) ProcessPaymentFailed(evt event.PaymentFailedEvent) error {
	log.Printf("Processing payment failure for OrderID: %s, Reason: %s", evt.OrderID, evt.FailureReason)

	orderID, err := uuid.Parse(evt.OrderID)
	if err != nil {
		log.Printf("Invalid OrderID format: %s", evt.OrderID)
		return err
	}

	// Update order status to PAYMENT_FAILED (you may need to add this status)
	if err := s.orderRepository.UpdateOrderStatus(orderID, "PAYMENT_FAILED"); err != nil {
		log.Printf("Failed to update order status to PAYMENT_FAILED: %v", err)
		return err
	}

	log.Printf("Order %s status updated to PAYMENT_FAILED", evt.OrderID)
	return nil
}
