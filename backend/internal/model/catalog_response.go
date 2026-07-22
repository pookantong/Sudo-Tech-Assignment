package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShowtimeResponse struct {
	ID         primitive.ObjectID `json:"id"`
	MovieID    primitive.ObjectID `json:"movie_id"`
	HallID     primitive.ObjectID `json:"hall_id"`
	HallName   string             `json:"hall_name"`
	CinemaID   primitive.ObjectID `json:"cinema_id"`
	CinemaName string             `json:"cinema_name"`
	StartsAt   time.Time          `json:"starts_at"`
	Price      float64            `json:"price"`
}