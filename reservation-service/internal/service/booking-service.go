package service

import (
	"errors"
	"reservation/internal/dto"
	"reservation/internal/models"
	"reservation/internal/repository"
	"time"
)

type BookingService interface {
	GetUserReservations(userID uint) ([]models.Reservation, error)
	CreateReservation(reservation *dto.ReservationCreate) (*models.ReservationDetails, error)
	ReservationCancel(id uint, reason string) (*models.ReservationDetails, error)
	GetByID(id uint) (*models.ReservationDetails, error)
}

type bookingService struct {
	repo repository.BookingRepo
}

func NewBookingServ(repo repository.BookingRepo) BookingService {
	return &bookingService{repo: repo}
}

var (
	ErrClientID            = errors.New("client ID must be greater than zero")
	ErrOwnerID             = errors.New("owner ID must be greater than zero")
	ErrStartAtEmpty        = errors.New("start time must be provided")
	ErrEndAtEmpty          = errors.New("end time must be provided")
	ErrStartAtAfterEndAt   = errors.New("start time must be before end time")
	ErrStartAtInPast       = errors.New("start time cannot be in the past")
	ErrNegativePrice       = errors.New("price cannot be negative and not be zero")
	ErrStatusEmpty         = errors.New("status must be provided")
	ErrReservationNotFound = errors.New("reservation not found")
)

func (r *bookingService) GetUserReservations(userID uint) ([]models.Reservation, error) {
	reservations, err := r.repo.GetUserReservations(userID)

	if err != nil {
		return nil, err
	}

	return reservations, nil
}

func (r *bookingService) GetByID(id uint) (*models.ReservationDetails, error) {
	reservation, err := r.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (r *bookingService) CreateReservation(reservation *dto.ReservationCreate) (*models.ReservationDetails, error) {
	if reservation.ClientID <= 0 {
		return nil, ErrClientID
	}

	if reservation.OwnerID <= 0 {
		return nil, ErrOwnerID
	}

	if reservation.StartAt.IsZero() {
		return nil, ErrStartAtEmpty
	}

	if reservation.EndAt.IsZero() {
		return nil, ErrEndAtEmpty
	}

	if !reservation.StartAt.Before(reservation.EndAt) {
		return nil, ErrStartAtAfterEndAt
	}

	if reservation.StartAt.Before(time.Now()) {
		return nil, ErrStartAtInPast
	}

	if reservation.Price <= 0 {
		return nil, ErrNegativePrice
	}

	if reservation.Status == "" {
		return nil, ErrStatusEmpty
	}

	newReservation := &models.ReservationDetails{
		ClientID: reservation.ClientID,
		OwnerID:  reservation.OwnerID,
		StartAt:  reservation.StartAt,
		EndAt:    reservation.EndAt,
		Price:    float64(reservation.Price),
		Status:   models.Status(reservation.Status),
		Duration: reservation.EndAt.Sub(reservation.StartAt),
	}

	if err := r.repo.Create(newReservation); err != nil {
		return nil, err
	}

	return newReservation, nil
}

func (r *bookingService) ReservationCancel(id uint, reason string) (*models.ReservationDetails, error) {
	reservation, err := r.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if reservation.Status == models.Cancelled || reservation.Status == models.Completed {
		return nil, errors.New("cannot cancel reservation")
	}

	reservation.Status = models.Cancelled
	reservation.ReasonForCancel = reason

	if err := r.repo.Save(reservation); err != nil {
		return nil, err
	}

	return reservation, nil
}
