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
	// Используем мапу для явного указания полей, которые нужно обновить
	// Это позволяет обновлять поля в 0 или пустую строку
	updateData := map[string]interface{}{
		"venue_type":           venue.VenueType,
		"owner_id":             venue.OwnerID,
		"is_active":            venue.IsActive,
		"hour_price":           venue.HourPrice,
		"district":             venue.District,
		"monday_enabled":       venue.Weekdays.Monday.Enabled,
		"monday_start_time":    venue.Weekdays.Monday.StartTime,
		"monday_end_time":      venue.Weekdays.Monday.EndTime,
		"tuesday_enabled":      venue.Weekdays.Tuesday.Enabled,
		"tuesday_start_time":   venue.Weekdays.Tuesday.StartTime,
		"tuesday_end_time":     venue.Weekdays.Tuesday.EndTime,
		"wednesday_enabled":    venue.Weekdays.Wednesday.Enabled,
		"wednesday_start_time": venue.Weekdays.Wednesday.StartTime,
		"wednesday_end_time":   venue.Weekdays.Wednesday.EndTime,
		"thursday_enabled":     venue.Weekdays.Thursday.Enabled,
		"thursday_start_time":  venue.Weekdays.Thursday.StartTime,
		"thursday_end_time":    venue.Weekdays.Thursday.EndTime,
		"friday_enabled":       venue.Weekdays.Friday.Enabled,
		"friday_start_time":    venue.Weekdays.Friday.StartTime,
		"friday_end_time":      venue.Weekdays.Friday.EndTime,
		"saturday_enabled":     venue.Weekdays.Saturday.Enabled,
		"saturday_start_time":  venue.Weekdays.Saturday.StartTime,
		"saturday_end_time":    venue.Weekdays.Saturday.EndTime,
		"sunday_enabled":       venue.Weekdays.Sunday.Enabled,
		"sunday_start_time":    venue.Weekdays.Sunday.StartTime,
		"sunday_end_time":      venue.Weekdays.Sunday.EndTime,
	}

	if err := r.db.Model(venue).Updates(updateData).Error; err != nil {
		r.logger.Error("Ошибка обновления площадки", "id", venue.ID, "error", err)
		return err
	}
	return nil
}

func (r *venueRepository) Delete(id uint) error {
	// Используем стандартный soft delete GORM через db.Delete()
	// GORM автоматически проставит DeletedAt
	result := r.db.Delete(&models.Venue{}, id)
	if result.Error != nil {
		r.logger.Error("Ошибка удаления площадки", "id", id, "error", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
