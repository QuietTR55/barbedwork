package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Permission struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type Role struct {
	ID          int           `json:"id"`
	Name        string        `json:"name"`
	Description *string       `json:"description"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Permissions []Permission  `json:"permissions"`
}

type WorkspaceUserRole struct {
	WorkspaceID pgtype.UUID `json:"workspace_id"`
	UserID      pgtype.UUID `json:"user_id"`
	RoleID      int         `json:"role_id"`
	AssignedAt  time.Time   `json:"assigned_at"`
}

// Request/Response DTOs
type CreateRoleRequest struct {
	Name          string  `json:"name" validate:"required,max=255"`
	Description   *string `json:"description" validate:"max=500"`
	PermissionIDs []int   `json:"permission_ids"`
}

type UpdateRoleRequest struct {
	Name          string  `json:"name" validate:"required,max=255"`
	Description   *string `json:"description" validate:"max=500"`
	PermissionIDs []int   `json:"permission_ids"`
}

type AssignRoleRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
	RoleID int    `json:"role_id" validate:"required"`
}

type RoleWithUserCount struct {
	Role      Role `json:"role"`
	UserCount int  `json:"user_count"`
}
