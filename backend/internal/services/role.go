package services

import (
	"backend/internal/models"
	"backend/internal/repos"
	"context"
	"errors"
	"fmt"
)

var (
	ErrRoleNotFound     = errors.New("role not found")
	ErrInvalidRoleName  = errors.New("role name must be between 1 and 255 characters")
	ErrRoleNameExists   = errors.New("role name already exists")
	ErrInvalidUserID    = errors.New("invalid user ID")
	ErrInvalidRoleID    = errors.New("invalid role ID")
	ErrRoleAssignmentExists = errors.New("role assignment already exists")
)

type RoleService struct {
	roleRepo *repos.RoleRepo
	userRepo *repos.UserRepo
}

func NewRoleService(roleRepo *repos.RoleRepo, userRepo *repos.UserRepo) *RoleService {
	return &RoleService{
		roleRepo: roleRepo,
		userRepo: userRepo,
	}
}

// GetAllPermissions retrieves all available permissions
func (s *RoleService) GetAllPermissions(ctx context.Context) ([]models.Permission, error) {
	return s.roleRepo.GetAllPermissions(ctx)
}

// GetAllRoles retrieves all roles with their permissions
func (s *RoleService) GetAllRoles(ctx context.Context) ([]models.Role, error) {
	return s.roleRepo.GetAllRoles(ctx)
}

// GetRoleByID retrieves a specific role by ID
func (s *RoleService) GetRoleByID(ctx context.Context, roleID int) (*models.Role, error) {
	if roleID <= 0 {
		return nil, ErrInvalidRoleID
	}

	role, err := s.roleRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		if err.Error() == "role not found" {
			return nil, ErrRoleNotFound
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return role, nil
}

// CreateRole creates a new role with validation
func (s *RoleService) CreateRole(ctx context.Context, req models.CreateRoleRequest) (*models.Role, error) {
	// Validate input
	if err := s.validateRoleName(req.Name); err != nil {
		return nil, err
	}

	if req.Description != nil && len(*req.Description) > 500 {
		return nil, errors.New("description must be less than 500 characters")
	}

	// Validate permission IDs
	if err := s.validatePermissionIDs(ctx, req.PermissionIDs); err != nil {
		return nil, err
	}

	// Create role
	role, err := s.roleRepo.CreateRole(ctx, req.Name, req.Description, req.PermissionIDs)
	if err != nil {
		if isUniqueConstraintError(err) {
			return nil, ErrRoleNameExists
		}
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return role, nil
}

// UpdateRole updates an existing role
func (s *RoleService) UpdateRole(ctx context.Context, roleID int, req models.UpdateRoleRequest) (*models.Role, error) {
	// Validate role ID
	if roleID <= 0 {
		return nil, ErrInvalidRoleID
	}

	// Validate input
	if err := s.validateRoleName(req.Name); err != nil {
		return nil, err
	}

	if req.Description != nil && len(*req.Description) > 500 {
		return nil, errors.New("description must be less than 500 characters")
	}

	// Validate permission IDs
	if err := s.validatePermissionIDs(ctx, req.PermissionIDs); err != nil {
		return nil, err
	}

	// Check if role exists
	_, err := s.roleRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		if err.Error() == "role not found" {
			return nil, ErrRoleNotFound
		}
		return nil, fmt.Errorf("failed to check role existence: %w", err)
	}

	// Update role
	role, err := s.roleRepo.UpdateRole(ctx, roleID, req.Name, req.Description, req.PermissionIDs)
	if err != nil {
		if isUniqueConstraintError(err) {
			return nil, ErrRoleNameExists
		}
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	return role, nil
}

// DeleteRole deletes a role and all its assignments
func (s *RoleService) DeleteRole(ctx context.Context, roleID int) error {
	if roleID <= 0 {
		return ErrInvalidRoleID
	}

	// Check if role exists
	_, err := s.roleRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		if err.Error() == "role not found" {
			return ErrRoleNotFound
		}
		return fmt.Errorf("failed to check role existence: %w", err)
	}

	// TODO: Add check to prevent deletion of system roles (Owner, Admin, etc.)
	// This would depend on your business logic

	return s.roleRepo.DeleteRole(ctx, roleID)
}

// GetWorkspaceUserRoles retrieves all user role assignments for a workspace
func (s *RoleService) GetWorkspaceUserRoles(ctx context.Context, workspaceID string) ([]models.WorkspaceUserRole, error) {
	if workspaceID == "" {
		return nil, errors.New("workspace ID is required")
	}

	return s.roleRepo.GetWorkspaceUserRoles(ctx, workspaceID)
}

// AssignRoleToUser assigns a role to a user in a workspace
func (s *RoleService) AssignRoleToUser(ctx context.Context, workspaceID, userID string, roleID int) (*models.WorkspaceUserRole, error) {
	// Validate inputs
	if workspaceID == "" {
		return nil, errors.New("workspace ID is required")
	}
	if userID == "" {
		return nil, ErrInvalidUserID
	}
	if roleID <= 0 {
		return nil, ErrInvalidRoleID
	}

	// Verify user exists
	_, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Verify role exists
	_, err = s.roleRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		if err.Error() == "role not found" {
			return nil, ErrRoleNotFound
		}
		return nil, fmt.Errorf("failed to verify role: %w", err)
	}

	// Assign role
	userRole, err := s.roleRepo.AssignRoleToUser(ctx, workspaceID, userID, roleID)
	if err != nil {
		if err.Error() == "role assignment already exists" {
			return nil, ErrRoleAssignmentExists
		}
		return nil, fmt.Errorf("failed to assign role: %w", err)
	}

	return userRole, nil
}

// RemoveRoleFromUser removes a role assignment from a user in a workspace
func (s *RoleService) RemoveRoleFromUser(ctx context.Context, workspaceID, userID string, roleID int) error {
	// Validate inputs
	if workspaceID == "" {
		return errors.New("workspace ID is required")
	}
	if userID == "" {
		return ErrInvalidUserID
	}
	if roleID <= 0 {
		return ErrInvalidRoleID
	}

	err := s.roleRepo.RemoveRoleFromUser(ctx, workspaceID, userID, roleID)
	if err != nil {
		if err.Error() == "role assignment not found" {
			return errors.New("role assignment not found")
		}
		return fmt.Errorf("failed to remove role assignment: %w", err)
	}

	return nil
}

// Helper functions

func (s *RoleService) validateRoleName(name string) error {
	if len(name) == 0 || len(name) > 255 {
		return ErrInvalidRoleName
	}
	return nil
}

func (s *RoleService) validatePermissionIDs(ctx context.Context, permissionIDs []int) error {
	if len(permissionIDs) == 0 {
		return nil // Empty permissions are allowed
	}

	// Get all available permissions to validate IDs
	availablePermissions, err := s.roleRepo.GetAllPermissions(ctx)
	if err != nil {
		return fmt.Errorf("failed to get available permissions: %w", err)
	}

	// Create a map for quick lookup
	validIDs := make(map[int]bool)
	for _, perm := range availablePermissions {
		validIDs[perm.ID] = true
	}

	// Validate each permission ID
	for _, id := range permissionIDs {
		if !validIDs[id] {
			return fmt.Errorf("invalid permission ID: %d", id)
		}
	}

	return nil
}

func isUniqueConstraintError(err error) bool {
	// This is a simple check - you might want to use a more robust method
	// based on your database driver's error types
	return err != nil && (
		err.Error() == "UNIQUE constraint failed" ||
		err.Error() == "duplicate key value violates unique constraint")
}
