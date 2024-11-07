package controllers

import (
	"database/sql"
	"mind_tips/internal/models"
	"mind_tips/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(db *sql.DB, c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.LogError(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}

	// ユーザー名の重複チェック
	var existingUser string
	err := db.QueryRow("SELECT name FROM users WHERE name = ?", user.Name).Scan(&existingUser)
	if err == nil {
		utils.LogError(c, http.StatusConflict, nil, "name already exists")
		return
	} else if err != sql.ErrNoRows {
		utils.LogError(c, http.StatusInternalServerError, err, "Failed to retrieve user")
		return
	}

	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.LogError(c, http.StatusInternalServerError, err, "Failed to hash password")
		return
	}

	// ユーザーをデータベースに挿入
	_, err = db.Exec("INSERT INTO users (name, password) VALUES (?, ?)", user.Name, hashedPassword)
	if err != nil {
		utils.LogError(c, http.StatusInternalServerError, err, "Failed to create user")
		return
	}

	// 成功時のレスポンス
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "User created successfully",
		"data":    gin.H{"name": user.Name},
	})
}
