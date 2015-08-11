package main

import (
	"encoding/json"
	"flag"
	"log"

	amqp "github.com/replaygaming/amqp-consumer"
)

var (
	env         = flag.String("env", "development", "Environment: development or production")
	amqpURL     = flag.String("amqp-url", "amqp://guest:guest@localhost:5672/metrics", "AMQP URL")
	gaGameKey   = flag.String("ga-game-key", "", "GameAnalytics GameKey")
	gaSecretKey = flag.String("ga-secret-key", "", "GameAnalytics SecretKey")
)

// Event contains the generic information received from Replay Poker
type Event struct {
	Version   uint
	Type      string `json:"event"`
	UserID    string `json:"user_id"`
	Timestamp string
	Session   struct {
		UUID   string `json:"uuid"`
		Number uint
	}
	Properties json.RawMessage
}

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
	c, err := amqp.NewConsumer(*amqpURL, "metrics_ex", "fanout", "", "", "metrics")
	if err != nil {
		log.Fatalf("[FATAL] AMQP consumer failed %s", err)
	}
	messages, err := c.Consume("")
	if err != nil {
		log.Fatalf("[FATAL] AMQP queue failed %s", err)
	}

	// Start event adapters
	ga := &GameAnalytics{
		GameKey:     *gaGameKey,
		SecretKey:   *gaSecretKey,
		Environment: *env,
	}
	adapters := []Adapter{ga}
	chans := make([]chan<- Event, len(adapters))

	for i, a := range adapters {
		c, err := a.Start()
		if err != nil {
			log.Fatalf("[FATAL] Adapter failed to start %s", err)
		}
		chans[i] = c
	}

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
