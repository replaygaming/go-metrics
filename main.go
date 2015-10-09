package main

import (
	"flag"
	"log"
	"os"

	amqp "github.com/replaygaming/amqp-consumer"
	"github.com/replaygaming/go-metrics/internal/amplitude"
)

var logger = log.New(os.Stdout, "[METRICS] ", 0)

// Adapter is the interface required to start a new service to receive incoming
// events and forward them to the correct API
type Adapter interface {
	Start() (chan<- []byte, error)
}

func main() {
	var (
		amqpURL = flag.String("amqp-url",
			"amqp://guest:guest@localhost:5672/metrics", "AMQP URL")
		amqpQueue       = flag.String("amqp-queue", "metrics", "AMQP Queue name")
		amplitudeAPIKey = flag.String("amplitude-api-key", "", "Amplitude API Key")
	)
	flag.Parse()

	// Start consumer queue
	c, err := amqp.NewConsumer(*amqpURL, "metrics_ex", "fanout", *amqpQueue, "",
		"metrics")
	if err != nil {
		logger.Fatalf("[FATAL] AMQP consumer failed %s", err)
	}
	messages, err := c.Consume(*amqpQueue)
	if err != nil {
		logger.Fatalf("[FATAL] AMQP queue failed %s", err)
	}

	// Start event adapters
	a := amplitude.NewClient(*amplitudeAPIKey)

	adapters := []Adapter{a}
	chans := make([]chan<- []byte, len(adapters))

	for i, a := range adapters {
		c, err := a.Start()
		if err != nil {
			logger.Fatalf("[FATAL] Adapter failed to start %s", err)
		}
		chans[i] = c
	}

	logger.Printf("[INFO] Starting metrics service %s", os.Args[1:])

	// Listen for incoming events
	for m := range messages {
		for _, c := range chans {
			c <- m.Body
		}
		m.Ack(false)
	}
	c.Done <- nil
}
