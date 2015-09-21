package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	amqp "github.com/replaygaming/amqp-consumer"
)

var (
	env             = flag.String("env", "development", "Environment: development or production")
	amqpURL         = flag.String("amqp-url", "amqp://guest:guest@localhost:5672/metrics", "AMQP URL")
	amqpQueue       = flag.String("amqp-queue", "metrics", "AMQP Queue name")
	amplitudeAPIKey = flag.String("amplitude-api-key", "", "Amplitude API Key")
)

// Adapter is the interface required to start a new service to receive incoming
// events and forward them to the correct API
type Adapter interface {
	Start() (chan<- Event, error)
}

func init() {
	flag.Parse()
}

func main() {
	// Start consumer queue
	c, err := amqp.NewConsumer(*amqpURL, "metrics_ex", "fanout", *amqpQueue, "", "metrics")
	if err != nil {
		log.Fatalf("[FATAL] AMQP consumer failed %s", err)
	}
	messages, err := c.Consume(*amqpQueue)
	if err != nil {
		log.Fatalf("[FATAL] AMQP queue failed %s", err)
	}

	// Start event adapters
	amplitude := &Amplitude{
		APIKey: *amplitudeAPIKey,
	}

	adapters := []Adapter{amplitude}
	chans := make([]chan<- Event, len(adapters))

	for i, a := range adapters {
		c, err := a.Start()
		if err != nil {
			log.Fatalf("[FATAL] Adapter failed to start %s", err)
		}
		chans[i] = c
	}

	log.Printf("[INFO] start %s", os.Args[1:])

	// Listen for incoming events
	for m := range messages {
		e := &Event{}
		err := json.Unmarshal(m.Body, e)
		if err != nil {
			log.Printf("[WARN] JSON conversion failed %s", err)
		} else {
			for _, c := range chans {
				c <- *e
			}
		}
		m.Ack(false)
	}
	c.Done <- nil
}
