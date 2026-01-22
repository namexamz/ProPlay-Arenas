package dto

import (
	"reservation/internal/models"
	"time"
)

type ReservationCreate struct {
	VenueID  uint          `json:"venue_id" binding:"required,min=1"`
	ClientID uint          `json:"client_id" binding:"required,min=1"`
	OwnerID  uint          `json:"owner_id" binding:"required,min=1"`
	StartAt  time.Time     `json:"start_at" binding:"required"`
	EndAt    time.Time     `json:"end_at" binding:"required"`

	Status   models.Status `json:"status" binding:"required"`
}

type ReservationCancel struct {
	Reason string `json:"reason" binding:"required"`
}

type ReservationUpdate struct {
	VenueID  *uint      `json:"venue_id,omitempty"`
	ClientID *uint      `json:"client_id,omitempty"`
	OwnerID  *uint      `json:"owner_id,omitempty"`
	StartAt  *time.Time `json:"start_at,omitempty"`
	EndAt    *time.Time `json:"end_at,omitempty"`
	Price    *float64   `json:"price_cents,omitempty"`
}

type BookingCreatedEvent struct {
	EventID   string    `json:"event_id"`
	CreatedAt time.Time `json:"created_at"`

	BookingID uint      `json:"booking_id"`
	VenueID   uint      `json:"venue_id"`
	ClientID  uint      `json:"client_id"`
	OwnerID   uint      `json:"owner_id"`
	StartAt   time.Time `json:"start_at"`
	EndAt     time.Time `json:"end_at"`

	Price  float64       `json:"price_cents"`
	Status models.Status `json:"status"`
}

type BookingCancelledEvent struct {
	EventID   string    `json:"event_id"`
	CreatedAt time.Time `json:"created_at"`

	BookingID uint          `json:"booking_id"`
	Reason    string        `json:"reason"`
	Status    models.Status `json:"status"`
}

type ResponsVenueServ struct {
	ID        uint      `json:"id"`
	OwnerID   uint      `json:"owner_id"`
	StartAt   time.Time `json:"start_at"`
	EndAt     time.Time `json:"end_at"`
	HourPrice float64   `json:"price_cents"`
}

// DayScheduleDTO - DTO для расписания одного дня недели (совместимо с venue-service)
type DayScheduleDTO struct {
	Enabled   bool    `json:"enabled"`
	StartTime *string `json:"start_time,omitempty"`
	EndTime   *string `json:"end_time,omitempty"`
}

// WeekdaysDTO - DTO для расписания всех дней недели (совместимо с venue-service)
type WeekdaysDTO struct {
	Monday    DayScheduleDTO `json:"monday"`
	Tuesday   DayScheduleDTO `json:"tuesday"`
	Wednesday DayScheduleDTO `json:"wednesday"`
	Thursday  DayScheduleDTO `json:"thursday"`
	Friday    DayScheduleDTO `json:"friday"`
	Saturday  DayScheduleDTO `json:"saturday"`
	Sunday    DayScheduleDTO `json:"sunday"`
}

// VenueScheduleResp - расписание в ответе от venue-service
type VenueScheduleResp struct {
	Weekdays WeekdaysDTO `json:"weekdays"`
}

// ResponsVenueServFull - расширенный ответ с расписанием (используется при получении данных от venue-service)
type ResponsVenueServFull struct {
	ID        uint              `json:"id"`
	OwnerID   uint              `json:"owner_id"`
	StartAt   time.Time         `json:"start_at"`
	EndAt     time.Time         `json:"end_at"`
	HourPrice float64           `json:"price_cents"`
	Weekdays  VenueScheduleResp `json:"weekdays"`
}

// AvailableSlot - свободный временной отрезок площадки на дату
type AvailableSlot struct {
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
}
