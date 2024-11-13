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

func UpdateUser(db *sql.DB, c *gin.Context, authMiddleware *jwt.GinJWTMiddleware) {
	claims := jwt.ExtractClaims(c)
	currentUserName, ok := claims["user_name"].(string)
	if !ok {
		log.Printf("Failed to extract username from claims: %v", claims)
		utils.LogError(c, http.StatusUnauthorized, nil, "User not authorized")
		return
	}

	var request models.UpdateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.LogError(c, http.StatusBadRequest, err, "Invalid request format")
		return
	}

	// 入力値の検証
	if request.NewUserName == "" || request.NewBio == "" {
		utils.LogError(c, http.StatusBadRequest, nil, "Username and bio cannot be empty")
		return
	}

	// トランザクションを開始
	tx, err := db.Begin()
	if err != nil {
		utils.LogError(c, http.StatusInternalServerError, err, "Failed to begin transaction")
		return
	}
	defer tx.Rollback() // トランザクションの自動ロールバック

	// 新しいユーザー名が既に存在するか確認（現在のユーザー以外で）
	var exists bool
	err = tx.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM user WHERE user_name = ? AND user_name != ?)",
		request.NewUserName,
		currentUserName,
	).Scan(&exists)
	if err != nil {
		utils.LogError(c, http.StatusInternalServerError, err, "Failed to check username availability")
		return
	}
	if exists {
		utils.LogError(c, http.StatusConflict, nil, "Username already exists")
		return
	}

	// 現在のユーザー情報を取得（更新前の確認）
	var currentUser models.UserResponse
	err = tx.QueryRow(
		"SELECT user_name, bio FROM user WHERE user_name = ?",
		currentUserName,
	).Scan(&currentUser.UserName, &currentUser.Bio)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.LogError(c, http.StatusNotFound, err, "User not found")
			return
		}
		utils.LogError(c, http.StatusInternalServerError, err, "Failed to retrieve current user info")
		return
	}

	// ユーザー情報を更新
	result, err := tx.Exec(
		"UPDATE user SET user_name = ?, bio = ? WHERE user_name = ?",
		request.NewUserName,
		request.NewBio,
		currentUserName,
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
		utils.LogError(c, http.StatusNotFound, nil, "User not found or no changes made")
		return
	}

	// トランザクションをコミット
	if err := tx.Commit(); err != nil {
		utils.LogError(c, http.StatusInternalServerError, err, "Failed to commit transaction")
		return
	}

	// 新しいトークンを生成するためのデータを準備
	newClaims := map[string]interface{}{
		"user_id":   claims["user_id"],
		"user_name": request.NewUserName,
	}

	newToken, expire, err := authMiddleware.TokenGenerator(newClaims)
	if err != nil {
		utils.LogError(c, http.StatusInternalServerError, err, "Failed to generate new token")
		return
	}

	// 成功時のレスポンス
	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "Update successful",
		"new_token": newToken,
		"expire_at": expire,
		"data": gin.H{
			"previous": gin.H{
				"user_name": currentUser.UserName,
				"bio":       currentUser.Bio,
			},
			"current": gin.H{
				"user_name": request.NewUserName,
				"bio":       request.NewBio,
			},
		},
	})
}
