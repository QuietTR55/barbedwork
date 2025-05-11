package models

import (
	"database/sql"

	"github.com/google/uuid"
)

type Workspace struct {
	Id        uuid.UUID
	ImagePath sql.NullString
	Name      string
}

type WorkspaceFullData struct {
	Id        uuid.UUID
	ImagePath sql.NullString
	Name      string
	Users     []User
}