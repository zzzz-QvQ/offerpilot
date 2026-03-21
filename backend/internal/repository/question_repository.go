package repository

import (
	"strings"

	"offerpilot/backend/internal/model"

	"gorm.io/gorm"
)

type QuestionQuery struct {
	Keyword   string
	Category  string
	Frequency string
	Page      int
	PageSize  int
	UserID    uint64
}

type QuestionListItem struct {
	ID             uint64
	Title          string
	Category       string
	Frequency      string
	Difficulty     string
	Content        string
	StandardAnswer string
	IsFavorite     bool
}

type QuestionRepository interface {
	List(query QuestionQuery) ([]QuestionListItem, int64, error)
	GetByID(userID uint64, questionID uint64) (*QuestionListItem, error)
	AddFavorite(userID uint64, questionID uint64) error
	RemoveFavorite(userID uint64, questionID uint64) error
}

type questionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) List(query QuestionQuery) ([]QuestionListItem, int64, error) {
	var total int64
	baseQuery := r.db.Model(&model.Question{}).
		Select(`
			questions.id,
			questions.title,
			questions.category,
			questions.frequency,
			questions.difficulty,
			questions.content,
			questions.standard_answer,
			CASE WHEN question_favorites.id IS NULL THEN false ELSE true END AS is_favorite
		`).
		Joins("LEFT JOIN question_favorites ON question_favorites.question_id = questions.id AND question_favorites.user_id = ?", query.UserID)

	baseQuery = applyQuestionFilters(baseQuery, query)

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []QuestionListItem
	offset := (query.Page - 1) * query.PageSize
	if err := baseQuery.
		Order("questions.updated_at DESC").
		Limit(query.PageSize).
		Offset(offset).
		Scan(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *questionRepository) GetByID(userID uint64, questionID uint64) (*QuestionListItem, error) {
	var item QuestionListItem
	err := r.db.Model(&model.Question{}).
		Select(`
			questions.id,
			questions.title,
			questions.category,
			questions.frequency,
			questions.difficulty,
			questions.content,
			questions.standard_answer,
			CASE WHEN question_favorites.id IS NULL THEN false ELSE true END AS is_favorite
		`).
		Joins("LEFT JOIN question_favorites ON question_favorites.question_id = questions.id AND question_favorites.user_id = ?", userID).
		Where("questions.id = ?", questionID).
		First(&item).Error
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *questionRepository) AddFavorite(userID uint64, questionID uint64) error {
	favorite := model.QuestionFavorite{
		UserID:     userID,
		QuestionID: questionID,
	}

	return r.db.Where(model.QuestionFavorite{UserID: userID, QuestionID: questionID}).FirstOrCreate(&favorite).Error
}

func (r *questionRepository) RemoveFavorite(userID uint64, questionID uint64) error {
	return r.db.Where("user_id = ? AND question_id = ?", userID, questionID).Delete(&model.QuestionFavorite{}).Error
}

func applyQuestionFilters(db *gorm.DB, query QuestionQuery) *gorm.DB {
	if keyword := strings.TrimSpace(query.Keyword); keyword != "" {
		likeKeyword := "%" + keyword + "%"
		db = db.Where("questions.title LIKE ? OR questions.content LIKE ?", likeKeyword, likeKeyword)
	}

	if category := strings.TrimSpace(query.Category); category != "" {
		db = db.Where("questions.category = ?", category)
	}

	if frequency := strings.TrimSpace(query.Frequency); frequency != "" {
		db = db.Where("questions.frequency = ?", frequency)
	}

	return db
}
