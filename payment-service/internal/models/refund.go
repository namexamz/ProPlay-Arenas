package models

import (
	"gorm.io/gorm"
)

type RefundStatus string

const (
	RefundStatusPending   RefundStatus = "pending"
	RefundStatusCompleted RefundStatus = "completed"
	RefundStatusFailed    RefundStatus = "failed"
)

type Refund struct {
	gorm.Model
	PaymentID uint         `gorm:"index"`
	Amount    int64        `gorm:"column:amount"`
	Reason    string       `gorm:"column:reason;type:text"`
	Status    RefundStatus `gorm:"column:status"`
	Payment   *Payment     `gorm:"foreignKey:PaymentID;references:ID"`
}
