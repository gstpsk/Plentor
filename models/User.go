package models

type User struct {
	_id   string `json:"id, omitempty"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Hash  string `json:"hash"`
}
