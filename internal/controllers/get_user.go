package controllers

import (
	"database/sql"
	"mind_tips/internal/models"
	"mind_tips/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUser(db *sql.DB, c *gin.Context) {
	userName, exists := c.Get("user_name")
	if !exists {
		utils.LogError(c, http.StatusUnauthorized, nil, "User not authorized")
		return
	}

	var user models.UserResponse
	err := db.QueryRow("SELECT user_name, bio FROM user WHERE user_name = ?", userName).Scan(&user.UserName, &user.Bio)
	if err != nil {
		utils.LogError(c, http.StatusInternalServerError, err, "Failed to retrieve user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user, "message": "User found", "status": "success"})
}
