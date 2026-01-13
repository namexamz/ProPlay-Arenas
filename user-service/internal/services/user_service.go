package service

import (
	"errors"
	"fmt"
	"log/slog"
	"user-service/internal/dto"
	"user-service/internal/models"
	"user-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
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

func NewUserService(logger *slog.Logger, repository repository.UserRepository) UserService {
	return &userService{logger: logger, repository: repository}
}

func (s *userService) Create(req dto.CreateUserRequest) (*models.User, error) {
	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("ошибка при хешировании пароля: %w", err)
	}

	role := models.RoleClient
	if req.Role != "" {
		role = models.Role(req.Role)
	}

	user := &models.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: passwordHash,
		Role:     role,
	}

	if err := s.repository.Create(user); err != nil {
		return nil ,fmt.Errorf("ошибка при создании пользователя: %w", err)
	}
	return user, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *userService) GetByID(id uint) (*models.User, error) {
	user, err := s.repository.GetByID(id)
	if err != nil {
		return nil, errors.New("Ошибка при получении пользавателя по ID")
	}
	return user, err
}

func (s *userService) GetByEmail(email string) (*models.User, error) {
	user, err := s.repository.GetByEmail(email)
	if err != nil {
		return nil, errors.New("Ошибка при получении пользавателя по Email")
	}
	return user, err
}

func (s *userService) Update(id uint, dto *dto.UpdateUserRequest) (*models.User, error) {
	user, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	if dto.FullName != nil {
		user.FullName = *dto.FullName
	}
	if dto.Email != nil {
		user.Email = *dto.Email
	}

	if err := s.repository.Update(user); err != nil {
		return nil, err
	}
	return user,nil
}

func (s *userService) Delete(id uint) error {
	_, err := s.repository.GetByID(id)
	if err != nil {
		return err
	}

	return s.repository.Delete(id)
}
