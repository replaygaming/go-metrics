package amplitude

import (
	"fmt"
	"log"
	"os"

	"github.com/replaygaming/amplitude"
	"github.com/replaygaming/go-metrics/internal/event"
)

var logger = log.New(os.Stdout, "[AMPLITUDE] ", log.Lshortfile)

type parsingError struct {
	jsonErr error
	blob    []byte
}

func (e parsingError) Error() string {
	return fmt.Sprintf("JSON conversion failed %s (%q)", e.jsonErr, e.blob)
}

type notImplementedError struct {
	eventType string
}

func (e notImplementedError) Error() string {
	return fmt.Sprintf("Event type %s not implemented", e.eventType)
}

// Amplitude implements the adapter interface. It translates the events
// received and forwards them to the Amplitude HTTP API.
type Amplitude struct {
	client  amplitude.Client
	events  chan event.Event
	batcher Batcher
}

// NewClient return a new Amplitude client with default values.
func NewClient(apiKey string) *Amplitude {
	c := amplitude.NewClient(apiKey)
	return &Amplitude{
		events:  make(chan event.Event),
		batcher: NewBatcher(c),
	}
}

// Start starts a new Amplitude client and prepares to receive incoming events.
func (a *Amplitude) Start() (chan<- event.Event, error) {
	logger.Println("[INFO] Starting Amplitude in production mode")
	b := a.batcher
	go func() {
		for e := range a.events {
			event, err := newEvent(e)
			if err != nil {
				if _, ok := err.(notImplementedError); !ok {
					logger.Printf("[WARN] %s", err)
				}
				continue
			}
			b.Batch(event)
		}
	}()
	return a.events, nil
}
