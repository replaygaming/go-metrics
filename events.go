package main

import "encoding/json"

// Event contains the generic information received from Replay Poker
type Event struct {
	Version   uint
	Type      string `json:"event"`
	UserID    string `json:"user_id"`
	Timestamp string
	Session   struct {
		UUID   string `json:"uuid"`
		Number uint
	}
	Properties json.RawMessage
}

// SessionEnd lists properties for 'session_end' event sent by Replay Poker
type SessionEnd struct {
	Length int `json:"session_length"`
}

// ChipsPurchase lists properties for 'chips_purchase' event sent by Replay
// Poker
type ChipsPurchase struct {
	Amount   float64
	Type     string
	Provider string
	Number   uint `json:"transaction_count"`
}

// TournamentRegistration lists properties for 'tournament_registration' event
// sent by Replay Poker
type TournamentRegistration struct {
	Type string
	Game string
}
