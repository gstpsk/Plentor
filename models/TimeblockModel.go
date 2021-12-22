package models

type Timeblock struct {
	Date  string `json:"date"`
	From  string `json:"from"`
	Until string `json:"until"`
}
