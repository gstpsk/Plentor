package models

type Event struct {
	Name       string      `json:"name"`
	Location   string      `json:"location"`
	Timeblocks []Timeblock `json:"timeblocks"`
}
