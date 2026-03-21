package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"offerpilot/backend/internal/model"
	"offerpilot/backend/internal/repository"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

var ErrInterviewSessionNotFound = errors.New("interview session not found")

type CreateInterviewSessionRequest struct {
	Mode    string                 `json:"mode"`
	JobRole string                 `json:"jobRole"`
	Config  map[string]interface{} `json:"config"`
}

type InterviewMessageItem struct {
	ID      string `json:"id"`
	Role    string `json:"role"`
	Content string `json:"content"`
	Time    string `json:"time"`
}

type InterviewQuestionItem struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

type InterviewStatusItem struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type InterviewScoreItem struct {
	Label string `json:"label"`
	Score int    `json:"score"`
}

type KnowledgeHitItem struct {
	Name    string `json:"name"`
	Level   string `json:"level"`
	Summary string `json:"summary"`
}

type InterviewSessionDetail struct {
	SessionID         string                  `json:"sessionId"`
	CurrentQuestionID string                  `json:"currentQuestionId"`
	Questions         []InterviewQuestionItem `json:"questions"`
	Messages          []InterviewMessageItem  `json:"messages"`
	StatusItems       []InterviewStatusItem   `json:"statusItems"`
	ScoreItems        []InterviewScoreItem    `json:"scoreItems"`
	KnowledgeHits     []KnowledgeHitItem      `json:"knowledgeHits"`
	AgentStage        string                  `json:"agentStage,omitempty"`
}

type CreateInterviewSessionResponse struct {
	SessionID string `json:"sessionId"`
}

type SubmitInterviewMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

type SubmitInterviewMessageResponse struct {
	Message           InterviewMessageItem    `json:"message"`
	CurrentQuestionID string                  `json:"currentQuestionId"`
	Questions         []InterviewQuestionItem `json:"questions"`
	StatusItems       []InterviewStatusItem   `json:"statusItems"`
	ScoreItems        []InterviewScoreItem    `json:"scoreItems"`
	KnowledgeHits     []KnowledgeHitItem      `json:"knowledgeHits"`
}

type StateEventPayload struct {
	CurrentQuestionID string                  `json:"currentQuestionId,omitempty"`
	Questions         []InterviewQuestionItem `json:"questions,omitempty"`
	StatusItems       []InterviewStatusItem   `json:"statusItems,omitempty"`
	ScoreItems        []InterviewScoreItem    `json:"scoreItems,omitempty"`
	AgentStage        string                  `json:"agentStage"`
}

type ReferenceEventPayload struct {
	KnowledgeHits []KnowledgeHitItem `json:"knowledgeHits"`
}

type DeltaEventPayload struct {
	Content string `json:"content"`
}

type DoneEventPayload struct {
	CurrentQuestionID string                  `json:"currentQuestionId,omitempty"`
	Questions         []InterviewQuestionItem `json:"questions,omitempty"`
	StatusItems       []InterviewStatusItem   `json:"statusItems,omitempty"`
	ScoreItems        []InterviewScoreItem    `json:"scoreItems,omitempty"`
	KnowledgeHits     []KnowledgeHitItem      `json:"knowledgeHits,omitempty"`
	AgentStage        string                  `json:"agentStage"`
	Message           *InterviewMessageItem   `json:"message,omitempty"`
}

type InterviewStreamEmitter func(eventType string, payload any) error

type InterviewService interface {
	CreateSession(userID uint64, req CreateInterviewSessionRequest) (*CreateInterviewSessionResponse, error)
	GetSessionDetail(userID, sessionID uint64) (*InterviewSessionDetail, error)
	SubmitMessage(userID, sessionID uint64, req SubmitInterviewMessageRequest) (*SubmitInterviewMessageResponse, error)
	StreamSession(ctx context.Context, userID, sessionID uint64, emit InterviewStreamEmitter) error
	FinishSession(userID, sessionID uint64) error
}

type interviewService struct {
	interviewRepo repository.InterviewRepository
}

