package model

type TrainingSession struct {
	ID          uint64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID      uint64  `gorm:"column:user_id;not null;index" json:"userId"`
	Title       string  `gorm:"column:title;size:255;not null" json:"title"`
	Category    string  `gorm:"column:category;size:100;not null" json:"category"`
	Score       float64 `gorm:"column:score;type:decimal(5,2);not null;default:0" json:"score"`
	DurationMin int     `gorm:"column:duration_min;not null;default:0" json:"durationMin"`
	SessionType string  `gorm:"column:session_type;size:50;not null;default:''" json:"sessionType"`
	BaseModel
}

func (TrainingSession) TableName() string {
	return "training_sessions"
}
