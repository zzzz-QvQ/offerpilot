package model

type Question struct {
	ID              uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title           string `gorm:"column:title;size:255;not null" json:"title"`
	Content         string `gorm:"column:content;type:text;not null" json:"content"`
	Category        string `gorm:"column:category;size:100;not null;index" json:"category"`
	Difficulty      string `gorm:"column:difficulty;size:50;not null" json:"difficulty"`
	Frequency       string `gorm:"column:frequency;size:50;not null;index" json:"frequency"`
	Tags            string `gorm:"column:tags;type:text" json:"tags"`
	StandardAnswer  string `gorm:"column:standard_answer;type:text" json:"standardAnswer"`
	ReferenceSource string `gorm:"column:reference_source;size:255" json:"referenceSource"`
	BaseModel
}

func (Question) TableName() string {
	return "questions"
}