func NewInterviewService(interviewRepo repository.InterviewRepository) InterviewService {
	return &interviewService{interviewRepo: interviewRepo}
}

func (s *interviewService) CreateSession(userID uint64, req CreateInterviewSessionRequest) (*CreateInterviewSessionResponse, error) {
	mode := req.Mode
	if mode == "" {
		mode = "general"
	}

	jobRole := req.JobRole
	if jobRole == "" {
		jobRole = "frontend engineer"
	}

	session := &model.InterviewSession{
		UserID:      userID,
		Mode:        mode,
		JobRole:     jobRole,
		Status:      "active",
		TotalRounds: 1,
		StartedAt:   time.Now(),
	}
	if err := s.interviewRepo.CreateSession(session); err != nil {
		return nil, err
	}

	firstQuestion := s.buildInitialQuestion(mode, jobRole)
	message := &model.InterviewMessage{
		SessionID:     session.ID,
		Role:          "interviewer",
		RoundNo:       1,
		QuestionID:    "q-1",
		Content:       firstQuestion,
		ScoreSnapshot: s.mustMarshalScoreSnapshot(78, 80, 76),
	}
	if err := s.interviewRepo.CreateMessage(message); err != nil {
		return nil, err
	}

	return &CreateInterviewSessionResponse{SessionID: strconv.FormatUint(session.ID, 10)}, nil
}

func (s *interviewService) GetSessionDetail(userID, sessionID uint64) (*InterviewSessionDetail, error) {
	session, err := s.interviewRepo.GetSessionByID(sessionID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInterviewSessionNotFound
		}
		return nil, err
	}

	messages, err := s.interviewRepo.ListMessages(session.ID)
	if err != nil {
		return nil, err
	}

	detail := s.buildSessionDetail(session, messages)
	detail.AgentStage = "session_created"
	return detail, nil
}

func (s *interviewService) SubmitMessage(userID, sessionID uint64, req SubmitInterviewMessageRequest) (*SubmitInterviewMessageResponse, error) {
	session, err := s.interviewRepo.GetSessionByID(sessionID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInterviewSessionNotFound
		}
		return nil, err
	}

	latestRound, err := s.interviewRepo.GetLatestRoundNo(session.ID)
	if err != nil {
		return nil, err
	}

	candidateMessage := &model.InterviewMessage{
		SessionID:     session.ID,
		Role:          "candidate",
		RoundNo:       latestRound,
		QuestionID:    fmt.Sprintf("q-%d", latestRound),
		Content:       req.Content,
		ScoreSnapshot: s.mustMarshalScoreSnapshot(80, 79, 78),
	}
	if err := s.interviewRepo.CreateMessage(candidateMessage); err != nil {
		return nil, err
	}

	nextRound := latestRound + 1
	interviewerReply := &model.InterviewMessage{
		SessionID:     session.ID,
		Role:          "interviewer",
		RoundNo:       nextRound,
		QuestionID:    fmt.Sprintf("q-%d", nextRound),
		Content:       s.buildFollowupQuestion(req.Content, nextRound),
		ScoreSnapshot: s.mustMarshalScoreSnapshot(82, 81, 80),
	}
	if err := s.interviewRepo.CreateMessage(interviewerReply); err != nil {
		return nil, err
	}

	session.TotalRounds = nextRound
	if err := s.interviewRepo.UpdateSession(session); err != nil {
		return nil, err
	}

	messages, err := s.interviewRepo.ListMessages(session.ID)
	if err != nil {
		return nil, err
	}

	detail := s.buildSessionDetail(session, messages)
	return &SubmitInterviewMessageResponse{
		Message:           toMessageItem(*interviewerReply),
		CurrentQuestionID: detail.CurrentQuestionID,
		Questions:         detail.Questions,
		StatusItems:       detail.StatusItems,
		ScoreItems:        detail.ScoreItems,
		KnowledgeHits:     detail.KnowledgeHits,
	}, nil
}

