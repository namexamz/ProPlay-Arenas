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
	Price    float64       `json:"price_cents" binding:"required,min=0.01"`
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
	ID      uint `json:"id"`
	OwnerID uint `json:"owner_id"`
}
