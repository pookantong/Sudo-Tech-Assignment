package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
)

type BookingStatus string

const (
	BookingPending BookingStatus = "PENDING"
	BookingSuccess BookingStatus = "SUCCESS"
	BookingTimeout BookingStatus = "TIMEOUT"
	BookingFailed  BookingStatus = "FAILED"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GoogleSub string             `bson:"google_sub" json:"-"`
	Email     string             `bson:"email" json:"email"`
	Name      string             `bson:"name" json:"name"`
	Role      Role               `bson:"role" json:"role"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type Cinema struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type Movie struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title    string             `bson:"title" json:"title"`
	ImageURL string             `bson:"image_url" json:"image_url"`
}

type Hall struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CinemaID primitive.ObjectID `bson:"cinema_id" json:"cinema_id"`
	Name     string             `bson:"name" json:"name"`
	SeatRows int                `bson:"seat_rows" json:"seat_rows"`
	SeatCols int                `bson:"seat_cols" json:"seat_cols"`
}

type Seat struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	HallID primitive.ObjectID `bson:"hall_id" json:"hall_id"`
	Label  string             `bson:"label" json:"label"`
	Row    int                `bson:"row" json:"row"`
	Col    int                `bson:"col" json:"col"`
}

type SeatStatus string

const (
	SeatAvailable SeatStatus = "AVAILABLE"
	SeatLocked    SeatStatus = "LOCKED"
	SeatBooked    SeatStatus = "BOOKED"
)

type Showtime struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MovieID  primitive.ObjectID `bson:"movie_id" json:"movie_id"`
	CinemaID primitive.ObjectID `bson:"cinema_id" json:"cinema_id"`
	HallID   primitive.ObjectID `bson:"hall_id" json:"hall_id"`
	StartsAt time.Time          `bson:"starts_at" json:"starts_at"`
	Price    float64            `bson:"price" json:"price"`
}

type Booking struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	ShowtimeID primitive.ObjectID `bson:"showtime_id" json:"showtime_id"`
	SeatID     primitive.ObjectID `bson:"seat_id" json:"seat_id"`
	SeatLabel  string             `bson:"seat_label" json:"seat_label"`
	Price      float64            `bson:"price" json:"price"`
	Status     BookingStatus      `bson:"status" json:"status"`
	LockToken  string             `bson:"lock_token" json:"-"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

type AuditEventType string

const (
	EventBookingSuccess AuditEventType = "BOOKING_SUCCESS"
	EventBookingTimeout AuditEventType = "BOOKING_TIMEOUT"
	EventSeatReleased   AuditEventType = "SEAT_RELEASED"
	EventSystemError    AuditEventType = "SYSTEM_ERROR"
)

type AuditLog struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Type      AuditEventType         `bson:"type" json:"type"`
	UserID    primitive.ObjectID     `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Metadata  map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt time.Time              `bson:"created_at" json:"created_at"`
}