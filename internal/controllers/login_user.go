package controllers

import (
	"database/sql"
	"mind_tips/internal/auth"
	"mind_tips/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func LoginUser(db *sql.DB, c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// ユーザーの情報をデータベースから取得
	var storedUser models.User
	err := db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", user.Username).Scan(&storedUser.ID, &storedUser.Username, &storedUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Invalid username or password",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Internal server error",
				"error":   err.Error(),
			})
		}
		return
	}

	// パスワードの比較
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Invalid username or password",
		})
		return
	}

	// JWTトークンを生成
	token, err := auth.GenerateJWT(storedUser.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to generate token",
			"error":   err.Error(),
		})
		return
	}

	// 成功時のレスポンス
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login successful",
		"data":    gin.H{"token": token},
	})
}
