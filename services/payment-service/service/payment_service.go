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
func (s *PaymentService) ProcessOrderCreatedEvent(event events.OrderCreatedEvent) error {
	log.Printf("Processing payment for order: %s, amount: %.2f %s", event.OrderID, event.Amount, event.Currency)

	// Create payment record in database
	payment := &models.Payment{
		ID:      uuid.New(),
		OrderID: event.OrderID,
		UserID:  event.UserID,
		Amount:  event.Amount,
		Status:  models.PaymentStatusPending,
	}

	if err := s.repo.CreatePayment(payment); err != nil {
		log.Printf("Failed to create payment record: %v", err)
		return err
	}

	// Process payment with Stripe
	paymentIntent, err := s.processStripePayment(event)
	if err != nil {
		log.Printf("Stripe payment failed for order %s: %v", event.OrderID, err)

		// Update payment status to failed
		s.repo.UpdateStatus(event.OrderID, models.PaymentStatusFailed)

		// Publish payment failed event
		failedEvent := events.PaymentFailedEvent{
			OrderID:       event.OrderID,
			CustomerID:    event.UserID,
			PaymentID:     payment.ID,
			FailureReason: err.Error(),
			FailureCode:   "stripe_error",
			Timestamp:     time.Now(),
		}

		if pubErr := s.rabbitMQClient.PublishPaymentFailed(failedEvent); pubErr != nil {
			log.Printf("Failed to publish payment failed event: %v", pubErr)
			return pubErr
		}
		return nil // Payment failed but event was published successfully
	}

	// Payment succeeded
	s.repo.UpdateStatus(event.OrderID, models.PaymentStatusSuccess)
	s.repo.UpdateIntentId(event.OrderID, paymentIntent.ID)

	// Publish payment success event
	successEvent := events.PaymentSuccessEvent{
		OrderID:               event.OrderID,
		UserID:                event.UserID,
		Amount:                event.Amount,
		Currency:              event.Currency,
		StripePaymentIntentID: paymentIntent.ID,
		StripeChargeID:        getChargeID(paymentIntent),
		Timestamp:             time.Now(),
	}

	if err := s.rabbitMQClient.PublishPaymentSuccess(successEvent); err != nil {
		log.Printf("Failed to publish payment success event: %v", err)
		return err
	}

	log.Printf("Payment successful for order: %s, PaymentIntent: %s", event.OrderID, paymentIntent.ID)
	return nil
}

// processStripePayment processes the payment via Stripe
func (s *PaymentService) processStripePayment(event events.OrderCreatedEvent) (*stripeGo.PaymentIntent, error) {
	log.Printf("Creating Stripe payment for order: %s, amount: %.2f %s", event.OrderID, event.Amount, event.Currency)

	// Use the payment method from the event, or default to test card
	paymentMethodID := event.PaymentMethodID
	if paymentMethodID == "" {
		paymentMethodID = "pm_card_visa" // Stripe test payment method for sandbox
	}

	currency := event.Currency
	if currency == "" {
		currency = "usd"
	}

	// Create and confirm payment intent in one call (for sandbox testing)
	paymentIntent, err := s.stripeClient.CreateAndConfirmPaymentIntent(
		event.OrderID.String(),
		event.Amount,
		currency,
		paymentMethodID,
	)

	if err != nil {
		return nil, err
	}

	// Check if payment was successful
	if paymentIntent.Status != stripeGo.PaymentIntentStatusSucceeded {
		log.Printf("Payment not succeeded, status: %s", paymentIntent.Status)
	}

	return paymentIntent, nil
}

// getChargeID extracts the charge ID from a payment intent
func getChargeID(pi *stripeGo.PaymentIntent) string {
	if pi.LatestCharge != nil {
		return pi.LatestCharge.ID
	}
	return ""
}

// HandleStripeWebhook processes Stripe webhook events
// This is called when Stripe sends payment confirmation via webhook
func (s *PaymentService) HandleStripeWebhook(paymentIntentID string, status string, orderID string) error {
	log.Printf("Received Stripe webhook for PaymentIntent: %s, Status: %s", paymentIntentID, status)

	switch status {
	case "succeeded":
		// Get payment intent details from Stripe
		pi, err := s.stripeClient.GetPaymentIntent(paymentIntentID)
		if err != nil {
			return err
		}

		orderUUID, err := uuid.Parse(orderID)
		if err != nil {
			return err
		}

		// Update payment status
		s.repo.UpdateStatus(orderUUID, models.PaymentStatusSuccess)

		// Publish success event
		successEvent := events.PaymentSuccessEvent{
			OrderID:               orderUUID,
			StripePaymentIntentID: pi.ID,
			StripeChargeID:        getChargeID(pi),
			Amount:                float64(pi.Amount) / 100, // Convert from cents
			Currency:              string(pi.Currency),
			Timestamp:             time.Now(),
		}

		return s.rabbitMQClient.PublishPaymentSuccess(successEvent)

	case "payment_failed":
		orderUUID, _ := uuid.Parse(orderID)
		s.repo.UpdateStatus(orderUUID, models.PaymentStatusFailed)

		failedEvent := events.PaymentFailedEvent{
			OrderID:       orderUUID,
			FailureReason: "Payment failed",
			FailureCode:   "card_declined",
			Timestamp:     time.Now(),
		}

		return s.rabbitMQClient.PublishPaymentFailed(failedEvent)
	}

	return nil
}
