package repository

import (
	"offerpilot/backend/internal/model"
	"time"

	"gorm.io/gorm"
)

type InterviewRepository interface {
	CreateSession(session *model.InterviewSession) error
	GetSessionByID(sessionID, userID uint64) (*model.InterviewSession, error)
	UpdateSession(session *model.InterviewSession) error
	ListMessages(sessionID uint64) ([]model.InterviewMessage, error)
	CreateMessage(message *model.InterviewMessage) error
	GetLatestRoundNo(sessionID uint64) (int, error)
}

type interviewRepository struct {
	db *gorm.DB
}

func NewInterviewRepository(db *gorm.DB) InterviewRepository {
	return &interviewRepository{db: db}
}

func (r *interviewRepository) CreateSession(session *model.InterviewSession) error {
	if session.StartedAt.IsZero() {
		session.StartedAt = time.Now()
	}
	return r.db.Create(session).Error
}

func (r *interviewRepository) GetSessionByID(sessionID, userID uint64) (*model.InterviewSession, error) {
	var session model.InterviewSession
	if err := r.db.Where("id = ? AND user_id = ?", sessionID, userID).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *interviewRepository) UpdateSession(session *model.InterviewSession) error {
	return r.db.Save(session).Error
}

func (r *interviewRepository) ListMessages(sessionID uint64) ([]model.InterviewMessage, error) {
	var messages []model.InterviewMessage
	if err := r.db.Where("session_id = ?", sessionID).Order("created_at ASC, id ASC").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *interviewRepository) CreateMessage(message *model.InterviewMessage) error {
	return r.db.Create(message).Error
}

func (r *interviewRepository) GetLatestRoundNo(sessionID uint64) (int, error) {
	var latest model.InterviewMessage
	if err := r.db.Where("session_id = ?", sessionID).Order("round_no DESC, id DESC").First(&latest).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}
	return latest.RoundNo, nil
}
