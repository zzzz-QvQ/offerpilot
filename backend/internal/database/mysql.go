package database

import (
	"offerpilot/backend/internal/config"
	"offerpilot/backend/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
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
		&model.InterviewSession{},
		&model.InterviewMessage{},
		&model.ReviewReport{},
		&model.ProjectPolishRecord{},
	); err != nil {
		return nil, err
	}

	if err := ensureDemoUser(db); err != nil {
		return nil, err
	}

	return db, nil
}

func ensureDemoUser(db *gorm.DB) error {
	const (
		email    = "test@example.com"
		password = "password"
		nickname = "Test User"
	)

	var existing model.User
	err := db.Where("email = ?", email).First(&existing).Error
	if err == nil {
		return nil
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return db.Create(&model.User{
		Email:        email,
		PasswordHash: string(passwordHash),
		Nickname:     nickname,
	}).Error
}
