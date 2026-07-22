package bootstrap

import (
	"context"
	"time"

	"cinema-booking-backend/internal/auth"
	"cinema-booking-backend/internal/booking"
	"cinema-booking-backend/internal/catalog"
	"cinema-booking-backend/internal/config"
	"cinema-booking-backend/internal/handler"
	"cinema-booking-backend/internal/repository"
	"cinema-booking-backend/internal/router"
	"cinema-booking-backend/internal/ws"

	"github.com/gin-gonic/gin"
)

type App struct {
	Router        *gin.Engine
	CancelWorkers context.CancelFunc
}

func NewApp(cfg *config.Config) *App {
	// Infrastructure
	mongoClient, db := NewMongo(cfg)
	redisClient, seatLock := NewRedis(cfg)
	producer, consumer := NewKafka(cfg)

	// Repository
	bookingRepo := repository.NewBookingRepository(db)
	userRepo := repository.NewUserRepository(db)
	catalogRepo := repository.NewCatalogRepository(db)

	// WebSocket
	hub := ws.NewHub()

	timeoutWatcher := booking.NewTimeoutWatcher(
		redisClient,
		mongoClient,
		bookingRepo,
		producer,
		hub,
	)

	googleVerifier := auth.NewGoogleVerifier(cfg.GoogleClientID)

	tokenIssuer := auth.NewTokenIssuer(
		cfg.JWTSecret,
		24*time.Hour,
	)

	// Services
	bookingService := booking.NewService(
		seatLock,
		bookingRepo,
		catalogRepo,
		producer,
		hub,
	)

	authService := auth.NewService(
		googleVerifier,
		tokenIssuer,
		userRepo,
	)

	catalogService := catalog.NewService(
		catalogRepo,
		bookingRepo,
		seatLock,
	)

	// Handlers
	bookingHandler := handler.NewBookingHandler(
		bookingService,
		hub,
	)
	authHandler := handler.NewAuthHandler(
		authService,
	)

	catalogHandler := handler.NewCatalogHandler(
		catalogService,
	)

	// Background Workers
	bgCtx, cancel := context.WithCancel(context.Background())

	RunWorkers(
		bgCtx,
		timeoutWatcher,
		consumer,
	)

	// HTTP Router
	r := router.NewRouter(
		cfg,
		authHandler,
		bookingHandler,
		catalogHandler,
	)

	return &App{
		Router:        r,
		CancelWorkers: cancel,
	}
}
