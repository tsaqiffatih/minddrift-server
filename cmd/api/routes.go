package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tsaqiffatih/minddrift-server/internal/handler"
)

func RegisterRoutes(r *gin.Engine, handlers map[string]interface{}) {
	api := r.Group("/api")

	userHandler := handlers["user"].(*handler.UserHandler)
	// articleHandler := handlers["article"].(*handler.ArticleHandler)
	// imageHandler := handlers["image"].(*handler.ImageHandler)

	// User Routes
	userRoutes := api.Group("/users")
	{
		userRoutes.POST("/register", userHandler.RegisterUser)
		userRoutes.POST("/login", userHandler.LoginUser)
		// userRoutes.GET("/:id", userHandler.GetUserByID)
		// userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
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
