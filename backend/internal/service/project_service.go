package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"offerpilot/backend/internal/model"
	"offerpilot/backend/internal/repository"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

var ErrProjectRecordNotFound = errors.New("project polish record not found")

type CreateProjectPolishRequest struct {
	ProjectName       string `json:"projectName" binding:"required"`
	Background        string `json:"background"`
	ProjectBackground string `json:"projectBackground"`
	TechStack         string `json:"techStack"`
	Responsibilities  string `json:"responsibilities"`
	Responsibility    string `json:"responsibility"`
	Difficulties      string `json:"difficulties"`
	Achievements      string `json:"achievements"`
	Role              string `json:"role"`
}

type ProjectFollowupItem struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type ProjectPolishResponse struct {
	ID                   string                `json:"id"`
	ResumeDescription    string                `json:"resumeDescription"`
	InterviewDescription string                `json:"interviewDescription"`
	Highlights           []string              `json:"highlights"`
	Difficulties         []string              `json:"difficulties"`
	Followups            []ProjectFollowupItem `json:"followups"`
}

type ProjectHistoryItem struct {
	ID          string `json:"id"`
	ProjectName string `json:"projectName"`
	CreatedAt   string `json:"createdAt"`
	Summary     string `json:"summary"`
}

type projectRawInput struct {
	ProjectName      string `json:"projectName"`
	Background       string `json:"background"`
	TechStack        string `json:"techStack"`
	Responsibilities string `json:"responsibilities"`
	Difficulties     string `json:"difficulties"`
	Achievements     string `json:"achievements"`
	Role             string `json:"role"`
}

type projectGeneratedContent struct {
	ResumeVersion    string                `json:"resumeVersion"`
	InterviewVersion string                `json:"interviewVersion"`
	Highlights       []string              `json:"highlights"`
	Difficulties     []string              `json:"difficulties"`
	Followups        []ProjectFollowupItem `json:"followups"`
}

type ProjectService interface {
	Create(userID uint64, req CreateProjectPolishRequest) (*ProjectPolishResponse, error)
	GetHistory(userID uint64) ([]ProjectHistoryItem, error)
	GetDetail(userID, recordID uint64) (*ProjectPolishResponse, error)
}

type projectService struct {
	projectRepo repository.ProjectRepository
	llmService  LLMService
}

func NewProjectService(projectRepo repository.ProjectRepository, llmService LLMService) ProjectService {
	return &projectService{projectRepo: projectRepo, llmService: llmService}
}

func (s *projectService) Create(userID uint64, req CreateProjectPolishRequest) (*ProjectPolishResponse, error) {
	normalized := s.normalizeInput(req)
	generated, err := s.generateContent(normalized)
	if err != nil {
		return nil, err
	}

	rawInputPayload, err := json.Marshal(normalized)
	if err != nil {
		return nil, err
	}
	highlightsPayload, err := json.Marshal(generated.Highlights)
	if err != nil {
		return nil, err
	}
	followupsPayload, err := json.Marshal(generated.Followups)
	if err != nil {
		return nil, err
	}

	record := &model.ProjectPolishRecord{
		UserID:           userID,
		ProjectName:      normalized.ProjectName,
		RawInput:         string(rawInputPayload),
		ResumeVersion:    generated.ResumeDescription,
		InterviewVersion: generated.InterviewDescription,
		Highlights:       string(highlightsPayload),
		Followups:        string(followupsPayload),
	}
	if err := s.projectRepo.Create(record); err != nil {
		return nil, err
	}

	generated.ID = strconv.FormatUint(record.ID, 10)
	return generated, nil
}

func (s *projectService) GetHistory(userID uint64) ([]ProjectHistoryItem, error) {
	records, err := s.projectRepo.ListByUser(userID, 20)
	if err != nil {
		return nil, err
	}

	items := make([]ProjectHistoryItem, 0, len(records))
	for _, record := range records {
		items = append(items, ProjectHistoryItem{
			ID:          strconv.FormatUint(record.ID, 10),
			ProjectName: record.ProjectName,
			CreatedAt:   record.CreatedAt.Format("2006-01-02 15:04"),
			Summary:     summarize(record.ResumeVersion, 56),
		})
	}

	return items, nil
}

func (s *projectService) GetDetail(userID, recordID uint64) (*ProjectPolishResponse, error) {
	record, err := s.projectRepo.GetByID(userID, recordID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectRecordNotFound
		}
		return nil, err
	}

	return s.fromRecord(record), nil
}

