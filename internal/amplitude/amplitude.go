package amplitude

import (
	"encoding/json"
	"log"

	"github.com/replaygaming/amplitude"
	"github.com/replaygaming/metrics/internal/event"
)

// Amplitude implements the adapter interface. It translates the events
// received and forwards them to the Amplitude HTTP_API.
type Amplitude struct {
	APIKey string
	events chan event.Event
}

// Start starts a new Amplitude server and prepares to receive incoming events.
func (a *Amplitude) Start() (chan<- event.Event, error) {
	var server *amplitude.Server
	server = amplitude.NewServer(a.APIKey)
	log.Println("[INFO AMPLITUDE] Starting Amplitude in production mode")
	a.events = make(chan event.Event)
	a.listen(server)
	return a.events, nil
}

// listen receives incoming events and translate them to Amplitude.
func (a *Amplitude) listen(server *amplitude.Server) {
	go func() {
		for in := range a.events {
			if in.Version != 1 {
				log.Printf("[WARN AMPLITUDE] Event Version not supported %v", in)
				continue
			}
			out := amplitude.Event{
				EventType: in.Type,
				UserID:    in.UserID,
			}
			var err error
			switch in.Type {
			case "session_start", "session_end":
				//not implemented
			case "chips_purchase":
				prop := &event.ChipsPurchase{}
				if err := json.Unmarshal(in.Properties, prop); err != nil {
					log.Printf("[WARN AMPLITUDE] JSON conversion failed %s (%q)", err,
						&in.Properties)
					continue
				}
				out.Revenue = prop.Amount
				_, err = server.SendEvent(out)
			case "tournament_registration", "hand_played":
				if err := json.Unmarshal(in.Properties, &out.EventProperties); err != nil {
					log.Printf("[WARN AMPLITUDE] JSON conversion failed %s (%q)", err,
						&in.Properties)
					continue
				}
				_, err = server.SendEvent(out)
			default:
				log.Printf("[WARN AMPLITUDE] Unknown event type %s (%v)", in.Type, in)
			}
			if err != nil {
				log.Printf("[ERROR AMPLITUDE] %s", err)
			}
		}
	}()
}
