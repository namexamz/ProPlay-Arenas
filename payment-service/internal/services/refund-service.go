package services

import (
	"log/slog"
	"time"

	"payment-service/internal/dto"
	"payment-service/internal/models"
	"payment-service/internal/repository"
)

type RefundService interface {
	CreateRefund(paymentID uint, req *dto.RefundRequest) (*models.Refund, error)
	GetRefundByID(id uint) (*models.Refund, error)
	GetRefundsByPaymentID(paymentID uint) ([]models.Refund, error)
}

type RefundServiceImpl struct {
	refundRepo  repository.RefundRepository
	paymentRepo repository.PaymentRepository
	logger      *slog.Logger
}

func NewRefundService(refundRepo repository.RefundRepository, paymentRepo repository.PaymentRepository) RefundService {
	return &RefundServiceImpl{
		refundRepo:  refundRepo,
		paymentRepo: paymentRepo,
		logger:      slog.Default(),
	}
}

func (s *RefundServiceImpl) CreateRefund(paymentID uint, req *dto.RefundRequest) (*models.Refund, error) {
	if req == nil {
		s.logger.Error("пустой запрос на создание возврата")
		return nil, ErrEmptyRequest
	}
	if req.Amount <= 0 {
		s.logger.Error("некорректная сумма возврата", "amount", req.Amount)
		return nil, ErrInvalidAmount
	}

	payment, err := s.paymentRepo.GetPaymentByID(paymentID)
	if err != nil {
		s.logger.Error("ошибка получения платежа для возврата", "payment_id", paymentID, "error", err)
		return nil, err
	}

	if payment.Status != models.PaymentStatusCompleted {
		s.logger.Error("платеж не завершен, возврат невозможен", "payment_id", paymentID, "status", payment.Status)
		return nil, ErrPaymentNotComplete
	}

	available := payment.Amount - payment.RefundedAmount
	if req.Amount > available {
		s.logger.Error("сумма возврата превышает доступную", "payment_id", paymentID, "amount", req.Amount, "available", available)
		return nil, ErrRefundAmountExceed
	}

	refund := &models.Refund{
		PaymentID: paymentID,
		Amount:    req.Amount,
		Reason:    req.Reason,
		Status:    models.RefundStatusCompleted,
	}

	if err := s.refundRepo.CreateRefund(refund); err != nil {
		s.logger.Error("ошибка сохранения возврата", "error", err)
		return nil, err
	}

	payment.RefundedAmount += req.Amount
	if payment.RefundedAmount >= payment.Amount {
		payment.Status = models.PaymentStatusRefunded
		now := time.Now()
		payment.RefundedAt = &now
	}

	if err := s.paymentRepo.UpdatePayment(payment); err != nil {
		s.logger.Error("ошибка обновления платежа при возврате", "payment_id", payment.ID, "error", err)
		return nil, err
	}

	s.logger.Info("возврат создан", "refund_id", refund.ID, "payment_id", payment.ID)
	return refund, nil
}

func (s *RefundServiceImpl) GetRefundByID(id uint) (*models.Refund, error) {
	refund, err := s.refundRepo.GetRefundByID(id)
	if err != nil {
		s.logger.Error("ошибка получения возврата по id", "refund_id", id, "error", err)
		return nil, err
	}
	return refund, nil
}

func (s *RefundServiceImpl) GetRefundsByPaymentID(paymentID uint) ([]models.Refund, error) {
	refunds, err := s.refundRepo.GetRefundsByPaymentID(paymentID)
	if err != nil {
		s.logger.Error("ошибка получения возвратов по платежу", "payment_id", paymentID, "error", err)
		return nil, err
	}
	return refunds, nil
}
