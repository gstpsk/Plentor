package models

type Event struct {
	Id         string      `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId     string      `json:"user_id"`
	Name       string      `json:"name"`
	Location   string      `json:"location"`
	Timeblocks []Timeblock `json:"timeblocks"`
}
