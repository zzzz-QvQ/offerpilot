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

var ErrReviewNotFound = errors.New("review not found")

type ReviewSummaryItem struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Hint  string `json:"hint"`
}

type ReviewRadarItem struct {
	Name     string `json:"name"`
	Score    int    `json:"score"`
	FullMark int    `json:"fullMark"`
}

type ReviewWeaknessItem struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type ReviewTrendItem struct {
	Date  string `json:"date"`
	Score int    `json:"score"`
}

type ReviewSuggestionItem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type ReviewHistoryItem struct {
	SessionID  string `json:"sessionId"`
	Title      string `json:"title"`
	FinishedAt string `json:"finishedAt"`
	Score      int    `json:"score"`
}

type ReviewReportResponse struct {
	OverallScore    int                    `json:"overallScore"`
	SummaryItems    []ReviewSummaryItem    `json:"summaryItems"`
	RadarItems      []ReviewRadarItem      `json:"radarItems"`
	WeaknessItems   []ReviewWeaknessItem   `json:"weaknessItems"`
	TrendItems      []ReviewTrendItem      `json:"trendItems"`
	SuggestionItems []ReviewSuggestionItem `json:"suggestionItems"`
}

type reviewLLMResult struct {
	OverallSummary string   `json:"overallSummary"`
	OverallScore   int      `json:"overallScore"`
	Technical      int      `json:"technical"`
	Expression     int      `json:"expression"`
	Logic          int      `json:"logic"`
	Depth          int      `json:"depth"`
	Project        int      `json:"project"`
	WeakPoints     []string `json:"weakPoints"`
	Suggestions    []string `json:"suggestions"`
}

type ReviewService interface {
	GetReport(userID, sessionID uint64) (*ReviewReportResponse, error)
	GetHistory(userID uint64) ([]ReviewHistoryItem, error)
}

type reviewService struct {
	reviewRepo        repository.ReviewRepository
	llmService        LLMService
	questionRetrieval QuestionRetrievalService
}

func NewReviewService(reviewRepo repository.ReviewRepository, llmService LLMService, questionRetrieval QuestionRetrievalService) ReviewService {
	return &reviewService{
		reviewRepo:        reviewRepo,
		llmService:        llmService,
		questionRetrieval: questionRetrieval,
	}
}

func (s *reviewService) GetReport(userID, sessionID uint64) (*ReviewReportResponse, error) {
	report, err := s.reviewRepo.GetBySessionID(sessionID)
	if err == nil {
		return s.fromPersistedReport(report), nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	session, messages, err := s.reviewRepo.GetSessionWithMessages(userID, sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReviewNotFound
		}
		return nil, err
	}

	return s.buildLLMReport(session, messages)
}

func (s *reviewService) GetHistory(userID uint64) ([]ReviewHistoryItem, error) {
	items, err := s.reviewRepo.ListHistory(userID, 20)
	if err != nil {
		return nil, err
	}

	result := make([]ReviewHistoryItem, 0, len(items))
	for _, item := range items {
		finishedAt := "Not finished"
		if item.EndedAt != nil && *item.EndedAt > 0 {
			finishedAt = time.UnixMilli(*item.EndedAt).Format("2006-01-02 15:04")
		}
		result = append(result, ReviewHistoryItem{
			SessionID:  strconv.FormatUint(item.SessionID, 10),
			Title:      item.Title,
			FinishedAt: finishedAt,
			Score:      item.Score,
		})
	}

	return result, nil
}

func (s *reviewService) fromPersistedReport(report *model.ReviewReport) *ReviewReportResponse {
	weakPoints := splitLines(report.WeakPoints)
	suggestions := splitLines(report.Suggestions)

	return &ReviewReportResponse{
		OverallScore: report.OverallScore,
		SummaryItems: []ReviewSummaryItem{
			{Label: "Overall Score", Value: fmt.Sprintf("%d", report.OverallScore), Hint: "Generated from stored review report"},
			{Label: "Technical", Value: fmt.Sprintf("%d", report.TechnicalScore), Hint: "Technical accuracy and depth"},
			{Label: "Expression", Value: fmt.Sprintf("%d", report.ExpressionScore), Hint: "Clarity and communication"},
		},
		RadarItems: []ReviewRadarItem{
			{Name: "technical", Score: report.TechnicalScore, FullMark: 100},
			{Name: "expression", Score: report.ExpressionScore, FullMark: 100},
			{Name: "logic", Score: report.LogicScore, FullMark: 100},
			{Name: "depth", Score: report.DepthScore, FullMark: 100},
			{Name: "project", Score: report.ProjectScore, FullMark: 100},
		},
		WeaknessItems:   makeWeaknessItems(weakPoints),
		TrendItems:      makeTrendItems(report.OverallScore),
		SuggestionItems: makeSuggestionItems(suggestions, weakPoints),
	}
}

