package service

import (
	"errors"
	"log/slog"
	"user-service/internal/dto"
	"user-service/internal/models"
	"user-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
	db         *gorm.DB
	repository repository.UserRepository
}

func NewUserService(logger *slog.Logger, db *gorm.DB, repository repository.UserRepository) UserService {
	return &userService{logger: logger, db: db, repository: repository}
}

func (s *userService) Create(req dto.CreateUserRequest) (*models.User, error) {
	PasswordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, errors.New("Оибка при хеширование пароля")
	}

	role := models.RoleClient
	if req.Role != "" {
		role = models.Role(req.Role)
	}

	user := &models.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: PasswordHash,
		Role:     role,
	}

	if err := s.repository.Create(user); err != nil {
		return nil, errors.New("Ошибка при создании пользавателя")
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
		return nil, errors.New("Ошибка при добавление пользавателя по ID",)
	}
	return user, err
}

func (s *userService) GetByEmail(email string) (*models.User, error) {
	user, err := s.repository.GetByEmail(email)
	if err != nil {
		return nil, errors.New("Ошибка при создании пользавателя по Email")
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
	return user, err
}

func (s *userService) Delete(id uint) error {
	_, err := s.repository.GetByID(id)
	if err != nil {
		return err
	}

	return s.repository.Delete(id)
}
