package main

import (
	cons "github.com/replaygaming/consumer"
	"github.com/replaygaming/go-metrics/internal/amplitude"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "[METRICS] ", 0)

// Adapter is the interface required to start a new service to receive incoming
// events and forward them to the correct API
type Adapter interface {
	Start() (chan<- []byte, error)
}

var (
	topic           string
	subscription    string
	amplitudeAPIKey string
)

func warn(message string, err error) {
	logger.Printf("[WARN] %s: %s", message, err)
}

func fatal(message string, v ...interface{}) {
	logger.Fatalf("[FATAL] "+message, v...)
}

func info(message string, v ...interface{}) {
	logger.Printf("[INFO] "+message, v...)

}

func init() {
	info("Initializing App")

	topic = os.Getenv("METRICS_TOPIC")
	if topic == "" {
		topic = "metrics"
	}

	subscription = os.Getenv("METRICS_SUBSCRIPTION")
	if subscription == "" {
		subscription = "metrics_workers"
	}

	amplitudeAPIKey = os.Getenv("AMPLITUDE_API_KEY")
	if amplitudeAPIKey == "" {
		fatal("Amplitude API Key not provided")
	}

	info("INIT - Topic = %s, Subscription = %s, Amplitude API Key = %s", topic, subscription, amplitudeAPIKey)
}

func main() {
	info("Starting Service")

	// Start consumer queue
	consumer := cons.NewConsumer(topic, subscription)

	// Start consuming messages from queue
	messages, err := consumer.Consume()
	if err != nil {
		fatal("Failed to consume messages: %s", err)
	}

	// Configure Amplitude client and adapter channels
	amplitudeClient := amplitude.NewClient(amplitudeAPIKey)
	adapters := []Adapter{amplitudeClient}
	channels := make([]chan<- []byte, len(adapters))

	for i, a := range adapters {
		c, err := a.Start()
		if err != nil {
			fatal("Adapter failed to start: %s", err)
		}
		channels[i] = c
	}

	info("Service Started")

	// Listen for incoming events
	for m := range messages {
		for _, c := range channels {
			c <- m.Data()
		}
		m.Done(true)
	}
}
