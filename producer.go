package main

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func main() {
	// Dial to RabbitMQ.
	amqpConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Panicf("Failed to connect to RabbitMQ: %s", err)
	}
	defer amqpConn.Close()

	// Open a channel.
	amqpCh, err := amqpConn.Channel()
	if err != nil {
		log.Panicf("Failed to open a channel: %s", err)
	}
	defer amqpCh.Close()

	// Check publish confirms.
	amqpConfirms := make(chan amqp.Confirmation)
	amqpCh.NotifyPublish(amqpConfirms)
	go func() {
		for confirm := range amqpConfirms {
			if confirm.Ack {
				log.Printf("Confirmed delivery with delivery tag: %d", confirm.DeliveryTag)
			} else {
				log.Printf("Rejected delivery with delivery tag: %d", confirm.DeliveryTag)
			}
		}
	}()

	if err = amqpCh.Confirm(false); err != nil {
		log.Panicf("Failed to enable publish confirms: %s", err)
	}

	// Declare an exchange.
	err = amqpCh.ExchangeDeclare(
		"event", // name
		"topic", // type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Panicf("Failed to declare an exchange: %s", err)
	}

	// Publish a message.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message := "Hello message!"
	err = amqpCh.PublishWithContext(
		ctx,
		"event",         // exchange
		"event.created", // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(message),
		},
	)
	if err != nil {
		log.Panicf("Failed to publish a message: %s", err)
	}
	log.Printf("Published a message: %s", message)
}
