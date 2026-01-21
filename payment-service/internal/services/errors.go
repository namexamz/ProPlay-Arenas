package services

import "errors"

var (
	ErrEmptyRequest       = errors.New("пустой запрос")
	ErrInvalidMethod      = errors.New("недопустимый метод оплаты")
	ErrInvalidAmount      = errors.New("сумма должна быть больше нуля")
	ErrEmptyCurrency      = errors.New("валюта не указана")
	ErrPaymentNotPending  = errors.New("платеж не в статусе pending")
	ErrPaymentNotComplete = errors.New("платеж не завершен")
	ErrRefundAmountExceed = errors.New("сумма возврата превышает доступную")
)
