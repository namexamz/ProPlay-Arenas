package services

import (
	"errors"
	"log/slog"
	"time"
	"venue-service/internal/models"
	"venue-service/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrVenueNotFound = errors.New("venue not found")
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

type ScheduleUpdate struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type VenueService interface {
	GetByID(id uint) (*models.Venue, error)
	GetList(filter VenueFilter) ([]models.Venue, error)
	GetByOwnerID(ownerID uint) ([]models.Venue, error)
	Create(venue *models.Venue) error
	Update(id uint, venue *models.Venue) error
	Delete(id uint) error
	GetSchedule(id uint) (*models.Venue, error)
	UpdateSchedule(id uint, schedule ScheduleUpdate) error
}

type venueService struct {
	repository repository.VenueRepository
	logger     *slog.Logger
}

func NewVenueService(repository repository.VenueRepository, logger *slog.Logger) VenueService {
	return &venueService{
		repository: repository,
		logger:     logger.With("layer", "service"),
	}
}

func (s *venueService) GetByID(id uint) (*models.Venue, error) {
	venue, err := s.repository.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVenueNotFound
		}
		s.logger.Error("Ошибка получения площадки по ID", "id", id, "error", err)
		return nil, err
	}

	return venue, nil
}

func (s *venueService) Create(v *models.Venue) error {
	if err := s.repository.Create(v); err != nil {
		s.logger.Error("Ошибка создания площадки", "venue_type", v.VenueType, "owner_id", v.OwnerID, "error", err)
		return err
	}
	return nil
}

func (s *venueService) GetList(filter VenueFilter) ([]models.Venue, error) {
	repoFilter := repository.VenueFilter{
		District:  filter.District,
		VenueType: filter.VenueType,
		HourPrice: filter.HourPrice,
		IsActive:  filter.IsActive,
		OwnerID:   filter.OwnerID,
		Page:      filter.Page,
		Limit:     filter.Limit,
	}
	venues, err := s.repository.GetList(repoFilter)
	if err != nil {
		s.logger.Error("Ошибка получения списка площадок", "error", err)
		return nil, err
	}
	return venues, nil
}

func (s *venueService) GetByOwnerID(ownerID uint) ([]models.Venue, error) {
	venues, err := s.repository.GetByOwnerID(ownerID)
	if err != nil {
		s.logger.Error("Ошибка получения площадок владельца", "owner_id", ownerID, "error", err)
		return nil, err
	}
	return venues, nil
}

func (s *venueService) Update(id uint, venue *models.Venue) error {
	// Проверяем существование площадки
	existingVenue, err := s.repository.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrVenueNotFound
		}
		s.logger.Error("Ошибка получения площадки для обновления", "id", id, "error", err)
		return err
	}

	existingVenue.VenueType = venue.VenueType
	existingVenue.OwnerID = venue.OwnerID
	existingVenue.IsActive = venue.IsActive
	existingVenue.HourPrice = venue.HourPrice
	existingVenue.District = venue.District
	existingVenue.StartTime = venue.StartTime
	existingVenue.EndTime = venue.EndTime
	existingVenue.Weekdays = venue.Weekdays

	if err := s.repository.Update(existingVenue); err != nil {
		s.logger.Error("Ошибка обновления площадки", "id", id, "error", err)
		return err
	}
	return nil
}

func (s *venueService) Delete(id uint) error {
	// Проверяем существование площадки
	_, err := s.repository.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrVenueNotFound
		}
		return err
	}

	if err := s.repository.Delete(id); err != nil {
		s.logger.Error("Ошибка деактивации площадки", "id", id, "error", err)
		return err
	}
	return nil
}

func (s *venueService) GetSchedule(id uint) (*models.Venue, error) {
	venue, err := s.repository.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrVenueNotFound
		}
		s.logger.Error("Ошибка получения расписания площадки", "id", id, "error", err)
		return nil, err
	}
	return venue, nil
}

func (s *venueService) UpdateSchedule(id uint, schedule ScheduleUpdate) error {
	venue, err := s.repository.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrVenueNotFound
		}
		return err
	}

	// Обновляем только время работы
	venue.StartTime = schedule.StartTime
	venue.EndTime = schedule.EndTime

	if err := s.repository.Update(venue); err != nil {
		s.logger.Error("Ошибка обновления расписания", "id", id, "error", err)
		return err
	}
	return nil
}
