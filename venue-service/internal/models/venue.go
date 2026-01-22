package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type VenueType string

const (
	VenueFootball   VenueType = "football"
	VenueBasketball VenueType = "basketball"
	VenueTennis     VenueType = "tennis"
	VenueGym        VenueType = "gym"
	VenueSwimming   VenueType = "swimming"
)

func (vt VenueType) IsValid() bool {
	switch vt {
	case VenueFootball, VenueBasketball, VenueTennis, VenueGym, VenueSwimming:
		return true
	default:
		return false
	}
}

func (vt VenueType) String() string {
	return string(vt)
}

// DaySchedule структура для расписания одного дня недели
// Если Enabled = false, то StartTime и EndTime должны быть nil
type DaySchedule struct {
	Enabled   bool       `json:"enabled" gorm:"column:enabled;default:true"`              // Включен ли день для бронирования
	StartTime *time.Time `json:"start_time,omitempty" gorm:"column:start_time;type:time"` // Время начала работы (nil если disabled)
	EndTime   *time.Time `json:"end_time,omitempty" gorm:"column:end_time;type:time"`     // Время окончания работы (nil если disabled)
}

// Weekdays структура для дней недели, когда можно делать бронирования
// Каждый день имеет свое расписание (время начала и окончания работы)
type Weekdays struct {
	Monday    DaySchedule `json:"monday" gorm:"embedded;embeddedPrefix:monday_"`       // Понедельник
	Tuesday   DaySchedule `json:"tuesday" gorm:"embedded;embeddedPrefix:tuesday_"`     // Вторник
	Wednesday DaySchedule `json:"wednesday" gorm:"embedded;embeddedPrefix:wednesday_"` // Среда
	Thursday  DaySchedule `json:"thursday" gorm:"embedded;embeddedPrefix:thursday_"`   // Четверг
	Friday    DaySchedule `json:"friday" gorm:"embedded;embeddedPrefix:friday_"`       // Пятница
	Saturday  DaySchedule `json:"saturday" gorm:"embedded;embeddedPrefix:saturday_"`   // Суббота
	Sunday    DaySchedule `json:"sunday" gorm:"embedded;embeddedPrefix:sunday_"`       // Воскресенье
}

type Venue struct {
	gorm.Model
	VenueType VenueType `json:"venue_type" gorm:"column:venue_type;type:varchar(50);not null"`
	OwnerID   uint      `json:"owner_id" gorm:"column:owner_id;not null;index"`
	IsActive  bool      `json:"is_active" gorm:"column:is_active;default:true"`
	HourPrice int       `json:"hour_price" gorm:"column:hour_price;not null;check:hour_price >= 0"`
	District  string    `json:"district" gorm:"column:district;type:varchar(50);not null"`
	Weekdays  Weekdays  `json:"weekdays" gorm:"embedded"` // Дни недели для бронирования с расписанием
}

func (Venue) TableName() string {
	return "venues"
}

// validateVenue проверяет валидность данных площадки
func (v *Venue) validateVenue() error {
	if !v.VenueType.IsValid() {
		return fmt.Errorf("неверный тип площадки: %s", v.VenueType)
	}

	// Проверяем расписание для каждого дня недели
	days := []struct {
		name     string
		schedule DaySchedule
	}{
		{"понедельник", v.Weekdays.Monday},
		{"вторник", v.Weekdays.Tuesday},
		{"среда", v.Weekdays.Wednesday},
		{"четверг", v.Weekdays.Thursday},
		{"пятница", v.Weekdays.Friday},
		{"суббота", v.Weekdays.Saturday},
		{"воскресенье", v.Weekdays.Sunday},
	}

	for _, day := range days {
		if day.schedule.Enabled {
			// Если день включен, StartTime и EndTime должны быть заданы
			if day.schedule.StartTime == nil || day.schedule.EndTime == nil {
				return fmt.Errorf("для %s время начала и окончания должны быть указаны", day.name)
			}
			// Проверка, что время начала раньше времени окончания
			if !day.schedule.StartTime.Before(*day.schedule.EndTime) {
				return fmt.Errorf("для %s время начала должно быть раньше времени окончания", day.name)
			}
		} else {
			// Если день выключен, StartTime и EndTime должны быть nil
			if day.schedule.StartTime != nil || day.schedule.EndTime != nil {
				return fmt.Errorf("для %s время начала и окончания должны быть nil, если день выключен", day.name)
			}
		}
	}

	return nil
}

func (v *Venue) BeforeCreate(tx *gorm.DB) error {
	return v.validateVenue()
}

func (v *Venue) BeforeUpdate(tx *gorm.DB) error {
	return v.validateVenue()
}
