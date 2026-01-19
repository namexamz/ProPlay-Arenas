package dto

import (
	"reservation/internal/models"
	"time"
)

type ReservationCreate struct {
	VenueID  uint          `json:"venue_id" binding:"required"`
	ClientID uint          `json:"client_id"`
	OwnerID  uint          `json:"owner_id"`
	StartAt  time.Time     `json:"start_at"`
	EndAt    time.Time     `json:"end_at"`
	Price    float64       `json:"price_cents,omitempty"`
	Status   models.Status `json:"status"`
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
