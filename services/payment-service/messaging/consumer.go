package messaging

import (
	"context"
	"encoding/json"
	"log"
	"payment-service/messaging/events"

	amqp "github.com/rabbitmq/amqp091-go"
)

// OrderEventHandler defines the callback function for processing order events
type OrderEventHandler func(event events.OrderCreatedEvent) error

// ConsumeOrderEvents starts consuming order.created events from RabbitMQ
// The handler function is called for each received event
func (c *RabbitmqClientImpl) ConsumeOrderEvents(ctx context.Context, handler OrderEventHandler) error {
	// Set QoS (Quality of Service) - process one message at a time
	err := c.Channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return err
	}

	msgs, err := c.Channel.Consume(
		OrderCreatedQueue, // queue
		"payment-service", // consumer tag
		false,             // auto-ack (false = manual ack for reliability)
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	if err != nil {
		return err
	}

	// Start consuming in a goroutine
	go func() {
		log.Printf("Started consuming order.created events from queue: %s", OrderCreatedQueue)

		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping order events consumer (context cancelled)...")
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Println("Order events channel closed")
					return
				}
				c.processOrderMessage(msg, handler)
			}
		}
	}()

	return nil
}

// processOrderMessage handles a single order.created message
func (c *RabbitmqClientImpl) processOrderMessage(msg amqp.Delivery, handler OrderEventHandler) {
	log.Printf("Received message from queue: %s", OrderCreatedQueue)

	var event events.OrderCreatedEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("Error unmarshaling order event: %v", err)
		// Reject the message and don't requeue (bad message format)
		msg.Nack(false, false)
		return
	}

	log.Printf("Processing order.created event - OrderID: %s, Amount: %.2f %s",
		event.OrderID, event.Amount, event.Currency)

	// Call the handler to process the event
	if err := handler(event); err != nil {
		log.Printf("Error processing order event for OrderID %s: %v", event.OrderID, err)
		// Requeue the message for retry
		msg.Nack(false, true)
		return
	}

	// Acknowledge the message - successfully processed
	if err := msg.Ack(false); err != nil {
		log.Printf("Error acknowledging message: %v", err)
	} else {
		log.Printf("Successfully processed and acknowledged order.created event for OrderID: %s", event.OrderID)
	}
}
