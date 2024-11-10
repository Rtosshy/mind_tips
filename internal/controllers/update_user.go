package controllers

import (
	"database/sql"
	"mind_tips/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateUser(db *sql.DB, c *gin.Context) {
	userName, exists := c.Get("user_name")
	if !exists {
		utils.LogError(c, http.StatusUnauthorized, nil, "User not authorized")
		return
	}

	result, err := db.Exec(
		"UPDATE user SET user_name = ?, bio = ? WHERE user_name = ?",
		c.PostForm("new_user_name"),
		c.PostForm("new_bio"),
		userName,
	)
	if err != nil {
		utils.LogError(c, http.StatusInternalServerError, err, "Failed to update user")
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
		"message": "Update successful",
		"status":  "success",
	})
}
