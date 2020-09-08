package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	go client()
	go server()

	// Keep main thread alive
	var s string
	_, _ = fmt.Scanln(&s)
}

func client() {
	conn, ch, q := getQueue()
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		q.Name,
		"",  // Name of the consumer, empty string automatically assigns one
		true, // If the task processing the message can fail, this should be false
		false,
		false,
		false,
		nil,
	)
	failOnErr(err, "Failed to register a consumer")

	for msg := range msgs {
		log.Printf("Received message: %s", msg.Body)
	}
}

func server() {
	conn, ch, q := getQueue()
	defer conn.Close()
	defer ch.Close()

	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("Hello RabbitMQ"),
	}

	for {
		_ = ch.Publish(
			"",  // Default exchange
			q.Name,
			false,
			false,
			msg,
		)
	}
}

func getQueue() (*amqp.Connection, *amqp.Channel, *amqp.Queue) {
	conn, err := amqp.Dial("amqp://guest@localhost:5672")
	failOnErr(err, "Failed to establish connection to RabbitMQ")

	ch, err := conn.Channel()
	failOnErr(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,  // Delete messages when no consumer is active on the queue
		false,
		false,
		nil,
	)
	failOnErr(err, "Failed to declare a queue")

	return conn, ch, &q
}

func failOnErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", err, msg)
	}
}
