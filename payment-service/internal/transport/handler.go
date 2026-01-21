package transport

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"

	"payment-service/internal/dto"
	"payment-service/internal/models"
	"payment-service/internal/repository"
	"payment-service/internal/services"
)

type PaymentHandler struct {
	paymentService services.PaymentService
	refundService  services.RefundService
	logger         *slog.Logger
}

func NewPaymentHandler(paymentService services.PaymentService, refundService services.RefundService, logger *slog.Logger) *PaymentHandler {
	if logger == nil {
		logger = slog.Default()
	}
	return &PaymentHandler{
		paymentService: paymentService,
		refundService:  refundService,
		logger:         logger,
	}
}

func (h *PaymentHandler) RegisterRoutes(rg *gin.RouterGroup) {
	payments := rg.Group("/payments")
	{
		payments.POST("", h.CreatePayment)
		payments.GET("", h.GetPaymentsHistory)
		payments.GET("/:id", h.GetPaymentByID)
		payments.POST("/:id/refund", h.CreateRefund)
	}

	bookings := rg.Group("/bookings")
	{
		bookings.GET("/:id/payment", h.GetPaymentByBookingID)
	}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "VALIDATION_ERROR", "400", err.Error())
		return
	}

	payment, err := h.paymentService.CreatePayment(&req)
	if err != nil {
		if isClientError(err) {
			writeError(c, http.StatusBadRequest, "PAYMENT_ERROR", "400", err.Error())
			return
		}
		h.logger.Error("failed to create payment", "error", err)
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "500", "internal server error")
		return
	}

	c.JSON(http.StatusCreated, paymentToResponse(payment))
}

func (h *PaymentHandler) GetPaymentByID(c *gin.Context) {
	id, err := parseUintID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "INVALID_ID", "400", "invalid payment id")
		return
	}

	payment, err := h.paymentService.GetPaymentByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(c, http.StatusNotFound, "NOT_FOUND", "404", "payment not found")
			return
		}
		h.logger.Error("failed to get payment by id", "error", err, "payment_id", id)
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "500", "internal server error")
		return
	}

	c.JSON(http.StatusOK, paymentToResponse(payment))
}

func (h *PaymentHandler) GetPaymentsHistory(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		writeError(c, http.StatusBadRequest, "MISSING_PARAM", "400", "missing user_id")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(c, http.StatusBadRequest, "INVALID_UUID", "400", "invalid user_id")
		return
	}

	limit, offset := parseLimitOffset(c)
	payments, total, err := h.paymentService.GetPaymentsByUserID(userID, limit, offset)
	if err != nil {
		h.logger.Error("failed to get payment history", "error", err, "user_id", userID)
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "500", "internal server error")
		return
	}

	resp := dto.PaymentHistoryResponse{
		Payments: make([]dto.PaymentResponse, len(payments)),
		Total:    total,
		Count:    len(payments),
	}
	for i, p := range payments {
		payment := p
		resp.Payments[i] = *paymentToResponse(&payment)
	}

	c.JSON(http.StatusOK, resp)
}

func (h *PaymentHandler) CreateRefund(c *gin.Context) {
	paymentID, err := parseUintID(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "INVALID_ID", "400", "invalid payment id")
		return
	}

	var req dto.RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "VALIDATION_ERROR", "400", err.Error())
		return
	}

	refund, err := h.refundService.CreateRefund(paymentID, &req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(c, http.StatusNotFound, "NOT_FOUND", "404", "payment not found")
			return
		}
		if isClientError(err) {
			writeError(c, http.StatusBadRequest, "REFUND_ERROR", "400", err.Error())
			return
		}
		h.logger.Error("failed to create refund", "error", err, "payment_id", paymentID)
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "500", "internal server error")
		return
	}

	c.JSON(http.StatusCreated, refundToResponse(refund))
}

func (h *PaymentHandler) GetPaymentByBookingID(c *gin.Context) {
	bookingID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "INVALID_UUID", "400", "invalid booking id")
		return
	}

	payment, err := h.paymentService.GetPaymentByBookingID(bookingID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(c, http.StatusNotFound, "NOT_FOUND", "404", "payment not found")
			return
		}
		h.logger.Error("failed to get payment by booking id", "error", err, "booking_id", bookingID)
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "500", "internal server error")
		return
	}

	c.JSON(http.StatusOK, paymentToResponse(payment))
}

func paymentToResponse(payment *models.Payment) *dto.PaymentResponse {
	return &dto.PaymentResponse{
		ID:             payment.ID,
		BookingID:      payment.BookingID,
		UserID:         payment.UserID,
		Amount:         payment.Amount,
		Currency:       payment.Currency,
		Method:         payment.Method,
		Status:         payment.Status,
		RefundedAmount: payment.RefundedAmount,
		PaidAt:         payment.PaidAt,
		RefundedAt:     payment.RefundedAt,
		CreatedAt:      payment.CreatedAt,
		UpdatedAt:      payment.UpdatedAt,
	}
}

func refundToResponse(refund *models.Refund) *dto.RefundResponse {
	return &dto.RefundResponse{
		ID:        refund.ID,
		PaymentID: refund.PaymentID,
		Amount:    refund.Amount,
		Reason:    refund.Reason,
		Status:    refund.Status,
		CreatedAt: refund.CreatedAt,
		UpdatedAt: refund.UpdatedAt,
	}
}

func writeError(c *gin.Context, status int, errCode, code, message string) {
	c.JSON(status, dto.ErrorResponse{
		Error:   errCode,
		Code:    code,
		Message: message,
	})
}

func parseUintID(raw string) (uint, error) {
	id64, err := strconv.ParseUint(raw, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id64), nil
}

func parseLimitOffset(c *gin.Context) (int, int) {
	limit := 10
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	return limit, offset
}

func isClientError(err error) bool {
	return errors.Is(err, services.ErrEmptyRequest) ||
		errors.Is(err, services.ErrInvalidMethod) ||
		errors.Is(err, services.ErrInvalidAmount) ||
		errors.Is(err, services.ErrEmptyCurrency) ||
		errors.Is(err, services.ErrPaymentNotPending) ||
		errors.Is(err, services.ErrPaymentNotComplete) ||
		errors.Is(err, services.ErrRefundAmountExceed)
}
