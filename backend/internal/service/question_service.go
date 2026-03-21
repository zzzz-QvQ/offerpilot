package service

import (
	"errors"
	"strconv"
	"strings"

	"offerpilot/backend/internal/repository"

	"gorm.io/gorm"
)

var ErrQuestionNotFound = errors.New("question not found")

type QuestionListParams struct {
	Keyword   string
	Category  string
	Frequency string
	Page      int
	PageSize  int
}

type QuestionItem struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Category     string   `json:"category"`
	Frequency    string   `json:"frequency"`
	IsFavorite   bool     `json:"isFavorite"`
	Summary      string   `json:"summary"`
	Difficulty   string   `json:"difficulty"`
	AnswerPoints []string `json:"answerPoints"`
}

type ToggleFavoriteResponse struct {
	QuestionID string `json:"questionId"`
	IsFavorite bool   `json:"isFavorite"`
}

type QuestionListResponse struct {
	List     []QuestionItem `json:"list"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
}

type QuestionService interface {
	List(userID uint64, params QuestionListParams) (*QuestionListResponse, error)
	GetDetail(userID uint64, questionID uint64) (*QuestionItem, error)
	AddFavorite(userID uint64, questionID uint64) (*ToggleFavoriteResponse, error)
	RemoveFavorite(userID uint64, questionID uint64) (*ToggleFavoriteResponse, error)
}

type questionService struct {
	questionRepo repository.QuestionRepository
}

func NewQuestionService(questionRepo repository.QuestionRepository) QuestionService {
	return &questionService{questionRepo: questionRepo}
}

func (s *questionService) List(userID uint64, params QuestionListParams) (*QuestionListResponse, error) {
	page := params.Page
	if page <= 0 {
		page = 1
	}

	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	items, total, err := s.questionRepo.List(repository.QuestionQuery{
		Keyword:   params.Keyword,
		Category:  params.Category,
		Frequency: params.Frequency,
		Page:      page,
		PageSize:  pageSize,
		UserID:    userID,
	})
	if err != nil {
		return nil, err
	}

	result := make([]QuestionItem, 0, len(items))
	for _, item := range items {
		result = append(result, toQuestionItem(item))
	}

	return &QuestionListResponse{
		List:     result,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *questionService) GetDetail(userID uint64, questionID uint64) (*QuestionItem, error) {
	item, err := s.questionRepo.GetByID(userID, questionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrQuestionNotFound
		}
		return nil, err
	}

	result := toQuestionItem(*item)
	return &result, nil
}

func (s *questionService) AddFavorite(userID uint64, questionID uint64) (*ToggleFavoriteResponse, error) {
	if _, err := s.GetDetail(userID, questionID); err != nil {
		return nil, err
	}

	if err := s.questionRepo.AddFavorite(userID, questionID); err != nil {
		return nil, err
	}

	return &ToggleFavoriteResponse{QuestionID: strconv.FormatUint(questionID, 10), IsFavorite: true}, nil
}

func (s *questionService) RemoveFavorite(userID uint64, questionID uint64) (*ToggleFavoriteResponse, error) {
	if _, err := s.GetDetail(userID, questionID); err != nil {
		return nil, err
	}

	if err := s.questionRepo.RemoveFavorite(userID, questionID); err != nil {
		return nil, err
	}

	return &ToggleFavoriteResponse{QuestionID: strconv.FormatUint(questionID, 10), IsFavorite: false}, nil
}

func toQuestionItem(item repository.QuestionListItem) QuestionItem {
	answerSource := item.StandardAnswer
	if strings.TrimSpace(answerSource) == "" {
		answerSource = item.Content
	}

	return QuestionItem{
		ID:           strconv.FormatUint(item.ID, 10),
		Title:        item.Title,
		Category:     item.Category,
		Frequency:    item.Frequency,
		IsFavorite:   item.IsFavorite,
		Summary:      buildSummary(item.Content),
		Difficulty:   item.Difficulty,
		AnswerPoints: splitAnswerPoints(answerSource),
	}
}

func buildSummary(content string) string {
	content = strings.TrimSpace(content)
	if content == "" {
		return "暂无题目摘要"
	}

	runes := []rune(content)
	if len(runes) <= 80 {
		return content
	}

	return string(runes[:80]) + "..."
}

func splitAnswerPoints(text string) []string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(text, "\n")
	points := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(strings.TrimLeft(line, "-•0123456789.、 "))
		if trimmed != "" {
			points = append(points, trimmed)
		}
	}

	if len(points) > 0 {
		return points
	}

	parts := strings.Split(text, "；")
	if len(parts) <= 1 {
		parts = strings.Split(text, ";")
	}

	fallback := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			fallback = append(fallback, trimmed)
		}
	}

	if len(fallback) > 0 {
		return fallback
	}

	return []string{"暂无标准答案要点"}
}
