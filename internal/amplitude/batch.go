package amplitude

import (
	"log"

	"github.com/replaygaming/amplitude"
)

const batchSize = 10

type batch struct {
	client *amplitude.Client
	events []amplitude.Event
}

func (b *batch) Send() {
	_, err := b.client.SendEvent(b.events...)
	if err != nil {
		log.Printf("[AMPLITUDE ERROR] %s", err)
	}
}

func (b *batch) Add(e amplitude.Event) {
	b.events = append(b.events, e)
}

func (b *batch) Full() bool {
	return len(b.events) == batchSize
}

func newBatch(c *amplitude.Client) *batch {
	return &batch{
		client: c,
		events: make([]amplitude.Event, batchSize),
	}
}
