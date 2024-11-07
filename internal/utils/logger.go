package utils

import (
	"log"

	"github.com/gin-gonic/gin"
)

// ログ関数
func LogError(c *gin.Context, statusCode int, err error, message string) {
	// errがnilでない場合
	if err != nil {
		log.Printf("Error: %v, StatusCode: %d, Message: %s", err, statusCode, message) // エラーがある場合にログ出力
		c.JSON(statusCode, gin.H{
			"status":  "error",
			"message": message,
			"error":   err.Error(),
		})
	} else {
		// errがnilの場合
		log.Printf("StatusCode: %d, Message: %s", statusCode, message) // エラーがない場合にもログを出力
		c.JSON(statusCode, gin.H{
			"status":  "error",
			"message": message,
		})
	}
}
