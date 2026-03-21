package handler

import (
	"errors"
	"net/http"
	"offerpilot/backend/internal/middleware"
	"offerpilot/backend/internal/pkg/response"
	"offerpilot/backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QuestionHandler struct {
	questionService service.QuestionService
}

func NewQuestionHandler(questionService service.QuestionService) *QuestionHandler {
	return &QuestionHandler{questionService: questionService}
}

func (h *QuestionHandler) List(c *gin.Context) {
	userID, ok := getRequestUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	result, err := h.questionService.List(userID, service.QuestionListParams{
		Keyword:   c.Query("keyword"),
		Category:  c.Query("category"),
		Frequency: c.Query("frequency"),
		Page:      page,
		PageSize:  pageSize,
	})
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50000, "failed to fetch questions")
		return
	}

	response.Success(c, result)
}

func (h *QuestionHandler) Detail(c *gin.Context) {
	userID, ok := getRequestUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
		return
	}

	questionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40000, "invalid question id")
		return
	}

	result, err := h.questionService.GetDetail(userID, questionID)
	if err != nil {
		if errors.Is(err, service.ErrQuestionNotFound) {
			response.Error(c, http.StatusNotFound, 40401, "question not found")
			return
		}

		response.Error(c, http.StatusInternalServerError, 50000, "failed to fetch question detail")
		return
	}

	response.Success(c, result)
}

func (h *QuestionHandler) AddFavorite(c *gin.Context) {
	userID, ok := getRequestUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
		return
	}

	questionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40000, "invalid question id")
		return
	}

	result, err := h.questionService.AddFavorite(userID, questionID)
	if err != nil {
		if errors.Is(err, service.ErrQuestionNotFound) {
			response.Error(c, http.StatusNotFound, 40401, "question not found")
			return
		}

		response.Error(c, http.StatusInternalServerError, 50000, "failed to add favorite")
		return
	}

	response.Success(c, result)
}

func (h *QuestionHandler) RemoveFavorite(c *gin.Context) {
	userID, ok := getRequestUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
		return
	}

	questionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40000, "invalid question id")
		return
	}

	result, err := h.questionService.RemoveFavorite(userID, questionID)
	if err != nil {
		if errors.Is(err, service.ErrQuestionNotFound) {
			response.Error(c, http.StatusNotFound, 40401, "question not found")
			return
		}

		response.Error(c, http.StatusInternalServerError, 50000, "failed to remove favorite")
		return
	}

	response.Success(c, result)
}

func getRequestUserID(c *gin.Context) (uint64, bool) {
	value, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		return 0, false
	}

	userID, ok := value.(uint64)
	return userID, ok
}
