package transport

import (
	"fmt"
	"venue-service/internal/models"
	"venue-service/internal/validation"
)

// DayScheduleDTO - DTO для расписания одного дня недели
// Время представлено в формате "HH:MM" как строки
type DayScheduleDTO struct {
	Enabled   bool    `json:"enabled"`              // Включен ли день (false - валидное значение)
	StartTime *string `json:"start_time,omitempty"` // Формат "HH:MM" (nil если disabled)
	EndTime   *string `json:"end_time,omitempty"`   // Формат "HH:MM" (nil если disabled)
}

// WeekdaysDTO - DTO для дней недели
type WeekdaysDTO struct {
	Monday    DayScheduleDTO `json:"monday" binding:"required"`
	Tuesday   DayScheduleDTO `json:"tuesday" binding:"required"`
	Wednesday DayScheduleDTO `json:"wednesday" binding:"required"`
	Thursday  DayScheduleDTO `json:"thursday" binding:"required"`
	Friday    DayScheduleDTO `json:"friday" binding:"required"`
	Saturday  DayScheduleDTO `json:"saturday" binding:"required"`
	Sunday    DayScheduleDTO `json:"sunday" binding:"required"`
}

// VenueDTO - DTO для запросов (Create/Update) и ответов
// Для PUT (Update) все поля обязательны - это полное обновление записи
type VenueDTO struct {
	ID        uint             `json:"id,omitempty"` // Только в ответах
	VenueType models.VenueType `json:"venue_type" binding:"required"`
	OwnerID   uint             `json:"owner_id" binding:"required"`
	IsActive  bool             `json:"is_active"`
	HourPrice int              `json:"hour_price" binding:"required"`
	District  string           `json:"district" binding:"required"`
	Weekdays  WeekdaysDTO      `json:"weekdays" binding:"required"`
}

// ScheduleDTO - DTO для расписания работы площадки (ответ)
// Теперь возвращает полное расписание всех дней недели
type ScheduleDTO struct {
	Weekdays WeekdaysDTO `json:"weekdays"`
}

// ScheduleUpdateDTO - DTO для обновления расписания (запрос)
// Теперь принимает полное расписание всех дней недели
type ScheduleUpdateDTO struct {
	Weekdays WeekdaysDTO `json:"weekdays" binding:"required"`
}

// toDayScheduleDTO конвертирует DaySchedule модели в DTO
func toDayScheduleDTO(schedule models.DaySchedule) DayScheduleDTO {
	dto := DayScheduleDTO{
		Enabled: schedule.Enabled,
	}
	if schedule.Enabled && schedule.StartTime != nil && schedule.EndTime != nil {
		startTimeStr := schedule.StartTime.Format("15:04")
		endTimeStr := schedule.EndTime.Format("15:04")
		dto.StartTime = &startTimeStr
		dto.EndTime = &endTimeStr
	}
	return dto
}

