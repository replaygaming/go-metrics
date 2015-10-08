package amplitude

import (
	"bytes"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/replaygaming/amplitude"
)

type delayedClient struct {
	wg      sync.WaitGroup
	payload []byte
}

func TestNewQueue(t *testing.T) {
	q := NewQueue(1)
	if len(q) != 0 {
		t.Error("Expected queue to be empty")
	}
	if cap(q) != 1 {
		t.Error("Expected queue to be capped at 1")
	}
}

func TestQueue_Key(t *testing.T) {
	var q Queue
	if q.Key() != "events" {
		t.Error("Wrong key for queue payload")
	}
}

func TestQueue_Value(t *testing.T) {
	q := NewQueue(2)
	q = append(q, []byte(`{"even_type":"hand_played"}`))
	expected := []byte(`[{"even_type":"hand_played"}]`)
	result, _ := q.Value()

	if !bytes.Equal(expected, result) {
		t.Errorf("Expected single queue value to equal %q\ngot: %q", expected, result)
	}

	q = append(q, []byte(`{"even_type":"purchase"}`))
	expected = []byte(`[{"even_type":"hand_played"},{"even_type":"purchase"}]`)
	result, _ = q.Value()
	if !bytes.Equal(expected, result) {
		t.Errorf("Expected single queue value to equal %q\ngot: %q", expected, result)
	}
}

func (c *delayedClient) Send(p amplitude.Payload) ([]byte, error) {
	c.payload, _ = p.Value()
	if !bytes.Equal(c.payload, []byte("[]")) {
		c.wg.Done()
	}
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
	e := []byte(`{"even_type":"hand_played"}`)
	c := &delayedClient{}
	c.wg.Add(1)
	b := &TimeoutBatcher{
		client:    c,
		queueSize: 1,
		timeout:   1 * time.Nanosecond,
		in:        make(chan []byte),
	}
	b.start()
	b.Batch(e)
	c.wg.Wait()
	expected := []byte(`[{"even_type":"hand_played"}]`)
	result := c.payload
	if !bytes.Equal(expected, result) {
		t.Errorf("Expected payload to equal %q\ngot: %q", expected, result)
	}
}

func TestBatcher_BatchTimeout(t *testing.T) {
	e := []byte(`{"even_type":"hand_played"}`)
	c := &delayedClient{}
	c.wg.Add(1)
	b := &TimeoutBatcher{
		client:    c,
		queueSize: 2,
		timeout:   1 * time.Nanosecond,
		in:        make(chan []byte),
	}
	b.start()
	b.Batch(e)
	c.wg.Wait()
	expected := []byte(`[{"even_type":"hand_played"}]`)
	result := c.payload
	if !bytes.Equal(expected, result) {
		t.Errorf("Expected payload to equal %s\ngot: %s", expected, result)
	}
}
