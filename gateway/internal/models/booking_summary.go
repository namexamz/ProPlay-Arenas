package models

import "encoding/json"

type BookingSummaryResponse struct {
	Booking      json.RawMessage `json:"booking"`
	Venue        json.RawMessage `json:"venue,omitempty"`
	Payment      json.RawMessage `json:"payment,omitempty"`
	VenueError   *string         `json:"venue_error,omitempty"`
	PaymentError *string         `json:"payment_error,omitempty"`
}

type BookingSummaryLookup struct {
	VenueID uint `json:"venue_id"`
}
