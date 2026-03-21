package repository

import (
	"offerpilot/backend/internal/model"

	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(record *model.ProjectPolishRecord) error
	GetByID(userID, recordID uint64) (*model.ProjectPolishRecord, error)
	ListByUser(userID uint64, limit int) ([]model.ProjectPolishRecord, error)
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(record *model.ProjectPolishRecord) error {
	return r.db.Create(record).Error
}

func (r *projectRepository) GetByID(userID, recordID uint64) (*model.ProjectPolishRecord, error) {
	var record model.ProjectPolishRecord
	if err := r.db.Where("id = ? AND user_id = ?", recordID, userID).First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *projectRepository) ListByUser(userID uint64, limit int) ([]model.ProjectPolishRecord, error) {
	var records []model.ProjectPolishRecord
	query := r.db.Where("user_id = ?", userID).Order("created_at DESC, id DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}
