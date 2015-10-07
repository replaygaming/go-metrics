package amplitude

import (
	"bytes"
	"testing"

	"github.com/replaygaming/go-metrics/internal/event"
)

func TestNewEvent_GenericEvent(t *testing.T) {
	events := []event.Event{
		{Type: "tournament_registration", UserID: "1", Properties: []byte(`{"GameID":"2"}`)},
		{Type: "hand_played", UserID: "1", Properties: []byte(`{"GameID":"2"}`)},
	}
	for _, in := range events {
		out, _ := newEvent(in)
		if out.EventType != in.Type {
			t.Errorf("Expected event type to be %s\ngot: %s", out.EventType, in.Type)
		}
		if out.UserID != in.UserID {
			t.Errorf("Expected user id to be %s\ngot: %s", out.UserID, in.UserID)
		}
		id, _ := out.EventProperties["GameID"]
		if bytes.Equal(*id, []byte("2")) {
			t.Errorf("Expected event property to be %s\ngot: %s", "2", *id)
		}
	}
}

func TestNewEvent_ChipsPurchase_SetRevenue(t *testing.T) {
	in := event.Event{
		Type:       "chips_purchase",
		Properties: []byte(`{"amount":9.99}`),
	}
	out, _ := newEvent(in)
	expected := 9.99
	result := out.Revenue
	if expected != result {
		t.Errorf("Expected event revenue to be %.2f\ngot: %.2f", expected, result)
	}
}

func TestNewEvent_ParsingError(t *testing.T) {
	events := []event.Event{
		{Type: "chips_purchase", Properties: []byte("")},
		{Type: "tournament_registration", Properties: []byte("")},
		{Type: "hand_played", Properties: []byte("")},
	}
	for _, in := range events {
		_, err := newEvent(in)
		if err == nil {
			t.Error("Expected parsing error, got none")
		}
	}
}

func TestNewEvent_NotImplementedError(t *testing.T) {
	events := []event.Event{
		{Type: "session_start", Properties: []byte("")},
		{Type: "session_end", Properties: []byte("")},
		{Type: "other", Properties: []byte("")},
	}
	for _, in := range events {
		_, err := newEvent(in)
		if err == nil {
			t.Error("Expected error, got none")
		}
	}
}
