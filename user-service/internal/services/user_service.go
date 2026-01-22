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
	GetPublicProfile(id uint) (*dto.UserResponse, error)
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
		logger: logger.With(
			slog.String("layer", "service"),
			slog.String("service", "user"),
		),
		repository: repository,
	}
}

func (s *userService) GetMe(id uint) (*models.User, error) {
	s.logger.Debug("get current user started",
		slog.Uint64("user_id", uint64(id)),
	)

	user, err := s.repository.GetByID(id)
	if err != nil {
		s.logger.Warn("get current user failed",
			slog.Uint64("user_id", uint64(id)),
			slog.Any("error", err),
		)
		return nil, err
	}

	s.logger.Debug("get current user completed",
		slog.Uint64("user_id", uint64(id)),
	)
	return user, nil
}

func (s *userService) UpdateMe(id uint, req dto.UpdateUserRequest) (*models.User, error) {
	s.logger.Info("update current user started",
		slog.Uint64("user_id", uint64(id)),
	)

	user, err := s.repository.GetByID(id)
	if err != nil {
		s.logger.Warn("update current user failed: user not found",
			slog.Uint64("user_id", uint64(id)),
			slog.Any("error", err),
		)
		return nil, err
	}

	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Email != nil {
		user.Email = *req.Email
	}

	if err := s.repository.Update(user); err != nil {
		s.logger.Error("update current user failed",
			slog.Uint64("user_id", uint64(id)),
			slog.Any("error", err),
		)
		return nil, err
	}

	s.logger.Info("update current user completed",
		slog.Uint64("user_id", uint64(id)),
	)
	return user, nil
}

func (s *userService) BecomeOwner(id uint) (*models.User, error) {
	s.logger.Info("become owner started",
		slog.Uint64("user_id", uint64(id)),
	)

	user, err := s.repository.GetByID(id)
	if err != nil {
		s.logger.Warn("become owner failed: user not found",
			slog.Uint64("user_id", uint64(id)),
			slog.Any("error", err),
		)
		return nil, err
	}

	if user.Role != models.RoleClient {
		s.logger.Warn("become owner forbidden: invalid role",
			slog.Uint64("user_id", uint64(id)),
			slog.String("current_role", string(user.Role)),
		)
		return nil, errors.New("владельцем может стать только клиент")
	}

	user.Role = models.RoleOwner

	if err := s.repository.Update(user); err != nil {
		s.logger.Error("become owner failed",
			slog.Uint64("user_id", uint64(id)),
			slog.Any("error", err),
		)
		return nil, err
	}

	s.logger.Info("become owner completed",
		slog.Uint64("user_id", uint64(id)),
		slog.String("new_role", string(user.Role)),
	)
	return user, nil
}
func (s *userService) GetPublicProfile(id uint) (*dto.UserResponse, error) {
    user, err := s.repository.GetByID(id)
    if err != nil {
        return nil, err
    }

    // формируем DTO
    return &dto.UserResponse{
        ID:       user.ID,           // здесь уже работает
        FullName: user.FullName,
        Email:    "",                // скрываем email для публичного профиля
        Role:     string(user.Role),
    }, nil
}