package main

import (
	"encoding/json"
	"log"

	"github.com/replaygaming/amplitude"
)

// Amplitude implements the adapter interface. It translates the events
// received and forwards them to the Amplitude HTTP_API.
type Amplitude struct {
	APIKey string
	events chan Event
}

// Start starts a new Amplitude server and prepares to receive incoming events.
func (a *Amplitude) Start() (chan<- Event, error) {
	var server *amplitude.Server
	server = amplitude.NewServer(a.APIKey)
	log.Println("[INFO AMPLITUDE] Starting Amplitude in production mode")
	a.events = make(chan Event)
	a.listen(server)
	return a.events, nil
}

// listen receives incoming events and translate them to Amplitude.
func (a *Amplitude) listen(server *amplitude.Server) {
	go func() {
		for e := range a.events {
			if e.Version != 1 {
				log.Printf("[WARN AMPLITUDE] Event Version not supported %v", e)
				continue
			}
			event := amplitude.Event{
				Type:   e.Type,
				UserID: e.UserID,
			}
			var err error
			switch e.Type {
			case "session_start", "session_end":
				//not implemented
			case "chips_purchase":
				prop := &ChipsPurchase{}
				if err := json.Unmarshal(e.Properties, prop); err != nil {
					log.Printf("[WARN AMPLITUDE] JSON conversion failed %s (%q)", err,
						&e.Properties)
					continue
				}
				event.Revenue = prop.Amount
				err = server.SendEvent(event)
			case "tournament_registration", "hand_played":
				if err := json.Unmarshal(e.Properties, &event.Properties); err != nil {
					log.Printf("[WARN AMPLITUDE] JSON conversion failed %s (%q)", err,
						&e.Properties)
					continue
				}
				err = server.SendEvent(event)
			default:
				log.Printf("[WARN AMPLITUDE] Unknown event type %s (%v)", e.Type, e)
			}
			if err != nil {
				log.Printf("[ERROR AMPLITUDE] %s", err)
			}
		}
	}()
}
