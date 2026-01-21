package repository

import (
	"log/slog"
	"venue-service/internal/models"

	"gorm.io/gorm"
)

type VenueFilter struct {
	District  string
	VenueType models.VenueType
	HourPrice int
	IsActive  *bool
	OwnerID   uint
	Page      int
	Limit     int
}

type VenueRepository interface {
	GetByID(id uint) (*models.Venue, error)
	GetList(filter VenueFilter) ([]models.Venue, error)
	GetByOwnerID(ownerID uint) ([]models.Venue, error)
	Create(venue *models.Venue) error
	Update(venue *models.Venue) error
	Delete(id uint) error
}

type venueRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewVenueRepository(db *gorm.DB, logger *slog.Logger) VenueRepository {
	return &venueRepository{
		db:     db,
		logger: logger.With("layer", "repository"),
	}
}

func (r *venueRepository) Create(venue *models.Venue) error {
	if err := r.db.Create(venue).Error; err != nil {
		r.logger.Error("Ошибка создания площадки", "error", err)
		return err
	}
	return nil
}

func (r *venueRepository) GetByID(id uint) (*models.Venue, error) {
	var venue models.Venue
	if err := r.db.First(&venue, id).Error; err != nil {
		r.logger.Error("Ошибка получения площадки по ID", "id", id, "error", err)
		return nil, err
	}
	return &venue, nil
}

func (r *venueRepository) GetList(filter VenueFilter) ([]models.Venue, error) {
	query := r.db.Model(&models.Venue{})
	if filter.District != "" {
		query = query.Where("district = ?", filter.District)
	}
	if filter.VenueType != "" {
		query = query.Where("venue_type = ?", filter.VenueType)
	}
	if filter.HourPrice > 0 {
		query = query.Where("hour_price = ?", filter.HourPrice)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	if filter.OwnerID > 0 {
		query = query.Where("owner_id = ?", filter.OwnerID)
	}

	// Пагинация
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
		if filter.Page > 0 {
			offset := (filter.Page - 1) * filter.Limit
			query = query.Offset(offset)
		}
	}

	// Сортировка по ID (новые сначала)
	query = query.Order("id DESC")

	var venues []models.Venue
	if err := query.Find(&venues).Error; err != nil {
		r.logger.Error("Ошибка получения списка площадок", "error", err)
		return nil, err
	}
	return venues, nil
}

func (r *venueRepository) GetByOwnerID(ownerID uint) ([]models.Venue, error) {
	var venues []models.Venue
	if err := r.db.Where("owner_id = ?", ownerID).Order("id DESC").Find(&venues).Error; err != nil {
		r.logger.Error("Ошибка получения площадок владельца", "owner_id", ownerID, "error", err)
		return nil, err
	}
	return venues, nil
}

func (r *venueRepository) Update(venue *models.Venue) error {
	if err := r.db.Model(venue).Updates(venue).Error; err != nil {
		r.logger.Error("Ошибка обновления площадки", "id", venue.ID, "error", err)
		return err
	}
	return nil
}

func (r *venueRepository) Delete(id uint) error {
	// Деактивация вместо удаления (soft delete)
	result := r.db.Model(&models.Venue{}).Where("id = ?", id).Update("is_active", false)
	if result.Error != nil {
		r.logger.Error("Ошибка деактивации площадки", "id", id, "error", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
