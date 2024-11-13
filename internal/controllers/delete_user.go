package controllers

import (
	"database/sql"
	"log"
	"mind_tips/internal/utils"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func DeleteUser(db *sql.DB, c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userName, ok := claims["user_name"].(string)
	if !ok {
		log.Printf("Failed to extract username from claims: %v", claims)
		utils.LogError(c, http.StatusUnauthorized, nil, "User not authorized")
		return
	}
	result, err := db.Exec("DELETE FROM user WHERE user_name = ?", userName)
	if err != nil {
		utils.LogError(c, http.StatusInternalServerError, err, "Failed to delete user")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.LogError(c, http.StatusInternalServerError, err, "Failed to retrieve affected rows")
		return
	}

	if rowsAffected == 0 {
		utils.LogError(c, http.StatusNotFound, nil, "User not found")
		return
	}

	// 成功時のレスポンス
	c.JSON(http.StatusOK, gin.H{
		"message": "Delete successful",
		"status":  "success",
	})
}
