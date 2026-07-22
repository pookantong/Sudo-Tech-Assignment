package handler

import (
	"net/http"

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