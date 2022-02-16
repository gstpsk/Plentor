package models

type Registration struct {
	Id      string `json:"_id,omitempty" bson:"_id,omitempty"`
	EventId string `json:"event_id"`
	Name    string `json:"name"`
	Date    string `json:"date"`  // yyyy-MM-dd
	From    string `json:"from"`  // HH:mm
	Until   string `json:"until"` // HH:mm
}
