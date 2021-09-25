package models

type User struct {
	Model
	User         string `json:"user"`
	PasswordHash string `json:"password_hash"`
}
