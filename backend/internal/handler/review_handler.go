package handler

import (
	"errors"
	"net/http"
	"offerpilot/backend/internal/pkg/response"
	"offerpilot/backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewService service.ReviewService
}

func NewReviewHandler(reviewService service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

func (h *ReviewHandler) GetReport(c *gin.Context) {
	userID, ok := getRequestUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
		return
	}

	sessionID, err := strconv.ParseUint(c.Param("sessionId"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40000, "invalid session id")
		return
	}

	result, err := h.reviewService.GetReport(userID, sessionID)
	if err != nil {
		if errors.Is(err, service.ErrReviewNotFound) {
			response.Error(c, http.StatusNotFound, 40421, "review report not found")
			return
		}

		response.Error(c, http.StatusInternalServerError, 50000, "failed to fetch review report")
		return
	}

	response.Success(c, result)
}

func (h *ReviewHandler) GetHistory(c *gin.Context) {
	userID, ok := getRequestUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
		return
	}

	result, err := h.reviewService.GetHistory(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50000, "failed to fetch review history")
		return
	}

	response.Success(c, result)
}
