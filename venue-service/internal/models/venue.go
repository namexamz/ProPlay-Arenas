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

// Weekdays структура для дней недели, когда можно делать бронирования
// true - можно делать бронь, false - нельзя
type Weekdays struct {
	Monday    bool `json:"monday" gorm:"column:monday;default:true"`       // Понедельник
	Tuesday   bool `json:"tuesday" gorm:"column:tuesday;default:true"`     // Вторник
	Wednesday bool `json:"wednesday" gorm:"column:wednesday;default:true"` // Среда
	Thursday  bool `json:"thursday" gorm:"column:thursday;default:true"`   // Четверг
	Friday    bool `json:"friday" gorm:"column:friday;default:true"`       // Пятница
	Saturday  bool `json:"saturday" gorm:"column:saturday;default:true"`   // Суббота
	Sunday    bool `json:"sunday" gorm:"column:sunday;default:true"`       // Воскресенье
}

type Venue struct {
	gorm.Model
	VenueType VenueType `json:"venue_type" gorm:"column:venue_type;type:varchar(50);not null"`
	OwnerID   uint      `json:"owner_id" gorm:"column:owner_id;not null;index"`
	IsActive  bool      `json:"is_active" gorm:"column:is_active;default:true"`
	HourPrice int       `json:"hour_price" gorm:"column:hour_price;not null;check:hour_price >= 0"`
	District  string    `json:"district" gorm:"column:district;type:varchar(50);not null"`
	StartTime time.Time `json:"start_time" gorm:"column:start_time;type:time;not null"` // Рабочее время начала
	EndTime   time.Time `json:"end_time" gorm:"column:end_time;type:time;not null"`     // Рабочее время окончания
	Weekdays  Weekdays  `json:"weekdays" gorm:"embedded"`                               // Дни недели для бронирования
}

func (Venue) TableName() string {
	return "venues"
}

// validateVenue проверяет валидность данных площадки
func (v *Venue) validateVenue() error {
	if !v.VenueType.IsValid() {
		return fmt.Errorf("неверный тип площадки: %s", v.VenueType)
	}

	// Проверка, что время начала раньше времени окончания
	// Сравниваем только часы и минуты (игнорируем дату)
	startTime := time.Date(0, 1, 1, v.StartTime.Hour(), v.StartTime.Minute(), v.StartTime.Second(), 0, time.UTC)
	endTime := time.Date(0, 1, 1, v.EndTime.Hour(), v.EndTime.Minute(), v.EndTime.Second(), 0, time.UTC)

	if !startTime.Before(endTime) {
		return fmt.Errorf("время начала должно быть раньше времени окончания")
	}

	return nil
}

func (v *Venue) BeforeCreate(tx *gorm.DB) error {
	return v.validateVenue()
}

func (v *Venue) BeforeUpdate(tx *gorm.DB) error {
	return v.validateVenue()
}
