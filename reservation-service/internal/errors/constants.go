package errors

import "errors"

var (
	ErrClientID                = errors.New("client ID must be greater than zero")
	ErrOwnerID                 = errors.New("owner ID must be greater than zero")
	ErrStartAtEmpty            = errors.New("start time must be provided")
	ErrEndAtEmpty              = errors.New("end time must be provided")
	ErrStartAtAfterEndAt       = errors.New("start time must be before end time")
	ErrStartAtInPast           = errors.New("start time cannot be in the past")
	ErrNegativePrice           = errors.New("price cannot be negative and not be zero")
	ErrStatusEmpty             = errors.New("status must be provided")
	ErrReservationNotFound     = errors.New("reservation not found")
	ErrInvalidStatus           = errors.New("invalid reservation status")
	ErrInvalidRole             = errors.New("вы не являетесь клиентом и не можете создать бронь")
	ErrCannotCancel            = errors.New("cannot cancel reservation")
	ErrOnlyPendingReservations = errors.New("only pending reservations can be updated")
	ErrNotOwner                = errors.New("you are not the owner of this venue")
	ErrForbidden			   = errors.New("forbidden access to the resource")
	ErrDuration             = errors.New("минимальная длительность бронирования - 1 час")
)
