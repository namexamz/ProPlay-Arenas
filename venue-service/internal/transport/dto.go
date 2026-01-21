package transport

import (
	"fmt"
	"venue-service/internal/models"
	"venue-service/internal/services"
	"venue-service/internal/validation"
)

// VenueDTO - DTO для запросов (Create/Update) и ответов
// Время представлено в формате "HH:MM" как строки
type VenueDTO struct {
	ID        uint             `json:"id,omitempty"` // Только в ответах
	VenueType models.VenueType `json:"venue_type"`
	OwnerID   uint             `json:"owner_id"`
	IsActive  bool             `json:"is_active"`
	HourPrice int              `json:"hour_price"`
	District  string           `json:"district"`
	StartTime string           `json:"start_time"` // Формат "HH:MM"
	EndTime   string           `json:"end_time"`   // Формат "HH:MM"
	Weekdays  models.Weekdays  `json:"weekdays"`
}

// ScheduleDTO - DTO для расписания работы площадки (ответ)
type ScheduleDTO struct {
	StartTime string `json:"start_time"` // Формат "HH:MM"
	EndTime   string `json:"end_time"`   // Формат "HH:MM"
}

// ScheduleUpdateDTO - DTO для обновления расписания (запрос)
type ScheduleUpdateDTO struct {
	StartTime string `json:"start_time" binding:"required"` // Формат "HH:MM"
	EndTime   string `json:"end_time" binding:"required"`   // Формат "HH:MM"
}

// ToVenueDTO конвертирует модель Venue в DTO (для ответов)
func ToVenueDTO(venue *models.Venue) VenueDTO {
	return VenueDTO{
		ID:        venue.ID,
		VenueType: venue.VenueType,
		OwnerID:   venue.OwnerID,
		IsActive:  venue.IsActive,
		HourPrice: venue.HourPrice,
		District:  venue.District,
		StartTime: venue.StartTime.Format("15:04"),
		EndTime:   venue.EndTime.Format("15:04"),
		Weekdays:  venue.Weekdays,
	}
}

// ToVenueDTOList конвертирует список моделей Venue в список DTO
func ToVenueDTOList(venues []models.Venue) []VenueDTO {
	dtoList := make([]VenueDTO, len(venues))
	for i := range venues {
		dtoList[i] = ToVenueDTO(&venues[i])
	}
	return dtoList
}

// FromVenueDTO конвертирует DTO в модель Venue
// Возвращает ошибку, если время имеет неверный формат
func FromVenueDTO(dto *VenueDTO) (*models.Venue, error) {
	// Парсим время начала
	startTime, err := validation.ValidateTime(dto.StartTime)
	if err != nil {
		return nil, fmt.Errorf("неверный формат start_time: %w", err)
	}

	// Парсим время окончания
	endTime, err := validation.ValidateTime(dto.EndTime)
	if err != nil {
		return nil, fmt.Errorf("неверный формат end_time: %w", err)
	}

	venue := &models.Venue{
		VenueType: dto.VenueType,
		OwnerID:   dto.OwnerID,
		IsActive:  dto.IsActive,
		HourPrice: dto.HourPrice,
		District:  dto.District,
		StartTime: startTime,
		EndTime:   endTime,
		Weekdays:  dto.Weekdays,
	}

	// Если есть ID (для обновления), устанавливаем его
	if dto.ID != 0 {
		venue.ID = dto.ID
	}

	return venue, nil
}

// ToScheduleDTO конвертирует модель Venue в ScheduleDTO
func ToScheduleDTO(venue *models.Venue) ScheduleDTO {
	return ScheduleDTO{
		StartTime: venue.StartTime.Format("15:04"),
		EndTime:   venue.EndTime.Format("15:04"),
	}
}

// FromScheduleUpdateDTO конвертирует ScheduleUpdateDTO в services.ScheduleUpdate
func FromScheduleUpdateDTO(dto *ScheduleUpdateDTO) (*services.ScheduleUpdate, error) {
	startTime, err := validation.ValidateTime(dto.StartTime)
	if err != nil {
		return nil, fmt.Errorf("неверный формат start_time: %w", err)
	}

	endTime, err := validation.ValidateTime(dto.EndTime)
	if err != nil {
		return nil, fmt.Errorf("неверный формат end_time: %w", err)
	}

	return &services.ScheduleUpdate{
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}
