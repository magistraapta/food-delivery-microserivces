package messaging

import (
	"context"
	"encoding/json"
	"log"
	"order-service/messaging/event"

	amqp "github.com/rabbitmq/amqp091-go"
)

// PaymentSuccessHandler defines the callback function for processing payment success events
type PaymentSuccessHandler func(event event.PaymentSuccessEvent) error

// PaymentFailedHandler defines the callback function for processing payment failed events
type PaymentFailedHandler func(event event.PaymentFailedEvent) error

// ConsumePaymentEvents starts consuming both payment.success and payment.failed events from RabbitMQ
func (c *RabbitMQClient) ConsumePaymentEvents(
	ctx context.Context,
	successHandler PaymentSuccessHandler,
	failedHandler PaymentFailedHandler,
) error {
	// Set QoS (Quality of Service) - process one message at a time
	err := c.Channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return err
	}

	// Start consuming payment.success events
	successMsgs, err := c.Channel.Consume(
		PaymentSuccessQueue,     // queue
		"order-service-success", // consumer tag
		false,                   // auto-ack (false = manual ack for reliability)
		false,                   // exclusive
		false,                   // no-local
		false,                   // no-wait
		nil,                     // args
	)
	if err != nil {
		return err
	}

	// Start consuming payment.failed events
	failedMsgs, err := c.Channel.Consume(
		PaymentFailedQueue,     // queue
		"order-service-failed", // consumer tag
		false,                  // auto-ack
		false,                  // exclusive
		false,                  // no-local
		false,                  // no-wait
		nil,                    // args
	)
	if err != nil {
		return err
	}

	// Start consuming in a goroutine
	go func() {
		log.Printf("Started consuming payment events from queues: %s, %s", PaymentSuccessQueue, PaymentFailedQueue)

		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping payment events consumer (context cancelled)...")
				return
			case msg, ok := <-successMsgs:
				if !ok {
					log.Println("Payment success channel closed")
					return
				}
				c.processPaymentSuccessMessage(msg, successHandler)
			case msg, ok := <-failedMsgs:
				if !ok {
					log.Println("Payment failed channel closed")
					return
				}
				c.processPaymentFailedMessage(msg, failedHandler)
			}
		}
	}()

	return nil
}

// processPaymentSuccessMessage handles a single payment.success message
func (c *RabbitMQClient) processPaymentSuccessMessage(msg amqp.Delivery, handler PaymentSuccessHandler) {
	log.Printf("Received payment.success message")

	var evt event.PaymentSuccessEvent
	if err := json.Unmarshal(msg.Body, &evt); err != nil {
		log.Printf("Error unmarshaling payment success event: %v", err)
		// Reject the message and don't requeue (bad message format)
		msg.Nack(false, false)
		return
	}

	log.Printf("Processing payment.success event - OrderID: %s, Amount: %.2f %s",
		evt.OrderID, evt.Amount, evt.Currency)

	// Call the handler to process the event
	if err := handler(evt); err != nil {
		log.Printf("Error processing payment success event for OrderID %s: %v", evt.OrderID, err)
		// Requeue the message for retry
		msg.Nack(false, true)
		return
	}

	// Acknowledge the message - successfully processed
	if err := msg.Ack(false); err != nil {
		log.Printf("Error acknowledging payment success message: %v", err)
	} else {
		log.Printf("Successfully processed payment.success event for OrderID: %s", evt.OrderID)
	}
}

// processPaymentFailedMessage handles a single payment.failed message
func (c *RabbitMQClient) processPaymentFailedMessage(msg amqp.Delivery, handler PaymentFailedHandler) {
	log.Printf("Received payment.failed message")

	var evt event.PaymentFailedEvent
	if err := json.Unmarshal(msg.Body, &evt); err != nil {
		log.Printf("Error unmarshaling payment failed event: %v", err)
		// Reject the message and don't requeue (bad message format)
		msg.Nack(false, false)
		return
	}

	log.Printf("Processing payment.failed event - OrderID: %s, Reason: %s",
		evt.OrderID, evt.FailureReason)

	// Call the handler to process the event
	if err := handler(evt); err != nil {
		log.Printf("Error processing payment failed event for OrderID %s: %v", evt.OrderID, err)
		// Requeue the message for retry
		msg.Nack(false, true)
		return
	}

	// Acknowledge the message - successfully processed
	if err := msg.Ack(false); err != nil {
		log.Printf("Error acknowledging payment failed message: %v", err)
	} else {
		log.Printf("Successfully processed payment.failed event for OrderID: %s", evt.OrderID)
	}
}

// PaymentTimeoutHandler defines the callback function for processing payment timeout events
type PaymentTimeoutHandler func(event event.PaymentTimeoutEvent) error

// ConsumePaymentTimeoutEvents starts consuming payment timeout events from RabbitMQ
// These events arrive after a 5-minute delay to check if payment was completed
func (c *RabbitMQClient) ConsumePaymentTimeoutEvents(ctx context.Context, handler PaymentTimeoutHandler) error {
	// Start consuming payment.timeout events
	timeoutMsgs, err := c.Channel.Consume(
		PaymentTimeoutQueue,     // queue (the final queue, not the delay queue)
		"order-service-timeout", // consumer tag
		false,                   // auto-ack
		false,                   // exclusive
		false,                   // no-local
		false,                   // no-wait
		nil,                     // args
	)
	if err != nil {
		return err
	}

	go func() {
		log.Printf("Started consuming payment timeout events from queue: %s", PaymentTimeoutQueue)

		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping payment timeout consumer (context cancelled)...")
				return
			case msg, ok := <-timeoutMsgs:
				if !ok {
					log.Println("Payment timeout channel closed")
					return
				}
				c.processPaymentTimeoutMessage(msg, handler)
			}
		}
	}()

	return nil
}

// processPaymentTimeoutMessage handles a single payment timeout message
func (c *RabbitMQClient) processPaymentTimeoutMessage(msg amqp.Delivery, handler PaymentTimeoutHandler) {
	log.Printf("Received payment.timeout message (5 minutes elapsed)")

	var evt event.PaymentTimeoutEvent
	if err := json.Unmarshal(msg.Body, &evt); err != nil {
		log.Printf("Error unmarshaling payment timeout event: %v", err)
		msg.Nack(false, false)
		return
	}

	log.Printf("Processing payment.timeout event - OrderID: %s, CreatedAt: %s",
		evt.OrderID, evt.CreatedAt)

	// Call the handler to process the event
	if err := handler(evt); err != nil {
		log.Printf("Error processing payment timeout event for OrderID %s: %v", evt.OrderID, err)
		// Requeue the message for retry
		msg.Nack(false, true)
		return
	}

	// Acknowledge the message - successfully processed
	if err := msg.Ack(false); err != nil {
		log.Printf("Error acknowledging payment timeout message: %v", err)
	} else {
		log.Printf("Successfully processed payment.timeout event for OrderID: %s", evt.OrderID)
	}
}
