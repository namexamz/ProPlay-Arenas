package repository

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"payment-service/internal/models"
)

var ErrNotFound = errors.New("не найдено")

var paymentRepoLogger = slog.Default()

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
	if err := r.db.Create(payment).Error; err != nil {
		paymentRepoLogger.Error("ошибка создания платежа", "error", err)
		return err
	}
	paymentRepoLogger.Info("платеж создан", "payment_id", payment.ID)
	return nil
}

func (r *PaymentRepositoryImpl) GetPaymentByID(id uint) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.First(&payment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			paymentRepoLogger.Warn("платеж не найден", "payment_id", id)
			return nil, fmt.Errorf("платеж не найден: %w", ErrNotFound)
		}
		paymentRepoLogger.Error("ошибка получения платежа по id", "payment_id", id, "error", err)
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepositoryImpl) GetPaymentsByUserID(userID uuid.UUID, limit, offset int) ([]models.Payment, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	var payments []models.Payment
	var total int64

	if err := r.db.Where("user_id = ?", userID).Model(&models.Payment{}).Count(&total).Error; err != nil {
		paymentRepoLogger.Error("ошибка подсчета платежей пользователя", "user_id", userID, "error", err)
		return nil, 0, err
	}

	if err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&payments).Error; err != nil {
		paymentRepoLogger.Error("ошибка получения платежей пользователя", "user_id", userID, "error", err)
		return nil, 0, err
	}

	return payments, total, nil
}

func (r *PaymentRepositoryImpl) GetPaymentByTransactionID(transactionID string) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.Where("transaction_id = ?", transactionID).First(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			paymentRepoLogger.Warn("платеж по transaction_id не найден", "transaction_id", transactionID)
			return nil, fmt.Errorf("платеж по transaction_id не найден: %w", ErrNotFound)
		}
		paymentRepoLogger.Error("ошибка получения платежа по transaction_id", "transaction_id", transactionID, "error", err)
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepositoryImpl) GetPaymentByBookingID(bookingID uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.Where("booking_id = ?", bookingID).First(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			paymentRepoLogger.Warn("платеж по booking_id не найден", "booking_id", bookingID)
			return nil, fmt.Errorf("платеж по booking_id не найден: %w", ErrNotFound)
		}
		paymentRepoLogger.Error("ошибка получения платежа по booking_id", "booking_id", bookingID, "error", err)
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepositoryImpl) UpdatePayment(payment *models.Payment) error {
	if err := r.db.Save(payment).Error; err != nil {
		paymentRepoLogger.Error("ошибка обновления платежа", "payment_id", payment.ID, "error", err)
		return err
	}
	paymentRepoLogger.Info("платеж обновлен", "payment_id", payment.ID)
	return nil
}

func (r *PaymentRepositoryImpl) DeletePayment(id uint) error {
	if err := r.db.Delete(&models.Payment{}, id).Error; err != nil {
		paymentRepoLogger.Error("ошибка удаления платежа", "payment_id", id, "error", err)
		return err
	}
	paymentRepoLogger.Info("платеж удален", "payment_id", id)
	return nil
}
