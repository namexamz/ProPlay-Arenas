package dto

import (
	"time"

	"github.com/google/uuid"

	"payment-service/internal/models"
)

type CreatePaymentRequest struct {
	BookingID uuid.UUID            `json:"booking_id" binding:"required"`
	UserID    uuid.UUID            `json:"user_id" binding:"required"`
	Amount    int64                `json:"amount" binding:"required,gt=0"`
	Currency  string               `json:"currency" binding:"required"`
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
	TransactionID  string               `json:"transaction_id"`
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

type RefundRequest struct {
	Amount int64  `json:"amount" binding:"required,gt=0"`
	Reason string `json:"reason" binding:"required,min=5,max=500"`
}

type RefundResponse struct {
	ID        uint                `json:"id"`
	PaymentID uint                `json:"payment_id"`
	Amount    int64               `json:"amount"`
	Reason    string              `json:"reason"`
	Status    models.RefundStatus `json:"status"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
