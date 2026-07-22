package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"cinema-booking-backend/internal/booking"
	"cinema-booking-backend/internal/cache"
	"cinema-booking-backend/internal/middleware"
	"cinema-booking-backend/internal/ws"
)

type BookingHandler struct {
	service  *booking.Service
	hub      *ws.Hub
	upgrader websocket.Upgrader
}

func NewBookingHandler(
	service *booking.Service,
	hub *ws.Hub,
) *BookingHandler {
	return &BookingHandler{
		service: service,
		hub:     hub,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(
				r *http.Request,
			) bool {
				return true
			},
		},
	}
}

// ====================
// Select Seat
// ====================

type selectSeatRequest struct {
	ShowtimeID string `json:"showtime_id" binding:"required"`
	SeatID     string `json:"seat_id" binding:"required"`
}

// POST /seats/select
func (h *BookingHandler) SelectSeat(
	c *gin.Context,
) {
	var req selectSeatRequest

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

	userID, err := primitive.ObjectIDFromHex(
		middleware.UserIDFromContext(c),
	)

	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid user",
			},
		)

		return
	}

	showtimeID, err := primitive.ObjectIDFromHex(
		req.ShowtimeID,
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid showtime_id",
			},
		)

		return
	}

	seatID, err := primitive.ObjectIDFromHex(
		req.SeatID,
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid seat_id",
			},
		)

		return
	}

	b, err := h.service.SelectSeat(
		c.Request.Context(),
		userID,
		showtimeID,
		seatID,
	)

	if err != nil {
		switch {
		case errors.Is(
			err,
			booking.ErrSeatTaken,
		):
			c.JSON(
				http.StatusConflict,
				gin.H{
					"error": "seat already locked or booked",
				},
			)

		case errors.Is(
			err,
			booking.ErrSeatNotFound,
		):
			c.JSON(
				http.StatusNotFound,
				gin.H{
					"error": "seat not found",
				},
			)

		case errors.Is(
			err,
			booking.ErrSeatNotInShowtime,
		):
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"error": "seat does not belong to this showtime",
				},
			)

		default:
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"error": "internal error",
				},
			)
		}

		return
	}

	c.JSON(
		http.StatusCreated,
		b,
	)
}

// ====================
// Payment
// ====================

type confirmPaymentRequest struct {
	BookingID string `json:"booking_id" binding:"required"`
	Outcome   string `json:"outcome"`
}

// POST /bookings/confirm
func (h *BookingHandler) ConfirmPayment(
	c *gin.Context,
) {
	var req confirmPaymentRequest

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

	bookingID, err := primitive.ObjectIDFromHex(
		req.BookingID,
	)

	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid booking_id",
			},
		)

		return
	}

	userID, err := primitive.ObjectIDFromHex(
		middleware.UserIDFromContext(c),
	)

	if err != nil {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid user",
			},
		)

		return
	}

	var opErr error

	if req.Outcome == "fail" {
		opErr = h.service.FailPayment(
			c.Request.Context(),
			bookingID,
			userID,
		)
	} else {
		opErr = h.service.ConfirmPayment(
			c.Request.Context(),
			bookingID,
			userID,
		)
	}

	switch {
	case opErr == nil:

	case errors.Is(
		opErr,
		booking.ErrBookingNotFound,
	):
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "booking not found",
			},
		)

		return

	case errors.Is(
		opErr,
		booking.ErrNotOwner,
	):
		c.JSON(
			http.StatusForbidden,
			gin.H{
				"error": "not your booking",
			},
		)

		return

	case errors.Is(
		opErr,
		cache.ErrLockNotOwned,
	):
		c.JSON(
			http.StatusConflict,
			gin.H{
				"error": "lock expired or already released, please retry seat selection",
			},
		)

		return

	default:
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": opErr.Error(),
			},
		)

		return
	}
	if req.Outcome == "fail" {
		c.JSON(
			http.StatusOK,
			gin.H{
				"status": "payment_failed",
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status": "confirmed",
		},
	)
}

// ====================
// Admin
// ====================

// GET /admin/bookings
func (h *BookingHandler) AdminListBookings(
	c *gin.Context,
) {
	filter := booking.AdminBookingFilter{
		UserID:     c.Query("user"),
		ShowtimeID: c.Query("showtime"),
		Status:     c.Query("status"),
	}

	results, err := h.service.ListForAdmin(
		c.Request.Context(),
		filter,
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
		results,
	)
}

// ====================
// WebSocket
// ====================

// GET /ws
func (h *BookingHandler) ServeWS(
	c *gin.Context,
) {
	conn, err := h.upgrader.Upgrade(
		c.Writer,
		c.Request,
		nil,
	)

	if err != nil {
		return
	}

	h.hub.Register(conn)

	go func() {
		defer h.hub.Unregister(conn)

		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()
}
