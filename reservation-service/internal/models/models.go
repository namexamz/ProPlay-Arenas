package models

import (
	"time"
)

type Status string

const (
	Pending   Status = "pending"
	Confirmed Status = "confirmed"
	Cancelled Status = "cancelled"
	Completed Status = "completed"
)

type ReservationDetails struct {
	Base
	VenueID         uint          `json:"venue_id"`
	ClientID        uint          `json:"client_id"`
	OwnerID         uint          `json:"owner_id"`
	StartAt         time.Time     `json:"start_at" gorm:"not null"`
	EndAt           time.Time     `json:"end_at" gorm:"not null"`
	Price           float64       `json:"price_cents,omitempty"`
	Duration        time.Duration `json:"duration_minutes,omitempty"`
	ReasonForCancel string        `json:"reason_for_cancel,omitempty"`
	Status          Status        `json:"status"`
}

type Reservation struct {
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
	Status  Status    `json:"status"`
}
