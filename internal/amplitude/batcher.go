package amplitude

import (
	"time"

	"github.com/replaygaming/amplitude"
)

// Batcher interface
type Batcher interface {
	Batch(amplitude.Event)
}

// TimeoutBatcher implements the batcher interface using a fixed queue size and
// a timeout limit. When either is reached, the batch is sent to the client.
type TimeoutBatcher struct {
	client    amplitude.Client
	events    []amplitude.Event
	in        chan (amplitude.Event)
	queueSize int
	timeout   time.Duration
}

// Send calls the client SendEvent immediately with the events already in the
// batch.
func (b *TimeoutBatcher) send(e []amplitude.Event) {
	if _, err := b.client.SendEvent(e...); err != nil {
		logger.Printf("[ERROR] %s", err)
	}
}

// Batch adds the event to the queue and automatically sends it once the batch
// in full or the timeout is reached.
func (b *TimeoutBatcher) Batch(e amplitude.Event) {
	go func() {
		b.in <- e
	}()
}

func (b *TimeoutBatcher) start() {
	queue := b.newQueue()
	tick := time.Tick(b.timeout)
	go func() {
		for {
			select {
			case e := <-b.in:
				queue = append(queue, e)
				if len(queue) == b.queueSize {
					b.send(queue)
					queue = b.newQueue()
				}
			case <-tick:
				if len(queue) > 0 {
					b.send(queue)
					queue = b.newQueue()
				}
			}
		}
	}()
}

func (b *TimeoutBatcher) newQueue() []amplitude.Event {
	return make([]amplitude.Event, 0, b.queueSize)
}

// NewBatcher returns a TimeoutBatcher with a timeout of 5 seconds and a batch
// size of 10 events.
func NewBatcher(c amplitude.Client) *TimeoutBatcher {
	b := &TimeoutBatcher{
		client:    c,
		queueSize: 10,
		timeout:   5 * time.Second,
		in:        make(chan amplitude.Event),
	}
	b.start()
	return b
}
