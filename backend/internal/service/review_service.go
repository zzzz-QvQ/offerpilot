package service

import (
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

type ReviewService interface {
	GetReport(userID, sessionID uint64) (*ReviewReportResponse, error)
	GetHistory(userID uint64) ([]ReviewHistoryItem, error)
}

type reviewService struct {
	reviewRepo repository.ReviewRepository
}

func NewReviewService(reviewRepo repository.ReviewRepository) ReviewService {
	return &reviewService{reviewRepo: reviewRepo}
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

	return s.buildPlaceholderReport(session, messages), nil
}

func (s *reviewService) GetHistory(userID uint64) ([]ReviewHistoryItem, error) {
	items, err := s.reviewRepo.ListHistory(userID, 20)
	if err != nil {
		return nil, err
	}

	result := make([]ReviewHistoryItem, 0, len(items))
	for _, item := range items {
		finishedAt := "未结束"
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

	response := &ReviewReportResponse{
		OverallScore: report.OverallScore,
		SummaryItems: []ReviewSummaryItem{
			{Label: "总体评分", Value: fmt.Sprintf("%d 分", report.OverallScore), Hint: "基于当前复盘记录生成"},
			{Label: "技术得分", Value: fmt.Sprintf("%d 分", report.TechnicalScore), Hint: "技术准确性与广度"},
			{Label: "表达得分", Value: fmt.Sprintf("%d 分", report.ExpressionScore), Hint: "表达与沟通清晰度"},
		},
		RadarItems: []ReviewRadarItem{
			{Name: "知识准确性", Score: report.TechnicalScore, FullMark: 100},
			{Name: "表达完整度", Score: report.ExpressionScore, FullMark: 100},
			{Name: "结构清晰度", Score: report.LogicScore, FullMark: 100},
			{Name: "追问应对", Score: report.DepthScore, FullMark: 100},
			{Name: "项目结合度", Score: report.ProjectScore, FullMark: 100},
		},
		WeaknessItems:   makeWeaknessItems(weakPoints),
		TrendItems:      makeTrendItems(report.OverallScore),
		SuggestionItems: makeSuggestionItems(suggestions, weakPoints),
	}

	return response
}

func (s *reviewService) buildPlaceholderReport(session *model.InterviewSession, messages []model.InterviewMessage) *ReviewReportResponse {
	candidateCount := 0
	for _, message := range messages {
		if message.Role == "candidate" {
			candidateCount++
		}
	}

	overallScore := 72 + min(candidateCount*3, 18)
	technical := clamp(overallScore+2, 0, 100)
	expression := clamp(overallScore-1, 0, 100)
	logic := clamp(overallScore-3, 0, 100)
	depth := clamp(overallScore-5, 0, 100)
	project := clamp(overallScore+4, 0, 100)

	weakPoints := []string{
		"性能优化回答仍偏抽象，建议补充真实优化指标。",
		"追问时的结构化表达还不够稳定。",
	}
	suggestions := []string{
		"建议增加多轮追问训练，重点强化技术决策解释。",
		"复盘时补充项目结果和量化收益，提高说服力。",
	}

	return &ReviewReportResponse{
		OverallScore: overallScore,
		SummaryItems: []ReviewSummaryItem{
			{Label: "总体评分", Value: fmt.Sprintf("%d 分", overallScore), Hint: "基于当前会话占位分析"},
			{Label: "面试模式", Value: session.Mode, Hint: "当前复盘来源会话模式"},
			{Label: "目标岗位", Value: session.JobRole, Hint: "当前复盘来源岗位方向"},
		},
		RadarItems: []ReviewRadarItem{
			{Name: "知识准确性", Score: technical, FullMark: 100},
			{Name: "表达完整度", Score: expression, FullMark: 100},
			{Name: "结构清晰度", Score: logic, FullMark: 100},
			{Name: "追问应对", Score: depth, FullMark: 100},
			{Name: "项目结合度", Score: project, FullMark: 100},
		},
		WeaknessItems:   makeWeaknessItems(weakPoints),
		TrendItems:      makeTrendItems(overallScore),
		SuggestionItems: makeSuggestionItems(suggestions, weakPoints),
	}
}

func makeWeaknessItems(items []string) []ReviewWeaknessItem {
	result := make([]ReviewWeaknessItem, 0, len(items))
	for idx, item := range items {
		result = append(result, ReviewWeaknessItem{
			Name:  shorten(item, 12),
			Score: 60 + idx*8,
		})
	}
	if len(result) == 0 {
		return []ReviewWeaknessItem{{Name: "暂无薄弱点", Score: 80}}
	}
	return result
}

func makeTrendItems(overall int) []ReviewTrendItem {
	return []ReviewTrendItem{
		{Date: "03-16", Score: clamp(overall-8, 0, 100)},
		{Date: "03-17", Score: clamp(overall-5, 0, 100)},
		{Date: "03-18", Score: clamp(overall-3, 0, 100)},
		{Date: "03-19", Score: clamp(overall-2, 0, 100)},
		{Date: "03-20", Score: overall},
	}
}

func makeSuggestionItems(suggestions, weakPoints []string) []ReviewSuggestionItem {
	result := make([]ReviewSuggestionItem, 0, len(suggestions)+len(weakPoints))
	for _, item := range suggestions {
		result = append(result, ReviewSuggestionItem{Title: shorten(item, 18), Description: item, Type: "建议"})
	}
	for _, item := range weakPoints {
		result = append(result, ReviewSuggestionItem{Title: shorten(item, 18), Description: item, Type: "薄弱点"})
	}
	return result
}

func splitLines(text string) []string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	parts := strings.Split(text, "\n")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(strings.TrimLeft(part, "-•0123456789.、 "))
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

func clamp(value, minValue, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func init() {
	_ = json.Valid
}
