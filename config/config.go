package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DatabaseURL    string
	Port           string
	JWTSecret      string
	JWTAudience    string
	SMTPHost       string
	SMTPPort       int
	SMTPUser       string
	SMTPPass       string
	MindDriftEmail string
	BaseURL        string
	FrontendURL    string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Warning: No .env file found, using system environment variables.")
	}

	port, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))

	config := &Config{
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		Port:           getEnv("PORT", "8080"),
		JWTSecret:      getEnv("JWT_SECRET", ""),
		JWTAudience:    getEnv("JWT_AUDIENCE", ""),
		SMTPHost:       getEnv("SMTP_HOST", ""),
		SMTPPort:       port,
		SMTPUser:       getEnv("SMTP_USER", ""),
		SMTPPass:       getEnv("SMTP_PASSWORD", ""),
		MindDriftEmail: getEnv("MINDDRIFT_EMAIL", ""),
		BaseURL:        getEnv("BASE_URL", ""),
		FrontendURL:    getEnv("FRONTEND_URL", ""),
	}

	if config.DatabaseURL == "" {
		log.Fatal("❌ DATABASE_URL is not set")
	}

	if config.JWTSecret == "" {
		log.Fatal("❌ JWT_SECRET is not set")
	}

	if config.SMTPHost == "" || config.SMTPUser == "" || config.SMTPPass == "" || config.BaseURL == "" || config.FrontendURL == "" {
		log.Fatal("❌ SMTP configuration is incomplete")
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func InitDB(cfg *Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	fmt.Println("✅ Database connected!")
	return db
}
