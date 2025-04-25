package services

import (
	"backend/internal/models"
	"backend/internal/repos"
)

type UserService struct {
	userRepo *repos.UserRepo
}

func NewUserService(userRepo *repos.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(username string, passwordHash string) error {
	return s.userRepo.CreateUser(username, passwordHash)
}

func (s *UserService) GetAllUsers() ([]*models.User, error) {
	return s.userRepo.GetAllUsers()
}
