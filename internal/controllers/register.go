package controllers

import (
	"database/sql"
	"mind_tips/internal/models"
	"mind_tips/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(db *sql.DB, c *gin.Context) {
	var user models.UserAuth
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.LogError(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}

	// ユーザー名の重複チェック
	var existingUser string
	err := db.QueryRow("SELECT user_name FROM user WHERE user_name = ?", user.UserName).Scan(&existingUser)
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
	_, err = db.Exec("INSERT INTO user (user_name, password) VALUES (?, ?)", user.UserName, hashedPassword)
	if err != nil {
		utils.LogError(c, http.StatusInternalServerError, err, "Failed to create user")
		return
	}

	// 成功時のレスポンス
	c.JSON(http.StatusCreated, gin.H{
		"data":    gin.H{"name": user.UserName},
		"message": "User created successfully",
		"status":  "success",
	})
}
