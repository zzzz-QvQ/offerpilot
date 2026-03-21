package model

type User struct {
	ID           uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Email        string `gorm:"column:email;size:255;not null;uniqueIndex" json:"email"`
	PasswordHash string `gorm:"column:password_hash;size:255;not null" json:"-"`
	Nickname     string `gorm:"column:nickname;size:100;not null" json:"nickname"`
	BaseModel
}

func (User) TableName() string {
	return "users"
}
