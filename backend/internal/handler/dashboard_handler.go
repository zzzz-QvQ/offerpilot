package handler

import (
	"net/http"
	"offerpilot/backend/internal/middleware"
	"offerpilot/backend/internal/pkg/response"
	"offerpilot/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardService service.DashboardService
}

func NewDashboardHandler(dashboardService service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

func (h *DashboardHandler) GetSummary(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
		return
	}

	result, err := h.dashboardService.GetSummary(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50000, "failed to fetch dashboard summary")
		return
	}

	response.Success(c, result)
}

func (h *DashboardHandler) GetRecentSessions(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
		return
	}

	result, err := h.dashboardService.GetRecentSessions(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50000, "failed to fetch recent sessions")
		return
	}

	response.Success(c, result)
}

func (h *DashboardHandler) GetWeakPoints(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
		return
	}

	result, err := h.dashboardService.GetWeakPoints(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50000, "failed to fetch weak points")
		return
	}

	response.Success(c, result)
}

func getCurrentUserID(c *gin.Context) (uint64, bool) {
	userIDValue, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		return 0, false
	}

	userID, ok := userIDValue.(uint64)
	return userID, ok
}
