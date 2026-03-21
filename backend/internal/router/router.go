package router

import (
	"offerpilot/backend/internal/handler"
	"offerpilot/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func New(authHandler *handler.AuthHandler, dashboardHandler *handler.DashboardHandler, questionHandler *handler.QuestionHandler, interviewHandler *handler.InterviewHandler, reviewHandler *handler.ReviewHandler, projectHandler *handler.ProjectHandler, jwtSecret string) *gin.Engine {
	engine := gin.Default()

	api := engine.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.Use(middleware.JWTAuth(jwtSecret))
			auth.GET("/me", authHandler.Me)
			auth.POST("/logout", authHandler.Logout)
		}

		dashboard := api.Group("/dashboard")
		dashboard.Use(middleware.JWTAuth(jwtSecret))
		{
			dashboard.GET("/summary", dashboardHandler.GetSummary)
			dashboard.GET("/recent-sessions", dashboardHandler.GetRecentSessions)
			dashboard.GET("/weak-points", dashboardHandler.GetWeakPoints)
		}

		questions := api.Group("/questions")
		questions.Use(middleware.JWTAuth(jwtSecret))
		{
			questions.GET("", questionHandler.List)
			questions.GET("/:id", questionHandler.Detail)
			questions.POST("/:id/favorite", questionHandler.AddFavorite)
			questions.DELETE("/:id/favorite", questionHandler.RemoveFavorite)
		}

		interview := api.Group("/interview")
		interview.Use(middleware.JWTAuth(jwtSecret))
		{
			interview.POST("/session", interviewHandler.CreateSession)
			interview.GET("/session/:id", interviewHandler.GetSession)
			interview.GET("/session/:id/stream", interviewHandler.StreamSession)
			interview.POST("/session/:id/message", interviewHandler.SubmitMessage)
			interview.POST("/session/:id/finish", interviewHandler.FinishSession)
		}

		review := api.Group("/review")
		review.Use(middleware.JWTAuth(jwtSecret))
		{
			review.GET("/history", reviewHandler.GetHistory)
			review.GET("/:sessionId", reviewHandler.GetReport)
		}

		project := api.Group("/project")
		project.Use(middleware.JWTAuth(jwtSecret))
		{
			project.POST("/polish", projectHandler.Create)
			project.GET("/polish/history", projectHandler.GetHistory)
			project.GET("/polish/:id", projectHandler.GetDetail)
		}
	}

	return engine
}