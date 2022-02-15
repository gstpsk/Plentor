package models

type Registration struct {
	Name  string `json:"name"`
	Date  string `json:"date"`  // yyyy-MM-dd
	From  string `json:"from"`  // HH:mm
	Until string `json:"until"` // HH:mm
}
