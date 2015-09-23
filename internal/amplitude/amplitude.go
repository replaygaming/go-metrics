package amplitude

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/replaygaming/amplitude"
	"github.com/replaygaming/metrics/internal/event"
)

// Amplitude implements the adapter interface. It translates the events
// received and forwards them to the Amplitude HTTP API.
type Amplitude struct {
	APIKey string
	events chan event.Event
}

// Start starts a new Amplitude client and prepares to receive incoming events.
func (a *Amplitude) Start() (chan<- event.Event, error) {
	var client *amplitude.Client
	client = amplitude.NewClient(a.APIKey)
	log.Println("[INFO AMPLITUDE] Starting Amplitude in production mode")
	a.events = make(chan event.Event)
	a.listen(client)
	return a.events, nil
}

// listen receives incoming events and translate them to Amplitude.
func (a *Amplitude) listen(client *amplitude.Client) {
	go func() {
		b := newBatch(client)
		tick := time.Tick(5 * time.Second)
		for {
			select {
			case e := <-a.events:
				event, err := newEvent(e)
				if err != nil {
					log.Printf("[AMPLITUDE WARN] %s", err)
					continue
				}
				b.Add(event)
				if b.Full() {
					go b.Send()
					b = newBatch(client)
				}
			case <-tick:
				go b.Send()
				b = newBatch(client)
			}
		}
	}()
}

func newEvent(in event.Event) (amplitude.Event, error) {
	out := amplitude.Event{
		EventType: in.Type,
		UserID:    in.UserID,
	}
	switch in.Type {
	case "chips_purchase":
		prop := &event.ChipsPurchase{}
		if err := json.Unmarshal(in.Properties, prop); err != nil {
			err = fmt.Errorf("JSON conversion failed %s (%q)", err, &in.Properties)
			return out, err
		}
		out.Revenue = prop.Amount
	case "tournament_registration", "hand_played":
		if err := json.Unmarshal(in.Properties, &out.EventProperties); err != nil {
			err = fmt.Errorf("JSON conversion failed %s (%q)", err, &in.Properties)
			return out, err
		}
	case "session_start", "session_end":
		//not implemented
	default:
		err := fmt.Errorf("Unknown event type %s (%v)", in.Type, in)
		return out, err
	}
	return out, nil
}
