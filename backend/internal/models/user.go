package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	Id       pgtype.UUID `json:"id"`
	Username string `json:"username"`
	PasswordHash string `json:"password_hash"`
	ImagePath    pgtype.Text `json:"image_path"`
	Permissions  []string `json:"permissions"`
}