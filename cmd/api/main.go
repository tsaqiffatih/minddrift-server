package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tsaqiffatih/minddrift-server/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  Warning: No .env file found, using system environment variables.")
	}

	// Get DATABASE_URL
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatalf("❌ DATABASE_URL is not set")
	}

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("✅ Database connected!")

	// Auto Migrate the database
	err = db.AutoMigrate(
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

	fmt.Println("✅ Database migration completed!")

	// Setup Router (using Gin)
	r := gin.Default()

	// Simple Health Check Endpoint
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
