package messaging

import (
	"context"
	"encoding/json"
	"log"
	"payment-service/messaging/events"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// PublishToExchange publishes an event to an exchange with a routing key
func (c *RabbitMQClient) PublishToExchange(exchange, routingKey string, event interface{}) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = c.Channel.PublishWithContext(
		ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent, // message will survive broker restart
			Timestamp:    time.Now(),
			Body:         body,
		},
	)
	if err != nil {
		return err
	}

	log.Printf("Published message to exchange: %s with routing key: %s", exchange, routingKey)
	return nil
}

// PublishPaymentSuccess publishes a payment.success event to the payment events exchange
// This event is consumed by Order Service to update the order status
func (c *RabbitMQClient) PublishPaymentSuccess(event events.PaymentSuccessEvent) error {
	log.Printf("Publishing payment.success event for OrderID: %s", event.OrderID)
	return c.PublishToExchange(PaymentEventsExchange, PaymentSuccessRoutingKey, event)
}

// PublishPaymentFailed publishes a payment.failed event to the payment events exchange
// This event is consumed by Order Service to update the order status
func (c *RabbitMQClient) PublishPaymentFailed(event events.PaymentFailedEvent) error {
	log.Printf("Publishing payment.failed event for OrderID: %s", event.OrderID)
	return c.PublishToExchange(PaymentEventsExchange, PaymentFailedRoutingKey, event)
}
