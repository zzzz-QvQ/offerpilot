package model

type ProjectPolishRecord struct {
	ID               uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID           uint64 `gorm:"column:user_id;not null;index" json:"userId"`
	ProjectName      string `gorm:"column:project_name;size:255;not null;index" json:"projectName"`
	RawInput         string `gorm:"column:raw_input;type:longtext;not null" json:"rawInput"`
	ResumeVersion    string `gorm:"column:resume_version;type:longtext;not null" json:"resumeVersion"`
	InterviewVersion string `gorm:"column:interview_version;type:longtext;not null" json:"interviewVersion"`
	Highlights       string `gorm:"column:highlights;type:longtext;not null" json:"highlights"`
	Followups        string `gorm:"column:followups;type:longtext;not null" json:"followups"`
	BaseModel
}

func (ProjectPolishRecord) TableName() string {
	return "project_polish_records"
}
