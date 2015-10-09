package amplitude

import (
	"bytes"
	"time"

	"github.com/replaygaming/amplitude"
)

// Queue implents the amplitude.Payload interface.
type Queue [][]byte

// Encode combine all payload received into a json array.
func (q Queue) Encode() ([]byte, error) {
	vals := [][]byte{
		[]byte("["),
		bytes.Join(q, []byte(",")),
		[]byte("]"),
	}
	return bytes.Join(vals, nil), nil
}

// NewQueue creates a new queue with the size passed.
func NewQueue(size int) Queue {
	return make(Queue, 0, size)
}

// Batcher interface
type Batcher interface {
	Batch([]byte)
}

// TimeoutBatcher implements the batcher interface using a fixed queue size and
// a timeout limit. When either is reached, the batch is sent to the client.
type TimeoutBatcher struct {
	client    amplitude.Client
	in        chan []byte
	queueSize int
	timeout   time.Duration
}

// Send calls the client Send immediately with the events already in the batch.
func (b *TimeoutBatcher) send(q Queue) {
	if _, err := b.client.Send(q); err != nil {
		logger.Printf("[ERROR] %s", err)
	}
}

// Batch adds the event to the queue and automatically sends it once the batch
// in full or the timeout is reached.
func (b *TimeoutBatcher) Batch(e []byte) {
	go func() {
		b.in <- e
	}()
}

func (b *TimeoutBatcher) start() {
	size := b.queueSize
	queue := NewQueue(size)
	tick := time.Tick(b.timeout)
	go func() {
		for {
			select {
			case e := <-b.in:
				queue = append(queue, e)
				if len(queue) == b.queueSize {
					b.send(queue)
					queue = NewQueue(size)
				}
			case <-tick:
				if len(queue) > 0 {
					b.send(queue)
					queue = NewQueue(size)
				}
			}
		}
	}()
}

// NewBatcher returns a TimeoutBatcher with a timeout of 5 seconds and a batch
// size of 10 events.
func NewBatcher(c amplitude.Client) *TimeoutBatcher {
	b := &TimeoutBatcher{
		client:    c,
		queueSize: 10,
		timeout:   5 * time.Second,
		in:        make(chan []byte),
	}
	b.start()
	return b
}
