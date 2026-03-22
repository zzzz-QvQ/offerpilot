package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"offerpilot/backend/internal/model"
	"offerpilot/backend/internal/pkg/vectorstore"
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
	QuestionID string  `json:"questionId"`
	Title      string  `json:"title"`
	Category   string  `json:"category"`
	Score      float64 `json:"score"`
	Snippet    string  `json:"snippet"`
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
	Accepted bool `json:"accepted"`
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

type interviewEvaluationResult struct {
	Summary         string `json:"summary"`
	FollowupFocus   string `json:"followupFocus"`
	ExpressionScore int    `json:"expressionScore"`
	AccuracyScore   int    `json:"accuracyScore"`
	StructureScore  int    `json:"structureScore"`
}

type interviewService struct {
	interviewRepo     repository.InterviewRepository
	llmService        LLMService
	questionRetrieval QuestionRetrievalService
}

func NewInterviewService(
	interviewRepo repository.InterviewRepository,
	llmService LLMService,
	questionRetrieval QuestionRetrievalService,
) InterviewService {
	return &interviewService{
		interviewRepo:     interviewRepo,
		llmService:        llmService,
		questionRetrieval: questionRetrieval,
	}
}

func (s *interviewService) CreateSession(userID uint64, req CreateInterviewSessionRequest) (*CreateInterviewSessionResponse, error) {
	mode := strings.TrimSpace(req.Mode)
	if mode == "" {
		mode = "general"
	}

	jobRole := strings.TrimSpace(req.JobRole)
	if jobRole == "" {
		jobRole = "frontend engineer"
	}

	firstQuestion, err := s.generateInitialQuestion(context.Background(), mode, jobRole)
	if err != nil {
		firstQuestion = s.buildFallbackInitialQuestion(mode, jobRole)
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

	message := &model.InterviewMessage{
		SessionID:     session.ID,
		Role:          "interviewer",
		RoundNo:       1,
		QuestionID:    "q-1",
		Content:       firstQuestion,
		ScoreSnapshot: s.mustMarshalScoreSnapshot(0, 0, 0),
	}
	if err := s.interviewRepo.CreateMessage(message); err != nil {
		return nil, err
	}

	return &CreateInterviewSessionResponse{SessionID: strconv.FormatUint(session.ID, 10)}, nil
}

func (s *interviewService) buildFallbackInitialQuestion(mode, jobRole string) string {
	mode = strings.TrimSpace(mode)
	jobRole = strings.TrimSpace(jobRole)

	if mode == "" {
		mode = "general"
	}
	if jobRole == "" {
		jobRole = "frontend engineer"
	}

	return fmt.Sprintf(
		"我们先进行一轮 %s 模式的 %s 面试。请你先做一个 1 分钟左右的自我介绍，并结合一个你最熟悉的前端项目，说明你的技术选型、核心职责以及遇到的一个主要难点。",
		mode,
		jobRole,
	)
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
	if latestRound <= 0 {
		latestRound = 1
	}

	candidateMessage := &model.InterviewMessage{
		SessionID:     session.ID,
		Role:          "candidate",
		RoundNo:       latestRound,
		QuestionID:    fmt.Sprintf("q-%d", latestRound),
		Content:       strings.TrimSpace(req.Content),
		ScoreSnapshot: s.mustMarshalScoreSnapshot(0, 0, 0),
	}
	if err := s.interviewRepo.CreateMessage(candidateMessage); err != nil {
		return nil, err
	}

	return &SubmitInterviewMessageResponse{Accepted: true}, nil
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

	candidateMessage, currentQuestionMessage, err := s.findRoundContext(messages)
	if err != nil {
		return err
	}

	baseDetail := s.buildSessionDetail(session, messages)
	if err := emit("state", StateEventPayload{
		CurrentQuestionID: baseDetail.CurrentQuestionID,
		Questions:         baseDetail.Questions,
		StatusItems:       s.buildStreamingStatusItems(session, "generating_question"),
		AgentStage:        "generating_question",
	}); err != nil {
		return err
	}

	if err := emit("state", StateEventPayload{
		CurrentQuestionID: baseDetail.CurrentQuestionID,
		Questions:         baseDetail.Questions,
		StatusItems:       s.buildStreamingStatusItems(session, "retrieving_knowledge"),
		AgentStage:        "retrieving_knowledge",
	}); err != nil {
		return err
	}

	knowledgeHits, ragContext, err := s.retrieveKnowledge(ctx, currentQuestionMessage.Content, candidateMessage.Content, session.JobRole)
	if err != nil {
		return err
	}
	if err := emit("reference", ReferenceEventPayload{KnowledgeHits: knowledgeHits}); err != nil {
		return err
	}

	evaluation, err := s.evaluateAnswer(ctx, session, currentQuestionMessage.Content, candidateMessage.Content, ragContext)
	if err != nil {
		return err
	}
	scoreItems := s.buildScoreItemsFromEvaluation(evaluation)

	if err := emit("state", StateEventPayload{
		CurrentQuestionID: baseDetail.CurrentQuestionID,
		Questions:         baseDetail.Questions,
		StatusItems:       s.buildStreamingStatusItems(session, "scoring"),
		ScoreItems:        scoreItems,
		AgentStage:        "scoring",
	}); err != nil {
		return err
	}

	if err := emit("state", StateEventPayload{
		CurrentQuestionID: baseDetail.CurrentQuestionID,
		Questions:         baseDetail.Questions,
		StatusItems:       s.buildStreamingStatusItems(session, "generating_followup"),
		ScoreItems:        scoreItems,
		AgentStage:        "generating_followup",
	}); err != nil {
		return err
	}

	nextRound := session.TotalRounds + 1
	questionID := fmt.Sprintf("q-%d", nextRound)
	var builder strings.Builder
	err = s.llmService.StreamText(ctx, StreamTextRequest{
		SystemPrompt: buildInterviewFollowupSystemPrompt(),
		Messages: []LLMMessage{
			{Role: "user", Content: buildInterviewFollowupUserPrompt(session.Mode, session.JobRole, currentQuestionMessage.Content, candidateMessage.Content, ragContext, evaluation)},
		},
	}, func(chunk StreamTextChunk) error {
		if chunk.Delta == "" {
			return nil
		}
		builder.WriteString(chunk.Delta)
		return emit("delta", DeltaEventPayload{Content: chunk.Delta})
	})
	if err != nil {
		return err
	}

	finalReply := strings.TrimSpace(builder.String())
	if finalReply == "" {
		return errors.New("llm returned empty follow-up response")
	}

	interviewerMessage := &model.InterviewMessage{
		SessionID:     session.ID,
		Role:          "interviewer",
		RoundNo:       nextRound,
		QuestionID:    questionID,
		Content:       finalReply,
		ScoreSnapshot: s.mustMarshalScoreSnapshot(evaluation.ExpressionScore, evaluation.AccuracyScore, evaluation.StructureScore),
	}
	if err := s.interviewRepo.CreateMessage(interviewerMessage); err != nil {
		return err
	}

	session.TotalRounds = nextRound
	if err := s.interviewRepo.UpdateSession(session); err != nil {
		return err
	}

	updatedMessages, err := s.interviewRepo.ListMessages(session.ID)
	if err != nil {
		return err
	}
	finalDetail := s.buildSessionDetail(session, updatedMessages)
	finalStatusItems := s.buildStreamingStatusItems(session, "completed")

	if err := emit("state", StateEventPayload{
		CurrentQuestionID: finalDetail.CurrentQuestionID,
		Questions:         finalDetail.Questions,
		StatusItems:       finalStatusItems,
		ScoreItems:        scoreItems,
		AgentStage:        "completed",
	}); err != nil {
		return err
	}

	return emit("done", DoneEventPayload{
		CurrentQuestionID: finalDetail.CurrentQuestionID,
		Questions:         finalDetail.Questions,
		StatusItems:       finalStatusItems,
		ScoreItems:        scoreItems,
		KnowledgeHits:     knowledgeHits,
		AgentStage:        "completed",
		Message:           pointerToMessageItem(toMessageItem(*interviewerMessage)),
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

func (s *interviewService) generateInitialQuestion(ctx context.Context, mode, jobRole string) (string, error) {
	response, err := s.llmService.GenerateText(ctx, GenerateTextRequest{
		SystemPrompt: buildInterviewInitialQuestionSystemPrompt(),
		Messages: []LLMMessage{{
			Role:    "user",
			Content: buildInterviewInitialQuestionUserPrompt(mode, jobRole),
		}},
	})
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(response.Content), nil
}

func (s *interviewService) retrieveKnowledge(ctx context.Context, currentQuestion, answer, jobRole string) ([]KnowledgeHitItem, string, error) {
	if s.questionRetrieval == nil {
		return []KnowledgeHitItem{}, "", nil
	}

	query := strings.TrimSpace(strings.Join([]string{currentQuestion, answer, jobRole}, "\n"))
	hits, err := s.questionRetrieval.SearchQuestions(ctx, query, 3)
	if err != nil {
		if isRetrievalUnavailable(err) {
			return []KnowledgeHitItem{}, "", nil
		}
		return nil, "", err
	}

	knowledgeHits := make([]KnowledgeHitItem, 0, len(hits))
	contextLines := make([]string, 0, len(hits))
	for idx, hit := range hits {
		knowledgeHits = append(knowledgeHits, KnowledgeHitItem{
			QuestionID: hit.QuestionID,
			Title:      hit.Title,
			Category:   hit.Category,
			Score:      hit.Score,
			Snippet:    summarizeText(hit.Text, 180),
		})
		contextLines = append(contextLines, fmt.Sprintf("[%d] %s | %s\n%s", idx+1, hit.Title, hit.Category, hit.Text))
	}

	return knowledgeHits, strings.Join(contextLines, "\n\n"), nil
}

func (s *interviewService) evaluateAnswer(ctx context.Context, session *model.InterviewSession, currentQuestion, answer, ragContext string) (*interviewEvaluationResult, error) {
	response, err := s.llmService.GenerateText(ctx, GenerateTextRequest{
		SystemPrompt: buildInterviewEvaluationSystemPrompt(),
		Messages: []LLMMessage{{
			Role:    "user",
			Content: buildInterviewEvaluationUserPrompt(session.Mode, session.JobRole, currentQuestion, answer, ragContext),
		}},
	})
	if err != nil {
		return nil, err
	}

	var result interviewEvaluationResult
	if err := unmarshalLooseJSON(response.Content, &result); err != nil {
		return nil, err
	}
	result.ExpressionScore = clampScore(result.ExpressionScore, 0, 100)
	result.AccuracyScore = clampScore(result.AccuracyScore, 0, 100)
	result.StructureScore = clampScore(result.StructureScore, 0, 100)
	return &result, nil
}

func (s *interviewService) findRoundContext(messages []model.InterviewMessage) (*model.InterviewMessage, *model.InterviewMessage, error) {
	for idx := len(messages) - 1; idx >= 0; idx-- {
		if messages[idx].Role != "candidate" {
			continue
		}
		candidate := messages[idx]
		for q := idx - 1; q >= 0; q-- {
			if messages[q].Role == "interviewer" {
				question := messages[q]
				return &candidate, &question, nil
			}
		}
		return &candidate, nil, errors.New("missing interviewer question context")
	}
	return nil, nil, errors.New("missing candidate answer for streaming round")
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
		ScoreItems:        []InterviewScoreItem{},
		KnowledgeHits:     []KnowledgeHitItem{},
	}
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

func (s *interviewService) buildScoreItemsFromEvaluation(evaluation *interviewEvaluationResult) []InterviewScoreItem {
	return []InterviewScoreItem{
		{Label: "Expression", Score: evaluation.ExpressionScore},
		{Label: "Accuracy", Score: evaluation.AccuracyScore},
		{Label: "Structure", Score: evaluation.StructureScore},
	}
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

func summarizeText(text string, maxLen int) string {
	runes := []rune(strings.TrimSpace(text))
	if len(runes) <= maxLen {
		return string(runes)
	}
	return string(runes[:maxLen]) + "..."
}

func unmarshalLooseJSON(content string, target any) error {
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)
	return json.Unmarshal([]byte(content), target)
}

func clampScore(value, minValue, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func isRetrievalUnavailable(err error) bool {
	return errors.Is(err, ErrEmbeddingNotConfigured) || errors.Is(err, vectorstore.ErrMilvusNotConfigured)
}
