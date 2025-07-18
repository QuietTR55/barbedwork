package repos

import (
	"backend/internal/models"
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepo struct {
	db *pgxpool.Pool
}

func NewRoleRepo(db *pgxpool.Pool) *RoleRepo {
	return &RoleRepo{db: db}
}

// GetAllPermissions retrieves all available permissions
func (r *RoleRepo) GetAllPermissions(ctx context.Context) ([]models.Permission, error) {
	query := `
		SELECT id, name, description 
		FROM permissions 
		ORDER BY name
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query permissions: %w", err)
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var p models.Permission
		if err := rows.Scan(&p.ID, &p.Name, &p.Description); err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, p)
	}

	return permissions, rows.Err()
}

// GetAllRoles retrieves all roles with their permissions
func (r *RoleRepo) GetAllRoles(ctx context.Context) ([]models.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at,
		       COALESCE(
		           json_agg(
		               json_build_object(
		                   'id', p.id,
		                   'name', p.name,
		                   'description', p.description
		               )
		           ) FILTER (WHERE p.id IS NOT NULL), 
		           '[]'::json
		       ) as permissions
		FROM roles r
		LEFT JOIN role_permissions rp ON r.id = rp.role_id
		LEFT JOIN permissions p ON rp.permission_id = p.id
		GROUP BY r.id, r.name, r.description, r.created_at, r.updated_at
		ORDER BY r.created_at
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query roles: %w", err)
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		var permissionsJSON []byte

		if err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
			&permissionsJSON,
		); err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}

		// Parse permissions JSON
		if len(permissionsJSON) > 0 && string(permissionsJSON) != "[]" {
			if err := json.Unmarshal(permissionsJSON, &role.Permissions); err != nil {
				return nil, fmt.Errorf("failed to parse permissions JSON: %w", err)
			}
		} else {
			role.Permissions = []models.Permission{}
		}

		roles = append(roles, role)
	}

	return roles, rows.Err()
}

// GetRoleByID retrieves a specific role with its permissions
func (r *RoleRepo) GetRoleByID(ctx context.Context, roleID int) (*models.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at,
		       COALESCE(
		           json_agg(
		               json_build_object(
		                   'id', p.id,
		                   'name', p.name,
		                   'description', p.description
		               )
		           ) FILTER (WHERE p.id IS NOT NULL), 
		           '[]'::json
		       ) as permissions
		FROM roles r
		LEFT JOIN role_permissions rp ON r.id = rp.role_id
		LEFT JOIN permissions p ON rp.permission_id = p.id
		WHERE r.id = $1
		GROUP BY r.id, r.name, r.description, r.created_at, r.updated_at
	`

	var role models.Role
	var permissionsJSON []byte

	err := r.db.QueryRow(ctx, query, roleID).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
		&permissionsJSON,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to query role: %w", err)
	}

	// Parse permissions JSON
	if len(permissionsJSON) > 0 && string(permissionsJSON) != "[]" {
		if err := json.Unmarshal(permissionsJSON, &role.Permissions); err != nil {
			return nil, fmt.Errorf("failed to parse permissions JSON: %w", err)
		}
	} else {
		role.Permissions = []models.Permission{}
	}

	return &role, nil
}

