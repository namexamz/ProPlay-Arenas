package dto

import (
	"reservation/internal/models"
	"time"
)

type ReservationCreate struct {
	ClientID uint          `json:"client_id"`
	OwnerID  uint          `json:"owner_id"`
	StartAt  time.Time     `json:"start_at"`
	EndAt    time.Time     `json:"end_at"`
	Price    float64       `json:"price_cents,omitempty"`
	Status   models.Status `json:"status"`
}
