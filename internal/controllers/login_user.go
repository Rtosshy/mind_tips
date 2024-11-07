package controllers

import (
	"database/sql"
	"mind_tips/internal/auth"
	"mind_tips/internal/models"
	"mind_tips/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func LoginUser(db *sql.DB, c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.LogError(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}

	// ユーザーの情報をデータベースから取得
	var storedUser models.User
	err := db.QueryRow("SELECT id, name, password FROM users WHERE name = ?", user.Name).Scan(&storedUser.ID, &storedUser.Name, &storedUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.LogError(c, http.StatusUnauthorized, err, "Invalid name or password")
		} else {
			utils.LogError(c, http.StatusInternalServerError, err, "Failed to retrieve user")
		}
		return
	}

	// パスワードの比較
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		utils.LogError(c, http.StatusUnauthorized, err, "Invalid name or password")
		return
	}

	// JWTトークンを生成
	token, err := auth.GenerateJWT(storedUser.Name)
	if err != nil {
		utils.LogError(c, http.StatusInternalServerError, err, "Failed to generate token")
		return
	}

	// 成功時のレスポンス
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login successful",
		"data":    gin.H{"token": token},
	})
}
