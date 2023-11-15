package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
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

	// Declare a queue.
	amqpQ, err := amqpCh.QueueDeclare(
		"",    // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Panicf("Failed to declare a queue: %s", err)
	}

	// Set QoS.
	err = amqpCh.Qos(1, 0, false)

	// Bind the queue to the exchange.
	err = amqpCh.QueueBind(
		amqpQ.Name,      // queue
		"event.created", // routing key
		"event",         // exchange
		false,
		nil,
	)
	if err != nil {
		log.Panicf("Failed to bind a queue: %s", err)
	}

	// Declare a consumer.
	messages, err := amqpCh.Consume(
		amqpQ.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Panicf("Failed to register a consumer: %s", err)
	}

	// Consume messages.
	var forever chan struct{}

	go func() {
		for msg := range messages {
			log.Printf("Received a message: %s", msg.Body)
			_ = msg.Ack(false)
		}
	}()

	<-forever
}
