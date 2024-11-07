package middleware

import (
	"mind_tips/internal/auth"
	"mind_tips/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JWT認証ミドルウェア
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			utils.LogError(c, http.StatusUnauthorized, nil, "Missing token")
			c.Abort()
			return
		}

		claims, err := auth.ValidateJWT(tokenString)
		if err != nil {
			utils.LogError(c, http.StatusUnauthorized, err, "Invalid token")
			c.Abort()
			return
		}

		// ユーザー名をコンテキストに設定
		c.Set("name", claims.Name)
		c.Next()
	}
}
