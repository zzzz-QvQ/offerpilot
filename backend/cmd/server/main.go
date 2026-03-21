package main

import (
	"context"
	"log"
	"offerpilot/backend/internal/config"
	"offerpilot/backend/internal/database"
	"offerpilot/backend/internal/handler"
	"offerpilot/backend/internal/pkg/vectorstore"
	"offerpilot/backend/internal/repository"
	"offerpilot/backend/internal/router"
	"offerpilot/backend/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := database.NewMySQL(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)
	questionRepo := repository.NewQuestionRepository(db)
	interviewRepo := repository.NewInterviewRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	projectRepo := repository.NewProjectRepository(db)

	llmService := service.NewLLMService(cfg.LLMBaseURL, cfg.LLMAPIKey, cfg.LLMModel)
	embeddingService := service.NewEmbeddingService(cfg.EmbeddingBaseURL, cfg.EmbeddingAPIKey, cfg.EmbeddingModel)
	milvusStore := vectorstore.NewMilvusStore(vectorstore.MilvusConfig{
		BaseURL:      cfg.MilvusBaseURL,
		Token:        cfg.MilvusToken,
		Database:     cfg.MilvusDatabase,
		Collection:   cfg.MilvusCollection,
		VectorDim:    cfg.MilvusVectorDim,
		VectorField:  "embedding",
		PrimaryField: "id",
		TextField:    "text",
	})
	questionRetrievalService := service.NewQuestionRetrievalService(embeddingService, milvusStore, questionRepo)
	_ = embeddingService

	if cfg.EnableQuestionIndexing {
		stats, err := questionRetrievalService.IndexQuestionsFromRepository(context.Background())
		if err != nil {
			log.Fatalf("failed to index questions into milvus: %v", err)
		}
		log.Printf("question indexing completed: total=%d indexed=%d", stats.TotalQuestions, stats.IndexedCount)
		return
	}

	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpireHours)
	dashboardService := service.NewDashboardService(dashboardRepo)
	questionService := service.NewQuestionService(questionRepo)
	interviewService := service.NewInterviewService(interviewRepo, llmService, questionRetrievalService)
	reviewService := service.NewReviewService(reviewRepo, llmService, questionRetrievalService)
	projectService := service.NewProjectService(projectRepo, llmService)

	authHandler := handler.NewAuthHandler(authService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)
	questionHandler := handler.NewQuestionHandler(questionService)
	interviewHandler := handler.NewInterviewHandler(interviewService)
	reviewHandler := handler.NewReviewHandler(reviewService)
	projectHandler := handler.NewProjectHandler(projectService)

	engine := router.New(authHandler, dashboardHandler, questionHandler, interviewHandler, reviewHandler, projectHandler, cfg.JWTSecret)
	if err := engine.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
