package repository

import (
	"log/slog"
	"user-service/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Delete(id uint) error
	Update(user *models.User) error
	UpdateRole(userID uint, role models.Role) error
}

type userRepository struct {
	logger *slog.Logger
	db     *gorm.DB
}

func NewUserRepository(logger *slog.Logger, db *gorm.DB) UserRepository {
	return &userRepository{
		logger: logger.With(slog.String("layer", "repository"), slog.String("entity", "user")),
		db:     db,
	}
}

func (r *userRepository) Create(user *models.User) error {
	r.logger.Info("create user started", slog.String("email", user.Email))

	if err := r.db.Create(user).Error; err != nil {
		r.logger.Error("create user failed",
			slog.String("email", user.Email),
			slog.Any("error", err),
		)
		return err
	}

	r.logger.Info("create user completed",
		slog.Uint64("user_id", uint64(user.ID)),
	)
	return nil
}

func (r *userRepository) GetByID(id uint) (*models.User, error) {
	r.logger.Debug("get user by id started", slog.Uint64("user_id", uint64(id)))

	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		r.logger.Warn("get user by id failed",
			slog.Uint64("user_id", uint64(id)),
			slog.Any("error", err),
		)
		return nil, err
	}

	r.logger.Debug("get user by id completed", slog.Uint64("user_id", uint64(id)))
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	r.logger.Debug("get user by email started", slog.String("email", email))

	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		r.logger.Warn("get user by email failed",
			slog.String("email", email),
			slog.Any("error", err),
		)
		return nil,err
	}

	r.logger.Debug("get user by email completed",
		slog.Uint64("user_id", uint64(user.ID)),
		slog.String("email", email),
	)
	return &user, nil
}

func (r *userRepository) Delete(id uint) error {
	r.logger.Info("delete user started", slog.Uint64("user_id", uint64(id)))

	if err := r.db.Delete(models.User{}, id).Error; err != nil {
		r.logger.Error("delete user failed",
			slog.Uint64("user_id", uint64(id)),
			slog.Any("error", err),
		)
		return err
	}

	r.logger.Info("delete user completed", slog.Uint64("user_id", uint64(id)))
	return nil
}

func (r *userRepository) Update(user *models.User) error {
	r.logger.Info("update user started", slog.Uint64("user_id", uint64(user.ID)))

	if err := r.db.Save(user).Error; err != nil {
		r.logger.Error("update user failed",
			slog.Uint64("user_id", uint64(user.ID)),
			slog.Any("error", err),
		)
		return err
	}

	r.logger.Info("update user completed", slog.Uint64("user_id", uint64(user.ID)))
	return nil
}

func (r *userRepository) UpdateRole(id uint, role models.Role) error {
	r.logger.Info("update user role started",
		slog.Uint64("user_id", uint64(id)),
		slog.String("role", string(role)),
	)

	if err := r.db.Model(models.User{}).
		Where("id = ?", id).
		Update("role", role).Error; err != nil {

		r.logger.Error("update user role failed",
			slog.Uint64("user_id", uint64(id)),
			slog.String("role", string(role)),
			slog.Any("error", err),
		)
		return err
	}

	r.logger.Info("update user role completed",
		slog.Uint64("user_id", uint64(id)),
		slog.String("role", string(role)),
	)
	return nil
}
