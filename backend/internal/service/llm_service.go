package service

import (
	"bufio"
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

var ErrLLMNotConfigured = errors.New("llm service is not configured")

type LLMMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GenerateTextRequest struct {
	Model        string       `json:"model,omitempty"`
	SystemPrompt string       `json:"systemPrompt,omitempty"`
	Messages     []LLMMessage `json:"messages,omitempty"`
	Temperature  *float64     `json:"temperature,omitempty"`
	MaxTokens    *int         `json:"maxTokens,omitempty"`
}

type GenerateTextResponse struct {
	Model        string `json:"model"`
	Content      string `json:"content"`
	FinishReason string `json:"finishReason,omitempty"`
	RawResponse  []byte `json:"-"`
}

type StreamTextRequest struct {
	Model        string       `json:"model,omitempty"`
	SystemPrompt string       `json:"systemPrompt,omitempty"`
	Messages     []LLMMessage `json:"messages,omitempty"`
	Temperature  *float64     `json:"temperature,omitempty"`
	MaxTokens    *int         `json:"maxTokens,omitempty"`
}

type StreamTextChunk struct {
	Model        string `json:"model,omitempty"`
	Delta        string `json:"delta"`
	FinishReason string `json:"finishReason,omitempty"`
	Done         bool   `json:"done"`
}

type StreamTextHandler func(chunk StreamTextChunk) error

type LLMService interface {
	GenerateText(ctx context.Context, req GenerateTextRequest) (*GenerateTextResponse, error)
	StreamText(ctx context.Context, req StreamTextRequest, handler StreamTextHandler) error
}

type llmService struct {
	baseURL    string
	apiKey     string
	model      string
	httpClient *http.Client
}

type openAIChatRequest struct {
	Model       string       `json:"model"`
	Messages    []LLMMessage `json:"messages"`
	Temperature *float64     `json:"temperature,omitempty"`
	MaxTokens   *int         `json:"max_tokens,omitempty"`
	Stream      bool         `json:"stream"`
}

type openAIChatResponse struct {
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

type openAIStreamResponse struct {
	Model   string `json:"model"`
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

func NewLLMService(baseURL, apiKey, model string) LLMService {
	return &llmService{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		model:   model,
		httpClient: &http.Client{
			Timeout: 90 * time.Second,
		},
	}
}

func (s *llmService) GenerateText(ctx context.Context, req GenerateTextRequest) (*GenerateTextResponse, error) {
	if err := s.validateConfig(); err != nil {
		return nil, err
	}

	payload, err := s.buildChatRequest(req.Model, req.SystemPrompt, req.Messages, req.Temperature, req.MaxTokens, false)
	if err != nil {
		return nil, err
	}

	body, err := s.doChatRequest(ctx, payload)
	if err != nil {
		return nil, err
	}

	var response openAIChatResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	if len(response.Choices) == 0 {
		return nil, errors.New("llm response contains no choices")
	}

	return &GenerateTextResponse{
		Model:        response.Model,
		Content:      response.Choices[0].Message.Content,
		FinishReason: response.Choices[0].FinishReason,
		RawResponse:  body,
	}, nil
}

func (s *llmService) StreamText(ctx context.Context, req StreamTextRequest, handler StreamTextHandler) error {
	if err := s.validateConfig(); err != nil {
		return err
	}
	if handler == nil {
		return errors.New("stream handler is required")
	}

	payload, err := s.buildChatRequest(req.Model, req.SystemPrompt, req.Messages, req.Temperature, req.MaxTokens, true)
	if err != nil {
		return err
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("llm stream request failed: %s", strings.TrimSpace(string(body)))
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" {
			return handler(StreamTextChunk{Done: true})
		}

		var payload openAIStreamResponse
		if err := json.Unmarshal([]byte(data), &payload); err != nil {
			return err
		}
		if len(payload.Choices) == 0 {
			continue
		}

		chunk := StreamTextChunk{
			Model:        payload.Model,
			Delta:        payload.Choices[0].Delta.Content,
			FinishReason: payload.Choices[0].FinishReason,
			Done:         payload.Choices[0].FinishReason != "",
		}
		if err := handler(chunk); err != nil {
			return err
		}
	}

	return scanner.Err()
}

func (s *llmService) validateConfig() error {
	if s.baseURL == "" || s.apiKey == "" || s.model == "" {
		return ErrLLMNotConfigured
	}
	return nil
}

func (s *llmService) buildChatRequest(model, systemPrompt string, messages []LLMMessage, temperature *float64, maxTokens *int, stream bool) (*openAIChatRequest, error) {
	resolvedModel := strings.TrimSpace(model)
	if resolvedModel == "" {
		resolvedModel = s.model
	}
	if resolvedModel == "" {
		return nil, ErrLLMNotConfigured
	}

	resolvedMessages := make([]LLMMessage, 0, len(messages)+1)
	if strings.TrimSpace(systemPrompt) != "" {
		resolvedMessages = append(resolvedMessages, LLMMessage{Role: "system", Content: systemPrompt})
	}
	resolvedMessages = append(resolvedMessages, messages...)
	if len(resolvedMessages) == 0 {
		return nil, errors.New("at least one message is required")
	}

	return &openAIChatRequest{
		Model:       resolvedModel,
		Messages:    resolvedMessages,
		Temperature: temperature,
		MaxTokens:   maxTokens,
		Stream:      stream,
	}, nil
}

func (s *llmService) doChatRequest(ctx context.Context, payload *openAIChatRequest) ([]byte, error) {
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/chat/completions", bytes.NewReader(reqBody))
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
		return nil, fmt.Errorf("llm request failed: %s", strings.TrimSpace(string(body)))
	}

	return body, nil
}
