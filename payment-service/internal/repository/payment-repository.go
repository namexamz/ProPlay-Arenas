package repository

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"payment-service/internal/models"
)

var ErrNotFound = errors.New("не найдено")

type PaymentRepository interface {
	CreatePayment(payment *models.Payment) error
	GetPaymentByID(id uint) (*models.Payment, error)
	GetPaymentsByUserID(userID uuid.UUID, limit, offset int) ([]models.Payment, int64, error)
	GetPaymentByTransactionID(transactionID string) (*models.Payment, error)
	GetPaymentByBookingID(bookingID uuid.UUID) (*models.Payment, error)
	UpdatePayment(payment *models.Payment) error
	DeletePayment(id uint) error
}

type PaymentRepositoryImpl struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &PaymentRepositoryImpl{db: db}
}

func (r *PaymentRepositoryImpl) CreatePayment(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *PaymentRepositoryImpl) GetPaymentByID(id uint) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.First(&payment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("платеж не найден: %w", ErrNotFound)
		}
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepositoryImpl) GetPaymentsByUserID(userID uuid.UUID, limit, offset int) ([]models.Payment, int64, error) {
	var payments []models.Payment
	var total int64

	if err := r.db.Where("user_id = ?", userID).Model(&models.Payment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

func (r *PaymentRepositoryImpl) GetPaymentByTransactionID(transactionID string) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.Where("transaction_id = ?", transactionID).First(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("платеж по transaction_id не найден: %w", ErrNotFound)
		}
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepositoryImpl) GetPaymentByBookingID(bookingID uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.Where("booking_id = ?", bookingID).First(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("платеж по booking_id не найден: %w", ErrNotFound)
		}
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepositoryImpl) UpdatePayment(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

func (r *PaymentRepositoryImpl) DeletePayment(id uint) error {
	return r.db.Delete(&models.Payment{}, id).Error
}