// fromDayScheduleDTO конвертирует DayScheduleDTO в модель
func fromDayScheduleDTO(dto DayScheduleDTO) (models.DaySchedule, error) {
	schedule := models.DaySchedule{
		Enabled: dto.Enabled,
	}

	if dto.Enabled {
		if dto.StartTime == nil || dto.EndTime == nil {
			return schedule, fmt.Errorf("start_time и end_time обязательны, если день включен")
		}

		startTime, err := validation.ValidateTime(*dto.StartTime)
		if err != nil {
			return schedule, fmt.Errorf("неверный формат start_time: %w", err)
		}

		endTime, err := validation.ValidateTime(*dto.EndTime)
		if err != nil {
			return schedule, fmt.Errorf("неверный формат end_time: %w", err)
		}

		schedule.StartTime = &startTime
		schedule.EndTime = &endTime
	}

	return schedule, nil
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
		Weekdays: WeekdaysDTO{
			Monday:    toDayScheduleDTO(venue.Weekdays.Monday),
			Tuesday:   toDayScheduleDTO(venue.Weekdays.Tuesday),
			Wednesday: toDayScheduleDTO(venue.Weekdays.Wednesday),
			Thursday:  toDayScheduleDTO(venue.Weekdays.Thursday),
			Friday:    toDayScheduleDTO(venue.Weekdays.Friday),
			Saturday:  toDayScheduleDTO(venue.Weekdays.Saturday),
			Sunday:    toDayScheduleDTO(venue.Weekdays.Sunday),
		},
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
	monday, err := fromDayScheduleDTO(dto.Weekdays.Monday)
	if err != nil {
		return nil, fmt.Errorf("понедельник: %w", err)
	}

	tuesday, err := fromDayScheduleDTO(dto.Weekdays.Tuesday)
	if err != nil {
		return nil, fmt.Errorf("вторник: %w", err)
	}

	wednesday, err := fromDayScheduleDTO(dto.Weekdays.Wednesday)
	if err != nil {
		return nil, fmt.Errorf("среда: %w", err)
	}

	thursday, err := fromDayScheduleDTO(dto.Weekdays.Thursday)
	if err != nil {
		return nil, fmt.Errorf("четверг: %w", err)
	}

	friday, err := fromDayScheduleDTO(dto.Weekdays.Friday)
	if err != nil {
		return nil, fmt.Errorf("пятница: %w", err)
	}

	saturday, err := fromDayScheduleDTO(dto.Weekdays.Saturday)
	if err != nil {
		return nil, fmt.Errorf("суббота: %w", err)
	}

	sunday, err := fromDayScheduleDTO(dto.Weekdays.Sunday)
	if err != nil {
		return nil, fmt.Errorf("воскресенье: %w", err)
	}

	venue := &models.Venue{
		VenueType: dto.VenueType,
		OwnerID:   dto.OwnerID,
		IsActive:  dto.IsActive,
		HourPrice: dto.HourPrice,
		District:  dto.District,
		Weekdays: models.Weekdays{
			Monday:    monday,
			Tuesday:   tuesday,
			Wednesday: wednesday,
			Thursday:  thursday,
			Friday:    friday,
			Saturday:  saturday,
			Sunday:    sunday,
		},
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
		Weekdays: WeekdaysDTO{
			Monday:    toDayScheduleDTO(venue.Weekdays.Monday),
			Tuesday:   toDayScheduleDTO(venue.Weekdays.Tuesday),
			Wednesday: toDayScheduleDTO(venue.Weekdays.Wednesday),
			Thursday:  toDayScheduleDTO(venue.Weekdays.Thursday),
			Friday:    toDayScheduleDTO(venue.Weekdays.Friday),
			Saturday:  toDayScheduleDTO(venue.Weekdays.Saturday),
			Sunday:    toDayScheduleDTO(venue.Weekdays.Sunday),
		},
	}
}

// FromScheduleUpdateDTO конвертирует ScheduleUpdateDTO в models.Weekdays
func FromScheduleUpdateDTO(dto *ScheduleUpdateDTO) (*models.Weekdays, error) {
	monday, err := fromDayScheduleDTO(dto.Weekdays.Monday)
	if err != nil {
		return nil, fmt.Errorf("понедельник: %w", err)
	}

	tuesday, err := fromDayScheduleDTO(dto.Weekdays.Tuesday)
	if err != nil {
		return nil, fmt.Errorf("вторник: %w", err)
	}

	wednesday, err := fromDayScheduleDTO(dto.Weekdays.Wednesday)
	if err != nil {
		return nil, fmt.Errorf("среда: %w", err)
	}

	thursday, err := fromDayScheduleDTO(dto.Weekdays.Thursday)
	if err != nil {
		return nil, fmt.Errorf("четверг: %w", err)
	}

	friday, err := fromDayScheduleDTO(dto.Weekdays.Friday)
	if err != nil {
		return nil, fmt.Errorf("пятница: %w", err)
	}

	saturday, err := fromDayScheduleDTO(dto.Weekdays.Saturday)
	if err != nil {
		return nil, fmt.Errorf("суббота: %w", err)
	}

	sunday, err := fromDayScheduleDTO(dto.Weekdays.Sunday)
	if err != nil {
		return nil, fmt.Errorf("воскресенье: %w", err)
	}

	return &models.Weekdays{
		Monday:    monday,
		Tuesday:   tuesday,
		Wednesday: wednesday,
		Thursday:  thursday,
		Friday:    friday,
		Saturday:  saturday,
		Sunday:    sunday,
	}, nil
}
