package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tsaqiffatih/minddrift-server/config"
	"github.com/tsaqiffatih/minddrift-server/internal/handler"
	"github.com/tsaqiffatih/minddrift-server/internal/model"
	"github.com/tsaqiffatih/minddrift-server/internal/repository"
	"github.com/tsaqiffatih/minddrift-server/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	db := config.InitDB(cfg)

	err := db.AutoMigrate(
		&model.User{},
		&model.Article{},
		&model.Category{},
		&model.Tag{},
		&model.ArticleVersion{},
		&model.Comment{},
		&model.Image{},
		&model.SEOMetadata{},
		&model.Analytic{},
		&model.Backup{},
	)
	if err != nil {
		log.Fatalf("❌ Database migration failed: %v", err)
	}

	userRepo := repository.NewUserRepository(db)

	userService := service.NewUserService(userRepo)

	userHandler := handler.NewUserHandler(userService)

	fmt.Println("✅ Database migration completed!")

	r := gin.Default()

	handlers := map[string]interface{}{
		"user": userHandler,
	}

	RegisterRoutes(r, handlers)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Server running at http://localhost:%s\n", port)
	r.Run(":" + port)
}
