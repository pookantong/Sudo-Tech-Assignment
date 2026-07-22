package mq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type BookingEventType string

const (
	EventBookingSuccess BookingEventType = "BOOKING_SUCCESS"
	EventBookingTimeout BookingEventType = "BOOKING_TIMEOUT"
	EventSeatReleased   BookingEventType = "SEAT_RELEASED"
)

type BookingEvent struct {
	Type BookingEventType `json:"type"`

	BookingID string `json:"booking_id"`
	UserID    string `json:"user_id"`

	ShowtimeID string `json:"showtime_id"`

	SeatID    string `json:"seat_id"`
	SeatLabel string `json:"seat_label"`

	OccurredAt time.Time `json:"occurred_at"`
}

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(
	brokers string,
	topic string,
) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokers),
			Topic:                  topic,
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
		},
	}
}

func (p *Producer) Publish(
	ctx context.Context,
	event BookingEvent,
) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(
		ctx,
		kafka.Message{
			Key:   []byte(event.BookingID),
			Value: body,
		},
	)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(
	brokers string,
	topic string,
	groupID string,
) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(
			kafka.ReaderConfig{
				Brokers: []string{brokers},
				Topic:   topic,
				GroupID: groupID,
			},
		),
	}
}

func (c *Consumer) Run(ctx context.Context) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf(
				"mq consumer stopped: %v",
				err,
			)

			return
		}

		var event BookingEvent

		if err := json.Unmarshal(
			msg.Value,
			&event,
		); err != nil {
			log.Printf(
				"mq consumer: bad message: %v",
				err,
			)

			continue
		}

		log.Printf(
			"[mock notification] type=%s user=%s seat=%s showtime=%s",
			event.Type,
			event.UserID,
			event.SeatLabel,
			event.ShowtimeID,
		)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}