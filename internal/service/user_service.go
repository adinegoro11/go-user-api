package service

import (
	"strings"

	"github.com/adinegoro11/go-user-api/internal/dto"
	"github.com/adinegoro11/go-user-api/internal/model"
	"github.com/adinegoro11/go-user-api/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) Me(userID uint) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *UserService) FindAll() ([]model.User, error) {
	return s.userRepo.FindAll()
}

func (s *UserService) FindByID(userID uint) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *UserService) Update(userID uint, req dto.UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	user.Name = strings.TrimSpace(req.Name)
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Delete(userID uint) error {
	return s.userRepo.Delete(userID)
}
