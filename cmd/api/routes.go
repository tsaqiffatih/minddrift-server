package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tsaqiffatih/minddrift-server/config"
	"github.com/tsaqiffatih/minddrift-server/internal/handler"
	"github.com/tsaqiffatih/minddrift-server/internal/middleware"
)

func RegisterRoutes(r *gin.Engine, handlers map[string]interface{}, cfg *config.Config) {
	api := r.Group("/api")

	authMiddleware := middleware.AuthMiddleware(cfg)

	userHandler := handlers["user"].(handler.UserHandler)
	// articleHandler := handlers["article"].(*handler.ArticleHandler)
	// imageHandler := handlers["image"].(*handler.ImageHandler)

	// User Routes
	userRoutes := api.Group("/users")
	{
		userRoutes.POST("/register", userHandler.RegisterUser)
		userRoutes.POST("/login", userHandler.LoginUser)
		userRoutes.POST("/auth/forgot-password", userHandler.RequestResetPassword)
		userRoutes.POST("/auth/validate-reset-token", userHandler.ValidateResetToken)
		userRoutes.POST("/auth/reset-password", userHandler.ResetPassword)

		userRoutes.GET("/verify-email", userHandler.VerifyEmail)
		userRoutes.GET("/resend-email", userHandler.ResendEmail)

		userRoutes.DELETE("/:id", authMiddleware, userHandler.DeleteUser)
		// userRoutes.GET("/:id", userHandler.GetUserByID)
		// userRoutes.PUT("/:id", userHandler.UpdateUser)
		// userRoutes.GET("/me", authMiddleware, userHandler.GetCurrentUser)
	}

	// // Article Routes
	// articleRoutes := api.Group("/articles")
	// {
	// 	articleRoutes.POST("/", articleHandler.CreateArticle)
	// 	articleRoutes.GET("/", articleHandler.GetAllArticles)
	// 	articleRoutes.GET("/:id", articleHandler.GetArticleByID)
	// 	articleRoutes.PUT("/:id", articleHandler.UpdateArticle)
	// 	articleRoutes.DELETE("/:id", articleHandler.DeleteArticle)
	// }

	// // Image Routes
	// imageRoutes := api.Group("/images")
	// {
	// 	imageRoutes.POST("/", imageHandler.UploadImage)
	// 	imageRoutes.GET("/:id", imageHandler.GetImageByID)
	// 	imageRoutes.DELETE("/:id", imageHandler.DeleteImage)
	// }
}
