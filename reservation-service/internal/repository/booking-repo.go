package repository

import (
	"errors"
	"reservation/internal/models"

	"gorm.io/gorm"
)

type BookingRepo interface {
	GetByID(id uint) (*models.ReservationDetails, error)
	GetUserReservations(userID uint) ([]models.Reservation, error)
	GetVenueBookings(venueID uint, ownerID uint) ([]models.ReservationDetails, error)
	Create(reservation *models.ReservationDetails) error
	Save(reservation *models.ReservationDetails) error
}

type gormBookingRepo struct {
	db *gorm.DB
}

func NewBookingRepo(db *gorm.DB) BookingRepo {
	return &gormBookingRepo{db: db}
}

var (
	ErrFindReservations = errors.New("reservations not found")
)

func (r *gormBookingRepo) GetUserReservations(userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation

	result := r.db.Model(&models.ReservationDetails{}).
		Where("client_id = ?", userID).
		Select("start_at", "end_at", "status").
		Find(&reservations)

	if result.Error != nil {
		return nil, result.Error
	}

	return reservations, nil
}

func (r *gormBookingRepo) Create(reservation *models.ReservationDetails) error {
	result := r.db.Create(reservation)
	return result.Error
}

func (r *gormBookingRepo) GetByID(id uint) (*models.ReservationDetails, error) {
	var reservation models.ReservationDetails

	result := r.db.First(&reservation, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &reservation, nil
}

func (r *gormBookingRepo) Save(reservation *models.ReservationDetails) error {
	result := r.db.Save(reservation)
	return result.Error
}

func (r *gormBookingRepo) GetVenueBookings(venueID uint, ownerID uint) ([]models.ReservationDetails, error) {
	var bookings []models.ReservationDetails

	result := r.db.Where("venue_id = ? AND owner_id = ?", venueID, ownerID).Find(&bookings)
	if result.Error != nil {
		return nil, result.Error
	}

	return bookings, nil
}