func (s *interviewService) StreamSession(ctx context.Context, userID, sessionID uint64, emit InterviewStreamEmitter) error {
	session, err := s.interviewRepo.GetSessionByID(sessionID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInterviewSessionNotFound
		}
		return err
	}

	messages, err := s.interviewRepo.ListMessages(session.ID)
	if err != nil {
		return err
	}
	if len(messages) == 0 {
		return ErrInterviewSessionNotFound
	}

	lastMessage := messages[len(messages)-1]
	detail := s.buildSessionDetail(session, messages)
	references := s.buildKnowledgeHits()
	scores := s.buildScoreItems()

	steps := []func() error{
		func() error {
			return emit("state", StateEventPayload{
				CurrentQuestionID: detail.CurrentQuestionID,
				Questions:         detail.Questions,
				StatusItems:       detail.StatusItems,
				AgentStage:        "session_created",
			})
		},
		func() error {
			return emit("state", StateEventPayload{
				CurrentQuestionID: detail.CurrentQuestionID,
				Questions:         detail.Questions,
				StatusItems:       s.buildStreamingStatusItems(session, "retrieving_knowledge"),
				AgentStage:        "retrieving_knowledge",
			})
		},
		func() error {
			return emit("reference", ReferenceEventPayload{KnowledgeHits: references})
		},
		func() error {
			return emit("state", StateEventPayload{
				CurrentQuestionID: detail.CurrentQuestionID,
				Questions:         detail.Questions,
				StatusItems:       s.buildStreamingStatusItems(session, "scoring"),
				ScoreItems:        scores,
				AgentStage:        "scoring",
			})
		},
	}

	for _, step := range steps {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := step(); err != nil {
			return err
		}
		time.Sleep(250 * time.Millisecond)
	}

	if err := emit("state", StateEventPayload{
		CurrentQuestionID: detail.CurrentQuestionID,
		Questions:         detail.Questions,
		StatusItems:       s.buildStreamingStatusItems(session, "generating_followup"),
		ScoreItems:        scores,
		AgentStage:        "generating_followup",
	}); err != nil {
		return err
	}
	time.Sleep(250 * time.Millisecond)

	for _, chunk := range s.splitTextForStream(lastMessage.Content) {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := emit("delta", DeltaEventPayload{Content: chunk}); err != nil {
			return err
		}
		time.Sleep(180 * time.Millisecond)
	}

	finalStatus := s.buildStreamingStatusItems(session, "completed")
	if err := emit("state", StateEventPayload{
		CurrentQuestionID: detail.CurrentQuestionID,
		Questions:         detail.Questions,
		StatusItems:       finalStatus,
		ScoreItems:        scores,
		AgentStage:        "completed",
	}); err != nil {
		return err
	}

	return emit("done", DoneEventPayload{
		CurrentQuestionID: detail.CurrentQuestionID,
		Questions:         detail.Questions,
		StatusItems:       finalStatus,
		ScoreItems:        scores,
		KnowledgeHits:     references,
		AgentStage:        "completed",
		Message:           pointerToMessageItem(toMessageItem(lastMessage)),
	})
}

func (s *interviewService) FinishSession(userID, sessionID uint64) error {
	session, err := s.interviewRepo.GetSessionByID(sessionID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInterviewSessionNotFound
		}
		return err
	}

	now := time.Now()
	session.Status = "finished"
	session.EndedAt = &now
	return s.interviewRepo.UpdateSession(session)
}

