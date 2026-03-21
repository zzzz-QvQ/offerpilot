package model

type ReviewReport struct {
	ID              uint64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	SessionID       uint64 `gorm:"column:session_id;not null;uniqueIndex" json:"sessionId"`
	OverallScore    int    `gorm:"column:overall_score;not null;default:0" json:"overallScore"`
	TechnicalScore  int    `gorm:"column:technical_score;not null;default:0" json:"technicalScore"`
	ExpressionScore int    `gorm:"column:expression_score;not null;default:0" json:"expressionScore"`
	LogicScore      int    `gorm:"column:logic_score;not null;default:0" json:"logicScore"`
	DepthScore      int    `gorm:"column:depth_score;not null;default:0" json:"depthScore"`
	ProjectScore    int    `gorm:"column:project_score;not null;default:0" json:"projectScore"`
	WeakPoints      string `gorm:"column:weak_points;type:text" json:"weakPoints"`
	Suggestions     string `gorm:"column:suggestions;type:text" json:"suggestions"`
	RawReportJSON   string `gorm:"column:raw_report_json;type:text" json:"rawReportJson"`
	BaseModel
}

func (ReviewReport) TableName() string {
	return "review_reports"
}
