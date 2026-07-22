package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"cinema-booking-backend/internal/model"
)

type CatalogRepository struct {
	cinemas  *mongo.Collection
	movies   *mongo.Collection
	halls    *mongo.Collection
	showtimes *mongo.Collection
	seats    *mongo.Collection
}

func NewCatalogRepository(db *mongo.Database) *CatalogRepository {
	return &CatalogRepository{
		cinemas:  db.Collection("cinemas"),
		movies:   db.Collection("movies"),
		halls:    db.Collection("halls"),
		showtimes: db.Collection("showtimes"),
		seats:     db.Collection("seats"),
	}
}

// ====================
// Cinema
// ====================

func (r *CatalogRepository) GetCinema(
	ctx context.Context,
	cinemaID primitive.ObjectID,
) (*model.Cinema, error) {
	var cinema model.Cinema

	err := r.cinemas.
		FindOne(
			ctx,
			bson.M{
				"_id": cinemaID,
			},
		).
		Decode(&cinema)

	if err != nil {
		return nil, err
	}

	return &cinema, nil
}

func (r *CatalogRepository) ListCinemas(
	ctx context.Context,
) ([]model.Cinema, error) {
	cur, err := r.cinemas.Find(
		ctx,
		bson.M{},
	)

	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	var cinemas []model.Cinema

	if err := cur.All(ctx, &cinemas); err != nil {
		return nil, err
	}

	return cinemas, nil
}

// ====================
// Movie
// ====================

func (r *CatalogRepository) ListMovies(
	ctx context.Context,
) ([]model.Movie, error) {
	cur, err := r.movies.Find(
		ctx,
		bson.M{},
	)

	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	var movies []model.Movie

	if err := cur.All(ctx, &movies); err != nil {
		return nil, err
	}

	return movies, nil
}

// ====================
// Hall
// ====================

func (r *CatalogRepository) GetHall(
	ctx context.Context,
	hallID primitive.ObjectID,
) (*model.Hall, error) {
	var hall model.Hall

	err := r.halls.
		FindOne(
			ctx,
			bson.M{
				"_id": hallID,
			},
		).
		Decode(&hall)

	if err != nil {
		return nil, err
	}

	return &hall, nil
}

func (r *CatalogRepository) ListHalls(
	ctx context.Context,
) ([]model.Hall, error) {
	cur, err := r.halls.Find(
		ctx,
		bson.M{},
	)

	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	var halls []model.Hall

	if err := cur.All(ctx, &halls); err != nil {
		return nil, err
	}

	return halls, nil
}

func (r *CatalogRepository) ListHallsByCinema(
	ctx context.Context,
	cinemaID primitive.ObjectID,
) ([]model.Hall, error) {
	cur, err := r.halls.Find(
		ctx,
		bson.M{
			"cinema_id": cinemaID,
		},
	)

	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	var halls []model.Hall

	if err := cur.All(ctx, &halls); err != nil {
		return nil, err
	}

	return halls, nil
}

// ====================
// Showtime
// ====================

func (r *CatalogRepository) ListShowtimesByMovie(
	ctx context.Context,
	movieID primitive.ObjectID,
) ([]model.Showtime, error) {
	cur, err := r.showtimes.Find(
		ctx,
		bson.M{
			"movie_id": movieID,
		},
	)

	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	var showtimes []model.Showtime

	if err := cur.All(ctx, &showtimes); err != nil {
		return nil, err
	}

	return showtimes, nil
}

func (r *CatalogRepository) ListShowtimesByMovieAndCinema(
	ctx context.Context,
	movieID primitive.ObjectID,
	cinemaID primitive.ObjectID,
) ([]model.Showtime, error) {
	cur, err := r.showtimes.Find(
		ctx,
		bson.M{
			"movie_id":  movieID,
			"cinema_id": cinemaID,
		},
	)

	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	var showtimes []model.Showtime

	if err := cur.All(ctx, &showtimes); err != nil {
		return nil, err
	}

	return showtimes, nil
}

func (r *CatalogRepository) GetShowtime(
	ctx context.Context,
	showtimeID primitive.ObjectID,
) (*model.Showtime, error) {
	var showtime model.Showtime

	err := r.showtimes.
		FindOne(
			ctx,
			bson.M{
				"_id": showtimeID,
			},
		).
		Decode(&showtime)

	if err != nil {
		return nil, err
	}

	return &showtime, nil
}

// ====================
// Seat
// ====================

func (r *CatalogRepository) ListSeatsByHall(
	ctx context.Context,
	hallID primitive.ObjectID,
) ([]model.Seat, error) {
	cur, err := r.seats.Find(
		ctx,
		bson.M{
			"hall_id": hallID,
		},
	)

	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	var seats []model.Seat

	if err := cur.All(ctx, &seats); err != nil {
		return nil, err
	}

	return seats, nil
}

func (r *CatalogRepository) GetSeat(
	ctx context.Context,
	seatID primitive.ObjectID,
) (*model.Seat, error) {
	var seat model.Seat

	err := r.seats.
		FindOne(
			ctx,
			bson.M{
				"_id": seatID,
			},
		).
		Decode(&seat)

	if err != nil {
		return nil, err
	}

	return &seat, nil
}

func (r *CatalogRepository) GenerateSeatsForHall(
	ctx context.Context,
	hallID primitive.ObjectID,
	rows int,
	cols int,
) error {
	docs := make(
		[]interface{},
		0,
		rows*cols,
	)

	for row := 0; row < rows; row++ {
		rowLabel := string(
			rune('A' + row),
		)

		for col := 1; col <= cols; col++ {
			docs = append(
				docs,
				model.Seat{
					HallID: hallID,
					Label: fmt.Sprintf(
						"%s%d",
						rowLabel,
						col,
					),
					Row: row,
					Col: col,
				},
			)
		}
	}

	if len(docs) == 0 {
		return nil
	}

	_, err := r.seats.InsertMany(
		ctx,
		docs,
	)

	return err
}