func (s *interviewService) buildSessionDetail(session *model.InterviewSession, messages []model.InterviewMessage) *InterviewSessionDetail {
	currentQuestionID := fmt.Sprintf("q-%d", session.TotalRounds)
	questions := make([]InterviewQuestionItem, 0, session.TotalRounds)
	for i := 1; i <= max(session.TotalRounds, 1); i++ {
		status := "pending"
		if i < session.TotalRounds || session.Status == "finished" {
			status = "completed"
		} else {
			status = "in_progress"
		}

		questions = append(questions, InterviewQuestionItem{
			ID:     fmt.Sprintf("q-%d", i),
			Title:  fmt.Sprintf("Round %d", i),
			Status: status,
		})
	}

	messageItems := make([]InterviewMessageItem, 0, len(messages))
	for _, message := range messages {
		messageItems = append(messageItems, toMessageItem(message))
	}

	return &InterviewSessionDetail{
		SessionID:         strconv.FormatUint(session.ID, 10),
		CurrentQuestionID: currentQuestionID,
		Questions:         questions,
		Messages:          messageItems,
		StatusItems:       s.buildStatusItems(session),
		ScoreItems:        s.buildScoreItems(),
		KnowledgeHits:     s.buildKnowledgeHits(),
	}
}

func (s *interviewService) buildInitialQuestion(mode, jobRole string) string {
	return fmt.Sprintf("You are in %s mode for a %s role. Start with a concise self introduction and highlight your strongest frontend capability.", mode, jobRole)
}

func (s *interviewService) buildFollowupQuestion(answer string, round int) string {
	preview := []rune(strings.TrimSpace(answer))
	if len(preview) > 40 {
		preview = preview[:40]
	}
	return fmt.Sprintf("You mentioned \"%s\". Please explain the technical decision, the implementation steps, and the measurable result behind it for round %d.", string(preview), round)
}

func (s *interviewService) buildStatusItems(session *model.InterviewSession) []InterviewStatusItem {
	statusText := "in_progress"
	if session.Status == "finished" {
		statusText = "completed"
	}
	return []InterviewStatusItem{
		{Label: "Session", Value: statusText},
		{Label: "Round", Value: fmt.Sprintf("%d", max(session.TotalRounds, 1))},
		{Label: "Mode", Value: session.Mode},
		{Label: "Role", Value: session.JobRole},
	}
}

func (s *interviewService) buildStreamingStatusItems(session *model.InterviewSession, stage string) []InterviewStatusItem {
	items := s.buildStatusItems(session)
	return append([]InterviewStatusItem{{Label: "Agent Stage", Value: stage}}, items...)
}

func (s *interviewService) buildScoreItems() []InterviewScoreItem {
	return []InterviewScoreItem{
		{Label: "Expression", Score: 82},
		{Label: "Accuracy", Score: 80},
		{Label: "Structure", Score: 79},
	}
}

func (s *interviewService) buildKnowledgeHits() []KnowledgeHitItem {
	return []KnowledgeHitItem{
		{Name: "project communication", Level: "high", Summary: "The answer ties the response back to concrete project ownership and impact."},
		{Name: "technical tradeoff", Level: "medium", Summary: "The candidate mentions the approach, but could quantify why it was chosen."},
		{Name: "result validation", Level: "low", Summary: "The answer still needs clearer metrics and final outcome verification."},
	}
}

func (s *interviewService) splitTextForStream(text string) []string {
	runes := []rune(strings.TrimSpace(text))
	if len(runes) == 0 {
		return []string{"No interviewer feedback generated."}
	}

	chunkSize := 18
	chunks := make([]string, 0, (len(runes)/chunkSize)+1)
	for start := 0; start < len(runes); start += chunkSize {
		end := start + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[start:end]))
	}
	return chunks
}

func (s *interviewService) mustMarshalScoreSnapshot(expression, accuracy, structure int) string {
	payload, _ := json.Marshal(map[string]int{
		"expression": expression,
		"accuracy":   accuracy,
		"structure":  structure,
	})
	return string(payload)
}

func toMessageItem(message model.InterviewMessage) InterviewMessageItem {
	createdAt := message.CreatedAt
	if createdAt == 0 {
		createdAt = time.Now().UnixMilli()
	}
	return InterviewMessageItem{
		ID:      strconv.FormatUint(message.ID, 10),
		Role:    message.Role,
		Content: message.Content,
		Time:    time.UnixMilli(createdAt).Format("15:04"),
	}
}

func pointerToMessageItem(item InterviewMessageItem) *InterviewMessageItem {
	return &item
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}