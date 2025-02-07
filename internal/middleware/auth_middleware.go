package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tsaqiffatih/minddrift-server/pkg/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token diperlukan"})
			c.Abort()
			return
		}

		userID, err := utils.VerifyJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
