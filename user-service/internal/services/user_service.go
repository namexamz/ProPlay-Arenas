package service

import (
	"errors"
	"fmt"
	"log/slog"

	"user-service/internal/dto"
	"user-service/internal/models"
	"user-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrInvalidRole   = errors.New("invalid role")
	ErrEmptyUpdate   = errors.New("no fields to update")
)


type UserService interface {
	Create(req dto.CreateUserRequest) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(id uint, dto *dto.UpdateUserRequest) (*models.User, error)
	Delete(id uint) error
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


func (s *userService) Create(req dto.CreateUserRequest) (*models.User, error) {
	passwordHash, err := hashPassword(req.Password)
	if err != nil {

		s.logger.Error("password hashing failed", slog.Any("error", err))
		return nil, fmt.Errorf("password hashing failed: %w", err)
	}

	role := models.RoleClient
	if req.Role != "" {
		switch models.Role(req.Role) {
		case models.RoleAdmin, models.RoleClient:
			role = models.Role(req.Role)
		default:
			return nil, ErrInvalidRole
		}
	}

	user := &models.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: passwordHash,
		Role:     role,
	}

	if err := s.repository.Create(user); err != nil {

		s.logger.Error("failed to create user", slog.Any("error", err))
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}


func (s *userService) GetByID(id uint) (*models.User, error) {
	user, err := s.repository.GetByID(id)
	if err != nil {
		// КРИТИЧНО:
		// нужно различать not found и реальную ошибку БД
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		s.logger.Error("failed to get user by id", slog.Any("error", err))
		return nil, err
	}

	return user, nil
}


func (s *userService) GetByEmail(email string) (*models.User, error) {
	user, err := s.repository.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		s.logger.Error("failed to get user by email", slog.Any("error", err))
		return nil, err
	}

	return user, nil
}

func (s *userService) Update(id uint, dto *dto.UpdateUserRequest) (*models.User, error) {
	if dto.FullName == nil && dto.Email == nil {
		return nil, ErrEmptyUpdate
	}

	user, err := s.repository.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if dto.FullName != nil {
		user.FullName = *dto.FullName
	}
	if dto.Email != nil {
		user.Email = *dto.Email
	}

	if err := s.repository.Update(user); err != nil {
		s.logger.Error("failed to update user", slog.Any("error", err))
		return nil, err
	}

	return user, nil
}

func (s *userService) Delete(id uint) error {
	err := s.repository.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}

		s.logger.Error("failed to delete user", slog.Any("error", err))
		return err
	}

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
