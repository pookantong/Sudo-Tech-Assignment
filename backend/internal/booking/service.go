package booking

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"cinema-booking-backend/internal/cache"
	"cinema-booking-backend/internal/model"
	"cinema-booking-backend/internal/mq"
	"cinema-booking-backend/internal/repository"
	"cinema-booking-backend/internal/ws"
)

var (
	ErrSeatTaken = errors.New(
		"seat is already locked or booked",
	)

	ErrSeatNotFound = errors.New(
		"seat not found",
	)

	ErrSeatNotInShowtime = errors.New(
		"seat does not belong to this showtime",
	)

	ErrBookingNotFound = errors.New(
		"booking not found",
	)

	ErrNotOwner = errors.New(
		"this booking does not belong to the requesting user",
	)
)

// NOTE: no more *mongo.Client field here. ConfirmSuccess used to need it to
// open a session/transaction (updating "bookings" + "seats" together), but
// since Seat status is now derived from Booking rather than stored on a
// per-showtime Seat document, ConfirmSuccess is a single collection write —
// no transaction needed, which also means no MongoDB replica-set
// requirement (a standalone mongod errors on session numbers with
// "Transaction numbers are only allowed on a replica set member or mongos").
type Service struct {
	lock        *cache.SeatLock
	repo        *repository.BookingRepository
	catalogRepo *repository.CatalogRepository
	producer    *mq.Producer
	hub         *ws.Hub
}

func NewService(
	lock *cache.SeatLock,
	repo *repository.BookingRepository,
	catalogRepo *repository.CatalogRepository,
	producer *mq.Producer,
	hub *ws.Hub,
) *Service {
	return &Service{
		lock:        lock,
		repo:        repo,
		catalogRepo: catalogRepo,
		producer:    producer,
		hub:         hub,
	}
}

func (s *Service) SelectSeat(
	ctx context.Context,
	userID primitive.ObjectID,
	showtimeID primitive.ObjectID,
	seatID primitive.ObjectID,
) (*model.Booking, error) {
	showtime, err := s.catalogRepo.GetShowtime(
		ctx,
		showtimeID,
	)
	if err != nil {
		if errors.Is(
			err,
			mongo.ErrNoDocuments,
		) {
			return nil, ErrBookingNotFound
		}

		return nil, err
	}

	seat, err := s.catalogRepo.GetSeat(
		ctx,
		seatID,
	)
	if err != nil {
		if errors.Is(
			err,
			mongo.ErrNoDocuments,
		) {
			return nil, ErrSeatNotFound
		}

		return nil, err
	}

	if seat.HallID != showtime.HallID {
		return nil, ErrSeatNotInShowtime
	}

	ownerToken := userID.Hex() + ":" + uuid.NewString()

	if err := s.lock.Acquire(
		ctx,
		showtimeID.Hex(),
		seatID.Hex(),
		ownerToken,
	); err != nil {
		if errors.Is(
			err,
			cache.ErrLockNotAcquired,
		) {
			return nil, ErrSeatTaken
		}

		return nil, err
	}

	booking := &model.Booking{
		UserID:     userID,
		ShowtimeID: showtimeID,
		SeatID:     seatID,
		SeatLabel:  seat.Label,
		Price:      showtime.Price,
		LockToken:  ownerToken,
	}

	if err := s.repo.CreatePending(
		ctx,
		booking,
	); err != nil {
		_ = s.lock.Release(
			ctx,
			showtimeID.Hex(),
			seatID.Hex(),
			ownerToken,
		)

		return nil, err
	}

	s.hub.Broadcast(
		ws.SeatUpdate{
			ShowtimeID: showtimeID,
			SeatID:     seatID,
			SeatLabel:  seat.Label,
			Status:     model.SeatLocked,
		},
	)

	return booking, nil
}

func (s *Service) getOwnedBooking(
	ctx context.Context,
	bookingID primitive.ObjectID,
	requestingUserID primitive.ObjectID,
) (*model.Booking, error) {
	booking, err := s.repo.GetByID(
		ctx,
		bookingID,
	)

	if err != nil {
		if errors.Is(
			err,
			repository.ErrNotFound,
		) {
			return nil, ErrBookingNotFound
		}

		return nil, err
	}

	if booking.UserID != requestingUserID {
		return nil, ErrNotOwner
	}

	return booking, nil
}

