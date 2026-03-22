package handler

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"offerpilot/backend/internal/pkg/response"
	"offerpilot/backend/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type InterviewHandler struct {
	interviewService service.InterviewService
}

func NewInterviewHandler(interviewService service.InterviewService) *InterviewHandler {
	return &InterviewHandler{interviewService: interviewService}
}

func (h *InterviewHandler) CreateSession(c *gin.Context) {
	userID, ok := getRequestUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
		return
	}

	var req service.CreateInterviewSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil && c.ContentType() == "application/json" {
		response.Error(c, http.StatusBadRequest, 40000, "invalid request parameters")
		return
	}

	result, err := h.interviewService.CreateSession(userID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50000, "failed to create interview session: "+err.Error())
		return
	}

	response.Success(c, result)
}

func (h *InterviewHandler) GetSession(c *gin.Context) {
	userID, ok := getRequestUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
		return
	}

	sessionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40000, "invalid session id")
		return
	}

	result, err := h.interviewService.GetSessionDetail(userID, sessionID)
	if err != nil {
		if errors.Is(err, service.ErrInterviewSessionNotFound) {
			response.Error(c, http.StatusNotFound, 40411, "interview session not found")
			return
		}

		response.Error(c, http.StatusInternalServerError, 50000, "failed to fetch interview session")
		return
	}

	response.Success(c, result)
}

func (h *InterviewHandler) SubmitMessage(c *gin.Context) {
	userID, ok := getRequestUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
		return
	}

	sessionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40000, "invalid session id")
		return
	}

	var req service.SubmitInterviewMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40000, "invalid request parameters")
		return
	}

	result, err := h.interviewService.SubmitMessage(userID, sessionID, req)
	if err != nil {
		if errors.Is(err, service.ErrInterviewSessionNotFound) {
			response.Error(c, http.StatusNotFound, 40411, "interview session not found")
			return
		}

		response.Error(c, http.StatusInternalServerError, 50000, "failed to submit interview message")
		return
	}

	response.Success(c, result)
}

func (h *InterviewHandler) StreamSession(c *gin.Context) {
	userID, ok := getRequestUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
		return
	}

	sessionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40000, "invalid session id")
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		response.Error(c, http.StatusInternalServerError, 50000, "streaming unsupported")
		return
	}

	writer := bufio.NewWriter(c.Writer)
	emit := func(eventType string, payload any) error {
		bytes, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(writer, "event: %s\n", eventType); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(writer, "data: %s\n\n", string(bytes)); err != nil {
			return err
		}
		if err := writer.Flush(); err != nil {
			return err
		}
		flusher.Flush()
		return nil
	}

	if err := h.interviewService.StreamSession(c.Request.Context(), userID, sessionID, emit); err != nil {
		if errors.Is(err, service.ErrInterviewSessionNotFound) {
			_ = emit("done", gin.H{"agentStage": "completed"})
			return
		}
		_ = emit("done", gin.H{"agentStage": "completed"})
	}
}

func (h *InterviewHandler) FinishSession(c *gin.Context) {
	userID, ok := getRequestUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, 40100, "unauthorized")
		return
	}

	sessionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40000, "invalid session id")
		return
	}

	if err := h.interviewService.FinishSession(userID, sessionID); err != nil {
		if errors.Is(err, service.ErrInterviewSessionNotFound) {
			response.Error(c, http.StatusNotFound, 40411, "interview session not found")
			return
		}

		response.Error(c, http.StatusInternalServerError, 50000, "failed to finish interview session")
		return
	}

	response.Success(c, gin.H{})
}
