package service

import (
	"fmt"
	"time"

	"offerpilot/backend/internal/repository"
)

type DashboardStatItem struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Hint  string `json:"hint"`
}

type RecommendTaskItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Level       string `json:"level"`
}

type DashboardSummaryResponse struct {
	StatCards      []DashboardStatItem `json:"statCards"`
	RecommendTasks []RecommendTaskItem `json:"recommendTasks"`
}

type RecentSessionItem struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Score        int    `json:"score"`
	DurationText string `json:"durationText"`
	FinishedAt   string `json:"finishedAt"`
	Tag          string `json:"tag"`
}

type WeakPointItem struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type DashboardService interface {
	GetSummary(userID uint64) (*DashboardSummaryResponse, error)
	GetRecentSessions(userID uint64) ([]RecentSessionItem, error)
	GetWeakPoints(userID uint64) ([]WeakPointItem, error)
}

type dashboardService struct {
	dashboardRepo repository.DashboardRepository
}

func NewDashboardService(dashboardRepo repository.DashboardRepository) DashboardService {
	return &dashboardService{dashboardRepo: dashboardRepo}
}

func (s *dashboardService) GetSummary(userID uint64) (*DashboardSummaryResponse, error) {
	summary, err := s.dashboardRepo.GetSummaryStats(userID)
	if err != nil {
		return nil, err
	}

	weakestCategory := summary.WeakestCategory
	if weakestCategory == "" {
		weakestCategory = "暂无数据"
	}

	response := &DashboardSummaryResponse{
		StatCards: []DashboardStatItem{
			{Label: "总训练次数", Value: fmt.Sprintf("%d", summary.TotalTrainingCount), Hint: "累计完成的训练场次"},
			{Label: "平均分", Value: fmt.Sprintf("%d 分", int(summary.AverageScore+0.5)), Hint: "最近阶段综合表现"},
			{Label: "本周训练次数", Value: fmt.Sprintf("%d", summary.WeeklyTrainingCount), Hint: "近 7 天训练活跃度"},
			{Label: "当前最大薄弱项", Value: weakestCategory, Hint: "建议优先安排专项强化"},
		},
		RecommendTasks: buildRecommendTasks(weakestCategory),
	}

	return response, nil
}

func (s *dashboardService) GetRecentSessions(userID uint64) ([]RecentSessionItem, error) {
	sessions, err := s.dashboardRepo.GetRecentSessions(userID, 6)
	if err != nil {
		return nil, err
	}

	if len(sessions) == 0 {
		return []RecentSessionItem{}, nil
	}

	items := make([]RecentSessionItem, 0, len(sessions))
	for _, session := range sessions {
		items = append(items, RecentSessionItem{
			ID:           fmt.Sprintf("%d", session.ID),
			Title:        session.Title,
			Score:        int(session.Score + 0.5),
			DurationText: fmt.Sprintf("%d 分钟", session.DurationMin),
			FinishedAt:   formatRecentTime(session.FinishedAt),
			Tag:          session.Tag,
		})
	}

	return items, nil
}

func (s *dashboardService) GetWeakPoints(userID uint64) ([]WeakPointItem, error) {
	weakPoints, err := s.dashboardRepo.GetWeakPoints(userID, 6)
	if err != nil {
		return nil, err
	}

	if len(weakPoints) == 0 {
		return []WeakPointItem{
			{Name: "性能优化", Value: 60},
			{Name: "浏览器原理", Value: 68},
			{Name: "工程化", Value: 72},
		}, nil
	}

	items := make([]WeakPointItem, 0, len(weakPoints))
	for _, item := range weakPoints {
		items = append(items, WeakPointItem{
			Name:  item.Category,
			Value: int(item.Score + 0.5),
		})
	}

	return items, nil
}

func buildRecommendTasks(weakestCategory string) []RecommendTaskItem {
	if weakestCategory == "" || weakestCategory == "暂无数据" {
		return []RecommendTaskItem{
			{
				ID:          "task-default-1",
				Title:       "开始首次专项训练",
				Description: "当前训练数据较少，建议先完成 1-2 轮专项训练，建立基础画像。",
				Level:       "中优先级",
			},
		}
	}

	return []RecommendTaskItem{
		{
			ID:          "task-1",
			Title:       weakestCategory + " 专项训练",
			Description: "针对当前最薄弱模块安排集中训练，优先补齐核心概念和标准回答结构。",
			Level:       "高优先级",
		},
		{
			ID:          "task-2",
			Title:       "最近训练复盘整理",
			Description: "回看近几次训练记录，整理扣分点并提炼统一答题框架。",
			Level:       "中优先级",
		},
	}
}

func formatRecentTime(t time.Time) string {
	if t.IsZero() {
		return "未知时间"
	}

	now := time.Now()
	if sameDay(now, t) {
		return "今天 " + t.Format("15:04")
	}

	yesterday := now.AddDate(0, 0, -1)
	if sameDay(yesterday, t) {
		return "昨天 " + t.Format("15:04")
	}

	return t.Format("01-02 15:04")
}

func sameDay(a, b time.Time) bool {
	aYear, aMonth, aDay := a.Date()
	bYear, bMonth, bDay := b.Date()
	return aYear == bYear && aMonth == bMonth && aDay == bDay
}
