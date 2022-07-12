package models

type User struct {
	UID       int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Sex       int64  `json:"sex"`
}