// CreateRole creates a new role with the specified permissions
func (r *RoleRepo) CreateRole(ctx context.Context, name string, description *string, permissionIDs []int) (*models.Role, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert role
	var roleID int
	query := `
		INSERT INTO roles (name, description) 
		VALUES ($1, $2) 
		RETURNING id
	`
	if err := tx.QueryRow(ctx, query, name, description).Scan(&roleID); err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	// Insert role permissions
	if len(permissionIDs) > 0 {
		permQuery := `INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2)`
		for _, permID := range permissionIDs {
			if _, err := tx.Exec(ctx, permQuery, roleID, permID); err != nil {
				return nil, fmt.Errorf("failed to assign permission %d to role: %w", permID, err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return the created role
	return r.GetRoleByID(ctx, roleID)
}

// UpdateRole updates an existing role
func (r *RoleRepo) UpdateRole(ctx context.Context, roleID int, name string, description *string, permissionIDs []int) (*models.Role, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Update role basic info
	query := `
		UPDATE roles 
		SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $3
	`
	if _, err := tx.Exec(ctx, query, name, description, roleID); err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	// Delete existing permissions
	if _, err := tx.Exec(ctx, "DELETE FROM role_permissions WHERE role_id = $1", roleID); err != nil {
		return nil, fmt.Errorf("failed to delete existing permissions: %w", err)
	}

	// Insert new permissions
	if len(permissionIDs) > 0 {
		permQuery := `INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2)`
		for _, permID := range permissionIDs {
			if _, err := tx.Exec(ctx, permQuery, roleID, permID); err != nil {
				return nil, fmt.Errorf("failed to assign permission %d to role: %w", permID, err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return the updated role
	return r.GetRoleByID(ctx, roleID)
}

// DeleteRole deletes a role and all its assignments
func (r *RoleRepo) DeleteRole(ctx context.Context, roleID int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Delete workspace user role assignments
	if _, err := tx.Exec(ctx, "DELETE FROM workspace_user_roles WHERE role_id = $1", roleID); err != nil {
		return fmt.Errorf("failed to delete role assignments: %w", err)
	}

	// Delete role permissions
	if _, err := tx.Exec(ctx, "DELETE FROM role_permissions WHERE role_id = $1", roleID); err != nil {
		return fmt.Errorf("failed to delete role permissions: %w", err)
	}

	// Delete role
	if _, err := tx.Exec(ctx, "DELETE FROM roles WHERE id = $1", roleID); err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	return tx.Commit(ctx)
}

// GetWorkspaceUserRoles retrieves all user role assignments for a workspace
func (r *RoleRepo) GetWorkspaceUserRoles(ctx context.Context, workspaceID string) ([]models.WorkspaceUserRole, error) {
	query := `
		SELECT workspace_id, user_id, role_id, assigned_at
		FROM workspace_user_roles
		WHERE workspace_id = $1
		ORDER BY assigned_at DESC
	`

	rows, err := r.db.Query(ctx, query, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query workspace user roles: %w", err)
	}
	defer rows.Close()

	var userRoles []models.WorkspaceUserRole
	for rows.Next() {
		var ur models.WorkspaceUserRole
		if err := rows.Scan(&ur.WorkspaceID, &ur.UserID, &ur.RoleID, &ur.AssignedAt); err != nil {
			return nil, fmt.Errorf("failed to scan workspace user role: %w", err)
		}
		userRoles = append(userRoles, ur)
	}

	return userRoles, rows.Err()
}

// AssignRoleToUser assigns a role to a user in a workspace
func (r *RoleRepo) AssignRoleToUser(ctx context.Context, workspaceID, userID string, roleID int) (*models.WorkspaceUserRole, error) {
	query := `
		INSERT INTO workspace_user_roles (workspace_id, user_id, role_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (workspace_id, user_id, role_id) DO NOTHING
		RETURNING workspace_id, user_id, role_id, assigned_at
	`

	var userRole models.WorkspaceUserRole
	err := r.db.QueryRow(ctx, query, workspaceID, userID, roleID).Scan(
		&userRole.WorkspaceID,
		&userRole.UserID,
		&userRole.RoleID,
		&userRole.AssignedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("role assignment already exists")
		}
		return nil, fmt.Errorf("failed to assign role to user: %w", err)
	}

	return &userRole, nil
}

// RemoveRoleFromUser removes a role assignment from a user in a workspace
func (r *RoleRepo) RemoveRoleFromUser(ctx context.Context, workspaceID, userID string, roleID int) error {
	query := `
		DELETE FROM workspace_user_roles 
		WHERE workspace_id = $1 AND user_id = $2 AND role_id = $3
	`

	result, err := r.db.Exec(ctx, query, workspaceID, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("role assignment not found")
	}

	return nil
}