func (s *Service) ConfirmPayment(
	ctx context.Context,
	bookingID primitive.ObjectID,
	requestingUserID primitive.ObjectID,
) error {
	booking, err := s.getOwnedBooking(
		ctx,
		bookingID,
		requestingUserID,
	)

	if err != nil {
		return err
	}

	locked, owner, err := s.lock.IsLocked(
		ctx,
		booking.ShowtimeID.Hex(),
		booking.SeatID.Hex(),
	)

	if err != nil {
		return err
	}

	if !locked || owner != booking.LockToken {
		return cache.ErrLockNotOwned
	}

	// No *mongo.Client / session param anymore — see the note on Service
	// above for why this no longer needs a transaction.
	if err := s.repo.ConfirmSuccess(
		ctx,
		booking.ID,
	); err != nil {
		return err
	}

	_ = s.lock.Release(
		ctx,
		booking.ShowtimeID.Hex(),
		booking.SeatID.Hex(),
		booking.LockToken,
	)

	_ = s.repo.LogEvent(
		ctx,
		model.EventBookingSuccess,
		booking.UserID,
		map[string]interface{}{
			"booking_id": booking.ID.Hex(),
			"seat_id":    booking.SeatID.Hex(),
			"seat_label": booking.SeatLabel,
		},
	)

	if err := s.producer.Publish(
		ctx,
		mq.BookingEvent{
			Type:       mq.EventBookingSuccess,
			BookingID:  booking.ID.Hex(),
			UserID:     booking.UserID.Hex(),
			ShowtimeID: booking.ShowtimeID.Hex(),
			SeatID:     booking.SeatID.Hex(),
			SeatLabel:  booking.SeatLabel,
			OccurredAt: time.Now(),
		},
	); err != nil {
		log.Printf(
			"mq publish failed (non-fatal): %v",
			err,
		)
	}

	s.hub.Broadcast(
		ws.SeatUpdate{
			ShowtimeID: booking.ShowtimeID,
			SeatID:     booking.SeatID,
			SeatLabel:  booking.SeatLabel,
			Status:     model.SeatBooked,
		},
	)

	return nil
}

func (s *Service) FailPayment(
	ctx context.Context,
	bookingID primitive.ObjectID,
	requestingUserID primitive.ObjectID,
) error {
	booking, err := s.getOwnedBooking(
		ctx,
		bookingID,
		requestingUserID,
	)

	if err != nil {
		return err
	}

	locked, owner, err := s.lock.IsLocked(
		ctx,
		booking.ShowtimeID.Hex(),
		booking.SeatID.Hex(),
	)

	if err != nil {
		return err
	}

	if !locked || owner != booking.LockToken {
		return cache.ErrLockNotOwned
	}

	if err := s.repo.MarkFailed(
		ctx,
		booking.ID,
	); err != nil {
		return err
	}

	if err := s.lock.Release(
		ctx,
		booking.ShowtimeID.Hex(),
		booking.SeatID.Hex(),
		booking.LockToken,
	); err != nil {
		log.Printf(
			"lock release failed (non-fatal, will expire via TTL): %v",
			err,
		)
	}

	_ = s.repo.LogEvent(
		ctx,
		model.EventSeatReleased,
		booking.UserID,
		map[string]interface{}{
			"booking_id": booking.ID.Hex(),
			"seat_id":    booking.SeatID.Hex(),
			"seat_label": booking.SeatLabel,
			"reason":     "payment_failed",
		},
	)

	if err := s.producer.Publish(
		ctx,
		mq.BookingEvent{
			Type:       mq.EventSeatReleased,
			BookingID:  booking.ID.Hex(),
			UserID:     booking.UserID.Hex(),
			ShowtimeID: booking.ShowtimeID.Hex(),
			SeatID:     booking.SeatID.Hex(),
			SeatLabel:  booking.SeatLabel,
			OccurredAt: time.Now(),
		},
	); err != nil {
		log.Printf(
			"mq publish failed (non-fatal): %v",
			err,
		)
	}

	s.hub.Broadcast(
		ws.SeatUpdate{
			ShowtimeID: booking.ShowtimeID,
			SeatID:     booking.SeatID,
			SeatLabel:  booking.SeatLabel,
			Status:     model.SeatAvailable,
		},
	)

	return nil
}

type AdminBookingFilter struct {
	UserID     string
	ShowtimeID string
	Status     string
}

func (s *Service) ListForAdmin(
	ctx context.Context,
	filter AdminBookingFilter,
) ([]model.Booking, error) {
	query := bson.M{}

	if filter.UserID != "" {
		if oid, err := primitive.ObjectIDFromHex(
			filter.UserID,
		); err == nil {
			query["user_id"] = oid
		}
	}

	if filter.ShowtimeID != "" {
		if oid, err := primitive.ObjectIDFromHex(
			filter.ShowtimeID,
		); err == nil {
			query["showtime_id"] = oid
		}
	}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	return s.repo.ListForAdmin(
		ctx,
		query,
	)
}

var ErrBookingForbidden = errors.New(
	"booking does not belong to user",
)
