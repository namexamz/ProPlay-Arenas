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
	return &userRepository{logger: logger, db: db}
}

func (r *userRepository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}
func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User

	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User

	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *userRepository) Delete(id uint) error {
	if err := r.db.Delete(models.User{},id).Error; err != nil {
		return err
	}
	return nil
}
func (r *userRepository) Update(user *models.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return err
	}
	return nil
}
func(r *userRepository)UpdateRole(id uint, role models.Role) error{
	return r.db.Model(models.User{}).Where("id = ?",id).Update("role",role).Error
}
