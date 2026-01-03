package service

import (
	"log"
	"payment-service/messaging"
	"payment-service/messaging/events"
	"payment-service/models"
	"payment-service/repository"
	"payment-service/stripe"
	"time"

	"github.com/google/uuid"
	stripeGo "github.com/stripe/stripe-go/v84"
)

type PaymentService struct {
	repo           repository.PaymentRepository
	stripeClient   *stripe.StripeClient
	rabbitMQClient *messaging.RabbitMQClient
}

func NewPaymentService(
	repo repository.PaymentRepository,
	stripeClient *stripe.StripeClient,
	rabbitMQClient *messaging.RabbitMQClient,
) *PaymentService {
	return &PaymentService{
		repo:           repo,
		stripeClient:   stripeClient,
		rabbitMQClient: rabbitMQClient,
	}
}

// ProcessOrderCreatedEvent handles incoming order.created events
// Creates a Stripe Checkout Session and stores the payment URL
func (s *PaymentService) ProcessOrderCreatedEvent(event events.OrderCreatedEvent) error {
	log.Printf("Processing payment for order: %s, amount: %.2f %s", event.OrderID, event.Amount, event.Currency)

	// Create payment record in database
	payment := &models.Payment{
		ID:       uuid.New(),
		OrderID:  event.OrderID,
		UserID:   event.UserID,
		Amount:   event.Amount,
		Currency: event.Currency,
		Status:   models.PaymentStatusPending,
	}

	if event.Currency == "" {
		payment.Currency = "usd"
	}

	if err := s.repo.CreatePayment(payment); err != nil {
		log.Printf("Failed to create payment record: %v", err)
		return err
	}

	// Create Stripe Checkout Session
	checkoutSession, err := s.stripeClient.CreateCheckoutSession(
		event.OrderID.String(),
		event.Amount,
		payment.Currency,
		"Food Order", // Product name
	)

	if err != nil {
		log.Printf("Failed to create Stripe Checkout Session for order %s: %v", event.OrderID, err)

		// Update payment status to failed
		s.repo.UpdateStatus(event.OrderID, models.PaymentStatusFailed)

		// Publish payment failed event
		failedEvent := events.PaymentFailedEvent{
			OrderID:       event.OrderID,
			CustomerID:    event.UserID,
			PaymentID:     payment.ID,
			FailureReason: err.Error(),
			FailureCode:   "checkout_session_failed",
			Timestamp:     time.Now(),
		}

		if pubErr := s.rabbitMQClient.PublishPaymentFailed(failedEvent); pubErr != nil {
			log.Printf("Failed to publish payment failed event: %v", pubErr)
			return pubErr
		}
		return nil
	}

	// Store checkout session info in database
	if err := s.repo.UpdateCheckoutSession(event.OrderID, checkoutSession.ID, checkoutSession.URL); err != nil {
		log.Printf("Failed to update checkout session for order %s: %v", event.OrderID, err)
	}

	// Publish checkout created event with the payment URL
	checkoutEvent := events.PaymentCheckoutCreatedEvent{
		OrderID:     event.OrderID,
		UserID:      event.UserID,
		Amount:      event.Amount,
		Currency:    payment.Currency,
		CheckoutURL: checkoutSession.URL,
		SessionID:   checkoutSession.ID,
		ExpiresAt:   time.Now().Add(5 * time.Minute), // Matches our timeout
		Timestamp:   time.Now(),
	}

	if err := s.rabbitMQClient.PublishPaymentCheckoutCreated(checkoutEvent); err != nil {
		log.Printf("Failed to publish payment checkout created event: %v", err)
		// Non-critical - checkout was created successfully
	}

	log.Printf("Checkout session created for order: %s, URL: %s", event.OrderID, checkoutSession.URL)
	return nil
}

// GetCheckoutURL retrieves the checkout URL for an order
func (s *PaymentService) GetCheckoutURL(orderID uuid.UUID) (string, error) {
	payment, err := s.repo.FindByOrderId(orderID)
	if err != nil {
		return "", err
	}
	return payment.CheckoutURL, nil
}

// GetPaymentByOrderID retrieves payment info for an order
func (s *PaymentService) GetPaymentByOrderID(orderID uuid.UUID) (*models.Payment, error) {
	return s.repo.FindByOrderId(orderID)
}

// HandleCheckoutSessionCompleted processes completed checkout sessions from Stripe webhook
func (s *PaymentService) HandleCheckoutSessionCompleted(session *stripeGo.CheckoutSession) error {
	log.Printf("Handling completed checkout session: %s", session.ID)

	// Get order ID from metadata
	orderIDStr, ok := session.Metadata["order_id"]
	if !ok {
		log.Printf("No order_id in checkout session metadata")
		return nil
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		log.Printf("Invalid order ID format: %s", orderIDStr)
		return err
	}

	// Get payment record
	payment, err := s.repo.FindByOrderId(orderID)
	if err != nil {
		log.Printf("Payment not found for order %s: %v", orderIDStr, err)
		return err
	}

	// Update payment status to success
	if err := s.repo.UpdateStatus(orderID, models.PaymentStatusSuccess); err != nil {
		log.Printf("Failed to update payment status: %v", err)
		return err
	}

	// Update payment intent ID if available
	if session.PaymentIntent != nil {
		s.repo.UpdatePaymentIntent(orderID, session.PaymentIntent.ID, "")
	}

	// Publish payment success event
	successEvent := events.PaymentSuccessEvent{
		OrderID:               orderID,
		UserID:                payment.UserID,
		Amount:                payment.Amount,
		Currency:              payment.Currency,
		StripePaymentIntentID: "",
		StripeChargeID:        "",
		Timestamp:             time.Now(),
	}

	if session.PaymentIntent != nil {
		successEvent.StripePaymentIntentID = session.PaymentIntent.ID
	}

	if err := s.rabbitMQClient.PublishPaymentSuccess(successEvent); err != nil {
		log.Printf("Failed to publish payment success event: %v", err)
		return err
	}

	log.Printf("Payment successful for order: %s via checkout session: %s", orderIDStr, session.ID)
	return nil
}

// HandleCheckoutSessionExpired processes expired checkout sessions
func (s *PaymentService) HandleCheckoutSessionExpired(session *stripeGo.CheckoutSession) error {
	log.Printf("Handling expired checkout session: %s", session.ID)

	orderIDStr, ok := session.Metadata["order_id"]
	if !ok {
		return nil
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return err
	}

	// Get payment record
	payment, err := s.repo.FindByOrderId(orderID)
	if err != nil {
		return err
	}

	// Only update if still pending
	if payment.Status == models.PaymentStatusPending {
		s.repo.UpdateStatus(orderID, models.PaymentStatusExpired)

		// Publish payment failed event
		failedEvent := events.PaymentFailedEvent{
			OrderID:       orderID,
			CustomerID:    payment.UserID,
			PaymentID:     payment.ID,
			FailureReason: "Checkout session expired",
			FailureCode:   "checkout_expired",
			Timestamp:     time.Now(),
		}

		if err := s.rabbitMQClient.PublishPaymentFailed(failedEvent); err != nil {
			log.Printf("Failed to publish payment failed event: %v", err)
		}
	}

	return nil
}
