package amplitude

import (
	"log"
	"os"

	"github.com/replaygaming/amplitude"
)

var logger = log.New(os.Stdout, "[AMPLITUDE] ", log.Lshortfile)

// Amplitude implements the adapter interface. It forwards events to the
// Amplitude HTTP API.
type Amplitude struct {
	client  amplitude.Client
	events  chan []byte
	batcher Batcher
}

// NewClient return a new Amplitude client with default values.
func NewClient(apiKey string) *Amplitude {
	c := amplitude.NewClient(apiKey)
	return &Amplitude{
		events:  make(chan []byte),
		batcher: NewBatcher(c),
	}
}

// Start starts a new Amplitude client and prepares to receive incoming events.
func (a *Amplitude) Start() (chan<- []byte, error) {
	logger.Println("[INFO] Starting Amplitude in production mode")
	b := a.batcher
	go func() {
		for e := range a.events {
			b.Batch(e)
		}
	}()
	return a.events, nil
}
