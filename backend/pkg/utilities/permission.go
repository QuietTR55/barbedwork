package utilities

import (
	"backend/internal/repos"
	"backend/internal/services"
	"context"
	"errors"
)


type PermissionChecker struct {
	userService *services.UserService
	userRepo    *repos.UserRepo
}

func NewPermissionChecker(userService *services.UserService, userRepo *repos.UserRepo) *PermissionChecker {
	return &PermissionChecker{
		userService: userService,
		userRepo:    userRepo,
	}
}

func (p *PermissionChecker) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	if p.userService == nil {
		return nil, errors.New("user service not initialized")
	}
	return p.userRepo.GetAllUserPermissions(ctx, userID)
}

func (p *PermissionChecker) CheckUserPermission(ctx context.Context, userID string, requiredPermission string) (bool, error) {
	permissions, err := p.GetUserPermissions(ctx, userID)
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