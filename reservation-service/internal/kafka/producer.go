package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"reservation/internal/dto"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

type Producer interface {
	PublishBookingCreated(ctx context.Context, evt dto.BookingCreatedEvent) error
	PublishBookingCancelled(ctx context.Context, evt dto.BookingCancelledEvent) error
	Close() error
}

const (
	TopicBookingCreated   = "booking.created"
	TopicBookingCancelled = "booking.cancelled"
)

type kafkaGoProducer struct {
	writer *kafkago.Writer
}

func NewProducer(brokers []string) Producer {
	return &kafkaGoProducer{
		writer: &kafkago.Writer{
			Addr:         kafkago.TCP(brokers...),
			Balancer:     &kafkago.LeastBytes{},
			RequiredAcks: kafkago.RequireOne,
		},
	}
}

func (p *kafkaGoProducer) PublishBookingCreated(ctx context.Context, evt dto.BookingCreatedEvent) error {
	return p.writeJSON(ctx, TopicBookingCreated, fmt.Sprintf("%d", evt.BookingID), evt)
}

func (p *kafkaGoProducer) PublishBookingCancelled (ctx context.Context, evt dto.BookingCancelledEvent)error {
	return p.writeJSON(ctx, TopicBookingCancelled, fmt.Sprintf("%d", evt.BookingID), evt)
}

func (p *kafkaGoProducer) writeJSON(ctx context.Context, topic, key string, v any) error {
	// 1. Превращаем структуру в JSON
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	// 2. Создаем сообщение для Kafka
	msg := kafkago.Message{
		Topic: topic,       // В какой топик отправить
		Key:   []byte(key), // Ключ (для порядка)
		Value: b,           // Сами данные (JSON)
		Time:  time.Now(),  // Время отправки
	}

	return p.writer.WriteMessages(ctx, msg)
}

func (p *kafkaGoProducer) Close() error {
	return p.writer.Close()
}