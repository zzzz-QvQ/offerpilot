package service

import (
	"errors"
	"offerpilot/backend/internal/model"
	"offerpilot/backend/internal/pkg/jwtutil"
	"offerpilot/backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUnauthorized       = errors.New("unauthorized")
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserProfile struct {
	ID       uint64 `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  UserProfile `json:"user"`
}

type AuthService interface {
	Login(input LoginRequest) (*LoginResponse, error)
	GetCurrentUser(userID uint64) (*UserProfile, error)
	Logout(userID uint64) error
}

type authService struct {
	userRepo       repository.UserRepository
	jwtSecret      string
	jwtExpireHours int
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string, jwtExpireHours int) AuthService {
	return &authService{
		userRepo:       userRepo,
		jwtSecret:      jwtSecret,
		jwtExpireHours: jwtExpireHours,
	}
}

func (s *authService) Login(input LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := jwtutil.GenerateToken(s.jwtSecret, s.jwtExpireHours, user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  toUserProfile(user),
	}, nil
}

func (s *authService) GetCurrentUser(userID uint64) (*UserProfile, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUnauthorized
		}
		return nil, err
	}

	profile := toUserProfile(user)
	return &profile, nil
}

func (s *authService) Logout(_ uint64) error {
	return nil
}

func toUserProfile(user *model.User) UserProfile {
	return UserProfile{
		ID:       user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
	}
}
