package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/replaygaming/gameanalytics"
)

type sessionEndProperties struct {
	Length int `json:"session_length"`
}

// GameAnalytics implements the adapter interface. It translates the events
// received and forwards them to the GameAnalytics SDK server.
type GameAnalytics struct {
	Environment string
	GameKey     string
	SecretKey   string
	events      chan Event
}

// Start starts a new GameAnalytics SDK server and prepares to receive incoming
// events.
func (a *GameAnalytics) Start() (chan<- Event, error) {
	var server *ga.Server
	if a.Environment == "production" {
		server = ga.NewServer(a.GameKey, a.SecretKey)
		log.Println("[INFO] Starting GameAnalytics in production mode")
	} else {
		server = ga.NewSandboxServer()
		log.Println("[INFO] Starting GameAnalytics in sandbox mode")
	}
	err := server.Start()
	if err != nil {
		return nil, fmt.Errorf("GameAnalytics (%v)", err)
	}
	a.events = make(chan Event)
	a.listen(server)
	return a.events, nil
}

// listen receives incoming events and translate them to GameAnalytics.
func (a *GameAnalytics) listen(server *ga.Server) {
	go func() {
		for e := range a.events {
			if e.Version != 1 {
				log.Printf("[WARN] Event Version not supported %v", e)
				continue
			}
			ts, err := a.getTimeDiff(e.Timestamp, server.TimestampOffset)
			if err != nil {
				log.Printf("[WARN] Event timestamp parsing failed %s (%e)", err, e)
				continue
			}
			shared := &ga.DefaultAnnotations{
				UserID:          e.UserID,
				SessionID:       e.Session.UUID,
				SessionNumber:   e.Session.Number,
				ClientTimestamp: ts,
			}
			shared.DefaultsRequiredValues()
			switch e.Type {
			case "session_start":
				user := ga.NewUserEvent(shared)
				server.SendEvent(user)
			case "session_end":
				ended := ga.NewSessionEndEvent(shared)
				s := &sessionEndProperties{}
				err := json.Unmarshal(e.Properties, s)
				if err != nil {
					log.Printf("[WARN] JSON conversion failed %s", err)
					continue
				}
				ended.Length = s.Length
				server.SendEvent(ended)
			default:
				log.Printf("[WARN] Unknown event type %s (%v)", e.Type, e)
			}
		}
	}()
}

func (a *GameAnalytics) getTimeDiff(ts string, offset int) (int, error) {
	timestamp, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return 0,
			fmt.Errorf("[WARN] Event timestamp parsing failed %s (%s)", ts, err)
	}
	now := int(timestamp.Unix())
	diff := now - offset
	return diff, nil
}
