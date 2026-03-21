package vectorstore

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

var ErrMilvusNotConfigured = errors.New("milvus is not configured")

type VectorDocument struct {
	ID     string                 `json:"id"`
	Vector []float64              `json:"vector"`
	Fields map[string]interface{} `json:"fields"`
}

type SearchResult struct {
	ID     string                 `json:"id"`
	Score  float64                `json:"score"`
	Fields map[string]interface{} `json:"fields"`
}

type MilvusStore interface {
	EnsureCollection(ctx context.Context) error
	Upsert(ctx context.Context, docs []VectorDocument) error
	Search(ctx context.Context, vector []float64, topK int, outputFields []string) ([]SearchResult, error)
}

type MilvusConfig struct {
	BaseURL      string
	Token        string
	Database     string
	Collection   string
	VectorDim    int
	VectorField  string
	PrimaryField string
	TextField    string
}

type milvusStore struct {
	cfg        MilvusConfig
	httpClient *http.Client
}

func NewMilvusStore(cfg MilvusConfig) MilvusStore {
	if cfg.VectorField == "" {
		cfg.VectorField = "embedding"
	}
	if cfg.PrimaryField == "" {
		cfg.PrimaryField = "id"
	}
	if cfg.TextField == "" {
		cfg.TextField = "text"
	}

	return &milvusStore{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (m *milvusStore) EnsureCollection(ctx context.Context) error {
	if err := m.validate(); err != nil {
		return err
	}

	body := map[string]interface{}{
		"collectionName": m.cfg.Collection,
		"dbName":         m.cfg.Database,
		"schema": map[string]interface{}{
			"autoID":             false,
			"enableDynamicField": true,
			"fields": []map[string]interface{}{
				{
					"fieldName": m.cfg.PrimaryField,
					"dataType":  "VarChar",
					"isPrimary": true,
					"maxLength": 128,
				},
				{
					"fieldName": "question_id",
					"dataType":  "VarChar",
					"maxLength": 128,
				},
				{
					"fieldName": "title",
					"dataType":  "VarChar",
					"maxLength": 512,
				},
				{
					"fieldName": "category",
					"dataType":  "VarChar",
					"maxLength": 128,
				},
				{
					"fieldName": m.cfg.TextField,
					"dataType":  "VarChar",
					"maxLength": 8192,
				},
				{
					"fieldName": m.cfg.VectorField,
					"dataType":  "FloatVector",
					"dimension": m.cfg.VectorDim,
				},
			},
		},
	}

	_, err := m.doRequest(ctx, http.MethodPost, "/v2/vectordb/collections/create", body, true)
	return err
}

func (m *milvusStore) Upsert(ctx context.Context, docs []VectorDocument) error {
	if err := m.validate(); err != nil {
		return err
	}
	if len(docs) == 0 {
		return nil
	}

	rows := make([]map[string]interface{}, 0, len(docs))
	for _, doc := range docs {
		row := map[string]interface{}{
			m.cfg.PrimaryField: doc.ID,
			m.cfg.VectorField:  doc.Vector,
		}
		for key, value := range doc.Fields {
			row[key] = value
		}
		rows = append(rows, row)
	}

	body := map[string]interface{}{
		"collectionName": m.cfg.Collection,
		"dbName":         m.cfg.Database,
		"data":           rows,
	}

	_, err := m.doRequest(ctx, http.MethodPost, "/v2/vectordb/entities/upsert", body, false)
	return err
}

func (m *milvusStore) Search(ctx context.Context, vector []float64, topK int, outputFields []string) ([]SearchResult, error) {
	if err := m.validate(); err != nil {
		return nil, err
	}
	if topK <= 0 {
		topK = 5
	}

	body := map[string]interface{}{
		"collectionName": m.cfg.Collection,
		"dbName":         m.cfg.Database,
		"data":           []interface{}{vector},
		"annsField":      m.cfg.VectorField,
		"limit":          topK,
		"outputFields":   outputFields,
	}

	responseBody, err := m.doRequest(ctx, http.MethodPost, "/v2/vectordb/entities/search", body, false)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, err
	}

	results := make([]SearchResult, 0, len(response.Data))
	for _, item := range response.Data {
		result := SearchResult{Fields: map[string]interface{}{}}
		if id, ok := item[m.cfg.PrimaryField]; ok {
			result.ID = fmt.Sprintf("%v", id)
		}
		if distance, ok := item["distance"].(float64); ok {
			result.Score = distance
		} else if score, ok := item["score"].(float64); ok {
			result.Score = score
		}
		for key, value := range item {
			if key == m.cfg.PrimaryField || key == "distance" || key == "score" {
				continue
			}
			result.Fields[key] = value
		}
		results = append(results, result)
	}

	return results, nil
}

func (m *milvusStore) validate() error {
	if m.cfg.BaseURL == "" || m.cfg.Collection == "" || m.cfg.VectorDim <= 0 {
		return ErrMilvusNotConfigured
	}
	return nil
}

func (m *milvusStore) doRequest(ctx context.Context, method, path string, body interface{}, ignoreExists bool) ([]byte, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	url := strings.TrimRight(m.cfg.BaseURL, "/") + path
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if m.cfg.Token != "" {
		req.Header.Set("Authorization", "Bearer "+m.cfg.Token)
	}

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		if ignoreExists && strings.Contains(strings.ToLower(string(respBody)), "already exists") {
			return respBody, nil
		}
		return nil, fmt.Errorf("milvus request failed: %s", strings.TrimSpace(string(respBody)))
	}

	return respBody, nil
}
