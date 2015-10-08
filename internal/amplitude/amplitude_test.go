package amplitude

import (
	"testing"
	"time"
)

var testKey = "abcde"

func TestNewClient(t *testing.T) {
	a := NewClient(testKey)
	if a.batcher == nil {
		t.Error("Expected batcher to be initialized")
	}
	if a.events == nil {
		t.Error("Expected events channel to be initialized")
	}
}

func TestAmplitude_Start(t *testing.T) {
	e := []byte(`{"event_type":"test","user_id":"1234"}`)
	c := &delayedClient{}
	c.wg.Add(1)
	b := &TimeoutBatcher{
		client:    c,
		queueSize: 1,
		timeout:   1 * time.Nanosecond,
		in:        make(chan []byte),
	}
	b.start()
	a := &Amplitude{
		events:  make(chan []byte),
		batcher: b,
	}
	a.Start()
	a.events <- []byte("")
	a.events <- e
	c.wg.Wait()
}
