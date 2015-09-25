package amplitude

import (
	"encoding/json"
	"fmt"

	"github.com/replaygaming/amplitude"
	"github.com/replaygaming/metrics/internal/event"
)

func newEvent(in event.Event) (amplitude.Event, error) {
	var err error
	out := amplitude.Event{
		EventType: in.Type,
		UserID:    in.UserID,
	}
	switch in.Type {
	case "chips_purchase":
		prop := &event.ChipsPurchase{}
		if err = json.Unmarshal(in.Properties, prop); err != nil {
			err = parsingError{jsonErr: err, blob: in.Properties}
		} else {
			out.Revenue = prop.Amount
		}
	case "tournament_registration", "hand_played":
		if err = json.Unmarshal(in.Properties, &out.EventProperties); err != nil {
			err = parsingError{jsonErr: err, blob: in.Properties}
		}
	case "session_start", "session_end":
		err = notImplementedError{}
	default:
		err = fmt.Errorf("Unknown event type %s (%v)", in.Type, in)
	}
	return out, err
}
