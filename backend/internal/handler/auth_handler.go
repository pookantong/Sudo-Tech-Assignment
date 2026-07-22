package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cinema-booking-backend/internal/auth"
)

type AuthHandler struct {
	authService *auth.Service
}

func NewAuthHandler(
	authService *auth.Service,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type googleLoginRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req googleLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := h.authService.GoogleLogin(
		c.Request.Context(),
		req.IDToken,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": result.Token,
		"user":  result.User,
	})
}