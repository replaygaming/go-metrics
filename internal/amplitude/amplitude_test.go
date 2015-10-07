package amplitude

import (
	"errors"
	"testing"
	"time"

	"github.com/replaygaming/amplitude"
	"github.com/replaygaming/go-metrics/internal/event"
)

var testKey = "abcde"

func TestParsingError_Error(t *testing.T) {
	err := parsingError{jsonErr: errors.New("Invalid character"), blob: []byte("P}")}
	expected := `JSON conversion failed Invalid character ("P}")`
	result := err.Error()
	if expected != result {
		t.Errorf("Expected error to be %s\ngot %s", expected, result)
	}
}

func TestNotImplementedError_Error(t *testing.T) {
	err := notImplementedError{eventType: "tournament_registration"}
	expected := "Event type tournament_registration not implemented"
	result := err.Error()
	if expected != result {
		t.Errorf("Expected error to be %s\ngot %s", expected, result)
	}
}

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
	e := event.Event{Type: "hand_played", Properties: []byte(`{"GameID":"2"}`)}
	c := &delayedClient{}
	c.wg.Add(1)
	b := &TimeoutBatcher{
		client:    c,
		queueSize: 1,
		timeout:   1 * time.Nanosecond,
		in:        make(chan amplitude.Event),
	}
	b.start()
	a := &Amplitude{
		events:  make(chan event.Event),
		batcher: b,
	}
	a.Start()
	a.events <- event.Event{}
	a.events <- e
	c.wg.Wait()
}
