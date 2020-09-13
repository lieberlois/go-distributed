package qutils

import (
	"github.com/streadway/amqp"
	"log"
)

const SensorDiscoveryExchange = "SensorDiscovery"
const PersistReadingsQueue = "PersistReading"
const WebappSourceExchange = "WebappSources"
const WebappReadingsExchange = "WebappReadings"
const WebappDiscoveryQueue = "WebappDiscovery"


func GetChannel(url string) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(url)
	failOnErr(err, "Failed to establish connection to message broker")

	ch, err := conn.Channel()
	failOnErr(err, "Failed to open a channel")

	return conn, ch
}

func GetQueue(name string, ch *amqp.Channel, autoDelete bool) *amqp.Queue {
	q, err := ch.QueueDeclare(
		name,
		false,
		autoDelete,
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
