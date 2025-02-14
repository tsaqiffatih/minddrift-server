package utils

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/tsaqiffatih/minddrift-server/config"
	"github.com/tsaqiffatih/minddrift-server/internal/model"
)

var (
	jwtSecret   = []byte(os.Getenv("JWT_SECRET"))
	jwtAudience = os.Getenv("JWT_AUDIENCE")
)

type Claims struct {
	UserID uuid.UUID      `json:"user_id"`
	Role   model.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(cfg *config.Config, userID uuid.UUID, role model.UserRole) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "minddrift",
			Subject:   userID.String(),
			Audience:  []string{cfg.JWTAudience},
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := []byte(cfg.JWTSecret)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Println("❌ Failed to generate JWT token:", err)
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(cfg *config.Config, tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

func GenerateTokenVerification(userID string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Println("❌ Failed to generate token:", err)
		return "", errors.New("Failed to generate token")
	}

	return tokenString, nil
}

func ParseTokenVerification(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		log.Println("❌ Failed to parse token:", err)
		return "", errors.New("Internal server error")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("Invalid token")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return "", errors.New("Token has expired")
		}
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("Invalid token payload")
	}

	return userID, nil
}
