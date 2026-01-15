package service

import (
	"errors"
	"log/slog"

	"user-service/internal/dto"
	"user-service/internal/models"
	"user-service/internal/repository"
)

type UserService interface {
	GetMe(id uint) (*models.User, error)
	UpdateMe(id uint, req dto.UpdateUserRequest) (*models.User, error)
	BecomeOwner(id uint) (*models.User, error)
	GetPublicProfile(id uint) (*models.User, error)
}

type userService struct {
	logger     *slog.Logger
	repository repository.UserRepository
}

func NewUserService(
	logger *slog.Logger,
	repository repository.UserRepository,
) UserService {
	return &userService{
		logger:     logger,
		repository: repository,
	}
}

func (s *userService) GetMe(id uint) (*models.User, error) {
	user, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) UpdateMe(id uint, req dto.UpdateUserRequest) (*models.User, error) {
	user, err := s.repository.GetByID(id)

	if err != nil {
		return nil, err
	}
	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Email != nil {
		user.Email = *req.Email
	}

	if err := s.repository.Update(user); err != nil {
		return nil, err
	}
	return user, nil

}

func (s *userService) BecomeOwner(id uint) (*models.User, error) {
	user, err := s.repository.GetByID(id)

	if err != nil {
		return nil, err
	}

	if user.Role != models.RoleClient {
		return nil, errors.New("Владельцем может стать только клиент.")
	}

	user.Role = models.RoleOwner

	if err := s.repository.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetPublicProfile(id uint) (*models.User, error) {
	user, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	users := &models.User{
		FullName: user.FullName,
		Role:     user.Role,
	}
	return users, nil
}
