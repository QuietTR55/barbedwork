package utilities

import (
	"context"
	"errors"
)

// Global variable to hold the user service instance
var userService UserServiceInterface

// Interface for user service to allow for dependency injection
type UserServiceInterface interface {
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)
}

// SetUserService sets the user service instance for the utilities package
func SetUserService(service UserServiceInterface) {
	userService = service
}

// GetUserPermissions retrieves all permissions for a user
func GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	if userService == nil {
		return nil, errors.New("user service not initialized")
	}
	return userService.GetUserPermissions(ctx, userID)
}

func CheckUserPermission(ctx context.Context, userID string, requiredPermission string) (bool, error) {
	permissions, err := GetUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, permission := range permissions {
		if permission == requiredPermission {
			return true, nil
		}
	}
	return false, nil
}