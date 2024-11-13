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
