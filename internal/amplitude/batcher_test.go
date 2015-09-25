package amplitude

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/replaygaming/amplitude"
)

type delayedClient struct {
	wg     sync.WaitGroup
	events []amplitude.Event
}

func (c *delayedClient) SendEvent(e ...amplitude.Event) ([]byte, error) {
	for range e {
		c.wg.Done()
	}
	c.events = e
	return nil, errors.New("Err client")
}

func TestNewBatcher(t *testing.T) {
	c := &amplitude.NoopClient{}
	b := NewBatcher(c)
	if c != b.client {
		t.Error("Expected batcher client to be assigned")
	}
}

func TestBatcher_BatchLimit(t *testing.T) {
	e := amplitude.Event{EventType: "hand_played"}
	c := &delayedClient{}
	c.wg.Add(1)
	b := &TimeoutBatcher{
		client:    c,
		queueSize: 1,
		timeout:   1 * time.Nanosecond,
		in:        make(chan amplitude.Event),
	}
	b.start()
	b.Batch(e)
	c.wg.Wait()
	if len(c.events) != 1 {
		t.Errorf("Expected 1 event to be sent\ngot: %d", len(c.events))
	}
	if c.events[0].EventType != e.EventType {
		t.Errorf("Expected event to equal %v\ngot: %v", e, c.events[0])
	}
}

func TestBatcher_BatchTimeout(t *testing.T) {
	e := amplitude.Event{EventType: "hand_played"}
	c := &delayedClient{}
	c.wg.Add(1)
	b := &TimeoutBatcher{
		client:    c,
		queueSize: 2,
		timeout:   1 * time.Nanosecond,
		in:        make(chan amplitude.Event),
	}
	b.start()
	b.Batch(e)
	c.wg.Wait()
	if len(c.events) != 1 {
		t.Errorf("Expected 1 event to be sent\ngot: %d", len(c.events))
	}
	if c.events[0].EventType != e.EventType {
		t.Errorf("Expected event to equal %v\ngot: %v", e, c.events[0])
	}
}
