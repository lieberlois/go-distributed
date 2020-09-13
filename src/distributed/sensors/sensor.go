package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"github.com/streadway/amqp"
	"go-distributed/src/distributed/dto"
	"go-distributed/src/distributed/qutils"
	"log"
	"math/rand"
	"strconv"
	"time"
)

var url = "amqp://guest:guest@localhost:5672"

var name = flag.String("name", "sensor", "name of the sensor")
var freq = flag.Uint("freq", 5, "update frequency in cycles/s")
var max = flag.Float64("max", 5., "maximum value")
var min = flag.Float64("min", 1., "minimum value")
var stepSize = flag.Float64("step", 0.1, "maximum allowable change per measurement")

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

var value = r.Float64()*(*max-*min) + *min
var nom = (*max-*min)/2 + *min

func main() {
	flag.Parse()

	conn, ch := qutils.GetChannel(url)
	defer conn.Close()
	defer ch.Close()

	dataQueue := qutils.GetQueue(*name, ch, false) // asserts that the queue was created properly (sending only needs exchange)

	publishQueueName(ch)

	discoveryQueue := qutils.GetQueue("", ch, true)
	_ = ch.QueueBind(
		discoveryQueue.Name,
		"",
		qutils.SensorDiscoveryExchange,
		false,
		nil,
	)

	go listenForDiscoverRequests(discoveryQueue.Name, ch)

	dur, _ := time.ParseDuration(strconv.Itoa(1000/(int(*freq))) + "ms")
	signal := time.Tick(dur)
	buf := new(bytes.Buffer)

	for range signal {
		calcValue()
		reading := dto.SensorMessage{
			Name:      *name,
			Value:     value,
			Timestamp: time.Now(),
		}
		buf.Reset()
		enc := gob.NewEncoder(buf)
		_ = enc.Encode(reading)

		msg := amqp.Publishing{
			Body: buf.Bytes(),
		}

		_ = ch.Publish(
			"",
			dataQueue.Name,
			false,
			false,
			msg,
		)

		log.Printf("Reading sent: %v", value)
	}
}

func listenForDiscoverRequests(name string, ch *amqp.Channel) {
	msgs, _ := ch.Consume(
		name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	for range msgs {
		publishQueueName(ch)
	}
}

func publishQueueName(ch *amqp.Channel) {
	// Sensor Queue Discovery
	msg := amqp.Publishing{
		Body: []byte(*name),
	}

	// Publish creation of new queue
	_ = ch.Publish(
		"amq.fanout",
		"",
		false,
		false,
		msg,
	)
}

func calcValue() {
	var maxStep, minStep float64

	if value < nom {
		maxStep = *stepSize
		minStep = -1 * *stepSize * (value - *min) / (nom - *min)
	} else {
		maxStep = *stepSize * (*max - value) / (*max - nom)
		minStep = -1 * *stepSize
	}

	value += r.Float64()*(maxStep-minStep) + minStep
}
