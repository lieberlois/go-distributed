package qutils

import (
	"github.com/streadway/amqp"
	"log"
)

const SensorListQueue = "SensorList"

func GetChannel(url string) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(url)
	failOnErr(err, "Failed to establish connection to message broker")

	ch, err := conn.Channel()
	failOnErr(err, "Failed to open a channel")

	return conn, ch
}

func GetQueue(name string, ch *amqp.Channel) *amqp.Queue {
	q, err := ch.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnErr(err, "Failed to declare a queue")

	return &q
}

func failOnErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", err, msg)
	}
}
