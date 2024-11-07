package controllers

import (
	"database/sql"
	"mind_tips/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUser(db *sql.DB, c *gin.Context) {
	name, exists := c.Get("name")
	if !exists {
		utils.LogError(c, http.StatusUnauthorized, nil, "User not authorized")
		return
	}

	type User struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	var user User
	err := db.QueryRow("SELECT id, name FROM users WHERE name = ?", name).Scan(&user.ID, &user.Name)
	if err != nil {
		utils.LogError(c, http.StatusInternalServerError, err, "Failed to retrieve user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User found", "user": user, "status": "success"})
}
