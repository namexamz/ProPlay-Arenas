package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"payment-service/internal/models"
)

type RefundRepository interface {
	CreateRefund(refund *models.Refund) error
	GetRefundByID(id uint) (*models.Refund, error)
	GetRefundsByPaymentID(paymentID uint) ([]models.Refund, error)
	UpdateRefund(refund *models.Refund) error
}

type RefundRepositoryImpl struct {
	db *gorm.DB
}

func NewRefundRepository(db *gorm.DB) RefundRepository {
	return &RefundRepositoryImpl{db: db}
}

func (r *RefundRepositoryImpl) CreateRefund(refund *models.Refund) error {
	return r.db.Create(refund).Error
}

func (r *RefundRepositoryImpl) GetRefundByID(id uint) (*models.Refund, error) {
	var refund models.Refund
	if err := r.db.Preload("Payment").First(&refund, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("возврат не найден: %w", ErrNotFound)
		}
		return nil, err
	}
	return &refund, nil
}

func (r *RefundRepositoryImpl) GetRefundsByPaymentID(paymentID uint) ([]models.Refund, error) {
	var refunds []models.Refund
	if err := r.db.Where("payment_id = ?", paymentID).Order("created_at DESC").Find(&refunds).Error; err != nil {
		return nil, err
	}
	return refunds, nil
}

func (r *RefundRepositoryImpl) UpdateRefund(refund *models.Refund) error {
	return r.db.Save(refund).Error
}
