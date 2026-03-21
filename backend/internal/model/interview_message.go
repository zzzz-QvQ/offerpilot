package model

type InterviewMessage struct {
	ID            uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	SessionID     uint64 `gorm:"column:session_id;not null;index" json:"sessionId"`
	Role          string `gorm:"column:role;size:50;not null" json:"role"`
	RoundNo       int    `gorm:"column:round_no;not null;default:1" json:"roundNo"`
	QuestionID    string `gorm:"column:question_id;size:100;not null;default:''" json:"questionId"`
	Content       string `gorm:"column:content;type:text;not null" json:"content"`
	ScoreSnapshot string `gorm:"column:score_snapshot;type:text" json:"scoreSnapshot"`
	CreatedAt     int64  `gorm:"column:created_at;autoCreateTime:milli" json:"createdAt"`
}

func (InterviewMessage) TableName() string {
	return "interview_messages"
}
