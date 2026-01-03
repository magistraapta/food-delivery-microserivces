package service

import (
	"log"
	"order-service/messaging"
	"order-service/messaging/event"
	"order-service/models"
	"order-service/repository"
	"time"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(order *models.Order) error
	GetOrderById(id uuid.UUID) (*models.Order, error)
	ProcessPaymentSuccess(evt event.PaymentSuccessEvent) error
	ProcessPaymentFailed(evt event.PaymentFailedEvent) error
	ProcessPaymentTimeout(evt event.PaymentTimeoutEvent) error
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

	// Schedule payment timeout check (5 minutes from now)
	timeoutEvt := event.PaymentTimeoutEvent{
		OrderID:   order.ID.String(),
		CreatedAt: time.Now(),
	}

	if err := s.rabbitMQClient.PublishPaymentTimeout(timeoutEvt); err != nil {
		log.Printf("Failed to schedule payment timeout for OrderID %s: %v", order.ID, err)
		// Non-critical error - order was created successfully
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

	// Update order status to PAYMENT_FAILED
	if err := s.orderRepository.UpdateOrderStatus(orderID, models.PAYMENT_FAILED); err != nil {
		log.Printf("Failed to update order status to PAYMENT_FAILED: %v", err)
		return err
	}

	log.Printf("Order %s status updated to PAYMENT_FAILED", evt.OrderID)
	return nil
}

// ProcessPaymentTimeout handles timeout events for orders that haven't been paid
// If the order is still PENDING after 5 minutes, it will be cancelled
func (s *OrderServiceImpl) ProcessPaymentTimeout(evt event.PaymentTimeoutEvent) error {
	log.Printf("Processing payment timeout for OrderID: %s (created at: %s)", evt.OrderID, evt.CreatedAt)

	orderID, err := uuid.Parse(evt.OrderID)
	if err != nil {
		log.Printf("Invalid OrderID format: %s", evt.OrderID)
		return err
	}

	// Get current order status
	order, err := s.orderRepository.GetOrderById(orderID)
	if err != nil {
		log.Printf("Failed to get order %s: %v", evt.OrderID, err)
		return err
	}

	// Only cancel if order is still PENDING (payment not completed)
	if order.Status == models.PENDING {
		log.Printf("Order %s is still PENDING after 5 minutes - cancelling order", evt.OrderID)

		if err := s.orderRepository.UpdateOrderStatus(orderID, models.CANCELLED); err != nil {
			log.Printf("Failed to cancel order %s: %v", evt.OrderID, err)
			return err
		}

		log.Printf("Order %s cancelled due to payment timeout", evt.OrderID)
	} else {
		log.Printf("Order %s status is %s - no action needed", evt.OrderID, order.Status)
	}

	return nil
}
