package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentStatus string
type PaymentMethod string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusFailed    PaymentStatus = "failed"

	MethodCard PaymentMethod = "card"
	MethodCash PaymentMethod = "cash"
)

type Payment struct {
	gorm.Model
	BookingID      uuid.UUID     `gorm:"type:uuid;index"`
	UserID         uuid.UUID     `gorm:"type:uuid;index"`
	Amount         int64         `gorm:"column:amount"`
	Currency       string        `gorm:"column:currency"`
	Method         PaymentMethod `gorm:"column:method"`
	Status         PaymentStatus `gorm:"column:status"`
	TransactionID  string        `gorm:"column:transaction_id;uniqueIndex"`
	RefundedAmount int64         `gorm:"column:refunded_amount;default:0"`
	PaidAt         *time.Time    `gorm:"column:paid_at"`
	RefundedAt     *time.Time    `gorm:"column:refunded_at"`
}
