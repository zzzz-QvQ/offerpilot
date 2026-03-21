package handler

import (
	"errors"
	"net/http"
	"offerpilot/backend/internal/middleware"
	"offerpilot/backend/internal/pkg/response"
	"offerpilot/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40000, "invalid request parameters")
		return
	}

	result, err := h.authService.Login(req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			response.Error(c, http.StatusUnauthorized, 40101, "email or password is incorrect")
			return
		}

		response.Error(c, http.StatusInternalServerError, 50000, "login failed")
		return
	}

	response.Success(c, result)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userIDValue, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
		return
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
		return
	}

	profile, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
			return
		}

		response.Error(c, http.StatusInternalServerError, 50000, "failed to fetch current user")
		return
	}

	response.Success(c, profile)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userIDValue, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
		return
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
		return
	}

	if err := h.authService.Logout(userID); err != nil {
		response.Error(c, http.StatusInternalServerError, 50000, "logout failed")
		return
	}

	response.Success(c, gin.H{})
}
