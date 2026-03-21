package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"offerpilot/backend/internal/model"
	"offerpilot/backend/internal/repository"
	"strconv"
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

type InterviewService interface {
	CreateSession(userID uint64, req CreateInterviewSessionRequest) (*CreateInterviewSessionResponse, error)
	GetSessionDetail(userID, sessionID uint64) (*InterviewSessionDetail, error)
	SubmitMessage(userID, sessionID uint64, req SubmitInterviewMessageRequest) (*SubmitInterviewMessageResponse, error)
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
		mode = "综合模拟"
	}

	jobRole := req.JobRole
	if jobRole == "" {
		jobRole = "前端工程师"
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

	return s.buildSessionDetail(session, messages), nil
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
		status := "待开始"
		if i < session.TotalRounds {
			status = "已完成"
		} else if session.Status == "finished" {
			status = "已完成"
		} else {
			status = "进行中"
		}

		questions = append(questions, InterviewQuestionItem{
			ID:     fmt.Sprintf("q-%d", i),
			Title:  fmt.Sprintf("第 %d 轮问题", i),
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
	return fmt.Sprintf("当前模式为%s，请先做一个 2 分钟的自我介绍，并说明你在%s岗位上最拿手的能力。", mode, jobRole)
}

func (s *interviewService) buildFollowupQuestion(answer string, round int) string {
	preview := []rune(answer)
	if len(preview) > 24 {
		preview = preview[:24]
	}
	return fmt.Sprintf("你刚才提到“%s”，请继续展开说明你的技术决策、落地过程以及最终结果。", string(preview))
}

func (s *interviewService) buildStatusItems(session *model.InterviewSession) []InterviewStatusItem {
	statusText := "模拟中"
	if session.Status == "finished" {
		statusText = "已结束"
	}
	return []InterviewStatusItem{
		{Label: "会话状态", Value: statusText},
		{Label: "当前轮次", Value: fmt.Sprintf("第 %d 轮", max(session.TotalRounds, 1))},
		{Label: "模式", Value: session.Mode},
		{Label: "岗位", Value: session.JobRole},
	}
}

func (s *interviewService) buildScoreItems() []InterviewScoreItem {
	return []InterviewScoreItem{
		{Label: "表达完整度", Score: 82},
		{Label: "知识准确性", Score: 80},
		{Label: "结构清晰度", Score: 79},
	}
}

func (s *interviewService) buildKnowledgeHits() []KnowledgeHitItem {
	return []KnowledgeHitItem{
		{Name: "项目表达", Level: "高", Summary: "能够围绕项目经历组织答案。"},
		{Name: "技术决策", Level: "中", Summary: "提到了方案，但还可以更量化。"},
		{Name: "结果复盘", Level: "低", Summary: "建议补充指标和结果验证。"},
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
	return InterviewMessageItem{
		ID:      strconv.FormatUint(message.ID, 10),
		Role:    message.Role,
		Content: message.Content,
		Time:    time.UnixMilli(message.CreatedAt).Format("15:04"),
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
