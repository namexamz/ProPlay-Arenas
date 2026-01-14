package service

import (
	"reservation/internal/models"
	"reservation/internal/repository"


)

type BookingService interface {
	
	GetUserReservations(userID uint) ([]models.Reservation, error)
}

type bookingService struct {
	repo repository.BokingRepo
}

func NewBookingServ(repo repository.BokingRepo) BookingService {
	return &bookingService{repo: repo}
}



func (r *bookingService) GetUserReservations(userID uint) ([]models.Reservation, error) {
	reservations, err := r.repo.GetUserReservations(userID)

	if err != nil {
		return nil, err
	}

	return reservations, nil
}
