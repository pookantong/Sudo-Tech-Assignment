package router

import (
	"net/http"

	"cinema-booking-backend/internal/config"
	"cinema-booking-backend/internal/handler"
	"cinema-booking-backend/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	cfg *config.Config,
	authHandler *handler.AuthHandler,
	bookingHandler *handler.BookingHandler,
	catalogHandler *handler.CatalogHandler,
) *gin.Engine {

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			cfg.FrontendURL,
		},
		AllowMethods: []string{
			"GET",
			"POST",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
	}))

	// Health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Auth
	r.POST(
		"/auth/google",
		authHandler.GoogleLogin,
	)

	// Public Catalog
	r.GET(
		"/movies",
		catalogHandler.ListMovies,
	)

	r.GET(
		"/movies/:id/showtimes",
		catalogHandler.ListShowtimes,
	)

	r.GET(
		"/showtimes/:id/seats",
		catalogHandler.GetSeatMap,
	)

	r.GET("/ws", bookingHandler.ServeWS)

	// Authenticated routes
	authed := r.Group("/")
	authed.Use(
		middleware.RequireAuth(cfg.JWTSecret),
	)

	authed.POST(
		"/seats/select",
		bookingHandler.SelectSeat,
	)

	authed.POST(
		"/bookings/confirm",
		bookingHandler.ConfirmPayment,
	)

	// Admin routes
	admin := authed.Group("/admin")

	admin.Use(
		middleware.RequireAdmin(),
	)

	admin.GET(
		"/bookings",
		bookingHandler.AdminListBookings,
	)

	admin.POST(
		"/movies",
		catalogHandler.CreateMovie,
	)

	admin.POST(
		"/showtimes",
		catalogHandler.CreateShowtime,
	)

	return r
}
