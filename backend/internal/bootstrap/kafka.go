package bootstrap

import (
	"cinema-booking-backend/internal/config"
	"cinema-booking-backend/internal/mq"
)

func NewKafka(cfg *config.Config) (*mq.Producer, *mq.Consumer) {
	producer := mq.NewProducer(
		cfg.KafkaBrokers,
		cfg.KafkaTopic,
	)

	consumer := mq.NewConsumer(
		cfg.KafkaBrokers,
		cfg.KafkaTopic,
		"notification-service",
	)

	return producer, consumer
}