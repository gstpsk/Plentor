package models

type Event struct {
	UserId     string      `json:"user_id"`
	Name       string      `json:"name"`
	Location   string      `json:"location"`
	Timeblocks []Timeblock `json:"timeblocks"`
}
