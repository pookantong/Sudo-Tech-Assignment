package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"cinema-booking-backend/internal/catalog"
)

type CatalogHandler struct {
	service *catalog.Service
}

func NewCatalogHandler(
	service *catalog.Service,
) *CatalogHandler {
	return &CatalogHandler{
		service: service,
	}
}

// ====================
// Movies
// ====================

func (h *CatalogHandler) ListMovies(
	c *gin.Context,
) {
	movies, err := h.service.ListMovies(
		c.Request.Context(),
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "internal error",
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		movies,
	)
}

type createMovieRequest struct {
	Title string `json:"title" binding:"required"`
}

func (h *CatalogHandler) CreateMovie(
	c *gin.Context,
) {
	var req createMovieRequest

	if err := c.ShouldBindJSON(
		&req,
	); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	movie, err := h.service.CreateMovie(
		c.Request.Context(),
		req.Title,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "internal error",
			},
		)

		return
	}

	c.JSON(
		http.StatusCreated,
		movie,
	)
}

// ====================
// Halls
// ====================

func (h *CatalogHandler) ListHalls(
	c *gin.Context,
) {
	halls, err := h.service.ListHalls(
		c.Request.Context(),
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "internal error",
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		halls,
	)
}

type createHallRequest struct {
	Name     string `json:"name" binding:"required"`
	SeatRows int    `json:"seat_rows" binding:"required,min=1"`
	SeatCols int    `json:"seat_cols" binding:"required,min=1"`
}

func (h *CatalogHandler) CreateHall(
	c *gin.Context,
) {
	var req createHallRequest

	if err := c.ShouldBindJSON(
		&req,
	); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	hall, err := h.service.CreateHall(
		c.Request.Context(),
		req.Name,
		req.SeatRows,
		req.SeatCols,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "internal error",
			},
		)

		return
	}

	c.JSON(
		http.StatusCreated,
		hall,
	)
}

// ====================
// Showtimes
// ====================

func (h *CatalogHandler) ListShowtimes(
	c *gin.Context,
) {
	movieID, err := primitive.ObjectIDFromHex(
		c.Param("id"),
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid movie id",
			},
		)

		return
	}

	showtimes, err := h.service.ListShowtimesByMovie(
		c.Request.Context(),
		movieID,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "internal error",
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		showtimes,
	)
}

type createShowtimeRequest struct {
	MovieID  string    `json:"movie_id" binding:"required"`
	HallID   string    `json:"hall_id" binding:"required"`
	StartsAt time.Time `json:"starts_at" binding:"required"`
	Price    float64   `json:"price" binding:"required"`
}

func (h *CatalogHandler) CreateShowtime(
	c *gin.Context,
) {
	var req createShowtimeRequest

	if err := c.ShouldBindJSON(
		&req,
	); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	movieID, err := primitive.ObjectIDFromHex(
		req.MovieID,
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid movie_id",
			},
		)

		return
	}

	hallID, err := primitive.ObjectIDFromHex(
		req.HallID,
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid hall_id",
			},
		)

		return
	}

	showtime, err := h.service.CreateShowtime(
		c.Request.Context(),
		movieID,
		hallID,
		req.StartsAt,
		req.Price,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "internal error",
			},
		)

		return
	}

	c.JSON(
		http.StatusCreated,
		showtime,
	)
}

// ====================
// Seat Map
// ====================

func (h *CatalogHandler) GetSeatMap(
	c *gin.Context,
) {
	showtimeID, err := primitive.ObjectIDFromHex(
		c.Param("id"),
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid showtime id",
			},
		)

		return
	}

	seats, err := h.service.GetSeatMap(
		c.Request.Context(),
		showtimeID,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "internal error",
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		seats,
	)
}