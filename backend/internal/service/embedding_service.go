package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var ErrEmbeddingNotConfigured = errors.New("embedding service is not configured")

type EmbeddingRequest struct {
	Model string   `json:"model,omitempty"`
	Input []string `json:"input"`
}

type EmbeddingVector struct {
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
}

type EmbeddingResponse struct {
	Model       string            `json:"model"`
	Vectors     []EmbeddingVector `json:"vectors"`
	RawResponse []byte            `json:"-"`
}

type EmbeddingService interface {
	Embed(ctx context.Context, req EmbeddingRequest) (*EmbeddingResponse, error)
}

type embeddingService struct {
	baseURL    string
	apiKey     string
	model      string
	httpClient *http.Client
}

type openAIEmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type openAIEmbeddingResponse struct {
	Model string `json:"model"`
	Data  []struct {
		Index     int       `json:"index"`
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
}

func NewEmbeddingService(baseURL, apiKey, model string) EmbeddingService {
	return &embeddingService{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		model:   model,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (s *embeddingService) Embed(ctx context.Context, req EmbeddingRequest) (*EmbeddingResponse, error) {
	if s.baseURL == "" || s.apiKey == "" || s.model == "" {
		return nil, ErrEmbeddingNotConfigured
	}
	if len(req.Input) == 0 {
		return nil, errors.New("embedding input is required")
	}

	resolvedModel := strings.TrimSpace(req.Model)
	if resolvedModel == "" {
		resolvedModel = s.model
	}

	payload, err := json.Marshal(openAIEmbeddingRequest{
		Model: resolvedModel,
		Input: req.Input,
	})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/embeddings", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("embedding request failed: %s", strings.TrimSpace(string(body)))
	}

	var response openAIEmbeddingResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	vectors := make([]EmbeddingVector, 0, len(response.Data))
	for _, item := range response.Data {
		vectors = append(vectors, EmbeddingVector{
			Index:     item.Index,
			Embedding: item.Embedding,
		})
	}

	return &EmbeddingResponse{
		Model:       response.Model,
		Vectors:     vectors,
		RawResponse: body,
	}, nil
}