func (s *reviewService) buildLLMReport(session *model.InterviewSession, messages []model.InterviewMessage) (*ReviewReportResponse, error) {
	referenceContext, _ := s.buildReferenceContext(context.Background(), session, messages)
	conversation := s.formatMessages(messages)

	response, err := s.llmService.GenerateText(context.Background(), GenerateTextRequest{
		SystemPrompt: buildReviewAnalysisSystemPrompt(),
		Messages: []LLMMessage{{
			Role:    "user",
			Content: buildReviewAnalysisUserPrompt(session.Mode, session.JobRole, conversation, referenceContext),
		}},
	})
	if err != nil {
		return nil, err
	}

	var generated reviewLLMResult
	if err := unmarshalLooseJSON(response.Content, &generated); err != nil {
		return nil, err
	}

	technical := reviewClampScore(generated.Technical)
	expression := reviewClampScore(generated.Expression)
	logic := reviewClampScore(generated.Logic)
	depth := reviewClampScore(generated.Depth)
	project := reviewClampScore(generated.Project)
	overall := reviewClampScore(generated.OverallScore)
	if overall == 0 {
		overall = reviewClampScore((technical + expression + logic + depth + project) / 5)
	}

	weakPoints := trimLines(generated.WeakPoints)
	suggestions := trimLines(generated.Suggestions)

	return &ReviewReportResponse{
		OverallScore: overall,
		SummaryItems: []ReviewSummaryItem{
			{Label: "Overall Score", Value: fmt.Sprintf("%d", overall), Hint: strings.TrimSpace(generated.OverallSummary)},
			{Label: "Interview Mode", Value: session.Mode, Hint: "Session metadata"},
			{Label: "Target Role", Value: session.JobRole, Hint: "Session metadata"},
		},
		RadarItems: []ReviewRadarItem{
			{Name: "technical", Score: technical, FullMark: 100},
			{Name: "expression", Score: expression, FullMark: 100},
			{Name: "logic", Score: logic, FullMark: 100},
			{Name: "depth", Score: depth, FullMark: 100},
			{Name: "project", Score: project, FullMark: 100},
		},
		WeaknessItems:   makeWeaknessItems(weakPoints),
		TrendItems:      makeTrendItems(overall),
		SuggestionItems: makeSuggestionItems(suggestions, weakPoints),
	}, nil
}

func (s *reviewService) buildReferenceContext(ctx context.Context, session *model.InterviewSession, messages []model.InterviewMessage) (string, error) {
	if s.questionRetrieval == nil {
		return "", nil
	}

	queryParts := []string{session.JobRole, session.Mode}
	for _, message := range messages {
		if message.Role == "candidate" || message.Role == "interviewer" {
			queryParts = append(queryParts, message.Content)
		}
	}

	hits, err := s.questionRetrieval.SearchQuestions(ctx, strings.Join(queryParts, "\n"), 3)
	if err != nil {
		if isRetrievalUnavailable(err) {
			return "", nil
		}
		return "", err
	}
	if len(hits) == 0 {
		return "", nil
	}

	parts := make([]string, 0, len(hits))
	for idx, hit := range hits {
		parts = append(parts, fmt.Sprintf("[%d] %s | %s\n%s", idx+1, hit.Title, hit.Category, hit.Text))
	}
	return strings.Join(parts, "\n\n"), nil
}

func (s *reviewService) formatMessages(messages []model.InterviewMessage) string {
	lines := make([]string, 0, len(messages))
	for _, message := range messages {
		role := strings.Title(message.Role)
		lines = append(lines, fmt.Sprintf("%s: %s", role, strings.TrimSpace(message.Content)))
	}
	return strings.Join(lines, "\n")
}

func makeWeaknessItems(items []string) []ReviewWeaknessItem {
	result := make([]ReviewWeaknessItem, 0, len(items))
	for idx, item := range items {
		result = append(result, ReviewWeaknessItem{
			Name:  shorten(item, 24),
			Score: 60 + idx*8,
		})
	}
	if len(result) == 0 {
		return []ReviewWeaknessItem{{Name: "No obvious weak point", Score: 80}}
	}
	return result
}

func makeTrendItems(overall int) []ReviewTrendItem {
	return []ReviewTrendItem{
		{Date: "03-16", Score: reviewClampScore(overall - 8)},
		{Date: "03-17", Score: reviewClampScore(overall - 5)},
		{Date: "03-18", Score: reviewClampScore(overall - 3)},
		{Date: "03-19", Score: reviewClampScore(overall - 2)},
		{Date: "03-20", Score: overall},
	}
}

func makeSuggestionItems(suggestions, weakPoints []string) []ReviewSuggestionItem {
	result := make([]ReviewSuggestionItem, 0, len(suggestions)+len(weakPoints))
	for _, item := range suggestions {
		result = append(result, ReviewSuggestionItem{Title: shorten(item, 24), Description: item, Type: "suggestion"})
	}
	for _, item := range weakPoints {
		result = append(result, ReviewSuggestionItem{Title: shorten(item, 24), Description: item, Type: "weakness"})
	}
	return result
}

func splitLines(text string) []string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	parts := strings.Split(text, "\n")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(strings.TrimLeft(part, "-• 0123456789.、"))
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func trimLines(items []string) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func shorten(text string, maxLen int) string {
	runes := []rune(text)
	if len(runes) <= maxLen {
		return text
	}
	return string(runes[:maxLen]) + "..."
}

func reviewClampScore(value int) int {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value
}

func init() {
	_ = json.Valid
}
