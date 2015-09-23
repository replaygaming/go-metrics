package event

import "encoding/json"

// Event contains the generic information received from Replay Poker
type Event struct {
	Type      string `json:"event"`
	UserID    string `json:"user_id"`
	Timestamp string
	Session   struct {
		UUID   string `json:"uuid"`
		Number uint
	}
	Properties json.RawMessage
}

// ChipsPurchase lists properties for 'chips_purchase' event sent by Replay
// Poker
type ChipsPurchase struct {
	Amount   float64
	Type     string
	Provider string
	Number   uint `json:"transaction_count"`
}
