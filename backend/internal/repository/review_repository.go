package repository

import (
	"offerpilot/backend/internal/model"

	"gorm.io/gorm"
)

type ReviewHistoryItem struct {
	SessionID uint64
	Title     string
	Score     int
	EndedAt   *int64
}

type ReviewRepository interface {
	GetBySessionID(sessionID uint64) (*model.ReviewReport, error)
	ListHistory(userID uint64, limit int) ([]ReviewHistoryItem, error)
	GetSessionWithMessages(userID, sessionID uint64) (*model.InterviewSession, []model.InterviewMessage, error)
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) GetBySessionID(sessionID uint64) (*model.ReviewReport, error) {
	var report model.ReviewReport
	if err := r.db.Where("session_id = ?", sessionID).First(&report).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *reviewRepository) ListHistory(userID uint64, limit int) ([]ReviewHistoryItem, error) {
	var items []ReviewHistoryItem
	err := r.db.Table("interview_sessions").
		Select(`
			interview_sessions.id AS session_id,
			CONCAT(interview_sessions.mode, ' - ', interview_sessions.job_role) AS title,
			COALESCE(review_reports.overall_score, 0) AS score,
			UNIX_TIMESTAMP(interview_sessions.ended_at) * 1000 AS ended_at
		`).
		Joins("LEFT JOIN review_reports ON review_reports.session_id = interview_sessions.id").
		Where("interview_sessions.user_id = ?", userID).
		Order("interview_sessions.updated_at DESC").
		Limit(limit).
		Scan(&items).Error
	return items, err
}

func (r *reviewRepository) GetSessionWithMessages(userID, sessionID uint64) (*model.InterviewSession, []model.InterviewMessage, error) {
	var session model.InterviewSession
	if err := r.db.Where("id = ? AND user_id = ?", sessionID, userID).First(&session).Error; err != nil {
		return nil, nil, err
	}

	var messages []model.InterviewMessage
	if err := r.db.Where("session_id = ?", sessionID).Order("created_at ASC, id ASC").Find(&messages).Error; err != nil {
		return nil, nil, err
	}

	return &session, messages, nil
}
