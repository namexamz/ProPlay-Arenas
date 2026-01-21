package dto

import (
	"time"

	"github.com/google/uuid"

	"payment-service/internal/models"
)

type CreatePaymentRequest struct {
	BookingID uuid.UUID            `json:"booking_id" binding:"required,oneof=card cash"`
	UserID    uuid.UUID            `json:"user_id" binding:"required,oneof=card cash"`
	Amount    int64                `json:"amount" binding:"required,gt=0"`
	Currency  string               `json:"currency" binding:"required,oneof=card cash"`
	Method    models.PaymentMethod `json:"method" binding:"required,oneof=card cash"`
}

type PaymentResponse struct {
	ID             uint                 `json:"id"`
	BookingID      uuid.UUID            `json:"booking_id"`
	UserID         uuid.UUID            `json:"user_id"`
	Amount         int64                `json:"amount"`
	Currency       string               `json:"currency"`
	Method         models.PaymentMethod `json:"method"`
	Status         models.PaymentStatus `json:"status"`
	RefundedAmount int64                `json:"refunded_amount"`
	PaidAt         *time.Time           `json:"paid_at"`
	RefundedAt     *time.Time           `json:"refunded_at"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
}

type PaymentHistoryResponse struct {
	Payments []PaymentResponse `json:"payments"`
	Total    int64             `json:"total"`
	Count    int               `json:"count"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
