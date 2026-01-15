package models

import (
	"fmt"

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

type Venue struct {
	gorm.Model
	VenueType VenueType `json:"venue_type" gorm:"column:venue_type;type:varchar(50);not null"`
	OwnerID   uint      `json:"owner_id" gorm:"column:owner_id;not null;index"`
	IsActive  bool      `json:"is_active" gorm:"column:is_active;default:true"`
	HourPrice int       `json:"hour_price" gorm:"column:hour_price;not null;check:hour_price >= 0"`
}

func (Venue) TableName() string {
	return "venues"
}

func (v *Venue) BeforeCreate(tx *gorm.DB) error {
	if !v.VenueType.IsValid() {
		return fmt.Errorf("неверный тип площадки: %s", v.VenueType)
	}
	return nil
}

func (v *Venue) BeforeUpdate(tx *gorm.DB) error {
	if !v.VenueType.IsValid() {
		return fmt.Errorf("неверный тип площадки: %s", v.VenueType)
	}
	return nil
}
