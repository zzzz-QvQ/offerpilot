package repository

import (
	"time"

	"offerpilot/backend/internal/model"

	"gorm.io/gorm"
)

type DashboardSummaryStats struct {
	TotalTrainingCount  int64
	AverageScore        float64
	WeeklyTrainingCount int64
	WeakestCategory     string
}

type DashboardRecentSession struct {
	ID          uint64
	Title       string
	Score       float64
	DurationMin int
	FinishedAt  time.Time
	Tag         string
}

type DashboardWeakPoint struct {
	Category string
	Score    float64
}

type DashboardRepository interface {
	GetSummaryStats(userID uint64) (DashboardSummaryStats, error)
	GetRecentSessions(userID uint64, limit int) ([]DashboardRecentSession, error)
	GetWeakPoints(userID uint64, limit int) ([]DashboardWeakPoint, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetSummaryStats(userID uint64) (DashboardSummaryStats, error) {
	var summary DashboardSummaryStats

	if err := r.db.Model(&model.TrainingSession{}).
		Where("user_id = ?", userID).
		Count(&summary.TotalTrainingCount).Error; err != nil {
		return summary, err
	}

	if err := r.db.Model(&model.TrainingSession{}).
		Where("user_id = ?", userID).
		Select("COALESCE(AVG(score), 0)").
		Scan(&summary.AverageScore).Error; err != nil {
		return summary, err
	}

	weekStart := time.Now().AddDate(0, 0, -7)
	if err := r.db.Model(&model.TrainingSession{}).
		Where("user_id = ? AND created_at >= ?", userID, weekStart).
		Count(&summary.WeeklyTrainingCount).Error; err != nil {
		return summary, err
	}

	var weakest model.WeakPoint
	err := r.db.Where("user_id = ?", userID).
		Order("score ASC, hit_count DESC, updated_at DESC").
		First(&weakest).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return summary, nil
		}
		return summary, err
	}

	summary.WeakestCategory = weakest.Category
	return summary, nil
}

func (r *dashboardRepository) GetRecentSessions(userID uint64, limit int) ([]DashboardRecentSession, error) {
	var sessions []DashboardRecentSession
	err := r.db.Model(&model.TrainingSession{}).
		Where("user_id = ?", userID).
		Select("id, title, score, duration_min, updated_at AS finished_at, category AS tag").
		Order("updated_at DESC").
		Limit(limit).
		Scan(&sessions).Error
	return sessions, err
}

func (r *dashboardRepository) GetWeakPoints(userID uint64, limit int) ([]DashboardWeakPoint, error) {
	var weakPoints []DashboardWeakPoint
	err := r.db.Model(&model.WeakPoint{}).
		Where("user_id = ?", userID).
		Select("category, score").
		Order("score ASC, hit_count DESC, updated_at DESC").
		Limit(limit).
		Scan(&weakPoints).Error
	return weakPoints, err
}
