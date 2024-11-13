package controllers

import (
	"database/sql"
	"log"
	"mind_tips/internal/models"
	"mind_tips/internal/utils"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func GetUser(db *sql.DB, c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userName, ok := claims["user_name"].(string)
	if !ok {
		log.Printf("Failed to extract username from claims: %v", claims)
		utils.LogError(c, http.StatusUnauthorized, nil, "User not authorized")
		return
	}

	var user models.UserResponse
	err := db.QueryRow("SELECT user_name, bio FROM user WHERE user_name = ?", userName).Scan(&user.UserName, &user.Bio)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.LogError(c, http.StatusNotFound, err, "User not found")
			return
		}
		utils.LogError(c, http.StatusInternalServerError, err, "Failed to retrieve user")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    user,
		"message": "User found",
		"status":  "success",
	})
}
