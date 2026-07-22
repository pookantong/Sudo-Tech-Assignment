package catalog

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"cinema-booking-backend/internal/cache"
	"cinema-booking-backend/internal/model"
	"cinema-booking-backend/internal/repository"
)

type Service struct {
	repo        *repository.CatalogRepository
	bookingRepo *repository.BookingRepository
	lock        *cache.SeatLock
}

func NewService(
	repo *repository.CatalogRepository,
	bookingRepo *repository.BookingRepository,
	lock *cache.SeatLock,
) *Service {
	return &Service{
		repo:        repo,
		bookingRepo: bookingRepo,
		lock:        lock,
	}
}

// ====================
// Movie
// ====================

func (s *Service) ListMovies(
	ctx context.Context,
) ([]model.Movie, error) {
	return s.repo.ListMovies(ctx)
}

// ====================
// Hall
// ====================

func (s *Service) ListHalls(
	ctx context.Context,
) ([]model.Hall, error) {
	return s.repo.ListHalls(ctx)
}

func (s *Service) GetHall(
	ctx context.Context,
	hallID primitive.ObjectID,
) (*model.Hall, error) {
	return s.repo.GetHall(ctx, hallID)
}

// ====================
// Showtime
// ====================

func (s *Service) ListShowtimesByMovie(
	ctx context.Context,
	movieID primitive.ObjectID,
) ([]model.ShowtimeResponse, error) {
	showtimes, err := s.repo.ListShowtimesByMovie(ctx, movieID)
	if err != nil {
		return nil, err
	}

	result := make(
		[]model.ShowtimeResponse,
		0,
		len(showtimes),
	)

	for _, showtime := range showtimes {
		hall, err := s.repo.GetHall(
			ctx,
			showtime.HallID,
		)

		if err != nil {
			return nil, err
		}

		cinema, err := s.repo.GetCinema(
			ctx,
			hall.CinemaID,
		)
		if err != nil {
			return nil, err
		}

		result = append(
			result,
			model.ShowtimeResponse{
				ID:         showtime.ID,
				MovieID:    showtime.MovieID,
				HallID:     showtime.HallID,
				HallName:   hall.Name,
				CinemaID:   cinema.ID,
				CinemaName: cinema.Name,
				StartsAt:   showtime.StartsAt,
				Price:      showtime.Price,
			},
		)
	}

	return result, nil
}

// ====================
// Seat Map
// ====================

type SeatView struct {
	ID     string           `json:"id"`
	Label  string           `json:"label"`
	Row    int              `json:"row"`
	Col    int              `json:"col"`
	Status model.SeatStatus `json:"status"`
	Price  float64          `json:"price"`
}

func (s *Service) GetSeatMap(
	ctx context.Context,
	showtimeID primitive.ObjectID,
) ([]SeatView, error) {
	// 1. หา Showtime
	showtime, err := s.repo.GetShowtime(
		ctx,
		showtimeID,
	)
	if err != nil {
		return nil, err
	}

	// 2. ดึงที่นั่งของ Hall
	seats, err := s.repo.ListSeatsByHall(
		ctx,
		showtime.HallID,
	)
	if err != nil {
		return nil, err
	}

	// 3. ดึงที่นั่งที่ถูกจองแล้ว
	bookedSeatIDs, err := s.bookingRepo.ListBookedSeatIDsByShowtime(
		ctx,
		showtimeID,
	)
	if err != nil {
		return nil, err
	}

	// 4. Timeout สำหรับ Redis
	lockCtx, cancel := context.WithTimeout(
		ctx,
		5*time.Second,
	)
	defer cancel()

	views := make(
		[]SeatView,
		0,
		len(seats),
	)

	for _, seat := range seats {
		status := model.SeatAvailable

		// เช็ก BOOKED ก่อน
		if bookedSeatIDs[seat.ID] {
			status = model.SeatBooked
		} else {
			// ถ้ายังไม่ BOOKED ค่อยเช็ก Redis LOCKED
			locked, _, err := s.lock.IsLocked(
				lockCtx,
				showtimeID.Hex(),
				seat.ID.Hex(),
			)

			if err == nil && locked {
				status = model.SeatLocked
			}
		}

		views = append(
			views,
			SeatView{
				ID:     seat.ID.Hex(),
				Label:  seat.Label,
				Row:    seat.Row,
				Col:    seat.Col,
				Status: status,
				Price:  showtime.Price,
			},
		)
	}

	return views, nil
}
