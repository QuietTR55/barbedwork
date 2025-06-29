package services

import (
	"backend/internal/models"
	"backend/internal/repos"
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
)

type UserService struct {
	userRepo *repos.UserRepo
}

func NewUserService(userRepo *repos.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(ctx context.Context, username string, passwordHash string) (*models.User, error) {
	//check name length
	if len(username) > 50 {
		return nil, errors.New("username must be less than 50 characters")
	}
	return s.userRepo.CreateUser(ctx, username, passwordHash)
}

func (s *UserService) Login(ctx context.Context, username string, password string) (*models.User, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		fmt.Println("Error getting user by username:", err)
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		fmt.Println("Password comparison failed:", err)
		return nil, ErrInvalidPassword
	}
	return user, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	return s.userRepo.GetAllUsers(ctx)
}

func (s *UserService) GetUser(ctx context.Context, userID string) (*models.User, error) {
	return s.userRepo.GetUserByID(ctx, userID)
}
