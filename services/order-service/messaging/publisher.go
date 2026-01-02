package messaging

import (
	"context"
	"encoding/json"
	"log"
	"order-service/messaging/event"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// PublishToExchange publishes an event to an exchange with a routing key
func (c *RabbitMQClient) PublishToExchange(exchange, routingKey string, evt interface{}) error {
	body, err := json.Marshal(evt)
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

// PublishOrderCreated publishes an order.created event to the order events exchange
// This event is consumed by Payment Service to initiate payment processing
func (c *RabbitMQClient) PublishOrderCreated(evt event.OrderCreatedEvent) error {
	log.Printf("Publishing order.created event for OrderID: %s", evt.OrderID)
	return c.PublishToExchange(OrderEventsExchange, OrderCreatedRoutingKey, evt)
}
