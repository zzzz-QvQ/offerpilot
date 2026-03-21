package model

type QuestionFavorite struct {
	ID         uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID     uint64 `gorm:"column:user_id;not null;index:idx_user_question,unique" json:"userId"`
	QuestionID uint64 `gorm:"column:question_id;not null;index:idx_user_question,unique" json:"questionId"`
	BaseModel
}

func (QuestionFavorite) TableName() string {
	return "question_favorites"
}
