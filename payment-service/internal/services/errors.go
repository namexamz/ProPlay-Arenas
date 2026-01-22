package services

import "errors"

var (
	ErrEmptyRequest       = errors.New("пустой запрос")
	ErrInvalidAmount      = errors.New("сумма должна быть больше нуля")
	ErrInvalidMethod      = errors.New("недопустимый метод оплаты")
	ErrPaymentNotComplete = errors.New("платеж не завершен")
	ErrRefundAmountExceed = errors.New("сумма возврата превышает доступную")
)
