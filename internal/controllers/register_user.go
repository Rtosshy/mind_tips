package controllers

import (
	"database/sql"
	"mind_tips/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(db *sql.DB, c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// ユーザー名の重複チェック
	var existingUser string
	err := db.QueryRow("SELECT name FROM users WHERE name = ?", user.Name).Scan(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"status":  "error",
			"message": "name already exists",
		})
		return
	} else if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Database error",
			"error":   err.Error(),
		})
		return
	}

	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to hash password",
			"error":   err.Error(),
		})
		return
	}

	// ユーザーをデータベースに挿入
	_, err = db.Exec("INSERT INTO users (name, password) VALUES (?, ?)", user.Name, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create user",
			"error":   err.Error(),
		})
		return
	}

	// 成功時のレスポンス
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "User created successfully",
		"data":    gin.H{"name": user.Name},
	})
}
