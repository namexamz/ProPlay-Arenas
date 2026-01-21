package repository

import (
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/gorm"

	"payment-service/internal/models"
)

type RefundRepository interface {
	CreateRefund(refund *models.Refund) error
	GetRefundByID(id uint) (*models.Refund, error)
	GetRefundsByPaymentID(paymentID uint) ([]models.Refund, error)
	UpdateRefund(refund *models.Refund) error
}

var refundRepoLogger = slog.Default()

type RefundRepositoryImpl struct {
	db *gorm.DB
}

func NewRefundRepository(db *gorm.DB) RefundRepository {
	return &RefundRepositoryImpl{db: db}
}

func (r *RefundRepositoryImpl) CreateRefund(refund *models.Refund) error {
	if err := r.db.Create(refund).Error; err != nil {
		refundRepoLogger.Error("ошибка создания возврата", "error", err)
		return err
	}
	refundRepoLogger.Info("возврат создан", "refund_id", refund.ID)
	return nil
}

func (r *RefundRepositoryImpl) GetRefundByID(id uint) (*models.Refund, error) {
	var refund models.Refund
	if err := r.db.Preload("Payment").First(&refund, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			refundRepoLogger.Warn("возврат не найден", "refund_id", id)
			return nil, fmt.Errorf("возврат не найден: %w", ErrNotFound)
		}
		refundRepoLogger.Error("ошибка получения возврата по id", "refund_id", id, "error", err)
		return nil, err
	}
	return &refund, nil
}

func (r *RefundRepositoryImpl) GetRefundsByPaymentID(paymentID uint) ([]models.Refund, error) {
	var refunds []models.Refund
	if err := r.db.Where("payment_id = ?", paymentID).Order("created_at DESC").Find(&refunds).Error; err != nil {
		refundRepoLogger.Error("ошибка получения возвратов по платежу", "payment_id", paymentID, "error", err)
		return nil, err
	}
	return refunds, nil
}

func (r *RefundRepositoryImpl) UpdateRefund(refund *models.Refund) error {
	if err := r.db.Save(refund).Error; err != nil {
		refundRepoLogger.Error("ошибка обновления возврата", "refund_id", refund.ID, "error", err)
		return err
	}
	refundRepoLogger.Info("возврат обновлен", "refund_id", refund.ID)
	return nil
}
