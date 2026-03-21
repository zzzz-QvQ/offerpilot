package database

import (
	"offerpilot/backend/internal/config"
	"offerpilot/backend/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQL(cfg config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.TrainingSession{},
		&model.WeakPoint{},
		&model.Question{},
		&model.QuestionFavorite{},
	); err != nil {
		return nil, err
	}

	return db, nil
}