func (s *projectService) normalizeInput(req CreateProjectPolishRequest) projectRawInput {
	background := strings.TrimSpace(req.Background)
	if background == "" {
		background = strings.TrimSpace(req.ProjectBackground)
	}

	responsibilities := strings.TrimSpace(req.Responsibilities)
	if responsibilities == "" {
		responsibilities = strings.TrimSpace(req.Responsibility)
	}

	role := strings.TrimSpace(req.Role)
	if role == "" {
		role = "frontend owner"
	}

	achievements := strings.TrimSpace(req.Achievements)
	if achievements == "" {
		achievements = "delivered core features and supported stable launch"
	}

	difficulties := strings.TrimSpace(req.Difficulties)
	if difficulties == "" {
		difficulties = "complex interaction stability, performance tuning and browser compatibility"
	}

	return projectRawInput{
		ProjectName:      strings.TrimSpace(req.ProjectName),
		Background:       defaultString(background, "product-oriented project delivery and iteration"),
		TechStack:        defaultString(strings.TrimSpace(req.TechStack), "Vue 3, TypeScript, Vite"),
		Responsibilities: defaultString(responsibilities, "owned solution design, implementation and release coordination"),
		Difficulties:     difficulties,
		Achievements:     achievements,
		Role:             role,
	}
}

func (s *projectService) generateContent(input projectRawInput) (*ProjectPolishResponse, error) {
	response, err := s.llmService.GenerateText(context.Background(), GenerateTextRequest{
		SystemPrompt: buildProjectPolishSystemPrompt(),
		Messages: []LLMMessage{{
			Role:    "user",
			Content: buildProjectPolishUserPrompt(input),
		}},
	})
	if err != nil {
		return nil, err
	}

	var generated projectGeneratedContent
	if err := unmarshalLooseJSON(response.Content, &generated); err != nil {
		return nil, err
	}

	return &ProjectPolishResponse{
		ResumeDescription:    strings.TrimSpace(generated.ResumeVersion),
		InterviewDescription: strings.TrimSpace(generated.InterviewVersion),
		Highlights:           ensureStringSlice(trimStringSlice(generated.Highlights), []string{"Completed core business modules and shipped them with stable delivery quality."}),
		Difficulties:         ensureStringSlice(trimStringSlice(generated.Difficulties), []string{fmt.Sprintf("Challenge and solution: %s", input.Difficulties)}),
		Followups:            ensureFollowups(trimFollowups(generated.Followups)),
	}, nil
}

func (s *projectService) fromRecord(record *model.ProjectPolishRecord) *ProjectPolishResponse {
	var rawInput projectRawInput
	var highlights []string
	var followups []ProjectFollowupItem

	_ = json.Unmarshal([]byte(record.RawInput), &rawInput)
	_ = json.Unmarshal([]byte(record.Highlights), &highlights)
	_ = json.Unmarshal([]byte(record.Followups), &followups)

	difficultyItems := make([]string, 0)
	for _, item := range splitParagraphs(rawInput.Difficulties) {
		difficultyItems = append(difficultyItems, fmt.Sprintf("Challenge and solution: %s", item))
	}
	if len(difficultyItems) == 0 {
		difficultyItems = []string{"Challenge and solution: split complex requirements into smaller modules and validate quality before release."}
	}

	return &ProjectPolishResponse{
		ID:                   strconv.FormatUint(record.ID, 10),
		ResumeDescription:    record.ResumeVersion,
		InterviewDescription: record.InterviewVersion,
		Highlights:           ensureStringSlice(highlights, []string{"Completed core business modules and shipped them with stable delivery quality."}),
		Difficulties:         difficultyItems,
		Followups:            ensureFollowups(followups),
	}
}

func splitParagraphs(text string) []string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	parts := strings.FieldsFunc(text, func(r rune) bool {
		return r == '\n' || r == ';'
	})
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func ensureStringSlice(items, fallback []string) []string {
	if len(items) == 0 {
		return fallback
	}
	return items
}

func ensureFollowups(items []ProjectFollowupItem) []ProjectFollowupItem {
	if len(items) > 0 {
		return items
	}
	return []ProjectFollowupItem{{
		Question: "What contribution best shows your ownership in this project?",
		Answer:   "I would answer it in four parts: problem background, decision making, implementation process and measurable result.",
	}}
}

func summarize(text string, maxLen int) string {
	clean := strings.TrimSpace(strings.ReplaceAll(text, "\n", " "))
	runes := []rune(clean)
	if len(runes) <= maxLen {
		return clean
	}
	return string(runes[:maxLen]) + "..."
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func trimStringSlice(items []string) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func trimFollowups(items []ProjectFollowupItem) []ProjectFollowupItem {
	result := make([]ProjectFollowupItem, 0, len(items))
	for _, item := range items {
		question := strings.TrimSpace(item.Question)
		answer := strings.TrimSpace(item.Answer)
		if question == "" || answer == "" {
			continue
		}
		result = append(result, ProjectFollowupItem{Question: question, Answer: answer})
	}
	return result
}
