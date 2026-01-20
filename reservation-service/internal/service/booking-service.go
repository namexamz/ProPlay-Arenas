package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reservation/internal/dto"
	"reservation/internal/kafka"
	"reservation/internal/models"
	"reservation/internal/repository"
	"time"

	"github.com/google/uuid"
)

type BookingService interface {
	GetUserReservations(userID uint) ([]models.Reservation, error)
	CreateReservation(reservation *dto.ReservationCreate, claims *models.Claims) (*models.ReservationDetails, error)
	ReservationCancel(id uint, reason string) (*models.ReservationDetails, error)
	GetByID(id uint) (*models.ReservationDetails, error)
	ReservationUpdate(id uint, reservation *dto.ReservationUpdate) (*models.ReservationDetails, error)
}

type bookingService struct {
	repo     repository.BookingRepo
	producer kafka.Producer
}

func NewBookingServ(repo repository.BookingRepo, producer kafka.Producer) BookingService {
	return &bookingService{repo: repo, producer: producer}
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
	ErrInvalidStatus       = errors.New("invalid reservation status")
	ErrInvalidRole         = errors.New("вы не являетесь клиентом и не можете создать бронь")
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

func (r *bookingService) CreateReservation(reservation *dto.ReservationCreate, claims *models.Claims) (*models.ReservationDetails, error) {

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

	if claims.Role != models.RoleClient && claims.Role != models.RoleAdmin {
		return nil, ErrInvalidRole
	}

	reservation.ClientID = claims.UserID

	newReservation := &models.ReservationDetails{
		ClientID: reservation.ClientID,
		VenueID:  reservation.VenueID,
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

	evt := dto.BookingCreatedEvent{
		EventID:   uuid.NewString(),
		CreatedAt: time.Now(),
		BookingID: newReservation.ID,
		VenueID:   newReservation.VenueID,
		ClientID:  newReservation.ClientID,
		OwnerID:   newReservation.OwnerID,
		StartAt:   newReservation.StartAt,
		EndAt:     newReservation.EndAt,
		Price:     newReservation.Price,
		Status:    newReservation.Status,
	}

	if err := r.producer.PublishBookingCreated(context.Background(), evt); err != nil {
		log.Printf("Ошибка отправки события в Kafka: %v", err)
		return newReservation, fmt.Errorf("бронь создана (id=%d), но не удалось отправить событие в Kafka: %w", newReservation.ID, err)
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

	evt := dto.BookingCancelledEvent{
		EventID:   uuid.NewString(),
		CreatedAt: time.Now(),
		BookingID: reservation.ID,
		Reason:    reason,
		Status:    reservation.Status,
	}

	if err := r.producer.PublishBookingCancelled(context.Background(), evt); err != nil {
		log.Printf("Ошибка отправки события отмены в Kafka: %v", err)
	}

	return reservation, nil
}

func (r *bookingService) ReservationUpdate(id uint, reservation *dto.ReservationUpdate) (*models.ReservationDetails, error) {

	reserv, err := r.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if reservation.ClientID != nil && *reservation.ClientID <= 0 {
		return nil, ErrClientID
	}

	if reservation.OwnerID != nil && *reservation.OwnerID <= 0 {
		return nil, ErrOwnerID
	}

	if reservation.StartAt != nil && reservation.StartAt.IsZero() {
		return nil, ErrStartAtEmpty
	}

	if reservation.EndAt != nil && reservation.EndAt.IsZero() {
		return nil, ErrEndAtEmpty
	}

	// Определяем финальные значения для валидации (не мутируя reserv заранее)
	finalStartAt := reserv.StartAt
	if reservation.StartAt != nil {
		finalStartAt = *reservation.StartAt
	}

	finalEndAt := reserv.EndAt
	if reservation.EndAt != nil {
		finalEndAt = *reservation.EndAt
	}

	// Проверяем, что StartAt < EndAt для итогового диапазона
	if !finalStartAt.Before(finalEndAt) {
		return nil, ErrStartAtAfterEndAt
	}

	// Проверяем, что finalStartAt не в прошлом (независимо от того, обновляется ли он)
	if finalStartAt.Before(time.Now()) {
		return nil, ErrStartAtInPast
	}

	if reservation.Price != nil && *reservation.Price <= 0 {
		return nil, ErrNegativePrice
	}

	if reserv.Status != models.Pending {
		return nil, errors.New("only pending reservations can be updated. Current status: " + string(reserv.Status))
	}

	if reservation.VenueID != nil {
		reserv.VenueID = *reservation.VenueID
	}

	if reservation.ClientID != nil {
		reserv.ClientID = *reservation.ClientID
	}

	if reservation.OwnerID != nil {
		reserv.OwnerID = *reservation.OwnerID
	}

	if reservation.StartAt != nil {
		reserv.StartAt = *reservation.StartAt
	}

	if reservation.EndAt != nil {
		reserv.EndAt = *reservation.EndAt
	}

	if reservation.Price != nil {
		reserv.Price = *reservation.Price
	}

	if err := r.repo.Save(reserv); err != nil {
		return nil, err
	}

	reserv.Duration = reserv.EndAt.Sub(reserv.StartAt)

	return reserv, nil

}
