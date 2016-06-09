package main

import (
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

var (
	amqpURL         string
	amqpQueue       string
	amplitudeAPIKey string
)

func warn(message string, err error) {
	logger.Printf("[WARN] %s: %s", message, err)
}

func fatal(message string, err error) {
	logger.Fatalf("[FATAL] %s: %s", message, err)
}

func init() {
	logger.Printf("[INFO] Initializing App")

	amqpURL = os.Getenv("AMQP_URL")
	amqpQueue = os.Getenv("AMQP_QUEUE")
	amplitudeAPIKey = os.Getenv("AMPLITUDE_API_KEY")

	logger.Printf("[INFO] INIT - AMPQ URL = %s, AMQP Queue = %s, Amplitude API Key = %s", amqpURL, amqpQueue, amplitudeAPIKey)
}

func main() {
	// Start consumer queue
	// NewConsumer(amqpURI, exchange, exchangeType, queueName, key, ctag string) (*Consumer, error)
	c, err := amqp.NewConsumer(amqpURL, "metrics_ex", "fanout", amqpQueue, "", "metrics")
	if err != nil {
		fatal("AMQP Consumer Failed", err)
	}

	messages, err := c.Consume(amqpQueue)
	if err != nil {
		fatal("AMQP Queue Failed", err)
	}

	// Start event adapters
	a := amplitude.NewClient(amplitudeAPIKey)

	adapters := []Adapter{a}
	channels := make([]chan<- []byte, len(adapters))

	for i, a := range adapters {
		c, err := a.Start()
		if err != nil {
			logger.Fatalf("[FATAL] Adapter failed to start %s", err)
		}
		channels[i] = c
	}

	logger.Printf("[INFO] Starting Service")

	// Listen for incoming events
	for m := range messages {
		for _, c := range channels {
			c <- m.Body
		}
		m.Ack(false)
	}
	c.Done <- nil
}
