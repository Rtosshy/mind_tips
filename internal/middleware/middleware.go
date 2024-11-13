package middleware

import (
	"database/sql"
	"log"
	"mind_tips/internal/auth"

	jwt "github.com/appleboy/gin-jwt/v2"
)

func JWTMiddleware(db *sql.DB) *jwt.GinJWTMiddleware {
	authMiddleware, err := auth.NewJWTMiddleware(db)
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	return authMiddleware
}

// JWT認証ミドルウェア
// func JWTMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		tokenString := c.GetHeader("Authorization")
// 		if tokenString == "" {
// 			utils.LogError(c, http.StatusUnauthorized, nil, "Missing token")
// 			c.Abort()
// 			return
// 		}

// 		claims, err := auth.ValidateJWT(tokenString)
// 		if err != nil {
// 			utils.LogError(c, http.StatusUnauthorized, err, "Invalid token")
// 			c.Abort()
// 			return
// 		}

// 		// ユーザー名をコンテキストに設定
// 		c.Set("user_name", claims.UserName)
// 		c.Next()
// 	}
// }
