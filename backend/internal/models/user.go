package models

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	PasswordHash string `json:"password_hash"`
}