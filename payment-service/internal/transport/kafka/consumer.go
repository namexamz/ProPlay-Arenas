package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/google/uuid"

	"payment-service/internal/config"
	"payment-service/internal/dto"
	"payment-service/internal/models"
	"payment-service/internal/services"
)

type Consumer struct {
	paymentService services.PaymentService
	refundService  services.RefundService
	logger         *slog.Logger
	brokers        []string
	groupID        string
	createdTopic   string
	cancelledTopic string
}

type BookingCreatedEvent struct {
	BookingID uuid.UUID            `json:"booking_id"`
	UserID    uuid.UUID            `json:"user_id"`
	Amount    int64                `json:"amount"`
	Method    models.PaymentMethod `json:"method"`
}

type BookingCancelledEvent struct {
	BookingID uuid.UUID `json:"booking_id"`
}

func NewConsumerFromEnv(paymentService services.PaymentService, refundService services.RefundService, logger *slog.Logger) *Consumer {
	if logger == nil {
		logger = slog.Default()
	}

	brokers := splitBrokers(config.GetEnv("KAFKA_BROKERS", ""))
	return &Consumer{
		paymentService: paymentService,
		refundService:  refundService,
		logger:         logger,
		brokers:        brokers,
		groupID:        config.GetEnv("KAFKA_GROUP_ID", "payment-service"),
		createdTopic:   config.GetEnv("KAFKA_TOPIC_BOOKING_CREATED", "booking.created"),
		cancelledTopic: config.GetEnv("KAFKA_TOPIC_BOOKING_CANCELLED", "booking.cancelled"),
	}
}

func (c *Consumer) Start(ctx context.Context) {
	if len(c.brokers) == 0 {
		c.logger.Warn("Kafka brokers not configured, consumer disabled")
		return
	}

	go c.consumeBookingCreated(ctx)
	go c.consumeBookingCancelled(ctx)
}

func (c *Consumer) consumeBookingCreated(ctx context.Context) {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: c.brokers,
		GroupID: c.groupID,
		Topic:   c.createdTopic,
	})
	defer reader.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			c.logger.Error("ошибка чтения сообщения booking.created", "error", err)
			continue
		}

		var event BookingCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			c.logger.Error("ошибка парсинга booking.created", "error", err)
			continue
		}

		if event.BookingID == uuid.Nil || event.UserID == uuid.Nil {
			c.logger.Error("некорректные данные booking.created", "booking_id", event.BookingID, "user_id", event.UserID)
			continue
		}

		if event.Amount <= 0 {
			c.logger.Error("некорректная сумма в booking.created", "amount", event.Amount)
			continue
		}

		if event.Method == "" {
			event.Method = models.MethodCard
		}

		req := dto.CreatePaymentRequest{
			BookingID: event.BookingID,
			UserID:    event.UserID,
			Amount:    event.Amount,
			Currency:  "RUB",
			Method:    event.Method,
		}

		if _, err := c.paymentService.CreatePendingPayment(&req); err != nil {
			c.logger.Error("ошибка создания pending платежа из booking.created", "error", err, "booking_id", event.BookingID)
			continue
		}
	}
}

func (c *Consumer) consumeBookingCancelled(ctx context.Context) {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: c.brokers,
		GroupID: c.groupID,
		Topic:   c.cancelledTopic,
	})
	defer reader.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			c.logger.Error("ошибка чтения сообщения booking.cancelled", "error", err)
			continue
		}

		var event BookingCancelledEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			c.logger.Error("ошибка парсинга booking.cancelled", "error", err)
			continue
		}

		if event.BookingID == uuid.Nil {
			c.logger.Error("некорректные данные booking.cancelled", "booking_id", event.BookingID)
			continue
		}

		payment, err := c.paymentService.GetPaymentByBookingID(event.BookingID)
		if err != nil {
			c.logger.Error("ошибка получения платежа по booking_id для возврата", "error", err, "booking_id", event.BookingID)
			continue
		}

		if payment.Status != models.PaymentStatusCompleted {
			continue
		}

		remaining := payment.Amount - payment.RefundedAmount
		if remaining <= 0 {
			continue
		}

		if _, err := c.refundService.CreateRefund(payment.ID, &dto.RefundRequest{
			Amount: remaining,
			Reason: "отмена бронирования",
		}); err != nil {
			c.logger.Error("ошибка создания возврата по booking.cancelled", "error", err, "payment_id", payment.ID)
			continue
		}
	}
}

func splitBrokers(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
