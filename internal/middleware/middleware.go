package middleware

import (
	"fmt"
	"mind_tips/internal/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JWT認証ミドルウェア
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		claims, err := auth.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// ユーザー名をコンテキストに設定
		fmt.Println("Username from token:", claims.Name)
		c.Set("name", claims.Name)
		c.Next()
	}
}
