package booking

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"cinema-booking-backend/internal/model"
	"cinema-booking-backend/internal/mq"
	"cinema-booking-backend/internal/repository"
	"cinema-booking-backend/internal/ws"
)

type TimeoutWatcher struct {
	redisClient *redis.Client
	mongoClient *mongo.Client
	repo        *repository.BookingRepository
	producer    *mq.Producer
	hub         *ws.Hub
}

func NewTimeoutWatcher(
	redisClient *redis.Client,
	mongoClient *mongo.Client,
	repo *repository.BookingRepository,
	producer *mq.Producer,
	hub *ws.Hub,
) *TimeoutWatcher {
	return &TimeoutWatcher{
		redisClient: redisClient,
		mongoClient: mongoClient,
		repo:        repo,
		producer:    producer,
		hub:         hub,
	}
}

func (w *TimeoutWatcher) Run(
	ctx context.Context,
) {
	pubsub := w.redisClient.PSubscribe(
		ctx,
		"__keyevent@*__:expired",
	)

	defer pubsub.Close()

	log.Println(
		"timeout watcher: listening for expired seat locks",
	)

	for {
		select {
		case <-ctx.Done():
			return

		case msg, ok := <-pubsub.Channel():
			if !ok {
				return
			}

			w.handleExpiredKey(
				ctx,
				msg.Payload,
			)
		}
	}
}

func (w *TimeoutWatcher) handleExpiredKey(
	ctx context.Context,
	key string,
) {
	const prefix = "lock:seat:"

	if !strings.HasPrefix(
		key,
		prefix,
	) {
		return
	}

	parts := strings.SplitN(
		strings.TrimPrefix(key, prefix),
		":",
		2,
	)

	if len(parts) != 2 {
		return
	}

	showtimeID, err := primitive.ObjectIDFromHex(
		parts[0],
	)

	if err != nil {
		return
	}

	seatID, err := primitive.ObjectIDFromHex(
		parts[1],
	)

	if err != nil {
		return
	}

	bookings, err := w.repo.ListForAdmin(
		ctx,
		bson.M{
			"showtime_id": showtimeID,
			"seat_id":     seatID,
			"status":      model.BookingPending,
		},
	)

	if err != nil {
		log.Printf(
			"timeout watcher: find pending booking failed: %v",
			err,
		)

		return
	}

	if len(bookings) == 0 {
		return
	}

	pending := bookings[0]

	timeoutCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	// PENDING → TIMEOUT
	if err := w.repo.MarkTimeout(
		timeoutCtx,
		pending.ID,
	); err != nil {
		log.Printf(
			"timeout watcher: mark timeout failed: %v",
			err,
		)

		_ = w.repo.LogEvent(
			timeoutCtx,
			model.EventSystemError,
			pending.UserID,
			map[string]interface{}{
				"op":         "MarkTimeout",
				"booking_id": pending.ID.Hex(),
				"error":      err.Error(),
			},
		)

		return
	}

	// Audit: Booking timeout
	if err := w.repo.LogEvent(
		timeoutCtx,
		model.EventBookingTimeout,
		pending.UserID,
		map[string]interface{}{
			"booking_id": pending.ID.Hex(),
			"seat_id":    pending.SeatID.Hex(),
			"seat_label": pending.SeatLabel,
		},
	); err != nil {
		log.Printf(
			"timeout watcher: audit timeout failed: %v",
			err,
		)
	}

	// Audit: Seat released
	if err := w.repo.LogEvent(
		timeoutCtx,
		model.EventSeatReleased,
		pending.UserID,
		map[string]interface{}{
			"booking_id": pending.ID.Hex(),
			"seat_id":    pending.SeatID.Hex(),
			"seat_label": pending.SeatLabel,
		},
	); err != nil {
		log.Printf(
			"timeout watcher: audit release failed: %v",
			err,
		)
	}

	// Kafka
	if err := w.producer.Publish(
		timeoutCtx,
		mq.BookingEvent{
			Type:       mq.EventBookingTimeout,
			BookingID:  pending.ID.Hex(),
			UserID:     pending.UserID.Hex(),
			ShowtimeID: pending.ShowtimeID.Hex(),
			SeatID:     pending.SeatID.Hex(),
			SeatLabel:  pending.SeatLabel,
			OccurredAt: time.Now(),
		},
	); err != nil {
		log.Printf(
			"mq publish failed (non-fatal): %v",
			err,
		)
	}

	// WebSocket
	w.hub.Broadcast(
		ws.SeatUpdate{
			ShowtimeID: pending.ShowtimeID,
			SeatID:     pending.SeatID,
			SeatLabel:  pending.SeatLabel,
			Status:     model.SeatAvailable,
		},
	)
}
