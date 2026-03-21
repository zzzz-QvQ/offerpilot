package handler

import (
	"errors"
	"net/http"
	"offerpilot/backend/internal/pkg/response"
	"offerpilot/backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projectService service.ProjectService
}

func NewProjectHandler(projectService service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

func (h *ProjectHandler) Create(c *gin.Context) {
	userID, ok := getRequestUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
		return
	}

	var req service.CreateProjectPolishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40000, "invalid request payload")
		return
	}

	result, err := h.projectService.Create(userID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50000, "failed to generate project polish result")
		return
	}

	response.Success(c, result)
}

func (h *ProjectHandler) GetHistory(c *gin.Context) {
	userID, ok := getRequestUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
		return
	}

	result, err := h.projectService.GetHistory(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50000, "failed to fetch project polish history")
		return
	}

	response.Success(c, result)
}

func (h *ProjectHandler) GetDetail(c *gin.Context) {
	userID, ok := getRequestUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
		return
	}

	recordID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40000, "invalid project polish id")
		return
	}

	result, err := h.projectService.GetDetail(userID, recordID)
	if err != nil {
		if errors.Is(err, service.ErrProjectRecordNotFound) {
			response.Error(c, http.StatusNotFound, 40431, "project polish record not found")
			return
		}

		response.Error(c, http.StatusInternalServerError, 50000, "failed to fetch project polish detail")
		return
	}

	response.Success(c, result)
}
