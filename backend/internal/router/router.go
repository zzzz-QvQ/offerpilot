package router

import (
	"offerpilot/backend/internal/handler"
	"offerpilot/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func New(authHandler *handler.AuthHandler, dashboardHandler *handler.DashboardHandler, questionHandler *handler.QuestionHandler, jwtSecret string) *gin.Engine {
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
	}

	return engine
}
