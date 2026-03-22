package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	Port                   string
	MySQLDSN               string
	JWTSecret              string
	JWTExpireHours         int
	LLMBaseURL             string
	LLMAPIKey              string
	LLMModel               string
	EmbeddingBaseURL       string
	EmbeddingAPIKey        string
	EmbeddingModel         string
	MilvusBaseURL          string
	MilvusToken            string
	MilvusDatabase         string
	MilvusCollection       string
	MilvusVectorDim        int
	EnableQuestionIndexing bool
}

func Load() Config {
	loadDotEnv()

	return Config{
		Port:                   getEnv("APP_PORT", "8080"),
		MySQLDSN:               getEnv("MYSQL_DSN", "root:password@tcp(127.0.0.1:3306)/offerpilot?charset=utf8mb4&parseTime=True&loc=Local"),
		JWTSecret:              getEnv("APP_JWT_SECRET", "change_me"),
		JWTExpireHours:         getEnvAsInt("APP_JWT_EXPIRE_HOURS", 72),
		LLMBaseURL:             getEnv("LLM_BASE_URL", "https://api.openai.com/v1"),
		LLMAPIKey:              getEnv("LLM_API_KEY", ""),
		LLMModel:               getEnv("LLM_MODEL", "gpt-4o-mini"),
		EmbeddingBaseURL:       getEnv("EMBEDDING_BASE_URL", "https://api.openai.com/v1"),
		EmbeddingAPIKey:        getEnv("EMBEDDING_API_KEY", ""),
		EmbeddingModel:         getEnv("EMBEDDING_MODEL", "text-embedding-3-small"),
		MilvusBaseURL:          getEnv("MILVUS_BASE_URL", "http://127.0.0.1:19530"),
		MilvusToken:            getEnv("MILVUS_TOKEN", ""),
		MilvusDatabase:         getEnv("MILVUS_DATABASE", "default"),
		MilvusCollection:       getEnv("MILVUS_COLLECTION", "question_knowledge"),
		MilvusVectorDim:        getEnvAsInt("MILVUS_VECTOR_DIM", 1536),
		EnableQuestionIndexing: getEnvAsBool("ENABLE_QUESTION_INDEXING", false),
	}
}

func loadDotEnv() {
	candidates := []string{
		".env",
		filepath.Join("backend", ".env"),
	}

	for _, path := range candidates {
		file, err := os.Open(path)
		if err != nil {
			continue
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			value = strings.Trim(value, `"'`)

			if key == "" {
				continue
			}

			if _, exists := os.LookupEnv(key); !exists {
				_ = os.Setenv(key, value)
			}
		}

		return
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvAsBool(key string, fallback bool) bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if value == "" {
		return fallback
	}

	switch value {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}
