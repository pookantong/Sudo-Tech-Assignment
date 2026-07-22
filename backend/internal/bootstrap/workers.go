package bootstrap

import (
	"context"

	"cinema-booking-backend/internal/booking"
	"cinema-booking-backend/internal/mq"
)

func RunWorkers(
	ctx context.Context,
	timeoutWatcher *booking.TimeoutWatcher,
	consumer *mq.Consumer,
) {
	go timeoutWatcher.Run(ctx)
	go consumer.Run(ctx)
}