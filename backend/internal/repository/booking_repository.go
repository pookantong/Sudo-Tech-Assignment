package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"cinema-booking-backend/internal/model"
)

var ErrNotFound = errors.New("document not found")

type BookingRepository struct {
	bookings *mongo.Collection
	audits   *mongo.Collection
}

func NewBookingRepository(
	db *mongo.Database,
) *BookingRepository {
	return &BookingRepository{
		bookings: db.Collection("bookings"),
		audits:   db.Collection("audit_logs"),
	}
}

// ====================
// Booking
// ====================

func (r *BookingRepository) CreatePending(
	ctx context.Context,
	b *model.Booking,
) error {
	now := time.Now()

	b.Status = model.BookingPending
	b.CreatedAt = now
	b.UpdatedAt = now

	res, err := r.bookings.InsertOne(ctx, b)
	if err != nil {
		return err
	}

	b.ID = res.InsertedID.(primitive.ObjectID)

	return nil
}

func (r *BookingRepository) ConfirmSuccess(ctx context.Context, bookingID primitive.ObjectID) error {
	_, err := r.bookings.UpdateOne(ctx,
		bson.M{"_id": bookingID},
		bson.M{"$set": bson.M{"status": model.BookingSuccess, "updated_at": time.Now()}},
	)
	return err
}

func (r *BookingRepository) MarkFailed(
	ctx context.Context,
	bookingID primitive.ObjectID,
) error {
	_, err := r.bookings.UpdateOne(
		ctx,
		bson.M{
			"_id":    bookingID,
			"status": model.BookingPending,
		},
		bson.M{
			"$set": bson.M{
				"status":     model.BookingFailed,
				"updated_at": time.Now(),
			},
		},
	)

	return err
}

func (r *BookingRepository) MarkTimeout(
	ctx context.Context,
	bookingID primitive.ObjectID,
) error {
	_, err := r.bookings.UpdateOne(
		ctx,
		bson.M{
			"_id":    bookingID,
			"status": model.BookingPending,
		},
		bson.M{
			"$set": bson.M{
				"status":     model.BookingTimeout,
				"updated_at": time.Now(),
			},
		},
	)

	return err
}

func (r *BookingRepository) GetByID(
	ctx context.Context,
	bookingID primitive.ObjectID,
) (*model.Booking, error) {
	var b model.Booking

	err := r.bookings.
		FindOne(
			ctx,
			bson.M{"_id": bookingID},
		).
		Decode(&b)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &b, nil
}

func (r *BookingRepository) ListForAdmin(
	ctx context.Context,
	filter bson.M,
) ([]model.Booking, error) {
	cur, err := r.bookings.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []model.Booking

	if err := cur.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// ====================
// Seat Availability
// ====================

// Find booking ที่กำลัง lock หรือจองที่นั่งอยู่
func (r *BookingRepository) FindActiveBookingBySeat(
	ctx context.Context,
	showtimeID primitive.ObjectID,
	seatID primitive.ObjectID,
) (*model.Booking, error) {
	var booking model.Booking

	err := r.bookings.FindOne(
		ctx,
		bson.M{
			"showtime_id": showtimeID,
			"seat_id":     seatID,
			"status": bson.M{
				"$in": []model.BookingStatus{
					model.BookingPending,
					model.BookingSuccess,
				},
			},
		},
	).Decode(&booking)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err
	}

	return &booking, nil
}

// ====================
// Audit Log
// ====================

func (r *BookingRepository) LogEvent(
	ctx context.Context,
	eventType model.AuditEventType,
	userID primitive.ObjectID,
	metadata map[string]interface{},
) error {
	_, err := r.audits.InsertOne(
		ctx,
		model.AuditLog{
			Type:      eventType,
			UserID:    userID,
			Metadata:  metadata,
			CreatedAt: time.Now(),
		},
	)

	return err
}

func (r *BookingRepository) ListBookedSeatIDsByShowtime(
	ctx context.Context,
	showtimeID primitive.ObjectID,
) (map[primitive.ObjectID]bool, error) {
	cursor, err := r.bookings.Find(
		ctx,
		bson.M{
			"showtime_id": showtimeID,
			"status":      model.BookingSuccess,
		},
	)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	bookedSeatIDs := make(map[primitive.ObjectID]bool)

	for cursor.Next(ctx) {
		var booking model.Booking

		if err := cursor.Decode(&booking); err != nil {
			return nil, err
		}

		bookedSeatIDs[booking.SeatID] = true
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return bookedSeatIDs, nil
}