package main

import (
	"log"
	"offerpilot/backend/internal/config"
	"offerpilot/backend/internal/database"
	"offerpilot/backend/internal/handler"
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

	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpireHours)
	dashboardService := service.NewDashboardService(dashboardRepo)
	questionService := service.NewQuestionService(questionRepo)
	interviewService := service.NewInterviewService(interviewRepo)
	reviewService := service.NewReviewService(reviewRepo)
	projectService := service.NewProjectService(projectRepo)

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
