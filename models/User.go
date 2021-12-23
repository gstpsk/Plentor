package models

type User struct {
	Id    string `json:"_id, omitempty" bson:"_id,omitempty"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Hash  string `json:"hash"`
}
