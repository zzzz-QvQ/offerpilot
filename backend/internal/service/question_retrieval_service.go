package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"offerpilot/backend/internal/pkg/vectorstore"
	"offerpilot/backend/internal/repository"
)

type QuestionVectorDocument struct {
	QuestionID string    `json:"questionId"`
	Title      string    `json:"title"`
	Category   string    `json:"category"`
	Text       string    `json:"text"`
	Embedding  []float64 `json:"embedding,omitempty"`
}

type QuestionKnowledgeHit struct {
	QuestionID string  `json:"questionId"`
	Title      string  `json:"title"`
	Category   string  `json:"category"`
	Text       string  `json:"text"`
	Score      float64 `json:"score"`
}

type QuestionIndexStats struct {
	TotalQuestions int `json:"totalQuestions"`
	IndexedCount   int `json:"indexedCount"`
}

type QuestionRetrievalService interface {
	EnsureCollection(ctx context.Context) error
	IndexQuestions(ctx context.Context, docs []QuestionVectorDocument) error
	IndexQuestionsFromRepository(ctx context.Context) (*QuestionIndexStats, error)
	SearchQuestions(ctx context.Context, query string, topK int) ([]QuestionKnowledgeHit, error)
}

type questionRetrievalService struct {
	embeddingService EmbeddingService
	vectorStore      vectorstore.MilvusStore
	questionRepo     repository.QuestionRepository
}

func NewQuestionRetrievalService(
	embeddingService EmbeddingService,
	vectorStore vectorstore.MilvusStore,
	questionRepo repository.QuestionRepository,
) QuestionRetrievalService {
	return &questionRetrievalService{
		embeddingService: embeddingService,
		vectorStore:      vectorStore,
		questionRepo:     questionRepo,
	}
}

func (s *questionRetrievalService) EnsureCollection(ctx context.Context) error {
	return s.vectorStore.EnsureCollection(ctx)
}

func (s *questionRetrievalService) IndexQuestions(ctx context.Context, docs []QuestionVectorDocument) error {
	if len(docs) == 0 {
		return nil
	}

	inputs := make([]string, 0, len(docs))
	for _, doc := range docs {
		inputs = append(inputs, doc.Text)
	}

	embeddings, err := s.embeddingService.Embed(ctx, EmbeddingRequest{Input: inputs})
	if err != nil {
		return err
	}

	vectorDocs := make([]vectorstore.VectorDocument, 0, len(docs))
	for idx, doc := range docs {
		if idx >= len(embeddings.Vectors) {
			break
		}
		vectorDocs = append(vectorDocs, vectorstore.VectorDocument{
			ID:     doc.QuestionID,
			Vector: embeddings.Vectors[idx].Embedding,
			Fields: map[string]interface{}{
				"question_id": doc.QuestionID,
				"title":       doc.Title,
				"category":    doc.Category,
				"text":        doc.Text,
			},
		})
	}

	return s.vectorStore.Upsert(ctx, vectorDocs)
}

func (s *questionRetrievalService) IndexQuestionsFromRepository(ctx context.Context) (*QuestionIndexStats, error) {
	items, err := s.questionRepo.ListAllForIndexing()
	if err != nil {
		return nil, err
	}

	docs := make([]QuestionVectorDocument, 0, len(items))
	for _, item := range items {
		docs = append(docs, BuildQuestionVectorDocument(
			strconv.FormatUint(item.ID, 10),
			item.Title,
			item.Category,
			item.Content,
			item.StandardAnswer,
			item.Tags,
		))
	}

	if err := s.EnsureCollection(ctx); err != nil {
		return nil, err
	}
	if err := s.IndexQuestions(ctx, docs); err != nil {
		return nil, err
	}

	return &QuestionIndexStats{
		TotalQuestions: len(items),
		IndexedCount:   len(docs),
	}, nil
}

func (s *questionRetrievalService) SearchQuestions(ctx context.Context, query string, topK int) ([]QuestionKnowledgeHit, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []QuestionKnowledgeHit{}, nil
	}

	embeddings, err := s.embeddingService.Embed(ctx, EmbeddingRequest{Input: []string{query}})
	if err != nil {
		return nil, err
	}
	if len(embeddings.Vectors) == 0 {
		return []QuestionKnowledgeHit{}, nil
	}

	results, err := s.vectorStore.Search(ctx, embeddings.Vectors[0].Embedding, topK, []string{"question_id", "title", "category", "text"})
	if err != nil {
		return nil, err
	}

	hits := make([]QuestionKnowledgeHit, 0, len(results))
	for _, result := range results {
		hits = append(hits, QuestionKnowledgeHit{
			QuestionID: getStringField(result.Fields, "question_id", result.ID),
			Title:      getStringField(result.Fields, "title", ""),
			Category:   getStringField(result.Fields, "category", ""),
			Text:       getStringField(result.Fields, "text", ""),
			Score:      result.Score,
		})
	}

	return hits, nil
}

func BuildQuestionVectorDocument(questionID, title, category, content, standardAnswer, tags string) QuestionVectorDocument {
	combined := strings.TrimSpace(strings.Join([]string{
		title,
		"Category: " + strings.TrimSpace(category),
		"Content: " + strings.TrimSpace(content),
		"Standard Answer: " + strings.TrimSpace(standardAnswer),
		"Tags: " + strings.TrimSpace(tags),
	}, "\n"))

	return QuestionVectorDocument{
		QuestionID: questionID,
		Title:      title,
		Category:   category,
		Text:       combined,
	}
}

func getStringField(fields map[string]interface{}, key, fallback string) string {
	if value, exists := fields[key]; exists {
		return fmt.Sprintf("%v", value)
	}
	return fallback
}
