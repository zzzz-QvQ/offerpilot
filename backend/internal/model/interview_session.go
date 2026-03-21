package model

import "time"

type InterviewSession struct {
	ID          uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID      uint64     `gorm:"column:user_id;not null;index" json:"userId"`
	Mode        string     `gorm:"column:mode;size:50;not null" json:"mode"`
	JobRole     string     `gorm:"column:job_role;size:100;not null" json:"jobRole"`
	Status      string     `gorm:"column:status;size:50;not null;default:'active'" json:"status"`
	TotalRounds int        `gorm:"column:total_rounds;not null;default:0" json:"totalRounds"`
	StartedAt   time.Time  `gorm:"column:started_at;not null" json:"startedAt"`
	EndedAt     *time.Time `gorm:"column:ended_at" json:"endedAt"`
	BaseModel
}

func (InterviewSession) TableName() string {
	return "interview_sessions"
}
