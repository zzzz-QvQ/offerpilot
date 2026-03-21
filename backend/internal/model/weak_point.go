package model

type WeakPoint struct {
	ID       uint64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID   uint64  `gorm:"column:user_id;not null;index" json:"userId"`
	Category string  `gorm:"column:category;size:100;not null" json:"category"`
	Score    float64 `gorm:"column:score;type:decimal(5,2);not null;default:0" json:"score"`
	HitCount int     `gorm:"column:hit_count;not null;default:0" json:"hitCount"`
	BaseModel
}

func (WeakPoint) TableName() string {
	return "weak_points"
}
